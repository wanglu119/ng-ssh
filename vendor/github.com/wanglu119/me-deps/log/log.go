package log

import (
	"os"
	
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// Very verbose messages for debugging specific issues
	LevelDebug = "debug"
	// Default log level, informational
	LevelInfo = "info"
	// Warnings are messages about possible issues
	LevelWarn = "warn"
	// Errors are messages about things we know are problems
	LevelError = "error"
)

// Type and function aliases from zap to limit the libraries scope into MM code
type Field = zapcore.Field

var Int64 = zap.Int64
var Int32 = zap.Int32
var Int = zap.Int
var Uint32 = zap.Uint32
var String = zap.String
var Any = zap.Any
var Err = zap.Error
var NamedErr = zap.NamedError
var Bool = zap.Bool
var Duration = zap.Duration

type LoggerConfiguration struct {
	EnableConsole bool
	ConsoleJson   bool
	ConsoleLevel  string
	EnableFile    bool
	FileJson      bool
	FileLevel     string
	FileLocation  string
}

type Logger struct {
	*zap.Logger
	consoleLevel zap.AtomicLevel
	fileLevel    zap.AtomicLevel
	fileLogger *lumberjack.Logger
}

var logger *Logger

func getZapLevel(level string) zapcore.Level {
	switch level {
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func makeEncoder(json bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
//	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	if json {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func NewLogger(config *LoggerConfiguration) *Logger {
	cores := []zapcore.Core{}
	
	consoleLevel := zap.NewAtomicLevelAt(getZapLevel(config.ConsoleLevel))
	fileLevel :=    zap.NewAtomicLevelAt(getZapLevel(config.FileLevel))
	
	logger := &Logger{
		consoleLevel: consoleLevel,
		fileLevel:    fileLevel,
	}
	
	if config.EnableConsole {
		writer := zapcore.Lock(os.Stderr)
		core := zapcore.NewCore(makeEncoder(config.ConsoleJson), writer, consoleLevel)
		cores = append(cores, core)
	}

	if config.EnableFile {
		logger.fileLogger = &lumberjack.Logger{
			Filename: config.FileLocation,
			MaxSize:  100,
			MaxBackups: 30,
			MaxAge: 7,
			Compress: true,
		}
		writer := zapcore.AddSync(logger.fileLogger)
		core := zapcore.NewCore(makeEncoder(config.FileJson), writer, fileLevel)
		cores = append(cores, core)
	}
	
	combinedCore := zapcore.NewTee(cores...)
	zlog := zap.New(combinedCore,zap.AddCaller())

	
	logger.Logger = zlog

	return logger
}

func (l *Logger) SetConsoleLevel(level string) {
	l.consoleLevel.SetLevel(getZapLevel(level))
}

func (l *Logger) SetFileLevel(level string) {
	l.fileLevel.SetLevel(getZapLevel(level))
}

func (l *Logger) SetFileLocation(fileLocation string) {
	l.fileLogger.Filename = fileLocation
	l.fileLogger.Rotate()
}

func GetLogger() *Logger {
	return logger
}

func init() {
	config := &LoggerConfiguration {
		EnableConsole: true,
		ConsoleLevel: "info",
		EnableFile: true,
		FileLevel: "info",
		FileJson: true,
	}
	logger = NewLogger(config)
}




