package storage

import (
	"database/sql"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alseiitov/real-time-forum/internal/configs"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func ConnectDB(c *configs.Conf) (*Database, error) {
	if err := checkPath(c.Backend.Database.Path); err != nil {
		return nil, err
	}

	db, err := sql.Open(c.Backend.Database.Driver, c.GetDBFilePath())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := Database{db: db}

	checkTables(&database, c.Backend.Database.Schema)

	return &database, nil
}

func checkPath(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	}
	return nil
}

func checkTables(db *Database, path string) error {
	schema, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	schemes := strings.Split(string(schema), ";\n")

	for _, s := range schemes {
		_, err = db.Exec(string(s))
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) Close() error {
	return db.db.Close()
}

func (db *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(args...)
}

func (db *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Query(args...)
}

func (db *Database) QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args...), nil
}
