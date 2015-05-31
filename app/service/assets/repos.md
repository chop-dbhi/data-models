# Repositories

{{range .}}- {{.URL}}
    - Branch: {{.Branch}}
    - Commit: {{.CommitSHA1}}
    - Commit Date: {{.CommitTime}}
    - Fetched: {{.FetchTime}}
{{end}}
