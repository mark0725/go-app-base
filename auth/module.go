package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	base_app "github.com/mark0725/go-app-base/app"
	base_config "github.com/mark0725/go-app-base/config"
	base_log "github.com/mark0725/go-app-base/logger"
	base_web "github.com/mark0725/go-app-base/web"
)

const APP_MODULE_NAME string = "base-auth"

type AppModule struct{}

type AppModuleConfig struct {
	Org base_config.OrgConfig    `json:"org"`
	App base_config.AppConfigApp `json:"app"`
	Log base_log.LogConfig       `json:"log"`
}

var logger = base_log.GetLogger(APP_MODULE_NAME)
var g_appConfig *AppModuleConfig

func init() {
	base_app.AppModuleRegister(APP_MODULE_NAME, &AppModule{}, []string{},
		base_app.AppModuleRegisterOptionWithConfigType(&AppModuleConfig{}),
	)
}

func (m *AppModule) Init(appConfig interface{}, depends []string) error {
	if v, ok := appConfig.(*AppModuleConfig); !ok {
		logger.Error("invalid app config")
		return errors.New("invalid app config")
	} else {
		g_appConfig = v
	}

	logger.Tracef("AppModule %s init ... ", APP_MODULE_NAME)

	base_web.EndPointRegister(APP_MODULE_NAME, func(group string, r *gin.RouterGroup) {
		InitRoutes(group, r)
	})

	return nil
}
