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

type DbConfig struct {
	Host     string
	User     string
	Name     string
	Password string
}

type Config struct {
	RawMails string
}

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

func (c Config) makeQ(db *sql.DB, name string) Querier {
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

func (c Config) q(db *sql.DB, name string) Querier {
	f, ok := cachedQueriers[name]
	if ok == false {
		f = c.makeQ(db, name)
		cachedQueriers[name] = f
	}

	return f
}

type Store struct {
	Config Config
	Db     *sql.DB
}

func (store Store) Query(name string, args ...interface{}) (*sql.Rows, error) {
	f := store.Config.q(store.Db, name)
	return f(args...)
}

func (store Store) Register(name string, qs string) {
	namedQueriesMut.Lock()
	namedQueries[name] = qs
	namedQueriesMut.Unlock()
}

func connString(dbc DbConfig) string {
	return "host=" + dbc.Host + " user=" + dbc.User + " dbname=" + dbc.Name + " password=" + dbc.Password
}

func NewStore(dbc DbConfig) Store {
	db, err := sql.Open("postgres", connString(dbc))
	if err != nil {
		log.Fatal(err)
	}

	return Store{
		Config: Config{
			RawMails: "raw_emails",
		},
		Db: db,
	}
}
