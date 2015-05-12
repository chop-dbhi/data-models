# {{.Label}}

- ID: {{.Name}}/{{.Version}}
- Tables: {{len .Tables}}

## Table of Contents
{{range $t := .Tables.List}}- [{{$t.Name}}](#{{$t.Name}}){{range $f := $t.Fields.List}}
    - [{{$f.Name}}](#{{$t.Name}}-{{$f.Name}}){{end}}
{{end}}

{{range $t:= .Tables.List}}
## {{$t.Name}} {#{{$t.Name}}}

{{$t.Description}}

- Fields: {{len $t.Fields}}

{{range $f := $t.Fields.List}}
#### {{$f.Name}} {#{{$t.Name}}-{{$f.Name}}}
{{if $f.RefField}}*References: [{{$f.RefField.Table.Name}}](#{{$f.RefField.Table.Name}}) / {{$f.RefField.Name}}*{{end}}

Type | Length | Precision | Scale | Description
-----|--------|-----------|-------|------------
{{$f.Type}} | {{$f.Length}} | {{$f.Precision}} | {{$f.Scale}} | {{$f.Description}}

{{if $f.Mappings}}##### Mappings
Model | Version | Table | Field | Comment
------|---------|-------|-------|--------
{{range $m := $f.Mappings}}{{$m.Field.Table.Model.Name}} | {{$m.Field.Table.Model.Version}} | {{$m.Field.Table.Name}} | {{$m.Field.Name}} | {{$m.Comment}}
{{end}}
{{end}}
{{end}}
{{end}}
