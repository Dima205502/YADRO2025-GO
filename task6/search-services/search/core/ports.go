package core

import "context"

type Searcher interface {
	Search(context.Context, int, string) ([]Comics, error)
}

type DB interface {
	Find(context.Context, []string, int) ([]Comics, error)
}

type Words interface {
	Norm(context.Context, string) ([]string, error)
}
