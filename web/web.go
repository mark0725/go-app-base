package web

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var endPoints = map[string]func(group string, r *gin.RouterGroup){}

func EndPointRegister(module string, f func(group string, r *gin.RouterGroup)) {
	endPoints[module] = f
}

func StartWebServe(ctx context.Context, conf *WebConfig) error {
	for _, server := range conf.Servers {
		router := gin.Default()

		writer := WebLogWriter{
			buf: bytes.NewBuffer([]byte{}),
		}

		if server.AccessLog != "" {
			accessLogFile, err := os.OpenFile(server.AccessLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				logger.Errorf("error opening access log file: %s error: %v", server.AccessLog, err)
			}
			writer.accessLog = accessLogFile
		}

		if server.ErrorLog != "" {
			errorLogFile, err := os.OpenFile(server.ErrorLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				logger.Errorf("error opening error log file: %s error: %v", server.ErrorLog, err)
			}
			writer.errorLog = errorLogFile
		}

		router.Use(gin.LoggerWithWriter(writer))
		router.Use(gin.RecoveryWithWriter(io.MultiWriter(writer, os.Stderr)))

		if server.StaticDir != "" {
			router.Static("/static", server.StaticDir)
		}

		if server.TemplateDir != "" {
			router.LoadHTMLGlob(server.TemplateDir)
		}

		for _, pt := range server.Endpoints {
			r := router.Group(pt.Path)
			endPoints[pt.Module](pt.Group, r)
		}

		router.NoRoute(func(c *gin.Context) {
			if c.Request.URL.Path != "/" {
				c.File("./static" + c.Request.URL.Path)
			} else {
				c.File("./static/index.html")
			}
		})

		listenConfs, err := ParseListenConfig(server.Listen)
		if err != nil {
			logger.Errorf("ParseListenConfig %s error: %v", server.Listen, err)
			return err
		}
		for _, listenConf := range listenConfs {
			if err := StartListener(ctx, listenConf, router); err != nil {
				logger.Error("StartWeb error:", err)
				return err
			}
		}

	}

	return nil
}

// type WebServer struct {
// 	Config    *config.WebConfig
// 	Listeners []*WebListener
// 	Handlers  *gin.Engine
// 	ctx       context.Context
// }

// type WebListener struct {
// 	httpServers []*http.Server
// 	ListenConf  *ListenConfig
// }

func StartListener(ctx context.Context, listenConf *ListenConfig, r *gin.Engine) error {
	addr := fmt.Sprintf("%s:%d", listenConf.Address, listenConf.Port)
	logger.Info("Starting web server on:", addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 在 Goroutine 中启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("gin server start error: ", err)
		}
	}()

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Server Shutdown:", err)
		} else {
			logger.Info("Server exiting: ", addr)
		}
	}()

	return nil

}
