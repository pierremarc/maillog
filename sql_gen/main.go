// Optional is a tool that generates 'optional' type wrappers around a given type T.
//
// Typically this process would be run using go generate, like this:
//
//	//go:generate optional -type=Foo
//
// running this command
//
//	optional -type=Foo
//
package main

import (
	"bytes"
	"flag"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type generator struct {
	packageName string
	queries     []string
}

func queryName(n string) string {
	basename := strings.Split(n, ".")[0]
	parts := strings.Split(basename, "-")
	name := ""
	for _, p := range parts {
		name += strings.Title(p)
	}
	return name
}

type queryMap map[string]string

func (g *generator) generate() ([]byte, error) {
	templatePath := path.Join(os.Getenv("GOPATH"),
		"src/github.com/pierremarc/maillog",
		"sql_gen/queries.tpl")
	var queries = make(queryMap)
	tb, _ := ioutil.ReadFile(templatePath)
	t := template.Must(template.New("sql").Parse(string(tb)))

	rootPath := path.Join(os.Getenv("GOPATH"),
		"src/github.com/pierremarc/maillog/sql_gen/sql")

	log.Printf("Looking for sql files in %s", rootPath)

	sqlfiles, _ := filepath.Glob(rootPath + "/*.sql")

	for _, fn := range sqlfiles {
		qn := queryName(filepath.Base(fn))
		log.Printf("Got %s in %s", qn, fn)
		bs, _ := ioutil.ReadFile(fn)
		queries[qn] = string(bs)
	}

	data := struct {
		Timestamp time.Time
		Queries   queryMap
	}{
		time.Now().UTC(),
		queries,
	}
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("option: ")

	outputName := flag.String("output", "queries.go", "output file name")

	flag.Parse()

	pkg, err := build.Default.ImportDir(".", 0)
	if err != nil {
		log.Fatal(err)
	}

	var (
		g generator
	)

	g.queries = flag.Args()
	g.packageName = pkg.Name

	src, err := g.generate()
	if err != nil {
		log.Fatal(err)
	}

	if err = ioutil.WriteFile(*outputName, src, 0644); err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
