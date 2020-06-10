// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// This is a very simple ftpd server using this library as an example
// and as something to run tests against.
package main

import (
	"MPDCDS_FTPServer/conf"
	"MPDCDS_FTPServer/logger"
	"MPDCDS_FTPServer/utils"
	"flag"
	"log"
	"strconv"

	"MPDCDS_FTPServer/file-driver"
	"MPDCDS_FTPServer/server"
)

func main() {
	ip, _ := utils.ExternalIP()

	var (
		//root = flag.String("root", "/tmp", "Root directory to serve")
		port = flag.Int("port", conf.Sysconfig.FTPCmdPort, "Port")
		//host = flag.String("host", conf.Sysconfig.FTPHost, "Host")
		host = flag.String("host", ip.String(), "Host")
	)
	flag.Parse()

	factory := &filedriver.FileDriverFactory{
		Perm: server.NewSimplePerm("user", "group"),
	}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     *port,
		Hostname: *host,
		Auth:     &server.SimpleAuth{},
	}

	logger.GetLogger().Info("Starting ftp server on " + opts.Hostname + ":" + strconv.Itoa(opts.Port))
	//log.Printf("Username %v, Password %v", *user, *pass)
	server := server.NewServer(opts)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
