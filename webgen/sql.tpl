package {{ .Package }}

import (
	"encoding/base64"
)

func base64decodeSql(input string) string {
    content, err := base64.StdEncoding.DecodeString(input)
    if err != nil {
        return input
    }
    return string(content)
}

{{ $Prefix := .Prefix }}
{{ range $Name, $Query :=  .Files }}
const {{ $Prefix }}{{ $Name }} = "{{ $Name }}"
{{ end }}

func RegisterQueries(store Store) {
	{{ range $Name, $Query :=  .Files }}
	store.Register({{ $Prefix }}{{ $Name }}, base64decodeSql(`{{ base64 $Query }}`))
	{{ end }}
}