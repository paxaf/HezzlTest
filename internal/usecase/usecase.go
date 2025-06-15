package usecase

import (
	"context"

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
	DeleteProject(ctx context.Context, id int) error
	AddProject(ctx context.Context, item *entity.Project) error
	UpdateProject(ctx context.Context, item *entity.Project) error
	GetProjects(ctx context.Context, key string) ([]entity.Project, error)
	GetProject(ctx context.Context, key string, id int) (*entity.Project, error)
}

func New(repo repository.Repository) *usecase {
	return &usecase{repo: repo}
}
