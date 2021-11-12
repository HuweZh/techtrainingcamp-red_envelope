package commons

//更新mysql的类型
const (
	INSERTENVELOPE      = 0 //插入新的红包
	UPDATEENVELOPESTATE = 1 //更新红包的打开状态
	UPDATEUSER          = 2 //更新user的抢红包次数
	UPDATEAMOUNT        = 3 //更新用户的钱包总数
)

//更新redis的类型
const (
	SET    = 0 //设置键值对
	RPUSH  = 1 //列表中插入新值
	EXPIRE = 2 //设置过期时间
	LREM   = 3 //从列表中删除一个元素
)

//返回数据的状态码
const (
	OK         = 0 //正常
	NOTGETONE  = 1 //没有抢到红包
	REPEATOPEN = 2 //重复打开相同的红包
	NOCHANCE   = 3 //没有抢红包的机会了
	NOENVELOPE = 4 //此用户的钱包列表为空
)

//返回数据的信息
const (
	SUCCESS     = "success"
	RUNOUTOF    = "抢红包次数用完"
	OPENED      = "红包已经打开，请勿重复打开"
	HAVEZERO    = "此用户没有红包"
	SNATCHERROR = "未抢到红包"
)
