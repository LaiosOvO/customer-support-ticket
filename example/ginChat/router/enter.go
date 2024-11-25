package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/example/ginChat/service/message"
	"github.com/flipped-aurora/gin-vue-admin/server/example/ginChat/service/user"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	r := gin.Default()

	// 首页

	// 用户
	r.POST("/user/getUserList", user.GetUserList)
	r.POST("/user/createUser", user.CreateUser)
	r.POST("/user/updateUser", user.UpdateUser)
	r.POST("/user/deleteUser", user.DeleteUser)
	r.POST("/user/login", user.Login)

	//发送消息
	r.GET("/user/sendMsg", message.SendMsg)
	//发送消息
	//r.GET("/user/sendUserMsg", service.SendUserMsg)

	// 发送redis消息
	r.POST("/user/redisMsg", message.RedisMsg)
	return r
}
