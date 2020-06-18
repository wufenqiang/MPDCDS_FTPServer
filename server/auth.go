package server

import (
	"MPDCDS_FTPServer/thrift/client"
	"MPDCDS_FTPServer/utils"
	"context"
)

// Auth is an interface to auth your ftp user login.
type Auth interface {
	CheckPasswd(string, string) (int16, string, error, string)
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
func (a *SimpleAuth) CheckPasswd(user string, password string) (int16, string, error, string) {

	tClient, tTransport := client.Connect()
	ctx := context.Background()
	auth, err := tClient.Auth(ctx, user, password)
	//关闭tTransport
	client.Close(tTransport)

	auth0 := utils.Auth{auth}
	status := auth0.Auth2Status()
	msg := auth0.Auth2Msg()
	token := auth0.Auth2Token()

	return status, token, err, msg
}
