package repository

import (
	"context"

	"github.com/paxaf/HezzlTest/internal/entity"
)

type Postgres interface {
	GetItemsByName(ctx context.Context, name string) ([]entity.Goods, error)
	GetItemsByProject(ctx context.Context, projectId int) ([]entity.Goods, error)
	GetItem(ctx context.Context, goodsId int) (*entity.Goods, error)
	GetAllItems(ctx context.Context) ([]entity.Goods, error)
	CreateItem(ctx context.Context, item *entity.Goods) error
	UpdateItem(ctx context.Context, item *entity.Goods) error
	DeleteItem(ctx context.Context, id int) error
	DeleteProject(ctx context.Context, id int) error
	AddProject(ctx context.Context, item *entity.Project) error
	UpdateProject(ctx context.Context, item *entity.Project) error
	GetProject(ctx context.Context, id int) (*entity.Project, error)
	GetProjects(ctx context.Context) ([]entity.Project, error)
}

type Redis interface {
	RedisGetItems(key string) ([]entity.Goods, error)
	RedisGetItem(key string) (*entity.Goods, error)
	RedisSetItem(key string, item interface{}) error
	RedisGetProjects(key string) ([]entity.Project, error)
	RedisGetProject(key string) (*entity.Project, error)
	CleanCache() error
}

type Repository interface {
	Postgres
	Redis
}

type Repo struct {
	Redis
	Postgres
}

func New(redis Redis, pgpool Postgres) *Repo {
	return &Repo{
		Redis:    redis,
		Postgres: pgpool,
	}
}
