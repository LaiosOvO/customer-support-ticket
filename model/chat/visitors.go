package chat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
	"time"
)

type Visitor struct {
	gorm.Model
	Name      string `json:"name"`
	Avator    string `json:"avator"`
	SourceIp  string `json:"source_ip"`
	ToId      string `json:"to_id"`
	VisitorId string `json:"visitor_id"`
	Status    uint   `json:"status"`
	Refer     string `json:"refer"`
	City      string `json:"city"`
	ClientIp  string `json:"client_ip"`
	Extra     string `json:"extra"`
}

func (Visitor) TableName() string {
	return "visitor"
}

func CreateVisitor(name, avator, sourceIp, toId, visitorId, refer, city, clientIp, extra string) {
	v := &Visitor{
		Name:      name,
		Avator:    avator,
		SourceIp:  sourceIp,
		ToId:      toId,
		VisitorId: visitorId,
		Status:    1,
		Refer:     refer,
		City:      city,
		ClientIp:  clientIp,
		Extra:     extra,
	}
	v.UpdatedAt = time.Now()
	global.GVA_DB.Create(v)
}

func FindVisitorByVistorId(visitorId string) Visitor {
	var v Visitor
	global.GVA_DB.Where("visitor_id = ?", visitorId).First(&v)
	return v
}

func FindVisitors(page uint, pagesize uint) []Visitor {
	offset := (page - 1) * pagesize
	if offset < 0 {
		offset = 0
	}
	var visitors []Visitor
	global.GVA_DB.Offset(int(offset)).Limit(int(pagesize)).Order("status desc, updated_at desc").Find(&visitors)
	return visitors
}

func FindVisitorsByKefuId(page uint, pagesize uint, kefuId string) []Visitor {
	offset := (page - 1) * pagesize
	if offset <= 0 {
		offset = 0
	}
	var visitors []Visitor
	global.GVA_DB.Where("to_id=?", kefuId).Offset(int(offset)).Limit(int(pagesize)).Order("updated_at desc").Find(&visitors)
	return visitors
}

func FindVisitorsOnline() []Visitor {
	var visitors []Visitor
	global.GVA_DB.Where("status = ?", 1).Find(&visitors)
	return visitors
}

func UpdateVisitorStatus(visitorId string, status uint) {
	visitor := Visitor{}
	global.GVA_DB.Model(&visitor).Where("visitor_id = ?", visitorId).Update("status", status)
}

func UpdateVisitor(name, avatar, visitorId string, status uint, clientIp string, sourceIp string, refer, extra string) {
	// 更新的字段和值
	updates := map[string]interface{}{
		"Status":    status,
		"ClientIp":  clientIp,
		"SourceIp":  sourceIp,
		"Refer":     refer,
		"Extra":     extra,
		"Name":      name,
		"Avatar":    avatar, // 注意这里变量名是 avatar 而不是 avator
		"UpdatedAt": time.Now(),
	}
	// 使用 Model 指定要更新的记录类型，Where 指定更新条件，Updates 指定要更新的字段和值
	global.GVA_DB.Model(&Visitor{}).Where("visitor_id = ?", visitorId).Updates(updates)
}

func UpdateVisitorKefu(visitorId string, kefuId string) {
	visitor := Visitor{}
	global.GVA_DB.Model(&visitor).Where("visitor_id = ?", visitorId).Update("to_id", kefuId)
}

// 查询条数
func CountVisitors() int64 {
	var count int64
	global.GVA_DB.Model(&Visitor{}).Count(&count)
	return count
}

// 查询条数
func CountVisitorsByKefuId(kefuId string) int64 {
	var count int64
	global.GVA_DB.Model(&Visitor{}).Where("to_id=?", kefuId).Count(&count)
	return count
}
