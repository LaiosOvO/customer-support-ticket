package models

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
)

type Visitor struct {
	gorm.Model
	Name   string `json: "name"`
	Avator string `json: "avator"`

	SourceIp  string `json: "source_ip"`
	ToId      string `json: "to_ip"`
	VisitorId string `json: "visitor_id"`
	Status    uint   `json: "status"`
	Refer     string `json: "refer"`
	City      string `json: "client_ip"`
	Extra     string `json: "extra"`
}

func FindVisitorsOnline() []Visitor {
	var visitors []Visitor
	global.GVA_DB.Where("status = ?", 1).Find(&visitors)

	return visitors
}

func FindVisitorByVistorId(visitor_id string) Visitor {
	var visitor Visitor

	global.GVA_DB.Where("visitor_id = ?", visitor_id).Find(&visitor)
	return visitor
}

func UpdateVisitorStatus(visitorId string, status uint) {
	visitor := Visitor{}
	global.GVA_DB.Model(&visitor).Where("visitor_id = ?", visitorId).Update("status", status)
}
