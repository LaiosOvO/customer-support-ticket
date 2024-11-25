package models

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// node
type Node struct {
	Conn          *websocket.Conn
	Addr          string
	FirstTime     uint64
	HeartbeatTime uint64
	LoginTime     uint64
	DataQueue     chan []byte
	GroupSets     set.Interface
}

func (node *Node) Hearbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
	return
}

// 消息
type Message struct {
	gorm.Model
	UserId     int64
	TargetId   int64
	Type       int
	Media      int
	Content    string
	CreateTime uint64
	ReadTime   uint64
	Pic        string
	Url        string
	Desc       string
	Amount     int
}

// 映射关系 [用户id] -- socket连接对象
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)

	// 获取websocket连接
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(writer, request, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(),
		HeartbeatTime: currentTime,
		LoginTime:     currentTime,
		DataQueue:     make(chan []byte, 1024),
		GroupSets:     set.New(set.ThreadSafe),
	}

	rwLocker.Lock()
	clientMap[userId] = node
	defer rwLocker.Unlock()

	fmt.Println(conn)
	go sendProc(node)
	go recvProc(node)

}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws]sendProc >>>>> msg:", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)

			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := Message{}
		err = json.Unmarshal(data, msg)
		if err != nil {
			fmt.Println(err)
		}
		if msg.Type == 3 {
			currentTime := uint64(time.Now().Unix())
			node.Hearbeat(currentTime)
		} else {
			dispatch(data)
			//
			fmt.Println("[ws] receivedProc <<<<<<", string(data))
		}
	}
}

func dispatch(data []byte) {

	msg := Message{}
	msg.CreateTime = uint64(time.Now().Unix())
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch msg.Type {
	case 1: // 私信
		fmt.Println("dispatch data :", string(data))
		sendMsg(msg.TargetId, data)
	case 2:
		// sendGroupMsg(msg.TargetId,data)
	}

}

func sendMsg(userId int64, msg []byte) {

	rwLocker.RLock()
	node, ok := clientMap[userId]

	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg)

	//ctx := context.Background()
	//targetIdStr := strconv.Itoa(int(userId))
	//userIdStr := strconv.Itoa(int(jsonMsg.UserId))
	jsonMsg.CreateTime = uint64(time.Now().Unix())
	// 在线缓存中查询用户
	r := "userId1"
	if r != "" {
		if ok {
			// 发送信息
			fmt.Println("sendMsg >>> userId: ", userId, "  msg:", string(msg))
			node.DataQueue <- msg
		}
	}
	defer rwLocker.RUnlock()

	fmt.Println("ctx, targetIdStr, userIdStr")
	//fmt.Println(ctx, targetIdStr, userIdStr)
	//var key string
	//if userId > jsonMsg.UserId {
	//	key = ""
	//}

	// 将消息存到数据库中

}
