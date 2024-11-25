package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	KefuId    string `json: "kefu_id"`
	VisitorId string `json: "visitor_id"`
	Content   string `json: "content"`
	MesType   string `json: "mes_type"`
	Status    string `json: "status"`
}
