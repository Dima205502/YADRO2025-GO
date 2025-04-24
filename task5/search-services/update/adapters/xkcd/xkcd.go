package xkcd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type Info struct {
	ID          int    `json:"num"`
	URL         string `json:"img"`
	Title       string `json:"title"`
	Description string `json:"alt"`
}

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
	c.log.Debug("Get", "id", id, "ctx", ctx)
	if id < 1 {
		return core.XKCDInfo{}, core.ErrBadArguments
	}

	url := fmt.Sprintf("%s/%d/info.0.json", c.url, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.log.Error("create request failed", "err", err)
		return core.XKCDInfo{}, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("request failed", "err", err)
		return core.XKCDInfo{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return core.XKCDInfo{}, core.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return core.XKCDInfo{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error("failed to read response body", "err", err)
		return core.XKCDInfo{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var xkcd Info
	if err := json.Unmarshal(body, &xkcd); err != nil {
		c.log.Error("failed to unmarshal response", "err", err)
		return core.XKCDInfo{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.log.Info("Get", "id", id, "xcd", xkcd)

	return core.XKCDInfo{
		ID:          xkcd.ID,
		URL:         xkcd.URL,
		Title:       xkcd.Title,
		Description: xkcd.Description,
	}, nil
}

func (c Client) LastID(ctx context.Context) (int, error) {
	c.log.Debug("LastID", "ctx", ctx)

	url := fmt.Sprintf("%s/info.0.json", c.url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.log.Error("create request failed", "err", err)
		return 0, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("request failed", "err", err)
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, core.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error("failed to read response body", "err", err)
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var xkcd Info
	if err := json.Unmarshal(body, &xkcd); err != nil {
		c.log.Error("failed to unmarshal response", "err", err)
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return xkcd.ID, nil
}
