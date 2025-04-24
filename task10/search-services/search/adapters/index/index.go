//go:generate mockgen -package mock_index -destination ./mocks/index.go yadro.com/course/search/core Builder
package index

import (
	"context"
	"log/slog"
	"time"

	"yadro.com/course/search/core"
)

type Index struct {
	log        *slog.Logger
	builder    core.Builder
	ttl        time.Duration
	wordToID   map[string][]int
	idToComics map[int]core.Comics
}

func NewIndex(log *slog.Logger, builder core.Builder, ttl time.Duration) (*Index, error) {
	return &Index{
		log:        log,
		builder:    builder,
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
	wordToID, idToComics, err := i.builder.BuildIndex(ctx)
	if err != nil {
		i.log.Error("First builder index initiator failed", "error", err)
	}

	i.wordToID = wordToID
	i.idToComics = idToComics

	i.log.Info("Start index initiator")

	ticker := time.NewTicker(i.ttl)
	go func() {
		for {
			select {
			case <-ticker.C:
				if wordToID, idToComics, err = i.builder.BuildIndex(ctx); err != nil {
					i.log.Error("Index build failed", "error", err)
				} else {
					i.wordToID = wordToID
					i.idToComics = idToComics
					i.log.Info("Index build complete")
				}
			case <-ctx.Done():
				i.log.Info("Index initiator stopped")
				return
			}
		}
	}()
}
