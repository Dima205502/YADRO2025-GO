package search

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"yadro.com/course/api/core"
	searchpb "yadro.com/course/proto/search"
)

type Client struct {
	log    *slog.Logger
	client searchpb.SearchClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		client: searchpb.NewSearchClient(conn),
		log:    log,
	}, nil
}

func (c Client) Ping(ctx context.Context) error {
	if _, err := c.client.Ping(ctx, &emptypb.Empty{}); err != nil {
		c.log.Error("Ping", "error", err)
		return err
	}

	return nil
}

func (c Client) DbSearch(ctx context.Context, limit int, phrase string) ([]core.Comics, error) {
	c.log.Debug("DbSearch", "limit", limit, "phrase", phrase)

	searchReply, err := c.client.DbSearch(ctx, &searchpb.SearchRequest{Limit: int64(limit), Phrase: phrase})
	if err != nil {
		c.log.Error("Search", "error", err)
		return nil, err
	}

	comicsReply := searchReply.Comics

	var comics []core.Comics
	for _, x := range comicsReply {
		comics = append(comics, core.Comics{ID: int(x.Id), URL: x.Url})
	}

	return comics, nil
}

func (c Client) IndexSearch(ctx context.Context, limit int, phrase string) ([]core.Comics, error) {
	c.log.Debug("IndexSearch", "limit", limit, "phrase", phrase)

	searchReply, err := c.client.IndexSearch(ctx, &searchpb.SearchRequest{Limit: int64(limit), Phrase: phrase})
	if err != nil {
		c.log.Error("IndexSearch", "error", err)
		return nil, err
	}

	comicsReply := searchReply.Comics

	var comics []core.Comics
	for _, x := range comicsReply {
		comics = append(comics, core.Comics{ID: int(x.Id), URL: x.Url})
	}

	return comics, nil
}
