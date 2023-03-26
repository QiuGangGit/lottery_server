package svc

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"lottery_server/internal/config"
	"lottery_server/model"
)

func NewServiceContext(config config.Config) *ServiceContext {
	s := &ServiceContext{Config: config}
	s.Sqlite()
	return s
}

type ServiceContext struct {
	Config config.Config
	sqlite *gorm.DB
}

func (s *ServiceContext) Sqlite() *gorm.DB {
	if s.sqlite == nil {
		db, err := gorm.Open(sqlite.Open(s.Config.SqlitePath), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db.AutoMigrate(&model.User{}, &model.Prize{}, &model.DrawRecord{}, &model.UserDrawRecord{})
		model.Init(db)
		s.sqlite = db
	}
	return s.sqlite
}
