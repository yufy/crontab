package worker

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"github.com/yufy/crontab/internal/worker/job"
	"go.uber.org/zap"
)

type Option struct {
	Name string
}

type Worker struct {
	name       string
	logger     *zap.Logger
	manager    *job.Manager
	ctx        context.Context
	cancenFunc context.CancelFunc
}

func NewOption(v *viper.Viper) (*Option, error) {
	var (
		err error
		o   = new(Option)
	)

	if err = v.UnmarshalKey("worker", o); err != nil {
		return nil, err
	}

	return o, nil
}

func New(o *Option, logger *zap.Logger, manager *job.Manager) *Worker {
	ctx, cancel := context.WithCancel(context.TODO())
	return &Worker{
		name:       o.Name,
		logger:     logger.With(zap.String("type", "worker")),
		manager:    manager,
		ctx:        ctx,
		cancenFunc: cancel,
	}
}

func (w *Worker) Run() error {
	w.logger.Info("worker running...")
	go w.manager.Watch(w.ctx)
	go w.manager.WatchKill(w.ctx)
	return nil
}

func (w *Worker) AwaitSinal() {
	ch := make(chan os.Signal, 1)
	signal.Reset(syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	sig := <-ch

	w.logger.Info("receive a signal", zap.String("signal", sig.String()))
	w.logger.Info("stopping worker...")

	w.cancenFunc()

	os.Exit(0)
}
