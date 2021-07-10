package database

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB(driver, dir, fileName, schemesDir string) (*sql.DB, error) {
	isNewDB := !fileExists(filepath.Join(dir, fileName))
	if isNewDB {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	enableForeignKeys := "?_foreign_keys=on&cache=shared&mode=rwc"
	dataSourceName := filepath.Join(dir, fileName) + enableForeignKeys

	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	if isNewDB {
		if err = prepareDB(db, schemesDir); err != nil {
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

		_, err = stmt.Exec()
		if err != nil {
			return err
		}
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
