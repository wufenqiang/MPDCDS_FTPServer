package server

import (
	"crypto/subtle"
)

// Auth is an interface to auth your ftp user login.
type Auth interface {
	CheckPasswd(string, string) (bool, error)
}

var (
	_ Auth = &SimpleAuth{}
)

// SimpleAuth implements Auth interface to provide a memory user login auth
type SimpleAuth struct {
	//User     string
	//Password string
}

// CheckPasswd will check user's password
func (a *SimpleAuth) CheckPasswd(user, password string) (bool, error) {

	/**
	此处需要修改调用API的认证接口,待API开发
	*/
	var user0 = "admin"
	var password0 = "123456"

	var flag = constantTimeEquals(user, user0) && constantTimeEquals(password, password0)

	return flag, nil
}

func constantTimeEquals(a, b string) bool {
	return len(a) == len(b) && subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
