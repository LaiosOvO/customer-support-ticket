package v1

import "github.com/flipped-aurora/gin-vue-admin/server/api/v1/chat"

type ApiGroup struct {
	ChatApiGroup chat.ApiGroup
	//ExampleApiGroup example.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
