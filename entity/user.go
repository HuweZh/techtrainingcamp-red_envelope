package entity

type User struct {
	Uid       int `json:"uid"`
	Envelopes []int
}
