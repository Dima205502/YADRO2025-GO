package db

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	mock_dbops "yadro.com/course/search/adapters/db/mocks"
)

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestSearchByWord(t *testing.T) {
	testCase := []struct {
		name     string
		expected error
	}{
		{
			name:     "success",
			expected: nil,
		},
		{
			name:     "unexpected error",
			expected: errors.New("unexpected error"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDBops := mock_dbops.NewMockDBops(ctrl)
			mockDBops.EXPECT().
				SelectContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(tc.expected)

			db := DB{
				log:  logger,
				conn: mockDBops,
			}

			_, err := db.SearchByWord(context.Background(), "keywords")
			assert.ErrorIs(t, err, tc.expected)
		})
	}
}

func TestFetchComics_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mock_dbops.NewMockDBops(ctrl)

	d := &DB{
		log:  logger,
		conn: mockConn,
	}
	ctx := context.Background()
	id := 42

	mockConn.
		EXPECT().
		GetContext(ctx, gomock.Any(), gomock.Any(), id).
		DoAndReturn(func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
			r, ok := dest.(*comicRow)
			if !ok {
				return errors.New("unexpected dest type")
			}
			r.ID = id
			r.URL = "http://xkcd.com/img"
			// Имитируем установку массива ключевых слов.
			var arr pgtype.TextArray
			err := arr.Set([]string{"action", "thriller"})
			if err != nil {
				return err
			}
			r.Keywords = arr
			return nil
		})

	comic, keywords, err := d.FetchComics(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, comic.ID)
	require.Equal(t, "http://xkcd.com/img", comic.URL)
	require.Equal(t, []string{"action", "thriller"}, keywords)
}

func TestFetchComics_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mock_dbops.NewMockDBops(ctrl)
	d := DB{
		log:  logger,
		conn: mockConn,
	}
	ctx := context.Background()
	id := 100

	expectedErr := errors.New("unexpected error")
	mockConn.
		EXPECT().
		GetContext(ctx, gomock.Any(), gomock.Any(), id).
		Return(expectedErr)

	_, _, err := d.FetchComics(ctx, id)
	require.ErrorIs(t, err, expectedErr)
}

func TestGetMaxID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mock_dbops.NewMockDBops(ctrl)
	d := DB{
		log:  logger,
		conn: mockConn,
	}
	ctx := context.Background()

	// Ожидаем вызов GetContext, который записывает максимальный ID.
	mockConn.
		EXPECT().
		GetContext(ctx, gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
			ptr, ok := dest.(*int)
			if !ok {
				return errors.New("expected *int as dest")
			}
			*ptr = 100
			return nil
		})

	maxID, err := d.GetMaxID(ctx)
	require.NoError(t, err)
	require.Equal(t, 100, maxID)
}

func TestGetMaxID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mock_dbops.NewMockDBops(ctrl)
	d := DB{
		log:  logger,
		conn: mockConn,
	}
	ctx := context.Background()
	expectedErr := errors.New("unexpected error")

	mockConn.
		EXPECT().
		GetContext(ctx, gomock.Any(), gomock.Any()).
		Return(expectedErr)

	_, err := d.GetMaxID(ctx)
	require.ErrorIs(t, err, expectedErr)
}

func TestGetComics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mock_dbops.NewMockDBops(ctrl)
	d := DB{
		log:  logger,
		conn: mockConn,
	}
	ctx := context.Background()
	id := 7

	mockConn.
		EXPECT().
		GetContext(ctx, gomock.Any(), gomock.Any(), id).
		DoAndReturn(func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
			r, ok := dest.(*comicsInf)
			if !ok {
				return errors.New("unexpected dest type")
			}
			r.ID = id
			r.URL = "http://example.com/comic7.jpg"
			return nil
		})

	comic, err := d.GetComics(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, comic.ID)
	require.Equal(t, "http://example.com/comic7.jpg", comic.URL)
}

func TestGetComics_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mock_dbops.NewMockDBops(ctrl)

	expectedErr := errors.New("unexpected error")
	// Ожидаем, что GetContext вернет данные для метода GetComics.
	mockConn.
		EXPECT().
		GetContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(expectedErr)

	d := DB{
		log:  logger,
		conn: mockConn,
	}

	_, err := d.GetComics(context.Background(), 0)
	require.ErrorIs(t, err, expectedErr)
}
