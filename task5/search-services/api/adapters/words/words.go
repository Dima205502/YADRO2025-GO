package words

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	wordspb "yadro.com/course/proto/words"
)

type Client struct {
	log    *slog.Logger
	client wordspb.WordsClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		client: wordspb.NewWordsClient(conn),
		log:    log,
	}, nil
}

func (c Client) Norm(ctx context.Context, phrase string) ([]string, error) {
	req := &wordspb.WordsRequest{
		Phrase: phrase,
	}

	resp, err := c.client.Norm(ctx, req)
	if err != nil {
		c.log.Error("failed to normalize phrase", "error", err, "phrase", phrase)
		return nil, err
	}

	return resp.Words, nil
}

func (c Client) Ping(ctx context.Context) error {
	if _, err := c.client.Ping(ctx, &emptypb.Empty{}); err != nil {
		c.log.Error("failed to ping words service", "error", err)
		return err
	}

	return nil
}
