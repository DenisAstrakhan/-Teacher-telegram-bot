package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// loglevel возможность задать минимальный уровень который будет логироватся DEBUG INFO WARN ERROR
func NewLogger(loglevel string) (*zap.Logger, func() error, error) {
	lvl := zap.NewAtomicLevel()
	if err := lvl.UnmarshalText([]byte(loglevel)); err != nil {
		return nil, nil, fmt.Errorf("unmarshal log level: %w", err)
	}
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, nil, fmt.Errorf("mkdir log folder: %w", err)
	}
	timestamp := time.Now().UTC().Format("2006-01-02T15.04.05.000000")
	logFilePath := filepath.Join("logs", fmt.Sprintf("%s.log", timestamp))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %w", err)
	}
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15.04.05.000000")
	encoder := zapcore.NewConsoleEncoder(cfg)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), lvl),
		zapcore.NewCore(encoder, zapcore.AddSync(logFile), lvl),
	)
	logger := zap.New(
		core,
		zap.AddCaller(),                       // параметр логирует место из которого произведён сам лог
		zap.AddStacktrace(zapcore.ErrorLevel), //параметр говорит что стек вызовов надо показывать для уровней Error и выше
	)
	return logger, logFile.Close, nil
}
