package models

import "huhusw.com/red_envelope/commons"

type Envelope struct {
	EnvelopeId int `grom:"envelope_id" json:"envelope_id"`
	UserId     int `grom:"user_id" json:"user_id"`
	Value      int `grom:"value" json:"value"`
	Opened     int `grom:"opened" json:"opened"`
	SnatchTime int `grom:"snatch_time" json:"snatch_time"`
}

//默认操作的是envelopes表
//改变结构体的默认表名称
func (Envelope) TableName() string {
	return "envelope"
}

func GetEnvelope(id int) *Envelope {
	envelope := Envelope{}
	//从redis缓存汇总获取红包

	//从数据库中获取红包
	commons.GetDB().Where("envelope_id", id).First(&envelope)
	return &envelope
}

//更新红包是否打开
func UpdateState(id int, open int) {
	//更新单列
	commons.GetDB().Model(&Envelope{}).Where("envelope_id = ?", id).Update("opened", open)
}

func SaveEnvelope(envelope Envelope) {
	commons.GetDB().Create(envelope)
}
