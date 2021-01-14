package handlers

import (
	"fmt"
	"go-remind/logic"
	"go-remind/models"
	"go-remind/server"
	"go-remind/tools"
	"time"
)

var isFirst bool = true

func Scheduler() {
	var job logic.JobLogic
	for {
		if !isFirst {
			timer := time.NewTicker(1 * time.Hour)
			<-timer.C
		}
		isFirst = false

		// 获取接下来一小时内需要发送的任务列表
		now := tools.GetCurrTime()
		h, _ := time.ParseDuration("1h")
		jobs, err := job.GetJobsByTime(tools.TimeString(now), tools.TimeString(now.Add(1*h)))
		if err != nil {
			fmt.Printf("出错了：%v", err)
			return
		}
		// 任务通道
		ch := make(chan models.Job, 100)

		handleJob := func(ch <-chan models.Job) {
			for item := range ch {
				// 发送通知
				go HandleNotice(&item)
			}
		}
		// 处理任务
		go handleJob(ch)

		// 投递任务
		for _, job := range jobs {
			ch <- job
		}
	}
}

func HandleNotice(job *models.Job) {
	now := tools.GetCurrTime()
	noticeTime, _ := time.ParseInLocation(tools.TimeFormat,
		job.NoticeTime.Format(tools.TimeFormat), time.Local)
	diff := noticeTime.Sub(now)
	timer := time.NewTimer(diff)
	<-timer.C

	email := &server.EmailMsg{Job: job}
	err := server.Notice(email)
	//成功与否
	isOk := models.JobSuccess
	jobLogic := logic.JobLogic{}

	if err != nil {
		isOk = models.JobFail
		fmt.Printf("通知失败:%v", err)
	}
	_ = jobLogic.UpdateStatusById(job.Id, isOk)

}
