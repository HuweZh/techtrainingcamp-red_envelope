package models

//更新数据库的类型和数据
type updateData struct {
	Type int
	Data interface{}
}

//需要操作数据库更新的channel
var UpdateChann chan updateData

//初始化channel
func init() {
	UpdateChann = make(chan updateData, 1000)
	go update()
}

//传数据到channel
func SetData(t int, d interface{}) {
	UpdateChann <- updateData{t, d}
}

//开启gorouine，持续判断channel中的数据，做到读写分离
func update() {
	for {
		rowData := <-UpdateChann
		switch rowData.Type {
		//插入新红包
		case 0:
			d := rowData.Data.(Envelope)
			SaveEnvelope(d)
		//更新当前红包的状态，由未打开变成已打开
		case 1:
			d := rowData.Data.(Envelope)
			UpdateState(d.EnvelopeId, d.Opened)
		//更新当前用户的红包计数
		case 2:
			d := rowData.Data.(User)
			UpdateCurCount(d.UserId, d.CurCount)
		}
	}
}
