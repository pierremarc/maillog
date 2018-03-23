package main

import (
	"database/sql"
	"errors"
	"log"
	"sync"

	"strings"
	"text/template"

	_ "github.com/lib/pq"
)

type Querier func(args ...interface{}) (*sql.Rows, error)

// func (qf Query) prepared(db *sql.DB, qs string) func(...interface{}) (*sql.Rows, error) {
// 	return func(args ...interface{}) (*sql.Rows, error) {
// 		return db.Exec(qs, args)
// 	}
// }

func qes(db *sql.DB, qs string) Querier {
	return func(args ...interface{}) (*sql.Rows, error) {
		return db.Query(qs, args)
	}
}

var namedQueriesMut sync.Mutex
var namedQueries = map[string]string{}

var cachedQueriers = map[string]Querier{}

func noopQ(name string) Querier {
	return func(args ...interface{}) (*sql.Rows, error) {
		return nil, errors.New("Query Does Not Exist: " + name)
	}
}

func (c Tables) makeQ(db *sql.DB, name string) Querier {
	qs, ok := namedQueries[name]
	if ok == false {
		return noopQ(name)
	}

	var builder strings.Builder
	queryTemplate := template.New(name)
	parsed, err := queryTemplate.Parse(qs)
	if err != nil {
		return noopQ(name)
	}

	err = parsed.Execute(&builder, c)
	if err != nil {
		return noopQ(name)
	}

	qs = builder.String()

	return func(args ...interface{}) (*sql.Rows, error) {
		return db.Query(qs, args...)
	}
}

func (c Tables) q(db *sql.DB, name string) Querier {
	f, ok := cachedQueriers[name]
	if ok == false {
		f = c.makeQ(db, name)
		cachedQueriers[name] = f
	}

	return f
}

type Store struct {
	Tables Tables
	Db     *sql.DB
}

func (store Store) Query(name string, args ...interface{}) (*sql.Rows, error) {
	f := store.Tables.q(store.Db, name)
	return f(args...)
}

func (store Store) Register(name string, qs string) {
	namedQueriesMut.Lock()
	namedQueries[name] = qs
	namedQueriesMut.Unlock()
}

type queryFunc func(cb rowsCb, cells ...interface{}) (int, error)

func (store Store) QueryFunc(name string, args ...interface{}) queryFunc {
	f := func(cb rowsCb, cells ...interface{}) (int, error) {
		rows, err := store.Query(name, args...)
		if err != nil {
			return 0, err
		}
		WithRows(rows, cb, cells...)
		return 1, nil
	}

	return f
}

func connString(dbc DbConfig) string {
	return "host=" + dbc.Host + " user=" + dbc.User + " dbname=" + dbc.Name + " password=" + dbc.Password
}

func NewStore(dbc DbConfig, tables Tables) Store {
	log.Println(connString(dbc))
	db, err := sql.Open("postgres", connString(dbc))
	if err != nil {
		log.Fatal(err)
	}

	return Store{
		Tables: tables,
		Db:     db,
	}
}

type rowsCb func(args ...interface{})

func WithRows(rows *sql.Rows, cb rowsCb, args ...interface{}) {
	for rows.Next() {
		scanError := rows.Scan(args...)
		if scanError != nil {
			errMesg := scanError.Error()
			log.Printf("Error [WithRows] %s", errMesg)
		} else {
			cb(args...)
		}
	}
}

func RowCallback(f func()) rowsCb {
	return func(args ...interface{}) {
		f()
	}
}
