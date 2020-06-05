package client

import (
	"MPDCDS_FTPServer/conf"
	"MPDCDS_FTPServer/thrift/MPDCDS_BackendService"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"net"
	"os"
)

//创建客户端连接，获取连接对象
func Connect() (*MPDCDS_BackendService.MPDCDS_BackendServiceClient, thrift.TTransport) {
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, err := thrift.NewTSocket(net.JoinHostPort(conf.Sysconfig.NetworkAddr, conf.Sysconfig.ThriftPort))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
	}
	trans := thrift.NewTFramedTransport(transport)
	//useTransport := transportFactory.GetTransport(transport)
	client := MPDCDS_BackendService.NewMPDCDS_BackendServiceClientFactory(trans, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to "+conf.Sysconfig.NetworkAddr, conf.Sysconfig.ThriftPort, " ", err)
	}
	//defer transport.Close()
	return client, transport
}

//释放transport
func Close(transport thrift.TTransport) {
	defer transport.Close()
}
