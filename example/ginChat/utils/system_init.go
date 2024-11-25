package utils

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.AddConfigPath("example/ginChat/config")
	viper.SetConfigName("app")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("config init err")
	}
	fmt.Println("config app inited")
}

type UserBasic struct {
	gorm.Model
	Name          string    `gorm:"column:name"`
	PassWord      string    `gorm:"column:pass_word"`
	Phone         string    `gorm:"column:phone"`
	Email         string    `gorm:"column:email"`
	Identity      string    `gorm:"column:identity"`
	ClientIp      string    `gorm:"column:client_ip"`
	ClientPort    string    `gorm:"column:client_port"`
	LoginTime     time.Time `gorm:"column:login_time"`
	HeartbeatTime time.Time `gorm:"column:heartbeat_time"`
	LoginOutTime  time.Time `gorm:"column:login_out_time"`
	IsLogout      bool      `gorm:"column:is_logout"`
	DeviceInfo    string    `gorm:"column:device_info"`
	Salt          string    `gorm:"column:salt"`
	Avatar        string    `gorm:"column:avatar"`
}

func InitMySQL() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	var count int64
	if err := DB.Table("user_basic").Count(&count).Error; err != nil {
		log.Println("数据库查询失败：%v", err)
	}
	var res = []UserBasic{}
	DB.Table("user_basic").Find(&res)
	for _, ite := range res {
		fmt.Println(ite)
	}
	fmt.Println("count user talbe count = ", count)
	fmt.Println("MySql inited ...")
}

func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

}

const (
	PublishKey = "websocket"
)

// Publish 发布消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Publish 。。。。", msg)
	err = Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅Redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Subscribe 。。。。", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Subscribe 。。。。", msg.Payload)
	return msg.Payload, err
}
