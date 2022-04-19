package etcd

import (
	"context"
	"time"

	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Option struct {
	EndPoints   []string
	Username    string
	Password    string
	DialTimeout time.Duration
}

func NewOption(v *viper.Viper) (*Option, error) {
	var (
		err error
		o   = new(Option)
	)
	if err = v.UnmarshalKey("etcd", o); err != nil {
		return nil, err
	}
	o.DialTimeout *= time.Second

	return o, nil
}

func New(o *Option) (*clientv3.Client, error) {
	var (
		err    error
		client *clientv3.Client
	)

	if client, err = clientv3.New(clientv3.Config{
		Endpoints:   o.EndPoints,
		DialTimeout: o.DialTimeout,
		Username:    o.Username,
		Password:    o.Password,
	}); err != nil {
		return nil, err
	}
	// 启动的时候，去连接一下etcd，否则之后，etcd连接不上会导致程序阻塞
	ctx, cancel := context.WithTimeout(context.Background(), o.DialTimeout)
	defer cancel()
	_, err = client.Status(ctx, o.EndPoints[0])
	if err != nil {
		return nil, err
	}
	return client, nil
}
