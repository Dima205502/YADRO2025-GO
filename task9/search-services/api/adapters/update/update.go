//go:generate mockgen -package mock_update -destination ./mocks/update.go "yadro.com/course/proto/update" UpdateClient
package update

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"yadro.com/course/api/core"
	updatepb "yadro.com/course/proto/update"
)

type Client struct {
	log    *slog.Logger
	client updatepb.UpdateClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		client: updatepb.NewUpdateClient(conn),
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

func (c Client) Status(ctx context.Context) (core.UpdateStatus, error) {
	statusReply, err := c.client.Status(ctx, &emptypb.Empty{})
	if err != nil {
		c.log.Error("Status", "", err)
		return "", err
	}

	var status core.UpdateStatus

	switch statusReply.Status {
	case updatepb.Status_STATUS_IDLE:
		status = core.StatusUpdateIdle
	case updatepb.Status_STATUS_RUNNING:
		status = core.StatusUpdateRunning
	default:
		status = core.StatusUpdateUnknown
	}

	return status, nil
}

func (c Client) Stats(ctx context.Context) (core.UpdateStats, error) {
	statsReply, err := c.client.Stats(ctx, &emptypb.Empty{})

	if err != nil {
		c.log.Error("Stats", "error", err)
		return core.UpdateStats{}, err
	}

	return core.UpdateStats{
		WordsTotal:    int(statsReply.WordsTotal),
		WordsUnique:   int(statsReply.WordsUnique),
		ComicsTotal:   int(statsReply.ComicsTotal),
		ComicsFetched: int(statsReply.ComicsFetched),
	}, nil
}

func (c Client) Update(ctx context.Context) error {
	_, err := c.client.Update(ctx, &emptypb.Empty{})
	return err
}

func (c Client) Drop(ctx context.Context) error {
	_, err := c.client.Drop(ctx, &emptypb.Empty{})
	return err
}
