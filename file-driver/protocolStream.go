package filedriver

import (
	"github.com/pkg/errors"
	"io"
	"reflect"

	"strings"
)

func ReadFile(path string) (int64, io.ReadCloser, error) {
	pf := ProtocolFactory{path}
	return pf.getData()
}

type ReturnType struct {
	s int64
	f io.ReadCloser
	e error
}
type ProtocolFactory struct {
	path string
}

const protocolSplit = "://"

func (pf *ProtocolFactory) ProtocolSplit() string {
	return protocolSplit
}
func (pf *ProtocolFactory) head() (head string, err error) {
	if strings.Contains(pf.path, protocolSplit) {
		head := strings.Split(pf.path, protocolSplit)[0]
		if head == "" {
			err = errors.New(pf.path + "没有找到协议头")
		} else {
			err = nil
		}
		head = strings.ToLower(head)
		return head, err
	} else {
		panic(pf.path + "没有找到协议头规格信息,eg:[file://][http://]")
	}
}
func (pf *ProtocolFactory) thePath() (thepath string, err error) {
	head, e0 := pf.head()
	thepath = strings.SplitAfter(pf.path, head+protocolSplit)[1]
	return thepath, e0
}
func (pf *ProtocolFactory) getData() (int64, io.ReadCloser, error) {
	head, e0 := pf.head()
	if e0 != nil {
		return 0, nil, e0
	}
	//动态调用接口
	i := pf.callMethod(pf, head)
	//转换成返回结构体
	rf := pf.toReturnType(i)
	return rf.s, rf.f, rf.e
}
func (pf *ProtocolFactory) toReturnType(i interface{}) ReturnType {
	switch i.(type) {
	case ReturnType:
		return i.(ReturnType)
	default:
		panic("ToReturnType类型无法匹配")
		return ReturnType{0, nil, nil}
	}
}
func (pf *ProtocolFactory) callMethod(i interface{}, head string) interface{} {
	//设置搜索接口名
	var methodName = "GetData_" + head

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

	panic(pf.path + "没有实现协议头的数据读取类[GetData_" + head + "]")
}
