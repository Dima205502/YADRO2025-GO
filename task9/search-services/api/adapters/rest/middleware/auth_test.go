package middleware

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	mock_auth "yadro.com/course/api/adapters/rest/middleware/mocks"
)

func TestAuth_InvalidFormat(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	handler := Auth(next, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	resp := rec.Result()
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	require.Equal(t, "invalid authorization format\n", string(body))
	require.False(t, nextCalled)

	nextCalled = false
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "NoToken")
	rec = httptest.NewRecorder()
	handler(rec, req)
	resp = rec.Result()
	body, _ = io.ReadAll(resp.Body)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	require.Equal(t, "invalid authorization format\n", string(body))
	require.False(t, nextCalled)
}

func TestAuth_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVerifier := mock_auth.NewMockTokenVerifier(ctrl)

	const token = "bad-token"
	mockVerifier.
		EXPECT().
		Verify(token).
		Return(errors.New("verification failed"))

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})
	handler := Auth(next, mockVerifier)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.Header.Set("Authorization", "Token "+token)
	rec := httptest.NewRecorder()
	handler(rec, req)

	resp := rec.Result()
	body, _ := io.ReadAll(resp.Body)

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	require.Equal(t, "invalid token\n", string(body))
	require.False(t, nextCalled)
}

func TestAuth_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const token = "good-token"
	mockVerifier := mock_auth.NewMockTokenVerifier(ctrl)
	mockVerifier.
		EXPECT().
		Verify(token).
		Return(nil)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	handler := Auth(next, mockVerifier)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token "+token)
	rec := httptest.NewRecorder()
	handler(rec, req)

	resp := rec.Result()
	body, _ := io.ReadAll(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "OK", string(body))
	require.True(t, nextCalled)
}
