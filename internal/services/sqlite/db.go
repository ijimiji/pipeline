package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func New() *Database {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}

	return &Database{
		DB: db,
	}
}

type Database struct {
	*sql.DB
}

func (d *Database) Close() {
	d.Close()
}
