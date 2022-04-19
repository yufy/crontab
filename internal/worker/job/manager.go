package job

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/yufy/crontab/internal/pkg/constant"
	"github.com/yufy/crontab/internal/pkg/model"
	"github.com/yufy/crontab/internal/worker/scheduler"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Manager struct {
	logger *zap.Logger
	client *clientv3.Client
	s      *scheduler.Scheduler
}

func NewManager(logger *zap.Logger, client *clientv3.Client, s *scheduler.Scheduler) *Manager {
	return &Manager{
		logger: logger.With(zap.String("type", "worker.job.manager")),
		client: client,
		s:      s,
	}
}

func (m *Manager) Watch(ctx context.Context) error {
	getResp, err := m.client.Get(ctx, constant.SaveJobDir, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	// 将已有的任务添加到调度中
	for _, kv := range getResp.Kvs {
		var se scheduler.Event
		if err := json.Unmarshal(kv.Value, &se.Job); err != nil {
			continue
		}
		se.Typ = scheduler.SaveEvent
		m.s.Push(&se)
	}

	// 监听任务变化，同步到调度中
	ch := m.client.Watch(ctx, constant.SaveJobDir, clientv3.WithPrefix(), clientv3.WithRev(getResp.Header.Revision+1))
	for resp := range ch {
		for _, event := range resp.Events {
			var schedulerEvent scheduler.Event
			schedulerEvent.Job = &model.Job{}

			switch event.Type {
			case mvccpb.PUT:
				if err := json.Unmarshal(event.Kv.Value, schedulerEvent.Job); err != nil {
					m.logger.Error(err.Error())
					continue
				}
				schedulerEvent.Typ = scheduler.SaveEvent
			case mvccpb.DELETE:
				schedulerEvent.Job.Name = strings.TrimPrefix(string(event.Kv.Key), constant.SaveJobDir)
				schedulerEvent.Typ = scheduler.DeleteEvent
			}
			m.s.Push(&schedulerEvent)
		}
	}
	return nil
}

func (m *Manager) WatchKill(ctx context.Context) error {
	ch := m.client.Watch(ctx, constant.KillJobDir, clientv3.WithPrefix())
	for resp := range ch {
		for _, event := range resp.Events {
			var schedulerEvent scheduler.Event
			schedulerEvent.Job = &model.Job{}

			switch event.Type {
			case mvccpb.PUT:
				schedulerEvent.Job.Name = strings.TrimPrefix(string(event.Kv.Key), constant.KillJobDir)
				schedulerEvent.Typ = scheduler.KillEvent
				m.s.Push(&schedulerEvent)
			}
		}
	}
	return nil
}
