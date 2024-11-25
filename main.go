package main

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	_ "go.uber.org/automaxprocs"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title                       Gin-Vue-Admin Swagger API接口文档
// @version                     v2.7.3
// @description                 使用gin+vue进行极速开发的全栈开发基础平台
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /

func main() {
	global.GVA_VP = core.Viper() // 初始化Viper
	//fmt.Println(global.GVA_VP)
	fmt.Println(global.GVA_CONFIG)
	//initialize.OtherInit()
	//global.GVA_LOG = core.Zap() // 初始化zap日志库
	//zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	//initialize.Timer()
	initialize.DBList()
	//if global.GVA_DB != nil {
	//	initialize.RegisterTables() // 初始化表
	//	// 程序结束前关闭数据库链接
	//	db, _ := global.GVA_DB.DB()
	//	defer db.Close()
	//}
	//fmt.Println("invoke 的 db")
	//api := gorm.TestApi{}
	//api.GetExaCustomerList()

	//fmt.Println()
	go service.ServiceGroupApp.ChatServiceGroup.WSService.HandleAllMessageDispatch()

	//go ws.WebServerBackend()
	core.RunWindowsServer()
}
