package filedriver

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string
	Age  int64
}

func (p *Person) GetName() string {
	return p.Name
}
func (p *Person) GetAge() int64 {
	return p.Age
}

func TestCallMethod(t *testing.T) {
	hhg := Person{"hhg08", 222}
	//method_name := "GetName"
	method_name := "GetAge"
	name := CallMethod(hhg, method_name)
	fmt.Println(name)
}

func TestProtocolFactory(t *testing.T) {
	pf := ProtocolFactory{"file:///tmp/"}

	pf.getFile()
}
