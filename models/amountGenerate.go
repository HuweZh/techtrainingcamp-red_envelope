package models

//生成每个用户的红包金额

import (
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"huhusw.com/red_envelope/commons"
	"huhusw.com/red_envelope/logger"
)

type MoneySetting struct {
	SumEnvelope int     `json:"sum_envelope"`
	Probability float64 `json:"probability"`
	Money       int     `json:"money"`
}

//全局配置
var moneySetting = &MoneySetting{}

//创建channel，存储当前生成的红包金额
var moneyChan chan int

//红包的平均值
var mean int

//红包最小值，mean/8
var min int

//60s的定时器，查询redis中配置的变动
var tickerConfig = time.NewTicker(time.Second * 60)

//10s的定时器，写入当前的剩余金额和红包数量
var tickerMoney = time.NewTicker(time.Second * 10)

//初始化函数，结构体加载时调用，将数据库连接信息赋值，并初始化一个数据库连接
func init() {
	//加载数据库文件
	filePtr, err := os.Open("./config/config.json") //config的文件目录
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"全局金额配置出错": err.Error(),
		})
		// fmt.Printf("Open file failed [Err:%s]\n", err.Error())
		return
	}
	//关闭文件
	defer filePtr.Close()

	// 创建json解码器
	decoder := json.NewDecoder(filePtr)

	//读取配置文件中的信息，初始化全局的金额配置
	err = decoder.Decode(moneySetting)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"全局金额配置文件解码出错": err.Error(),
		})
		// fmt.Printf("json decode error [Err:%s]\n", err.Error())
	}
	//初始化全局配置
	commons.GetRedis().Set(context.Background(), "sum_envelope", moneySetting.SumEnvelope, 0)
	commons.GetRedis().Set(context.Background(), "probability", moneySetting.Probability, 0)
	commons.GetRedis().Set(context.Background(), "money", moneySetting.Money, 0)

	//平均数
	mean = moneySetting.Money / moneySetting.SumEnvelope
	//每个红包的最小值
	min = mean / 8

	//开启定时器任务，查询配置变量的变化
	go timeFuncConfig()
	go timeFuncMoney()

	//分配空间，有缓存的channel，缓存大小为1000个红包金额
	moneyChan = make(chan int, 10000)

	// 生成红包金额的协程
	go getRandomMoney()
}

//定时查看redis中的值，做到实时更新
func timeFuncConfig() {
	for _ = range tickerConfig.C {
		//读取redis中新配置的值，更新红包金额和数量
		new_money, _ := commons.GetRedis().Get(context.Background(), "new_money").Result()
		new_envelops, _ := commons.GetRedis().Get(context.Background(), "new_envelops").Result()
		new_probability, _ := commons.GetRedis().Get(context.Background(), "new_probability").Result()

		//更新参数
		if new_money != "" && new_envelops != "" && new_probability != "" {
			// fmt.Println("new_money", new_money)
			// fmt.Println("new_envelops", new_envelops)
			// fmt.Println("new_probability", new_probability)
			commons.GetRedis().Del(context.Background(), "new_money", "new_envelops", "new_probability")
			moneySetting.Money, _ = strconv.Atoi(new_money)
			moneySetting.SumEnvelope, _ = strconv.Atoi(new_envelops)
			moneySetting.Probability, _ = strconv.ParseFloat(new_probability, 32)
			//平均数
			mean = moneySetting.Money / moneySetting.SumEnvelope
			//每个红包的最小值
			min = mean / 8
		}
	}
}

//定时查看redis中的值，做到实时更新
func timeFuncMoney() {
	for _ = range tickerMoney.C {
		//设置剩余的红包金额和红包数量
		commons.GetRedis().Set(context.Background(), "sum_envelope", moneySetting.SumEnvelope, 0)
		commons.GetRedis().Set(context.Background(), "probability", moneySetting.Probability, 0)
		commons.GetRedis().Set(context.Background(), "money", moneySetting.Money, 0)
	}
}

//产生min-max之间的随机数
func randInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

// func GetProbability() float32 {
// 	return moneySetting.Probability
// }

/**
 * 随机红包
 */
func getRandomMoney() {
	//产生百分之一的误差
	// money := int(float32(moneySetting.Money) * 0.99)
	for {
		//循环生成随机红包，100个红包为一组进行生成，对应的金额是100*mean，保证生成的红包是稳定的
		randMoney(mean*100, 100)
		moneySetting.Money -= mean * 100
		moneySetting.SumEnvelope -= 100
		//资源耗尽时，只添加空红包
		if moneySetting.Money == 0 || moneySetting.SumEnvelope == 0 {
			moneyChan <- 0
		}
	}
}

/**
     * 随机红包
     * 原理： 第一个红包  额度 = [最小额度,总派发额度/总红包个数 * 2]内任意一个随机数  即为对应红包的额度
     *        第二个红包  额度 = [最小额度,(总派发额度-第一个红包的额度)/(总红包个数-1) * 2]内任意一个随机数  即为对应红包的额度
     *        第三个红包  额度 = [最小额度,(总派发额度-前2个红包的额度)/(总红包个数-2) * 2]内任意一个随机数  即为对应红包的额度
     *        如此循环直至第n-1个红包 额度 = [最小额度,(总派发额度-前n-2个红包的额度)/(总红包个数-(n-2)) * 2]内任意一个随机数  即为对应红包的额度
     *        第n个红包 额度 = 总额度 - 前n-1个红包的派发额度
     *        需要指定当前红包的个数、红包的总额度、红包最小额度
	 **/
func randMoney(amountMoney, count int) {
	//红包金额为0代表抢红包失败，这里countZero记录了这100个红包里有多少个金额为0的红包
	countZero := int(100 * (1 - moneySetting.Probability))
	//100个为一组，循环生成红包金额
	for i := 0; i < 99; i++ {
		max := amountMoney / (count - i) * 2
		curMoney := randInt(min, max)
		moneyChan <- curMoney
		amountMoney -= curMoney
		//每个人没抢到红包的概率，插入一个空红包
		rand.Seed(time.Now().UnixNano())
		if (float64(rand.Intn(1000))/1000.0) < (1-moneySetting.Probability) && countZero != 0 {
			moneyChan <- 0
			countZero -= 1
		}
	}
	//最后一个红包
	moneyChan <- amountMoney
	//最后剩余的空红包
	for i := 0; i < countZero; i++ {
		moneyChan <- 0
	}
}

//生成金额
func GetAmount() int {
	return <-moneyChan
}
