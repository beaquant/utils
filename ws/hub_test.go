package ws

// https://github.com/gorilla/websocket/tree/master/examples/chat

import (
	"log"
	"net/http"
	"testing"
)

func TestHub_Run(t *testing.T) {
	fn := func(message []byte, hub *Hub) error {
		log.Println("message:", string(message))
		hub.AddBroadcast(message)
		return nil
	}
	hub := NewHub(fn) //新建一个用户
	go hub.Run()      //开始获取用户中传送的数据

	http.HandleFunc("/ws", func(res http.ResponseWriter, r *http.Request) {
		ServeWs(hub, res, r)
	})
	hub.AddBroadcast([]byte(string("this is return message")))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}
}
