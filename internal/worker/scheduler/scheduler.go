package scheduler

import (
	"context"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/yufy/crontab/internal/pkg/model"
	"github.com/yufy/crontab/internal/worker/executor"
	"github.com/yufy/crontab/internal/worker/executor/locker"
	"github.com/yufy/crontab/internal/worker/store"

	"go.uber.org/zap"
)

type EventType int

const (
	SaveEvent EventType = iota
	DeleteEvent
	KillEvent
)

type Event struct {
	Typ EventType
	Job *model.Job
}

type Scheduler struct {
	logger     *zap.Logger
	eventCh    chan *Event
	resultCh   chan *executor.JobExecuteResult
	tables     map[string]*jobInfo
	exec       *executor.Executor
	execTables map[string]*executor.JobExecuteInfo
	store      *store.Store
}

type jobInfo struct {
	job      *model.Job
	expr     *cronexpr.Expression
	nextTime time.Time
}

func New(logger *zap.Logger, store *store.Store, exec *executor.Executor) *Scheduler {
	s := &Scheduler{
		logger:     logger.With(zap.String("type", "worker.scheduler")),
		eventCh:    make(chan *Event, 100),
		resultCh:   make(chan *executor.JobExecuteResult, 100),
		tables:     make(map[string]*jobInfo),
		exec:       exec,
		execTables: make(map[string]*executor.JobExecuteInfo),
		store:      store,
	}

	go s.scheduler()

	return s
}

func (s *Scheduler) scheduler() {
	var (
		duration time.Duration
		ticker   *time.Ticker
	)
	for {
		duration = s.tryRun()
		ticker = time.NewTicker(duration)

		select {
		case event := <-s.eventCh:
			{
				s.handleEvent(event)
			}
		case result := <-s.resultCh:
			{
				s.handleResult(result)
			}
		case <-ticker.C:
		}

		ticker.Reset(s.tryRun())
	}
}

func (s *Scheduler) tryRun() time.Duration {
	var nearTime *time.Time

	if len(s.tables) == 0 {
		return 1 * time.Second
	}

	now := time.Now()
	for _, jobInfo := range s.tables {
		if jobInfo.nextTime.Before(now) || jobInfo.nextTime.Equal(now) {
			s.run(jobInfo)
			jobInfo.nextTime = jobInfo.expr.Next(now)
		}

		if nearTime == nil || jobInfo.nextTime.Before(*nearTime) {
			nearTime = &jobInfo.nextTime
		}
	}

	return (*nearTime).Sub(now)
}

func (s *Scheduler) run(info *jobInfo) {
	if _, ok := s.execTables[info.job.Name]; ok {
		s.logger.Info(info.job.Name + " is still running, skip it.")
		return
	}

	executeInfo := &executor.JobExecuteInfo{
		Job:      info.job,
		PlanTime: info.nextTime,
		RealTime: time.Now(),
	}
	executeInfo.Ctx, executeInfo.CancelFunc = context.WithCancel(context.TODO())
	s.execTables[info.job.Name] = executeInfo

	s.exec.Execute(executeInfo, s.resultCh)
}

func (s *Scheduler) handleEvent(event *Event) {
	switch event.Typ {
	case SaveEvent:
		s.handleSave(event.Job)
	case DeleteEvent:
		s.handleDelete(event.Job)
	case KillEvent:
		s.handleKill(event.Job)
	}
}

func (s *Scheduler) Push(event *Event) {
	s.eventCh <- event
}

func (s *Scheduler) handleSave(job *model.Job) {
	expr, err := cronexpr.Parse(job.Expr)
	if err != nil {
		return
	}
	info := &jobInfo{
		job:      job,
		expr:     expr,
		nextTime: expr.Next(time.Now()),
	}
	s.tables[job.Name] = info
}

func (s *Scheduler) handleDelete(job *model.Job) {
	delete(s.tables, job.Name)
}

func (s *Scheduler) handleKill(job *model.Job) {
	if jobExecuteInfo, ok := s.execTables[job.Name]; ok {
		jobExecuteInfo.CancelFunc()
		delete(s.execTables, job.Name)
	}
}

func (s *Scheduler) handleResult(result *executor.JobExecuteResult) {
	delete(s.execTables, result.JobExecuteInfo.Job.Name)

	if result.Err == locker.ErrLocked {
		return
	}

	log := &model.Log{
		JobName:      result.JobExecuteInfo.Job.Name,
		Command:      result.JobExecuteInfo.Job.Command,
		Output:       string(result.Output),
		PlanTime:     result.JobExecuteInfo.PlanTime,
		ScheduleTime: result.JobExecuteInfo.RealTime,
		StartTime:    result.StartTime,
		EndTime:      result.EndTime,
	}
	if result.Err != nil {
		log.Err = result.Err.Error()
	}

	s.store.Push(log)
}
