package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type CustomerRouter struct{}

func (e *CustomerRouter) InitCustomerRouter(Router *gin.RouterGroup) {
	customerRouter := Router.Group("customer")
	//.Use(middleware.OperationRecord())
	testHandler := func(c *gin.Context) {
		response.OkWithMessage("test gin integeration", c)
	}
	customerRouterWithoutRecord := Router.Group("customer")
	{
		//customerRouter.POST("customer", exaCustomerApi.CreateExaCustomer)   // 创建客户
		//customerRouter.PUT("customer", exaCustomerApi.UpdateExaCustomer)    // 更新客户
		//customerRouter.DELETE("customer", exaCustomerApi.DeleteExaCustomer) // 删除客户

		customerRouter.POST("customer", testHandler)   // 创建客户
		customerRouter.PUT("customer", testHandler)    // 更新客户
		customerRouter.DELETE("customer", testHandler) // 删除客户
	}
	{
		customerRouterWithoutRecord.GET("customer", testHandler)     // 获取单一客户信息
		customerRouterWithoutRecord.GET("customerList", testHandler) // 获取客户列表
	}
}
