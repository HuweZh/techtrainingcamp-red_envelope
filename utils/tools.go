package utils

import (
	"math"
	"math/rand"
	"red_envelope/configure"
	"time"
)

// UnixToTime
//时间戳转换日期
func UnixToTime(timestamp int64) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05") //转换所需的模板
}

// TimeToUnix
//日期转换时间戳
func TimeToUnix(s string) int64 {
	template := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(template, s, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// GetCurrentTime
//获取当前时间戳
func GetCurrentTime() int64 {
	return time.Now().Unix()
}

var remainSize int64 = configure.TotalAmount //剩余的红包数量
var remainMoney int64 = configure.TotalMoney //剩余的钱

// GetRandomMoney
//https://www.zhihu.com/question/22625187/answer/85530416
func GetRandomMoney() int64 {
	if remainSize == 1 {
		remainSize--
		if remainMoney > configure.MaxMoney {
			remainMoney -= configure.MaxMoney
			return configure.MaxMoney
		}
		res := remainMoney
		remainMoney = 0
		return res
	}
	r := rand.Float64()
	min := float64(configure.MinMoney) * 0.01 //分转换为元
	max := float64(configure.MaxMoney) * 0.01
	money := r * max
	if money < min {
		money = min
	}
	money = math.Floor(money*100) / 100
	remainSize--
	ans := int64(money * 100)
	if remainMoney < ans {
		ans = remainMoney
	}
	remainMoney -= ans
	return ans
}
