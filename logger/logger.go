// github.com/bigwhite/experiments/tree/master/uber-zap-advanced-usage/demo1/pkg/log/log.go
package log

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	InfoLevel   Level = zap.InfoLevel   // 0, default level
	WarnLevel   Level = zap.WarnLevel   // 1
	ErrorLevel  Level = zap.ErrorLevel  // 2
	DPanicLevel Level = zap.DPanicLevel // 3, used in development log
	// PanicLevel logs a message, then panics
	PanicLevel Level = zap.PanicLevel // 4
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel Level = zap.FatalLevel // 5
	DebugLevel Level = zap.DebugLevel // -1
)

type Field = zap.Field

func (l *Logger) Debugf(template string, args ...any) {
	var fields []zap.Field
	for _, a := range args {
		f := zap.Any(template, a)
		fields = append(fields, f)
	}
	l.l.Debug(template, fields...)
}

func (l *Logger) Infof(template string, args ...any) {
	var fields []zap.Field
	for _, a := range args {
		f := zap.Any(template, a)
		fields = append(fields, f)
	}
	l.l.Info(template, fields...)
}

func (l *Logger) Warnf(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Errorf(template string, args ...any) {
	var fields []zap.Field
	for _, a := range args {
		f := zap.Any(template, a)
		fields = append(fields, f)
	}
	l.l.Error(template, fields...)
}

func (l *Logger) DPanicf(msg string, fields ...Field) {
	l.l.DPanic(msg, fields...)
}

func (l *Logger) Panicf(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}

func (l *Logger) Fatalf(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

// function variables for all field types
// in github.com/uber-go/zap/field.go

var (
	Skip       = zap.Skip
	Binary     = zap.Binary
	Bool       = zap.Bool
	Boolp      = zap.Boolp
	ByteString = zap.ByteString

	Float64   = zap.Float64
	Float64p  = zap.Float64p
	Float32   = zap.Float32
	Float32p  = zap.Float32p
	Durationp = zap.Durationp

	Any = zap.Any

	Infof   = std.Infof
	Warnf   = std.Warnf
	Errorf  = std.Errorf
	DPanicf = std.DPanicf
	Panicf  = std.Panicf
	Fatalf  = std.Fatalf
	Debugf  = std.Debugf
)

// not safe for concurrent use
func ResetDefault(l *Logger) {
	std = l
	Infof = std.Infof
	Warnf = std.Warnf
	Errorf = std.Errorf
	DPanicf = std.DPanicf
	Panicf = std.Panicf
	Fatalf = std.Fatalf
	Debugf = std.Debugf
}

type Logger struct {
	l     *zap.Logger // zap ensure that zap.Logger is safe for concurrent use
	level Level
}

var std = New(os.Stderr, InfoLevel)

func Default() *Logger {
	return std
}

// New create a new logger (not support log rotating).
func New(writer io.Writer, level Level) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	cfg := zap.NewProductionConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(writer),
		zapcore.Level(level),
	)
	logger := &Logger{
		l:     zap.New(core),
		level: level,
	}
	return logger
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func Sync() error {
	if std != nil {
		return std.Sync()
	}
	return nil
}
