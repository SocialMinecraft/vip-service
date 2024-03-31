package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Db struct {
	db *sql.DB
}

func Connect(url string) (*Db, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	var re Db
	re.db = db

	return &re, nil
}

func (r *Db) Close() {
	r.db.Close()
}
