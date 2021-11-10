package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"red_envelope/configure"
	"red_envelope/models"
	"strconv"
)

//全局变量
var SnowflakeNode *snowflake.Node = nil
var DB *gorm.DB = nil
var RDB *redis.Client = nil
var CTX = context.Background() //一个空的上下文

func Init() (err error) {
	//初始化雪花算法
	SnowflakeNode, err = snowflake.NewNode(configure.MachineId)
	if err != nil {
		return err
	}
	//初始化数据库连接
	var s string
	s = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configure.MysqlUser,
		configure.MysqlPass,
		configure.MysqlAddr,
		configure.MysqlDatabase,
	)
	DB, err = gorm.Open(mysql.Open(s), &gorm.Config{})
	if err != nil {
		return err
	}
	//初始化redis连接
	RDB = redis.NewClient(&redis.Options{
		Addr:     configure.RedisAddr,
		Password: configure.RedisPass,
		DB:       0, // use default DB
	})
	_, err = RDB.Ping(CTX).Result()
	if err != nil {
		return err
	}
	RDB.FlushDB(CTX)

	//首先减去数据库中已有的红包
	envelopeList := []models.Envelope{}
	DB.Find(&envelopeList)
	for _, envelope := range envelopeList {
		remainSize -= 1
		remainMoney -= int64(envelope.Value)
		//直接缓存用户表
		user := models.User{}
		key := strconv.FormatInt(envelope.Uid, 10) + "user" //我是直接缓存用户表的，就没有操作数据库了
		userJson, err := RDB.Get(CTX, key).Result()
		if err != redis.Nil {
			err = json.Unmarshal([]byte(userJson), &user)
			if err != nil {
				return err
			}
			user.CurCount += 1
		} else {
			user.Uid = envelope.Uid
			user.MaxCount = configure.MaxSnatch
			user.CurCount = 1
		}
		userJsonByte, _ := json.Marshal(user)
		if err := RDB.Set(CTX, key, userJsonByte, 0).Err(); err != nil { //更新redis
			return err
		}

	}

	err = InitRedEnvelope()

	return err
}

func AddMoney(addMoney int64, addSize int64) (err error) {
	remainMoney += addMoney
	remainSize += addSize
	return InitRedEnvelope()
}

// InitRedEnvelope 初始化生成全部红包放入redis中
func InitRedEnvelope() (err error) {
	//然后生成全部红包放入redis中
	for remainSize > 0 && remainMoney > 0 {
		var envelopeJson []byte
		envelope := models.Envelope{
			EnvelopeId: int64(SnowflakeNode.Generate() % configure.JSMAXNUM), //js整数的取值范围在[-2^53,2^53]
			Uid:        0,
			Value:      int32(GetRandomMoney()),
			Opened:     false,
			SnatchTime: 0,
			OpenedTime: 0,
		}
		envelopeJson, err = json.Marshal(envelope)
		if err != nil {
			return err
		}
		RDB.RPush(CTX, "allEnvelopeList", envelopeJson)
	}
	return
}
