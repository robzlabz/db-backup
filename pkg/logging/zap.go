package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()

	// Konfigurasi untuk output console only
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stdout"}

	// Konfigurasi encoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.Encoding = "console"

	// Buat logger
	logger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic(err)
	}

	Logger = logger
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}

// Sugar returns sugared logger
func Sugar() *zap.SugaredLogger {
	return GetLogger().Sugar()
}

// Info logs a message at InfoLevel
func Info(msg string, fields ...interface{}) {
	Sugar().Infow(msg, fields...)
}

// Error logs a message at ErrorLevel
func Error(msg string, fields ...interface{}) {
	Sugar().Errorw(msg, fields...)
}

// Debug logs a message at DebugLevel
func Debug(msg string, fields ...interface{}) {
	Sugar().Debugw(msg, fields...)
}

// Warn logs a message at WarnLevel
func Warn(msg string, fields ...interface{}) {
	Sugar().Warnw(msg, fields...)
}

// Infof logs a formatted message at InfoLevel
func Infof(template string, args ...interface{}) {
	Sugar().Infof(template, args...)
}

// Errorf logs a formatted message at ErrorLevel
func Errorf(template string, args ...interface{}) {
	Sugar().Errorf(template, args...)
}

// Debugf logs a formatted message at DebugLevel
func Debugf(template string, args ...interface{}) {
	Sugar().Debugf(template, args...)
}

// Warnf logs a formatted message at WarnLevel
func Warnf(template string, args ...interface{}) {
	Sugar().Warnf(template, args...)
}
