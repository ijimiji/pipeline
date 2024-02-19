package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
)

func New() *Database {
	db, err := otelsql.Open("sqlite3", "file::memory:?cache=shared",
		otelsql.WithDBName("sqlite"),
	)
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
