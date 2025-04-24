package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRate_Success(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	handler := Rate(nextHandler, 100)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)
	resp := rr.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, "OK", string(body))
	require.True(t, nextCalled)
}

func TestRate_Failure(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		_, _ = w.Write([]byte("should not be called"))
	})

	handler := Rate(nextHandler, 100)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ctx, cancel := context.WithCancel(req.Context())
	cancel()
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler(rr, req)
	resp := rr.Result()

	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, "server is going down\n", string(body))

	require.False(t, nextCalled)
}
