package usecase

import (
	"context"
	"log"

	"github.com/paxaf/HezzlTest/internal/entity"
	"github.com/paxaf/HezzlTest/internal/repository"
)

type usecase struct {
	repo repository.Repository
}

type Usecase interface {
	GetAllItems(ctx context.Context, key string) ([]entity.Goods, error)
	GetItem(ctx context.Context, key string, goodsId int) (*entity.Goods, error)
	GetItemsByProject(ctx context.Context, key string, projectId int) ([]entity.Goods, error)
	GetItemsByName(ctx context.Context, key string, name string) ([]entity.Goods, error)
	CreateItem(ctx context.Context, item *entity.Goods) error
	UpdateItem(ctx context.Context, item *entity.Goods) error
	DeleteItem(ctx context.Context, id int) error
}

func New(repo repository.Repository) *usecase {
	return &usecase{repo: repo}
}

func (uc *usecase) GetAllItems(ctx context.Context, key string) ([]entity.Goods, error) {
	res, err := uc.repo.RedisGetItems(key)
	if err == nil {
		return res, nil
	}
	res, err = uc.repo.GetAllItems(ctx)
	if err != nil {
		return nil, err
	}
	err = uc.repo.RedisSetItem(key, res)
	if err != nil {
		log.Printf("error set cache: %v", err)
	}
	return res, nil
}

func (uc *usecase) GetItem(ctx context.Context, key string, goodsId int) (*entity.Goods, error) {
	res, err := uc.repo.RedisGetItem(key)
	if err == nil {
		return res, nil
	}
	res, err = uc.repo.GetItem(ctx, goodsId)
	if err != nil {
		return nil, err
	}
	err = uc.repo.RedisSetItem(key, res)
	if err != nil {
		log.Printf("error set cache: %v", err)
	}
	return res, nil
}

func (uc *usecase) GetItemsByProject(ctx context.Context, key string, projectId int) ([]entity.Goods, error) {
	res, err := uc.repo.RedisGetItems(key)
	if err == nil {
		return res, nil
	}
	res, err = uc.repo.GetItemsByProject(ctx, projectId)
	if err != nil {
		return nil, err
	}
	err = uc.repo.RedisSetItem(key, res)
	if err != nil {
		log.Printf("error set cache: %v", err)
	}
	return res, nil
}

func (uc *usecase) GetItemsByName(ctx context.Context, key string, name string) ([]entity.Goods, error) {
	res, err := uc.repo.RedisGetItems(key)
	if err == nil {
		return res, nil
	}
	res, err = uc.repo.GetItemsByName(ctx, name)
	if err != nil {
		return nil, err
	}
	err = uc.repo.RedisSetItem(key, res)
	if err != nil {
		log.Printf("error set cache: %v", err)
	}
	return res, nil
}

func (uc *usecase) CreateItem(ctx context.Context, item *entity.Goods) error {
	err := uc.repo.CleanCache()
	if err != nil {
		return err
	}
	err = uc.repo.CreateItem(ctx, item)
	if err != nil {
		return err
	}
	return nil
}

func (uc *usecase) UpdateItem(ctx context.Context, item *entity.Goods) error {
	err := uc.repo.CleanCache()
	if err != nil {
		return err
	}
	err = uc.repo.UpdateItem(ctx, item)
	if err != nil {
		return err
	}
	return nil
}

func (uc *usecase) DeleteItem(ctx context.Context, id int) error {
	err := uc.repo.CleanCache()
	if err != nil {
		return err
	}
	err = uc.repo.DeleteItem(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
