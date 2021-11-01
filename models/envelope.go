package models

import (
	"encoding/json"
	"time"

	"huhusw.com/red_envelope/commons"
)

type Envelope struct {
	EnvelopeId commons.ID `grom:"envelope_id" json:"envelope_id"`
	UserId     int        `grom:"user_id" json:"user_id"`
	Value      int        `grom:"value" json:"value"`
	Opened     int        `grom:"opened" json:"opened"`
	SnatchTime int        `grom:"snatch_time" json:"snatch_time"`
}

//10s的定时器，用来执行定时任务
var ticker = time.NewTicker(time.Second * 10)

//批量插入红包的切片
var envelopes []Envelope

func init() {
	go func() {

		for _ = range ticker.C {
			// fmt.Printf("ticked at %v", time.Now())
			if len(envelopes) != 0 {
				// fmt.Println("定时器任务执行....")
				commons.GetDB().Create(&envelopes)
				envelopes = envelopes[0:0]
			}
		}
	}()
}

//默认操作的是envelopes表
//改变结构体的默认表名称
func (Envelope) TableName() string {
	return "envelope"
}

func GetEnve(uid int) Envelope {
	return Envelope{commons.GetID(), uid, 50, 0, int(time.Now().Unix())}
}

func GetEnvelope(id int) Envelope {
	envelope := Envelope{}
	//从redis缓存汇总获取红包

	//从数据库中获取红包
	commons.GetDB().Where("envelope_id", id).First(&envelope)
	return envelope
}

//编码json，存入redis
func (e Envelope) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

//更新红包是否打开
func UpdateState(id commons.ID, open int) {
	//更新单列
	commons.GetDB().Model(&Envelope{}).Where("envelope_id = ?", id).Update("opened", open)
}

//批量插入红包数据
func SaveEnvelope(envelope Envelope) {
	if len(envelopes) < 64 {
		envelopes = append(envelopes, envelope)
	} else {
		commons.GetDB().Create(&envelopes)
		envelopes = envelopes[0:0]
	}
}
