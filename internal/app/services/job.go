package services

import (
	"context"

	"github.com/yufy/crontab/internal/app/repositories"
	"github.com/yufy/crontab/internal/pkg/model"
	"go.uber.org/zap"
)

type JobService interface {
	List(ctx context.Context) ([]*model.Job, error)
	Save(ctx context.Context, job *model.Job) (*model.Job, error)
	Delete(ctx context.Context, name string) (*model.Job, error)
	Kill(ctx context.Context, name string) error
}

type defaultJobService struct {
	logger     *zap.Logger
	repository repositories.JobRepository
}

func NewJobService(logger *zap.Logger, repository repositories.JobRepository) JobService {
	return &defaultJobService{
		logger:     logger.With(zap.String("type", "services.Job")),
		repository: repository,
	}
}

func (s *defaultJobService) List(ctx context.Context) ([]*model.Job, error) {
	return s.repository.List(ctx)
}

func (s *defaultJobService) Save(ctx context.Context, job *model.Job) (*model.Job, error) {
	return s.repository.Save(ctx, job)
}

func (s *defaultJobService) Delete(ctx context.Context, name string) (*model.Job, error) {
	return s.repository.Delete(ctx, name)
}

func (s *defaultJobService) Kill(ctx context.Context, name string) error {
	return s.repository.Kill(ctx, name)
}
