package models

import "gorm.io/gorm"

type Wallet struct {
	Uid   int64 `json:"uid"`
	Money int32 `json:"money"`
}

func (Wallet) TableName() string {
	return "wallet"
}


func InsertWallet(DB *gorm.DB, newWallet *Wallet) error {
	return DB.Create(newWallet).Error
}

func UpdateWalletByUid(DB *gorm.DB, uid int64, data *map[string]interface{}) error {
	return DB.Model(&Wallet{}).Where("uid=?", uid).Updates(data).Error
}
