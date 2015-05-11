# {{.Label}}

*ID: {{.Name}}/{{.Version}}*

## Tables
{{range $table := .Tables.List}}- [{{$table.Name}}](#{{$table.Name}})
{{end}}
{{range $table := .Tables.List}}
### {{$table.Name}}

Name | Type | Description
-----|------|------------
{{range $field := $table.Fields.List}}{{$field.Name}} | {{$field.Schema.Type}} | {{$field.Description}}
{{end}}
{{end}}
