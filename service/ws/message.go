package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net/http"
)

type Node struct {
	Conn         *websocket.Conn
	Addr         string
	FirstTime    uint64
	HearbeatTime uint64
	LoginTime    uint64
	DataQueue    chan []byte
	GroupSets    set.Interface
}

//type Message struct {
//	Is_Kefu bool
//	Type    int
//}

func sendUserChat(c *gin.Context) {

}

func Chat(writer http.ResponseWriter, request *http.Request) {

}

func sendProc(node *Node) {

	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("send data >>>> data", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func handleAllQueueSend() {
	for {
		data := <-message
		var typeMsg TypeMessage
		json.Unmarshal(data.content, &typeMsg)
		conn := data.conn

		if typeMsg.Type == nil || typeMsg.Data == nil {
			continue
		}

		msgType := typeMsg.Type.(string)
		log.Println("客户端受到信息: ", typeMsg)

		switch msgType {
		case "ping":
			msg := TypeMessage{
				Type: "pong",
			}

			str, _ := json.Marshal(msg)
			data.Mux.Lock()
			defer data.Mux.Unlock()
			conn.WriteMessage(websocket.TextMessage, str)
		case "inputing":
			data1 := typeMsg.Data.(map[string]interface{})
			isKefu := data1["is_kefu"]
			//from := data1["from"].(string)

			to := data1["to"].(string)

			if isKefu != nil && isKefu == "yes" {
				//visitorId := data1["from_id"].(string)
				//kefuInfo := models.FindVisitorByVistorId(visitorId)
				//VisitorMessage(visitorId, data1["content"].(string), kefuInfo)
			} else {
				OneKefuMessage(to, data.content)

			}
		}
	}
}

func recevProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		dispatch(data)
	}
}

func kefuReceive() {

}

func userReceive() {

}

func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("dispatch data error", err)
	}

	if msg.is_kefu {
		kefu2user(&msg)
	} else {
		user2kefu(&msg)
	}
}

func kefu2user(msg *Message) {

}

func user2kefu(msg *Message) {

}
