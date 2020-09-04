package conf

type sysconfig struct {

	//thrift 服务ip
	ThriftHost string `json:"ThriftHost"`
	ThriftPort string `json:"ThriftPort"`

	//日志存储地址、级别
	LoggerPath  string `json:"LoggerPath"`
	LoggerLevel string `json:"LoggerLevel"`

	//日志中显示相关密文
	ShadeInLog bool `json:ShadeInLog`

	//ftp相关配置
	//"FTPHost":"10.16.39.75",
	//FTPHost    string `json:FTPHost`
	FTPCmdPort int `json:FTPCmdPort`

	//网盘根目录
	NetworkDisk string `json:NetworkDisk`
}

const ProjectName = "MPDCDS_FTPServer"
