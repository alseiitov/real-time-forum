package hash

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

type PasswordHasher interface {
	Hash(password string) string
}

type hasher struct {
	salt string
}

func NewHasher(salt string) (*hasher, error) {
	if salt == "" {
		return nil, errors.New("PASSWORD_SALT is empty")
	}
	return &hasher{salt: salt}, nil
}

func (h *hasher) Hash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}
