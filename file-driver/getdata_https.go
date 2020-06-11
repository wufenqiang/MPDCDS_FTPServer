package filedriver

import "MPDCDS_FTPServer/utils"

func (pf *ProtocolFactory) GetData_https() ReturnType {
	f, e := utils.HttpGet(pf.path)
	return ReturnType{f, e}
}
