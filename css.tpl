package main
{{ $Prefix := .Prefix }}
{{ range $Name, $Css :=  .Files }}
const {{ $Prefix }}{{ $Name }} = `{{ $Css }}`
{{ end }}
