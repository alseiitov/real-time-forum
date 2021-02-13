package hash

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hashed string) bool
}

type bcryptHasher struct {
	salt string
}

func NewBcryptHasher(salt string) *bcryptHasher {
	return &bcryptHasher{salt: salt}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	salted := h.salt + password

	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(salted), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(encryptedPass), nil
}

func (h *bcryptHasher) Compare(password, hashed string) bool {
	salted := h.salt + password

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(salted))
	if err != nil {
		return false
	}

	return true
}
