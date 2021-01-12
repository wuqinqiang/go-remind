package db

import (
	"fmt"
	"go-remind/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	Url            = "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&parseTime=True&timeout=5s"
	MaxOpen        = 10 // 最大打开数
	MaxIdle        = 2  // 最大保留连接数
	LifeMinuteTime = 5  // 连接可重用最大时间
)

var Gorm *gorm.DB

func InitDb(conf *config.Db) {
	Gorm, err := gorm.Open(mysql.Open(
		fmt.Sprintf(
			Url, conf.User, conf.Password,
			conf.Address, conf.Port, conf.DbName)), &gorm.Config{})
	if err != nil {
		fmt.Printf("open db:%v", err)
	}
	sqlDb, err := Gorm.DB()
	if err != nil {
		fmt.Printf("sql Db:%v", err)
	}
	sqlDb.SetMaxOpenConns(MaxOpen)
	sqlDb.SetConnMaxIdleTime(MaxIdle)
	sqlDb.SetConnMaxLifetime(LifeMinuteTime * time.Minute)
}
