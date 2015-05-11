package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

const mappingsDir = "mappings"

var (
	dataModels    ModelIndex
	modelMappings map[string]*Mapping
)

func rebuildCache() {
	logrus.Debugf("parse: rebuilding cache")

	models := findModels()
	newDataModels := make(ModelIndex)

	var (
		err   error
		model *Model
	)

	for _, path := range models {
		if model, err = parseModel(path); err == nil {
			newDataModels.Add(model)
		}
	}

	dataModels = newDataModels

	//parseMappings()
}

func parseMappings() {
	//mapDir := filepath.Join(repoDir, "mappings")
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

	tableList := make([]map[string]string, 0)
	tableFields := make(map[string][]map[string]string)

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

		doc := records[0]

		if _, ok := doc["ref_field"]; ok {
			logrus.Debugf("parse: adding fields for %s", doc["table"])
			tableFields[doc["table"]] = records
		} else if _, ok := doc["field"]; !ok {
			logrus.Debug("parse: adding tables")
			tableList = append(tableList, records...)
		}

		return nil
	})

	var (
		ok        bool
		doc       map[string]string
		t         *Table
		f         *Field
		s         *Schema
		fieldList []map[string]string
		fields    FieldIndex
	)

	// Combine and link.
	tables := make(TableIndex)

	// Fields that has references to other fields.
	refs := make([]*Field, 0)

	for _, doc = range tableList {
		fields = make(FieldIndex)

		t = &Table{
			Name:        doc["table"],
			Description: doc["description"],
			Label:       doc["label"],
			Fields:      fields,
			attrs:       doc,
			model:       model,
		}

		tables.Add(t)
		logrus.Debugf("added table %s", t.Name)

		fieldList, ok = tableFields[t.Name]

		if !ok {
			continue
		}

		for _, doc = range fieldList {
			s = &Schema{}

			f = &Field{
				Name:        doc["field"],
				Description: doc["description"],
				Label:       doc["label"],
				RefTable:    doc["ref_table"],
				RefField:    doc["ref_field"],
				Schema:      s,
				table:       t,
				attrs:       doc,
			}

			t.Fields.Add(f)

			if doc["ref_table"] != "" {
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

		f.refTable = rt
		f.refField = rf
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

func findModels() []string {
	var (
		curPath    string
		modelPaths = make([]string, 0)
	)

	filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
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
