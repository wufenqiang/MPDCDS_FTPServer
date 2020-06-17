package logger

import (
	"MPDCDS_FTPServer/conf"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"time"
)

var zapLogger *zap.Logger

// 初始化日志 logger
func InitLog(loggerpath string, level string) *zap.Logger {
	infoPath := loggerpath + "/info/"
	errPath := loggerpath + "/error/"
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	config := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder, //将级别转换成大写
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
	encoder := zapcore.NewConsoleEncoder(config)
	// 设置级别
	logLevel := zap.DebugLevel
	switch level {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "panic":
		logLevel = zap.PanicLevel
	case "fatal":
		logLevel = zap.FatalLevel
	default:
		logLevel = zap.InfoLevel
	}
	// 实现两个判断日志等级的interface  可以自定义级别展示
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= logLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= logLevel
	})

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter(infoPath)
	warnWriter := getWriter(errPath)

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		// 将info及以下写入logPath,  warn及以上写入errPath
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
		//日志都会在console中展示
		zapcore.NewCore(zapcore.NewConsoleEncoder(config),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel)) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
}

func getWriter(FilePathName string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	hook, err := rotatelogs.New(
		FilePathName+conf.Sysconfig.ProjectName+"_%Y%m%d"+".log", // 没有使用go风格反人类的format格式.%Y%m%d%H
		//rotatelogs.WithLinkName(FilePathName),//// 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Hour*24*30),    // 保存30天
		rotatelogs.WithRotationTime(time.Hour*24), //切割频率 24小时
	)
	if err != nil {
		log.Println("日志启动异常")
		panic(err)
	}
	return hook
}

//初始化
func init() {
	//logpath:=path.Join(conf.Sysconfig.LoggerPath,conf.Sysconfig.ProjectName)
	//logger := InitLog(logpath, conf.Sysconfig.LoggerLevel)

	logger := InitLog(conf.Sysconfig.LoggerPath, conf.Sysconfig.LoggerLevel)

	//logger.Info("Logger init......\r\n")
	logger.Info("Logger init......")

	zapLogger = logger
}

//获取zap logger 对象
func GetLogger() *zap.Logger {
	return zapLogger
}
