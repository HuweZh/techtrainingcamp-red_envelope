package models

import (
	"encoding/json"

	"huhusw.com/red_envelope/commons"
)

type User struct {
	UserId     int `grom:"user_id"`
	MaxCount   int `grom:"max_count"`
	CurCount   int `grom:"cur_count"`
	CreateTime int `grom:"create_time"`
	Amount     int `grom:"amount"`
}

//默认操作的是users表
//改变结构体的默认表名称
func (User) TableName() string {
	return "user"
}

//根据用户id，获取用户信息
func GetUser(id int) User {
	user := User{}
	//从redis缓存汇总获取用户

	//从数据库中获取用户
	commons.GetDB().Where("user_id", id).First(&user)
	return user
}

//根据用户id获取红包列表
func GetEnvelopeList(uid int) []Envelope {
	envelopes := []Envelope{}
	//从redis缓存汇总获取用户

	//从数据库中获取用户
	commons.GetDB().Where("user_id", uid).Find(&envelopes)
	return envelopes
}

//更新当前的抢红包次数
func UpdateCurCount(id int, cur int) {
	//更新单列
	commons.GetDB().Model(&User{}).Where("user_id = ?", id).Update("cur_count", cur)
}

//更新当前的抢红包次数
func UpdateAmount(id int, amount int) {
	//更新单列
	commons.GetDB().Model(&User{}).Where("user_id = ?", id).Update("amount", amount)
}

//编码json，存入redis
func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}
