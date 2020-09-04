package test

import (
	"MPDCDS_FTPServer/src/file-driver"
	"MPDCDS_FTPServer/src/ftp-server"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRealPath(t *testing.T) {
	factory := &file_driver.FileDriverFactory{
		Perm: ftp_server.NewSimplePerm("user", "group"),
	}
	driver, _ := factory.NewDriver()

	var path string = "/"
	paths := strings.Split(path, "/")
	dir, _ := os.Getwd()

	paths0 := strings.Split(dir, "\\")

	fmt.Println(runtime.GOOS)

	fmt.Println("path:" + path)
	fmt.Println("finalpath:" + filepath.Join(append([]string{paths0[0], driver.RealPath(path)}, paths...)...))

}
