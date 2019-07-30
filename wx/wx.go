package wx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type WxPush struct {
	url     string
	key     string
	msgQ    [][]string
	msgLock sync.RWMutex
}

const (
	FANGTANG_URL = "https://sc.ftqq.com/"
)

func NewWxPush(url, key string) *WxPush {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	push := new(WxPush)
	push.url = url
	push.key = key
	push.msgQ = make([][]string, 0)
	go push.sendWxWorker()
	return push
}

func (push *WxPush) sendWxQ() {
	push.msgLock.Lock()
	defer push.msgLock.Unlock()
	if len(push.msgQ) == 0 {
		return
	}

	msg := push.msgQ[0]
	params := url.Values{}
	params.Set("text", msg[0] /*+"_"+time.Now().Format("2006-01-02 15:04:05")*/)
	params.Set("desp", msg[1])

	var path = ""
	if !strings.Contains(push.url, "pushbear") {
		path = fmt.Sprintf(push.url+"%s.send?%s", push.key, params.Encode())
	} else {
		params.Set("sendkey", push.key)
		path = fmt.Sprintf(push.url+"sub?%s", params.Encode())
	}
	//fmt.Println("send wx:", url)
	retry := 10
	for {
		_, err := httpGet(http.DefaultClient, path)

		if err != nil {
			time.Sleep(1 * time.Second)
			if retry > 0 {
				retry--
			} else {
				fmt.Println("wx send fail...", err)
				break
			}
		} else {
			break
		}
	}
	push.msgQ = push.msgQ[1:]
}

func (push *WxPush) wxMsgQPush(subject, body string) {
	push.msgLock.Lock()
	defer push.msgLock.Unlock()
	push.msgQ = append(push.msgQ, []string{subject, body})
}

func (push *WxPush) sendWxWorker() {
	timer := time.NewTicker(time.Duration(1) * time.Second)
	for {
		select {
		case <-timer.C:
			go push.sendWxQ()
		}
	}
}

func (push *WxPush) sendWx(subject string, args ...interface{}) {
	body := fmt.Sprint(args)
	push.wxMsgQPush(subject, body)
}

func (push *WxPush) SendWxString(subject, body string) {
	push.wxMsgQPush(subject, body)
}

func httpGet(client *http.Client, reqUrl string) (map[string]interface{}, error) {
	respData, err := newHttpRequest(client, "GET", reqUrl, "", nil)
	if err != nil {
		//log.Println("NewHttpRequest Err", err)
		return nil, err
	}

	var bodyDataMap map[string]interface{}
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		return nil, err
	}
	return bodyDataMap, nil
}

func newHttpRequest(client *http.Client, reqType string, reqUrl string, postData string, requstHeaders map[string]string) ([]byte, error) {
	req, _ := http.NewRequest(reqType, reqUrl, strings.NewReader(postData))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")

	if requstHeaders != nil {
		for k, v := range requstHeaders {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode, string(bodyData)))
	}

	return bodyData, nil
}
