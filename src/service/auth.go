package service

import (
	"MPDCDS_FTPServer/src/thrift/thrift-client"
	"context"
	"gitlab.weather.com.cn/wufenqiang/MPDCDSPro/src/thrift/thriftcore"
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

	ctx := context.Background()
	authInfo := thriftcore.NewAuthInfo()
	authInfo.User = user
	authInfo.Password = password
	auth, err := thrift_client.ThriftClient.Auth(ctx, authInfo)
	//关闭tTransport
	thrift_client.ThriftClose()

	auth0 := thrift_client.AuthReturn{auth}
	status := auth0.Auth2Status()
	msg := auth0.Auth2Msg()
	token := auth0.Auth2Token()

	return status, token, err, msg
}
