package ws

import (
	"github.com/flipped-aurora/gin-vue-admin/server/example/gofly/models"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type User struct {
	Conn       *websocket.Conn
	Name       string
	Id         string
	Avator     string
	To_id      string
	Role_id    string
	Mux        sync.Mutex
	UpdateTime time.Time
}

type Message struct {
	conn        *websocket.Conn
	context     *gin.Context
	content     []byte
	messageType int
	Mux         sync.Mutex
}

type TypeMessage struct {
	Type interface{} `json: "type"`
	Data interface{} `json: "data"`
}

type ClientMessage struct {
	Name      string `json: "name"`
	Avator    string `json: "avator"`
	Id        string `json: "id"`
	VisitorId string `json: "visitor_id"`
	Group     string `json: "group"`
	Time      string `json: "time"`
	ToId      string `json: "time"`
	Content   string `json: "content"`
	City      string `json: "city"`
	ClientIp  string `json: "client_ip"`
	Refer     string `json: "refer"`
	IsKefu    string `json: "is_kefu"`
}

var ClientList = make(map[string]*User)
var KefuList = make(map[string]*User)

var message = make(chan *Message, 10)
var Mux sync.RWMutex

var upgrader = websocket.Upgrader{}

func init() {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,

		CheckOrigin: func(r *http.Request) bool { return true },
	}

	//go UpdateVisitorStatusCron()
}

func UpdateVisitorStatusCron() {

	for {
		visitors := models.FindVisitorsOnline()

		for _, visitor := range visitors {
			if visitor.VisitorId == "" {
				continue
			}
			_, ok := ClientList[visitor.VisitorId]
			if !ok {
				models.UpdateVisitorStatus(visitor.VisitorId, 0)
			}

			SendPingToKefuClient()
			time.Sleep(120 * time.Second)
		}
	}
}

// 定期检查客服的在线情况
func SendPingToKefuClient() {
	msg := TypeMessage{
		Type: "many pong",
	}

	str, _ := json.Marshal(msg)
	for kefuId, kefu := range KefuList {
		if kefu == nil {
			continue
		}

		kefu.Mux.Lock()
		defer kefu.Mux.Unlock()
		err := kefu.Conn.WriteMessage(websocket.TextMessage, str)
		if err != nil {
			log.Println("定时发送ping给客服,失败", err.Error())
			delete(KefuList, kefuId)
		}
	}
}

// 后端官博发送消息
func WebServerBackend() {
	log.Println("后台的消息接受初始化")
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
