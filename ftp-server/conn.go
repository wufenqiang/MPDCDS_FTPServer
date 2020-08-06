package ftp_server

import (
	"MPDCDS_FTPServer/conf"
	"MPDCDS_FTPServer/logger"
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	defaultWelcomeMessage = "Welcome to the Go MPDCDS_FTPServer"
)

/**
returns a random 20 char string that can be used as a unique session ID
*/
func newSessionID() string {
	hash := sha256.New()
	_, err := io.CopyN(hash, rand.Reader, 50)
	if err != nil {
		return "????????????????????"
	}
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr[0:20]
}

type Conn struct {
	conn          net.Conn
	controlReader *bufio.Reader
	controlWriter *bufio.Writer
	dataConn      DataSocket
	driver        Driver
	//logger        *zap.Logger
	server      *Server
	tlsConfig   *tls.Config
	sessionID   string
	namePrefix  string
	reqUser     string
	token       string
	user        string
	renameFrom  string
	lastFilePos int64
	appendData  bool
	closed      bool
	tls         bool
}

func (conn *Conn) LoginUser() string {
	return conn.user
}

func (conn *Conn) LoginToken() string {
	return conn.token
}

func (conn *Conn) IsLogin() bool {
	return len(conn.user) > 0
}

func (conn *Conn) PublicIp() string {
	return conn.server.PublicIp
}

func (conn *Conn) passiveListenIP() string {
	var listenIP string
	if len(conn.PublicIp()) > 0 {
		listenIP = conn.PublicIp()
	} else {
		listenIP = conn.conn.LocalAddr().(*net.TCPAddr).IP.String()
	}

	lastIdx := strings.LastIndex(listenIP, ":")
	if lastIdx <= 0 {
		return listenIP
	}
	return listenIP[:lastIdx]
}

func (conn *Conn) PassivePort() int {
	if len(conn.server.PassivePorts) > 0 {
		portRange := strings.Split(conn.server.PassivePorts, "-")

		if len(portRange) != 2 {
			log.Println("empty port")
			return 0
		}

		minPort, _ := strconv.Atoi(strings.TrimSpace(portRange[0]))
		maxPort, _ := strconv.Atoi(strings.TrimSpace(portRange[1]))

		return minPort + mrand.Intn(maxPort-minPort)
	}
	// let system automatically chose one port
	return 0
}

/**
Serve starts an endless loop that reads FTP commands from the client and
responds appropriately. terminated is a channel that will receive a true
message when the connection closes. This loop will be running inside a
goroutine, so use this channel to be notified when the connection can be
cleaned up.
*/
func (conn *Conn) Serve() {
	//conn.logger.Print(conn.sessionID, "Connection Established")
	var msg string
	if conf.Sysconfig.ShadeInLog {
		//msg=fmt.Sprintf("[%s][%s]%s\r\n",conn.sessionID,conn.conn.RemoteAddr(),"Connection Established")
		msg = fmt.Sprintf("[%s][%s]%s", conn.sessionID, conn.conn.RemoteAddr(), "Connection Established")
	} else {
		//msg=fmt.Sprintf("[%s]%s\r\n",conn.sessionID,"Connection Established")
		msg = fmt.Sprintf("[%s]%s", conn.sessionID, "Connection Established")
	}
	logger.GetLogger().Info(msg)

	// send welcome
	conn.writeMessage(220, conn.server.WelcomeMessage)
	// read commands
	for {
		line, err := conn.controlReader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				//conn.logger.Print(conn.sessionID, fmt.Sprint("read error:", err))
				logger.GetLogger().Error(conn.sessionID + ",read error:" + err.Error())
			}

			break
		}
		conn.receiveLine(line)
		// QUIT command closes connection, break to avoid error on reading from
		// closed socket
		if conn.closed == true {
			break
		}
	}
	conn.Close()
	//conn.logger.Print(conn.sessionID, "Connection Terminated")
	//logger.GetLogger().Info(conn.sessionID + " Connection Terminated")
	if conf.Sysconfig.ShadeInLog {
		logger.GetLogger().Info(conn.sessionID + "(" + conn.conn.RemoteAddr().String() + ") Connection Terminated")
	} else {
		logger.GetLogger().Info(conn.sessionID + " Connection Terminated")
	}
}

/**
Close will manually close this connection, even if the client isn't ready.
*/
func (conn *Conn) Close() {
	conn.conn.Close()
	conn.closed = true
	if conn.dataConn != nil {
		conn.dataConn.Close()
		conn.dataConn = nil
	}
}

func (conn *Conn) upgradeToTLS() error {
	//conn.logger.Print(conn.sessionID, "Upgrading connectiion to TLS")
	logger.GetLogger().Info(conn.sessionID + " Upgrading connectiion to TLS")
	tlsConn := tls.Server(conn.conn, conn.tlsConfig)
	err := tlsConn.Handshake()
	if err == nil {
		conn.conn = tlsConn
		conn.controlReader = bufio.NewReader(tlsConn)
		conn.controlWriter = bufio.NewWriter(tlsConn)
		conn.tls = true
	}
	return err
}

/**
receiveLine accepts a single line FTP command and co-ordinates an appropriate response.
*/
func (conn *Conn) receiveLine(line string) {
	command, param := conn.analysisLine(line)
	//conn.logger.PrintCommand(conn.sessionID, command, param)
	conn.PrintReceive(conn.sessionID, line)
	cmdObj := commands[strings.ToUpper(command)]
	if cmdObj == nil {
		conn.writeMessage(500, "Command not found")
		return
	}
	if cmdObj.RequireParam() && param == "" {
		conn.writeMessage(553, "action aborted, required param missing")
	} else if cmdObj.RequireAuth() && conn.user == "" {
		conn.writeMessage(530, "not logged in")
	} else {
		cmdObj.Execute(conn, param)
	}
}

/**
分析获取到的信息，获取命令及参数
*/
func (conn *Conn) analysisLine(line string) (string, string) {
	params := strings.SplitN(strings.Trim(line, "\r\n"), " ", 2)
	if len(params) == 1 {
		return params[0], ""
	}
	return params[0], strings.TrimSpace(params[1])
}

/**
writeMessage will send a standard FTP response back to the client.
*/
func (conn *Conn) writeMessage(code int, message string) (wrote int, err error) {
	line := fmt.Sprintf("%d %s\r\n", code, message)

	conn.PrintResponse(conn.sessionID, line)

	wrote, err = conn.controlWriter.WriteString(line)
	conn.controlWriter.Flush()
	return
}
func (conn *Conn) writeMessageMultiline(code int, message string) (wrote int, err error) {
	line := fmt.Sprintf("%d-%s\r\n%d END\r\n", code, message, code)

	conn.PrintResponse(conn.sessionID, line)

	wrote, err = conn.controlWriter.WriteString(line)
	conn.controlWriter.Flush()
	return
}

// buildPath takes a client supplied path or filename and generates a safe
// absolute path within their account sandbox.
//
//    buildpath("/")
//    => "/"
//    buildpath("one.txt")
//    => "/one.txt"
//    buildpath("/files/two.txt")
//    => "/files/two.txt"
//    buildpath("files/two.txt")
//    => "/files/two.txt"
//    buildpath("/../../../../etc/passwd")
//    => "/etc/passwd"
//
// The driver implementation is responsible for deciding how to treat this path.
// Obviously they MUST NOT just read the path off disk. The probably want to
// prefix the path with something to scope the users access to a sandbox.
func (conn *Conn) buildPath(filename string) (fullPath string) {
	if len(filename) > 0 && filename[0:1] == "/" {
		fullPath = filepath.Clean(filename)
	} else if len(filename) > 0 && filename != "-a" {
		fullPath = filepath.Clean(conn.namePrefix + "/" + filename)
	} else {
		fullPath = filepath.Clean(conn.namePrefix)
	}
	fullPath = strings.Replace(fullPath, "//", "/", -1)
	fullPath = strings.Replace(fullPath, string(filepath.Separator), "/", -1)
	return
}

/**
sendOutofbandData will send a string to the client via the currently open data socket. Assumes the socket is open and ready to be used.
Nlst
List
Mlsd
*/
func (conn *Conn) sendOutofbandData(data []byte) {
	bytes := len(data)
	if conn.dataConn != nil {
		conn.dataConn.Write(data)
		conn.dataConn.Close()

		//if conf.Sysconfig.ShadeInLog {
		//	logger.GetLogger().Info(conn.sessionID + "(DataPort:" + conn.dataConn.Host() + ":" + strconv.Itoa(conn.dataConn.Port()) + "------>)")
		//}

		conn.dataConn = nil
	}
	//message := "Closing data connection, sent " + strconv.Itoa(bytes) + " bytes"
	message := fmt.Sprintf("Closing data connection, sent %d bytes", bytes)
	conn.writeMessage(226, message)
}

/**
Retr
*/
func (conn *Conn) sendOutofBandDataWriter(data io.ReadCloser) error {
	conn.lastFilePos = 0
	bytes, err := io.Copy(conn.dataConn, data)
	if err != nil {
		conn.dataConn.Close()
		conn.dataConn = nil
		return err
	}
	//message := "Closing data connection, sent " + strconv.Itoa(int(bytes)) + " bytes"
	message := fmt.Sprintf("Closing data connection, sent %d bytes", bytes)
	conn.writeMessage(226, message)
	conn.dataConn.Close()

	//if conf.Sysconfig.ShadeInLog {
	//	//println("conn.PublicIp()="+conn.PublicIp())
	//	//println("conn.PassivePort()=" + strconv.Itoa(conn.PassivePort()))
	//	//println("conn.passiveListenIP()=" + conn.passiveListenIP())
	//
	//	//远程连接命令端口
	//	//println("conn.conn.RemoteAddr()="+conn.conn.RemoteAddr().String())
	//	logger.GetLogger().Info(conn.sessionID + "(DataPort:" + conn.dataConn.Host() + ":" + strconv.Itoa(conn.dataConn.Port()) + "------>)")
	//}

	conn.dataConn = nil

	return nil
}

func (conn *Conn) PrintReceive(sessionId string, line string) {
	var line0 string
	if conf.Sysconfig.ShadeInLog {
		if strings.Contains(line, "PASS") {
			//line0 = "PASS ******(隐藏)\r\n"
			line0 = "PASS ******(隐藏)"
		} else {
			line0 = line
		}
	} else {
		line0 = line
	}
	msg := fmt.Sprintf("[%s][%s %s %s:%d]%s", sessionId, conn.conn.RemoteAddr(), ">>>>>>", conn.server.Hostname, conn.server.Port, line0)
	logger.GetLogger().Info(msg)
}
func (conn *Conn) PrintResponse(sessionId string, line string) {
	msg := fmt.Sprintf("[%s][%s %s %s:%d]%s", sessionId, conn.conn.RemoteAddr(), "<<<<<<", conn.server.Hostname, conn.server.Port, line)
	logger.GetLogger().Info(msg)
}
