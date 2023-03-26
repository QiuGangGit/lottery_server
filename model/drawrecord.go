package model

import (
	"gorm.io/gorm"
	"time"
)

// DrawRecord 抽奖记录
type DrawRecord struct {
	Id int `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	// 奖品id
	PrizeId int `json:"prize_id" gorm:"column:prize_id"`
	// 中奖的用户id
	UserId string `json:"user_id" gorm:"column:user_id"`
	// 中奖时间
	DrawTime int64 `json:"draw_time" gorm:"column:draw_time"`
	// 开始时间
	StartTime int64 `json:"start_time" gorm:"column:start_time"`
}

func (m *DrawRecord) TableName() string {
	return "draw_record"
}

func initDrawRecord(tx *gorm.DB) {
	var prizes []Prize
	tx.Model(&Prize{}).Find(&prizes)
	for _, prize := range prizes {
		tx.Create(&DrawRecord{
			PrizeId:   prize.Id,
			UserId:    "",
			DrawTime:  0,
			StartTime: time.Now().Unix(),
		})
	}
}

// UserDrawRecord 用户抽奖记录
type UserDrawRecord struct {
	UserId       string `json:"user_id" gorm:"column:user_id"`
	DrawRecordId int    `json:"draw_record_id" gorm:"column:draw_record_id"`
	// 状态
	Status     int   `json:"status" gorm:"column:status"` // 0: 没开奖 1: 中奖 2: 没中奖
	CreateTime int64 `json:"create_time" gorm:"column:create_time"`
	EndTime    int64 `json:"end_time" gorm:"column:end_time"`
}

func (m *UserDrawRecord) TableName() string {
	return "user_draw_record"
}

func initUserDrawRecord(tx *gorm.DB) {

}
