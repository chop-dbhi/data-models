# {{.Label}}

*ID: {{.Name}}/{{.Version}}*

## Tables
{{range $t := .Tables.List}}- [{{$t.Name}}](#{{$t.Name}})
{{end}}
{{range $t:= .Tables.List}}
### {{$t.Name}}

{{$t.Description}}

Name | Type | References | Description
-----|------|------------|------------
{{range $f:= $t.Fields.List}}{{$f.Name}} | {{$f.Type}} | {{if $f.RefField}}[{{$f.RefField.Table.Name}}](#{{$f.RefField.Table.Name}}) {{$f.RefField.Name}}{{end}} | {{$f.Description}}
{{end}}
{{end}}
