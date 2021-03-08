package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/repository"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type UsersRefreshTokensInput struct {
	AccessToken  string
	RefreshToken string
}

func (s *UsersService) RefreshTokens(input UsersRefreshTokensInput) (Tokens, error) {
	sub, role, _ := s.tokenManager.Parse(input.AccessToken)

	err := s.repo.DeleteSession(sub, input.RefreshToken)
	if err != nil {
		if err == repository.ErrSessionNotFound {
			return Tokens{}, ErrSessionNotFound
		}
		return Tokens{}, err
	}

	return s.setSession(sub, role)
}

func (s *UsersService) setSession(userID, role int) (Tokens, error) {
	accessToken, err := s.tokenManager.NewJWT(userID, role)
	refreshToken := s.tokenManager.NewRefreshToken()
	if err != nil {
		return Tokens{}, err
	}

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	session := model.Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	return tokens, s.repo.SetSession(session)
}
