package main

import (
	"encoding/base64"
)

func base64decode(input string) string {
    content, err := base64.StdEncoding.DecodeString(input)
    if err != nil {
        // log.Printf("Error:base64.StdEncoding.DecodeString (%s)", err.Error())
        return input
    }
    return string(content)
}

{{ $Prefix := .Prefix }}
{{ range $Name, $Js :=  .Files }}
var {{ $Prefix }}{{ $Name }} = base64decode(`{{ base64 $Js }}`)
{{ end }}