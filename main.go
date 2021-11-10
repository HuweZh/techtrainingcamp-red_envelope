package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"red_envelope/configure"
	"red_envelope/middlewares"
	"red_envelope/routers"
	"red_envelope/utils"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api" //前面加一个，这样即便没使用也不会报没使用的错了
	"github.com/alibaba/sentinel-golang/core/flow"
	_ "net/http/pprof"
)

func main() {

	//性能优化分析,访问http://localhost:6060/debug/pprof/
	go func() {
		if configure.UseProfiler {
			err := http.ListenAndServe("localhost:6060", nil)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

	//初始化数据库，缓存等资源
	rand.Seed(time.Now().Unix())
	var err error
	err = utils.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	/*
	//初始化日志记录
	utils.InitLogger()

	//初始化RMQ
	utils.InitRMQ()

	//初始化sentinel配置和设置埋点(定义资源)
	//https://sentinelguard.io/zh-cn/docs/golang/flow-control.html
	err = sentinel.InitWithConfigFile("./configure/sentinel_configure.yaml")
	if err != nil {
		utils.MyLog.Error(err)
		return
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "request",
			Threshold:              2000, //流控阈值；如果字段 StatIntervalInMs 是1000(也就是1秒)，那么Threshold就表示QPS
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling, //流控效果为匀速排队flow.Throttling或直接拒绝flow.Reject
			StatIntervalInMs:       1000,            //规则对应的流量控制器的独立统计结构的统计周期
			MaxQueueingTimeMs:      1000,            //匀速排队的最大等待时间ms
		},
	})
	if err != nil {
		utils.MyLog.Error(err)
		return
	}

	//初始化gin
	r := gin.Default()
	r.Use(middlewares.RecoveryMiddleWare())
	r.Use(middlewares.Cors())
	r.Use(middlewares.AntiCrawler())
	gin.SetMode(configure.GIN_MODE)
	routers.ApiRoutersInit(r)
	routers.AdminRoutersInit(r)

	err = r.Run(":8080")
	if err != nil {
		utils.MyLog.Error(err)
		return
	}
	*/
}
