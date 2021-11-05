package models

//生成每个用户的红包金额

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"huhusw.com/red_envelope/commons"
)

type MoneySetting struct {
	SumEnvelope int64   `json:"sum_envelope"`
	Probability float32 `json:"probability"`
	Money       int64   `json:"money"`
}

//全局配置
var moneySetting = &MoneySetting{}

//创建channel，存储当前生成的红包金额
var moneyChan chan int64

//60s的定时器，查询redis中配置的变动
var tickerMoney = time.NewTicker(time.Second * 60)

//初始化函数，结构体加载时调用，将数据库连接信息赋值，并初始化一个数据库连接
func init() {
	//加载数据库文件
	filePtr, err := os.Open("./config/config.json") //config的文件目录
	if err != nil {
		fmt.Printf("Open file failed [Err:%s]\n", err.Error())
		return
	}
	//关闭文件
	defer filePtr.Close()

	// 创建json解码器
	decoder := json.NewDecoder(filePtr)

	//读取配置文件中的信息，初始化全局的金额配置
	err = decoder.Decode(moneySetting)
	if err != nil {
		fmt.Printf("json decode error [Err:%s]\n", err.Error())
	}
	// fmt.Println(moneySetting.SumEnvelope)
	// fmt.Println(moneySetting.Probability)
	// fmt.Println(moneySetting.Money)
	//初始化全局配置
	commons.GetRedis().Set(context.Background(), "sum_envelope", moneySetting.SumEnvelope, 0)
	commons.GetRedis().Set(context.Background(), "probability", moneySetting.Probability, 0)
	commons.GetRedis().Set(context.Background(), "money", moneySetting.Money, 0)

	//开启定时器任务，查询配置变量的变化
	go timeFunc()

	//分配空间，有缓存的channel，缓存大小为1000个红包金额
	moneyChan = make(chan int64, 10000)

	// 生成红包金额的协程
	go randomNormalInt64()
}

//定时查看redis中的值，做到实时更新
func timeFunc() {
	for _ = range tickerMoney.C {
		fmt.Printf("ticked at %v\n", time.Now())
	}
}

//基于截尾正态分布版红包,min 红包的最小值，max 红包的最大值
// func randomNormalInt64() {
// 	//红包金额的平均值，金额产生百分之一的误差
// 	var mean float64 = moneySetting.Money * 0.99 / float64(moneySetting.SumEnvelope)
// 	fmt.Println("mean = ", mean)
// 	//标准差，采用平均值/2 替代真实的样本集的标准差
// 	var sigma float64 = mean / 2
// 	fmt.Println("sigma = ", sigma)

// 	//全部红包的权重之和
// 	var sum_ float64 = 0.0

// 	// 每个红包额度权重
// 	var divideTable []float64 = make([]float64, moneySetting.SumEnvelope)

// 	//切片数组的索引
// 	var index int64 = 0
// 	//给每一个红包分配权重
// 	for ; index < moneySetting.SumEnvelope; index++ {
// 		randomWight := normalVariate(mean, sigma)
// 		if randomWight >= sigma && randomWight <= (2*sigma) {
// 			// fmt.Println("randomWight = ", randomWight)
// 			divideTable[index] = randomWight
// 			sum_ += randomWight
// 		}
// 	}
// 	var curSum float64 = 0.0
// 	index = 0
// 	//循环生成前n-1个红包金额
// 	for ; index < moneySetting.SumEnvelope-1; index++ {
// 		curMoney := int64(divideTable[index] / sum_ * moneySetting.Money)
// 		if curMoney < int64(mean/16) {
// 			curMoney = int64(mean / 16)
// 		}
// 		moneyChan <- int(curMoney)
// 		curSum += float64(curMoney)
// 	}
// 	//最后一个红包的额度
// 	curMoney := int64(moneySetting.Money - curSum)
// 	moneyChan <- int(curMoney)
// }

// func normalVariate(mean float64, sgima float64) float64 {
// 	//将时间戳设置成种子数
// 	// rand.Seed(time.Now().UnixNano())
// 	//映射正态分布概率
// 	z := 0.0
// 	for {
// 		u1 := rand.Float64()
// 		u2 := 1 - rand.Float64()
// 		z = 4 * math.Exp(-0.5) / math.Sqrt(2) * (u1 - 0.5) / u2
// 		zz := z * z / 4.0
// 		if zz <= -math.Log(u2) {
// 			break
// 		}
// 	}
// 	return (mean + z*sgima)
// }

//其实没有做到正态分布，更像一个平均分布，这个函数有问题，5亿的预算发了38亿出来
//正态分布随机数生产器：min:最小值，max:最大值，miu:期望值（均值），sigma:方差
// func RandomNormalInt64(min int64, max int64, miu int64, sigma int64) int64 {
func randomNormalInt64() {
	// if min >= max {
	// 	return 0
	// }
	// if miu < min {
	// 	miu = min
	// }
	// if miu > max {
	// 	miu = max
	// }
	//红包金额的平均值，金额产生百分之一的误差
	var miu int64 = int64(float64(moneySetting.Money) * 0.95 / float64(moneySetting.SumEnvelope))
	var sigma int64 = 2
	var min = (miu / 8)
	var max = miu * 16
	var x int64
	var y, dScope float64
	for {
		x = RandInt64(min, max)
		y = NormalFloat64(x, miu, sigma) * 100000
		dScope = float64(RandInt64(0, int64(NormalFloat64(x, miu, sigma)*100000)))
		//注意下传的是两个miu
		if dScope <= y {
			break
		}
	}
	moneyChan <- x
}

func TTT() int64 {
	// if min >= max {
	// 	return 0
	// }
	// if miu < min {
	// 	miu = min
	// }
	// if miu > max {
	// 	miu = max
	// }
	//红包金额的平均值，金额产生百分之一的误差
	var miu int64 = int64(float64(moneySetting.Money) * 0.95 / float64(moneySetting.SumEnvelope))
	var sigma int64 = 2
	var min = (miu / 8)
	var max = miu * 16
	var x int64
	var y, dScope float64
	for {
		x = RandInt64(min, max)
		y = NormalFloat64(x, miu, sigma) * 100000
		dScope = float64(RandInt64(0, int64(NormalFloat64(x, miu, sigma)*100000)))
		//注意下传的是两个miu
		if dScope <= y {
			break
		}
	}
	return x
}

//正态分布公式
func NormalFloat64(x int64, miu int64, sigma int64) float64 {
	randomNormal := 1 / (math.Sqrt(2*math.Pi) * float64(sigma)) * math.Pow(math.E, -math.Pow(float64(x-miu), 2)/(2*math.Pow(float64(sigma), 2)))
	//注意下是x-miu
	return randomNormal
}

//产生min-max之间的随机数
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func GetProbability() float32 {
	return moneySetting.Probability
}

//生成金额
func GetAmount() int64 {
	return <-moneyChan
}
