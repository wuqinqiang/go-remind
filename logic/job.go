package logic

import (
	"go-remind/db"
	"go-remind/models"
	"time"
)

type JobLogic struct{}

func NewJob(content string, sendTime time.Time, phone, email string) models.Job {
	return models.Job{
		Content:    content,
		NoticeTime: sendTime,
		Phone:      phone,
		Email:      email,
		Status:     models.JobWait,
	}
}

func (j *JobLogic) Insert(job models.Job) error {
	result := db.Gorm.Create(&job)
	return result.Error
}

func (j *JobLogic) GetJobsByTime(startTime string, endTime string) (jobs []models.Job, err error) {
	err = db.Gorm.Where("status=? and notice_time>=? and notice_time<=?", models.JobWait, startTime, endTime).
		Find(&jobs).Error
	return
}

func (j *JobLogic) UpdateStatusById(id, status int) error {
	return db.Gorm.Model(&models.Job{}).Where("id=?", id).Update("status", status).Error
}
