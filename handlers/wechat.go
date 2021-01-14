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
	"strings"
	"unicode/utf8"
)

var RequestFormatErr = errors.New("参数错误")

func Message(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("运行错误:%v", err)
		}
	}()
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
			//解析用户消息，恢复用户消息
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
	_ = server.Send()

}

func HandleMessage(content string) (string, error) {
	phone := tools.PhoneMatch.FindStringSubmatch(content)
	email := tools.EmailMatch.FindStringSubmatch(content)
	if phone == nil || email == nil {
		return "不留下联系方式我咋么联系上你", RequestFormatErr
	}
	mmp := tools.TimeMatch.FindAllStringSubmatch(content, -1)
	if mmp == nil {
		return "我得再升升级才能满足你的时间格式", RequestFormatErr
	}

	// 最大匹配到分
	if len(mmp) > 3 {
		mmp = mmp[:3]
	}
	var sendDate string
	for _, item := range mmp {
		//今天明天大后天
		if _, ok := tools.TimeDay[item[0]]; ok {
			sendDate = tools.TimeDay[item[0]]
			continue
		}
		//本身日期格式 2020-05-20 13:00
		if sendDate == "" {
			sendDate = item[0] + ""
		} else {
			lateTime := item[0]
			//19点20分
			if tools.TimeHMS[lateTime[utf8.RuneCountInString(lateTime)-1:]] {
				numberTime := lateTime[0 : utf8.RuneCountInString(lateTime)-1]
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

	sendTimer := tools.StringToTimer(sendDate + ":00")
	diff := sendTimer.Sub(tools.GetCurrTime())
	if diff < 0 {
		return "过期的时间就别让我通知了", RequestFormatErr
	}

	jobLogic := &logic.JobLogic{}
	job := logic.NewJob(content, sendTimer, phone[0], email[0])
	err := jobLogic.Insert(job)
	if err != nil {
		return "请检查输入内容", RequestFormatErr
	}

	if diff.Minutes() < 0 {
		return fmt.Sprintf("%s秒后短信提醒内容:%s", tools.Decimal(diff.Seconds()), content), nil
	}

	if diff.Hours() < 1 {
		//小于1个小时直接加入到定时器
		go func() {
			HandleNotice(job)
		}()
		return fmt.Sprintf("%s分钟后短信提醒内容:%s", tools.Decimal(diff.Minutes()), content), nil
	}
	return fmt.Sprintf("%s小时后短信提醒内容:%s", tools.Decimal(diff.Hours()), content), nil
}
