package task

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"welfare-sign/internal/pkg/log"
	"welfare-sign/internal/pkg/util"
	"welfare-sign/internal/pkg/wsgin"
)

// Task .
type Task interface {
	AddFunc(spec string, taskName string, cmd func(ctx context.Context) (wsgin.APICode, error))
	Run()
}

// Default .
func Default() Task {
	task := DefaultTask{}
	task.c = cron.New()
	return task
}

// DefaultTask .
type DefaultTask struct {
	c *cron.Cron
}

var ctx = context.WithValue(context.Background(), "req_ip", util.GetIP())

// AddFunc 添加定时任务
func (t DefaultTask) AddFunc(spec string, taskName string, cmd func(ctx context.Context) (wsgin.APICode, error)) {
	fmt.Println("task:", spec, "\t", taskName)
	f := func() {
		uuid, _ := util.NewV4()
		taskId := uuid.String()
		startTime := time.Now()
		ctx = context.WithValue(ctx, "req_id", taskId)

		log.Info(ctx, taskName, zap.Time("start_time", startTime))
		endTime := time.Now()
		latency := time.Since(startTime)
		code, err := cmd(ctx)
		if err != nil {
			log.Error(ctx, taskName, zap.Time("end_time", endTime), zap.Error(err), zap.Duration("latency", latency), zap.String("code", string(code)))
			return
		}
		log.Info(ctx, taskName, zap.Time("end_time", endTime), zap.Duration("latency", latency))
	}
	// 启动时触发一次
	f()
	_, err := t.c.AddFunc(spec, f)
	if err != nil {
		panic(err)
	}
}

// Run 启动
func (t DefaultTask) Run() {
	t.c.Run()
}
