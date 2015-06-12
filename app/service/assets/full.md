# {{.}}

{{if .Description}}{{.Description}}{{end}}

- ID: {{.Name}}/{{.Version}}
- URL: {{.URL}}

## Tables

{{range .Tables.List}}- [{{.}}](#{{.Slug}})
{{end}}

{{range .Tables.List}}## {{.}} {#{{.Slug}}}

{{.Description}}

**Fields**

{{range .Fields.List}}- [{{.}}](#{{.Slug}})
{{end}}
{{range .Fields.List}}#### {{.}} {#{{.Slug}}}

{{if .References}}*Refers to: [{{.References.Field.Table}}](#{{.References.Field.Table.Slug}}) / [{{.References.Field}}](#{{.References.Field.Slug}})*{{end}}

{{.Description}}

{{if .Type}}##### Schema

- Type: `{{.Type}}`{{if .Length}}
- Length: {{.Length}}{{end}}{{if .Precision}}
- Precision: {{.Precision}}{{end}}{{if .Scale}}
- Scale: {{.Scale}}{{end}}
{{end}}

{{if .Mappings}}##### Mappings

Model | Table | Field | Comment
------|-------|-------|--------
{{range .Mappings}}[{{.Field.Table.Model}}](/models/{{.Field.Table.Model.Path}}) | [{{.Field.Table}}](/models/{{.Field.Table.Model.Path}}#{{.Field.Table.Slug}}) | [{{.Field}}](/models/{{.Field.Table.Model.Path}}#{{.Field.Slug}}) | {{.Comment}}
{{end}}
{{end}}

{{if .InboundRefs}}##### Inbound References

*Total: {{len .InboundRefs}}*

Table | Field | Name
------|-------|-----
{{range .InboundRefs}}[{{.Field.Table}}](#{{.Field.Table.Slug}}) | [{{.Field}}](#{{.Field.Slug}}) | {{.}}
{{end}}
{{end}}

{{end}}
{{end}}
