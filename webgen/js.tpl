package {{ .Package }}

import (
	"encoding/base64"
)

func base64decodeJs(input string) string {
    content, err := base64.StdEncoding.DecodeString(input)
    if err != nil {
        return input
    }
    return string(content)
}

{{ $Prefix := .Prefix }}
{{ range $Name, $Js :=  .Files }}
var {{ $Prefix }}{{ $Name }} = base64decodeJs(`{{ base64 $Js }}`)
{{ end }}