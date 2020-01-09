package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod) //设置定时
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send: //这里send里面的值时run里面获取的，在这里才开始实际向前台传值
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) //设置写入的死亡时间，相当于http的超时
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{}) //如果取值出错，关闭连接，设置写入状态，和对应的数据
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage) //以io的形式写入数据，参数为数据类型
			if err != nil {
				return
			}
			w.Write(message) //写入数据，这个函数才是真正的想前台传递数据

			if err := w.Close(); err != nil { //关闭写入流
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) //心跳包，下面ping出错就会报错退出，断开这个连接
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	log.Println("start read message")
	defer func() {
		c.hub.unregister <- c //读取完毕后注销该client
		c.conn.Close()
		log.Println("websocket Close")
	}()
	c.conn.SetReadLimit(maxMessageSize)              //设置最大读取容量
	c.conn.SetReadDeadline(time.Now().Add(pongWait)) //设置读取死亡时间
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("websocket.IsUnexpectedCloseError:", err)
			}
			log.Println("websocket read message have err")
			break
		}
		if c.hub.OnMessage != nil {
			if err := c.hub.OnMessage(message, c.hub); err != nil {
				log.Println(err)
				break
			}
		}
	}
}
