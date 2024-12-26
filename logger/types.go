package logger

type LogConfig struct {
	Level     string                    `json:"level"`
	Default   LoggerConfig              `json:"default"`
	Appenders map[string]AppenderConfig `json:"appenders"`
	Loggers   map[string]LoggerConfig   `json:"loggers"`
}

type AppenderConfig struct {
	Type       string `json:"type"`
	Level      string `json:"level"`
	Path       string `json:"path"`
	Format     string `json:"format"`
	MaxSize    int    `json:"maxsize,omitempty"`    // 每个日志文件最大10MB
	MaxBackups int    `json:"maxbackups,omitempty"` // 保留的备份文件个数
	MaxAge     int    `json:"maxage,omitempty"`     // 保留备份文件的最大天数
	Compress   bool   `json:"compress,omitempty"`   // 是否压缩备份文件
}

type LoggerConfig struct {
	Level     string   `json:"level"`
	Appenders []string `json:"appenders"`
}
