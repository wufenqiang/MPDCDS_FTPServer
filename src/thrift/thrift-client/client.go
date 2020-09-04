package thrift_client

import (
	"MPDCDS_FTPServer/src/conf"
	"github.com/apache/thrift/lib/go/thrift"
	"gitlab.weather.com.cn/wufenqiang/MPDCDSPro/src/thrift/thrift-client-core"
	"gitlab.weather.com.cn/wufenqiang/MPDCDSPro/src/thrift/thriftcore"

	"net"
)

var ThriftClient *thriftcore.MPDCDSProServiceClient
var Transport thrift.TTransport

func init() {
	thrifthostport := net.JoinHostPort(conf.Sysconfig.ThriftHost, conf.Sysconfig.ThriftPort)
	ThriftClient, Transport = thrift_client_core.ConnectHostPort(thrifthostport)
}

func ThriftClose() {
	thrift_client_core.Close(Transport)
}
