package words

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"yadro.com/course/proto/words"
	mock_words "yadro.com/course/search/adapters/words/mocks"
	"yadro.com/course/search/core"
)

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestNormSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_words.NewMockWordsClient(ctrl)

	phrase := "The quick brown fox jumps over the lazy dog; however, the dog doesn’t even notice."
	expected := []string{"brown", "fox", "jump", "lazi", "dog", "howev", "quick", "doesn", "notic", "even"}

	mockClient.
		EXPECT().
		Norm(gomock.Any(), &words.WordsRequest{Phrase: phrase}, gomock.Any()).
		Return(&words.WordsReply{
			Words: []string{"brown", "fox", "jump", "lazi", "dog", "howev", "quick", "doesn", "notic", "even"},
		}, nil)

	c := Client{
		log:    logger,
		client: mockClient,
	}

	result, err := c.Norm(context.Background(), phrase)

	require.NoError(t, err)
	require.ElementsMatch(t, expected, result)
}

func TestNormResourceExhausted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_words.NewMockWordsClient(ctrl)

	mockClient.
		EXPECT().
		Norm(gomock.Any(), gomock.Any()).
		Return(nil, status.Error(codes.ResourceExhausted, ""))

	c := Client{
		log:    logger,
		client: mockClient,
	}

	result, err := c.Norm(context.Background(), "")

	require.Equal(t, core.ErrBadArguments, err)
	require.Nil(t, result)
}
