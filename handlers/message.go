package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/message"
	. "go-remind/config"
)

func Message(c *gin.Context) {
	fmt.Printf("数据：%v", All.Wechat)
	//配置微信
	config := &wechat.Config{
		AppID:          All.Wechat.AppID,
		AppSecret:      All.Wechat.AppSecret,
		Token:          All.Wechat.Token,
		EncodingAESKey: All.Wechat.EncodingAESKey,
	}
	wc := wechat.NewWechat(config)
	server := wc.GetServer(c.Request, c.Writer)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
		fmt.Printf("接收到数据:%v", message.MsgTypeEvent)
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
		return nil
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		fmt.Println(err.Error())
	}
}
