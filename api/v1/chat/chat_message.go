package chat

import "github.com/gin-gonic/gin"

type ChatMessageApi struct{}

func (chatService *ChatMessageApi) CreateVisitorConnection(c *gin.Context) {
	visitorService.NewVistorServer(c)
}

func (chatService *ChatMessageApi) CreateKefuConnection(c *gin.Context) {
	kefuService.NewKefuServer(c)
}
