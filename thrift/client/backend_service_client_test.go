package client

import (
	"MPDCDS_FTPServer/logger"
	"MPDCDS_FTPServer/thrift/MPDCDS_BackendService"
	"context"
	"fmt"
	"testing"
	"time"
)

/**需先启动MPDCDS_
 */

func TestAuth(t *testing.T) {
	tClient, tTransport := Connect()
	user := "tName1"
	password := "123456"
	res, err := tClient.Auth(context.Background(), user, password)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
	logger.GetLogger().Info(fmt.Sprintf(res.String()))
	tTransport.Close()
}

func TestLists(t *testing.T) {
	tClient, tTransport := Connect()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImEwNzE3ZmIwLTQ3YmQtNDdlNy1iMmJmLWFlN2RlMjM2MjhhYyIsInVzZXJuYW1lIjoidE5hbWUxIn0.LKoBOQkfc6_XtGrIPRAWgwUAkD1Zim7ltEzzdN5F0mQ"
	pwd := "/"
	res, err := tClient.Lists(context.Background(), token, pwd)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
	logger.GetLogger().Info(fmt.Sprintf(res.String()))
	tTransport.Close()
}

func TestDirAuth(t *testing.T) {
	tClient, tTransport := Connect()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImEwNzE3ZmIwLTQ3YmQtNDdlNy1iMmJmLWFlN2RlMjM2MjhhYyIsInVzZXJuYW1lIjoidE5hbWUxIn0.LKoBOQkfc6_XtGrIPRAWgwUAkD1Zim7ltEzzdN5F0mQ"
	absPath := "/"
	res, err := tClient.DirAuth(context.Background(), token, absPath)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
	logger.GetLogger().Info(fmt.Sprintf(res.String()))
	tTransport.Close()
}

func TestFile(t *testing.T) {
	tClient, tTransport := Connect()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImEwNzE3ZmIwLTQ3YmQtNDdlNy1iMmJmLWFlN2RlMjM2MjhhYyIsInVzZXJuYW1lIjoidE5hbWUxIn0.LKoBOQkfc6_XtGrIPRAWgwUAkD1Zim7ltEzzdN5F0mQ"
	absPath := "/code_ocf_1h/ser/data/ocf/1h"
	fileName := "ocf1h-1.txt"
	res, err := tClient.File(context.Background(), token, absPath, fileName)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
	logger.GetLogger().Info(fmt.Sprintf(res.String()))
	tTransport.Close()
}

func TestSaveDownFileInfo(t *testing.T) {
	tClient, tTransport := Connect()
	apidown := MPDCDS_BackendService.NewApiDownLoad()
	apidown.StartTime = time.Time{}.Format("2006-01-02 15:04:05")
	apidown.FileID = "bfc51fb5-2a41-4c24-8c2c-c05efeb2384e"
	apidown.AccessID = "6814f474-8e23-4949-a5ca-8e0de71a1666"
	apidown.EndTime = time.Time{}.Format("2006-01-02 15:04:05")
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImEwNzE3ZmIwLTQ3YmQtNDdlNy1iMmJmLWFlN2RlMjM2MjhhYyIsInVzZXJuYW1lIjoidE5hbWUxIn0.LKoBOQkfc6_XtGrIPRAWgwUAkD1Zim7ltEzzdN5F0mQ"
	res, err := tClient.SaveDownLoadFileInfo(context.Background(), token, apidown)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
	logger.GetLogger().Info(fmt.Sprintf(res.String()))
	tTransport.Close()
}
