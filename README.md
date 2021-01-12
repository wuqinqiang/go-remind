## 微信公众号提醒⏰



### 环境
**go环境自行安装**

### 增加配置
根目录下增加 config.json 文件 
```json
{
  "Wechat": {
    "AppID": "xxx",
    "AppSecret": "xxx",
    "Token": "xxx",
    "EncodingAESKey": "xxxxx"
  },
  "Db": {
    "Address": "127.0.0.1",
    "DbName": "remind",
    "User": "root",
    "Password": "Passw0rd",
    "Port": 3306
  },
  "Email": {
    "User": "1185079673@qq.com",
    "Pass": "xxx",
    "Host": "smtp.qq.com",
    "Port": 25
  }
}
```

**数据库**

**数据库文件放在根目录下**


### 运行
**在项目根目录下运行**
```go
go run main.go
```







