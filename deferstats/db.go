package deferstats

import (
	"database/sql"
	"time"
)

var (
	querylist []DeferDB

	selectThreshold int
)

type DB struct {
	Other *sql.DB
}

func NewDB(db *sql.DB) *DB {
	selectThreshold = 500

	return &DB{
		db,
	}
}

func (db *DB) logQuery(startTime time.Time, query string) {
	endTime := time.Now()
	t := int(((endTime.Sub(startTime)).Nanoseconds() / 1000000))

	ddb := DeferDB{
		Query: query,
		Time:  t,
	}

	if t >= selectThreshold {
		querylist = append(querylist, ddb)
	}
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	startTime := time.Now()
	defer db.logQuery(startTime, query)
	return db.Other.Query(query, args...)
}

func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	startTime := time.Now()
	defer db.logQuery(startTime, query)
	return db.Other.QueryRow(query, args...)
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	startTime := time.Now()
	defer db.logQuery(startTime, query)
	return db.Other.Exec(query, args...)
}
