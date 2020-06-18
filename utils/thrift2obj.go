package utils

import (
	"MPDCDS_FTPServer/thrift/MPDCDS_BackendService"
	"time"
)

func FileInfo2AbsPath(fileInfo *MPDCDS_BackendService.FileInfo) string {
	abspath := fileInfo.Data["file_address"]
	return abspath
}

func FileInfo2FileID(fileInfo *MPDCDS_BackendService.FileInfo) string {
	fileid := fileInfo.Data["file_id"]
	return fileid
}

func FileInfo2AccessId(fileInfo *MPDCDS_BackendService.FileInfo) string {
	fileid := fileInfo.Data["access_id"]
	return fileid
}

func Now2TimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
