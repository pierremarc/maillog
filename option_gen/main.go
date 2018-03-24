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
	"strings"
	"text/template"
	"time"
)

var (
	funcMap = template.FuncMap{
		"first": func(s string) string {
			return strings.ToLower(string(s[0]))
		},
	}
)

type typeMap map[string]string

type generator struct {
	packageName string
	types       typeMap
}

func (g *generator) generate(templatePath string) ([]byte, error) {
	bs, _ := ioutil.ReadFile(templatePath)
	t := template.Must(template.New("option").Parse(string(bs)))
	// t := template.Must(template.New("option").Parse(tmpl))

	data := struct {
		Timestamp time.Time
		Types     typeMap
	}{
		time.Now().UTC(),
		g.types,
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

	defaultTemplate := path.Join(os.Getenv("GOPATH"),
		"src/github.com/pierremarc/maillog",
		"option_gen/option-template.tpl")

	templatePath := flag.String("template", defaultTemplate, "template path")
	outputName := flag.String("output", "option.go", "output file name")

	flag.Parse()

	types := make(map[string]string)
	args := flag.Args()
	for _, pair := range args {
		parts := strings.Split(pair, "=")
		label := parts[0]
		typ := parts[1]
		types[label] = typ
	}

	// if len(*typeName) == 0 {
	// 	flag.Usage()
	// 	os.Exit(2)
	// }

	pkg, err := build.Default.ImportDir(".", 0)
	if err != nil {
		log.Fatal(err)
	}

	var (
		g generator
	)

	g.types = types
	g.packageName = pkg.Name

	src, err := g.generate(*templatePath)
	if err != nil {
		log.Fatal(err)
	}

	if err = ioutil.WriteFile(*outputName, src, 0644); err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
