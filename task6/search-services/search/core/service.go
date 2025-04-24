package core

import (
	"context"
	"log/slog"
)

type Service struct {
	log   *slog.Logger
	db    DB
	words Words
}

func NewService(log *slog.Logger, db DB, words Words) *Service {
	return &Service{
		log:   log,
		db:    db,
		words: words,
	}
}

func (s *Service) Search(ctx context.Context, limit int, phrase string) ([]Comics, error) {
	s.log.Debug("Search", "limit", limit, "phrase", phrase)

	words, err := s.words.Norm(ctx, phrase)
	if err != nil {
		s.log.Error("Search normalization", "err", err)
		return nil, err
	}

	s.log.Info("Search words", "words", words)

	comics, err := s.db.Find(ctx, words, limit)
	if err != nil { // уже логируется в db
		return nil, err
	}

	s.log.Info("Search", "result", comics)
	return comics, nil
}
