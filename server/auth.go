package server

import (
	"MPDCDS_FTPServer/thrift/MPDCDS_BackendService"
	"MPDCDS_FTPServer/thrift/client"
	"context"
)

// Auth is an interface to auth your ftp user login.
type Auth interface {
	CheckPasswd(string, string) (*MPDCDS_BackendService.Auth, error)
}

var (
	_Auth = &SimpleAuth{}
)

// SimpleAuth implements Auth interface to provide a memory user login auth
type SimpleAuth struct {
	//User     string
	//Password string
}

// CheckPasswd will check user's password
func (a *SimpleAuth) CheckPasswd(user string, password string) (*MPDCDS_BackendService.Auth, error) {

	tClient, tTransport := client.Connect()
	ctx := context.Background()
	auth, err := tClient.Auth(ctx, user, password)
	//关闭tTransport
	client.Close(tTransport)

	return auth, err
}
