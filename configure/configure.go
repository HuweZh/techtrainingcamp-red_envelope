/*
配置文件
*/
package configure

const (
	MaxSnatch        = 10        //最多抢红包数
	TotalMoney       = 100000000 //红包总金额
	TotalAmount      = 1000      //红包总个数
	MinMoney         = 1
	MaxMoney         = 2*TotalMoney/TotalAmount - MinMoney //这样就尽量做到把钱用完
	SnatchP          = 0.7                                 //抢到红包的概率
	MysqlAddr        = "172.16.10.188:3306"
	MysqlUser        = "root"
	MysqlPass        = "root"
	MysqlDatabase    = "red_envelope"
	MachineId        = 1 //分布式机器编号
	RedisAddr        = "172.16.70.164:6379"
	RedisPass        = ""
	JSMAXNUM         = 1 << 53 //js整数的取值范围在[-2^53,2^53]
	GIN_MODE         = "debug" //debug or release
	RmqNameseverAddr = "172.16.69.218:9876"
	RmqBrokerAddr    = "172.16.244.9:10911"
	UseProfiler      = true
)
