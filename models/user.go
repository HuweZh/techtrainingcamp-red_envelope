package models

type User struct {
	Uid      int64 `json:"uid"`
	MaxCount int32 `json:"max_count"` // 最多抢几次
	CurCount int32 `json:"cur_count"` // 当前为第几次抢
}

func (User) TableName() string {
	return "user"
}
