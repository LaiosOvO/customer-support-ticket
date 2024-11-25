package ws

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/chat"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type VisitorService struct{}

func (v *VisitorService) NewVistorServer(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade: ", err)
		return
	}

	// 获取get参数 创建ws连接
	vistorInfo := chat.FindVisitorByVistorId(c.Query("visitor_id"))
	if vistorInfo.VisitorId == "" {
		c.JSON(200, gin.H{
			"code": 404,
			"msg":  "访客不存在",
		})
	}

	//user := &User{
	//	Conn:       conn,
	//	Name:       vistorInfo.Name,
	//	Avator:     vistorInfo.Avator,
	//	Id:         vistorInfo.VisitorId,
	//	To_id:      vistorInfo.ToId,
	//	UpdateTime: time.Now(),
	//}

	// 存入访客map
	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:         conn,
		Addr:         conn.RemoteAddr().String(),
		HearbeatTime: currentTime,
		LoginTime:    currentTime,
		DataQueue:    make(chan []byte, 50),
	}
	ClientMap[c.Query("visitor_id")] = node

	go chat.UpdateVisitorStatus(vistorInfo.VisitorId, 1)
	//AddVisitorToList(user)
	// 接受消息

	go visitorReceiveMsg(conn, c)
}

// 接受用户的信息 转发给客服
func visitorReceiveMsg(conn *websocket.Conn, c *gin.Context) {
	for {
		var receive []byte
		messageType, receive, err := conn.ReadMessage()
		if err != nil {
			// 删除连接
			for _, visitor := range ClientList {
				if visitor.Conn == conn {
					log.Println("删除用户", visitor.Id)
					delete(ClientList, visitor.Id)
					VisitorOffline(visitor.To_id, visitor.Id, visitor.Name)
				}
			}
			log.Println(err)
			return
		}

		// 将消息放到信道
		message <- &Message{
			conn:        conn,
			content:     receive,
			context:     c,
			is_kefu:     false,
			messageType: messageType,
		}
	}
}

func VisitorMessage(visitorId, content string, kefuInfo chat.User) {
	msg := TypeMessage{
		Type: "message",
		Data: ClientMessage{
			Name:    kefuInfo.Nickname,
			Avator:  kefuInfo.Avator,
			Id:      kefuInfo.Name,
			Time:    time.Now().Format("2006-01-02 12:13:12"),
			ToId:    visitorId,
			Content: content,
			IsKefu:  "no",
		},
	}

	str, _ := json.Marshal(msg)
	visitor, ok := ClientMap[visitorId]
	if !ok || visitor == nil || visitor.Conn == nil {
		return
	}
	visitor.Conn.WriteMessage(websocket.TextMessage, str)
}

func OneVisitorMessage(toId string, str []byte) {
	visitor, ok := ClientMap[toId]
	if ok {
		log.Println("OneKefuMessage lock")
		visitor.Mux.Lock()
		defer visitor.Mux.Unlock()
		error := visitor.Conn.WriteMessage(websocket.TextMessage, str)
		if error != nil {
			log.Println("匿名用户发送消息失败", error)
		}
	}
}

func VisitorOffline(kefuId string, visitorId string, visitorName string) {
	chat.UpdateVisitorStatus(visitorId, 0)
	userInfo := make(map[string]string)
	userInfo["uid"] = visitorId
	userInfo["name"] = visitorName

	msg := TypeMessage{
		Type: "userOffline",
		Data: userInfo,
	}
	str, _ := json.Marshal(msg)

	OneKefuMessage(kefuId, str)
}

func AddVisitorToList(user *User) {
	oldUser, ok := ClientList[user.Id]
	if oldUser != nil || ok {
		msg := TypeMessage{
			Type: "close",
			Data: user.Id,
		}
		str, _ := json.Marshal(msg)
		if err := oldUser.Conn.WriteMessage(websocket.TextMessage, str); err != nil {
			oldUser.Conn.Close()
			user.UpdateTime = oldUser.UpdateTime
			delete(ClientList, user.Id)
		}
	}

	ClientList[user.Id] = user
	//lastMessage := chat.FindLastMessageByVisitorId(user.Id)
	userInfo := make(map[string]string)
	userInfo["uid"] = user.Id
	userInfo["username"] = user.Name
	userInfo["avator"] = user.Avator
	userInfo["last_message"] = ""
	if userInfo["last_message"] == "" {
		userInfo["last_message"] = "新访客"
	}

	msg := TypeMessage{
		Type: "userOnline",
		Data: userInfo,
	}

	str, _ := json.Marshal(msg)

	OneKefuMessage(user.To_id, str)
}
