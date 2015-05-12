# {{.Label}}

*ID: {{.Name}}/{{.Version}}*

## Tables
{{range $t := .Tables.List}}
### {{$t.Name}}

{{range $f := $t.Fields.List}}
{{if $f.Mappings}}
#### {{$f.Name}}
Model | Version | Table | Field | Comment
------|---------|-------|-------|--------
{{range $m := $f.Mappings}}{{$m.Field.Table.Model.Name}} | {{$m.Field.Table.Model.Version}} | {{$m.Field.Table.Name}} | {{$m.Field.Name}} | {{$m.Comment}}
{{end}}
{{end}}
{{end}}
{{end}}
