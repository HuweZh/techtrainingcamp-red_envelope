package models

import "huhusw.com/red_envelope/commons"

//更新数据库的类型和数据
type mysqlData struct {
	Type int
	Data interface{}
}

//需要操作数据库更新的channel
var mysqlChann chan mysqlData

//初始化channel
func init() {
	mysqlChann = make(chan mysqlData, 1000)
	go update()
}

//传数据到channel
func SetMysqlData(t int, d interface{}) {
	mysqlChann <- mysqlData{t, d}
}

//开启gorouine，持续判断channel中的数据，做到读写分离
func update() {
	for {
		rowData := <-mysqlChann
		switch rowData.Type {
		//插入新红包
		case commons.INSERTENVELOPE:
			d := rowData.Data.(Envelope)
			SaveEnvelope(d)
		//更新当前红包的状态，由未打开变成已打开
		case commons.UPDATEENVELOPESTATE:
			d := rowData.Data.(Envelope)
			UpdateState(d.EnvelopeId, d.Opened)
		//更新当前用户的红包计数
		case commons.UPDATEUSER:
			d := rowData.Data.(User)
			UpdateCurCount(d.UserId, d.CurCount)
		//更新用户的钱包金额
		case commons.UPDATEAMOUNT:
			d := rowData.Data.(User)
			UpdateAmount(d.UserId, d.Amount)
		}
	}
}
