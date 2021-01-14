package tools

import (
	"strconv"
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海

func GetCurrTime() time.Time {
	return time.Now().In(cstSh)
}

func TimeString(time time.Time) string {
	return time.Format(TimeFormat)
}

func StringToTimer(date string) time.Time {
	res, _ := time.ParseInLocation("2006-01-02 15:04:05", date, cstSh)
	return res
}

//浮点数转字符串截取
func Decimal(value float64) string {
	string := strconv.FormatFloat(value, 'f', 6, 64)
	return string[:4]
}
