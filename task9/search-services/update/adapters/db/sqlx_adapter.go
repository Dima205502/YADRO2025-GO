package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type sqlxAdapter struct {
	db *sqlx.DB
}

func NewSQLxAdapter(db *sqlx.DB) *sqlxAdapter {
	return &sqlxAdapter{db: db}
}

func (s *sqlxAdapter) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *sqlxAdapter) GetContext(ctx context.Context, dest interface{}, query string, args ...any) error {
	return s.db.GetContext(ctx, dest, query, args...)
}

func (s *sqlxAdapter) SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error {
	return s.db.SelectContext(ctx, dest, query, args...)
}

func (s *sqlxAdapter) GetDB() *sql.DB {
	return s.db.DB
}
