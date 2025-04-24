package core

import (
	"context"
)

type Searcher interface {
	DbSearch(context.Context, int, string) ([]Comics, error)
	IndexSearch(context.Context, int, string) ([]Comics, error)
}

type wordSearcher interface {
	SearchByWord(context.Context, string) ([]int, error)
	GetComics(context.Context, int) (Comics, error)
}

type Index interface {
	Start(context.Context)
	wordSearcher
}

type DB interface {
	wordSearcher
}

type Fetcher interface {
	GetMaxID(context.Context) (int, error)
	FetchComics(context.Context, int) (Comics, []string, error)
}

type Updater interface {
	UpdateIndex(context.Context, map[string][]int, map[int]Comics) error
}

type Words interface {
	Norm(context.Context, string) ([]string, error)
}
