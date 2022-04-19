package log

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Option struct {
	FileName   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Level      string
	Stdout     bool
}

func NewOption(v *viper.Viper) (*Option, error) {
	var (
		err error
		o   = new(Option)
	)
	if err = v.UnmarshalKey("log", o); err != nil {
		return nil, err
	}
	return o, nil
}

func New(o *Option) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(o.Level)
	if err != nil {
		return nil, err
	}

	cores := make([]zapcore.Core, 0, 2)

	fw := zapcore.AddSync(&lumberjack.Logger{
		Filename:   o.FileName,
		MaxSize:    o.MaxSize,
		MaxBackups: o.MaxBackups,
		MaxAge:     o.MaxAge,
	})

	je := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	jc := zapcore.NewCore(je, fw, level)
	cores = append(cores, jc)
	if o.Stdout {
		cw := zapcore.Lock(os.Stdout)
		ce := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		cc := zapcore.NewCore(ce, cw, level)
		cores = append(cores, cc)
	}
	core := zapcore.NewTee(cores...)
	logger := zap.New(core)

	zap.ReplaceGlobals(logger)

	return logger, nil
}
