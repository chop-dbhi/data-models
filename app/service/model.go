package main

import (
	"encoding/json"
	"sort"
)

type Models []*Model

func (m Models) Len() int {
	return len(m)
}

// TODO: handle version numbers correctly.
func (m Models) Less(i, j int) bool {
	a := m[i]
	b := m[j]

	if a.Name < b.Name {
		return true
	} else if a.Name > b.Name {
		return false
	}

	return a.Version < b.Version
}

func (m Models) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type ModelIndex map[string]map[string]*Model

func (mi ModelIndex) Add(m *Model) {
	if _, ok := mi[m.Name]; !ok {
		mi[m.Name] = make(map[string]*Model)
	}

	mi[m.Name][m.Version] = m
}

func (mi ModelIndex) Get(name, version string) *Model {
	if ms, ok := mi[name]; ok {
		return ms[version]
	}

	return nil
}

func (mi ModelIndex) List() Models {
	models := make(Models, 0)

	for _, mn := range mi {
		for _, m := range mn {
			models = append(models, m)
		}
	}

	sort.Sort(models)

	return models
}

type Tables []*Table

func (ts Tables) Len() int {
	return len(ts)
}

func (ts Tables) Less(i, j int) bool {
	a := ts[i]
	b := ts[j]

	return a.Name < b.Name
}

func (ts Tables) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

type TableIndex map[string]*Table

func (ti TableIndex) Add(t *Table) {
	ti[t.Name] = t
}

func (ti TableIndex) Get(name string) *Table {
	return ti[name]
}

func (ti TableIndex) List() Tables {
	tables := make(Tables, len(ti))

	i := 0
	for _, t := range ti {
		tables[i] = t
		i++
	}

	sort.Sort(tables)

	return tables
}

// MarshalJSON implements the json.Marshaler interface. The marshaled value
// is a sorted list of tables.
func (ti TableIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(ti.List())
}

type Fields []*Field

func (fs Fields) Len() int {
	return len(fs)
}

func (fs Fields) Less(i, j int) bool {
	a := fs[i]
	b := fs[j]

	return a.Name < b.Name
}

func (fs Fields) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

type FieldIndex map[string]*Field

func (fi FieldIndex) Add(f *Field) {
	fi[f.Name] = f
}

func (fi FieldIndex) Get(name string) *Field {
	return fi[name]
}

func (fi FieldIndex) List() Fields {
	fields := make(Fields, len(fi))

	i := 0
	for _, f := range fi {
		fields[i] = f
		i++
	}

	sort.Sort(fields)

	return fields
}

// MarshalJSON implements the json.Marshaler interface. The marshaled value
// is a sorted list of fields.
func (fi FieldIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(fi.List())
}

type Model struct {
	Label   string
	Name    string
	Version string
	URL     string
	Tables  TableIndex

	path string
}

type Table struct {
	Name        string
	Label       string
	Description string
	Fields      FieldIndex

	attrs map[string]string
	model *Model
}

type Field struct {
	Name        string
	Label       string
	Description string
	RefTable    string
	RefField    string
	Schema      *Schema `json:"-"`

	attrs    map[string]string
	table    *Table
	refTable *Table
	refField *Field
}

type Schema struct {
	Type      string
	Length    int
	Precision int
	Scale     int
	Default   string

	field *Field
}

type Mapping struct {
	Source *Model
	Target *Model
}
