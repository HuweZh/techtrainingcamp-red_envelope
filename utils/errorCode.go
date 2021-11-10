package utils

//定义错误码
const (
	CODE_SUCCESS                      = 0
	CODE_NOT_LOGIN_ERROR              = 1 //未登录
	CODE_PARAMETER_ERROR              = 2 //post参数错误
	CODE_STRING_TO_INT_ERROR          = 3 //string转int错误
	CODE_BINDJSON_ERROR               = 4
	CODE_MARSHAL_ERROR                = 5 //json.Marshal错误
	CODE_UNMARSHAL_ERROR              = 6 //json.Unmarshal
	CODE_REDIS_GET_ERROR              = 7
	CODE_REDIS_SET_ERROR              = 8
	CODE_OUT_OF_REDENVELOPES_ERROR    = 9 //总红包用尽
	CODE_INSERT_DB_ERROR              = 10
	CODE_OUT_OF_SNATCH_COUNT_ERROR    = 11 //抢红包次数用尽
	CODE_ENVELOPE_NOT_EXIST_ERROR     = 12
	CODE_UPDATE_DB_ERROR              = 13
	CODE_OTHER_ERROR                  = 14
	CODE_SERVER_INTERNAL_ERROR        = 15
	CODE_LIMIT_HAS_BEEN_REACHED_ERROR = 16
	CODE_DID_NOT_SNATCH               = 17
	CODE_ADD_MONEY_ERROR              = 18

	MSG_SUCCESS                      = "success"
	MSG_NOT_LOGIN_ERROR              = "not login"
	MSG_PARAMETER_ERROR              = "parameter error"
	MSG_STRING_TO_INT_ERROR          = "convert string to int error"
	MSG_BINDJSON_ERROR               = "bind json error"
	MSG_MARSHAL_ERROR                = "json marshal error"
	MSG_UNMARSHAL_ERROR              = "json unmarshal error"
	MSG_REDIS_GET_ERROR              = "redis get error"
	MSG_REDIS_SET_ERROR              = "redis set error"
	MSG_OUT_OF_REDENVELOPES_ERROR    = "out of red envelopes error"
	MSG_INSERT_DB_ERROR              = "insert db error"
	MSG_OUT_OF_SNATCH_COUNT_ERROR    = "out of snatch count"
	MSG_ENVELOPE_NOT_EXIST_ERROR     = "envelope not exist"
	MSG_UPDATE_DB_ERROR              = "update db error"
	MSG_OTHER_ERROR                  = "other errors"
	MSG_SERVER_INTERNAL_ERROR        = "server internal error"
	MSG_LIMIT_HAS_BEEN_REACHED_ERROR = "limit has been reached"
	MSG_DID_NOT_SNATCH               = "did not snatch red envelope"
	MSG_ADD_MONEY_ERROR              = "add money error"
)
