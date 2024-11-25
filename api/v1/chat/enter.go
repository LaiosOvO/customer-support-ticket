package chat

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	ChatMessageApi
}

// 调用的service
var (
	visitorService = service.ServiceGroupApp.ChatServiceGroup.VisitorService
	kefuService    = service.ServiceGroupApp.ChatServiceGroup.KeFuService
)
