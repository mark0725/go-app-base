package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

var g_Loggers map[string]*log.Logger = make(map[string]*log.Logger)
var g_defaultLogConf = LoggerConfig{Level: "info"}

func init() {
	logger := CreateLogger("default", &g_defaultLogConf)
	g_Loggers["default"] = logger
}

func LoggerInit(config *LogConfig) {
	log.SetFormatter(&LoggerFormatter{})
	//log.AddHook(FileLineHook{})
	logLevel := log.InfoLevel
	g_defaultLogConf.File = config.File
	g_defaultLogConf.Format = config.Format
	g_defaultLogConf.Level = config.Level

	if lvl, err := log.ParseLevel(config.Level); err == nil {
		logLevel = lvl
	}
	log.SetLevel(logLevel)

	//log.SetReportCaller(true)
	// log.SetFormatter(&log.TextFormatter{
	// 	CallerPrettyfier: func(f *runtime.Frame) (string, string) {
	// 		// 这里可以自定义包名和文件名的格式
	// 		filename := filepath.Base(f.File)
	// 		funcName := f.Function
	// 		return funcName, filename
	// 	},
	// 	FullTimestamp: true,
	// })

	initLoggers(config.Loggers)
}

func GetLogger(name string) *log.Logger {

	if logger, exists := g_Loggers[name]; exists {
		return logger
	}

	if l, exists := g_Loggers["default"]; exists {
		logger := log.New()
		logger.SetLevel(l.GetLevel())
		logger.SetFormatter(&LoggerFormatter{name: name})
		g_Loggers[name] = logger
		return logger
	}

	logger := CreateLogger(name, &g_defaultLogConf)
	g_Loggers[name] = logger

	return logger
}

func initLoggers(logsConfig map[string]LoggerConfig) {
	logger := CreateLogger("default", &g_defaultLogConf)
	g_Loggers["default"] = logger

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

		logger := CreateLogger(key, &loggerConf)
		g_Loggers[key] = logger

		//fmt.Printf("Key: %s, Value: %s\n", key, appConfig.Logs[key].Level)
	}
}

func CreateLogger(name string, conf *LoggerConfig) *log.Logger {

	logLevel := log.InfoLevel

	if lvl, err := log.ParseLevel(conf.Level); err == nil {
		logLevel = lvl
	}

	logger := log.New()
	logger.SetLevel(logLevel)
	logger.SetFormatter(&LoggerFormatter{name: name})
	//logger.AddHook(FileLineHook{})
	logger.SetReportCaller(true)
	if conf.File != "" {
		if file, err := os.OpenFile(conf.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
			log.Error("Failed to open log file: ", err)
		} else {
			mw := io.MultiWriter(os.Stdout, file)
			logger.SetOutput(mw)
		}

	}
	// logger.SetFormatter(&log.TextFormatter{
	// 	CallerPrettyfier: func(f *runtime.Frame) (string, string) {
	// 		// 这里可以自定义包名和文件名的格式
	// 		//filename := filepath.Base(f.File)
	// 		funcName := f.Function
	// 		return funcName, ""
	// 	},
	// 	FullTimestamp: true,
	// })

	return logger
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

type FileLineHook struct{}

func (hook FileLineHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook FileLineHook) Fire(entry *log.Entry) error {
	// Skip 4 levels to get the caller of the function which we're logging
	pc, file, line, ok := runtime.Caller(4)
	if !ok {
		return nil
	}

	fn := runtime.FuncForPC(pc)
	// Extract just the package and function name
	fnName := fn.Name()
	fnParts := strings.Split(fnName, "/")
	fnName = fnParts[len(fnParts)-1]

	// Trim down the file path to the last two segments for brevity
	fileParts := strings.Split(file, "/")
	file = strings.Join(fileParts[len(fileParts)-2:], "/")

	// Add to the entry's data
	entry.Data["file"] = file
	entry.Data["line"] = line
	entry.Data["func"] = fnName

	return nil
}
