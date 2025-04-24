package db

import (
	"context"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"yadro.com/course/update/core"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

type Stats struct {
	WordsTotal    int `db:"words_total"`
	WordsUnique   int `db:"words_unique"`
	ComicsFetched int `db:"comics_fetched"`
}

func New(log *slog.Logger, address string) (*DB, error) {

	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}

	return &DB{
		log:  log,
		conn: db,
	}, nil
}

func (db *DB) Add(ctx context.Context, comics core.Comics) error {
	_, err := db.conn.ExecContext(ctx,
		`INSERT INTO comics (comic_id, img_url, keywords) VALUES($1, $2, $3)`,
		comics.ID, comics.URL, comics.Words)

	return err
}

func (db *DB) Stats(ctx context.Context) (core.DBStats, error) {
	query := `
	SELECT 
		COUNT(*) AS words_total,
		COUNT(DISTINCT keyword) AS words_unique,
		(SELECT COUNT(*) FROM comics) AS comics_fetched
	FROM (
		SELECT unnest(keywords) AS keyword
		FROM comics
	) AS subquery;
	`
	var stats Stats
	err := db.conn.GetContext(ctx, &stats, query)

	return core.DBStats{
		WordsTotal:    stats.WordsTotal,
		WordsUnique:   stats.WordsUnique,
		ComicsFetched: stats.ComicsFetched,
	}, err
}

func (db *DB) IDs(ctx context.Context) ([]int, error) {
	var ids []int

	err := db.conn.SelectContext(ctx, &ids, `SELECT comic_id FROM comics`)
	if err != nil {
		db.log.Error("failed to fetch comic IDs", "error", err)
		return nil, fmt.Errorf("fetch comic IDs: %w", err)
	}

	return ids, nil
}

func (db *DB) Drop(ctx context.Context) error {
	_, err := db.conn.ExecContext(ctx, `DELETE FROM comics`)
	return err
}
