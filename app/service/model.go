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
	Label   string
	Name    string
	Version string
	URL     string
	Tables  TableIndex

	path string
}

func (m *Model) String() string {
	return fmt.Sprintf("%s/%s", m.Name, m.Version)
}

type Table struct {
	Name        string
	Label       string
	Description string
	Fields      FieldIndex

	Model *Model `json:"-"`

	attrs Attrs
}

func (t *Table) String() string {
	return fmt.Sprintf("%s/%s", t.Model, t.Name)
}

type Field struct {
	Name        string
	Label       string
	Description string

	// Schema fields
	Type      string
	Length    string
	Precision string
	Scale     string
	Default   string

	RefTable *Table `json:"-"`
	RefField *Field `json:"-"`

	Table    *Table     `json:"-"`
	Mappings []*Mapping `json:"-"`

	attrs Attrs
}

func (f *Field) String() string {
	return fmt.Sprintf("%s/%s", f.Table, f.Name)
}

// A mapping defined a correspondence between two fields. A mapping points
// points to the opposing field and a *comment* which describes the nuances
// of the relationship between the two fields.
type Mapping struct {
	Field   *Field
	Comment string
}
