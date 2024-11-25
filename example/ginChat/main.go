package main

import (
	"github.com/flipped-aurora/gin-vue-admin/server/example/ginChat/router"
	"github.com/flipped-aurora/gin-vue-admin/server/example/ginChat/utils"
	"github.com/spf13/viper"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()

	r.Run(viper.GetString("port.server"))
}
