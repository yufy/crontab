package executor

import (
	"context"
	"math/rand"
	"os/exec"
	"time"

	"github.com/yufy/crontab/internal/pkg/model"
	"github.com/yufy/crontab/internal/worker/executor/locker"
	"go.uber.org/zap"
)

type JobExecuteInfo struct {
	Job        *model.Job
	PlanTime   time.Time
	RealTime   time.Time
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

type JobExecuteResult struct {
	JobExecuteInfo *JobExecuteInfo
	Output         []byte
	Err            error
	StartTime      time.Time
	EndTime        time.Time
}

type Executor struct {
	logger *zap.Logger
	locker *locker.Locker
}

func New(logger *zap.Logger, locker *locker.Locker) *Executor {
	return &Executor{
		logger: logger.With(zap.String("type", "worker.executor")),
		locker: locker,
	}
}

func (e *Executor) Execute(info *JobExecuteInfo, ch chan<- *JobExecuteResult) {
	go e.execute(info, ch)
}

func (e *Executor) execute(info *JobExecuteInfo, ch chan<- *JobExecuteResult) {
	result := new(JobExecuteResult)
	result.JobExecuteInfo = info

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	result.StartTime = time.Now()
	err := e.locker.Try(info.Job.Name)
	if err != nil {
		result.Err = err
		result.EndTime = time.Now()
	} else {
		defer e.locker.Release(info.Job.Name)

		result.StartTime = time.Now()

		cmd := exec.CommandContext(info.Ctx, "C:\\Program Files\\Git\\bin\\bash", "-c", info.Job.Command)
		result.Output, result.Err = cmd.CombinedOutput()
		result.EndTime = time.Now()
	}

	ch <- result
}
