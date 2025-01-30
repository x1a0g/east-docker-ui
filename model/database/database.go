package database

import (
	v1 "east-docker-ui/model/database/v1"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

var Base string

func init() {
	Base = v1.LV1()

	DB, err = gorm.Open(sqlite.Open(Base), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	logger, _ := zap.NewDevelopment()
	logger.Info("Database connected")
}
