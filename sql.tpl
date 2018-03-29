package main
{{ $Prefix := .Prefix }}
{{ range $Name, $Query :=  .Files }}
const {{ $Prefix }}{{ $Name }} = "{{ $Name }}"
{{ end }}

func RegisterQueries(store Store) {
	{{ range $Name, $Query :=  .Files }}
	store.Register({{ $Prefix }}{{ $Name }}, `{{ $Query }}`)
	{{ end }}
}