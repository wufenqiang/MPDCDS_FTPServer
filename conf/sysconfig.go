package conf

import (
	"github.com/json-iterator/go"
	"io/ioutil"
)

var Sysconfig = &sysconfig{}

func init() {
	//指定对应的json配置文件
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("Sys config read err")
	}
	err = jsoniter.Unmarshal(b, Sysconfig)
	if err != nil {
		panic(err)
	}
}

type sysconfig struct {

	//thrift 服务ip
	NetworkAddr string `json:"NetworkAddr"`
	ThriftPort  string `json:"ThriftPort"`

	//日志存储地址、级别
	LoggerPath  string `json:"LoggerPath"`
	LoggerLevel string `json:"LoggerLevel"`
}
