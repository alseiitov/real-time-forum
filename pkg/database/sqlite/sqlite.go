package sqlite

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/config"
)

func ConnectDB(c *config.Conf) (*sql.DB, error) {

	return &sql.DB{}, nil
}
