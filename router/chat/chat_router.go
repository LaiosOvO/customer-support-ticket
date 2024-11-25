package chat

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type ChatRouter struct{}

var testHandler = func(c *gin.Context) {
	response.OkWithMessage("test gin integeration", c)
}

func (e *ChatRouter) InitChatRouter(Router *gin.RouterGroup) {
	chatRouter := Router.Group("chat")

	{
		chatRouter.GET("/ws/visitor", v1.ApiGroupApp.ChatApiGroup.ChatMessageApi.CreateVisitorConnection) // 匿名用户创建websocket连接
		chatRouter.GET("/ws/kefu", v1.ApiGroupApp.ChatApiGroup.ChatMessageApi.CreateKefuConnection)       // 匿名用户创建websocket连接
		chatRouter.PUT("customer", testHandler)                                                           // 更新客户
		chatRouter.DELETE("customer", testHandler)                                                        // 删除客户
	}

	chatRouterWithoutRecord := Router.Group("message")
	{
		chatRouterWithoutRecord.GET("customer", testHandler)     // 获取单一客户信息
		chatRouterWithoutRecord.GET("customerList", testHandler) // 获取客户列表
	}
}
