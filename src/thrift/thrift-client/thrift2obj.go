package thrift_client

import (
	"gitlab.weather.com.cn/wufenqiang/MPDCDSPro/src/thrift/thriftcore"
)

type FileReturn struct {
	*thriftcore.FileReturn
}

func (fileReturn FileReturn) FileInfo2Status() int16 {
	return fileReturn.Status
}
func (fileReturn FileReturn) FileInfo2Msg() string {
	return fileReturn.Msg
}
func (fileReturn FileReturn) FileInfo2Data() map[string]string {
	return fileReturn.Data
}
func (fileReturn FileReturn) FileInfo2FileAddress() string {
	abspath := fileReturn.Data["file_address"]
	return abspath
}
func (fileReturn FileReturn) FileInfo2FileID() string {
	fileid := fileReturn.Data["file_id"]
	return fileid
}
func (fileReturn FileReturn) FileInfo2AccessId() string {
	fileid := fileReturn.Data["access_id"]
	return fileid
}

type AuthReturn struct {
	*thriftcore.AuthReturn
}

func (auth AuthReturn) Auth2Status() int16 {
	return auth.Status
}
func (auth AuthReturn) Auth2Token() string {
	return auth.Token
}
func (auth AuthReturn) Auth2Msg() string {
	return auth.Msg
}

type DirAuthReturn struct {
	*thriftcore.DirAuthReturn
}

func (dirAuthReturn DirAuthReturn) DirAuth2Status() int16 {
	return dirAuthReturn.Status
}
func (dirAuthReturn DirAuthReturn) DirAuth2Msg() string {
	return dirAuthReturn.Msg
}

type ListsReturn struct {
	*thriftcore.ListsReturn
}

func (listsReturn ListsReturn) FileDirInfo2Status() int16 {
	return listsReturn.Status
}
func (listsReturn ListsReturn) FileDirInfo2Msg() string {
	return listsReturn.Msg
}
func (listsReturn ListsReturn) FileDirInfo2Data() []map[string]string {
	return listsReturn.Data
}
