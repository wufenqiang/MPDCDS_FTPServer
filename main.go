// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// This is a very simple ftpd server using this library as an example
// and as something to run tests against.
package main

import (
	"MPDCDS_FTPServer/src/conf"
	"MPDCDS_FTPServer/src/file-driver"
	"MPDCDS_FTPServer/src/ftp-server"
	"MPDCDS_FTPServer/src/logger"
	"MPDCDS_FTPServer/src/service"
	"flag"
)

func main() {

	//ftpip, _ := utils.ExternalIP()
	var (
		//root = flag.String("root", "/tmp", "Root directory to serve")
		port = flag.Int("port", conf.Sysconfig.FTPCmdPort, "Port")
		//host = flag.String("host", conf.Sysconfig.FTPHost, "Host")
		//host = flag.String("host", ftpip.String(), "Host")
		host = flag.String("host", "0.0.0.0", "Host")
	)
	flag.Parse()

	factory := &file_driver.FileDriverFactory{
		Perm: ftp_server.NewSimplePerm("user", "group"),
	}

	opts := &ftp_server.ServerOpts{
		Factory:  factory,
		Port:     *port,
		Hostname: *host,
		Auth:     &service.SimpleAuth{},
	}

	server := ftp_server.NewServer(opts)
	err := server.ListenAndServe()
	if err != nil {
		logger.GetLogger().Fatal("Error starting server:" + err.Error())
	}
}

//func init()  {
//	thrifthostport:=net.JoinHostPort(conf.Sysconfig.ThriftHost, conf.Sysconfig.ThriftPort)
//	thrift_client.ConnectHostPort(thrifthostport)
//}
