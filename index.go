package main

import (
	"log"
	"strconv"

	"github.com/blevesearch/bleve"
)

type IndexRecord struct {
	Subject string
	Body    string
}

type Index struct {
	i bleve.Index
}

func MakeIndex(p string) Index {
	var idx Index
	if pathExists(p) {
		i, err := bleve.Open(p)
		if err != nil {
			log.Fatalf("No Index: %v", err)
		}
		idx.i = i
	} else {
		m := bleve.NewIndexMapping()
		i, err := bleve.New(p, m)
		if err != nil {
			log.Fatalf("No Index: %v", err)
		}
		idx.i = i
	}

	return idx
}

func MakeMemoryIndex() Index {
	m := bleve.NewIndexMapping()
	i, err := bleve.NewMemOnly(m)
	if err != nil {
		log.Fatalf("No Index: %v", err)
	}

	return Index{i}
}

func (idx Index) Push(id string, rec IndexRecord) {
	idx.i.Index(id, rec)
}

func (idx Index) Query(term string) []int {
	query := bleve.NewMatchQuery(term)
	search := bleve.NewSearchRequest(query)
	searchResults, err := idx.i.Search(search)
	if err != nil {
		log.Printf("Search Error %v", err)
		return make([]int, 0)
	}

	results := make([]int, 0)
	for _, hit := range searchResults.Hits {
		ResultInt64From(strconv.ParseInt(hit.ID, 10, 32)).
			Map(func(id int64) {
				results = append(results, int(id))
			})
	}

	return results
}
