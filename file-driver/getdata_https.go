package filedriver

import "MPDCDS_FTPServer/utils"

func (pf *ProtocolFactory) GetData_https() ReturnType {
	s, f, e := utils.HttpGet(pf.path)
	return ReturnType{s, f, e}
}
