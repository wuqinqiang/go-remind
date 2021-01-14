package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/cache"
	"github.com/silenceper/wechat/v2"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	. "go-remind/config"
	"go-remind/logic"
	"go-remind/tools"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var RequestFormatErr = errors.New("参数错误")

var timeDay = map[string]string{
	"今天":  getDateString(0),
	"明天":  getDateString(1),
	"后天":  getDateString(2),
	"大后天": getDateString(3),
}

var timeHMS = map[string]bool{
	"点": true,
	"分": true,
}

type contentRegexp struct {
	*regexp.Regexp
}

func Message(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("运行错误:%v", err)
		}
	}()
	//使用 memcache 保存access_token，也可选择redis或自定义cache
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:          ConfAll.Wechat.AppID,
		AppSecret:      ConfAll.Wechat.AppSecret,
		Token:          ConfAll.Wechat.Token,
		EncodingAESKey: ConfAll.Wechat.EncodingAESKey,
		Cache:          memory,
	}

	officialAccount := wc.GetOfficialAccount(cfg)
	// 传入request和responseWriter
	server := officialAccount.GetServer(c.Request, c.Writer)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
		switch msg.MsgType {
		case message.MsgTypeText:
			//回复消息：演示回复用户发送的消息
			res, _ := HandleMessage(msg.Content)
			text := message.NewText(res)
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		case message.MsgTypeVoice:
			text := message.NewVoice(msg.Content)
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		default:
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("我睡着了，听不懂你在说啥")}
		}
	})
	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()

}

//计算日期
func getDateString(count int) string {
	t := time.Now()
	newTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	//通知时间
	noticeTime := newTime.AddDate(0, 0, count)
	logDay := noticeTime.Format("2006-01-02")
	return logDay
}

//时间匹配
var myexp = contentRegexp{regexp.MustCompile(
	`(今天|明天|后天|大后天|[\d]{4}-[\d]{2}-[\d]{2}\s[\d]{2}:[\d]{2}|[\d]{8}\s[\d]{1,2}:[\d]{1,2}|[[\d]{1,2}:[\d]{1,2}|[\d]{1,2}(个月|小时|点|分钟|分|秒|周|天))`,
)}

//手机号匹配
var phone = contentRegexp{regexp.MustCompile(
	`(1[356789]\d)(\d{4})(\d{4})`,
)}

func HandleMessage(content string) (string, error) {
	phone := phone.FindStringSubmatch(content)
	if phone == nil {
		return "不留下联系方式我咋么联系上你", RequestFormatErr
	}
	mmp := myexp.FindAllStringSubmatch(content, -1)
	fmt.Println(mmp)
	if mmp == nil {
		return "我得再升升级才能满足你的时间格式", RequestFormatErr
	}
	//最多只有三位 时 分 秒
	if len(mmp) > 3 {
		mmp = mmp[:3]
	}
	var sendDate string
	for _, item := range mmp {
		//今天明天大后天
		if _, ok := timeDay[item[0]]; ok {
			sendDate = timeDay[item[0]]
			continue
		}
		//本身日期格式 2020-05-20 13:00
		if sendDate == "" {
			sendDate = item[0]
		} else {
			lateTime := item[0]
			//19点20分
			if timeHMS[lateTime[utf8.RuneCountInString(lateTime)-1:]] {
				var numberTime string = lateTime[0 : utf8.RuneCountInString(lateTime)-1]
				if lateTime[utf8.RuneCountInString(lateTime)-1:] == "分" {
					sendDate = strings.Replace(sendDate, ":00", ":"+numberTime, 1)
					continue
				}
				sendDate += " " + numberTime + ":00"
				continue
			}
			sendDate = sendDate + " " + lateTime
		}
	}
	fmt.Printf("发送的时间：%v", sendDate)

	sendTimer := tools.StringToTimer(sendDate)
	diff := sendTimer.Sub(tools.GetCurrTime())
	if diff < 0 {
		return "过期的时间就别让我通知了", RequestFormatErr
	}

	jobLogic := &logic.JobLogic{}
	job := logic.NewJob(content, sendTimer, phone[0], "")
	err := jobLogic.Insert(job)
	if err != nil {
		return "请检查", RequestFormatErr
	}
	fmt.Printf("插入的值是：%v\n", job.Id)
	if diff.Minutes() < 0 {
		return fmt.Sprintf("%s秒后短信提醒内容:%s", tools.Decimal(diff.Seconds()), content), nil
	}
	if diff.Hours() < 1 {
		go func() {
			HandleNotice(job)
		}()
		return fmt.Sprintf("%s分钟后短信提醒内容:%s", tools.Decimal(diff.Minutes()), content), nil
	}
	return fmt.Sprintf("%s小时后短信提醒内容:%s", tools.Decimal(diff.Hours()), content), nil
}
