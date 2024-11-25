package chat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	gorm.Model
	KefuId    string `json:"kefu_id"`
	VisitorId string `json:"visitor_id"`
	Content   string `json:"content"`
	MesType   string `json:"mes_type"`
	Status    string `json:"status"`
}

func (Message) TableName() string {
	return "message"
}

func CreateMessage(kefu_id string, visitor_id string, content string, mes_type string) {
	global.GVA_DB.Exec("set names utf8mb4")
	v := &Message{
		KefuId:    kefu_id,
		VisitorId: visitor_id,
		Content:   content,
		MesType:   mes_type,
		Status:    "unread",
	}
	v.UpdatedAt = time.Now()
	global.GVA_DB.Create(v)
}
