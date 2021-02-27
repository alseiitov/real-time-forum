package service

import (
	"log"
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func (s *UsersService) DeleteExpiredSessions() {
	for {
		err := s.repo.DeleteExpiredSessions()
		if err != nil {
			log.Println(err.Error())
		}
		time.Sleep(time.Minute * 10)
	}
}

type UsersRefreshTokensInput struct {
	AccessToken  string
	RefreshToken string
}

func (s *UsersService) RefreshTokens(input UsersRefreshTokensInput) (Tokens, error) {
	sub, role, _ := s.tokenManager.Parse(input.AccessToken)

	err := s.repo.DeleteSession(sub, input.RefreshToken)
	if err != nil {
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
