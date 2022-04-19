package application

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/yufy/crontab/internal/pkg/transport/http"
)

type Application struct {
	name       string
	logger     *zap.Logger
	httpServer *http.Server
}

type Option func(*Application)

func HttpServerOption(svr *http.Server) Option {
	return func(a *Application) {
		a.httpServer = svr
	}
}

func New(name string, logger *zap.Logger, options ...Option) *Application {
	app := &Application{
		name:   name,
		logger: logger.With(zap.String("type", "Application")),
	}

	for _, o := range options {
		o(app)
	}

	return app
}

func (a *Application) Start() error {
	if a.httpServer != nil {
		if err := a.httpServer.Start(); err != nil {
			return errors.Wrap(err, "http server start error")
		}
	}
	return nil
}

func (a *Application) AwaitSinal() {
	ch := make(chan os.Signal, 1)
	signal.Reset(syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	s := <-ch
	a.logger.Info("receive a signal", zap.String("signal", s.String()))
	if a.httpServer != nil {
		if err := a.httpServer.Stop(); err != nil {
			a.logger.Warn("stop http server error", zap.Error(err))
		}
	}

	os.Exit(0)
}
