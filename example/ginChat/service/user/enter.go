package user

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/example/ginChat/models"
	"github.com/flipped-aurora/gin-vue-admin/server/example/ginChat/utils"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"time"
)

func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	c.JSON(200, gin.H{
		"code":    0,
		"message": "",
		"data":    data,
	})
}

func CreateUser(c *gin.Context) {
	user := models.UserBasic{}

	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("identity")
	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)

	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名不能为空",
			"data":    user,
		})
	}

	if data.Name != "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名已注册",
			"data":    user,
		})
	}

	if password != repassword {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "两次密码不一样",
			"data":    user,
		})
	}

	//
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	user.LoginTime = time.Now()
	user.LoginOutTime = time.Now()
	user.HeartbeatTime = time.Now()

	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "新增用户成功",
		"data":    user,
	})
}

func DeleteUser(c *gin.Context) {

	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "删除用户成功",
		"data":    user,
	})
}

func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Avatar = c.PostForm("icon")
	user.Email = c.PostForm("email")

	if false {

	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    200,
			"message": "修改用户成功",
			"data":    user,
		})
	}
}

func Login(c *gin.Context) {
	data := models.UserBasic{}

	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")

	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "该用户不存在",
			"data":    data,
		})
		return
	}

	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "密码不正确",
			"data":    data,
		})
		return
	}

	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAnePwd(name, pwd)

	c.JSON(200, gin.H{
		"code":    200, //  0成功   -1失败
		"message": "登陆成功",
		"data":    data,
	})
}
