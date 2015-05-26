package main

import (
	"bytes"
	"fmt"
	"io"
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

func diffKeys(attrs Attrs) []string {
	keys := make([]string, 0)

	for k, _ := range attrs {
		switch k {
		case "model", "version":
			continue
		default:
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)

	return keys
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

// DiffStrings compares two sorted slices of strings and returns the diff.
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

// Normalize a string for comparison to ignore insignificant changes.
func normalize(s string) string {
	s = strings.ToLower(s)
	s = nonAlphaNumRe.ReplaceAllString(s, " ")
	return strings.Trim(s, " ")
}

// DiffAttrs compare a pair of attributes.
func DiffAttrs(adoc, bdoc Attrs) *Diff {
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

	akeys := diffKeys(adoc)
	bkeys := diffKeys(bdoc)

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

func DiffModels(out io.Writer, a, b *Model) {
	buff := bytes.NewBuffer(nil)

	fieldStats := Stats{}

	atables := a.Tables
	btables := b.Tables

	tableDiff := DiffStrings(atables.Names(), btables.Names())

	fmt.Fprintln(buff, "# Tables")

	fmt.Fprintln(buff, "\n```")
	tableDiff.Write(buff, Fdiff)
	fmt.Fprintln(buff, "```")

	// Fields of new tables.
	for _, k := range tableDiff.Added {
		fieldStats.Added += len(btables.Get(k).attrs)
	}

	// Fields of removed tables.
	for _, k := range tableDiff.Removed {
		fieldStats.Removed += len(atables.Get(k).attrs)
	}

	// Matches are recursed
	var (
		afields, bfields FieldIndex
		adoc, bdoc       Attrs
	)

	fmt.Fprintln(buff, "\n# Fields\n")

	var diff *Diff

	for _, k := range tableDiff.Matches {
		afields = a.Tables.Get(k).Fields
		bfields = b.Tables.Get(k).Fields

		fmt.Fprintln(buff)
		fmt.Fprintf(buff, "## %s\n", k)

		diff = DiffStrings(afields.Names(), bfields.Names())

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
				adoc = afields.Get(f).attrs
				bdoc = bfields.Get(f).attrs

				diff = DiffAttrs(adoc, bdoc)

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
	fmt.Fprintln(out, fmt.Sprintf("# %s &rarr; %s\n", a.Label, b.Label))

	fmt.Fprintf(out, "- Tables: ")
	tableDiff.Stats().Write(out, Fall)

	fmt.Fprintf(out, "\n- Fields: ")
	fieldStats.Write(out, Fdiff)

	fmt.Fprintln(out, "\n")

	io.Copy(out, buff)
}
