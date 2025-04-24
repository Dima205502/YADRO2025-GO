package aaa

import (
	"io"
	"os"
	"testing"
	"time"

	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func setEnv(t *testing.T, key, value string) func() {
	orig := os.Getenv(key)
	err := os.Setenv(key, value)
	require.NoError(t, err)
	return func() {
		_ = os.Setenv(key, orig)
	}
}

func unsetEnv(t *testing.T, key string) func() {
	orig := os.Getenv(key)
	err := os.Unsetenv(key)
	require.NoError(t, err)
	return func() {
		_ = os.Setenv(key, orig)
	}
}

func TestNewSuccess(t *testing.T) {
	restoreUser := setEnv(t, "ADMIN_USER", "admin")
	defer restoreUser()
	restorePass := setEnv(t, "ADMIN_PASSWORD", "password")
	defer restorePass()

	aaaInst, err := New(time.Minute, logger)
	require.NoError(t, err)
	require.Equal(t, "password", aaaInst.users["admin"])
}

func TestNew_MissingAdminUser(t *testing.T) {
	restorePass := setEnv(t, "ADMIN_PASSWORD", "password")
	defer restorePass()
	restoreUser := unsetEnv(t, "ADMIN_USER")
	defer restoreUser()

	_, err := New(time.Minute, logger)
	require.Error(t, err)
	require.Contains(t, err.Error(), "could not get admin user")
}

func TestNew_MissingAdminPassword(t *testing.T) {
	restoreUser := setEnv(t, "ADMIN_USER", "admin")
	defer restoreUser()
	restorePass := unsetEnv(t, "ADMIN_PASSWORD")
	defer restorePass()

	_, err := New(time.Minute, logger)
	require.Error(t, err)
	require.Contains(t, err.Error(), "could not get admin password")
}

func TestLogin_Success(t *testing.T) {
	restoreUser := setEnv(t, "ADMIN_USER", "admin")
	defer restoreUser()
	restorePass := setEnv(t, "ADMIN_PASSWORD", "password")
	defer restorePass()

	instance, err := New(5*time.Minute, logger)
	require.NoError(t, err)

	tokenString, err := instance.Login("admin", "password")
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	parsed, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	require.NoError(t, err)
	claims, ok := parsed.Claims.(*jwt.RegisteredClaims)
	require.True(t, ok)
	require.Equal(t, "superuser", claims.Subject)

	require.WithinDuration(t, time.Now().Add(5*time.Minute), claims.ExpiresAt.Time, time.Second)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	restoreUser := setEnv(t, "ADMIN_USER", "admin")
	defer restoreUser()
	restorePass := setEnv(t, "ADMIN_PASSWORD", "password")
	defer restorePass()

	instance, err := New(5*time.Minute, logger)
	require.NoError(t, err)

	_, err = instance.Login("admin", "wrongpassword")
	require.Error(t, err)
	require.Equal(t, "invalid credentials", err.Error())
}

func TestVerify_Success(t *testing.T) {
	restoreUser := setEnv(t, "ADMIN_USER", "admin")
	defer restoreUser()
	restorePass := setEnv(t, "ADMIN_PASSWORD", "password")
	defer restorePass()

	instance, err := New(5*time.Minute, logger)
	require.NoError(t, err)

	tokenString, err := instance.Login("admin", "password")
	require.NoError(t, err)

	err = instance.Verify(tokenString)
	require.NoError(t, err)
}

func TestVerify_InvalidTokenString(t *testing.T) {
	instance := AAA{
		tokenTTL: 5 * time.Minute,
		users:    map[string]string{"dummy": "dummy"},
		log:      logger,
	}

	err := instance.Verify("not a valid token")
	require.Error(t, err)
	require.Equal(t, "invalid token", err.Error())
}

func TestVerify_InvalidRole(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		Subject:   "not superuser",
	})
	tokenStr, err := token.SignedString([]byte(secretKey))
	require.NoError(t, err)

	instance := AAA{
		tokenTTL: 5 * time.Minute,
		users:    map[string]string{"admin": "password"},
		log:      logger,
	}

	err = instance.Verify(tokenStr)
	require.Error(t, err)
	require.Equal(t, "invalid role", err.Error())
}

func TestVerify_ExpiredToken(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-5 * time.Minute)),
		Subject:   "superuser",
	})
	tokenStr, err := token.SignedString([]byte(secretKey))
	require.NoError(t, err)

	instance := AAA{
		tokenTTL: 5 * time.Minute,
		users:    map[string]string{"admin": "password"},
		log:      logger,
	}

	err = instance.Verify(tokenStr)
	require.Error(t, err)
	require.Equal(t, "invalid token", err.Error())
}
