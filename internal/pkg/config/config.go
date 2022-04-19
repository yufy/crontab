package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func New(path string) (*viper.Viper, error) {
	var (
		v   *viper.Viper
		err error
	)
	v = viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(path)

	if err = v.ReadInConfig(); err != nil {
		return nil, err
	}
	fmt.Printf("use config file -> %s\n", v.ConfigFileUsed())

	return v, nil
}
