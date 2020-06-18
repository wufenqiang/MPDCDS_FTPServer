package utils

import (
	"MPDCDS_FTPServer/thrift/MPDCDS_BackendService"
	"time"
)

func Now2TimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

type FileInfo struct {
	*MPDCDS_BackendService.FileInfo
}

func (fileInfo FileInfo) FileInfo2Status() int16 {
	return fileInfo.Status
}
func (fileInfo FileInfo) FileInfo2Msg() string {
	return fileInfo.Msg
}
func (fileInfo FileInfo) FileInfo2Data() map[string]string {
	return fileInfo.Data
}
func (fileInfo FileInfo) FileInfo2FileAddress() string {
	abspath := fileInfo.Data["file_address"]
	return abspath
}
func (fileInfo FileInfo) FileInfo2FileID() string {
	fileid := fileInfo.Data["file_id"]
	return fileid
}
func (fileInfo FileInfo) FileInfo2AccessId() string {
	fileid := fileInfo.Data["access_id"]
	return fileid
}

type Auth struct {
	*MPDCDS_BackendService.Auth
}

func (auth Auth) Auth2Status() int16 {
	return auth.Status
}
func (auth Auth) Auth2Token() string {
	return auth.Token
}
func (auth Auth) Auth2Msg() string {
	return auth.Msg
}

type DirAuth struct {
	*MPDCDS_BackendService.DirAuth
}

func (dirAuth DirAuth) DirAuth2Status() int16 {
	return dirAuth.Status
}
func (dirAuth DirAuth) DirAuth2Msg() string {
	return dirAuth.Msg
}

type FileDirInfo struct {
	*MPDCDS_BackendService.FileDirInfo
}

func (fileDirInfo FileDirInfo) FileDirInfo2Status() int16 {
	return fileDirInfo.Status
}
func (fileDirInfo FileDirInfo) FileDirInfo2Msg() string {
	return fileDirInfo.Msg
}
func (fileDirInfo FileDirInfo) FileDirInfo2Data() []map[string]string {
	return fileDirInfo.Data
}
