package app

import (
	"context"
	"encoding/json"

	base_config "github.com/mark0725/go-app-base/config"
	base_log "github.com/mark0725/go-app-base/logger"
)

var logger = base_log.GetLogger("app")
var g_AppConfigOrig map[string]any

func GetAppConfigOrig() map[string]any {
	return g_AppConfigOrig
}

type Appliction struct {
	AppConfig base_config.IAppConfig
	Options   *ApplicationOptions
}

func NewApplication(appConfig base_config.IAppConfig, options *ApplicationOptions) *Appliction {
	if options == nil {
		options = NewApplicationOptions()
	}

	return &Appliction{
		Options:   options,
		AppConfig: appConfig,
	}
}

func ConfigInit(envPrefix string, configFile string, appConfig any) {
	config, err := base_config.LoadConfig(
		base_config.WithEnv(envPrefix),
		base_config.WithConfigfile(configFile),
	)
	if err != nil {
		logger.Error("Error loading config:", err)
		return
	}

	logger.Tracef("app config: %#v\n", config)
	g_AppConfigOrig = config

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		logger.Error("Error marshalling to JSON:", err)
		return
	}

	logger.Tracef("app config json: %s\n", jsonData)

	errjson := json.Unmarshal([]byte(jsonData), appConfig)
	if errjson != nil {
		logger.Error("Error unmarshalling JSON:", err)
		return
	}

	logger.Debugf("app config: %#v\n", appConfig)
}

func (app *Appliction) AppInit() error {
	//base_log.SetLogLevel(base_log.LogLevelInfo)
	logger.Info("App init ...")
	ConfigInit(app.Options.envPrefix, app.Options.configFile, app.AppConfig)

	logCfg := &app.AppConfig.GetAppConfig().Log
	base_log.LoggerInit(logCfg)

	err := InitializeModules(app.AppConfig)
	if err != nil {
		logger.Errorf("initialize error: %v", err)
		return err
	}

	logger.Info("App init end.")
	return nil
}

func (app *Appliction) Run(ctx context.Context) error {
	// ctx := context.Background()
	// defer ctx.Done()
	logger.Info("App starting ...")
	modules := GetReadyModules()
	logger.Tracef("Ready modules %v", modules)
	for _, mod := range modules {
		//logger.Tracef("Module %s", mod)
		serves := GetServeNamesByModule(mod)
		for _, serve := range serves {
			logger.Tracef("Module %s serve %s", mod, serve)
			if s, ok := app.AppConfig.GetAppConfig().Serves[mod]; ok {
				if c, ok := s[serve]; ok {
					if !c.Enable {
						logger.Infof("Skip %s serve: %s", mod, serve)
						continue
					}
				}
			}

			logger.Debugf("Starting %s serve: %s", mod, serve)
			if err := StartServe(ctx, mod, serve); err != nil {
				logger.Errorf("Start %s serve: %s err: %v", mod, serve, err)
				return err
			}
			logger.Debugf("%s serve: %s Started", mod, serve)
		}
	}

	logger.Info("App started")

	return nil
}
