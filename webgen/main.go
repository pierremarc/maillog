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
	"bytes"
	"encoding/base64"
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
	prefix      string
	what        string
}

var (
	funcMap = template.FuncMap{
		"base64": func(s string) string {
			return base64.StdEncoding.EncodeToString([]byte(s))
		},
	}
)

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
	templatePath := path.Join(os.Getenv("GOPATH"),
		"src/github.com/pierremarc/maillog/webgen", g.what+".tpl")
	rootPath := path.Join(".", g.what)
	globPat := "/*." + g.what

	log.Printf("Loading template: %s", templatePath)
	tb, _ := ioutil.ReadFile(templatePath)
	nt := template.New(g.what)
	nt.Funcs(funcMap)
	t := template.Must(nt.Parse(string(tb)))

	log.Printf("Looking for files matching `%s`", rootPath+globPat)
	fileNames, _ := filepath.Glob(rootPath + globPat)

	for _, fn := range fileNames {
		qn := varName(filepath.Base(fn))
		log.Printf("Got %s in %s", qn, fn)
		bs, _ := ioutil.ReadFile(fn)
		files[qn] = string(bs)
	}

	data := struct {
		Package   string
		Timestamp time.Time
		Prefix    string
		Files     fileMap
	}{
		g.packageName,
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
	log.Printf("Output written to %s", *outputName)
}
