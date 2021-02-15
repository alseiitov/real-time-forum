package domain

import "time"

type Session struct {
	UserID       int
	RefreshToken string
	ExpiresAt    time.Time
}
