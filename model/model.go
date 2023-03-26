package model

import "gorm.io/gorm"

func Init(tx *gorm.DB) {
	initUser(tx)
	initPrize(tx)
	initDrawRecord(tx)
	initUserDrawRecord(tx)
}
