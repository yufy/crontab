//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/yufy/crontab/internal/pkg/config"
	"github.com/yufy/crontab/internal/pkg/database"
	"github.com/yufy/crontab/internal/pkg/etcd"
	"github.com/yufy/crontab/internal/pkg/log"
	"github.com/yufy/crontab/internal/worker"
	"github.com/yufy/crontab/internal/worker/executor"
	"github.com/yufy/crontab/internal/worker/executor/locker"
	"github.com/yufy/crontab/internal/worker/job"
	"github.com/yufy/crontab/internal/worker/scheduler"
	"github.com/yufy/crontab/internal/worker/store"
)

func CreateWorker(cf string) (*worker.Worker, error) {
	panic(wire.Build(
		config.New,
		log.NewOption, log.New,
		database.NewOption, database.New,
		store.New,
		etcd.NewOption, etcd.New,
		locker.New,
		executor.New,
		scheduler.New,
		job.NewManager,
		worker.NewOption, worker.New,
	))
}
