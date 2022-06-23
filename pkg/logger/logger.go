package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	// DefaultLevel the default log level
	DefaultLevel = zapcore.InfoLevel

	// DefaultTimeLayout the default time layout;
	DefaultTimeLayout = time.RFC3339
)

// option 模式
type Option func(*option)

type option struct {
	level          zapcore.Level
	fields         map[string]string
	file           io.Writer
	timeLayout     string
	disableConsole bool
}

func WithDebugLevel() Option {
	return func(o *option) {
		o.level = zapcore.DebugLevel
	}
}

func WithInfoLevel() Option {
	return func(o *option) {
		o.level = zapcore.InfoLevel
	}
}

// WithWarnLevel only greater than 'level' will output
func WithWarnLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.WarnLevel
	}
}

// WithErrorLevel only greater than 'level' will output
func WithErrorLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.ErrorLevel
	}
}

func WithDisableConsole() Option {
	return func(o *option) {
		o.disableConsole = true
	}
}

func WithTimeLayout(timeLayout string) Option {
	return func(o *option) {
		o.timeLayout = timeLayout
	}
}

func WithFileRotationP(file string) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}
	return func(o *option) {
		o.file = &lumberjack.Logger{
			Filename:   file, // 文件路径
			MaxSize:    128,  // 单个文件最大尺寸 默认单位M
			MaxBackups: 300,  // 最大备份数
			MaxAge:     30,   // 最大时间默认单位：天
			LocalTime:  true, // 使用本地时间
			Compress:   true, // 是否有所 默认不压缩
		}
	}
}

func WithField(key, value string) Option {
	return func(o *option) {
		o.fields[key] = value
	}
}

func NewJSONLogger(opts ...Option) (*zap.Logger, error) {
	opt := &option{
		level:  DefaultLevel,
		fields: make(map[string]string),
	}

	for _, f := range opts {
		f(opt)
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "function",
		StacktraceKey:  "Stacktrace",
		SkipLineEnding: false,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format(DefaultTimeLayout))
		},
		EncodeDuration:      zapcore.MillisDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder, // 全路径编码器
		EncodeName:          nil,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "",
	}
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= opt.level && level < zapcore.ErrorLevel
	})

	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= opt.level && level >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee()

	if !opt.disableConsole {
		core = zapcore.NewTee(
			zapcore.NewCore(jsonEncoder,
				zapcore.NewMultiWriteSyncer(zapcore.Lock(os.Stdout)),
				lowPriority,
			),
			zapcore.NewCore(jsonEncoder,
				zapcore.NewMultiWriteSyncer(zapcore.Lock(os.Stderr)),
				highPriority,
			),
		)
	}

	if opt.file != nil {
		core = zapcore.NewTee(core,
			zapcore.NewCore(jsonEncoder,
				zapcore.AddSync(opt.file),
				zap.LevelEnablerFunc(func(level zapcore.Level) bool {
					return level >= opt.level
				}),
			),
		)
	}

	logger := zap.New(core, zap.AddCaller(), zap.ErrorOutput(zapcore.Lock(os.Stderr)))

	for key, value := range opt.fields {
		logger = logger.WithOptions(zap.Fields(zapcore.Field{Key: key, Type: zapcore.StringType, String: value}))
	}
	return logger, nil
}
