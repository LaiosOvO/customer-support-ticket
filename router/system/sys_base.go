package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	testHandler := func(c *gin.Context) {
		response.OkWithMessage("test gin integeration", c)
	}

	baseRouter := Router.Group("base")
	{
		//baseRouter.POST("login", baseApi.Login)
		//baseRouter.POST("captcha", baseApi.Captcha)
		baseRouter.POST("login", testHandler)
		baseRouter.POST("captcha", testHandler)

	}
	return baseRouter
}
