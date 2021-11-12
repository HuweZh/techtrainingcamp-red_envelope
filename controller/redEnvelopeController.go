package controller

import (
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"huhusw.com/red_envelope/commons"
	"huhusw.com/red_envelope/models"
)

/*
红包相关接口控制器
*/
type RedEnvelopeController struct {
}

type RequestParamter struct {
	Uid         int `json:"uid"`
	Envelope_id int `json:"envelope_id"`
}

//抢红包业务逻辑
func (con RedEnvelopeController) Snatch(c *gin.Context) {
	//获取请求中的数据
	value, _ := c.Get("user")
	user := value.(models.User)

	// 存储此红包为该用户的一个红包
	var envelope models.Envelope = models.GetEnve(user.UserId)

	if envelope.Value == 0 {
		//没抢到红包
		commons.R(c, commons.NOTGETONE, commons.SNATCHERROR, nil)
	} else {
		//存入数据库
		models.SetMysqlData(commons.INSERTENVELOPE, envelope)
		models.SetRedisData(commons.RPUSH, "uid"+strconv.Itoa(user.UserId), envelope, 0)
		models.SetRedisData(commons.EXPIRE, "uid"+strconv.Itoa(user.UserId), nil, 600*1000000000)

		//构建返回的数据
		data := gin.H{
			"envelope_id": envelope.EnvelopeId,
			"max_count":   user.MaxCount,
			"cur_count":   user.CurCount,
		}

		//返回数据
		commons.R(c, commons.OK, commons.SUCCESS, data)
	}
}

//打开红包业务逻辑
func (con RedEnvelopeController) Open(c *gin.Context) {
	//获取请求参数
	value, _ := c.Get("envelope")
	envelope := value.(models.Envelope)

	//构建返回的数据
	data := gin.H{
		"value": envelope.Value,
	}
	commons.R(c, commons.OK, commons.SUCCESS, data)
}

//获取钱包列表业务逻辑
func (con RedEnvelopeController) GetWalletList(c *gin.Context) {
	//获取请求携带的参数
	value, _ := c.Get("envelopes")
	amount := c.GetInt("amount")
	envelopes := value.([]models.Envelope)

	//对红包按照时间戳排序
	sort.Slice(envelopes, func(i, j int) bool {
		return envelopes[i].SnatchTime < envelopes[j].SnatchTime
	})

	//构建返回数据，并返回
	data := gin.H{
		"amount":        amount,
		"envelope_list": envelopes,
	}
	commons.R(c, commons.OK, commons.SUCCESS, data)
}
