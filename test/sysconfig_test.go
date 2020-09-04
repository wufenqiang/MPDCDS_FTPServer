package test

import (
	"MPDCDS_FTPServer/src/conf"
	"fmt"
	"testing"
)

func TestSysconfig(t *testing.T) {

	fmt.Println(conf.Sysconfig.FTPCmdPort)
}
