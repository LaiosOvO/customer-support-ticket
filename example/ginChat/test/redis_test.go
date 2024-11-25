package test

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/example/ginChat/utils"
	"testing"
	"time"
)

var ctx context.Context

func init() {
	ctx = context.Background()
}

// TestPublish 测试发布消息到redis
func TestPublish(t *testing.T) {
	msg := "当前时间: " + time.Now().Format("15:04:05")
	err := utils.Publish(ctx, utils.PublishKey, msg)
	if err != nil {

		fmt.Println(err)
	}
}
