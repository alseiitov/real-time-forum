package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	sjwt "github.com/alseiitov/simple-jwt"
	uuid "github.com/satori/go.uuid"
)

type TokenManager interface {
	NewJWT(userID int, role int) (string, error)
	NewRefreshToken() string
	Parse(token string) (int, int, error)
}

type Manager struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewManager(secretKey string, accessTokenTTL, refreshTokenTTL time.Duration) (*Manager, error) {
	if secretKey == "" {
		return nil, errors.New("empty JWT secret key")
	}

	if accessTokenTTL == 0 {
		return nil, errors.New("empty accessTokenTTL for JWT")
	}

	if refreshTokenTTL == 0 {
		return nil, errors.New("empty refreshTokenTTL for JWT")
	}

	return &Manager{
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}, nil
}

func (m *Manager) NewJWT(userID int, role int) (string, error) {
	jwt := sjwt.New()
	jwt.SetPayload("sub", fmt.Sprintf("%v", userID))
	jwt.SetPayload("exp", fmt.Sprintf("%v", time.Now().Add(m.accessTokenTTL).Unix()))
	jwt.SetPayload("role", fmt.Sprintf("%v", role))

	token, err := jwt.Sign(m.secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *Manager) NewRefreshToken() string {
	token := uuid.NewV4()
	return token.String()
}

func (m *Manager) Parse(token string) (int, int, error) {
	jwt, err := sjwt.Parse(token)
	if err != nil {
		return -1, -1, err
	}

	err = jwt.Verify(token, m.secretKey)
	if err != nil {
		return -1, -1, err
	}

	subData, ok := jwt.Payload("sub")
	if !ok {
		return -1, -1, errors.New("empty sub")
	}
	sub, err := strconv.Atoi(fmt.Sprintf("%v", subData))
	if err != nil {
		return -1, -1, err
	}

	roleData, ok := jwt.Payload("role")
	if !ok {
		return -1, -1, errors.New("empty role")
	}

	role, err := strconv.Atoi(fmt.Sprintf("%v", roleData))
	if err != nil {
		return -1, -1, err
	}

	expData, ok := jwt.Payload("exp")
	if !ok {
		return -1, -1, errors.New("empty exp")
	}

	exp, err := strconv.ParseInt(fmt.Sprintf("%v", expData), 10, 64)
	if err != nil {
		return -1, -1, err
	}
	tm := time.Unix(exp, 0)

	if time.Now().After(tm) {
		return -1, -1, errors.New("token has expired")
	}

	return sub, role, nil
}
