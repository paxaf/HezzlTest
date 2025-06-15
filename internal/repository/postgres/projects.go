package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/paxaf/HezzlTest/internal/entity"
)

const (
	queryGetProjects   = `SELECT id, name, created_at FROM projects`
	queryGetProject    = `SELECT id, name, created_at FROM projects WHERE id = $1`
	queryUpdateProject = `UPDATE projects SET name = $1 WHERE id = $2`
	queryLockProjects  = `LOCK TABLE projects IN ACCESS EXCLUSIVE MODE`
	queryAddProject    = `INSERT INTO projects(name) VALUES($1)`
	queryDeleteProject = `DELETE FROM projects WHERE id = $1`
)

func (r *PgPool) GetProjects(ctx context.Context) ([]entity.Project, error) {
	rows, err := r.db.Query(ctx, queryGetProjects)
	if err != nil {
		return nil, fmt.Errorf("error get all projects: %w", err)
	}
	defer rows.Close()
	var res []entity.Project
	for rows.Next() {
		var val entity.Project
		err = rows.Scan(
			&val.Id,
			&val.Name,
			&val.Created_at,
		)
		if err != nil {
			return nil, fmt.Errorf("error scan into projevts struct")
		}
		res = append(res, val)
	}
	return res, nil
}

func (r *PgPool) GetProject(ctx context.Context, id int) (*entity.Project, error) {
	var val entity.Project
	err := r.db.QueryRow(ctx, queryGetProject, id).Scan(
		&val.Id,
		&val.Name,
		&val.Created_at,
	)

	if err != nil {
		return nil, entity.ErrNotFound
	}

	return &val, nil
}

func (r *PgPool) UpdateProject(ctx context.Context, item *entity.Project) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return fmt.Errorf("failed begin tx: %w", err)
	}

	defer execTx(ctx, tx, &err)

	_, err = tx.Exec(ctx, queryLockProjects)
	if err != nil {
		return fmt.Errorf("failed while locking projects: %w", err)
	}
	_, err = tx.Exec(ctx, queryUpdateProject,
		item.Name,
		item.Id,
	)
	if err != nil {
		return fmt.Errorf("failed update project: %w", err)
	}
	return nil
}

func (r *PgPool) AddProject(ctx context.Context, item *entity.Project) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return fmt.Errorf("failed begin tx: %w", err)
	}

	defer execTx(ctx, tx, &err)

	_, err = tx.Exec(ctx, queryLockProjects)
	if err != nil {
		return fmt.Errorf("failed while locking projects: %w", err)
	}
	_, err = tx.Exec(ctx, queryAddProject, item.Name)
	if err != nil {
		return fmt.Errorf("failed create project: %w", err)
	}
	return nil
}

func (r *PgPool) DeleteProject(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return fmt.Errorf("failed begin tx: %w", err)
	}

	defer execTx(ctx, tx, &err)

	_, err = tx.Exec(ctx, queryLockProjects)
	if err != nil {
		return fmt.Errorf("failed while locking projects: %w", err)
	}
	_, err = tx.Exec(ctx, queryDeleteProject, id)
	if err != nil {
		return fmt.Errorf("failed delete project: %w", err)
	}
	return nil
}
