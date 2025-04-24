package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"

	searchpb "yadro.com/course/proto/search"
	mock_grpc "yadro.com/course/search/adapters/grpc/mocks"
	"yadro.com/course/search/core"
)

func TestPing(t *testing.T) {
	srv := NewServer(nil)

	resp, err := srv.Ping(context.Background(), &emptypb.Empty{})

	require.NoError(t, err)
	require.Nil(t, resp)
}

func TestDbSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSearcher := mock_grpc.NewMockSearcher(ctrl)
	limit := 10
	phrase := "test query"

	expectedComics := []core.Comics{
		{ID: 1, URL: "http://example.com/1"},
		{ID: 2, URL: "http://example.com/2"},
	}

	mockSearcher.EXPECT().
		DbSearch(gomock.Any(), limit, phrase).
		Return(expectedComics, nil)

	srv := NewServer(mockSearcher)

	req := &searchpb.SearchRequest{
		Limit:  int64(limit),
		Phrase: phrase,
	}
	reply, err := srv.DbSearch(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, reply)

	expected := &searchpb.SearchReply{
		Comics: []*searchpb.Comics{
			{Id: int64(expectedComics[0].ID), Url: expectedComics[0].URL},
			{Id: int64(expectedComics[1].ID), Url: expectedComics[1].URL},
		},
	}
	require.Equal(t, expected, reply)
}

func TestIndexSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSearcher := mock_grpc.NewMockSearcher(ctrl)
	limit := 5
	phrase := "another query"

	expectedComics := []core.Comics{
		{ID: 3, URL: "http://example.com/3"},
		{ID: 4, URL: "http://example.com/4"},
		{ID: 5, URL: "http://example.com/5"},
	}

	mockSearcher.EXPECT().
		IndexSearch(gomock.Any(), limit, phrase).
		Return(expectedComics, nil)

	srv := NewServer(mockSearcher)

	req := &searchpb.SearchRequest{
		Limit:  int64(limit),
		Phrase: phrase,
	}
	reply, err := srv.IndexSearch(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, reply)
	require.Len(t, reply.Comics, len(expectedComics))

	for i, comic := range expectedComics {
		require.Equal(t, int64(comic.ID), reply.Comics[i].GetId())
		require.Equal(t, comic.URL, reply.Comics[i].GetUrl())
	}
}
