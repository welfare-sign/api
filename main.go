package main

import (
	"context"
	http2 "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"welfare-sign/internal/pkg/config"
	"welfare-sign/internal/pkg/log"
	"welfare-sign/internal/pkg/util"
	"welfare-sign/internal/server"
	"welfare-sign/internal/service"
	"welfare-sign/internal/task"
)

// @title 福利签API文档
// @version 1.0
// @description 福利签API文档
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	svc := service.New()
	httpSrv := server.New(svc)
	quit := make(chan os.Signal, 1)
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http2.ErrServerClosed {
			log.Error(context.Background(), "welfare-sign start error", zap.Error(err))
			quit <- os.Kill
		}
	}()
	log.Info(context.Background(), "welfare-sign start", zap.String("server_start_time", util.TimeFormat(time.Now())))

	if viper.GetBool(config.KeyTaskEnable) {
		// 触发定时任务
		go task.Run(svc)
	}

	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Error(context.Background(), "welfare-sign closing error", zap.Error(err))
	}
	svc.Close()
	log.Info(context.Background(), "welfare-sign exit", zap.String("server_stop_time", util.TimeFormat(time.Now())))
}
