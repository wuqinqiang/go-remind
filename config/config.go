package config

import (
	"encoding/json"
	"os"
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

var ConfAll *Configuration

func LoadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	ConfAll = &Configuration{}
	err = decoder.Decode(ConfAll)
	if err != nil {
		return err
	}
	return nil
}
