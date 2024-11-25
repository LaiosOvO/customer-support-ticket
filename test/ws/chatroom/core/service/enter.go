package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/test/ws/chatroom/core/models"
	"github.com/gin-gonic/gin"
)

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
