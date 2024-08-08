package config

type ConfigOptions struct {
	prefix     string
	configFile string
}

type ConfigOption func(*ConfigOptions)

func WithEnv(prefix string) ConfigOption {
	return func(c *ConfigOptions) {
		c.prefix = prefix
	}
}

func WithConfigfile(path string) ConfigOption {
	return func(c *ConfigOptions) {
		c.configFile = path
	}
}

func buildOptions(opts ...ConfigOption) *ConfigOptions {
	// 设置默认值
	options := &ConfigOptions{}

	// 应用选项函数进行定制化设置
	for _, opt := range opts {
		opt(options)
	}

	return options
}
