package models

type Envelope struct {
	EnvelopeId int `json:"envelope_id"`
	UserId     int
	Value      int
	State      int
	SnatchTime int
}

//默认操作的是envelopes表
//改变结构体的默认表名称
func (envelope Envelope) TableName() string {
	return "envelope"
}
