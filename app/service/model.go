package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Attrs map[string]string

type TableFieldIndex map[string]map[string]Attrs

func (i TableFieldIndex) Add(t, f string, a Attrs) {
	t = strings.ToLower(t)
	f = strings.ToLower(f)

	if _, ok := i[t]; !ok {
		i[t] = make(map[string]Attrs)
	}

	i[t][f] = a
}

func (i TableFieldIndex) Get(t, f string) Attrs {
	t = strings.ToLower(t)
	f = strings.ToLower(f)

	if _, ok := i[t]; !ok {
		return nil
	}

	return i[t][f]
}

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

func (mi ModelIndex) Keys() []string {
	keys := make([]string, len(mi))

	i := 0

	for k, _ := range mi {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	return keys
}

func (mi ModelIndex) Add(m *Model) {
	n := strings.ToLower(m.Name)
	v := strings.ToLower(m.Version)

	if _, ok := mi[n]; !ok {
		mi[n] = make(map[string]*Model)
	}

	mi[n][v] = m
}

func (mi ModelIndex) Get(n, v string) *Model {
	n = strings.ToLower(n)
	v = strings.ToLower(v)

	if ms, ok := mi[n]; ok {
		return ms[v]
	}

	return nil
}

func (mi ModelIndex) Versions(n string) Models {
	n = strings.ToLower(n)

	if _, ok := mi[n]; !ok {
		return nil
	}

	models := make(Models, 0)

	for _, m := range mi[n] {
		models = append(models, m)
	}

	sort.Sort(models)

	return models
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
	ti[strings.ToLower(t.Name)] = t
}

func (ti TableIndex) Get(n string) *Table {
	return ti[strings.ToLower(n)]
}

func (ti TableIndex) Names() []string {
	names := make([]string, len(ti))

	i := 0
	for s, _ := range ti {
		names[i] = s
		i++
	}

	sort.Strings(names)

	return names
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
	fi[strings.ToLower(f.Name)] = f
}

func (fi FieldIndex) Get(n string) *Field {
	return fi[strings.ToLower(n)]
}

func (ti FieldIndex) Names() []string {
	names := make([]string, len(ti))

	i := 0
	for s, _ := range ti {
		names[i] = s
		i++
	}

	sort.Strings(names)

	return names
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
	Label       string
	Name        string
	Version     string
	Description string
	URL         string
	Tables      TableIndex

	schema *Schema

	path string
}

func (m *Model) MarshalJSON() ([]byte, error) {
	aux := map[string]interface{}{
		"name":        m.Name,
		"version":     m.Version,
		"description": m.Description,
		"url":         m.URL,
		"tables":      m.Tables,
	}

	return json.Marshal(aux)
}

func (m *Model) String() string {
	return fmt.Sprintf("%s/%s", m.Name, m.Version)
}

type Table struct {
	Name        string
	Label       string
	Description string
	Fields      FieldIndex

	Model *Model

	attrs Attrs
}

func (t *Table) MarshalJSON() ([]byte, error) {
	aux := map[string]interface{}{
		"model":       t.Model.Name,
		"version":     t.Model.Version,
		"name":        t.Name,
		"description": t.Description,
		"fields":      t.Fields,
	}

	return json.Marshal(aux)
}

func (t *Table) String() string {
	return fmt.Sprintf("%s/%s", t.Model, t.Name)
}

type Field struct {
	Name        string
	Label       string
	Description string
	Required    bool

	// Schema fields
	Type      string
	Length    string
	Precision string
	Scale     string
	Default   string

	Table    *Table
	Mappings []*Mapping

	// The field this field was renamed from in the previous version.
	RenamedFrom *Field

	// The field this field was renamed to in the next version.
	RenamedTo *Field

	// The field this field references.
	References *Reference

	// Fields that reference this field.
	InboundRefs []*Reference

	attrs Attrs
}

func (f *Field) String() string {
	return fmt.Sprintf("%s/%s", f.Table, f.Name)
}

func (f *Field) MarshalJSON() ([]byte, error) {
	aux := map[string]interface{}{
		"name":        f.Name,
		"table":       f.Table.Name,
		"description": f.Description,
		"required":    f.Required,
		"type":        f.Type,
		"length":      json.Number(f.Length),
		"precision":   json.Number(f.Precision),
		"scale":       json.Number(f.Scale),
		"default":     f.Default,
	}

	return json.Marshal(aux)
}

// A mapping defined a correspondence between two fields. A mapping points
// points to the opposing field and a *comment* which describes the nuances
// of the relationship between the two fields.
type Mapping struct {
	Field   *Field
	Comment string
}

// Reference declares that the source field is a reference to the target field.
type Reference struct {
	Name  string
	Field *Field

	attrs Attrs
}

func (r *Reference) MarshalJSON() ([]byte, error) {
	aux := map[string]string{
		"name":  r.Name,
		"table": r.Field.Table.Name,
		"field": r.Field.Name,
	}

	return json.Marshal(aux)
}

// Schema contains constraints and indexes for a model.
type Schema struct {
	// Schematic components.
	PrimaryKeys  map[string]*PrimaryKey
	ForeignKeys  []*ForeignKey
	NotNullables []*NotNullable
	Uniques      map[string]*Unique
	Indexes      map[string]*Index
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	pks := make([]*PrimaryKey, len(s.PrimaryKeys))
	uniqs := make([]*Unique, len(s.Uniques))
	indexes := make([]*Index, len(s.Indexes))

	i := 0

	for _, pk := range s.PrimaryKeys {
		pks[i] = pk
		i++
	}

	i = 0

	for _, un := range s.Uniques {
		uniqs[i] = un
		i++
	}

	i = 0

	for _, idx := range s.Indexes {
		indexes[i] = idx
		i++
	}

	aux := map[string]interface{}{
		"indexes": indexes,
		"constraints": map[string]interface{}{
			"foreign_keys": s.ForeignKeys,
			"primary_keys": pks,
			"uniques":      uniqs,
			"not_null":     s.NotNullables,
		},
	}

	return json.Marshal(aux)
}

func (s *Schema) AddPrimaryKey(a Attrs) {
	if s.PrimaryKeys == nil {
		s.PrimaryKeys = make(map[string]*PrimaryKey)
	}

	n := a["name"]

	if pk, ok := s.PrimaryKeys[n]; !ok {
		s.PrimaryKeys[n] = &PrimaryKey{
			Name:   n,
			Table:  a["table"],
			Fields: []string{a["field"]},
		}
	} else {
		pk.Fields = append(pk.Fields, a["field"])
	}
}

func (s *Schema) AddForeignKey(a Attrs) {
	s.ForeignKeys = append(s.ForeignKeys, &ForeignKey{
		Name:        a["name"],
		SourceTable: a["table"],
		SourceField: a["field"],
		TargetTable: a["ref_table"],
		TargetField: a["ref_field"],
	})
}

func (s *Schema) AddNotNullable(a Attrs) {
	s.NotNullables = append(s.NotNullables, &NotNullable{
		Table: a["table"],
		Field: a["field"],
	})
}

func (s *Schema) AddUnique(a Attrs) {
	if s.Uniques == nil {
		s.Uniques = make(map[string]*Unique)
	}

	n := a["name"]

	if un, ok := s.Uniques[n]; !ok {
		s.Uniques[n] = &Unique{
			Name:   n,
			Table:  a["table"],
			Fields: []string{a["field"]},
		}
	} else {
		un.Fields = append(un.Fields, a["field"])
	}
}

func (s *Schema) AddIndex(a Attrs) {
	if s.Indexes == nil {
		s.Indexes = make(map[string]*Index)
	}

	n := a["name"]

	if idx, ok := s.Indexes[n]; !ok {
		var uniq bool

		switch strings.ToLower(a["unique"]) {
		case "yes", "y", "1":
			uniq = true
		}

		s.Indexes[n] = &Index{
			Name:   n,
			Order:  a["order"],
			Table:  a["table"],
			Unique: uniq,
			Fields: []string{a["field"]},
		}
	} else {
		idx.Fields = append(idx.Fields, a["field"])
	}
}

// PrimaryKey is a constraint which declares the field values uniquely define
// a record in the respective table.
type PrimaryKey struct {
	Name   string   `json:"name"`
	Table  string   `json:"table"`
	Fields []string `json:"fields"`
}

// Unique is a constraint which declares the field values be unique for
// a record in the respective table.
type Unique struct {
	Name   string   `json:"name"`
	Table  string   `json:"table"`
	Fields []string `json:"fields"`
}

// ForeignKey is a constraint which declares the field values are constrained
// to values in the referenced fields.
type ForeignKey struct {
	Name        string `json:"name"`
	SourceTable string `json:"source_table"`
	SourceField string `json:"source_field"`
	TargetTable string `json:"target_table"`
	TargetField string `json:"target_field"`
}

// NotNullable is a constraint which declares the field cannot be a null.
type NotNullable struct {
	Table string `json:"table"`
	Field string `json:"field"`
}

// Index represents a schematic index for one or more fields.
type Index struct {
	Name   string   `json:"name"`
	Unique bool     `json:"unique"`
	Order  string   `json:"order"`
	Table  string   `json:"table"`
	Fields []string `json:"fields"`
}
