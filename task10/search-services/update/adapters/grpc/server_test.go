package grpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	updatepb "yadro.com/course/proto/update"
	"yadro.com/course/update/adapters/grpc"
	mock_grpc "yadro.com/course/update/adapters/grpc/mocks"
	"yadro.com/course/update/core"
)

func newMockUpdater(t *testing.T) *mock_grpc.MockUpdater {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})
	return mock_grpc.NewMockUpdater(ctrl)
}

func TestServer_Ping(t *testing.T) {
	upd := newMockUpdater(t)
	srv := grpc.NewServer(upd)

	resp, err := srv.Ping(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)
	require.Nil(t, resp)
}

func TestStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpd := mock_grpc.NewMockUpdater(ctrl)
	mockUpd.
		EXPECT().
		Status(gomock.Any()).
		Return(core.StatusIdle)

	srv := grpc.NewServer(mockUpd)
	reply, err := srv.Status(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)
	require.Equal(t, updatepb.Status_STATUS_IDLE, reply.Status)

	mockUpd.
		EXPECT().
		Status(gomock.Any()).
		Return(core.StatusRunning)

	srv = grpc.NewServer(mockUpd)
	reply, err = srv.Status(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)
	require.Equal(t, updatepb.Status_STATUS_RUNNING, reply.Status)

	mockUpd.
		EXPECT().
		Status(gomock.Any()).
		Return(core.ServiceStatus("something else"))

	srv = grpc.NewServer(mockUpd)
	reply, err = srv.Status(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)
	require.Equal(t, updatepb.Status_STATUS_UNSPECIFIED, reply.Status)
}

func TestServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpd := mock_grpc.NewMockUpdater(ctrl)
	mockUpd.
		EXPECT().
		Update(gomock.Any()).
		Return(nil)

	srv := grpc.NewServer(mockUpd)
	_, err := srv.Update(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)

	mockUpd.
		EXPECT().
		Update(gomock.Any()).
		Return(core.ErrAlreadyExists)

	srv = grpc.NewServer(mockUpd)
	_, err = srv.Update(context.Background(), &emptypb.Empty{})
	grpcErr, ok := status.FromError(err)
	require.True(t, ok, "ожидается grpc.Status")
	require.Equal(t, codes.AlreadyExists, grpcErr.Code())
	require.Equal(t, "update already runs", grpcErr.Message())

	otherErr := errors.New("update failed")
	mockUpd.
		EXPECT().
		Update(gomock.Any()).
		Return(otherErr)

	srv = grpc.NewServer(mockUpd)
	_, err = srv.Update(context.Background(), &emptypb.Empty{})
	require.Error(t, err)
	require.Equal(t, otherErr.Error(), err.Error())
}

func TestStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var expectedStats core.ServiceStats

	expectedStats.WordsTotal = 100
	expectedStats.WordsUnique = 80
	expectedStats.ComicsTotal = 50
	expectedStats.ComicsFetched = 45

	mockUpd := mock_grpc.NewMockUpdater(ctrl)
	mockUpd.
		EXPECT().
		Stats(gomock.Any()).
		Return(expectedStats, nil)

	srv := grpc.NewServer(mockUpd)
	reply, err := srv.Stats(context.Background(), &emptypb.Empty{})

	require.NoError(t, err)
	require.Equal(t, int64(expectedStats.WordsTotal), reply.WordsTotal)
	require.Equal(t, int64(expectedStats.WordsUnique), reply.WordsUnique)
	require.Equal(t, int64(expectedStats.ComicsTotal), reply.ComicsTotal)
	require.Equal(t, int64(expectedStats.ComicsFetched), reply.ComicsFetched)
}

func TestDrop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpd := mock_grpc.NewMockUpdater(ctrl)
	mockUpd.
		EXPECT().
		Drop(gomock.Any()).
		Return(nil)

	srv := grpc.NewServer(mockUpd)
	_, err := srv.Drop(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)

	dropErr := errors.New("drop failed")
	mockUpd.
		EXPECT().
		Drop(gomock.Any()).
		Return(dropErr)

	srv = grpc.NewServer(mockUpd)
	_, err = srv.Drop(context.Background(), &emptypb.Empty{})

	require.Error(t, err)
	require.Equal(t, dropErr.Error(), err.Error())
}
