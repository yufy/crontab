package app

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yufy/crontab/internal/app/middlewares"
	"github.com/yufy/crontab/internal/pkg/application"
	"github.com/yufy/crontab/internal/pkg/transport/http"
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
)

type Option struct {
	Name string
	Mode string
}

type InitRouter func(r *gin.Engine)

func NewOption(v *viper.Viper) (*Option, error) {
	var (
		err error
		o   = new(Option)
	)

	if err = v.UnmarshalKey("app", o); err != nil {
		return nil, err
	}
	return o, nil
}

func NewRouter(o *Option, logger *zap.Logger, init InitRouter) *gin.Engine {
	gin.SetMode(o.Mode)
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, false))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(middlewares.Translations())

	init(r)

	return r
}

func New(o *Option, logger *zap.Logger, httpServer *http.Server) *application.Application {
	app := application.New(o.Name, logger, application.HttpServerOption(httpServer))

	return app
}
