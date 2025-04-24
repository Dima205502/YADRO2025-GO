package update

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	updatepb "yadro.com/course/proto/update"

	"log/slog"

	"go.uber.org/mock/gomock"

	mock_update "yadro.com/course/api/adapters/update/mocks"
	"yadro.com/course/api/core"
)

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_update.NewMockUpdateClient(ctrl)

	mockClient.EXPECT().
		Ping(gomock.Any(), &emptypb.Empty{}, gomock.Any()).
		Return(&emptypb.Empty{}, nil)

	cl := Client{
		log:    logger,
		client: mockClient,
	}

	err := cl.Ping(context.Background())
	require.NoError(t, err)
}

func TestStatus(t *testing.T) {
	testCase := []struct {
		name     string
		reply    *updatepb.StatusReply
		expected core.UpdateStatus
	}{
		{
			name:     "UpdateIdle",
			reply:    &updatepb.StatusReply{Status: updatepb.Status_STATUS_IDLE},
			expected: core.StatusUpdateIdle,
		},
		{
			name:     "UpdateRunning",
			reply:    &updatepb.StatusReply{Status: updatepb.Status_STATUS_RUNNING},
			expected: core.StatusUpdateRunning,
		},
		{
			name:     "UpdateUnknown",
			reply:    &updatepb.StatusReply{Status: updatepb.Status_STATUS_UNSPECIFIED},
			expected: core.StatusUpdateUnknown,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mock_update.NewMockUpdateClient(ctrl)

			mockClient.EXPECT().
				Status(gomock.Any(), &emptypb.Empty{}, gomock.Any()).
				Return(tc.reply, nil)

			c := Client{
				log:    logger,
				client: mockClient,
			}

			status, err := c.Status(context.Background())
			require.NoError(t, err)
			require.Equal(t, tc.expected, status)
		})
	}
}

func TestStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_update.NewMockUpdateClient(ctrl)

	reply := &updatepb.StatsReply{
		WordsTotal:    100,
		WordsUnique:   80,
		ComicsTotal:   50,
		ComicsFetched: 45,
	}

	mockClient.EXPECT().
		Stats(gomock.Any(), &emptypb.Empty{}, gomock.Any()).
		Return(reply, nil)

	cl := Client{
		log:    logger,
		client: mockClient,
	}

	stats, err := cl.Stats(context.Background())

	require.NoError(t, err)
	require.Equal(t, core.UpdateStats{
		WordsTotal:    100,
		WordsUnique:   80,
		ComicsTotal:   50,
		ComicsFetched: 45,
	}, stats)
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_update.NewMockUpdateClient(ctrl)
	mockClient.EXPECT().
		Update(gomock.Any(), &emptypb.Empty{}, gomock.Any()).
		Return(&emptypb.Empty{}, nil)

	cl := Client{
		log:    logger,
		client: mockClient,
	}

	err := cl.Update(context.Background())
	require.NoError(t, err)
}

func TestClient_Drop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_update.NewMockUpdateClient(ctrl)
	mockClient.EXPECT().
		Drop(gomock.Any(), &emptypb.Empty{}, gomock.Any()).
		Return(&emptypb.Empty{}, nil)

	cl := Client{
		log:    logger,
		client: mockClient,
	}
	err := cl.Drop(context.Background())
	require.NoError(t, err)
}

func TestClient_Status_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_update.NewMockUpdateClient(ctrl)
	expectedErr := errors.New("status error")
	mockClient.EXPECT().
		Status(gomock.Any(), &emptypb.Empty{}, gomock.Any()).
		Return(nil, expectedErr)

	cl := Client{
		log:    logger,
		client: mockClient,
	}
	_, err := cl.Status(context.Background())
	require.Error(t, err)
	require.Equal(t, expectedErr, err)
}
