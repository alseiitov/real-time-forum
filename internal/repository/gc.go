package repository

import (
	"database/sql"
	"log"
	"time"
)

func DeleteExpiredSessions(db *sql.DB) {
	for {
		_, err := db.Exec("DELETE FROM sessions WHERE expires_at < $1", time.Now())
		if err != nil {
			log.Println(err)
		}

		time.Sleep(5 * time.Second)
	}
}
