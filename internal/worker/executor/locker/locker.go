package locker

import (
	"context"
	"errors"
	"sync"

	"github.com/yufy/crontab/internal/pkg/constant"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var (
	ErrLocked = errors.New("lock failed")
)

type Locker struct {
	logger     *zap.Logger
	kv         clientv3.KV
	lease      clientv3.Lease
	mu         sync.Mutex
	jobLockers map[string]*jobLockerInfo
}

type jobLockerInfo struct {
	leaseID clientv3.LeaseID
	cancel  context.CancelFunc
}

func New(logger *zap.Logger, client *clientv3.Client) *Locker {
	return &Locker{
		logger:     logger.With(zap.String("type", "worker.executor.locker")),
		kv:         clientv3.NewKV(client),
		lease:      clientv3.NewLease(client),
		jobLockers: make(map[string]*jobLockerInfo),
	}
}

func (l *Locker) Try(name string) error {
	// 1、创建租约
	leaseResp, err := l.lease.Grant(context.TODO(), 5)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.TODO())

	// 2、自动续租
	keepResp, err := l.lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		l.lease.Revoke(context.TODO(), leaseResp.ID)
		return err
	}
	go func() {
		for result := range keepResp {
			if result == nil {
				return
			}
		}
	}()

	// 3、创建事务
	tx := l.kv.Txn(context.TODO())
	key := constant.LockerJobDir + name
	// 4、事务抢锁
	tx.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "", clientv3.WithLease(leaseResp.ID))).
		Else(clientv3.OpGet(key))
	txResp, err := tx.Commit()
	if err != nil {
		l.lease.Revoke(context.TODO(), leaseResp.ID)
		cancel()
		return err
	}
	if !txResp.Succeeded {
		return ErrLocked
	}

	l.mu.Lock()
	l.jobLockers[name] = &jobLockerInfo{
		leaseID: leaseResp.ID,
		cancel:  cancel,
	}
	l.mu.Unlock()

	return nil
}

func (l *Locker) Release(name string) {
	if info, ok := l.jobLockers[name]; ok {
		l.lease.Revoke(context.TODO(), info.leaseID)
		info.cancel()

		l.mu.Lock()
		delete(l.jobLockers, name)
		l.mu.Unlock()
	}
}
