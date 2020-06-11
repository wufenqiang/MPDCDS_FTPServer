package filedriver

import "os"

func (pf *ProtocolFactory) GetData_file() ReturnType {
	path, e0 := pf.thePath()
	if e0 != nil {
		return ReturnType{0, nil, e0}
	}

	f, e1 := os.Open(path)
	if e1 != nil {
		return ReturnType{0, nil, e1}
	}

	info, e2 := f.Stat()
	if e2 != nil {
		return ReturnType{0, f, e2}
	}

	s := info.Size()

	return ReturnType{s, f, nil}
}
