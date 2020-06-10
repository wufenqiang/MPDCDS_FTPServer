package filedriver

import (
	"testing"
)

func TestProtocolFactory(t *testing.T) {
	pf := ProtocolFactory{"ftp://tmp/"}
	pf.getData()
}
