package main

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
)

func Test_qes(t *testing.T) {
	type args struct {
		db *sql.DB
		qs string
	}
	tests := []struct {
		name string
		args args
		want Querier
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := qes(tt.args.db, tt.args.qs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("qes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_noopQ(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want Querier
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := noopQ(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("noopQ() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_makeQ(t *testing.T) {
	type fields struct {
		RawMails string
	}
	type args struct {
		db   *sql.DB
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Querier
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				RawMails: tt.fields.RawMails,
			}
			if got := c.makeQ(tt.args.db, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.makeQ() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_q(t *testing.T) {
	type fields struct {
		RawMails string
	}
	type args struct {
		db   *sql.DB
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Querier
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				RawMails: tt.fields.RawMails,
			}
			if got := c.q(tt.args.db, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.q() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_Query(t *testing.T) {
	type fields struct {
		Config Config
		Db     *sql.DB
	}
	type args struct {
		name string
		args []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    sql.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := Store{
				Config: tt.fields.Config,
				Db:     tt.fields.Db,
			}
			got, err := store.Query(tt.args.name, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Store.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}
