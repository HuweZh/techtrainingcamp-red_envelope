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
	CODE_SUCCESS                   = 0  //成功
	CODE_NOT_LOGIN_ERROR           = 1  //用户未登录
	CODE_PARAMETER_ERROR           = 2  //post参数错误
	CODE_STRING_TO_INT_ERROR       = 3  //string转int错误
	CODE_BINDJSON_ERROR            = 4  //绑定json错误
	CODE_MARSHAL_ERROR             = 5  //string转json错误
	CODE_UNMARSHAL_ERROR           = 6  //json解析错误
	CODE_REDIS_GET_ERROR           = 7  //redis的get操作错误
	CODE_REDIS_SET_ERROR           = 8  //redis的set操作错误
	CODE_OUT_OF_REDENVELOPES_ERROR = 9  //总红包用尽
	CODE_INSERT_DB_ERROR           = 10 //插入数据库错误
	CODE_OUT_OF_SNATCH_COUNT_ERROR = 11 //抢红包次数用尽
	CODE_ENVELOPE_NOT_EXIST_ERROR  = 12 //红包不存在
	CODE_UPDATE_DB_ERROR           = 13 //更新数据库错误
	CODE_OTHER_ERROR               = 14 //其他错误
	CODE_SERVER_INTERNAL_ERROR     = 15 //服务器内部错误
	CODE_DID_NOT_SNATCH            = 16 //抢红包失败
	CODE_ADD_MONEY_ERROR           = 17 //添加总金额失败
	CODE_REPEAT_ENVELOPE           = 18 //重复打开红包
	CODE_NO_ENVELOPE               = 19 //用户没有红包
)

//返回数据的信息
const (
	SUCCESS     = "success"
	RUNOUTOF    = "抢红包次数用完"
	OPENED      = "红包已经打开，请勿重复打开"
	HAVEZERO    = "此用户没有红包"
	SNATCHERROR = "未抢到红包"
)
