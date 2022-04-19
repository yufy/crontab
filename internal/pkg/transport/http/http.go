package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Option struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server struct {
	app        string
	o          *Option
	logger     *zap.Logger
	httpServer http.Server
}

func NewOption(v *viper.Viper) (*Option, error) {
	var (
		err error
		o   = new(Option)
	)
	if err = v.UnmarshalKey("http", o); err != nil {
		return nil, err
	}
	o.ReadTimeout *= time.Second
	o.WriteTimeout *= time.Second

	return o, nil
}

func New(o *Option, logger *zap.Logger, r *gin.Engine) (*Server, error) {
	httpServer := http.Server{
		Addr:         o.Host + ":" + o.Port,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,
		Handler:      r,
	}

	return &Server{
		o:          o,
		logger:     logger.With(zap.String("type", "http.Server")),
		httpServer: httpServer,
	}, nil
}

func (s *Server) Application(name string) {
	s.app = name
}

func (s *Server) Start() error {
	s.logger.Info("http server starting...", zap.String("addr", s.httpServer.Addr))
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("start http server err", zap.Error(err))
			return
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("stopping http server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server error")
	}
	return nil
}
