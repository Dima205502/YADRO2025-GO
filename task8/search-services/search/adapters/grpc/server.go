package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	searchpb "yadro.com/course/proto/search"
	"yadro.com/course/search/core"
)

type Server struct {
	searchpb.UnimplementedSearchServer
	service core.Searcher
}

func NewServer(service core.Searcher) *Server {
	return &Server{service: service}
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Server) DbSearch(ctx context.Context, in *searchpb.SearchRequest) (*searchpb.SearchReply, error) {
	comics, err := s.service.DbSearch(ctx, int(in.GetLimit()), in.GetPhrase())
	if err != nil {
		return &searchpb.SearchReply{}, err
	}

	comicsResponse := make([]*searchpb.Comics, 0, len(comics))

	for _, x := range comics {
		comicsResponse = append(comicsResponse, &searchpb.Comics{Id: int64(x.ID), Url: x.URL})
	}

	return &searchpb.SearchReply{Comics: comicsResponse}, nil
}

func (s *Server) IndexSearch(ctx context.Context, in *searchpb.SearchRequest) (*searchpb.SearchReply, error) {
	comics, err := s.service.IndexSearch(ctx, int(in.GetLimit()), in.GetPhrase())
	if err != nil {
		return &searchpb.SearchReply{}, err
	}

	comicsResponse := make([]*searchpb.Comics, 0, len(comics))

	for _, x := range comics {
		comicsResponse = append(comicsResponse, &searchpb.Comics{Id: int64(x.ID), Url: x.URL})
	}

	return &searchpb.SearchReply{Comics: comicsResponse}, nil
}
