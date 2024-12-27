package web

import (
	base_log "github.com/mark0725/go-app-base/logger"
)

const APP_MODULE_NAME string = "web"

type AppModule struct{}

var logger = base_log.GetLogger(APP_MODULE_NAME)
