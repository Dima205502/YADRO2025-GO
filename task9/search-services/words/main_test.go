package main

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	wordspb "yadro.com/course/proto/words"
)

var (
	logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	s      = &server{logger: logger}
)

func TestNormSuccess(t *testing.T) {
	req := &wordspb.WordsRequest{Phrase: "The quick brown fox jumps over the lazy dog; however, the dog doesn’t even notice."}

	resp, err := s.Norm(context.Background(), req)

	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"brown", "fox", "jump", "lazi", "dog", "howev", "quick", "doesn", "notic", "even"}, resp.Words)
}

func TestNormSizeExceeded(t *testing.T) {
	longPhrase := strings.Repeat("a", phraseSizeLimit+1)
	_, err := s.Norm(context.Background(), &wordspb.WordsRequest{Phrase: longPhrase})
	assert.Equal(t, codes.ResourceExhausted, status.Code(err))
}
