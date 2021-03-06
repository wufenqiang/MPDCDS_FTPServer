package file_driver

import (
	"MPDCDS_FTPServer/src/conf"
	"MPDCDS_FTPServer/src/ftp-server"
	"MPDCDS_FTPServer/src/logger"
	"MPDCDS_FTPServer/src/thrift/thrift-client"
	"gitlab.weather.com.cn/wufenqiang/MPDCDSPro/src/protocol-stream"
	"gitlab.weather.com.cn/wufenqiang/MPDCDSPro/src/thrift/thriftcore"

	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileDriver struct {
	//RootPath string
	ftp_server.Perm
}

func (driver *FileDriver) RealPath(FtpRelativePath string) string {
	///**
	//通过API获取数据的真实路径
	//*/
	////var RootPath string = "/tmp"
	//p := protocol_stream.ProtocolFactory{path0}
	//var head, _ = p.Head()
	//var thepath, _ = p.ThePath()
	//var RootPath string = conf.Sysconfig.NetworkDisk
	//
	//paths := strings.Split(thepath, "/")
	//
	//protocolSplit := p.ProtocolSplit()
	//
	//finalpath := head + protocolSplit + filepath.Join(append([]string{RootPath}, paths...)...)
	//return finalpath
	var RootPath string = conf.Sysconfig.NetworkDisk
	return protocol_stream.FtpRelativePath2AbsoluteURL(RootPath, FtpRelativePath)
}
func (driver *FileDriver) Init(conn *ftp_server.Conn) {
	//driver.conn = conn
}
func (driver *FileDriver) ChangeDir(path string, token string) error {
	ext := filepath.Ext(path)
	if ext != "" {
		return errors.New("Not a directory")
	}

	dirAuthInfo := thriftcore.NewDirAuthInfo()
	dirAuthInfo.Token = token
	dirAuthInfo.AbsPath = path

	//获取操作对象
	//tClient, tTransport := utils.ThriftConnect()
	ctx := context.Background()
	dirAuth, err := thrift_client.ThriftClient.DirAuth(ctx, dirAuthInfo)
	//关闭tTransport
	thrift_client.ThriftClose()
	if err != nil {
		return err
	}

	dirAuth0 := thrift_client.DirAuthReturn{dirAuth}

	if dirAuth0.Status == 0 {
		return nil
	}
	message := dirAuth0.DirAuth2Msg()
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
func (driver *FileDriver) Stat(path string) (ftp_server.FileInfo, error) {
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
func (driver *FileDriver) ListDir(path string, callback func(ftp_server.FileInfo) error) error {
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

	s, f, e0 := protocol_stream.ReadProtocol(rPath)
	if e0 != nil {
		return s, f, e0
	}
	switch f.(type) {
	case *os.File:
		f0 := f.(*os.File)
		_, e1 := f0.Stat()
		if e1 != nil {
			return s, f, e1
		}
		f0.Seek(offset, os.SEEK_SET)
	default:
		if offset != 0 {
			var e2 error = errors.New("io.ReadCloser无法使用偏移量offset(" + strconv.FormatInt(offset, 10) + ")")
			logger.GetLogger().Warn(e2.Error())
			return s, f, e2
		}
	}
	return s, f, e0
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
	ftp_server.Perm
}

func (factory *FileDriverFactory) NewDriver() (ftp_server.Driver, error) {
	return &FileDriver{ /*factory.RootPath,*/ factory.Perm}, nil
}
