package main

import (
	"bytes"
	"flag"
	"go/build"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type generator struct {
	packageName string
	prefix      string
	what        string
}

func varName(n string) string {
	basename := strings.Split(n, ".")[0]
	parts := strings.Split(basename, "-")
	name := ""
	for _, p := range parts {
		name += strings.Title(p)
	}
	return name
}

type fileMap map[string]string

func (g *generator) generate() ([]byte, error) {
	var files = make(fileMap)
	templatePath := path.Join(".", g.what+".tpl")
	rootPath := path.Join(".", g.what)
	globPat := "/*." + g.what

	tb, _ := ioutil.ReadFile(templatePath)
	t := template.Must(template.New(g.what).Parse(string(tb)))

	log.Printf("Looking for files matching `%s`", rootPath+globPat)
	fileNames, _ := filepath.Glob(rootPath + globPat)

	for _, fn := range fileNames {
		qn := varName(filepath.Base(fn))
		log.Printf("Got %s in %s", qn, fn)
		bs, _ := ioutil.ReadFile(fn)
		files[qn] = string(bs)
	}

	data := struct {
		Timestamp time.Time
		Prefix    string
		Files     fileMap
	}{
		time.Now().UTC(),
		g.prefix,
		files,
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

	outputName := flag.String("output", "queries.go", "Output file name")
	what := flag.String("what", "", "What - will build template name (what.tpl), dirpath (./what/) and glob pattern (*.what)")
	prefix := flag.String("prefix", "", "Variable name prefix")

	flag.Parse()

	pkg, err := build.Default.ImportDir(".", 0)
	if err != nil {
		log.Fatal(err)
	}

	var (
		g generator
	)

	// g.queries = flag.Args()
	g.what = *what
	g.packageName = pkg.Name

	if *prefix != "" {
		g.prefix = *prefix
	} else {
		g.prefix = strings.Title(*what)
	}

	src, err := g.generate()
	if err != nil {
		log.Fatal(err)
	}

	if err = ioutil.WriteFile(*outputName, src, 0644); err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
