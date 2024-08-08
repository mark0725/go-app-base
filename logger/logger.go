package logger

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var g_Loggers map[string]*log.Logger = make(map[string]*log.Logger)

// func init() {
// 	g_Loggers = make(map[string]*log.Logger)
// }

func LoggerInit(level string, configs map[string]LogConfig) {
	//log.SetFormatter(&MyFormatter{})
	//log.AddHook(FileLineHook{})
	logLevel := log.InfoLevel

	if lvl, err := log.ParseLevel(level); err == nil {
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

	initLoggers(level, configs)

	// if _, exists := g_Loggers["default"]; !exists {
	// 	logger := CreateLogger("default", log.InfoLevel)
	// 	g_Loggers["default"] = logger
	// }

	// log.Debug("This is a debug message") // 不会打印，因为默认级别是Info
	// log.Info("This is an info message")
	// log.Warn("This is a warning message")
	// log.Error("This is an error message")
}

func GetLogger(name string) *log.Logger {

	if logger, exists := g_Loggers[name]; exists {
		return logger
	}

	if logger, exists := g_Loggers["default"]; exists {
		return logger
	}

	//only for test
	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	g_Loggers["default"] = logger

	return logger
}

func initLoggers(level string, logsConfig map[string]LogConfig) {

	defaultLogConf := LogConfig{
		Level: level,
	}

	logger := CreateLogger("default", defaultLogConf)
	g_Loggers["default"] = logger

	keys := make([]string, 0, len(logsConfig))
	for key := range logsConfig {
		keys = append(keys, key)
	}
	sort.Strings(keys) // 对键进行排序

	for _, key := range keys {

		logger := CreateLogger(key, logsConfig[key])
		g_Loggers[key] = logger

		//fmt.Printf("Key: %s, Value: %s\n", key, appConfig.Logs[key].Level)
	}
}

func CreateLogger(name string, conf LogConfig) *log.Logger {

	logLevel := log.InfoLevel

	if lvl, err := log.ParseLevel(conf.Level); err == nil {
		logLevel = lvl
	}

	logger := log.New()
	logger.SetLevel(logLevel)
	//logger.SetFormatter(&MyFormatter{})
	//logger.AddHook(FileLineHook{})
	logger.SetReportCaller(true)
	logger.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// 这里可以自定义包名和文件名的格式
			//filename := filepath.Base(f.File)
			funcName := f.Function
			return funcName, ""
		},
		FullTimestamp: true,
	})

	return logger
}

// MyFormatter is a custom log formatter.
type MyFormatter struct{}

// Format formats the log entry.
func (f *MyFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b bytes.Buffer

	// Add the timestamp
	b.WriteString(entry.Time.Format(time.RFC3339))
	b.WriteString(" ")

	// Add the log level
	b.WriteString(entry.Level.String())
	b.WriteString(" ")

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
