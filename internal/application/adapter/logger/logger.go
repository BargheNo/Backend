package loggerimpl

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	constants *bootstrap.Constants
	zap       *zap.Logger
}

var loggerInstance logger.Logger

func NewLogger(config *bootstrap.Logger, constants *bootstrap.Constants) (*Logger, error) {
	if config.LogLevel == "" {
		config.LogLevel = constants.LogLevel.Info
	}

	var level zapcore.Level
	switch config.LogLevel {
	case constants.LogLevel.Debug:
		level = zapcore.DebugLevel
	case constants.LogLevel.Info:
		level = zapcore.InfoLevel
	case constants.LogLevel.Warn:
		level = zapcore.WarnLevel
	case constants.LogLevel.Error:
		level = zapcore.ErrorLevel
	case constants.LogLevel.Fatal:
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	consoleOutput, err := strconv.ParseBool(config.ConsoleOutput)
	if err != nil {
		consoleOutput = true
		log.Println("Error during checking console output enable. set it to true by default")
	}
	if consoleOutput {
		consoleCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	if config.LogFile != "" {
		file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		fileCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(file),
			level,
		)
		cores = append(cores, fileCore)
	}

	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	logger := &Logger{
		constants: constants,
		zap:       zapLogger,
	}

	loggerInstance = logger
	return logger, nil
}

func GetLogger() logger.Logger {
	if loggerInstance == nil {
		loggerInstance = &Logger{
			zap: zap.NewExample(),
		}
	}
	return loggerInstance
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

func (l *Logger) WithFields(fields map[string]interface{}) logger.Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return &Logger{
		zap: l.zap.With(zapFields...),
	}
}

func (l *Logger) Close() {
	l.zap.Sync()
}
