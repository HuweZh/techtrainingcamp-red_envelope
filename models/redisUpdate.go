package models

import (
	"context"
	"time"

	"huhusw.com/red_envelope/commons"
)

//更新数据库的类型和数据
type redisData struct {
	Type   int
	Key    string
	Value  interface{}
	Period int
}

//需要操作数据库更新的channel
var redisChann chan redisData

//初始化channel
func init() {
	redisChann = make(chan redisData, 1000)
	go updateRedis()
}

//传数据到channel
func SetRedisData(t int, k string, v interface{}, p int) {
	redisChann <- redisData{t, k, v, p}
}

//开启gorouine，持续判断channel中的数据，做到读写分离
func updateRedis() {
	for {
		rowData := <-redisChann
		switch rowData.Type {
		//redis set，设置user的键和值
		case commons.SET:
			key := rowData.Key
			value := rowData.Value.(User)
			expire := rowData.Period
			commons.GetRedis().Set(context.Background(), key, value, time.Duration(expire))
		//redis rpush，将红包push到user对应的列表中
		case commons.RPUSH:
			key := rowData.Key
			value := rowData.Value.(Envelope)
			commons.GetRedis().RPush(context.Background(), key, value)
		//设置过期时间，针对的是user对应的列表
		case commons.EXPIRE:
			key := rowData.Key
			expire := rowData.Period
			commons.GetRedis().Expire(context.Background(), key, time.Duration(expire))
		//删除指定键中的指定元素，针对的是红包从未打开变成打开
		case commons.LREM:
			key := rowData.Key
			value := rowData.Value.(Envelope)
			commons.GetRedis().LRem(context.Background(), key, 1, value)
		}
	}
}
