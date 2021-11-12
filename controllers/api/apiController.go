package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"net/http"
	"red_envelope/configure"
	"red_envelope/models"
	"red_envelope/utils"
	"strconv"
)

// ApiController
//把方法挂在结构体下面,方便以后继承
type ApiController struct{}

type GetWalletListRequestBody struct {
	Uid int64 `json:"uid"`
}

type OpenRequestBody struct {
	Uid        int64 `json:"uid"`
	EnvelopeId int64 `json:"envelope_id"`
}

func (con ApiController) SnatchHandler(c *gin.Context) {
	var request GetWalletListRequestBody
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_BINDJSON_ERROR,
			"msg":  utils.MSG_BINDJSON_ERROR,
		})
		panic(err)
	}

	//先检查抢红包次数是否用尽
	user := models.User{}
	key := strconv.FormatInt(request.Uid, 10) + "user" //我是直接缓存用户表的，就没有操作数据库了
	userJson, err := utils.RDB.Get(utils.CTX, key).Result()
	if err != redis.Nil {
		err = json.Unmarshal([]byte(userJson), &user)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": utils.CODE_UNMARSHAL_ERROR,
				"msg":  utils.MSG_UNMARSHAL_ERROR,
			})
			panic(err)
		}
	} else {
		user.Uid = request.Uid
		user.MaxCount = configure.MaxSnatch
		user.CurCount = 0
	}
	user.CurCount += 1
	if user.CurCount > user.MaxCount { //抢红包次数是否用尽
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_OUT_OF_SNATCH_COUNT_ERROR,
			"msg":  utils.MSG_OUT_OF_SNATCH_COUNT_ERROR,
		})
		return
	}
	userJsonByte, _ := json.Marshal(user)
	if err := utils.RDB.Set(utils.CTX, key, userJsonByte, 0).Err(); err != nil { //更新redis
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_REDIS_SET_ERROR,
			"msg":  utils.MSG_REDIS_SET_ERROR,
		})
		panic(err)
	}

	//生成随机数是否抢到红包
	r := rand.Float64()
	if r > configure.SnatchP {//没抢到红包
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_DID_NOT_SNATCH,
			"msg":  utils.MSG_DID_NOT_SNATCH,
		})
		return
	}

	//抢红包
	if utils.RDB.LLen(utils.CTX, "allEnvelopeList").Val() == 0 { //红包抢完了
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_OUT_OF_REDENVELOPES_ERROR,
			"msg":  utils.MSG_OUT_OF_REDENVELOPES_ERROR,
		})
		return
	}
	envelopeJson := utils.RDB.LPop(utils.CTX, "allEnvelopeList").Val()
	newEnvelope := models.Envelope{}
	err = json.Unmarshal([]byte(envelopeJson), &newEnvelope)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_UNMARSHAL_ERROR,
			"msg":  utils.MSG_UNMARSHAL_ERROR,
		})
		panic(err)
	}
	newEnvelope.Uid = request.Uid
	newEnvelope.SnatchTime = utils.GetCurrentTime()

	//先更新数据库再更新redis
	envelopeList := []models.Envelope{}
	envelopeList = models.GetEnvelopesByUid(utils.DB, request.Uid)
	key = strconv.FormatInt(request.Uid, 10) + "wallet" //红包列表key为uid+"wallet"
	envelopeList = append(envelopeList, newEnvelope)
	if err := models.InsertEnvelope(utils.DB, &newEnvelope); err != nil { //更新数据库
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_INSERT_DB_ERROR,
			"msg":  utils.MSG_INSERT_DB_ERROR,
		})
		panic(err)
	}
	//转为json格式然后更新redis
	jsonEnvelopeListByte, err := json.Marshal(envelopeList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_MARSHAL_ERROR,
			"msg":  utils.MSG_MARSHAL_ERROR,
		})
		panic(err)
	}
	jsonEnvelopeList := string(jsonEnvelopeListByte)
	if err := utils.RDB.Set(utils.CTX, key, jsonEnvelopeList, 0).Err(); err != nil { //更新redis
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_REDIS_SET_ERROR,
			"msg":  utils.MSG_REDIS_SET_ERROR,
		})
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": utils.CODE_SUCCESS,
		"msg":  utils.MSG_SUCCESS,
		"data": gin.H{
			"envelope_id": newEnvelope.EnvelopeId,
			"max_count":   user.MaxCount,
			"cur_count":   user.CurCount,
		},
	})
	return
}

func (con ApiController) OpenHandler(c *gin.Context) {
	var request OpenRequestBody
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_BINDJSON_ERROR,
			"msg":  utils.MSG_BINDJSON_ERROR,
		})
		panic(err)
	}
	//从redis中获取该用户的所有红包
	envelopeList := getUserEnvelopeList(request.Uid, c)
	if envelopeList == nil {
		return
	}

	//查找该红包是否存在
	index := 0
	for index < len(envelopeList) {
		if envelopeList[index].EnvelopeId == request.EnvelopeId {
			break
		}
		index += 1
	}
	if index == len(envelopeList) {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_ENVELOPE_NOT_EXIST_ERROR,
			"msg":  utils.MSG_ENVELOPE_NOT_EXIST_ERROR,
		})
		return
	}

	//如果红包没打开则更新redis和数据库
	if !envelopeList[index].Opened {
		envelopeList[index].Opened = true
		envelopeList[index].OpenedTime = utils.GetCurrentTime()
		jsonEnvelopeList, err := json.Marshal(envelopeList)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": utils.CODE_MARSHAL_ERROR,
				"msg":  utils.MSG_MARSHAL_ERROR,
			})
			panic(err)
		}
		key := strconv.FormatInt(request.Uid, 10) + "wallet"                                     //key为uid+"wallet"
		if err := utils.RDB.Set(utils.CTX, key, string(jsonEnvelopeList), 0).Err(); err != nil { //更新redis
			c.JSON(http.StatusOK, gin.H{
				"code": utils.CODE_REDIS_SET_ERROR,
				"msg":  utils.MSG_REDIS_SET_ERROR,
			})
			panic(err)
		}

		if err := models.UpdateEnvelopeByEnvelopeId(utils.DB, request.EnvelopeId, &map[string]interface{}{"opened": true, "opened_time": envelopeList[index].OpenedTime}); err != nil { //更新数据库
			c.JSON(http.StatusOK, gin.H{
				"code": utils.CODE_UPDATE_DB_ERROR,
				"msg":  utils.MSG_UPDATE_DB_ERROR,
			})
			panic(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.CODE_SUCCESS,
		"msg":  utils.MSG_SUCCESS,
		"data": gin.H{
			"value": envelopeList[index].Value,
		},
	})
	return
}

func (con ApiController) GetWalletListHandler(c *gin.Context) {
	var request GetWalletListRequestBody
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.CODE_BINDJSON_ERROR,
			"msg":  utils.MSG_BINDJSON_ERROR,
		})
		panic(err)
	}
	envelopeList := getUserEnvelopeList(request.Uid, c)
	if envelopeList == nil {
		return
	}

	dataEnvelopeList := []gin.H{}
	var dataAmount int32 = 0
	for _, envelope := range envelopeList { //下划线表示忽略返回值
		if envelope.Opened {
			dataAmount += envelope.Value
			dataEnvelopeList = append(dataEnvelopeList, gin.H{
				"envelope_id": envelope.EnvelopeId,
				"value":       envelope.Value, //拆开的红包显示金额
				"opened":      true,
				"snatch_time": envelope.SnatchTime,
			})
		} else {
			dataEnvelopeList = append(dataEnvelopeList, gin.H{
				"envelope_id": envelope.EnvelopeId,
				"opened":      false,
				"snatch_time": envelope.SnatchTime,
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": utils.CODE_SUCCESS,
		"msg":  utils.MSG_SUCCESS,
		"data": gin.H{
			"amount":        dataAmount,
			"envelope_list": dataEnvelopeList,
		},
	})

	//钱包信息发往消息队列处理
// 	wallet := models.Wallet{
// 		Uid:   request.Uid,
// 		Money: dataAmount,
// 	}
// 	walletJson, err := json.Marshal(wallet)
// 	if err == nil {
// 		rmqMsg := utils.RocketMqMessage{
// 			Topic:        "wallet",
// 			MessageBytes: walletJson,
// 		}
// 		utils.SendToRMQ(rmqMsg)
// 	}
}

func getUserEnvelopeList(uid int64, c *gin.Context) []models.Envelope {
	envelopeList := []models.Envelope{}
	//先从redis中查询
	key := strconv.FormatInt(uid, 10) + "wallet" //key为uid+"wallet"
	jsonEnvelopeList, err := utils.RDB.Get(utils.CTX, key).Result()
	if err != redis.Nil {
		err = json.Unmarshal([]byte(jsonEnvelopeList), &envelopeList)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": utils.CODE_UNMARSHAL_ERROR,
				"msg":  utils.MSG_UNMARSHAL_ERROR,
			})
			panic(err)
			return nil
		}
	} else {
		//redis查询不到再从数据库中查询，并更新缓存
		envelopeList = models.GetEnvelopesByUid(utils.DB, uid)
		//转为json格式然后更新redis
		jsonEnvelopeList, err := json.Marshal(envelopeList)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": utils.CODE_MARSHAL_ERROR,
				"msg":  utils.MSG_MARSHAL_ERROR,
			})
			panic(err)
			return nil
		}
		utils.RDB.Set(utils.CTX, key, string(jsonEnvelopeList), 0) //更新redis
	}
	return envelopeList
}
