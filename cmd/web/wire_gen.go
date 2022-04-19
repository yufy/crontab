// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
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

// Injectors from wire.go:

func CreateApp(cf string) (*application.Application, error) {
	viper, err := config.New(cf)
	if err != nil {
		return nil, err
	}
	option, err := app.NewOption(viper)
	if err != nil {
		return nil, err
	}
	logOption, err := log.NewOption(viper)
	if err != nil {
		return nil, err
	}
	logger, err := log.New(logOption)
	if err != nil {
		return nil, err
	}
	httpOption, err := http.NewOption(viper)
	if err != nil {
		return nil, err
	}
	etcdOption, err := etcd.NewOption(viper)
	if err != nil {
		return nil, err
	}
	client, err := etcd.New(etcdOption)
	if err != nil {
		return nil, err
	}
	jobRepository := repositories.NewJobRepository(logger, client)
	jobService := services.NewJobService(logger, jobRepository)
	jobController := controllers.NewJobController(logger, jobService)
	initRouter := app.RegisterRouter(jobController)
	engine := app.NewRouter(option, logger, initRouter)
	server, err := http.New(httpOption, logger, engine)
	if err != nil {
		return nil, err
	}
	applicationApplication := app.New(option, logger, server)
	return applicationApplication, nil
}
