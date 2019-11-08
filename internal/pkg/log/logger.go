package log

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tidwall/pretty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"welfare-sign/internal/pkg/config"
)

// log constant
const (
	RequestIDKey       = "req_id"
	RequestClientIPKey = "req_ip"
)

// Logger logger
type Logger interface {
	// Fatal 级别
	Fatal(ctx context.Context, title string, field ...zap.Field)
	// Error 级别
	Error(ctx context.Context, title string, field ...zap.Field)
	// Warn 级别
	Warn(ctx context.Context, title string, field ...zap.Field)
	// Info 级别
	Info(ctx context.Context, title string, field ...zap.Field)
	// Debug 级别
	Debug(ctx context.Context, title string, field ...zap.Field)
	Close()
}

// NullLogger null logger
type NullLogger struct{}

// Fatal .
func (l NullLogger) Fatal(ctx context.Context, title string, field ...zap.Field) {}

// Error .
func (l NullLogger) Error(ctx context.Context, title string, field ...zap.Field) {}

// Warn .
func (l NullLogger) Warn(ctx context.Context, title string, field ...zap.Field) {}

// Info .
func (l NullLogger) Info(ctx context.Context, title string, field ...zap.Field) {}

// Debug .
func (l NullLogger) Debug(ctx context.Context, title string, field ...zap.Field) {}

// Close .
func (l NullLogger) Close() {}

var (
	// DefaultLogger 默认普通日志记录器
	DefaultLogger Logger
	// logFieldsPool 日志字段池
	logFieldsPool sync.Pool
	// jsonEncoder json格式编码
	jsonEncoder = zapcore.NewJSONEncoder(encoder)

	entry = zapcore.Entry{}
)

// encoder zap格式化配置
var encoder = zapcore.EncoderConfig{
	TimeKey:        "timestamp",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "line",
	MessageKey:     "title",
	StacktraceKey:  "detail",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     EpochMillisTimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.FullCallerEncoder,
}

// EpochMillisTimeEncoder 日志时间字段 毫秒级单位整形
func EpochMillisTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	nano := t.UnixNano()
	millis := int64(nano) / int64(time.Millisecond)
	enc.AppendInt64(millis)
}

// logger 日志中心
type logger struct {
	FileWriter *lumberjack.Logger //普通日志 文件记录
	Logger     *zap.Logger        //普通日志记录器
}

// GetFieldsWithPool .
func GetFieldsWithPool() (m []zap.Field) {
	o, _ := logFieldsPool.Get().(*[]zap.Field)
	return *o
}

// ResetFields .
func ResetFields(m []zap.Field) {
	m = m[0:1]
	logFieldsPool.Put(&m)
}

func (l *logger) log(ctx context.Context, level zapcore.Level, msg string, f ...zap.Field) {
	fields := GetFieldsWithPool()
	defer ResetFields(fields)

	clientIP := ctx.Value(RequestClientIPKey)
	if clientIP != nil {
		if ip, ok := clientIP.(string); ok {
			fields = append(fields, zap.String("client_ip", ip))
		}
	}
	_reqID := ""
	reqID := ctx.Value(RequestIDKey)
	if reqID != nil {
		_reqID = reqID.(string)
	}

	bf := bytes.Buffer{}
	bf.WriteByte('{')
	f = append(f, zap.String(RequestIDKey, _reqID))

	jbf, err := jsonEncoder.EncodeEntry(entry, f)
	if err != nil {
		log.Println(errors.WithMessage(err, "log_json_marshal"))
	}
	bf.Write(jbf.Bytes())

	bts := bf.Bytes()
	fields = append(fields, zap.ByteString("context", bts[1:len(bts)-1]))
	if len(msg) < 5 {
		msg += "****"
	}
	switch level {
	case zap.FatalLevel:
		l.Logger.Fatal(msg, fields...)
	case zap.ErrorLevel:
		l.Logger.Error(msg, fields...)
	case zap.WarnLevel:
		l.Logger.Warn(msg, fields...)
	case zap.InfoLevel:
		l.Logger.Info(msg, fields...)
	default:
		l.Logger.Debug(msg, fields...)
	}
}

// Fatal 级别
func (l *logger) Fatal(ctx context.Context, title string, field ...zap.Field) {
	l.log(ctx, zap.FatalLevel, title, field...)
}

// Error 级别
func (l *logger) Error(ctx context.Context, title string, field ...zap.Field) {
	l.log(ctx, zap.ErrorLevel, title, field...)
}

// Warn 级别
func (l *logger) Warn(ctx context.Context, title string, field ...zap.Field) {
	l.log(ctx, zap.WarnLevel, title, field...)
}

// Info 级别
func (l *logger) Info(ctx context.Context, title string, field ...zap.Field) {
	l.log(ctx, zap.InfoLevel, title, field...)
}

// Debug 级别
func (l *logger) Debug(ctx context.Context, title string, field ...zap.Field) {
	l.log(ctx, zap.DebugLevel, title, field...)
}

// Close .
func (l *logger) Close() {
	if l.Logger != nil {
		err := l.Logger.Sync()
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	if !viper.GetBool(config.KeyLogEnable) {
		DefaultLogger = NullLogger{}
		return
	}

	logPath := viper.GetString(config.KeyLogPath)
	DefaultLogger = NewLogger(logPath)
}

// NewLogger 创建一个新的日志记录对象
func NewLogger(logpath string) Logger {
	project := viper.GetString(config.KeyLogProject)
	if logpath == "" || project == "" {
		panic(errors.New("invalid logpath or project"))
	}

	logFieldsPool = sync.Pool{
		New: func() interface{} {
			m := make([]zap.Field, 1, 6)
			m[0] = zap.String("project", project)
			return &m
		},
	}
	year, month, day := time.Now().Date()
	filename := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	l := &logger{
		FileWriter: &lumberjack.Logger{
			Filename:   fmt.Sprintf("%s%s.log", logpath, filename),
			MaxSize:    500,  // 日志文件最大容量，单位M
			MaxBackups: 10,   // 日志文件最大数量
			MaxAge:     28,   // 日志文件最大保存天数
			Compress:   true, // 是否压缩
		},
	}

	// 日志记录级别
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.FatalLevel
	})
	cores := make([]zapcore.Core, 0, 2)

	// 配置普通日志文件记录
	cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.AddSync(l.FileWriter), highPriority))
	// 控制台输出
	cores = append(cores, zapcore.NewCore(jsonEncoder, &JSONColorOut{}, highPriority))

	l.Logger = zap.New(zapcore.NewTee(cores...), zap.AddStacktrace(zap.ErrorLevel))

	return l
}

// JSONColorOut .
type JSONColorOut struct {
	count int
}

func (j *JSONColorOut) Write(p []byte) (n int, err error) {
	j.count++
	return fmt.Printf("%d\n%s\n", j.count, pretty.Color(pretty.Pretty(p), pretty.TerminalStyle))
}

// Sync .
func (j *JSONColorOut) Sync() error {
	return nil
}

// Fatal error
func Fatal(ctx context.Context, title string, field ...zap.Field) {
	DefaultLogger.Fatal(ctx, title, field...)
}

// Error error
func Error(ctx context.Context, title string, field ...zap.Field) {
	DefaultLogger.Error(ctx, title, field...)
}

// Warn 级别
func Warn(ctx context.Context, title string, field ...zap.Field) {
	DefaultLogger.Warn(ctx, title, field...)
}

// Info 级别
func Info(ctx context.Context, title string, field ...zap.Field) {
	DefaultLogger.Info(ctx, title, field...)
}

// Debug 级别
func Debug(ctx context.Context, title string, field ...zap.Field) {
	DefaultLogger.Debug(ctx, title, field...)
}
