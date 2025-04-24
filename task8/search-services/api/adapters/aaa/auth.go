package aaa

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "something secret here" // token sign key
const adminRole = "superuser"             // token subject

// Authentication, Authorization, Accounting
type AAA struct {
	users    map[string]string
	tokenTTL time.Duration
	log      *slog.Logger
}

func New(tokenTTL time.Duration, log *slog.Logger) (AAA, error) {
	const adminUser = "ADMIN_USER"
	const adminPass = "ADMIN_PASSWORD"

	user, ok := os.LookupEnv(adminUser)
	if !ok {
		return AAA{}, fmt.Errorf("could not get admin user from enviroment")
	}
	password, ok := os.LookupEnv(adminPass)
	if !ok {
		return AAA{}, fmt.Errorf("could not get admin password from enviroment")
	}

	return AAA{
		users:    map[string]string{user: password},
		tokenTTL: tokenTTL,
		log:      log,
	}, nil
}

func (a AAA) Login(name, password string) (string, error) {
	storedPass, ok := a.users[name]
	if !ok || storedPass != password {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenTTL)),
		Subject:   adminRole,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		a.log.Error("failed to sign token", "error", err)
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

func (a AAA) Verify(tokenString string) error {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		a.log.Error("token verification failed", "error", err)
		return errors.New("invalid token")
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.Subject != adminRole {
			return errors.New("invalid role")
		}

		if claims.ExpiresAt.Before(time.Now()) {
			return errors.New("token expired")
		}
	} else {
		return errors.New("invalid token claims")
	}

	return nil
}
