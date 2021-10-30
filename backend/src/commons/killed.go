package commons

type Killed struct {
	SeckillID int
	UserPhone int
	State     int
}

func (killed Killed) TableName() string {
	return "killed"
}
