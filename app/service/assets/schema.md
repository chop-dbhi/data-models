# {{.Label}}

*ID: {{.Name}}/{{.Version}}*

## Tables
{{range $t := .Tables.List}}- [{{$t.Name}}](#{{$t.Name}})
{{end}}
{{range $t := .Tables.List}}
### {{$t.Name}}

Name | Type | Length | Precision | Scale | Default
-----|------|--------|-----------|-------|--------
{{range $f := $t.Fields.List}}{{$f.Name}} | {{$f.Type}} | {{$f.Length}} | {{$f.Precision}} | {{$f.Scale}} | {{$f.Default}}
{{end}}
{{end}}
