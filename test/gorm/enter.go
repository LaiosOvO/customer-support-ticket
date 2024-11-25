package gorm

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type TestApi struct {
}

// customer test
func (e *TestApi) GetExaCustomerList() {
	fmt.Println("*********************")
	fmt.Println("测试gorm映射列表")
	fmt.Println("*********************")
	//db := global.GVA_DB.Model(&example.ExaCustomer{})
	//var CustomerList []example.ExaCustomer
	//err := db.Preload("SysUser").Find(&CustomerList).Error
	//if err != nil {
	//	panic(err)
	//}

	//for _, ite := range CustomerList {
	//	fmt.Println(ite.ID)
	//	fmt.Println(ite.SysUser)
	//}
	//
	//err = db.Find(&CustomerList).Error
	//for _, ite := range CustomerList {
	//	fmt.Println(ite.CustomerName)
	//	fmt.Println(ite.SysUser)
	//}

	res := []system.SysUser{}
	global.GVA_DB.Find(&res)
	for _, ite := range res {
		fmt.Println(ite)
	}
}
