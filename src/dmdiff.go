package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

const (
	Fadded int = 1 << iota
	Fremoved
	Fchanges
	Fmatches

	Fdiff = Fadded | Fremoved | Fchanges
	Fall  = Fdiff | Fmatches
)

var tableVar, fieldVar string

type Doc map[string]string

func (doc Doc) Keys() []string {
	keys := make([]string, 0)

	for k, _ := range doc {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

type Index map[string]map[string]Doc

func (idx Index) Keys(path ...string) []string {
	var (
		v  interface{}
		ok bool
	)

	// Start with the index itself.
	v = idx

	for _, k := range path {
		if v, ok = idx[k]; !ok {
			log.Fatalf("bad index path: %v", path)
		}
	}

	keys := make([]string, 0)

	switch x := v.(type) {
	case Index:
		for k, _ := range x {
			keys = append(keys, k)
		}
	case map[string]Doc:
		for k, _ := range x {
			keys = append(keys, k)
		}
	default:
		log.Fatalf("cannot get keys of type: %T", v)
	}

	sort.Strings(keys)

	return keys
}

// Add adds a document to the index.
func (idx Index) Add(doc Doc) {
	table := doc[tableVar]
	field := doc[fieldVar]

	if table == "" {
		log.Fatalf("no table defined: %v", doc)
	}

	if _, ok := idx[table]; !ok {
		idx[table] = make(map[string]Doc)
	}

	idx[table][field] = doc
}

// Stats holds counts of a Diff.
type Stats struct {
	Added   int
	Removed int
	Changes int
	Matches int
}

func (s *Stats) Total() int {
	return s.Added + s.Removed + s.Changes + s.Matches
}

func (s *Stats) Write(w io.Writer, f int) {
	if f&Fadded > 0 {
		fmt.Fprintf(w, "%d added, ", s.Added)
	}

	if f&Fremoved > 0 {
		fmt.Fprintf(w, "%d removed, ", s.Removed)
	}

	if f&Fchanges > 0 {
		fmt.Fprintf(w, "%d changed, ", s.Changes)
	}

	if f&Fmatches > 0 {
		fmt.Fprintf(w, "%d matched", s.Matches)
	}
}

// Diff holds information about the differences between two things.
type Diff struct {
	Added   []string
	Removed []string
	Matches []string
	Changes map[string][2]string
}

// Stats returns stats for this diff.
func (d *Diff) Stats() *Stats {
	s := Stats{}

	if d.Added != nil {
		s.Added = len(d.Added)
	}

	if d.Removed != nil {
		s.Removed = len(d.Removed)
	}

	if d.Changes != nil {
		s.Changes = len(d.Changes)
	}

	if d.Matches != nil {
		s.Matches = len(d.Matches)
	}

	return &s
}

func (d *Diff) Write(w io.Writer, f int) {
	if d.Matches != nil && f&Fmatches > 0 {
		for _, n := range d.Matches {
			fmt.Fprintf(w, "  %s\n", n)
		}
	}

	if d.Added != nil && f&Fadded > 0 {
		for _, n := range d.Added {
			fmt.Fprintf(w, "+ %s\n", n)
		}
	}

	if d.Removed != nil && f&Fremoved > 0 {
		for _, n := range d.Removed {
			fmt.Fprintf(w, "- %s\n", n)
		}
	}

	if d.Matches != nil && f&Fchanges > 0 {
		for k, v := range d.Changes {
			fmt.Fprintf(w, "~ %s\n    - %s\n    + %s\n", k, v[0], v[1])
		}
	}
}

// Diff compares two string slices and r
func DiffStrings(akeys, bkeys []string) *Diff {
	var (
		ad, bd bool
		ai, bi int
		ak, bk string
	)

	added := make([]string, 0)
	removed := make([]string, 0)
	matches := make([]string, 0)

	for {
		// Exhausted?
		ad = ai < len(akeys)
		bd = bi < len(bkeys)

		// Both are done.
		if !ad && !bd {
			break
		}

		// B is done.
		if !bd {
			removed = append(removed, akeys[ai:]...)
			break
		}

		// A is done.
		if !ad {
			added = append(added, bkeys[bi:]...)
			break
		}

		// Get keys to compare.
		ak = akeys[ai]
		bk = bkeys[bi]

		// File names match.
		if ak == bk {
			matches = append(matches, ak)
			ai++
			bi++

			// A is lexicographically smaller than B which means the B does not exist. Increment A.
		} else if ak < bk {
			removed = append(removed, ak)
			ai++

			// B is lexicographically smaller than A which means the A does not exist. Increment B.
		} else {
			added = append(added, bk)
			bi++
		}
	}

	return &Diff{
		Added:   added,
		Removed: removed,
		Matches: matches,
	}
}

var nonAlphaNumRe = regexp.MustCompile(`[^a-zA-Z0-9]`)

func normalize(s string) string {
	s = strings.ToLower(s)
	s = nonAlphaNumRe.ReplaceAllString(s, " ")
	return strings.Trim(s, " ")
}

// Diff compares two string slices and r
func DiffDocs(adoc, bdoc Doc) *Diff {
	var (
		ad, bd bool
		ai, bi int
		ak, bk string
		av, bv string
	)

	added := make([]string, 0)
	removed := make([]string, 0)
	matches := make([]string, 0)
	changes := make(map[string][2]string)

	akeys := adoc.Keys()
	bkeys := bdoc.Keys()

	for {
		// Exhausted?
		ad = ai < len(akeys)
		bd = bi < len(bkeys)

		// Both are done.
		if !ad && !bd {
			break
		}

		// B is done.
		if !bd {
			removed = append(removed, akeys[ai:]...)
			break
		}

		// A is done.
		if !ad {
			added = append(added, bkeys[bi:]...)
			break
		}

		// Get keys to compare.
		ak = akeys[ai]
		bk = bkeys[bi]

		// Keys match.
		if ak == bk {
			// Normalize the text a bit.
			av = normalize(adoc[ak])
			bv = normalize(bdoc[bk])

			if av == bv {
				matches = append(matches, ak)
			} else {
				changes[ak] = [2]string{adoc[ak], bdoc[bk]}
			}

			ai++
			bi++
			// A is lexicographically smaller than B which means the B does not exist. Increment A.
		} else if ak < bk {
			removed = append(removed, ak)
			ai++

			// B is lexicographically smaller than A which means the A does not exist. Increment B.
		} else {
			added = append(added, bk)
			bi++
		}
	}

	return &Diff{
		Added:   added,
		Removed: removed,
		Matches: matches,
		Changes: changes,
	}
}

// LoadFile reads and parses the file into an index of documents.
func LoadFile(path string) map[string]Doc {
	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	var doc Doc

	header := rows[0]
	docs := make(map[string]Doc)

	for _, row := range rows[1:] {
		doc = make(Doc)

		for i, f := range header {
			doc[f] = row[i]
		}

		docs[doc[fieldVar]] = doc
	}

	return docs
}

// Filenames returns a slice of filenames that will be loaded.
func Filenames(dir string) []string {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	fns := make([]string, 0)

	for _, fi := range files {
		name := fi.Name()

		// Ignore the tables file.
		if name == "tables.csv" {
			continue
		}

		if !fi.IsDir() && path.Ext(name) == ".csv" {
			fns = append(fns, name)
		}
	}

	return fns
}

/*
type Path []string

func (p *Path) String() string {
	return fmt.Sprintf("%s", *p)
}

func (p *Path) Set(v string) error {
	*p = append(*p, v)
	return nil
}
*/

func main() {
	/*
		path := Path{}

		flag.Var(&id, "id", "One or more order columns that denote the path of the hierarchy in the data model, e.g. Table, Field. The full path is expected to uniquely identify a record.")

		flag.Parse()
	*/

	flag.StringVar(&tableVar, "table", "", "Name of the table variable.")
	flag.StringVar(&fieldVar, "field", "", "Name of the field variable.")

	flag.Parse()

	if tableVar == "" || fieldVar == "" {
		flag.Usage()
		os.Exit(0)
	}

	// New document to hold the output.
	buff := &bytes.Buffer{}

	args := flag.Args()

	// a and b directory diff the version.
	adir := args[0]
	bdir := args[1]

	afiles := Filenames(adir)
	bfiles := Filenames(bdir)

	aindex := make(Index)
	bindex := make(Index)

	// Index files in A
	for _, fn := range afiles {
		docs := LoadFile(path.Join(adir, fn))

		for _, doc := range docs {
			aindex.Add(doc)
		}
	}

	// Index files in B
	for _, fn := range bfiles {
		docs := LoadFile(path.Join(bdir, fn))

		for _, doc := range docs {
			bindex.Add(doc)
		}
	}

	fieldStats := Stats{}

	// Compare tables in A and B
	atables := aindex.Keys()
	btables := bindex.Keys()

	tableDiff := DiffStrings(atables, btables)

	fmt.Fprintln(buff, "# Tables")

	fmt.Fprintln(buff, "\n```")
	tableDiff.Write(buff, Fdiff)
	fmt.Fprintln(buff, "```")

	// Fields of new tables.
	for _, k := range tableDiff.Added {
		fieldStats.Added += len(bindex[k])
	}

	// Fields of removed tables.
	for _, k := range tableDiff.Removed {
		fieldStats.Removed += len(aindex[k])
	}

	// Matches are recursed
	var (
		afields, bfields []string
		adoc, bdoc       Doc
	)

	fmt.Fprintln(buff, "\n# Fields\n")

	var diff *Diff

	for _, k := range tableDiff.Matches {
		afields = aindex.Keys(k)
		bfields = bindex.Keys(k)

		fmt.Fprintln(buff)
		fmt.Fprintf(buff, "## %s\n", k)

		diff = DiffStrings(afields, bfields)

		fmt.Fprintln(buff, "**All fields**\n")

		diff.Stats().Write(buff, Fall)

		fmt.Fprintln(buff, "\n```")
		diff.Write(buff, Fdiff)
		fmt.Fprintln(buff, "```")

		fieldStats.Added += len(diff.Added)
		fieldStats.Removed += len(diff.Removed)

		// Diff the matched fields.
		if len(diff.Matches) > 0 {
			for _, f := range diff.Matches {
				adoc = aindex[k][f]
				bdoc = bindex[k][f]

				diff = DiffDocs(adoc, bdoc)

				if len(diff.Added) > 0 || len(diff.Removed) > 0 || len(diff.Changes) > 0 {
					fmt.Fprintln(buff)
					fmt.Fprintf(buff, "### `%s`\n", f)

					diff.Stats().Write(buff, Fall)

					fmt.Fprintln(buff, "\n```")
					diff.Write(buff, Fdiff)
					fmt.Fprintln(buff, "```")

					if len(diff.Changes) > 0 {
						fieldStats.Changes += len(diff.Changes)
					}
				} else {
					fieldStats.Matches += 1
				}
			}
		}
	}

	// Print to stdout
	fmt.Println("# Overview\n")

	fmt.Printf("- Tables: ")
	tableDiff.Stats().Write(os.Stdout, Fall)

	fmt.Printf("\n- Fields: ")
	fieldStats.Write(os.Stdout, Fdiff)

	fmt.Println("\n")

	io.Copy(os.Stdout, buff)
}
