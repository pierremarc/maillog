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
	"errors"
	"flag"
	"fmt"
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

// https://lawlessguy.wordpress.com/2013/07/23/filling-a-slice-using-command-line-flags-in-go-golang/
type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type typeMap map[string]string

var basicTypes = typeMap{
	"Bool":       "bool",
	"String":     "string",
	"Int":        "int",
	"Int8":       "int8",
	"Int16":      "int16",
	"Int32":      "int32",
	"Int64":      "int64",
	"UInt":       "uint",
	"UInt8":      "uint8",
	"UInt16":     "uint16",
	"UInt32":     "uint32",
	"UInt64":     "uint64",
	"UintPtr":    "uintptr",
	"Byte":       "byte",
	"Rune":       "rune",
	"Float32":    "float32",
	"Float64":    "float64",
	"Complex64":  "complex64",
	"Complex128": "complex128",
}

type generator struct {
	packageName string
	types       typeMap
	imports     []string
}

func (g *generator) generate(templatePath string) ([]byte, error) {
	bs, _ := ioutil.ReadFile(templatePath)
	t := template.Must(template.New("option").Parse(string(bs)))

	data := struct {
		PackageName string
		Timestamp   time.Time
		Types       typeMap
		Imports     []string
	}{
		g.packageName,
		time.Now().UTC(),
		g.types,
		g.imports,
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

	var imports stringSlice

	templateType := flag.String("type", "none", "type to generate [ option | result | array], (required)")
	basics := flag.Bool("basics", false, "generate for basic types")
	flag.Var(&imports, "import", "a package to import, can be repeated")
	outputName := flag.String("output", "", "output file name, default is <type>.go")

	flag.Parse()

	if "none" == *templateType {
		log.Fatal(errors.New("type argument is required"))
	}

	types := make(map[string]string)
	args := flag.Args()
	for _, pair := range args {
		parts := strings.Split(pair, "=")
		label := parts[0]
		typ := parts[1]
		types[label] = typ
	}

	if true == *basics {
		for k, v := range basicTypes {
			types[k] = v
		}
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
	g.imports = imports

	templatepath := path.Join(os.Getenv("GOPATH"),
		"src/github.com/pierremarc/maillog/futil", *templateType)

	src, err := g.generate(templatepath)
	if err != nil {
		log.Fatal(err)
	}

	outPath := fmt.Sprintf("%s.go", *templateType)
	if "" != *outputName {
		outPath = *outputName
	}

	if err = ioutil.WriteFile(outPath, src, 0644); err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
