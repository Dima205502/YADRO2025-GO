package index

import (
	"context"
	"io"
	"testing"
	"time"

	"log/slog"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mock_index "yadro.com/course/search/adapters/index/mocks"
	"yadro.com/course/search/core"
	// Импортируйте путь, соответствующий расположению сгенерированного мока
)

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestNewIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBuilder := mock_index.NewMockBuilder(ctrl)

	idx, err := NewIndex(logger, mockBuilder, time.Second)
	require.NoError(t, err)
	require.NotNil(t, idx)
	require.Empty(t, idx.wordToID)
	require.Empty(t, idx.idToComics)
}

func TestSearchByWord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBuilder := mock_index.NewMockBuilder(ctrl)
	idx, _ := NewIndex(logger, mockBuilder, time.Second)

	idx.wordToID = map[string][]int{
		"hello": {1, 2, 3},
	}

	IDs, err := idx.SearchByWord(context.Background(), "hello")
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3}, IDs)

	IDs, err = idx.SearchByWord(context.Background(), "world")
	require.NoError(t, err)
	require.Empty(t, IDs)
}

func TestGetComics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBuilder := mock_index.NewMockBuilder(ctrl)
	idx, _ := NewIndex(logger, mockBuilder, time.Second)

	comic := core.Comics{
		ID:  10,
		URL: "http://example.com/img",
	}
	idx.idToComics = map[int]core.Comics{
		10: comic,
	}

	ret, err := idx.GetComics(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, comic, ret)

	_, err = idx.GetComics(context.Background(), 20)
	require.Error(t, err)
	require.Equal(t, core.ErrNotFound, err)
}
