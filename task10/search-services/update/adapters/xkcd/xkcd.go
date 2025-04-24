package xkcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"yadro.com/course/update/core"
)

type Client struct {
	log    *slog.Logger
	client http.Client
	url    string
}

type ComicsInfo struct {
	ID         int    `json:"num"`
	URL        string `json:"img"`
	Title      string `json:"title"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
}

const lastPath = "/info.0.json"

func NewClient(url string, timeout time.Duration, log *slog.Logger) (*Client, error) {
	log.Debug("New Client", "url", url, "timeout", timeout)

	if url == "" {
		return nil, fmt.Errorf("empty base url specified")
	}

	return &Client{
		client: http.Client{Timeout: timeout},
		log:    log,
		url:    url,
	}, nil
}

func (c Client) Get(ctx context.Context, id int) (core.XKCDInfo, error) {
	return c.get(ctx, fmt.Sprintf("%s/%d/%s", c.url, id, lastPath))
}

func (c Client) LastID(ctx context.Context) (int, error) {
	comics, err := c.get(ctx, c.url+lastPath)
	if err != nil {
		return 0, err
	}
	return comics.ID, nil
}

func (c Client) get(ctx context.Context, url string) (core.XKCDInfo, error) {
	c.log.Debug("get", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return core.XKCDInfo{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return core.XKCDInfo{}, fmt.Errorf("failed to request comics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return core.XKCDInfo{}, core.ErrNotFound
	}

	var info ComicsInfo
	if err = json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return core.XKCDInfo{}, fmt.Errorf("failed to decode comics: %w", err)
	}

	return core.XKCDInfo{
		ID:          info.ID,
		URL:         info.URL,
		Description: info.Title + " " + info.SafeTitle + " " + info.Transcript + " " + info.Alt,
	}, nil
}
