package db

import (
	"database/sql"
)

func (db *DB) Query(query string, args ...interface{}) (*Rows, error) {
}

func (db *DB) QueryRow(query string, args ...interface{}) *Row {
}

func (db *DB) Exec(query string, args ...interface{}) (Result, error) {
}
