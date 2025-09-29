package web

import (
	"context"

	config "github.com/mark0725/go-app-base/config"

	base_app "github.com/mark0725/go-app-base/app"
	base_log "github.com/mark0725/go-app-base/logger"
)

const APP_MODULE_NAME string = "web"

type AppModule struct{}

var logger = base_log.GetLogger(APP_MODULE_NAME)
var g_appConfig *config.AppConfig
var g_webServe WebServe

func init() {
	base_app.AppModuleRegister(APP_MODULE_NAME, &AppModule{}, []string{})
}

func (m *AppModule) Init(appConfig interface{}, depends []string) error {
	g_appConfig = appConfig.(*config.AppConfig)

	logger.Tracef("AppModule %s init ... ", APP_MODULE_NAME)

	if err := InitWeb(g_appConfig); err != nil {
		logger.Error("InitWeb error:", err)
		return err
	}

	if err := base_app.RegisterServe("web", APP_MODULE_NAME, &g_webServe); err != nil {
		logger.Error("RegisterServe error:", err)
		return err
	}

	return nil
}

func InitWeb(appConfig *config.AppConfig) error {

	logger.Debug("Web initialized successfully")

	return nil
}

type WebServe struct {
	ready bool
	done  <-chan struct{}
}

func (s *WebServe) Done() <-chan struct{} {
	return s.done
}
func (s *WebServe) Start(ctx context.Context) error {
	// if err := StartWebServe(ctx, &g_appConfig.Web); err != nil {
	// 	logger.Error("StartWeb error:", err)
	// 	return err
	// }
	s.ready = true
	return nil
}

func (s *WebServe) Ready() bool {
	return s.ready
}
func (s *WebServe) Stop() error {
	return nil
}
