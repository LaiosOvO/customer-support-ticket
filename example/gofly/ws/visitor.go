package ws

import (
	"github.com/flipped-aurora/gin-vue-admin/server/example/gofly/models"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func NewVistorServer(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade: ", err)
		return
	}

	// 获取get参数 创建ws连接
	vistorInfo := models.FindVisitorByVistorId(c.Query("visitor_id"))
	if vistorInfo.VisitorId == "" {
		c.JSON(200, gin.H{
			"code": 404,
			"msg":  "访客不存在",
		})
	}

	user := &User{
		Conn:       conn,
		Name:       vistorInfo.Name,
		Avator:     vistorInfo.Avator,
		Id:         vistorInfo.VisitorId,
		To_id:      vistorInfo.ToId,
		UpdateTime: time.Now(),
	}

	go models.UpdateVisitorStatus(vistorInfo.VisitorId, 1)

	AddVisitorToList(user)

	// 接受消息
	for {
		var receive []byte
		messageType, receive, err := conn.ReadMessage()
		if err != nil {
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

		message <- &Message{
			conn:        conn,
			content:     receive,
			context:     c,
			messageType: messageType,
		}
	}
}

func VisitorMessage(visitorId, content string, kefuInfo models.User) {

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
	visitor, ok := ClientList[visitorId]
	if !ok || visitor == nil || visitor.Conn == nil {
		return
	}
	visitor.Conn.WriteMessage(websocket.TextMessage, str)
}

func VisitorOffline(kefuId string, visitorId string, visitorName string) {

	models.UpdateVisitorStatus(visitorId, 0)
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
	//lastMessage := models.FindLastMessageByVisitorId(user.Id)
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
