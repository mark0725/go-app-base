package app

import (
	base_config "github.com/mark0725/go-app-base/config"
)

type ApplicationOptions struct {
	appName    string
	version    string
	envPrefix  string
	configFile string
}

func NewApplicationOptions() *ApplicationOptions {
	return &ApplicationOptions{
		appName:    "app",
		version:    "0.0.1",
		envPrefix:  "APP",
		configFile: "config.yaml",
	}
}

func (opiton *ApplicationOptions) AppName(appName string) *ApplicationOptions {
	opiton.appName = appName
	return opiton
}
func (opiton *ApplicationOptions) Version(version string) *ApplicationOptions {
	opiton.version = version
	return opiton
}

func (opiton *ApplicationOptions) EnvPrefix(envPrefix string) *ApplicationOptions {
	opiton.envPrefix = envPrefix
	return opiton
}

func (opiton *ApplicationOptions) ConfigFile(configFile string) *ApplicationOptions {
	opiton.configFile = configFile
	return opiton
}

func (opiton *ApplicationOptions) build() *base_config.AppConfig {
	return &base_config.AppConfig{}
}
