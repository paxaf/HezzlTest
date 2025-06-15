package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgPool struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PgPool {
	return &PgPool{db: pool}
}

func (r *PgPool) Close() {
	r.db.Close()
}
