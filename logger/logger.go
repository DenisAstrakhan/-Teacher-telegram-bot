package logger

import (
	"TeacherBot/domain"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MultiCore пишет сразу во все выходы
type MultiCore struct {
	console    zapcore.Core
	file       zapcore.Core
	clickhouse zapcore.Core
	fields     map[string]any
}

func (m *MultiCore) Enabled(level zapcore.Level) bool {
	return true
}

func (m *MultiCore) With(fields []zapcore.Field) zapcore.Core {
	clone := *m
	if clone.fields == nil {
		clone.fields = make(map[string]any)
	}
	for _, f := range fields {
		clone.fields[f.Key] = f.Interface
	}

	clone.console = clone.console.With(fields)
	clone.file = clone.file.With(fields)
	if clone.clickhouse != nil {
		clone.clickhouse = clone.clickhouse.With(fields)
	}

	return &clone
}

func (m *MultiCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if m.Enabled(entry.Level) {
		return ce.AddCore(entry, m)
	}
	return ce
}

func (m *MultiCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// Берём поля из вызова
	allFields := fields

	// Добавляем поля из With()
	if len(m.fields) > 0 {
		for k, v := range m.fields {
			allFields = append(allFields, zap.Any(k, v))
		}
	}

	// Отправляем во все Core
	if err := m.console.Write(entry, allFields); err != nil {
		fmt.Printf("Console write error: %v\n", err)
	}
	if err := m.file.Write(entry, allFields); err != nil {
		fmt.Printf("File write error: %v\n", err)
	}
	if m.clickhouse != nil {
		if err := m.clickhouse.Write(entry, allFields); err != nil {
			fmt.Printf("ClickHouse write error: %v\n", err)
		}
	}

	return nil
}

func (m *MultiCore) Sync() error {
	if err := m.console.Sync(); err != nil {
		return err
	}
	if err := m.file.Sync(); err != nil {
		return err
	}
	if m.clickhouse != nil {
		return m.clickhouse.Sync()
	}
	return nil
}

// fieldToString извлекает строковое значение из zapcore.Field.
// Важно: для zap.String значение хранится в поле String, а не в Interface (там nil).
func fieldToString(f zapcore.Field) string {
	switch f.Type {
	case zapcore.StringType:
		return f.String
	case zapcore.StringerType:
		if f.Interface != nil {
			return fmt.Sprintf("%v", f.Interface)
		}
	case zapcore.ReflectType:
		if f.Interface != nil {
			return fmt.Sprintf("%v", f.Interface)
		}
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Uint64Type, zapcore.Uint32Type:
		return fmt.Sprintf("%d", f.Integer)
	case zapcore.Float64Type, zapcore.Float32Type:
		return fmt.Sprintf("%v", f.Interface)
	case zapcore.BoolType:
		return strconv.FormatBool(f.Integer != 0)
	case zapcore.TimeType:
		if f.Interface != nil {
			return fmt.Sprintf("%v", f.Interface)
		}
	}
	if f.Interface != nil {
		return fmt.Sprintf("%v", f.Interface)
	}
	return f.String
}

type DataBaseCore struct {
	repository domain.UserLoggerRepository
	level      zapcore.Level
	service    string
	fields     map[string]any
}

func NewDataBaseCore(repository domain.UserLoggerRepository, ctx context.Context, service string, level zapcore.Level) (*DataBaseCore, error) {
	return &DataBaseCore{
		repository: repository,
		level:      level,
		service:    service,
		fields:     make(map[string]any),
	}, nil
}

func (c *DataBaseCore) Enabled(level zapcore.Level) bool { return level >= c.level }
func (c *DataBaseCore) With(fields []zapcore.Field) zapcore.Core {
	clone := *c
	if clone.fields == nil {
		clone.fields = make(map[string]any)
	}
	for _, f := range fields {
		clone.fields[f.Key] = fieldToString(f)
	}
	return &clone
}
func (c *DataBaseCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}
func (c *DataBaseCore) Sync() error { return nil }

func (c *DataBaseCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	var userID string
	var message string

	// Проверяем c.fields (из With())
	for k, v := range c.fields {
		if k == "user_id" && v != nil {
			userID = fmt.Sprintf("%v", v)
		}
		if k == "message" && v != nil {
			message = fmt.Sprintf("%v", v)
		}
	}

	for _, f := range fields {
		//проверяе что userID ещё не установлен
		if f.Key == "user_id" && userID == "" {
			userID = fieldToString(f)
		}
		if f.Key == "message" && message == "" {
			message = fieldToString(f)
		}
	}

	if message == "" {
		message = entry.Message
	}

	return c.repository.Write(entry.Time, entry.Level.String(), c.service, userID, message)
}

// loglevel возможность задать минимальный уровень который будет логироватся DEBUG INFO WARN ERROR
func NewLogger(repository domain.UserLoggerRepository, loglevel string) (*zap.Logger, func() error, error) {
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
	// сoздаём конфиг для логера
	cfg := zap.NewProductionEncoderConfig()
	// меняем то как будет выглябеть время при логировании
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15.04.05.000000")
	// создаём подЪядра логера
	encoder := zapcore.NewConsoleEncoder(cfg)
	service := os.Getenv("SERVICE_NAME")
	if service == "" {
		service = "my-app"
	}

	clickhouseCore, err := NewDataBaseCore(
		repository,
		context.Background(),
		service,
		lvl.Level(),
	)
	if err != nil {
		fmt.Printf("ClickHouse core init failed: %v\n", err)
	}
	// Единый Core, который пишет во всё
	multi := &MultiCore{
		console:    zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), lvl),
		file:       zapcore.NewCore(encoder, zapcore.AddSync(logFile), lvl),
		clickhouse: clickhouseCore,
		fields:     make(map[string]any),
	}
	//инициализируем логер
	logger := zap.New(
		multi,
		zap.AddCaller(),                       // параметр логирует место из которого произведён сам лог
		zap.AddStacktrace(zapcore.ErrorLevel), //параметр говорит что стек вызовов надо показывать для уровней Error и выше
	)
	return logger, logFile.Close, nil
}
