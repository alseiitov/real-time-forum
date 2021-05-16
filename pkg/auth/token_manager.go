package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	sjwt "github.com/alseiitov/simple-jwt"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrEmptyJWTSigningKey   = errors.New("JWT_SIGNING_KEY is empty")
	ErrEmptyAccessTokenTTL  = errors.New("empty accessTokenTTL for JWT")
	ErrEmptyRefreshTokenTTL = errors.New("empty refreshTokenTTL for JWT")

	ErrEmptySub     = errors.New("empty sub")
	ErrEmptyRole    = errors.New("empty role")
	ErrEmptyExp     = errors.New("empty exp")
	ErrExpiredToken = errors.New("token has expired")
)

type TokenManager interface {
	NewJWT(userID int, role int) (string, error)
	NewRefreshToken() string
	Parse(token string) (int, int, error)
}

type Manager struct {
	jwtSigningKey   string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewManager(jwtSigningKey string, accessTokenTTL, refreshTokenTTL time.Duration) (*Manager, error) {
	if jwtSigningKey == "" {
		return nil, ErrEmptyJWTSigningKey
	}

	if accessTokenTTL == 0 {
		return nil, ErrEmptyAccessTokenTTL
	}

	if refreshTokenTTL == 0 {
		return nil, ErrEmptyRefreshTokenTTL
	}

	return &Manager{
		jwtSigningKey:   jwtSigningKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}, nil
}

func (m *Manager) NewJWT(userID int, role int) (string, error) {
	jwt := sjwt.New()
	jwt.SetPayload("sub", fmt.Sprintf("%v", userID))
	jwt.SetPayload("exp", fmt.Sprintf("%v", time.Now().Add(m.accessTokenTTL).Unix()))
	jwt.SetPayload("role", fmt.Sprintf("%v", role))

	token, err := jwt.Sign(m.jwtSigningKey)
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

	err = jwt.Verify(token, m.jwtSigningKey)
	if err != nil {
		return -1, -1, err
	}

	subData, ok := jwt.Payload("sub")
	if !ok {
		return -1, -1, ErrEmptySub
	}
	sub, err := strconv.Atoi(fmt.Sprintf("%v", subData))
	if err != nil {
		return -1, -1, err
	}

	roleData, ok := jwt.Payload("role")
	if !ok {
		return -1, -1, ErrEmptyRole
	}

	role, err := strconv.Atoi(fmt.Sprintf("%v", roleData))
	if err != nil {
		return -1, -1, err
	}

	expData, ok := jwt.Payload("exp")
	if !ok {
		return -1, -1, ErrEmptyExp
	}

	exp, err := strconv.ParseInt(fmt.Sprintf("%v", expData), 10, 64)
	if err != nil {
		return -1, -1, err
	}
	tm := time.Unix(exp, 0)

	if time.Now().After(tm) {
		return sub, role, ErrExpiredToken
	}

	return sub, role, nil
}
