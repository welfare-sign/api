package main

import (
	"context"
	http2 "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"welfare-sign/internal/pkg/log"
	"welfare-sign/internal/pkg/util"
	"welfare-sign/internal/server/http"
	"welfare-sign/internal/service"
)

// @title 福利签API文档
// @version 1.0
// @description 福利签API文档
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /v1
func main() {
	svc := service.New()
	httpSrv := http.New(svc)
	quit := make(chan os.Signal, 1)
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http2.ErrServerClosed {
			log.Error(context.Background(), "welfare-sign start error", zap.Error(err))
			quit <- os.Kill
		}
	}()
	log.Info(context.Background(), "welfare-sign start", zap.String("server_start_time", util.TimeFormat(time.Now())))

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
