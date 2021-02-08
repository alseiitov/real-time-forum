package sqlite

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/alseiitov/real-time-forum/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB(c *config.Conf) (*sql.DB, error) {
	driver := c.GetDBDriver()
	fileName := c.GetDBFilePath()
	newDB := !fileExists(fileName)

	db, err := sql.Open(driver, fileName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if newDB {
		if err = prepareDB(db, c.GetDBSchemesDir()); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

func prepareDB(db *sql.DB, schemesDir string) error {
	schemes, err := readSchemes(schemesDir)
	if err != nil {
		return err
	}

	for _, scheme := range schemes {
		stmt, err := db.Prepare(scheme)
		if err != nil {
			return err
		}

		stmt.Exec()
		stmt.Close()
	}

	return nil
}

func readSchemes(schemesDir string) ([]string, error) {
	var schemes []string

	files, err := ioutil.ReadDir(schemesDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileName := filepath.Join(schemesDir, file.Name())
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}

		schemes = append(schemes, string(data))
	}
	return schemes, nil
}
