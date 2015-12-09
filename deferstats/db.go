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

// Tx wraps sql.Tx to provide latency timing against various db
// operations
type Tx struct {
	Db    *DB
	Other *sql.Tx
}

// Stmt wraps sql.Stmt to provide latency timing against various db
// operations
type Stmt struct {
	Db       *DB
	QueryStr string
	Other    *sql.Stmt
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

// NewDB is a constructor for wrapped database
func NewDB(db *sql.DB) *DB {
	selectThreshold = 500

	return &DB{
		db,
	}
}

// NewTx is a constructor for wrapped transaction
func NewTx(db *DB, tx *sql.Tx) *Tx {
	return &Tx{
		Db:    db,
		Other: tx,
	}
}

// NewStmt is a constructor for wrapped statement
func NewStmt(db *DB, querystr string, stmt *sql.Stmt) *Stmt {
	return &Stmt{
		Db:       db,
		QueryStr: querystr,
		Other:    stmt,
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

// Query is a method for returning query results
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	startTime := time.Now()
	defer db.logQuery(startTime, query)
	return db.Other.Query(query, args...)
}

// QueryRow is a method for returning query row
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	startTime := time.Now()
	defer db.logQuery(startTime, query)
	return db.Other.QueryRow(query, args...)
}

// Exec is method for executing query
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	startTime := time.Now()
	defer db.logQuery(startTime, query)
	return db.Other.Exec(query, args...)
}

// Prepare is a method for preparing query
func (db *DB) Prepare(query string) (*Stmt, error) {
	startTime := time.Now()
	defer db.logQuery(startTime, query)
	stmt, err := db.Other.Prepare(query)
	return NewStmt(db, query, stmt), err
}

// Begin is a method for beginning transaction
func (db *DB) Begin() (*Tx, error) {
	startTime := time.Now()
	defer db.logQuery(startTime, "begin")
	tx, err := db.Other.Begin()
	return NewTx(db, tx), err
}

// Commit is a method for committing transaction
func (tx *Tx) Commit() error {
	startTime := time.Now()
	defer tx.Db.logQuery(startTime, "commit")
	return tx.Other.Commit()
}

// Rollback is a method for rollbacking transaction
func (tx *Tx) Rollback() error {
	startTime := time.Now()
	defer tx.Db.logQuery(startTime, "rollback")
	return tx.Other.Rollback()
}

// Query is a method for returning query results
func (tx *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	startTime := time.Now()
	defer tx.Db.logQuery(startTime, query)
	return tx.Other.Query(query, args...)
}

// QueryRow is a method for returning query row
func (tx *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	startTime := time.Now()
	defer tx.Db.logQuery(startTime, query)
	return tx.Other.QueryRow(query, args...)
}

// Exec is method for executing query
func (tx *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	startTime := time.Now()
	defer tx.Db.logQuery(startTime, query)
	return tx.Other.Exec(query, args...)
}

// Prepare is a method for preparing query
func (tx *Tx) Prepare(query string) (*Stmt, error) {
	startTime := time.Now()
	defer tx.Db.logQuery(startTime, query)
	stmt, err := tx.Other.Prepare(query)
	return NewStmt(tx.Db, query, stmt), err
}

// Query is a method for returning query results
func (stmt *Stmt) Query(args ...interface{}) (*sql.Rows, error) {
	startTime := time.Now()
	defer stmt.Db.logQuery(startTime, stmt.QueryStr)
	return stmt.Other.Query(args...)
}

// QueryRow is a method for returning query row
func (stmt *Stmt) QueryRow(args ...interface{}) *sql.Row {
	startTime := time.Now()
	defer stmt.Db.logQuery(startTime, stmt.QueryStr)
	return stmt.Other.QueryRow(args...)
}

// Exec is method for executing query
func (stmt *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	startTime := time.Now()
	defer stmt.Db.logQuery(startTime, stmt.QueryStr)
	return stmt.Other.Exec(args...)
}
