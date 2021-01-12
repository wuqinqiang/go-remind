package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Wechat struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
}

type Db struct {
	Address  string
	DbName   string
	User     string
	Password string
	Port     int
}

type Email struct {
	User string
	Pass string
	Host string
	Port int
}
type Teng struct {
	SECRETID   string
	SecretKey  string
	SDKAppID   string
	TemplateID string
}

type Configuration struct {
	Wechat *Wechat
	Db     *Db
	Email  *Email
	Teng   *Teng
}

var once sync.Once

var All *Configuration

func LoadConfig() {
	once.Do(func() {
		file, err := os.Open("config.json")
		if err != nil {
			log.Fatalln("can not open the file")
		}
		decoder := json.NewDecoder(file)
		All = &Configuration{}
		err = decoder.Decode(All)
		if err != nil {
			log.Fatalln("Cannot get configuration from file", err)
		}
	})
}
