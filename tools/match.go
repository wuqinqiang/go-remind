package tools

import (
	"regexp"
	"time"
)

var TimeDay = map[string]string{
	"今天":  getDateString(0),
	"明天":  getDateString(1),
	"后天":  getDateString(2),
	"大后天": getDateString(3),
}

var TimeHMS = map[string]bool{
	"点": true,
	"分": true,
}

//时间匹配
var TimeMatch = contentRegexp{regexp.MustCompile(
	`(今天|明天|后天|大后天|[\d]{4}-[\d]{2}-[\d]{2}\s[\d]{2}:[\d]{2}|[\d]{8}\s[\d]{1,2}:[\d]{1,2}|[[\d]{1,2}:[\d]{1,2}|[\d]{1,2}(个月|小时|点|分钟|分|秒|周|天))`,
)}

//手机号匹配
var PhoneMatch = contentRegexp{regexp.MustCompile(
	`(1[356789]\d)(\d{4})(\d{4})`,
)}

//邮箱匹配
var EmailMatch = contentRegexp{regexp.MustCompile(
	`(\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*)`,
)}


//计算日期
func getDateString(count int) string {
	t := time.Now()
	newTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	//通知时间
	noticeTime := newTime.AddDate(0, 0, count)
	logDay := noticeTime.Format("2006-01-02")
	return logDay
}


type contentRegexp struct {
	*regexp.Regexp
}
