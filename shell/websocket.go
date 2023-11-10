package shell

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// 执行测试， 测试， 测试
func TestWebsocketClient() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	uri := url.URL{Scheme: "ws", Host: "127.0.0.1", Path: "/api/eml/ws", RawQuery: "ak=api_u98gygzzh40yfwuylgi6paan1qrnt1zt&zone=h3.ink"}
	logrus.Info("connecting to ", uri.String())

	ccc, _, err := websocket.DefaultDialer.Dial(uri.String(), nil)
	if err != nil {
		log.Panic("dial:", err)
	}
	defer ccc.Close()

	// ccc.SetPingHandler(func(appData string) error {
	// 	logrus.Info("ping:", appData)
	// 	return nil
	// })
	done := make(chan struct{})
	go func() {
		defer close(done) // 关闭管道
		for {
			msgtype, message, err := ccc.ReadMessage()
			if err != nil {
				logrus.Error("read:", err.Error())
				return
			}
			if msgtype == websocket.PingMessage {
				// 不会处理， 被上级拦截器处理了
				ccc.WriteMessage(websocket.PongMessage, []byte("pong"))
			}
			logrus.Infof("recv: %s", string(message))
		}
	}()
	for {
		select {
		case <-done: // 服务端关闭
			return
		case <-interrupt:
			logrus.Info("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := ccc.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logrus.Error("write close:", err.Error())
			}
			select {
			case <-done: // 服务端关闭
			case <-time.After(time.Second): // 超时
			}
			return
		}
	}
}
