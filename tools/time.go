package tools

import "time"

const (
	TimeFormat = "2006-01-02 15:04:05"
)

func GetCurrTime() string {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海

	return time.Now().In(cstSh).Format(TimeFormat)
}

func GetCurrentDay() string {
	return time.Now().Format("2006-01-02")
}
