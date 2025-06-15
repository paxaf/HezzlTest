package usecase

import (
	"context"

	"github.com/paxaf/HezzlTest/internal/entity"
	"github.com/paxaf/HezzlTest/internal/logger"
)

func (uc *usecase) GetProject(ctx context.Context, key string, id int) (*entity.Project, error) {
	res, err := uc.repo.RedisGetProject(key)
	if err == nil {
		return res, nil
	}
	res, err = uc.repo.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}
	err = uc.repo.RedisSetItem(key, res)
	if err != nil {
		logger.Error("error set cache", err)
	}
	return res, nil
}

func (uc *usecase) GetProjects(ctx context.Context, key string) ([]entity.Project, error) {
	res, err := uc.repo.RedisGetProjects(key)
	if err == nil {
		return res, nil
	}
	res, err = uc.repo.GetProjects(ctx)
	if err != nil {
		return nil, err
	}
	err = uc.repo.RedisSetItem(key, res)
	if err != nil {
		logger.Error("error set cache", err)
	}
	return res, nil
}

func (uc *usecase) UpdateProject(ctx context.Context, item *entity.Project) error {
	err := uc.repo.CleanCache()
	if err != nil {
		return err
	}
	err = uc.repo.UpdateProject(ctx, item)
	if err != nil {
		return err
	}
	uc.repo.LogEvent(entity.NewProjectEvent(entity.Update, *item))
	return nil
}

func (uc *usecase) AddProject(ctx context.Context, item *entity.Project) error {
	err := uc.repo.CleanCache()
	if err != nil {
		return err
	}
	err = uc.repo.AddProject(ctx, item)
	if err != nil {
		return err
	}
	uc.repo.LogEvent(entity.NewProjectEvent(entity.Update, *item))
	return nil
}

func (uc *usecase) DeleteProject(ctx context.Context, id int) error {
	err := uc.repo.CleanCache()
	if err != nil {
		return err
	}
	err = uc.repo.DeleteProject(ctx, id)
	if err != nil {
		return err
	}
	uc.repo.LogEvent(entity.NewProjectEvent(entity.Update, entity.Project{Id: id}))
	return nil
}
