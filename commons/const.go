package commons

//更新mysql的类型
const (
	INSERTENVELOPE      = 0
	UPDATEENVELOPESTATE = 1
	UPDATEUSER          = 2
	UPDATEAMOUNT        = 3
)

//更新redis的类型
const (
	SET    = 0
	RPUSH  = 1
	EXPIRE = 2
	LREM   = 3
)

//返回数据的状态码
const (
	OK        = 0
	BASEERROR = 1
)

//返回数据的信息
const (
	SUCCESS     = "success"
	RUNOUTOF    = "抢红包次数用完"
	OPENED      = "红包已经打开，请勿重复打开"
	HAVEZERO    = "此用户没用红包"
	SNATCHERROR = "未抢到红包"
)
