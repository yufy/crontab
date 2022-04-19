package store

import (
	"database/sql"

	"github.com/yufy/crontab/internal/pkg/model"
	"go.uber.org/zap"
)

type Store struct {
	db     *sql.DB
	logger *zap.Logger
	ch     chan *model.Log
}

func New(db *sql.DB, logger *zap.Logger) *Store {
	s := &Store{
		db:     db,
		logger: logger.With(zap.String("type", "worker.store")),
		ch:     make(chan *model.Log, 100),
	}

	go s.handle()

	return s
}

func (s *Store) handle() {
	sql := "INSERT INTO logs(job_name, command, output, error, plan_time, schedule_time, start_time, end_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	for log := range s.ch {
		if _, err := s.db.Exec(sql,
			log.JobName, log.Command, log.Output, log.Err, log.PlanTime.UnixMilli(), log.ScheduleTime.UnixMilli(), log.StartTime.UnixMilli(), log.EndTime.UnixMilli()); err != nil {
			s.logger.Warn("insert logs error", zap.String("error", err.Error()))
		}
	}
}

func (s *Store) Push(l *model.Log) {
	s.ch <- l
}
