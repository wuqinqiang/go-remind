package server

import (
	"fmt"
	. "go-remind/config"
	"go-remind/models"
	"go-remind/tools/try"
	"gopkg.in/gomail.v2"
)

type Message interface {
	SendMessage() error
}

type EmailMsg struct {
	Job *models.Job
}

func (email *EmailMsg) SendMessage() error {
	sendMail := gomail.NewMessage()
	sendMail.SetHeader(`From`, ConfAll.Email.User)
	sendMail.SetHeader(`To`, "1185079673@qq.com")
	sendMail.SetHeader(`Subject`, "来自吴亲库里的温馨提醒")
	sendMail.SetBody(`text/html`, email.Job.Content)
	err := gomail.NewDialer(
		ConfAll.Email.Host, ConfAll.Email.Port, ConfAll.Email.User,
		ConfAll.Email.Pass).DialAndSend(sendMail)
	if err != nil {
		return err
	}
	return nil
}

type SmsMsg struct {
	Job models.Job
}

func (email *SmsMsg) SendMessage() error {
	fmt.Printf("发送短信：%v\n", email.Job)
	return nil
}

func Notice(msg Message) error {
	return try.Do(func(attempt int) (retry bool, err error) {
		err = msg.SendMessage()
		if err != nil {
			return attempt < try.MaxRetries, err
		}
		return true, nil
	})
}
