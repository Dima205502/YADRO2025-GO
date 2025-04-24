package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgtype"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"yadro.com/course/search/core"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
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

func (db *DB) SearchByWord(ctx context.Context, keyword string) ([]int, error) {
	query := `
	SELECT comics_id 
	FROM comics
	WHERE $1 = ANY(keywords)
	`

	var IDs []int
	err := db.conn.SelectContext(
		ctx, &IDs,
		query,
		keyword,
	)

	return IDs, err
}

func (db *DB) FetchComics(ctx context.Context, id int) (core.Comics, []string, error) {
	query := `
        SELECT comics_id, img_url, keywords 
        FROM comics 
        WHERE comics_id = $1
    `

	var res struct {
		ID       int              `db:"comics_id"`
		Keywords pgtype.TextArray `db:"keywords"`
		URL      string           `db:"img_url"`
	}

	err := db.conn.GetContext(ctx, &res, query, id)
	if err != nil {
		db.log.Error("Fetch keywords error", "error", err, "id", id)
		return core.Comics{}, nil, fmt.Errorf("fetch keywords: %w", err)
	}

	var keywords []string

	if err := res.Keywords.AssignTo(&keywords); err != nil {
		return core.Comics{}, nil, fmt.Errorf("convert keywords: %w", err)
	}

	return core.Comics{ID: res.ID, URL: res.URL}, keywords, nil
}

func (db *DB) GetMaxID(ctx context.Context) (int, error) {
	query := `
	SELECT COALESCE(MAX(comics_id), 0)
	FROM comics
	`
	var maxID int
	err := db.conn.GetContext(ctx, &maxID, query)
	return maxID, err
}

func (db *DB) GetComics(ctx context.Context, id int) (core.Comics, error) {
	query := `
	SELECT comics_id, img_url
	FROM comics
	WHERE comics_id = $1
	`

	var comics struct {
		ID  int    `db:"comics_id"`
		URL string `db:"img_url"`
	}

	err := db.conn.GetContext(ctx, &comics, query, id)

	return core.Comics{ID: comics.ID, URL: comics.URL}, err
}
