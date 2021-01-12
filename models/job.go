package models

import (
	"gorm.io/gorm"
	"time"
)

var (
	// 待通知
	JobWait = 2
	// 通知成功
	JobSuccess = 1
)

type Job struct {
	gorm.Model
	Id         int64
	Content    string
	CreatedAt  time.Time
	NoticeTime time.Time
	Status     int8
	Phone      string
	Email      string
}

func (Job) TableName() string {
	return "jobs"
}
