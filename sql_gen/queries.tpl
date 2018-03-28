package main

{{ range $Name, $Query :=  .Queries }}
const Query{{ $Name }} = "{{ $Name }}"
{{ end }}

func RegisterQueries(store Store) {
	{{ range $Name, $Query :=  .Queries }}
	store.Register(Query{{ $Name }}, `{{ $Query }}`)
	{{ end }}
}