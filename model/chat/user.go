package chat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avator   string `json:"avator"`
	RoleName string `json:"role_name" sql:"-"`
	RoleId   string `json:"role_id" sql:"-"`
}

func (User) TableName() string {
	return "user"
}

func CreateUser(name string, password string, avator string, nickname string) uint {
	user := &User{
		Name:     name,
		Password: password,
		Avator:   avator,
		Nickname: nickname,
	}
	user.UpdatedAt = time.Now()
	global.GVA_DB.Create(user)
	return user.ID
}

func UpdateUser(id string, name string, password string, avator string, nickname string) {
	user := &User{
		Name:     name,
		Avator:   avator,
		Nickname: nickname,
	}
	user.UpdatedAt = time.Now()
	if password != "" {
		user.Password = password
	}
	//global.GVA_DB.Model(&User{}).Where("id = ?", id).Update(user)

}

func UpdateUserPass(name string, pass string) {
	user := &User{
		Password: pass,
	}
	user.UpdatedAt = time.Now()
	global.GVA_DB.Model(user).Where("name = ?", name).Update("Password", pass)
}

func UpdateUserAvator(name string, avator string) {
	user := &User{
		Avator: avator,
	}
	user.UpdatedAt = time.Now()
	global.GVA_DB.Model(user).Where("name = ?", name).Update("Avator", avator)
}
func FindUser(username string) User {
	var user User
	global.GVA_DB.Where("name = ?", username).First(&user)
	return user
}

func FindUserById(id interface{}) User {
	var user User
	global.GVA_DB.Select("user.*").Where("user.id = ?", id).First(&user)
	return user
}
