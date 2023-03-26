package model

import "gorm.io/gorm"

// Prize 抽奖池
type Prize struct {
	Id    int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name  string `json:"name" gorm:"column:name"`
	Desc  string `json:"desc" gorm:"column:desc"`
	Count int    `json:"count" gorm:"column:count"`
	Icon  string `json:"icon" gorm:"column:icon"`
}

func (m *Prize) TableName() string {
	return "prize"
}

func initPrize(tx *gorm.DB) {
	tx.Create(&Prize{
		Name:  "iPhone14ProMax",
		Count: 100,
		Icon:  "https://xxim-1312910328.cos.ap-guangzhou.myqcloud.com/tmp/iphone.png",
		Desc:  "美国科技公司苹果公司于2020年10月发布的智能手机，是iPhone 14系列的一员。作为苹果公司的旗舰机型，iPhone 14 Pro Max拥有许多出色的特性和功能",
	})
	tx.Create(&Prize{
		Name:  "iPad",
		Count: 100,
		Icon:  "https://xxim-1312910328.cos.ap-guangzhou.myqcloud.com/tmp/ipad.png",
		Desc:  "iPad是一款由美国科技公司苹果公司推出的平板电脑，于2010年首次发布。它采用了iOS操作系统，与iPhone等苹果公司的智能手机共享相似的用户界面和应用程序生态系统，但其屏幕更大、处理器更强大，适用于更多的生产力和娱乐应用。",
	})
	tx.Create(&Prize{
		Name:  "MacBookPro",
		Count: 100,
		Icon:  "https://xxim-1312910328.cos.ap-guangzhou.myqcloud.com/tmp/mbp.png",
		Desc:  "MacBook Pro是苹果公司的高端笔记本电脑系列，具备强大的性能、高清的屏幕和长久的电池续航。",
	})
	tx.Create(&Prize{
		Name:  "AirPods",
		Count: 100,
		Icon:  "https://xxim-1312910328.cos.ap-guangzhou.myqcloud.com/tmp/airpods.png",
		Desc:  "AirPods是苹果公司推出的无线蓝牙耳机，具有自动连接、智能感应、高质量音效和长久电池续航等特点。",
	})
	tx.Create(&Prize{
		Name:  "AppleWatch",
		Count: 100,
		Icon:  "https://xxim-1312910328.cos.ap-guangzhou.myqcloud.com/tmp/watch.png",
		Desc:  "Apple Watch是苹果公司推出的智能手表，能够监测健康状况、提醒日程、支持支付、内置音乐播放器等。",
	})
}
