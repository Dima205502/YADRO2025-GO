package index

import (
	"context"
	"log/slog"
	"time"

	"yadro.com/course/search/core"
)

type Index struct {
	log        *slog.Logger
	updater    core.Updater
	ttl        time.Duration
	wordToID   map[string][]int
	idToComics map[int]core.Comics
}

func NewIndex(log *slog.Logger, updater core.Updater, ttl time.Duration) (*Index, error) {
	return &Index{
		log:        log,
		updater:    updater,
		ttl:        ttl,
		wordToID:   make(map[string][]int),
		idToComics: make(map[int]core.Comics),
	}, nil
}

func (i *Index) SearchByWord(_ context.Context, word string) ([]int, error) {
	IDs := make([]int, 0)

	IDs = append(IDs, i.wordToID[word]...)

	return IDs, nil
}

func (i *Index) GetComics(_ context.Context, id int) (core.Comics, error) {
	comics, ok := i.idToComics[id]
	if ok {
		return comics, nil
	}
	return core.Comics{}, core.ErrNotFound
}

func (i *Index) Start(ctx context.Context) {
	err := i.updater.UpdateIndex(ctx, i.wordToID, i.idToComics)
	if err != nil {
		i.log.Error("First update index initiator failed", "error", err)
	}

	i.log.Info("Start index initiator")

	ticker := time.NewTicker(i.ttl)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := i.updater.UpdateIndex(ctx, i.wordToID, i.idToComics); err != nil {
					i.log.Error("Index update failed", "error", err)
				} else {
					i.log.Info("Index update complete")
				}
			case <-ctx.Done():
				i.log.Info("Index initiator stopped")
				return
			}
		}
	}()
}
