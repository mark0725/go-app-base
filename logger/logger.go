package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var g_Loggers map[string]*log.Logger = make(map[string]*log.Logger)
var g_defaultLogConf = LoggerConfig{Level: "info"}
var g_defaultLoggers = make(map[string]*log.Logger)
var g_LoggerChangedHandles = make(map[string]func())

func init() {
	logger := CreateLogger("default", &g_defaultLogConf)
	g_Loggers["default"] = logger
}

type LogLevelType string

const (
	LogLevelTrace LogLevelType = "trace"
	LogLevelDebug LogLevelType = "debug"
	LogLevelInfo  LogLevelType = "info"
	LogLevelWarn  LogLevelType = "warn"
	LogLevelError LogLevelType = "error"
	LogLevelFatal LogLevelType = "fatal"
	LogLevelPanic LogLevelType = "panic"
)

func SetLogLevel(level LogLevelType) {
	logLevel := log.InfoLevel

	if lvl, err := log.ParseLevel(string(level)); err == nil {
		logLevel = lvl
	}

	log.SetLevel(logLevel)
	g_defaultLogConf.Level = string(level)
	for _, logger := range g_defaultLoggers {
		logger.SetLevel(logLevel)
	}
}

func LoggerInit(config *LogConfig) {
	log.SetFormatter(&LoggerFormatter{})
	//log.AddHook(FileLineHook{})
	logLevel := log.InfoLevel
	g_defaultLogConf = config.LoggerConfig

	if lvl, err := log.ParseLevel(config.Level); err == nil {
		logLevel = lvl
	}
	log.SetLevel(logLevel)

	initLoggers(config.Loggers)

	for name, logger := range g_defaultLoggers {
		configLogger(name, &g_defaultLogConf, logger)
	}
}

func OnLoggerChanged(name string, f func()) {
	g_LoggerChangedHandles[name] = f
}

func GetLogger(name string) *log.Logger {
	if logger, exists := g_Loggers[name]; exists {
		return logger
	}

	logger := CreateLogger(name, &g_defaultLogConf)
	g_Loggers[name] = logger
	g_defaultLoggers[name] = logger

	return logger
}

func initLoggers(logsConfig map[string]LoggerConfig) {
	keys := make([]string, 0, len(logsConfig))
	for key := range logsConfig {
		keys = append(keys, key)
	}
	sort.Strings(keys) // 对键进行排序

	for _, key := range keys {
		loggerConf := logsConfig[key]
		if loggerConf.File == "" {
			loggerConf.File = g_defaultLogConf.File
		}
		if loggerConf.Format == "" {
			loggerConf.Format = g_defaultLogConf.Format
		}
		if loggerConf.Level == "" {
			loggerConf.Level = g_defaultLogConf.Level
		}

		if l, exists := g_Loggers[key]; exists {
			configLogger(key, &loggerConf, l)
			g_Loggers[key] = l
			delete(g_defaultLoggers, key)
			if h, exists := g_LoggerChangedHandles[key]; exists {
				h()
			}
		} else {
			logger := CreateLogger(key, &loggerConf)
			g_Loggers[key] = logger
		}

		//fmt.Printf("Key: %s, Value: %s\n", key, appConfig.Logs[key].Level)
	}
}

func CreateLogger(name string, conf *LoggerConfig) *log.Logger {
	logger := log.New()
	configLogger(name, conf, logger)
	return logger
}
func configLogger(name string, conf *LoggerConfig, logger *log.Logger) {

	logLevel := log.InfoLevel

	if lvl, err := log.ParseLevel(conf.Level); err == nil {
		logLevel = lvl
	}

	logger.SetLevel(logLevel)
	logger.SetFormatter(&LoggerFormatter{name: name})
	//logger.AddHook(FileLineHook{})
	logger.SetReportCaller(true)
	if conf.File != "" {
		logRotate := &lumberjack.Logger{
			Filename:   conf.File,       // 日志文件路径
			MaxSize:    conf.MaxAge,     // 每个日志文件最大10MB
			MaxBackups: conf.MaxBackups, // 保留的备份文件个数
			MaxAge:     conf.MaxAge,     // 保留备份文件的最大天数
			Compress:   conf.Compress,   // 是否压缩备份文件
		}

		mw := io.MultiWriter(os.Stdout, logRotate)
		logger.SetOutput(mw)

	}

}

// LoggerFormatter is a custom log formatter.
type LoggerFormatter struct {
	name string
}

// Format formats the log entry.
func (f *LoggerFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b bytes.Buffer

	// Add the timestamp
	b.WriteString("[")
	b.WriteString(entry.Time.Format("2006-01-02T15:04:05"))
	b.WriteString("] [")

	// Add the log level
	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteString("] ")

	if f.name != "" {
		b.WriteString("[")
		b.WriteString(f.name)
		b.WriteString("] ")
	}

	if entry.Caller != nil {
		fnParts := strings.Split(entry.Caller.Function, "/")
		fnName := fnParts[len(fnParts)-1]
		b.WriteString(fnName)
		b.WriteString(" ")
	}

	// Add the fields
	for key, value := range entry.Data {
		b.WriteString(fmt.Sprintf("%s=%v ", key, value))
	}

	// Add the message
	b.WriteString(entry.Message)
	b.WriteString("\n")

	return b.Bytes(), nil
}
