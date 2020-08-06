package protocol_stream

func (pf *ProtocolFactory) GetData_https() ReturnType {
	s, f, e := HttpGet(pf.path)
	return ReturnType{s, f, e}
}
