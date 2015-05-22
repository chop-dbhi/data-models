# {{.Label}}

- ID: {{.Name}}/{{.Version}}
- URL: {{.URL}}

## Tables

{{range $t := .Tables.List}}- [{{$t.Name}}](#{{$t.Name}})
{{end}}

{{range $t:= .Tables.List}}## {{$t.Name}} {#{{$t.Name}}}

{{$t.Description}}

**Fields**

{{range $f := $t.Fields.List}}- [{{$f.Name}}](#{{$t.Name}}-{{$f.Name}})
{{end}}
{{range $f := $t.Fields.List}}#### {{$f.Name}} {#{{$t.Name}}-{{$f.Name}}}

{{if $f.RefField}}*References: [{{$f.RefField.Table.Name}}](#{{$f.RefField.Table.Name}}) / {{$f.RefField.Name}}*{{end}}

{{$f.Description}}

{{if $f.Type}}##### Schema

- Type: `{{$f.Type}}`{{if $f.Length}}
- Length: {{$f.Length}}{{end}}{{if $f.Precision}}
- Precision: {{$f.Precision}}{{end}}{{if $f.Scale}}
- Scale: {{$f.Scale}}{{end}}
{{end}}

{{if $f.Mappings}}##### Mappings

Model | Version | Table | Field | Comment
------|---------|-------|-------|--------
{{range $m := $f.Mappings}}[{{$m.Field.Table.Model.Name}}](/models/{{$m.Field.Table.Model.Name}}) | [{{$m.Field.Table.Model.Version}}](/models/{{$m.Field.Table.Model.Name}}/{{$m.Field.Table.Model.Version}}) | [{{$m.Field.Table.Name}}](/models/{{$m.Field.Table.Model.Name}}/{{$m.Field.Table.Model.Version}}#{{$m.Field.Table.Name}}) | [{{$m.Field.Name}}](/models/{{$m.Field.Table.Model.Name}}/{{$m.Field.Table.Model.Version}}#{{$m.Field.Table.Name}}-{{$m.Field.Name}}) | {{$m.Comment}}
{{end}}
{{end}}
{{end}}
{{end}}
