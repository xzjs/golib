package golib

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

//IConf .
type IConf interface {
	GetConf(section string, key string) (value string)
}

//Confer .
type Confer struct {
	Cfg *ini.File
}

var conf IConf

//Conf .
func Conf() IConf {
	if conf == nil {
		c, err := ini.Load("conf.ini")
		if err != nil {
			fmt.Println("failed to read conf.ini ", err)
			os.Exit(1)
		}
		conf = &Confer{
			Cfg: c,
		}

	}
	return conf
}

//SetConf 设置单利
func SetConf(c IConf) {
	conf = c
}

// GetConf 获取配置值
func (c *Confer) GetConf(section string, key string) (value string) {
	return c.Cfg.Section(section).Key(key).String()
}
