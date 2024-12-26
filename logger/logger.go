package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var g_Loggers map[string]*log.Logger = make(map[string]*log.Logger)
var g_LogAppenders map[string]AppenderConfig = make(map[string]AppenderConfig)
var g_defaultLogConf = LoggerConfig{Level: "info", Appenders: []string{"console"}}
var g_defaultLoggers = make(map[string]*log.Logger)
var g_LoggerChangedHandles = make(map[string]func())

func init() {
	g_LogAppenders["console"] = AppenderConfig{
		Type:  "console",
		Level: "trace",
	}

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

	appenders := g_defaultLogConf.Appenders
	if len(config.Default.Appenders) > 0 {
		appenders = config.Default.Appenders
	}
	g_defaultLogConf = LoggerConfig{
		Level:     config.Level,
		Appenders: appenders,
	}

	for name, appender := range config.Appenders {
		g_LogAppenders[name] = appender
	}

	if lvl, err := log.ParseLevel(config.Level); err == nil {
		logLevel = lvl
	}
	log.SetLevel(logLevel)

	initLoggers(config.Loggers)
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

	for key, l := range g_Loggers {
		if conf, exists := logsConfig[key]; exists {
			if conf.Level == "" {
				conf.Level = g_defaultLogConf.Level
			}
			if len(conf.Appenders) == 0 {
				conf.Appenders = g_defaultLogConf.Appenders
			}

			configLogger(key, &conf, l)
			delete(g_defaultLoggers, key)
		} else {
			configLogger(key, &g_defaultLogConf, l)
		}
	}

	for key, conf := range logsConfig {
		if _, exists := g_Loggers[key]; !exists {
			if conf.Level == "" {
				conf.Level = g_defaultLogConf.Level
			}
			if len(conf.Appenders) == 0 {
				conf.Appenders = g_defaultLogConf.Appenders
			}

			logger := CreateLogger(key, &conf)
			g_Loggers[key] = logger
		}
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
	//logger.SetFormatter(&LoggerFormatter{name: name})
	logger.SetReportCaller(true)
	logger.Out = io.Discard
	logger.Hooks = make(log.LevelHooks)
	appenders := conf.Appenders
	if len(appenders) == 0 {
		appenders = g_defaultLogConf.Appenders
	}
	for _, appender := range appenders {
		if appenderConfig, exists := g_LogAppenders[appender]; exists {
			if appenderConfig.Level == "" {
				appenderConfig.Level = g_defaultLogConf.Level
			}

			appendLogLevel := logLevel
			if lvl, err := log.ParseLevel(appenderConfig.Level); err == nil {
				if lvl < appendLogLevel {
					appendLogLevel = lvl
				}
			}

			logger.AddHook(NewLogAppender(name, appendLogLevel, &LoggerFormatter{name: name}, &appenderConfig))

		}
	}
}

type LogAppender struct {
	name      string
	conf      *AppenderConfig
	levels    []log.Level
	formatter log.Formatter
	writer    io.Writer
}

func NewLogAppender(name string, level log.Level, formatter log.Formatter, conf *AppenderConfig) *LogAppender {
	levels := []log.Level{}
	for i := 0; i <= int(level); i++ {
		levels = append(levels, log.Level(i))
	}

	logWriter := io.Writer(os.Stderr)
	switch conf.Type {
	case "file":
		logWriter = &lumberjack.Logger{
			Filename:   conf.Path,       // 日志文件路径
			MaxSize:    conf.MaxAge,     // 每个日志文件最大10MB
			MaxBackups: conf.MaxBackups, // 保留的备份文件个数
			MaxAge:     conf.MaxAge,     // 保留备份文件的最大天数
			Compress:   conf.Compress,   // 是否压缩备份文件
		}
	case "module":
		logFile := filepath.Join(conf.Path, fmt.Sprintf("%s.log", name))
		logWriter = &lumberjack.Logger{
			Filename:   logFile,         // 日志文件路径
			MaxSize:    conf.MaxAge,     // 每个日志文件最大10MB
			MaxBackups: conf.MaxBackups, // 保留的备份文件个数
			MaxAge:     conf.MaxAge,     // 保留备份文件的最大天数
			Compress:   conf.Compress,   // 是否压缩备份文件
		}
	case "console":
		logWriter = os.Stdout
	default:

	}

	return &LogAppender{
		name:      name,
		conf:      conf,
		levels:    levels,
		formatter: formatter,
		writer:    logWriter,
	}
}

func (hook *LogAppender) Fire(entry *log.Entry) error {
	b, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.writer.Write(b)
	return err
}

func (hook *LogAppender) Levels() []log.Level {
	return hook.levels
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
