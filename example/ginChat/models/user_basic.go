package models

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/example/ginChat/utils"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string    `json:"name" `
	PassWord      string    `json:"pass_word"`
	Phone         string    `json:"phone"  valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string    `valid:"email" `
	Avatar        string    `json:"avatar" `
	Identity      string    `json:"identity" `
	ClientIp      string    `json:"client_ip" `
	ClientPort    string    `json:"client_port"`
	Salt          string    `json:"salt"`
	LoginTime     time.Time `json:"login_time"`
	HeartbeatTime time.Time `json:"heartbeat_time"`
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	GetUserList2()
	var data []*UserBasic

	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func GetUserList2() []*UserBasic {
	data := []*UserBasic{}

	utils.DB.Table("user_basic").Find(&data)
	for _, v := range data {
		fmt.Println(*v)
		fmt.Println(v.Name)
	}
	return data
}

func FindUserByNameAnePwd(name string, passowrd string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and pass_word =?", name, passowrd).First(&user)

	//
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.Md5Encode(str)

	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	return user
}
func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)

	return user
}

func CreateUser(user UserBasic) (err error) {
	err = utils.DB.Create(&user).Error
	if err != nil {
		fmt.Println("新增数据失败", err)
	}
	return err
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email, Avatar: user.Avatar})
}
