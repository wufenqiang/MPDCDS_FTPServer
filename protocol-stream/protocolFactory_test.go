package protocol_stream

import (
	"io"

	"os"
	"testing"
)

func TestProtocolFactory(t *testing.T) {
	pf := ProtocolFactory{"file://G:/tmp/id_rsa.pub"}
	_, r, _ := pf.getData()

	w, _ := os.Create("G:/tmp/test.pub")

	i, e := io.Copy(w, r)

	println(i)
	println(e)
}