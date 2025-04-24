package core

import (
	"context"
	"fmt"
	"log/slog"
)

type IndexBuilder struct {
	log     *slog.Logger
	fetcher Fetcher
}

func NewIndexBuilder(log *slog.Logger, fetcher Fetcher) (*IndexBuilder, error) {
	return &IndexBuilder{log: log, fetcher: fetcher}, nil
}

func (i *IndexBuilder) BuildIndex(ctx context.Context) (map[string][]int, map[int]Comics, error) {
	wordToID := make(map[string][]int)
	idToComics := make(map[int]Comics)

	maxID, err := i.fetcher.GetMaxID(ctx)
	if err != nil {
		return wordToID, idToComics, fmt.Errorf("couldn't get MaxID, %w", err)
	}

	for id := 1; id <= maxID; id++ {
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

	return wordToID, idToComics, nil
}
