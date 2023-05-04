package logx

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
)

func Init(level string, fileOut bool, filename string, maxSize, maxBackups, maxAge int, compress bool) {
	core := zapcore.NewCore(
		getEncoder(),
		getLogWriter(fileOut, filename, maxSize, maxBackups, maxAge, compress),
		caseLevel(level),
	)
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(fileOut bool, filename string, maxSize, maxBackups, maxAge int, compress bool) zapcore.WriteSyncer {
	if fileOut {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   fmt.Sprintf("./static/%s", filename),
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		}
		ws := io.MultiWriter(lumberJackLogger, os.Stdout)
		return zapcore.AddSync(ws)
	}
	ws := io.MultiWriter(os.Stdout)
	return zapcore.AddSync(ws)
}

func caseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
