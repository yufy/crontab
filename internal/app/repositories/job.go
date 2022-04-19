package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/yufy/crontab/internal/pkg/constant"
	"github.com/yufy/crontab/internal/pkg/model"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type JobRepository interface {
	List(ctx context.Context) ([]*model.Job, error)
	Save(ctx context.Context, job *model.Job) (*model.Job, error)
	Delete(ctx context.Context, name string) (*model.Job, error)
	Kill(ctx context.Context, name string) error
}

type defaultJobRepository struct {
	logger *zap.Logger
	client *clientv3.Client
}

func NewJobRepository(logger *zap.Logger, client *clientv3.Client) JobRepository {
	return &defaultJobRepository{
		logger: logger.With(zap.String("type", "repositories.Job")),
		client: client,
	}
}

func (r *defaultJobRepository) List(ctx context.Context) ([]*model.Job, error) {
	r.logger.Info("list crontab jobs")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	getResp, err := r.client.Get(ctx, constant.SaveJobDir, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	jobs := make([]*model.Job, 0, getResp.Count)
	for _, item := range getResp.Kvs {
		job := new(model.Job)
		if err := json.Unmarshal(item.Value, job); err != nil {
			r.logger.Warn("json decode job err", zap.Error(err))
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *defaultJobRepository) Save(ctx context.Context, job *model.Job) (*model.Job, error) {
	r.logger.Info("save crontab job")

	body, err := json.Marshal(job)
	if err != nil {
		r.logger.Warn("json encode job err", zap.Error(err))
		return nil, err
	}

	key := constant.SaveJobDir + job.Name
	putResp, err := r.client.Put(ctx, key, string(body), clientv3.WithPrevKV())
	if err != nil {
		r.logger.Warn("put job to etcd err", zap.Error(err))
		return nil, err
	}
	if putResp.PrevKv != nil {
		old := new(model.Job)
		if err := json.Unmarshal(putResp.PrevKv.Value, old); err != nil {
			// 解码老的job出错，不影响程序
			return nil, nil
		}
		return old, nil
	}
	return nil, nil
}

func (r *defaultJobRepository) Delete(ctx context.Context, name string) (*model.Job, error) {
	r.logger.Info("delete crontab job", zap.String("name", name))

	key := constant.SaveJobDir + name
	delResp, err := r.client.Delete(ctx, key, clientv3.WithPrevKV())
	if err != nil {
		r.logger.Warn("delete job err", zap.Error(err))
		return nil, err
	}

	if len(delResp.PrevKvs) > 0 {
		var oldJob model.Job
		if err := json.Unmarshal(delResp.PrevKvs[0].Value, &oldJob); err == nil {
			return &oldJob, nil
		}
	}

	return nil, nil
}

func (r *defaultJobRepository) Kill(ctx context.Context, name string) error {
	key := constant.KillJobDir + name

	leaseResp, err := r.client.Grant(ctx, 1)
	if err != nil {
		return err
	}

	if _, err := r.client.Put(ctx, key, "", clientv3.WithLease(leaseResp.ID)); err != nil {
		return err
	}

	return nil
}
