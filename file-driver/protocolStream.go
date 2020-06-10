package filedriver

import (
	"MPDCDS_FTPServer/logger"
	"os"
	"reflect"

	"strings"
)

func ReadFile(path string) (*os.File, error) {
	pf := ProtocolFactory{path}
	return pf.getFile()
}

type ProtocolFactory struct {
	path string
}

func (pf *ProtocolFactory) head() string {
	head := strings.Split(pf.path, "://")[0]
	if head == "" {
		logger.GetLogger().Error(pf.path + "没有协议头")
	}
	return strings.ToLower(head)
}
func (pf *ProtocolFactory) getFile() (*os.File, error) {
	head := pf.head()
	rf := ToReturnType(CallMethod(pf, "GetFile_"+head))
	return rf.getF(), rf.getE()
}

type ReturnType struct {
	f *os.File
	e error
}

func (rt ReturnType) getF() *os.File {
	return rt.f
}
func (rt ReturnType) getE() error {
	return rt.e
}

func ToReturnType(i interface{}) ReturnType {
	switch i.(type) {
	case ReturnType:
		return i.(ReturnType)
	default:
		logger.GetLogger().Error("ToReturnType类型无法匹配")
		return ReturnType{nil, nil}
	}
}

func CallMethod(i interface{}, methodName string) interface{} {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(i)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(i))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// check for method on value
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// check for method on pointer
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	if finalMethod.IsValid() {
		return finalMethod.Call([]reflect.Value{})[0].Interface()
	}

	// return or panic, method not found of either type
	return ReturnType{nil, nil}
}
func (pf *ProtocolFactory) GetFile_file() ReturnType {
	f, e := os.Open(pf.path)
	return ReturnType{f, e}
}
func (pf *ProtocolFactory) GetFile_http() ReturnType {
	return ReturnType{nil, nil}
}
