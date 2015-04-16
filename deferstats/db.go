package deferstats

import (
	"database/sql"
	"sync"
	"time"
)

// FIXME
var (
	// Querylist is the list of db_queries that will be sent to
	// deferpanic
	Querylist deferDBList

	// selectThreshold is the size in milliseconds that if a query goes
	// over it will be added to the Querylist
	selectThreshold int
)

// DB wraps sql.DB to provide latency timing against various db
// operations
type DB struct {
	Other *sql.DB
}

// DeferDB holds the query and latency for each sql query whose
// threshold was overran
type DeferDB struct {
	Query string `json:"Query"`
	Time  int    `json:"Time"`
}

// deferDBList is used to keep a list of DeferDB objects
// and interact with them in a thread-safe manner
type deferDBList struct {
	lock sync.RWMutex
	list []DeferDB
}

// Add adds a DeferDB object to the list
func (d *deferDBList) Add(item DeferDB) {
	d.lock.Lock()
	d.list = append(d.list, item)
	d.lock.Unlock()
}

// List returns a copy of the list
func (d *deferDBList) List() []DeferDB {
	d.lock.RLock()
	list := make([]DeferDB, len(d.list))
	for i, v := range d.list {
		list[i] = v
	}
	d.lock.RUnlock()
	return list
}

// Reset removes all entries from the list
func (d *deferDBList) Reset() {
	d.lock.Lock()
	d.list = []DeferDB{}
	d.lock.Unlock()
}

func NewDB(db *sql.DB) *DB {
	selectThreshold = 500

	return &DB{
		db,
	}
}

// logQuery takes a startTime and a query string and if it is over the
// selectThreshold than it appends to a long running query list
func (db *DB) logQuery(startTime time.Time, query string) {
	endTime := time.Now()
	t := int(((endTime.Sub(startTime)).Nanoseconds() / 1000000))

	ddb := DeferDB{
		Query: query,
		Time:  t,
	}

	if t >= selectThreshold {
		Querylist.Add(ddb)
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
