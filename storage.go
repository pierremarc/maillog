package main

import (
	"database/sql"
	"errors"
	"fmt"
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

func noopQ(name string, err string) Querier {
	return func(args ...interface{}) (*sql.Rows, error) {
		return nil, errors.New(fmt.Sprintf("Noop(%s): %s", name, err))
	}
}

func (c Tables) makeQ(db *sql.DB, name string) Querier {
	qs, ok := namedQueries[name]
	log.Printf("Tables.makeQ %s %v", name, ok)
	if ok == false {
		return noopQ(name, "Not in namedQueries table")
	}

	var builder strings.Builder
	queryTemplate := template.New(name)
	parsed, err := queryTemplate.Parse(qs)
	if err != nil {
		return noopQ(name, err.Error())
	}

	err = parsed.Execute(&builder, c)
	if err != nil {
		return noopQ(name, err.Error())
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
	log.Printf("Store.Register %s", name)
	namedQueriesMut.Lock()
	namedQueries[name] = qs
	namedQueriesMut.Unlock()
	log.Printf("Register Success %s", name)
}

type queryFunc func(cb rowsCb, cells ...interface{}) ResultBool

func (qf queryFunc) Exec() {
	qf(RowCallback(func() {}))
}

func (store Store) QueryFunc(name string, args ...interface{}) queryFunc {
	f := func(cb rowsCb, cells ...interface{}) ResultBool {
		rows, err := store.Query(name, args...)
		if err != nil {
			return ErrBool(err)
		}
		WithRows(rows, cb, cells...)
		return OkBool(true)
	}

	return f
}

func connString(dbc DbConfig) string {
	return "host=" + dbc.Host + " user=" + dbc.User + " dbname=" + dbc.Name + " password=" + dbc.Password
}

func NewStore(dbc DbConfig, tables Tables) Store {
	log.Println(connString(dbc))
	return ResultSqlDBFrom(sql.Open("postgres", connString(dbc))).
		FoldStoreF(
			func(err error) Store {
				log.Fatal(err)
				return Store{}
			},
			func(db *sql.DB) Store {
				return Store{
					Tables: tables,
					Db:     db,
				}
			})

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