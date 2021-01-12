package models

import (
	"time"
)

var (
	// 通知成功
	JobSuccess = 1
	// 待通知
	JobWait = 2
	// 通知失败
	JobFail = 3
)

type Job struct {
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
