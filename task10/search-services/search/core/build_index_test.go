package core

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestBuildIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFetcher := NewMockFetcher(ctrl)

	mockFetcher.
		EXPECT().
		GetMaxID(gomock.Any()).
		Return(3, nil)

	expectedWordToID := map[string][]int{
		"hell": {1, 3},
		"word": {1, 2},
		"run":  {2},
		"job":  {2, 3},
	}

	expectedIdToComics := map[int]Comics{
		1: Comics{
			ID:  1,
			URL: "http:xkcd.com/1.img",
		},
		2: Comics{
			ID:  2,
			URL: "http:xkcd.com/2.img",
		},
		3: Comics{
			ID:  3,
			URL: "http:xkcd.com/3.img",
		},
	}

	gomock.InOrder(
		mockFetcher.
			EXPECT().
			FetchComics(gomock.Any(), 1).
			Return(Comics{
				ID:  1,
				URL: "http:xkcd.com/1.img",
			}, []string{"hell", "word"}, nil),
		mockFetcher.
			EXPECT().
			FetchComics(gomock.Any(), 2).
			Return(Comics{
				ID:  2,
				URL: "http:xkcd.com/2.img",
			}, []string{"word", "run", "job"}, nil),
		mockFetcher.
			EXPECT().
			FetchComics(gomock.Any(), 3).
			Return(Comics{
				ID:  3,
				URL: "http:xkcd.com/3.img",
			}, []string{"hell", "job"}, nil),
	)

	builder, err := NewIndexBuilder(logger, mockFetcher)
	require.NoError(t, err)

	wordToID, idToComic, err := builder.BuildIndex(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedWordToID, wordToID)
	require.Equal(t, expectedIdToComics, idToComic)
}
