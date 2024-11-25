package service

import "github.com/flipped-aurora/gin-vue-admin/server/service/ws"

type ServiceGroup struct {
	//SystemServiceGroup  system.ServiceGroup
	ChatServiceGroup ws.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
