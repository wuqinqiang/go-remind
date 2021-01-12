package db

import (
	"database/sql"
	"fmt"
	"go-remind/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	Url            = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s"
	MaxOpen        = 5 // 最大打开数
	MaxIdle        = 2 // 最大保留连接数
	LifeMinuteTime = 5 // 连接可重用最大时间
)

var Gorm *gorm.DB

func InitDb(conf *config.Db) error {
	var err error
	var sqlDb *sql.DB
	Gorm, err = gorm.Open(mysql.Open(
		fmt.Sprintf(
			Url, conf.User, conf.Password,
			conf.Address, conf.Port, conf.DbName)), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDb, err = Gorm.DB()

	if err != nil {
		return err
	}
	sqlDb.SetMaxOpenConns(MaxOpen)
	sqlDb.SetConnMaxIdleTime(MaxIdle)
	sqlDb.SetConnMaxLifetime(LifeMinuteTime * time.Minute)
	return nil
}
