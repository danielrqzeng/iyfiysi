// gen by iyfiysi at 2020-08-02 22:41:29.8925854 +0800 CST m=+13.803003301
package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
)

/**
 * 获取日志
 * filePath 日志文件路径
 * level 日志级别
 * maxSize 每个日志文件保存的最大尺寸 单位：M
 * maxBackups 日志文件最多保存多少个备份
 * maxAge 文件最多保存多少天
 * compress 是否压缩
 * serviceName 服务名
 */
func NewLogger(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, serviceName string) *zap.Logger {
	core := newCore(filePath, level, maxSize, maxBackups, maxAge, compress)
	//return zap.New(core, zap.AddCaller(), zap.Development(), zap.Fields(zap.String("serviceName", serviceName)))
	return zap.New(core, zap.AddCaller(), zap.Development()) //, zap.Fields(zap.String("serviceName", serviceName)))
}

/**
 * zapcore构造
 */
func newCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) zapcore.Core {
	//日志文件路径配置2
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		//zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewConsoleEncoder(encoderConfig),
		//zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)), // 打印到和文件
		atomicLevel,                                         // 日志级别
	)
}

func NewJsonLogger(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, serviceName string) *zap.Logger {
	core := newJsonCore(filePath, level, maxSize, maxBackups, maxAge, compress)
	//return zap.New(core, zap.AddCaller(), zap.Development(), zap.Fields(zap.String("serviceName", serviceName)))
	return zap.New(core, zap.AddCaller(), zap.Development()) //, zap.Fields(zap.String("serviceName", serviceName)))
}

/**
 * zapcore构造
 */
func newJsonCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) zapcore.Core {
	//日志文件路径配置2
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		//zapcore.NewConsoleEncoder(encoderConfig),
		//zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)), // 打印到和文件
		atomicLevel,                                         // 日志级别
	)
}

var MainLogger *zap.Logger
var HttpLogger *zap.Logger
var ApiLogger *zap.Logger
var DBLogger *zap.Logger
var GrpcLogger *zap.Logger

func init() {
	MainLogger = NewLogger("./logs/main.log", zapcore.DebugLevel, 128, 30, 7, true, "Main")
	HttpLogger = NewLogger("./logs/http.log", zapcore.InfoLevel, 128, 30, 7, true, "http")
	GrpcLogger = NewLogger("./logs/grpc.log", zapcore.InfoLevel, 128, 30, 7, true, "grpc")
	ApiLogger = NewLogger("./logs/api.log", zapcore.InfoLevel, 128, 30, 7, true, "api")
	DBLogger = NewLogger("./logs/db.log", zapcore.InfoLevel, 128, 30, 7, true, "db")
}

func InitLogger(logDir string, logLever zapcore.Level) {
	MainLogger = NewLogger(path.Join(logDir, "/main.log"), logLever, 32, 30, 7, true, "Main")
	HttpLogger = NewLogger(path.Join(logDir, "/http.log"), logLever, 32, 30, 7, true, "http")
	GrpcLogger = NewLogger(path.Join(logDir, "./logs/grpc.log"), zapcore.InfoLevel, 128, 30, 7, true, "grpc")
	ApiLogger = NewLogger(path.Join(logDir, "/api.log"), logLever, 23, 30, 7, true, "api")
	DBLogger = NewLogger(path.Join(logDir, "/db.log"), logLever, 32, 30, 7, true, "db")
}
