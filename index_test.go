package main

import (
	"reflect"
	"testing"

	"github.com/blevesearch/bleve"
)

type testrec struct {
	Content string
}

func TestIndex_Query(t *testing.T) {
	m := bleve.NewIndexMapping()
	idx, err := bleve.NewMemOnly(m)
	if err != nil {
		t.Errorf("Could Not Create Index %v", err)
	}

	idx.Index("12", testrec{"foo bar"})
	idx.Index("13", testrec{"something else"})

	type fields struct {
		i bleve.Index
	}
	type args struct {
		term string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		{
			"Bleve Query",
			fields{idx},
			args{"bar"},
			[]int{12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				i: tt.fields.i,
			}
			if got := idx.Query(tt.args.term); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}
