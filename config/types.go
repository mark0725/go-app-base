package config

import (
	base_logger "github.com/mark0725/go-app-base/logger"
)

type OrgConfig struct {
	OrgId string `json:"org_id"`
}

type AppConfigApp struct {
	AppHome string `json:"app_home"`
}

type AppServeConfig struct {
	Enable bool `json:"enable"`
}

type AppConfig struct {
	Org    OrgConfig                            `json:"org"`
	App    AppConfigApp                         `json:"app"`
	Log    base_logger.LogConfig                `json:"log"`
	Serves map[string]map[string]AppServeConfig `json:"serves"`
}

type IAppConfig interface {
	GetAppConfig() *AppConfig
}
