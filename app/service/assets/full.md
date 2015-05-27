# {{.Label}}

- ID: {{.Name}}/{{.Version}}
- URL: {{.URL}}

## Tables

{{range .Tables.List}}- [{{.Name}}](#{{.Name}})
{{end}}

{{range .Tables.List}}## {{.Name}} {#{{.Name}}}

{{.Description}}

**Fields**

{{range .Fields.List}}- [{{.Name}}](#{{.Table.Name}}-{{.Name}})
{{end}}
{{range .Fields.List}}#### {{.Name}} {#{{.Table.Name}}-{{.Name}}}

{{if .References}}*Refers to: [{{.References.Field.Table.Name}}](#{{.References.Field.Table.Name}}) / [{{.References.Field.Name}}](#{{.References.Field.Table.Name}}-{{.References.Field.Name}})*{{end}}

{{.Description}}

{{if .Type}}##### Schema

- Type: `{{.Type}}`{{if .Length}}
- Length: {{.Length}}{{end}}{{if .Precision}}
- Precision: {{.Precision}}{{end}}{{if .Scale}}
- Scale: {{.Scale}}{{end}}
{{end}}

{{if .Mappings}}##### Mappings

Model | Version | Table | Field | Comment
------|---------|-------|-------|--------
{{range .Mappings}}[{{.Field.Table.Model.Name}}](/models/{{.Field.Table.Model.Name}}) | [{{.Field.Table.Model.Version}}](/models/{{.Field.Table.Model.Name}}/{{.Field.Table.Model.Version}}) | [{{.Field.Table.Name}}](/models/{{.Field.Table.Model.Name}}/{{.Field.Table.Model.Version}}#{{.Field.Table.Name}}) | [{{.Field.Name}}](/models/{{.Field.Table.Model.Name}}/{{.Field.Table.Model.Version}}#{{.Field.Table.Name}}-{{.Field.Name}}) | {{.Comment}}
{{end}}
{{end}}

{{if .InboundRefs}}##### Inbound References

*Total: {{len .InboundRefs}}*

Table | Field | Name
------|-------|-----
{{range .InboundRefs}}[{{.Field.Table.Name}}](#{{.Field.Table.Name}}) | [{{.Field.Name}}](#{{.Field.Table.Name}}-{{.Field.Name}}) | {{.Name}}
{{end}}
{{end}}

{{end}}
{{end}}
