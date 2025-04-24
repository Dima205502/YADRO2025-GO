package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

func TestNewService_InvalidConcurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDB(ctrl)
	xkcd := NewMockXKCD(ctrl)
	words := NewMockWords(ctrl)

	_, err := NewService(logger, db, xkcd, words, 0)
	require.Error(t, err)
}

func TestNewService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDB(ctrl)
	xkcd := NewMockXKCD(ctrl)
	words := NewMockWords(ctrl)

	_, err := NewService(logger, db, xkcd, words, 10)
	require.NoError(t, err)
}

func TestUpdate_XKCDLastIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDB(ctrl)
	xkcd := NewMockXKCD(ctrl)
	words := NewMockWords(ctrl)

	concurrency := 2
	svc, err := NewService(logger, db, xkcd, words, concurrency)
	require.NoError(t, err)

	ctx := context.Background()

	xkcd.EXPECT().LastID(gomock.Any()).Return(0, errors.New("xkcd error"))

	err = svc.Update(ctx)
	require.Error(t, err)
}

func TestUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDB(ctrl)
	xkcd := NewMockXKCD(ctrl)
	words := NewMockWords(ctrl)

	concurrency := 10
	svc, err := NewService(logger, db, xkcd, words, concurrency)
	require.NoError(t, err)

	ctx := context.Background()
	lastID := 1000

	xkcd.EXPECT().LastID(gomock.Any()).Return(lastID, nil)

	db.EXPECT().IDs(gomock.Any()).Return([]int{}, nil)

	for id := 1; id <= lastID; id++ {
		comic := XKCDInfo{
			ID:          id,
			URL:         fmt.Sprintf("url%d", id),
			Description: fmt.Sprintf("comic %d", id),
		}
		xkcd.EXPECT().Get(gomock.Any(), id).Return(comic, nil)

		words.EXPECT().Norm(gomock.Any(), comic.Description).
			Return([]string{"comic", fmt.Sprintf("%d", id)}, nil)

		comics := Comics{
			ID:    id,
			URL:   fmt.Sprintf("url%d", id),
			Words: []string{"comic", fmt.Sprintf("%d", id)},
		}
		db.EXPECT().Add(gomock.Any(), comics).Return(nil)
	}

	err = svc.Update(ctx)
	require.NoError(t, err)
}

func TestStats_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDB(ctrl)
	xkcd := NewMockXKCD(ctrl)
	words := NewMockWords(ctrl)

	concurrency := 1
	svc, err := NewService(logger, db, xkcd, words, concurrency)
	require.NoError(t, err)

	ctx := context.Background()
	lastID := 101
	expectedTotal := lastID - 1

	xkcd.EXPECT().LastID(gomock.Any()).Return(lastID, nil)

	dummyDBStats := DBStats{
		WordsTotal:    100,
		WordsUnique:   50,
		ComicsFetched: 20,
	}
	db.EXPECT().Stats(gomock.Any()).Return(dummyDBStats, nil)

	expected := ServiceStats{DBStats: dummyDBStats, ComicsTotal: expectedTotal}

	stats, err := svc.Stats(ctx)
	require.NoError(t, err)

	require.Equal(t, expected, stats)
}

func TestStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDB(ctrl)
	xkcd := NewMockXKCD(ctrl)
	words := NewMockWords(ctrl)

	concurrency := 1
	svc, err := NewService(logger, db, xkcd, words, concurrency)
	require.NoError(t, err)

	status1 := svc.Status(context.Background())
	require.Equal(t, StatusIdle, status1)

	svc.updateNow = true

	status2 := svc.Status(context.Background())
	require.Equal(t, StatusRunning, status2)
}

func TestDrop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDB(ctrl)
	xkcd := NewMockXKCD(ctrl)
	words := NewMockWords(ctrl)

	concurrency := 1
	svc, err := NewService(logger, db, xkcd, words, concurrency)
	require.NoError(t, err)

	ctx := context.Background()
	db.EXPECT().Drop(gomock.Any()).Return(nil)

	err = svc.Drop(ctx)
	require.NoError(t, err)
}
