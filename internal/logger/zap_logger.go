package logger

//https://www.cnblogs.com/zly-go/p/15500015.html
import (
	"fmt"
	"news-screen/internal/conf"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ log.Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		location = time.UTC
	}
	enc.AppendString(t.In(location).Format("2006-01-02 15:04:05.000"))
}

// GetZapLogger Logger 配置zap日志,将zap日志库引入
func GetZapLogger(logConfig *conf.Logging) log.Logger {
	//配置zap日志库的编码器
	encoder := zapcore.EncoderConfig{
		TimeKey:   "time",
		LevelKey:  "level",
		NameKey:   "logger",
		CallerKey: "caller",
		// MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeTime:     CustomTimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	return NewZapLogger(
		encoder,
		zap.NewAtomicLevelAt(ConvertLevel(logConfig.Level)),
		logConfig,
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development(),
	)
}
func GetPathName(path string, name string) string {
	sys := runtime.GOOS
	// 检查当前操作系统
	if sys == "windows" {
		// Windows 系统下使用当前目录
		path = "."
	} else if sys == "mac" {
		path = "."
	}
	// 创建目录
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		// log.Fatalf("Failed to create log directory: %v", err)
		// 路径不存在，打印当前运行路径
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}
		fmt.Printf("Current running path: %s\n", currentDir)

		// 将路径修改为 ../../logs/
		path = "../../logs"
		os.MkdirAll(path, os.ModePerm)
	}

	return filepath.Join(path, name)
}

func ConvertLevel(level string) zapcore.Level {
	// 解析日志级别
	var zapLevel zapcore.Level

	switch level {
	case "debug", "Debug":
		zapLevel = zapcore.DebugLevel
	case "info", "Info":
		zapLevel = zapcore.InfoLevel
	case "warn", "Warn":
		zapLevel = zapcore.WarnLevel
	case "error", "Error":
		zapLevel = zapcore.ErrorLevel
	case "dpanic", "DPanic":
		zapLevel = zapcore.DPanicLevel
	case "panic", "Panic":
		zapLevel = zapcore.PanicLevel
	case "fatal", "Fatal":
		zapLevel = zapcore.FatalLevel
	default:
		// 默认级别为 Info
		zapLevel = zapcore.InfoLevel
	}
	log.Info("log level is ", zapLevel)
	return zapLevel
}

// 日志自动切割，采用 lumberjack 实现的
func getLogWriter(logFileConf *conf.Logging_File) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   GetPathName(logFileConf.Path, logFileConf.Name),
		MaxSize:    (int)(logFileConf.MaxSize),    //日志的最大大小（M）
		MaxBackups: (int)(logFileConf.MaxBackups), //日志的最大保存数量
		MaxAge:     (int)(logFileConf.MaxAge),     //日志文件存储最大天数
		Compress:   logFileConf.Compress,          //是否执行压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

// NewZapLogger return a zap logger.
func NewZapLogger(encoder zapcore.EncoderConfig, level zap.AtomicLevel, logConfig *conf.Logging, opts ...zap.Option) *ZapLogger {
	//日志切割
	writeSyncer := getLogWriter(logConfig.File)
	var core zapcore.Core
	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder), // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writeSyncer)), // 打印到控制台和文件
		level, // 日志级别
	)

	//开发模式下打印到标准输出
	//// --根据配置文件判断输出到控制台还是日志文件--
	//if conf.GetConfig().GetString("project.mode") == "dev" {
	//	core = zapcore.NewCore(
	//		zapcore.NewConsoleEncoder(encoder),                      // 编码器配置
	//		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), // 打印到控制台
	//		level, // 日志级别
	//	)
	//} else {
	//	core = zapcore.NewCore(
	//		zapcore.NewJSONEncoder(encoder), // 编码器配置
	//		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writeSyncer)), // 打印到控制台和文件
	//		level, // 日志级别
	//	)
	//}
	zapLogger := zap.New(core, opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync}
}

// Log 实现log接口
func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Key values must appear in pairs: ", keyvals))
		return nil
	}

	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug("", data...)
	case log.LevelInfo:
		l.log.Info("", data...)
	case log.LevelWarn:
		l.log.Warn("", data...)
	case log.LevelError:
		l.log.Error("", data...)
	case log.LevelFatal:
		l.log.Fatal("", data...)
	}
	return nil
}
