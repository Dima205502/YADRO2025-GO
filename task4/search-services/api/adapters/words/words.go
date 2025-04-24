package words

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"yadro.com/course/api/core"
	wordspb "yadro.com/course/proto/words"
)

type Client struct {
	log    *slog.Logger
	client wordspb.WordsClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	log.Debug("NewClient", "address", address)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := wordspb.NewWordsClient(conn)

	return &Client{
		log:    log,
		client: client,
	}, nil
}

func (c Client) Norm(ctx context.Context, phrase string) ([]string, error) {
	req := &wordspb.WordsRequest{
		Phrase: phrase,
	}

	resp, err := c.client.Norm(ctx, req)
	if err != nil {
		if status.Code(err) == codes.ResourceExhausted {
			return nil, core.ErrBadArguments
		}

		c.log.Error("failed to normalize phrase", "error", err, "phrase", phrase)
		return nil, err
	}

	return resp.GetWords(), nil
}

func (c Client) Ping(ctx context.Context) error {
	if _, err := c.client.Ping(ctx, &emptypb.Empty{}); err != nil {
		c.log.Error("failed to ping words service", "error", err)
		return err
	}

	return nil
}
