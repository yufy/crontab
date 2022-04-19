//go:build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/yufy/crontab/internal/app"
	"github.com/yufy/crontab/internal/app/controllers"
	"github.com/yufy/crontab/internal/app/repositories"
	"github.com/yufy/crontab/internal/app/services"
	"github.com/yufy/crontab/internal/pkg/application"
	"github.com/yufy/crontab/internal/pkg/config"
	"github.com/yufy/crontab/internal/pkg/etcd"
	"github.com/yufy/crontab/internal/pkg/log"
	"github.com/yufy/crontab/internal/pkg/transport/http"
)

func CreateApp(cf string) (*application.Application, error) {
	panic(wire.Build(
		config.New,
		log.NewOption, log.New,
		etcd.NewOption, etcd.New,
		repositories.NewJobRepository,
		services.NewJobService,
		controllers.NewJobController,
		app.NewOption, app.RegisterRouter, app.NewRouter,
		http.NewOption, http.New,
		app.New,
	))
}
