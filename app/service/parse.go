package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var dataModels ModelIndex

var newlinesRe = regexp.MustCompile(`[\s]+`)

func stripNewlines(s string) string {
	return newlinesRe.ReplaceAllString(s, " ")
}

func rebuildCache(path string) {
	logrus.Debugf("parse: rebuilding cache")

	models := findModels(path)
	tmpDataModels := make(ModelIndex)

	var (
		err   error
		model *Model
	)

	// Process models in parallel
	wg := sync.WaitGroup{}
	wg.Add(len(models))

	for _, path := range models {
		go func(path string) {
			if model, err = parseModel(path); err == nil {
				tmpDataModels.Add(model)
			}

			wg.Done()
		}(path)
	}

	wg.Wait()

	// Parse mapping serially since it crosses the model boundary.
	parseMappings(tmpDataModels, filepath.Join(path, "mappings"))

	dataModels = tmpDataModels
}

func parseMappings(models ModelIndex, path string) {
	logrus.Debugf("parse: parsing %s", path)

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

		for _, r := range records {
			if sm = models.Get(r["source_model"], r["source_version"]); sm == nil {
				logrus.Warnf("parse: no model %s/%s", r["source_model"], r["source_version"])
				continue
			}

			if tm = models.Get(r["target_model"], r["target_version"]); tm == nil {
				logrus.Warnf("parse: no model %s/%s", r["target_model"], r["target_version"])
				continue
			}

			if st = sm.Tables.Get(r["source_table"]); st == nil {
				logrus.Warnf("parse: no table %s/%s", sm, r["source_table"])
				continue
			}

			if tt = tm.Tables.Get(r["target_table"]); tt == nil {
				logrus.Warnf("parse: no table %s/%s", tm, r["target_table"])
				continue
			}

			if sf = st.Fields.Get(r["source_field"]); sf == nil {
				logrus.Warnf("parse: no field %s/%s", st, r["source_field"])
				continue
			}

			if tf = tt.Fields.Get(r["target_field"]); tf == nil {
				logrus.Warnf("parse: no field %s/%s", tt, r["target_field"])
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

func loadModel(fn string) (*Model, error) {
	f, err := os.Open(fn)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	m := Model{}
	d := json.NewDecoder(f)

	if err = d.Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

// parseDefinition finds and parses all definitions files in the passed directory.
func parseFiles(model *Model, path string) (TableIndex, error) {
	logrus.Debugf("parse: parsing %s", path)

	var (
		ok    bool
		table string
	)

	tableList := make([]Attrs, 0)
	tableFields := make(map[string][]Attrs)
	schemata := make(TableFieldIndex)

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

		// Read all the records.
		records, err := r.ReadAll()

		if err != nil || len(records) == 0 {
			return nil
		}

		attrs := records[0]

		table = attrs["table"]

		// Table file.
		if _, ok := attrs["field"]; !ok {
			logrus.Debug("parse: adding tables")
			tableList = append(tableList, records...)
			return nil
		}

		// Fields file.
		if _, ok := attrs["ref_field"]; ok {
			logrus.Debugf("parse: adding fields for %s", attrs["table"])
			tableFields[table] = records
			return nil
		}

		// Schema
		if _, ok := attrs["precision"]; ok {
			logrus.Debugf("parse: augmenting schema data for %s", attrs["table"])

			for _, r := range records {
				schemata.Add(table, r["field"], r)
			}

			return nil
		}

		logrus.Debugf("parse: could not detect record type in %s", path)

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
	tables := make(TableIndex)

	// Fields that has references to other fields.
	refs := make([]*Field, 0)

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

		tables.Add(t)

		logrus.Debugf("added table %s", t.Name)

		fieldList, ok = tableFields[t.Name]

		if !ok {
			continue
		}

		for _, attrs = range fieldList {
			f = &Field{
				Name:        attrs["field"],
				Description: stripNewlines(attrs["description"]),
				Label:       attrs["label"],
				Mappings:    make([]*Mapping, 0),
				Table:       t,
				attrs:       attrs,
			}

			// Add schema information.
			if sattrs := schemata.Get(t.Name, f.Name); sattrs != nil {
				f.Type = sattrs["type"]
				f.Length = sattrs["length"]
				f.Precision = sattrs["precision"]
				f.Scale = sattrs["scale"]
				f.Default = sattrs["default"]
			}

			t.Fields.Add(f)

			// Defer settings up references.
			if attrs["ref_table"] != "" {
				refs = append(refs, f)
			}
		}
	}

	for _, f = range refs {
		rt := tables.Get(f.attrs["ref_table"])

		if rt == nil {
			logrus.Warnf("parse: could not reference table %s", f.attrs["ref_table"])
			continue
		}

		rf := rt.Fields.Get(f.attrs["ref_field"])

		if rf == nil {
			logrus.Warnf("parse: could not reference field %s", f.attrs["ref_field"])
			continue
		}

		f.RefTable = rt
		f.RefField = rf
	}

	return tables, nil
}

func parseModel(path string) (*Model, error) {
	metaFile := filepath.Join(path, "datamodel.json")

	var (
		m   *Model
		err error
	)

	if m, err = loadModel(metaFile); err != nil {
		return nil, err
	}

	tables, _ := parseFiles(m, path)

	m.Tables = tables

	return m, nil
}

func findModels(path string) []string {
	var (
		curPath    string
		modelPaths = make([]string, 0)
	)

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Ignore errors.
		if err != nil {
			return nil
		}

		// Ignore files.
		if !info.IsDir() {
			return nil
		}

		// Skip hidden directories.
		if path != "." && strings.HasPrefix(filepath.Base(path), ".") {
			return filepath.SkipDir
		}

		curPath = filepath.Join(path, "datamodel.json")

		// Queue path and skip descending it.
		if pathExists(curPath) {
			logrus.Debugf("found model %s", path)
			modelPaths = append(modelPaths, path)
		}

		return nil
	})

	return modelPaths
}
