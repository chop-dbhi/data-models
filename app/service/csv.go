package main

import (
	"encoding/csv"
	"io"
)

type MapCSVReader struct {
	fields []string

	csv *csv.Reader
}

func (r *MapCSVReader) zip(keys, values []string) map[string]string {
	m := make(map[string]string, len(keys))

	for i, k := range keys {
		m[k] = values[i]
	}

	return m
}

func (r *MapCSVReader) Fields() []string {
	if r.fields == nil {
		fields, err := r.csv.Read()

		if err != nil {
			return nil
		}

		r.fields = fields
	}

	return r.fields
}

func (r *MapCSVReader) Read() (map[string]string, error) {
	// First iteration.
	if r.fields == nil {
		r.Fields()
	}

	values, err := r.csv.Read()

	if err != nil {
		return nil, err
	}

	m := r.zip(r.fields, values)

	return m, err
}

func (r *MapCSVReader) ReadAll() ([]map[string]string, error) {
	records := make([]map[string]string, 0)

	for {
		record, err := r.Read()

		if record == nil {
			break
		}

		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func NewMapCSVReader(r io.Reader) *MapCSVReader {
	cr := csv.NewReader(r)

	cr.LazyQuotes = true
	cr.TrimLeadingSpace = true

	return &MapCSVReader{
		csv: cr,
	}
}
