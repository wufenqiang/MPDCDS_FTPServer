package filedriver

import "os"

func (pf *ProtocolFactory) GetData_file() ReturnType {
	path := pf.thePath()
	f, e := os.Open(path)
	return ReturnType{f, e}
}
