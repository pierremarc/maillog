package {{ .Package }}

import (
	"encoding/base64"
)

func base64decodeCss(input string) string {
    content, err := base64.StdEncoding.DecodeString(input)
    if err != nil {
        return input
    }
    return string(content)
}


{{ $Prefix := .Prefix }}
{{ range $Name, $Css :=  .Files }}
var {{ $Prefix }}{{ $Name }} = base64decodeCss(`{{ base64 $Css }}`)
{{ end }}
