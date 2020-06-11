package filedriver

import (
	"MPDCDS_FTPServer/conf"
	"MPDCDS_FTPServer/server"
	"MPDCDS_FTPServer/thrift/client"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileDriver struct {
	//RootPath string
	server.Perm
}

func (driver *FileDriver) RealPath(path string) string {
	/**
	通过API获取数据的真实路径
	*/
	//var RootPath string = "/tmp"
	p := ProtocolFactory{path}
	var head = p.head()
	var thepath string = p.thePath()
	var RootPath string = conf.Sysconfig.NetworkDisk

	paths := strings.Split(thepath, "/")

	protocolSplit := p.ProtocolSplit()

	finalpath := head + protocolSplit + filepath.Join(append([]string{RootPath}, paths...)...)
	return finalpath
}
func (driver *FileDriver) Init(conn *server.Conn) {
	//driver.conn = conn
}
func (driver *FileDriver) ChangeDir(path string, token string) error {
	ext := filepath.Ext(path)
	if ext != "" {
		return errors.New("Not a directory")
	}
	//获取操作对象
	tClient, tTransport := client.Connect()
	ctx := context.Background()
	dirAuth, err := tClient.DirAuth(ctx, token, path)
	//关闭tTransport
	client.Close(tTransport)
	if err != nil {
		return err
	}

	if dirAuth.Status == 0 {
		return nil
	}
	message := dirAuth.Msg
	return errors.New(message)
	//rPath := driver.realPath(path)
	//f, err := os.Lstat(rPath)
	//if err != nil {
	//	return err
	//}
	//if f.IsDir() {
	//	return nil
	//}
	//return errors.New("Not a directory")
}
func (driver *FileDriver) Stat(path string) (server.FileInfo, error) {
	basepath := driver.RealPath(path)
	rPath, err := filepath.Abs(basepath)
	if err != nil {
		return nil, err
	}
	f, err := os.Lstat(rPath)
	if err != nil {
		return nil, err
	}
	mode, err := driver.Perm.GetMode(path)
	if err != nil {
		return nil, err
	}
	if f.IsDir() {
		mode |= os.ModeDir
	}
	owner, err := driver.Perm.GetOwner(path)
	if err != nil {
		return nil, err
	}
	group, err := driver.Perm.GetGroup(path)
	if err != nil {
		return nil, err
	}
	return &FileInfo{f, mode, owner, group}, nil
}
func (driver *FileDriver) ListDir(path string, callback func(server.FileInfo) error) error {
	basepath := driver.RealPath(path)
	return filepath.Walk(basepath, func(f string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rPath, _ := filepath.Rel(basepath, f)
		if rPath == info.Name() {
			mode, err := driver.Perm.GetMode(rPath)
			if err != nil {
				return err
			}
			if info.IsDir() {
				mode |= os.ModeDir
			}
			owner, err := driver.Perm.GetOwner(rPath)
			if err != nil {
				return err
			}
			group, err := driver.Perm.GetGroup(rPath)
			if err != nil {
				return err
			}
			err = callback(&FileInfo{info, mode, owner, group})
			if err != nil {
				return err
			}
			if info.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})
}
func (driver *FileDriver) DeleteDir(path string) error {
	rPath := driver.RealPath(path)
	f, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return os.Remove(rPath)
	}
	return errors.New("Not a directory")
}
func (driver *FileDriver) DeleteFile(path string) error {
	rPath := driver.RealPath(path)
	f, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return os.Remove(rPath)
	}
	return errors.New("Not a file")
}
func (driver *FileDriver) Rename(fromPath string, toPath string) error {
	oldPath := driver.RealPath(fromPath)
	newPath := driver.RealPath(toPath)
	return os.Rename(oldPath, newPath)
}
func (driver *FileDriver) MakeDir(path string) error {
	rPath := driver.RealPath(path)
	return os.MkdirAll(rPath, os.ModePerm)
}
func (driver *FileDriver) GetFile(path string, offset int64) (int64, io.ReadCloser, error) {
	var rPath string
	if strings.HasPrefix(path, "file") {
		rPath = driver.RealPath(path)
	} else {
		rPath = path
	}

	f, err0 := ReadFile(rPath)
	if err0 != nil {
		return 0, nil, err0
	}
	switch f.(type) {
	case *os.File:
		f0 := f.(*os.File)
		info, err := f0.Stat()
		if err != nil {
			return 0, nil, err
		}
		f0.Seek(offset, os.SEEK_SET)
		return info.Size(), f, err
	default:
		return 0, f, nil
	}

}
func (driver *FileDriver) PutFile(destPath string, data io.Reader, appendData bool) (int64, error) {
	rPath := driver.RealPath(destPath)
	var isExist bool
	f, err := os.Lstat(rPath)
	if err == nil {
		isExist = true
		if f.IsDir() {
			return 0, errors.New("A dir has the same name")
		}
	} else {
		if os.IsNotExist(err) {
			isExist = false
		} else {
			return 0, errors.New(fmt.Sprintln("Put File error:", err))
		}
	}

	if appendData && !isExist {
		appendData = false
	}

	if !appendData {
		if isExist {
			err = os.Remove(rPath)
			if err != nil {
				return 0, err
			}
		}
		f, err := os.Create(rPath)
		if err != nil {
			return 0, err
		}
		defer f.Close()
		bytes, err := io.Copy(f, data)
		if err != nil {
			return 0, err
		}
		return bytes, nil
	}

	of, err := os.OpenFile(rPath, os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return 0, err
	}
	defer of.Close()

	_, err = of.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}

	bytes, err := io.Copy(of, data)
	if err != nil {
		return 0, err
	}

	return bytes, nil
}

type FileInfo struct {
	os.FileInfo

	mode  os.FileMode
	owner string
	group string
}

func (f *FileInfo) Mode() os.FileMode {
	return f.mode
}
func (f *FileInfo) Owner() string {
	return f.owner
}
func (f *FileInfo) Group() string {
	return f.group
}

type FileDriverFactory struct {
	//RootPath string
	server.Perm
}

func (factory *FileDriverFactory) NewDriver() (server.Driver, error) {
	return &FileDriver{ /*factory.RootPath,*/ factory.Perm}, nil
}
