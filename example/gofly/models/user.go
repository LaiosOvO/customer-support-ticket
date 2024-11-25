package models

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json: "name"`
	Password string `json: "password"`
	Nickname string `json: "nickname"`
	Avator   string `json: "avator"`
	RoleName string `json: "role_name" sql:"_"`
	RoleId   string `json: "role_id" sql:"_"`
}

func FindUserById(id interface{}) User {
	var user User
	global.GVA_DB.Select("*").Where("user.name = ?", id).First(&user)
	return user
}

func (User) TableName() string {
	return "user"
}
