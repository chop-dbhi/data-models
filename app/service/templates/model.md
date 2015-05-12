# {{.Label}}

*ID: {{.Name}}/{{.Version}}*

## Tables
{{range $t := .Tables.List}}- [{{$t.Name}}](#{{$t.Name}})
{{end}}
{{range $t:= .Tables.List}}
### {{$t.Name}}

Name | Type | Description
-----|------|------------
{{range $f:= $t.Fields.List}}{{$f.Name}} | {{$f.Type}} | {{$f.Description}}
{{end}}
{{end}}
