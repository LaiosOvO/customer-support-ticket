package ws

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/chat"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type WSService struct{}

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

type Node struct {
	Conn         *websocket.Conn
	Addr         string
	FirstTime    uint64
	HearbeatTime uint64
	LoginTime    uint64
	DataQueue    chan []byte
	GroupSets    set.Interface
	Mux          sync.Mutex
}

type Message struct {
	conn        *websocket.Conn
	context     *gin.Context
	content     []byte
	messageType int
	is_kefu     bool
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

var ClientMap = make(map[string]*Node)
var KefuMap = make(map[string]*Node)

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
		visitors := chat.FindVisitorsOnline()

		for _, visitor := range visitors {
			if visitor.VisitorId == "" {
				continue
			}
			_, ok := ClientList[visitor.VisitorId]
			if !ok {
				chat.UpdateVisitorStatus(visitor.VisitorId, 0)
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
func (wsservice *WSService) HandleAllMessageDispatch() {
	log.Println("后台的消息接受初始化")
	for {
		data := <-message

		fmt.Println(data)
		fmt.Println(data.content)

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

			to := data1["to"].(string)
			//from := data1["from"].(string)
			// 分发消息
			if data.is_kefu {
				OneVisitorMessage(to, data.content)
			} else {
				OneKefuMessage(to, data.content)
			}

		}
	}
}

func (wsservice *WSService) SendMessageV2(c *gin.Context) {
	fromId := c.PostForm("from_id")
	toId := c.PostForm("to_id")
	content := c.PostForm("content")
	cType := c.PostForm("type")
	//is_kefu := c.PostForm("is_kefu")
	is_kefu, _ := strconv.ParseBool(c.PostForm("is_kefu"))

	if content == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "内容不能为空",
		})
		return
	}

	var kefuInfo chat.User
	var vistorInfo chat.Visitor

	if is_kefu || cType == "kefu" {
		kefuInfo = chat.FindUser(fromId)
		vistorInfo = chat.FindVisitorByVistorId(toId)
	} else {
		kefuInfo = chat.FindUser(toId)
		vistorInfo = chat.FindVisitorByVistorId(fromId)
	}

	if kefuInfo.ID == 0 || vistorInfo.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "用户不存在",
		})
		return
	}

	chat.CreateMessage(kefuInfo.Name, vistorInfo.VisitorId, content, cType)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
