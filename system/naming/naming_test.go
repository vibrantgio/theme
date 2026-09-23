package naming

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

// fixture is the set both implementations are measured over: mixed case,
// digit runs with and without leading zeros, and names outside ASCII —
// accents, the two Japanese scripts. It is the same list on every platform,
// and on macOS the platform's own order is compared against the Go one over
// every pair of it (naming_darwin_test.go).
var fixture = []string{
	"note 10", "note 2", "Note 1", "note 1", "note 02", "note 002",
	"NOTE", "Note", "note",
	"Épée", "Epee", "über", "uber", "Uber", "Ärger", "Arger",
	"日本", "にほん", "漢字",
	"IMG_10.png", "IMG_2.png", "IMG_0002.png",
	"v1.10.0", "v1.9.0", "v1.2",
	"10 apples", "2 apples", "apples 2", "apple",
	"file-1", "file 1", "file_1", "file1",
	"", " x", "_x", "x", "Zebra", "9", "10", "1a",
}

// wantOrder is the fixture sorted, measured off macOS. It is the order every
// platform answers with, so a name-ordered list draws the same rows wherever
// the application runs.
var wantOrder = []string{
	"",
	" x",
	"_x",
	"1a",
	"2 apples",
	"9",
	"10",
	"10 apples",
	"apple",
	"apples 2",
	"Arger",
	"Ärger",
	"Epee",
	"Épée",
	"file 1",
	"file_1",
	"file-1",
	"file1",
	"IMG_2.png",
	"IMG_0002.png",
	"IMG_10.png",
	"note",
	"Note",
	"NOTE",
	"note 1",
	"Note 1",
	"note 2",
	"note 02",
	"note 002",
	"note 10",
	"uber",
	"Uber",
	"über",
	"v1.2",
	"v1.9.0",
	"v1.10.0",
	"x",
	"Zebra",
	"にほん",
	"日本",
	"漢字",
}

func TestFixtureSortsAsMeasured(t *testing.T) {
	// Both comparators are pinned to the one order, so the rows a list
	// draws do not depend on which platform answered.
	for _, c := range []struct {
		name    string
		compare func(a, b string) int
	}{
		{"the platform", Compare},
		{"the Go order", compareGo},
	} {
		got := append([]string(nil), fixture...)
		sort.SliceStable(got, func(i, j int) bool { return c.compare(got[i], got[j]) < 0 })
		if len(got) != len(wantOrder) {
			t.Fatalf("fixture has %d names, wantOrder %d", len(got), len(wantOrder))
		}
		for i := range got {
			if got[i] != wantOrder[i] {
				t.Fatalf("%s: row %d is %q, want %q\ngot:  %s\nwant: %s",
					c.name, i, got[i], wantOrder[i], strings.Join(got, " | "), strings.Join(wantOrder, " | "))
			}
		}
	}
}

func TestDigitsReadAsNumbers(t *testing.T) {
	for _, c := range [][2]string{
		{"note 2", "note 10"},
		{"9", "10"},
		{"v1.9.0", "v1.10.0"},
		{"IMG_2.png", "IMG_10.png"},
		{"2 apples", "10 apples"},
	} {
		if !Less(c[0], c[1]) {
			t.Errorf("%q does not sort before %q", c[0], c[1])
		}
		if Less(c[1], c[0]) {
			t.Errorf("%q sorts before %q as well", c[1], c[0])
		}
	}
}

func TestCaseDoesNotPartNames(t *testing.T) {
	got := []string{"note 3", "NOTE", "note 1", "Note", "note", "Note 2"}
	sort.SliceStable(got, func(i, j int) bool { return Less(got[i], got[j]) })
	want := []string{"note", "Note", "NOTE", "note 1", "Note 2", "note 3"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("case-insensitive order is %v, want %v", got, want)
	}
}

func TestLeadingZerosBreakTheTieAfterTheNumber(t *testing.T) {
	got := []string{"note 002", "note 10", "note 02", "note 2"}
	sort.SliceStable(got, func(i, j int) bool { return Less(got[i], got[j]) })
	want := []string{"note 2", "note 02", "note 002", "note 10"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("leading-zero order is %v, want %v", got, want)
	}
}

// A caller with a tie rule of its own needs Compare to answer zero for the
// names it means to tie on and nothing else, since a stable sort is what
// carries the rule.
func TestCompareIsZeroOnlyForEqualNames(t *testing.T) {
	for i, a := range fixture {
		for j, b := range fixture {
			if got := Compare(a, b); (got == 0) != (a == b) {
				t.Errorf("Compare(%q, %q) = %d; zero only where the names are equal (%d, %d)", a, b, got, i, j)
			}
		}
	}
}

func TestDropLeadingZeros(t *testing.T) {
	for _, c := range []struct {
		in    string
		want  string
		zeros []int
	}{
		{"note 002", "note 2", []int{2}},
		{"note 10", "note 10", []int{0}},
		{"000", "0", []int{2}},
		{"a0b00c", "a0b0c", []int{0, 1}},
		{"plain", "plain", nil},
	} {
		got, zeros := dropLeadingZeros(c.in)
		if got != c.want || fmt.Sprint(zeros) != fmt.Sprint(c.zeros) {
			t.Errorf("dropLeadingZeros(%q) = %q, %v; want %q, %v", c.in, got, zeros, c.want, c.zeros)
		}
	}
}
