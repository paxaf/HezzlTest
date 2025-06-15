package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/paxaf/HezzlTest/internal/entity"
)

const (
	queryGetItemsByProject = `SELECT id, project_id, name, description, priority, removed, created_at
	FROM GOODS
	WHERE project_id = $1`
	queryGetItem = `SELECT id, project_id, name, description, priority, removed, created_at
	FROM GOODS
	WHERE id = $1`
	queryGetItemsByName = `SELECT id, project_id, name, description, priority, removed, created_at
	FROM GOODS
	WHERE name ILIKE '%' || $1 || '%'`
	queryGetAllItems = `SELECT id, project_id, name, description, priority, removed, created_at
	FROM GOODS`
	queryLockGoods  = `LOCK TABLE goods IN ACCESS EXCLUSIVE MODE`
	queryCreateItem = `INSERT INTO GOODS (project_id, name, description, priority) 
	VALUES ($1, $2, $3, (SELECT COALESCE(MAX(priority), 0) + 1 FROM GOODS))`
	queryUpdateItem = `UPDATE GOODS SET name = $1, description = $2, priority = $3, removed = $4`
	queryDeleteItem = `DELETE FROM GOODS WHERE id = $1`
)

func execTx(ctx context.Context, tx pgx.Tx, errp *error) {
	if *errp != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			log.Printf("failed rollback tx: %v", *errp)
		}
	} else {
		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			log.Printf("failed commit tx: %v", *errp)
		}
	}
}

func (r *PgPool) GetItemsByProject(ctx context.Context, projectId int) ([]entity.Goods, error) {
	rows, err := r.db.Query(ctx, queryGetItemsByProject, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed get goods by project: %w", err)
	}
	defer rows.Close()
	var res []entity.Goods
	for rows.Next() {
		var item entity.Goods
		err = rows.Scan(
			&item.Id,
			&item.ProjectId,
			&item.Name,
			&item.Description,
			&item.Priority,
			&item.Removed,
			&item.Created_at,
		)
		if err != nil {
			return nil, fmt.Errorf("failed parse into sturct: %w", err)
		}
		res = append(res, item)
	}
	return res, nil
}

func (r *PgPool) GetItem(ctx context.Context, goodsId int) (*entity.Goods, error) {
	row := r.db.QueryRow(ctx, queryGetItem, goodsId)
	var item entity.Goods
	err := row.Scan(
		&item.Id,
		&item.ProjectId,
		&item.Name,
		&item.Description,
		&item.Priority,
		&item.Removed,
		&item.Created_at,
	)
	if err != nil {
		return nil, fmt.Errorf("failed parse into sturct: %w", err)
	}
	return &item, nil
}

func (r *PgPool) GetItemsByName(ctx context.Context, name string) ([]entity.Goods, error) {
	rows, err := r.db.Query(ctx, queryGetItemsByName, name)
	if err != nil {
		return nil, fmt.Errorf("failed get goods by name: %w", err)
	}
	defer rows.Close()
	var res []entity.Goods
	for rows.Next() {
		var item entity.Goods
		err = rows.Scan(
			&item.Id,
			&item.ProjectId,
			&item.Name,
			&item.Description,
			&item.Priority,
			&item.Removed,
			&item.Created_at,
		)
		if err != nil {
			return nil, fmt.Errorf("failed parse into sturct: %w", err)
		}
		res = append(res, item)
	}
	return res, nil
}

func (r *PgPool) GetAllItems(ctx context.Context) ([]entity.Goods, error) {
	rows, err := r.db.Query(ctx, queryGetAllItems)
	if err != nil {
		return nil, fmt.Errorf("failed get goods by project: %w", err)
	}
	defer rows.Close()
	var res []entity.Goods
	for rows.Next() {
		var item entity.Goods
		err = rows.Scan(
			&item.Id,
			&item.ProjectId,
			&item.Name,
			&item.Description,
			&item.Priority,
			&item.Removed,
			&item.Created_at,
		)
		if err != nil {
			return nil, fmt.Errorf("failed parse into sturct: %w", err)
		}
		res = append(res, item)
	}
	return res, nil
}

func (r *PgPool) CreateItem(ctx context.Context, item *entity.Goods) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return fmt.Errorf("failed begin tx: %w", err)
	}

	defer execTx(ctx, tx, &err)

	_, err = tx.Exec(ctx, queryLockGoods)
	if err != nil {
		return fmt.Errorf("failed while locking goods: %w", err)
	}
	_, err = tx.Exec(ctx, queryCreateItem,
		item.ProjectId,
		item.Name,
		item.Description,
	)
	if err != nil {
		return fmt.Errorf("failed create item: %w", err)
	}
	return nil
}

func (r *PgPool) UpdateItem(ctx context.Context, item *entity.Goods) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return fmt.Errorf("failed begin tx: %w", err)
	}

	defer execTx(ctx, tx, &err)

	_, err = tx.Exec(ctx, queryLockGoods)
	if err != nil {
		return fmt.Errorf("failed while locking goods: %w", err)
	}

	res, err := tx.Exec(ctx, queryUpdateItem,
		item.Name,
		item.Description,
		item.Priority,
		item.Removed,
	)
	if err != nil {
		return fmt.Errorf("failed update item: %w", err)
	}
	if res.RowsAffected() == 0 {
		err = errors.New("no rows affected")
		return fmt.Errorf("failed update item: %w", err)
	}
	return nil
}

func (r *PgPool) DeleteItem(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return fmt.Errorf("failed begin tx: %w", err)
	}

	defer execTx(ctx, tx, &err)

	_, err = tx.Exec(ctx, queryLockGoods)
	if err != nil {
		return fmt.Errorf("failed while locking goods: %w", err)
	}
	res, err := tx.Exec(ctx, queryDeleteItem, id)
	if err != nil {
		return fmt.Errorf("failed delete item: %w", err)
	}
	if res.RowsAffected() == 0 {
		err = errors.New("no rows affected")
		return fmt.Errorf("failed delete item: %w", err)
	}
	return nil
}
