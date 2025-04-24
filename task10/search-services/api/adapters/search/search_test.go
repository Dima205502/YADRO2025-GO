package search

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"

	"log/slog"

	mock_search "yadro.com/course/api/adapters/search/mocks"
	"yadro.com/course/api/core"
	searchpb "yadro.com/course/proto/search"
)

var logger = slog.New(slog.NewJSONHandler(io.Discard, nil))

func TestClient_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock_search.NewMockSearchClient(ctrl)

	mockClient.EXPECT().
		Ping(gomock.Any(), &emptypb.Empty{}, gomock.Any()).
		Return(&emptypb.Empty{}, nil)

	c := Client{
		log:    logger,
		client: mockClient,
	}

	err := c.Ping(context.Background())
	require.NoError(t, err)
}

func TestDbSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_search.NewMockSearchClient(ctrl)
	c := Client{
		log:    logger,
		client: mockClient,
	}

	ctx := context.Background()
	limit := 2
	phrase := "test phrase"

	reply := &searchpb.SearchReply{
		Comics: []*searchpb.Comics{
			{Id: 1, Url: "http://example.com/1"},
			{Id: 2, Url: "http://example.com/2"},
		},
	}

	expected := []core.Comics{
		core.Comics{ID: 1, URL: "http://example.com/1"},
		core.Comics{ID: 2, URL: "http://example.com/2"},
	}

	mockClient.EXPECT().
		DbSearch(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(reply, nil)

	comics, err := c.DbSearch(ctx, limit, phrase)
	require.NoError(t, err)
	require.Equal(t, expected, comics)
}

func TestIndexSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_search.NewMockSearchClient(ctrl)
	c := Client{
		log:    logger,
		client: mockClient,
	}

	ctx := context.Background()
	limit := 3
	phrase := "another search"

	reply := &searchpb.SearchReply{
		Comics: []*searchpb.Comics{
			{Id: 3, Url: "http://example.com/3"},
			{Id: 4, Url: "http://example.com/4"},
			{Id: 5, Url: "http://example.com/5"},
		},
	}

	mockClient.EXPECT().
		IndexSearch(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(reply, nil)

	expected := []core.Comics{
		core.Comics{ID: 3, URL: "http://example.com/3"},
		core.Comics{ID: 4, URL: "http://example.com/4"},
		core.Comics{ID: 5, URL: "http://example.com/5"},
	}

	comics, err := c.IndexSearch(ctx, limit, phrase)
	require.NoError(t, err)
	require.Equal(t, expected, comics)
}

func TestDbSearch_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_search.NewMockSearchClient(ctrl)
	c := Client{
		log:    logger,
		client: mockClient,
	}

	ctx := context.Background()
	limit := 1
	phrase := "error test"

	expectedErr := errors.New("db search error")
	mockClient.EXPECT().
		DbSearch(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, expectedErr)

	_, err := c.DbSearch(ctx, limit, phrase)

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
}

func TestIndexSearch_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_search.NewMockSearchClient(ctrl)
	c := Client{
		log:    logger,
		client: mockClient,
	}

	ctx := context.Background()
	limit := 1
	phrase := "error test"

	expectedErr := errors.New("index search error")
	mockClient.EXPECT().
		IndexSearch(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, expectedErr)

	_, err := c.IndexSearch(ctx, limit, phrase)

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
}
