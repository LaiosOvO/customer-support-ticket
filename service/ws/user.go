package ws

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/chat"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type KeFuService struct{}

func (KeFuService *KeFuService) NewKefuServer(c *gin.Context) {

	kefuId := c.Query("kefu_id")
	kefuInfo := chat.FindUser(kefuId)

	if kefuInfo.ID == 0 {
		c.JSON(200, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade: ", err)
		return
	}
	// 构建node
	// 将登陆的客服放入map
	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:         conn,
		Addr:         conn.RemoteAddr().String(),
		HearbeatTime: currentTime,
		LoginTime:    currentTime,
		DataQueue:    make(chan []byte, 50),
	}
	KefuMap[kefuInfo.Name] = node

	// 获取Get参数
	//var kefu User
	//kefu.Id = kefuInfo.Name
	//kefu.Name = kefuInfo.Nickname
	//kefu.Avator = kefuInfo.Avator
	//kefu.Role_id = kefuInfo.RoleId
	//kefu.Conn = conn
	//AddKeufuToList(&kefu)

	//receiveMsg(conn, c)
	go kefuReceiveMsg(conn, c)
}

// kefu下发消息到用户
func kefuReceiveMsg(conn *websocket.Conn, c *gin.Context) {
	for {
		var receive []byte
		messageType, receive, err := conn.ReadMessage()
		if err != nil {
			log.Println("ws/user.go", err)
			conn.Close()
			return
		}

		message <- &Message{
			conn:        conn,
			content:     receive,
			context:     c,
			messageType: messageType,
			is_kefu:     true,
		}
	}
}

//func AddKeufuToList(kefu *User) {
//	oldUser, ok := KefuList[kefu.Id]
//	if oldUser != nil || ok {
//		msg := TypeMessage{
//			Type: "close",
//			Data: kefu.Id,
//		}
//
//		str, _ := json.Marshal(msg)
//		if err := oldUser.Conn.WriteMessage(websocket.TextMessage, str); err != nil {
//			oldUser.Conn.Close()
//		}
//	}
//	KefuList[kefu.Id] = kefu
//}

// 给客服【toid】发送消息
func OneKefuMessage(toId string, str []byte) {
	kefu, ok := KefuMap[toId]
	if ok {
		log.Println("OneKefuMessage lock")
		kefu.Mux.Lock()
		defer kefu.Mux.Unlock()
		error := kefu.Conn.WriteMessage(websocket.TextMessage, str)
		if error != nil {
			log.Println("发送客服信息出错", error)
		}
	}
}

func KefuMessage(visitorId, content string, kefuInfo chat.User) {
	msg := TypeMessage{
		Type: "message",
		Data: ClientMessage{
			Name:    kefuInfo.Nickname,
			Avator:  kefuInfo.Avator,
			Id:      visitorId,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			ToId:    visitorId,
			Content: content,
			IsKefu:  "yes",
		},
	}

	str, _ := json.Marshal(msg)
	OneKefuMessage(kefuInfo.Name, str)
}
