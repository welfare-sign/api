package task

import (
	"context"

	"github.com/spf13/viper"

	"welfare-sign/internal/pkg/config"
	"welfare-sign/internal/pkg/log"
	"welfare-sign/internal/pkg/task"
	"welfare-sign/internal/service"
)

// Run 定时任务执行
func Run(svc *service.Service) {
	t := task.Default()
	t.AddFunc(viper.GetString(config.KeyTaskCheckinExpiredTimeStartInterval), "启动清除失效的任务", svc.FailureIssueRecord)
	log.Info(context.Background(), "task running")
	t.Run()
}
