package core

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"
)

// logger декларирован в файле build_index_test.go

func TestDbSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wordsMock := NewMockWords(ctrl)

	dbMock := NewMockDB(ctrl)

	indexDummy := NewMockIndex(ctrl)

	svc, err := NewService(logger, dbMock, indexDummy, wordsMock)
	require.NoError(t, err)

	ctx := context.Background()
	phrase := "Hello"

	wordsMock.EXPECT().
		Norm(ctx, phrase).
		Return([]string{"hello"}, nil)

	dbMock.EXPECT().
		SearchByWord(ctx, "hello").
		Return([]int{1, 2}, nil)

	dbMock.EXPECT().
		GetComics(ctx, 1).
		Return(Comics{ID: 1, URL: "http://a"}, nil)
	dbMock.EXPECT().
		GetComics(ctx, 2).
		Return(Comics{ID: 2, URL: "http://b"}, nil)

	expected := []Comics{Comics{ID: 1, URL: "http://a"}, Comics{ID: 2, URL: "http://b"}}

	results, err := svc.DbSearch(ctx, 10, phrase)
	require.NoError(t, err)
	require.Equal(t, results, expected)
}

func TestIndexSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wordsMock := NewMockWords(ctrl)
	indexMock := NewMockIndex(ctrl)

	dbDummy := NewMockDB(ctrl)

	svc, err := NewService(logger, dbDummy, indexMock, wordsMock)
	require.NoError(t, err)

	ctx := context.Background()
	phrase := "World"

	wordsMock.EXPECT().
		Norm(ctx, phrase).
		Return([]string{"world"}, nil)

	indexMock.EXPECT().
		SearchByWord(ctx, "world").
		Return([]int{3, 4}, nil)

	indexMock.EXPECT().
		GetComics(ctx, 3).
		Return(Comics{ID: 3, URL: "http://c"}, nil)
	indexMock.EXPECT().
		GetComics(ctx, 4).
		Return(Comics{ID: 4, URL: "http://d"}, nil)

	expected := []Comics{Comics{ID: 3, URL: "http://c"}, Comics{ID: 4, URL: "http://d"}}

	results, err := svc.IndexSearch(ctx, 10, phrase)
	require.NoError(t, err)
	require.Equal(t, results, expected)
}

func Test_search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wordsMock := NewMockWords(ctrl)

	dbDummy := NewMockDB(ctrl)

	indexMock := NewMockIndex(ctrl)

	searcherWordMock := NewMockwordSearcher(ctrl)

	svc, err := NewService(logger, dbDummy, indexMock, wordsMock)
	require.NoError(t, err)

	ctx := context.Background()
	phrase := "World"

	wordsMock.EXPECT().
		Norm(ctx, phrase).
		Return([]string{"world"}, nil)

	searcherWordMock.EXPECT().
		SearchByWord(ctx, "world").
		Return([]int{3, 4}, nil)

	searcherWordMock.
		EXPECT().
		GetComics(ctx, 3).
		Return(Comics{ID: 3, URL: "http://c"}, nil)

	searcherWordMock.EXPECT().
		GetComics(ctx, 4).
		Return(Comics{ID: 4, URL: "http://d"}, nil)

	expected := []Comics{Comics{ID: 3, URL: "http://c"}, Comics{ID: 4, URL: "http://d"}}

	res, err := svc.search(ctx, 10, phrase, searcherWordMock)

	require.NoError(t, err)
	require.Equal(t, expected, res)

}
