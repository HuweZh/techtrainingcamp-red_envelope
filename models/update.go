package models

import (
	"fmt"
)

//常亮，表示更新数据库的类型
const (
	INSERTENVELOPE      = 0
	UPDATEENVELOPESTATE = 1
	UPDATEUSER          = 2
)

//更新数据库的类型和数据
type UpdateData struct {
	Type int
	Data interface{}
}

//需要操作数据库更新的channel
var UpdateChann chan UpdateData

//初始化channel
func init() {
	UpdateChann = make(chan UpdateData, 1000)
	go update()
}

//传数据到channel
func SetData(u UpdateData) {
	UpdateChann <- u
}

//开启gorouine，持续判断channel中的数据，做到读写分离
func update() {
	for {
		rowData := <-UpdateChann
		switch rowData.Type {
		case 0:
			d := rowData.Data.(Envelope)
			SaveEnvelope(d)
		case 1:
			fmt.Println("case 1")
		}
	}
}
