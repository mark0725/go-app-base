package logger

type LogConfig struct {
	Level   string                  `json:"level"`
	File    string                  `json:"file"`
	Format  string                  `json:"format"`
	Loggers map[string]LoggerConfig `json:"loggers"`
}

type LoggerConfig struct {
	Level  string `json:"level"`
	File   string `json:"file"`
	Format string `json:"format"`
}
