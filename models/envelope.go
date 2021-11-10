package models

import (
	"gorm.io/gorm"
)

//首字母大写表示公有权限
//type User struct {
//	Uid int32 `json:"uid"` //转成json格式时，Uid会变成uid
//}

type Envelope struct {
	EnvelopeId int64 `json:"envelope_id"` // 红包id
	Uid        int64 `json:"uid"`         //拥有者
	Value      int32 `json:"value"`       //红包金额（分）
	Opened     bool  `json:"opened"`      //是否打开
	SnatchTime int64 `json:"snatch_time"` //抢到时间
	OpenedTime int64 `json:"opened_time"` //打开时间
}

// TableName 指定表名
func (Envelope) TableName() string {
	return "envelope"
}

func GetEnvelopesByUid(DB *gorm.DB, uid int64) []Envelope {
	envelopeList := []Envelope{}
	DB.Where("uid=?", uid).Find(&envelopeList)
	return envelopeList
}

func InsertEnvelope(DB *gorm.DB, newEnvelope *Envelope) error {
	return DB.Create(newEnvelope).Error
}

func UpdateEnvelopeByEnvelopeId(DB *gorm.DB, envelopeId int64, data *map[string]interface{}) error {
	return DB.Model(&Envelope{}).Where("envelope_id=?", envelopeId).Updates(data).Error
}
