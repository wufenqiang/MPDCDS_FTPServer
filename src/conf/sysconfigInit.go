package conf

import (
	"github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var Sysconfig = &sysconfig{}

func init() {
	dir, _ := os.Getwd()
	for {
		if strings.HasSuffix(dir, ProjectName) {
			pdir := filepath.Dir(dir)
			if pdir != "thrift" {
				break
			}
		}
		dir = filepath.Dir(dir)
	}

	//fmt.Println(dir)

	conffile := dir + "/config.json"

	//指定对应的json配置文件
	b, err := ioutil.ReadFile(conffile)
	if err != nil {
		panic("Sys config read err")
	}
	err = jsoniter.Unmarshal(b, Sysconfig)
	if err != nil {
		panic(err)
	}
}
