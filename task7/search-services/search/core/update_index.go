package core

import (
	"context"
	"fmt"
	"log/slog"
)

type IndexUpdater struct {
	log     *slog.Logger
	fetcher Fetcher
}

func NewIndexUpdater(log *slog.Logger, fetcher Fetcher) (*IndexUpdater, error) {
	return &IndexUpdater{log: log, fetcher: fetcher}, nil
}

func (i *IndexUpdater) UpdateIndex(ctx context.Context, wordToID map[string][]int, idToComics map[int]Comics) error {
	maxID, err := i.fetcher.GetMaxID(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get MaxID, %w", err)
	}

	for id := 1; id <= maxID; id++ {
		if _, ok := idToComics[id]; ok {
			continue
		}

		comics, keywords, err := i.fetcher.FetchComics(ctx, id)
		if err != nil {
			i.log.Error("Couldn't fetch comics from DB", "error", err, "id", id)
			continue
		}

		idToComics[id] = comics
		for _, word := range keywords {
			wordToID[word] = append(wordToID[word], id)
		}
	}

	return nil
}
