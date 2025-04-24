package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"log/slog"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mock_auth "yadro.com/course/api/adapters/rest/mocks/auth"
	mock_port "yadro.com/course/api/adapters/rest/mocks/ports"
	"yadro.com/course/api/core"
)

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPingerOk := mock_port.NewMockPinger(ctrl)
	mockPingerErr := mock_port.NewMockPinger(ctrl)

	mockPingerOk.
		EXPECT().
		Ping(gomock.Any()).
		Return(nil)

	mockPingerErr.
		EXPECT().
		Ping(gomock.Any()).
		Return(errors.New("ping failed"))

	pingers := map[string]core.Pinger{
		"service_ok":  mockPingerOk,
		"service_err": mockPingerErr,
	}

	expect := map[string]string{
		"service_ok":  "ok",
		"service_err": "unavailable",
	}

	handler := NewPingHandler(logger, pingers)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var pingResp struct {
		Replies map[string]string `json:"replies"`
	}

	err := json.NewDecoder(res.Body).Decode(&pingResp)
	require.NoError(t, err)
	require.Equal(t, expect, pingResp.Replies)
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockAuthenticator(ctrl)

	mockAuth.
		EXPECT().
		Login("testuser", "testpass").
		Return("token123", nil)

	handler := NewLoginHandler(logger, mockAuth)

	body := `{"name":"testuser", "password":"testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	tokenBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "token123", string(tokenBytes))
}

func TestLogin_InvalidJSON(t *testing.T) {
	handler := NewLoginHandler(logger, nil)
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("invalid json"))
	rec := httptest.NewRecorder()

	handler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// Тестирование UpdateHandler
func TestUpdateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockUpdater := mock_port.NewMockUpdater(ctrl)
	t.Run("StatusOK", func(t *testing.T) {
		mockUpdater.
			EXPECT().
			Update(gomock.Any()).
			Return(nil)

		handler := NewUpdateHandler(logger, mockUpdater)

		req := httptest.NewRequest(http.MethodPost, "/update", nil)
		rec := httptest.NewRecorder()

		handler(rec, req)
		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
		ctrl.Finish()
	})

	// Сценарий 2: updater.Update возвращает ошибку с кодом AlreadyExists (от grpc/status).
	t.Run("StatusAccepted", func(t *testing.T) {
		alreadyExistsErr := status.Error(codes.AlreadyExists, "already exists")
		ctrl = gomock.NewController(t)
		mockUpdater = mock_port.NewMockUpdater(ctrl)

		mockUpdater.
			EXPECT().
			Update(gomock.Any()).
			Return(alreadyExistsErr)

		handler := NewUpdateHandler(logger, mockUpdater)

		req := httptest.NewRequest(http.MethodPost, "/update", nil)
		rec := httptest.NewRecorder()
		handler(rec, req)
		require.Equal(t, http.StatusAccepted, rec.Result().StatusCode)
		ctrl.Finish()
	})

	t.Run("InternalServerError", func(t *testing.T) {
		ctrl = gomock.NewController(t)
		mockUpdater = mock_port.NewMockUpdater(ctrl)
		mockUpdater.
			EXPECT().
			Update(gomock.Any()).
			Return(errors.New("internal server"))

		handler := NewUpdateHandler(logger, mockUpdater)

		req := httptest.NewRequest(http.MethodPost, "/update", nil)
		rec := httptest.NewRecorder()
		handler(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
		ctrl.Finish()
	})

}

// Тестирование UpdateStatsHandler
func TestUpdateStatsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mock_port.NewMockUpdater(ctrl)

	stats := core.UpdateStats{
		WordsTotal:    100,
		WordsUnique:   80,
		ComicsFetched: 45,
		ComicsTotal:   50,
	}

	mockUpdater.
		EXPECT().
		Stats(gomock.Any()).
		Return(stats, nil)

	handler := NewUpdateStatsHandler(logger, mockUpdater)

	req := httptest.NewRequest(http.MethodGet, "/updatestats", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	var resp struct {
		WordsTotal    int `json:"words_total"`
		WordsUnique   int `json:"words_unique"`
		ComicsFetched int `json:"comics_fetched"`
		ComicsTotal   int `json:"comics_total"`
	}

	err := json.NewDecoder(res.Body).Decode(&resp)
	require.NoError(t, err)

	require.Equal(t, 100, resp.WordsTotal)
	require.Equal(t, 80, resp.WordsUnique)
	require.Equal(t, 45, resp.ComicsFetched)
	require.Equal(t, 50, resp.ComicsTotal)
}

// Тестирование UpdateStatusHandler
func TestUpdateStatusHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mock_port.NewMockUpdater(ctrl)

	mockUpdater.
		EXPECT().
		Status(gomock.Any()).
		Return(core.StatusUpdateRunning, nil)

	handler := NewUpdateStatusHandler(logger, mockUpdater)

	req := httptest.NewRequest(http.MethodGet, "/updatestatus", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	var resp map[string]core.UpdateStatus
	err := json.NewDecoder(res.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, core.StatusUpdateRunning, resp["status"])
}

func TestDropHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mock_port.NewMockUpdater(ctrl)

	mockUpdater.
		EXPECT().
		Drop(gomock.Any()).
		Return(nil)

	handler := NewDropHandler(logger, mockUpdater)

	req := httptest.NewRequest(http.MethodPost, "/drop", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	ctrl.Finish()

	ctrl = gomock.NewController(t)
	mockUpdater = mock_port.NewMockUpdater(ctrl)
	mockUpdater.
		EXPECT().
		Drop(gomock.Any()).
		Return(errors.New("drop error"))

	handler = NewDropHandler(logger, mockUpdater)

	req = httptest.NewRequest(http.MethodPost, "/drop", nil)
	rec = httptest.NewRecorder()
	handler(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
}

func TestSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSearcher := mock_port.NewMockSearcher(ctrl)
	handler := NewSearchHandler(logger, core.Searcher(mockSearcher))

	limit := 2
	phrase := "test"

	dbComics := []core.Comics{
		{ID: 1, URL: "http://a"},
		{ID: 2, URL: "http://b"},
	}

	mockSearcher.EXPECT().
		DbSearch(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(
			dbComics,
			nil)

	req := httptest.NewRequest(http.MethodGet, "/search?limit="+strconv.Itoa(limit)+"&phrase="+phrase, nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	expected := struct {
		Comics []struct {
			ID  int    `json:"id"`
			URL string `json:"url"`
		} `json:"comics"`
		Total int `json:"total"`
	}{
		Comics: []struct {
			ID  int    `json:"id"`
			URL string `json:"url"`
		}{
			{ID: 1, URL: "http://a"},
			{ID: 2, URL: "http://b"},
		},
		Total: 2,
	}

	var resp struct {
		Comics []struct {
			ID  int    `json:"id"`
			URL string `json:"url"`
		} `json:"comics"`
		Total int `json:"total"`
	}

	err := json.NewDecoder(res.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, expected, resp)
}

func TestSearchIndexHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSearcher := mock_port.NewMockSearcher(ctrl)
	handler := NewSearchIndexHandler(logger, core.Searcher(mockSearcher))

	limit := 3
	phrase := "index"

	mockSearcher.EXPECT().
		IndexSearch(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(
			[]core.Comics{
				{ID: 3, URL: "http://c"},
				{ID: 4, URL: "http://d"},
				{ID: 5, URL: "http://e"},
			},
			nil)

	req := httptest.NewRequest(http.MethodGet, "/searchindex?limit="+strconv.Itoa(limit)+"&phrase="+phrase, nil)
	rec := httptest.NewRecorder()

	handler(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	expected := struct {
		Comics []struct {
			ID  int    `json:"id"`
			URL string `json:"url"`
		} `json:"comics"`
		Total int `json:"total"`
	}{
		Comics: []struct {
			ID  int    `json:"id"`
			URL string `json:"url"`
		}{
			{ID: 3, URL: "http://c"},
			{ID: 4, URL: "http://d"},
			{ID: 5, URL: "http://e"},
		},
		Total: 3,
	}

	var resp struct {
		Comics []struct {
			ID  int    `json:"id"`
			URL string `json:"url"`
		} `json:"comics"`
		Total int `json:"total"`
	}
	err := json.NewDecoder(res.Body).Decode(&resp)
	require.NoError(t, err)

	require.Equal(t, expected, resp)
}
