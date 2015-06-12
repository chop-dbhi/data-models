package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const modelFileName = "models.csv"

var dataModelCache ModelIndex

var newlinesRe = regexp.MustCompile(`[\s]+`)

func stripNewlines(s string) string {
	return newlinesRe.ReplaceAllString(s, " ")
}

func rebuildCache() {
	logrus.Debugf("parse: rebuilding cache")

	wg := sync.WaitGroup{}
	wg.Add(len(registeredRepos))

	cache := make(ModelIndex)
	models := make(chan *Model)

	// Find models across repos.
	for _, r := range registeredRepos {
		go func(r *Repo) {
			for _, m := range findModels(r.path) {
				wg.Add(1)
				models <- m
			}

			wg.Done()
		}(r)
	}

	// Spawn 5 workers.
	for i := 0; i < 5; i++ {
		go func() {
			var m *Model

			for {
				select {
				case m = <-models:
					if m == nil {
						return
					}

					parseFiles(m)
					cache.Add(m)

					wg.Done()
				case <-time.After(time.Second):
					logrus.Warnf("model: timed out")
				}
			}
		}()
	}

	wg.Wait()

	close(models)

	// Parse mapping serially since it crosses the model boundary.
	for _, r := range registeredRepos {
		parseMappings(cache, r.path)
	}

	dataModelCache = cache
}

func parseMappings(models ModelIndex, path string) {
	// Load all the definitions files.
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Ignore errors.
		if err != nil {
			return nil
		}

		// Nothing to do with directories.
		if info.IsDir() {
			return nil
		}

		// Skip non-CSV files.
		if filepath.Ext(path) != ".csv" {
			return nil
		}

		f, err := os.Open(path)

		if err != nil {
			return nil
		}

		defer f.Close()

		r := NewMapCSVReader(f)

		if detectFileType(r.Fields()) != MappingsFile {
			return nil
		}

		logrus.Debugf("parse (%s): found mappings file", path)

		// Read all the records.
		records, err := r.ReadAll()

		if err != nil || len(records) == 0 {
			return nil
		}

		var (
			mp     *Mapping
			sm, tm *Model
			st, tt *Table
			sf, tf *Field
		)

		for lineno, r := range records {
			// 1 header + 1-indexed
			lineno += 2

			// Ignore incomplete mappings.
			if r["source_field"] == "" || r["target_field"] == "" {
				logrus.Infof("mapping (%s:%d): incomplete mapping", path, lineno)
				continue
			}

			if sm = models.Get(r["source_model"], r["source_version"]); sm == nil {
				logrus.Warnf("mapping (%s:%d): no model %s/%s", path, lineno, r["source_model"], r["source_version"])
				continue
			}

			if tm = models.Get(r["target_model"], r["target_version"]); tm == nil {
				logrus.Warnf("mapping (%s:%d): no model %s/%s", path, lineno, r["target_model"], r["target_version"])
				continue
			}

			if st = sm.Tables.Get(r["source_table"]); st == nil {
				logrus.Warnf("mapping (%s:%d): no table %s/%s", path, lineno, sm, r["source_table"])
				continue
			}

			if tt = tm.Tables.Get(r["target_table"]); tt == nil {
				logrus.Warnf("mapping (%s:%d): no table %s/%s", path, lineno, tm, r["target_table"])
				continue
			}

			if sf = st.Fields.Get(r["source_field"]); sf == nil {
				logrus.Warnf("mapping (%s:%d): no field %s/%s", path, lineno, st, r["source_field"])
				continue
			}

			if tf = tt.Fields.Get(r["target_field"]); tf == nil {
				logrus.Warnf("mapping (%s:%d): no field %s/%s", path, lineno, tt, r["target_field"])
				continue
			}

			// Bi-directional mapping.
			mp = &Mapping{
				Field:   sf,
				Comment: r["comment"],
			}

			tf.Mappings = append(tf.Mappings, mp)

			mp = &Mapping{
				Field:   tf,
				Comment: r["comment"],
			}

			sf.Mappings = append(sf.Mappings, mp)
		}

		return nil
	})
}

// parseFiles finds and parses all definitions files in the passed directory.
func parseFiles(model *Model) {
	var (
		ok        bool
		table     string
		tableList []Attrs
		refs      []*Reference
	)

	// Initialize
	schema := &Schema{}
	model.schema = schema

	tableFields := make(map[string][]Attrs)
	fieldSchemata := make(TableFieldIndex)

	// Load all the definitions files.
	filepath.Walk(model.path, func(path string, info os.FileInfo, err error) error {
		// Ignore errors.
		if err != nil {
			return nil
		}

		// Nothing to do with directories.
		if info.IsDir() {
			return nil
		}

		// Skip non-CSV files.
		if filepath.Ext(path) != ".csv" {
			return nil
		}

		f, err := os.Open(path)

		if err != nil {
			return nil
		}

		defer f.Close()

		r := NewMapCSVReader(f)

		fileType := detectFileType(r.Fields())

		if fileType == UnknownType {
			logrus.Warnf("parse (%s): could not detect file type", path)
			return nil
		}

		// Read all the records.
		records, err := r.ReadAll()

		if err != nil || len(records) == 0 {
			logrus.Warnf("parse (%s): error reading file", path)
			return nil
		}

		switch fileType {
		case TablesFile:
			logrus.Debugf("parse (%s): adding tables file", path)
			tableList = append(tableList, records...)

		case FieldsFile:
			logrus.Debugf("parse (%s): adding fields file", path)
			var tableRecords []Attrs

			for _, record := range records {
				table = record["table"]

				if tableRecords, ok = tableFields[table]; !ok {
					tableRecords = make([]Attrs, 0)
				}

				tableRecords = append(tableRecords, record)
				tableFields[table] = tableRecords
			}

		case SchemataFile:
			for _, r := range records {
				fieldSchemata.Add(r["table"], r["field"], r)
			}

		case ReferencesFile:
			for _, r := range records {
				refs = append(refs, &Reference{
					Name:  r["name"],
					attrs: r,
				})

				schema.AddForeignKey(r)
			}

		case ConstraintsFile:
			for _, r := range records {
				switch r["type"] {
				case "primary key":
					schema.AddPrimaryKey(r)

				case "not null":
					schema.AddNotNullable(r)
				}
			}

		case IndexesFile:
			for _, r := range records {
				schema.AddIndex(r)
			}
		}

		return nil
	})

	var (
		attrs     Attrs
		t         *Table
		f         *Field
		fieldList []Attrs
		fields    FieldIndex
	)

	// Combine and link.
	model.Tables = make(TableIndex)

	// Fields that has references to other fields.
	for _, attrs = range tableList {
		fields = make(FieldIndex)

		t = &Table{
			Name:        attrs["table"],
			Description: stripNewlines(attrs["description"]),
			Label:       attrs["label"],
			Fields:      fields,
			Model:       model,
			attrs:       attrs,
		}

		model.Tables.Add(t)

		fieldList, ok = tableFields[t.Name]

		if !ok {
			continue
		}

		var req bool

		for _, attrs = range fieldList {
			switch strings.ToLower(attrs["required"]) {
			case "yes", "y", "1":
				req = true
			default:
				req = false
			}

			f = &Field{
				Name:        attrs["field"],
				Description: stripNewlines(attrs["description"]),
				Label:       attrs["label"],
				Required:    req,
				Table:       t,
				attrs:       attrs,
			}

			// Add schema information.
			if sattrs := fieldSchemata.Get(t.Name, f.Name); sattrs != nil {
				f.Type = sattrs["type"]
				f.Length = sattrs["length"]
				f.Precision = sattrs["precision"]
				f.Scale = sattrs["scale"]
				f.Default = sattrs["default"]
			}

			t.Fields.Add(f)
		}
	}

	var (
		rt *Table
		rf *Field
	)

	// Add references.
	for _, ref := range refs {
		t = model.Tables.Get(ref.attrs["table"])

		if t == nil {
			logrus.Warnf("refs (%s): no source table `%s`", model.path, ref.attrs["table"])
			continue
		}

		f = t.Fields.Get(ref.attrs["field"])

		if f == nil {
			logrus.Warnf("refs (%s:%s): no source field `%s`", model.path, t.Name, ref.attrs["field"])
			continue
		}

		rt = model.Tables.Get(ref.attrs["ref_table"])

		if rt == nil {
			logrus.Warnf("refs (%s): could not reference table `%s` by %s", model.path, ref.attrs["ref_table"], f)
			continue
		}

		rf = rt.Fields.Get(ref.attrs["ref_field"])

		if rf == nil {
			logrus.Warnf("refs (%s): could not reference field `%s` by %s", model.path, ref.attrs["ref_field"], f)
			continue
		}

		ref.Field = rf

		// Add references
		f.References = ref

		// Add back references.
		rf.InboundRefs = append(rf.InboundRefs, &Reference{
			Name:  ref.Name,
			Field: f,
		})
	}
}

// findModels walks a path and looks for models.csv files which declare a
// data model. Files in the directory will be walked to find definition files.
func findModels(root string) []*Model {
	var models []*Model

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Ignore errors.
		if err != nil {
			return nil
		}

		// Nothing to do with directories.
		if info.IsDir() {
			return nil
		}

		// Skip non-CSV files.
		if filepath.Ext(path) != ".csv" {
			return nil
		}

		f, err := os.Open(path)

		if err != nil {
			return nil
		}

		defer f.Close()

		r := NewMapCSVReader(f)

		fileType := detectFileType(r.Fields())

		if fileType == ModelsFile {
			// Read only the first line.
			attrs, err := r.Read()

			if err != nil {
				logrus.Errorf("model (%s): error reading models files", path)
				return nil
			}

			m := Model{}

			// Set the path of where the model was found.
			m.Name = attrs["model"]
			m.Version = attrs["version"]
			m.Label = attrs["label"]
			m.Description = attrs["description"]
			m.URL = attrs["url"]
			m.path = filepath.Dir(path)

			models = append(models, &m)
		}

		return nil
	})

	return models
}
