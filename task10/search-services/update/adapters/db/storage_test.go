package db

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	mock_dbops "yadro.com/course/update/adapters/db/mocks"
	"yadro.com/course/update/core"
)

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestAdd(t *testing.T) {
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
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDBops := mock_dbops.NewMockDBops(ctrl)
			mockDBops.
				EXPECT().
				ExecContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, tc.expected)

			db := DB{
				log:  logger,
				conn: mockDBops,
			}

			err := db.Add(context.Background(), core.Comics{})
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestStats(t *testing.T) {
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
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDBops := mock_dbops.NewMockDBops(ctrl)
			mockDBops.
				EXPECT().
				GetContext(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(tc.expected)

			db := DB{
				log:  logger,
				conn: mockDBops,
			}

			_, err := db.Stats(context.Background())
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestIDs(t *testing.T) {
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
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDBops := mock_dbops.NewMockDBops(ctrl)
			mockDBops.
				EXPECT().
				SelectContext(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(tc.expected)

			db := DB{
				log:  logger,
				conn: mockDBops,
			}

			_, err := db.IDs(context.Background())
			assert.ErrorIs(t, err, tc.expected)
		})
	}
}

func TestDrop(t *testing.T) {
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
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDBops := mock_dbops.NewMockDBops(ctrl)
			mockDBops.
				EXPECT().
				ExecContext(gomock.Any(), gomock.Any()).
				Return(nil, tc.expected)

			db := DB{
				log:  logger,
				conn: mockDBops,
			}

			err := db.Drop(context.Background())
			assert.Equal(t, tc.expected, err)
		})
	}
}
