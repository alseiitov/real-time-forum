package model

import "time"

type Session struct {
	UserID       int       `json:"userID"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
