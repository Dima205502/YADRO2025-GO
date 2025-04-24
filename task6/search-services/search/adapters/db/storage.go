package db

import (
	"context"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"yadro.com/course/search/core"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

type Comics struct {
	ID  int    `db:"comics_id"`
	URL string `db:"img_url"`
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

func (db *DB) Find(ctx context.Context, words []string, limit int) ([]core.Comics, error) {
	/* Нахожу строки, у которых в keywords есть хотябы одно пересчение с words,
	сортирую по количеству слов в пересечении от большего к меньшему(ранжирование),
	дальше ставлю LIMIT, сортирую по comics_id и возвращаю.

	Функция array_intersect_count(arr1, arr2) создана в миграции номер 3
	*/

	const query = `
    SELECT comics_id, img_url
    FROM (
        SELECT 
            comics_id, 
            img_url,
            array_intersect_count(keywords, $1::text[]) AS matches 
        FROM 
            comics
        WHERE 
            keywords && $1::text[]
        ORDER BY 
            matches DESC
        LIMIT $2
    ) AS top_comics
    ORDER BY comics_id ASC
	`

	var comics []Comics

	err := db.conn.SelectContext(
		ctx,
		&comics,
		query,
		words,
		limit,
	)

	if err != nil {
		db.log.Error("Find comics error", "error", err)
		return nil, fmt.Errorf("find comics: %w", err)
	}

	coreComics := make([]core.Comics, 0)

	for _, x := range comics {
		coreComics = append(coreComics, core.Comics{ID: x.ID, URL: x.URL})
	}

	return coreComics, nil
}
