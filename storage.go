/*
 *  Copyright (C) 2018 Pierre Marchand <pierre.m@atelier-cartographique.be>
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as published by
 *  the Free Software Foundation, version 3 of the License.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"strings"
	"text/template"

	"github.com/jackc/pgx"
	// _ "github.com/lib/pq"
)

type Querier func(args ...interface{}) (*pgx.Rows, error)

// func (qf Query) prepared(db *pgx.DB, qs string) func(...interface{}) (*pgx.Rows, error) {
// 	return func(args ...interface{}) (*pgx.Rows, error) {
// 		return db.Exec(qs, args)
// 	}
// }

// func qes(db *pgx.DB, qs string) Querier {
// 	return func(args ...interface{}) (*pgx.Rows, error) {
// 		return db.Query(qs, args)
// 	}
// }

var namedQueriesMut sync.Mutex
var namedQueries = map[string]string{}

var cachedQueriers = map[string]Querier{}

func noopQ(name string, err string) Querier {
	return func(args ...interface{}) (*pgx.Rows, error) {
		return nil, errors.New(fmt.Sprintf("Noop(%s): %s", name, err))
	}
}

func (c Tables) makeQ(pool *pgx.ConnPool, name string) Querier {
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

	return func(args ...interface{}) (*pgx.Rows, error) {
		conn, err := pool.Acquire()
		if err != nil {
			log.Println("Error acquiring connection:", err)
		}
		defer pool.Release(conn)
		// log.Printf("Query %s", qs)
		// log.Println(args...)
		return pool.Query(qs, args...)
	}
}

func (c Tables) q(pool *pgx.ConnPool, name string) Querier {
	f, ok := cachedQueriers[name]
	if ok == false {
		f = c.makeQ(pool, name)
		cachedQueriers[name] = f
	}

	return f
}

type Store struct {
	Tables Tables
	pool   *pgx.ConnPool
}

func (store Store) Query(name string, args ...interface{}) (*pgx.Rows, error) {
	f := store.Tables.q(store.pool, name)
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

func (qf queryFunc) Exec(cells ...interface{}) ResultBool {
	return qf(RowCallback(func() {}), cells...)
}

func (store Store) QueryFunc(name string, args ...interface{}) queryFunc {
	f := func(cb rowsCb, cells ...interface{}) ResultBool {
		rows, err := store.Query(name, args...)
		if err != nil {
			log.Printf("Query Error(%s): %s", name, err.Error())
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

func connString(dbc DbConfig) pgx.ConnConfig {
	return pgx.ConnConfig{
		Host:     dbc.Host,
		User:     dbc.User,
		Database: dbc.Name,
		Password: dbc.Password,
	}
}

func NewStore(dbc DbConfig, tables Tables) Store {
	config := pgx.ConnPoolConfig{
		ConnConfig:     connString(dbc),
		MaxConnections: 24,
		AcquireTimeout: 30000000000,
	}
	return ResultConnPoolFrom(pgx.NewConnPool(config)).
		FoldStoreF(
			func(err error) Store {
				log.Fatal(err)
				return Store{}
			},
			func(pool *pgx.ConnPool) Store {
				return Store{
					Tables: tables,
					pool:   pool,
				}
			})

}

type rowsCb func(args ...interface{})

func WithRows(rows *pgx.Rows, cb rowsCb, args ...interface{}) {
	defer rows.Close()
	for rows.Next() {
		scanError := rows.Scan(args...)
		if scanError != nil {
			log.Printf("Error [WithRows] %s", scanError.Error())
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
