package controller

import (
	"fmt"

	"huhusw.com/red_envelope/commons"
)

func TestDB() {
	//存储数据库中的查询数据
	killList := []commons.Killed{}
	//查询数据库
	commons.GetDB().Find(&killList)
	//遍历输出
	for index, value := range killList {
		fmt.Printf("index = %v, value = %v\n", index, value)
	}
}
