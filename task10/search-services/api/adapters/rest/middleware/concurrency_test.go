package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// slowHandler эмулирует медленный обработчик, который ждёт закрытия канала proceed.
func slowHandler(proceed <-chan struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		<-proceed
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}

func TestConcurrency_AllowsUnderLimit(t *testing.T) {
	limit := 2
	proceed := make(chan struct{})
	nextHandler := slowHandler(proceed)
	handler := Concurrency(nextHandler, limit)

	var wg sync.WaitGroup
	responses := make([]*httptest.ResponseRecorder, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			handler(rec, req)
			responses[i] = rec
		}()
	}

	time.Sleep(10 * time.Millisecond)
	close(proceed)
	wg.Wait()

	for _, rec := range responses {
		resp := rec.Result()
		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "ok", string(body))
	}
}

func TestConcurrency_BlocksOverLimit(t *testing.T) {
	limit := 2
	proceed := make(chan struct{})
	nextHandler := slowHandler(proceed)
	handler := Concurrency(nextHandler, limit)

	var wg sync.WaitGroup

	totalRequests := 3
	responses := make([]*httptest.ResponseRecorder, totalRequests)
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			handler(rec, req)
			responses[i] = rec
		}()
	}

	time.Sleep(10 * time.Millisecond)

	close(proceed)
	wg.Wait()

	okCount := 0
	blockedCount := 0
	for _, rec := range responses {
		resp := rec.Result()
		if resp.StatusCode == http.StatusOK {
			okCount++
		} else if resp.StatusCode == http.StatusServiceUnavailable {
			blockedCount++
		}
	}

	require.Equal(t, 2, okCount)
	require.Equal(t, 1, blockedCount)
}
