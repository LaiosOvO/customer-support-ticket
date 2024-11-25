package models

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/example/ginChat/utils"
	"strconv"
	"sync"
)

// 读写锁
var rwLocker sync.RWMutex

// 获取缓存里面的消息
func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	rwLocker.RLock()
	//node, ok := clientMap[userIdA]
	rwLocker.RUnlock()
	//jsonMsg := Message{}
	//json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))
	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}
	//key = "msg_" + userIdStr + "_" + targetIdStr
	//rels, err := utils.Red.ZRevRange(ctx, key, 0, 10).Result()  //根据score倒叙

	var rels []string
	var err error
	if isRev {
		rels, err = utils.Red.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = utils.Red.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println(err) //没有找到
	}
	// 发送推送消息
	/**
	// 后台通过websoket 推送消息
	for _, val := range rels {
		fmt.Println("sendMsg >>> userID: ", userIdA, "  msg:", val)
		node.DataQueue <- []byte(val)
	}**/
	return rels
}
