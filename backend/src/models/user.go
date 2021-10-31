package models

type User struct {
	UserId     int `json:"uid"`
	MaxCount   int
	CurCount   int
	CreateTime int
}

//默认操作的是users表
//改变结构体的默认表名称
func (user User) TableName() string {
	return "user"
}
