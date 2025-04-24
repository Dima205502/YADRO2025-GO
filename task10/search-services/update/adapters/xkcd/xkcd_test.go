package xkcd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"yadro.com/course/update/core"
)

const clientURL = "https://xkcd.com"

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestNewClient(t *testing.T) {
	t.Run("empty url", func(t *testing.T) {
		_, err := NewClient("", 5*time.Second, logger)
		require.Error(t, err)
	})

	t.Run("valid url", func(t *testing.T) {
		client, err := NewClient(clientURL, 5*time.Second, logger)
		require.NoError(t, err)
		require.Equal(t, clientURL, client.url)
	})
}

func TestClient_get(t *testing.T) {
	var requestErr error = fmt.Errorf("network error")

	tests := []struct {
		name         string
		targetURL    string
		roundTripFn  roundTripFunc
		expectedInfo core.XKCDInfo
		expectedErr  error
	}{
		{
			name:      "success",
			targetURL: clientURL,
			roundTripFn: func(req *http.Request) (*http.Response, error) {
				info := ComicsInfo{
					ID:         101,
					URL:        clientURL + "/img.png",
					Title:      "Title",
					SafeTitle:  "SafeTitle",
					Transcript: "Transcript",
					Alt:        "Alt",
				}
				bodyBytes, _ := json.Marshal(info)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
				}, nil
			},
			expectedInfo: core.XKCDInfo{
				ID:          101,
				URL:         clientURL + "/img.png",
				Description: "Title SafeTitle Transcript Alt",
			},
			expectedErr: nil,
		},
		{
			name:      "not found",
			targetURL: clientURL + "/notfound",
			roundTripFn: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewReader(nil)),
				}, nil
			},
			expectedInfo: core.XKCDInfo{},
			expectedErr:  core.ErrNotFound,
		},
		{
			name:      "invalid json",
			targetURL: clientURL + "/invalidjson",
			roundTripFn: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("invalid json")),
				}, nil
			},
			expectedInfo: core.XKCDInfo{},
			expectedErr:  errors.New("failed to decode comics"),
		},
		{
			name:      "request error",
			targetURL: clientURL + "/requesterror",
			roundTripFn: func(req *http.Request) (*http.Response, error) {
				return nil, requestErr
			},
			expectedInfo: core.XKCDInfo{},
			expectedErr:  requestErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			httpClient := http.Client{
				Timeout:   5 * time.Second,
				Transport: tc.roundTripFn,
			}

			c := Client{
				log:    logger,
				client: httpClient,
				url:    clientURL,
			}

			info, err := c.get(context.Background(), tc.targetURL)

			if tc.name == "invalid json" {
				require.Error(t, err)
			} else {
				require.ErrorIs(t, err, tc.expectedErr)
				require.Equal(t, tc.expectedInfo, info)
			}
		})
	}
}

func TestLastID(t *testing.T) {
	comicsResp := ComicsInfo{
		ID:         789,
		URL:        clientURL + "/comics/789.png",
		Title:      "T",
		SafeTitle:  "ST",
		Transcript: "transcript",
		Alt:        "alt",
	}
	respBytes, err := json.Marshal(comicsResp)
	require.NoError(t, err)

	client := Client{
		log: logger,
		client: http.Client{
			Timeout: 5 * time.Second,
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(respBytes)),
				}, nil
			}),
		},
		url: clientURL,
	}

	ctx := context.Background()
	id, err := client.LastID(ctx)
	require.NoError(t, err)
	require.Equal(t, comicsResp.ID, id)
}

func TestLastID_Zero(t *testing.T) {
	client := Client{
		log: logger,
		client: http.Client{
			Timeout: 5 * time.Second,
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewReader(nil)),
				}, nil
			}),
		},
		url: clientURL,
	}

	ctx := context.Background()
	id, err := client.LastID(ctx)
	require.Error(t, err)
	require.Equal(t, 0, id)
}

func TestClient_Get(t *testing.T) {
	info := ComicsInfo{
		ID:         101,
		URL:        clientURL + "/img.png",
		Title:      "Title",
		SafeTitle:  "SafeTitle",
		Transcript: "Transcript",
		Alt:        "Alt",
	}
	bodyBytes, _ := json.Marshal(info)

	expected := core.XKCDInfo{
		ID:          101,
		URL:         clientURL + "/img.png",
		Description: "Title" + " " + "SafeTitle" + " " + "Transcript" + " " + "Alt",
	}

	client := Client{
		log: logger,
		client: http.Client{
			Timeout: 5 * time.Second,
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
				}, nil
			}),
		},
		url: clientURL,
	}

	res, err := client.Get(context.Background(), 101)
	require.NoError(t, err)
	require.Equal(t, expected, res)

}
