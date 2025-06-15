package repository

import (
	"context"

	"github.com/paxaf/HezzlTest/internal/entity"
)

type Postgres interface {
	GetItemsByName(ctx context.Context, name string) (*entity.Goods, error)
	GetItemsByProject(ctx context.Context, projectId int) (*[]entity.Goods, error)
	GetItem(ctx context.Context, goodsId int) (*entity.Goods, error)
	GetAllItems(ctx context.Context) (*[]entity.Goods, error)
	CreateItem(ctx context.Context, item *entity.Goods) error
	UpdateItem(ctx context.Context, item *entity.Goods) error
	DeleteItem(ctx context.Context, id int) error
}
