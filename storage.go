package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"strings"
	"text/template"

	// _ "github.com/lib/pq"
	_ "github.com/jackc/pgx/stdlib"
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

var queryMut sync.Mutex

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
		queryMut.Lock()
		defer queryMut.Unlock()
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
	pool   mxPool
}

func (store Store) Query(name string, args ...interface{}) (*sql.Rows, error) {
	f := store.Tables.q(store.Db, name)
	// store.pool.Get()
	// defer store.pool.Release()
	return f(args...)
}

func (store Store) Register(name string, qs string) {
	log.Printf("Store.Register %s", name)
	namedQueriesMut.Lock()
	namedQueries[name] = qs
	namedQueriesMut.Unlock()
}

type queryFunc func(cb rowsCb, cells ...interface{}) ResultBool

func (qf queryFunc) Exec(cells ...interface{}) {
	qf(RowCallback(func() {}), cells...)
}

func (store Store) QueryFunc(name string, args ...interface{}) queryFunc {
	f := func(cb rowsCb, cells ...interface{}) ResultBool {
		rows, err := store.Query(name, args...)
		if err != nil {
			return ErrBool(err)
		}
		defer func() {
			rows.Close()
		}()
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
	return ResultSqlDBFrom(sql.Open("pgx", connString(dbc))).
		FoldStoreF(
			func(err error) Store {
				log.Fatal(err)
				return Store{}
			},
			func(db *sql.DB) Store {
				return Store{
					Tables: tables,
					Db:     db,
					pool:   newMxPool(4),
				}
			})

}

type rowsCb func(args ...interface{})

type mxChan chan int

type mxPool struct {
	c mxChan
	n int
}

func newMxPool(n int) mxPool {
	c := make(mxChan, n)
	return mxPool{c, n}
}

func (p *mxPool) Get() int {
	log.Println("pool.Get Wait for a free mx")
	var c int
	select {
	case c = <-p.c:
	default:
		c = c + 1
	}
	return c
}

func (p *mxPool) Release() {
	select {
	case p.c <- 1:
	default:
		// let it go, let it go...
	}
}

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
	rows.Close()
}

func RowCallback(f func()) rowsCb {
	return func(args ...interface{}) {
		f()
	}
}
