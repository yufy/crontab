package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Option struct {
	Driver       string
	Host         string
	Port         string
	DBName       string
	Username     string
	Password     string
	Charset      string
	Collation    string
	ReadTimeout  string
	WriteTimeout string
}

func NewOption(v *viper.Viper) (*Option, error) {
	var (
		err error
		o   = new(Option)
	)

	if v.UnmarshalKey("database", o); err != nil {
		return nil, err
	}
	return o, nil
}

func New(o *Option, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open(o.Driver, generateDSN(o))
	if err != nil {
		logger.Error("load database failed", zap.String("error", err.Error()))
		return nil, err
	}
	logger.Info("load database success")
	return db, nil
}

func generateDSN(o *Option) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&collation=%s&readTimeout=%s&writeTimeout=%s&parseTime=true&loc=Local",
		o.Username, o.Password, o.Host, o.Port, o.DBName, o.Charset, o.Collation, o.ReadTimeout, o.WriteTimeout)
}
