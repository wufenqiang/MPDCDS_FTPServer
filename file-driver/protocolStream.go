package filedriver

import (
	"os"
	"reflect"

	"strings"
)

func ReadFile(path string) (*os.File, error) {
	pf := ProtocolFactory{path}
	return pf.getData()
}

type ReturnType struct {
	f *os.File
	e error
}
type ProtocolFactory struct {
	path string
}

const split = "://"

func (pf *ProtocolFactory) head() string {
	if strings.Contains(pf.path, split) {
		head := strings.Split(pf.path, split)[0]
		if head == "" {
			panic(pf.path + "没有找到协议头")
		}
		return strings.ToLower(head)
	} else {
		panic(pf.path + "没有找到协议头规格信息,eg:[file://][http://]")
	}
}
func (pf *ProtocolFactory) getData() (*os.File, error) {
	head := pf.head()
	//动态调用接口
	i := pf.CallMethod(pf, "GetData_"+head)
	//转换成返回结构体
	rf := pf.ToReturnType(i)
	return rf.f, rf.e
}
func (pf *ProtocolFactory) ToReturnType(i interface{}) ReturnType {
	switch i.(type) {
	case ReturnType:
		return i.(ReturnType)
	default:
		panic("ToReturnType类型无法匹配")
		return ReturnType{nil, nil}
	}
}
func (pf *ProtocolFactory) CallMethod(i interface{}, methodName string) interface{} {
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

	panic(pf.path + "没有实现协议头的数据读取类[GetData_" + pf.head() + "]")
}
func (pf *ProtocolFactory) GetData_file() ReturnType {
	head := pf.head()
	path := strings.SplitAfter(pf.path, head+split)[1]
	f, e := os.Open(path)
	return ReturnType{f, e}
}
func (pf *ProtocolFactory) GetData_http() ReturnType {
	return ReturnType{nil, nil}
}