package logger

type LogConfig struct {
	LoggerConfig
	Loggers map[string]LoggerConfig `json:"loggers"`
}

type LoggerConfig struct {
	Level      string `json:"level"`
	File       string `json:"file"`
	Format     string `json:"format"`
	MaxSize    int    `json:"maxsize,omitempty"`    // 每个日志文件最大10MB
	MaxBackups int    `json:"maxbackups,omitempty"` // 保留的备份文件个数
	MaxAge     int    `json:"maxage,omitempty"`     // 保留备份文件的最大天数
	Compress   bool   `json:"compress,omitempty"`   // 是否压缩备份文件
}
