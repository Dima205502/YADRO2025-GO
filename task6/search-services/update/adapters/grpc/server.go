package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	updatepb "yadro.com/course/proto/update"
	"yadro.com/course/update/core"
)

func NewServer(service core.Updater) *Server {
	return &Server{service: service}
}

type Server struct {
	updatepb.UnimplementedUpdateServer
	service core.Updater
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Server) Status(ctx context.Context, _ *emptypb.Empty) (*updatepb.StatusReply, error) {
	status := s.service.Status(ctx)

	var response updatepb.StatusReply

	switch status {
	case core.StatusIdle:
		response.Status = updatepb.Status_STATUS_IDLE
	case core.StatusRunning:
		response.Status = updatepb.Status_STATUS_RUNNING
	default:
		response.Status = updatepb.Status_STATUS_UNSPECIFIED
	}

	return &response, nil
}

func (s *Server) Update(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.service.Update(ctx); err != nil {
		if errors.Is(err, core.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "update already runs")
		}
		return nil, err
	}
	return nil, nil
}

func (s *Server) Stats(ctx context.Context, _ *emptypb.Empty) (*updatepb.StatsReply, error) {
	stats, err := s.service.Stats(ctx)

	return &updatepb.StatsReply{
			WordsTotal:    int64(stats.WordsTotal),
			WordsUnique:   int64(stats.WordsUnique),
			ComicsTotal:   int64(stats.ComicsTotal),
			ComicsFetched: int64(stats.ComicsFetched)},
		err
}

func (s *Server) Drop(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, s.service.Drop(ctx)
}
