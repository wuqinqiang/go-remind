package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/cache"
	"github.com/silenceper/wechat/v2"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	. "go-remind/config"
	"go-remind/logic"
	"go-remind/tools"
)

func Message(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("运行错误:%v", err)
		}
	}()
	//使用 memcache 保存access_token，也可选择redis或自定义cache
	wc := wechat.NewWechat()
	v := logic.JobLogic{}

	//m := models.Job{
	//	Content:    "666",
	//	NoticeTime: tools.GetCurrentDay(),
	//	Status:     2,
	//	Phone:      "131",
	//	Email:      "22",
	//}
	//v.Insert(m)

	list, err := v.GetJobsByTime("2020-01-12 20:41:58", "2021-01-13 20:41:58")
	for _, item := range list {
		//var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海

		fmt.Printf("时间:%v\n", item.NoticeTime.Format(tools.TimeFormat))
	}
	if err != nil {
		fmt.Printf("错误了:%v\n", err)
	}
	fmt.Printf("信息是:%v\n", list)

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
			res := message.NewText(msg.Content)
			//res := message.NewText(HandleMessage(msg.Content))
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: res}
		case message.MsgTypeVoice:
			text := message.NewVoice(msg.Content)
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		default:
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("我睡着了，听不懂你在说啥")}
		}
	})
	//处理消息接收以及回复
	err = server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()

}
