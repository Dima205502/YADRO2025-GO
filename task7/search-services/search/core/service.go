package core

import (
	"cmp"
	"context"
	"log/slog"
	"maps"
	"slices"
)

type Service struct {
	log   *slog.Logger
	db    DB
	index Index
	words Words
}

func NewService(log *slog.Logger, db DB, index Index, words Words) (*Service, error) {
	return &Service{
		log:   log,
		db:    db,
		index: index,
		words: words,
	}, nil
}

func (s *Service) DbSearch(ctx context.Context, limit int, phrase string) ([]Comics, error) {
	s.log.Debug("DbSearch", "limit", limit, "phrase", phrase)
	return s.search(ctx, limit, phrase, s.db)
}

func (s *Service) IndexSearch(ctx context.Context, limit int, phrase string) ([]Comics, error) {
	s.log.Debug("IndexSearch", "limit", limit, "phrase", phrase)
	return s.search(ctx, limit, phrase, s.index)
}

func (s *Service) search(ctx context.Context, limit int, phrase string, searcher wordSearcher) ([]Comics, error) {
	words, err := s.words.Norm(ctx, phrase)
	if err != nil {
		s.log.Error("search normalization", "error", err)
		return nil, err
	}

	IDToScore := make(map[int]int)

	for _, word := range words {
		IDs, err := searcher.SearchByWord(ctx, word)
		if err != nil {
			s.log.Error("searchByWord", "error", err)
			return nil, err
		}

		for _, id := range IDs {
			IDToScore[id]++
		}
	}

	sorted := slices.SortedFunc(maps.Keys(IDToScore), func(a, b int) int {
		if IDToScore[a] != IDToScore[b] {
			return cmp.Compare(IDToScore[b], IDToScore[a])
		}

		return cmp.Compare(a, b)
	})

	if len(sorted) < limit {
		limit = len(sorted)
	}
	sorted = sorted[:limit]

	ans := make([]Comics, 0, len(sorted))

	for _, id := range sorted {
		comics, err := searcher.GetComics(ctx, id)
		if err != nil {
			s.log.Error("Can't get comics by ID", "error", err, "id", id)
			return nil, err
		}

		ans = append(ans, comics)
	}

	return ans, nil
}
