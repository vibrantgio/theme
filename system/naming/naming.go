// Package naming orders names the way the platform's file browser orders
// them: case-insensitively, with a run of digits read as the number it
// spells, so "note 2" stands before "note 10" and "Note" beside "note".
//
// Every name-ordered list an application draws — a folder tree, a folder
// listing, a rail of feeds — asks [Compare] or [Less] for the order, so one
// order answers all of them and none of them spells out a rule of its own.
// Where two names compare equal the caller's own tie rule decides, and a
// stable sort is what carries it: a folder before a note of the same name
// is the tree gathering its folders first, not something this package knows.
//
// It is a package of its own rather than another file beside the colour
// shim in theme/system. That shim reads the appearance the package
// publishes as a stream; an order over names is neither the appearance nor
// a stream, and the two share only the fact that macOS answers for both.
// The promise the colour shim's header makes holds here as well: no window,
// no NSApplication, any thread.
//
// On macOS the answer is the platform's own — NSString's
// localizedStandardCompare:, which is what Finder sorts by — so the order a
// list draws is the order the reader sees everywhere else on the system.
// Elsewhere [compareGo] answers, over the same Unicode collation the
// platform collates with (CLDR root, x/text/collate) plus the numeric
// reading, so the order is one order on every platform. compareGo is also
// what the macOS test compares the platform against, over a fixture of
// mixed case, digits with leading zeros and names outside ASCII.
//
// One divergence between the two is measured, and it is the only one a sweep
// of case, digits, punctuation and scripts turned up: a name spelled with the
// sharp s, against one spelled with "ss" and differing in what follows,
// orders the pair the other way round on macOS. A list holding both spellings
// therefore draws one order on macOS and the other elsewhere.
package naming

import (
	"sync"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// Compare returns a negative number where a sorts before b, a positive
// number where it sorts after, and zero where the two names order alike.
//
// Zero is not "the same name": two names that differ only in a way the
// collation ignores still order, so zero is reached only by names that are
// equal byte for byte.
func Compare(a, b string) int { return platformCompare(a, b) }

// Less reports whether a sorts before b, for a sort.Slice or a slices.SortFunc
// predicate. It must be a stable sort where the caller has a tie rule of its
// own, since [Compare] answers zero only for names that are equal.
func Less(a, b string) bool { return Compare(a, b) < 0 }

// collator holds the one comparison compareGo is built on: the root
// collation, with digit runs read as numbers, and with an order forced
// between names that collate alike but are spelled differently — the three
// options localizedStandardCompare: carries.
//
// A collate.Collator keeps a working buffer and so is not safe for
// concurrent use; the mutex is that constraint, not contention management.
var (
	collatorMu sync.Mutex
	collator   = collate.New(language.Und, collate.Numeric, collate.Force)
)

// compareGo answers the platform's order without the platform. It differs
// from the collator alone in one measured place: where two digit runs spell
// the same number, macOS puts the one with fewer leading zeros first ("note
// 2" before "note 02" before "note 002") and the collator puts it last. So
// the names are collated with their leading zeros dropped, which is the
// comparison that decides everything else, and the zeros break the tie
// afterwards.
func compareGo(a, b string) int {
	ka, za := dropLeadingZeros(a)
	kb, zb := dropLeadingZeros(b)
	collatorMu.Lock()
	c := collator.CompareString(ka, kb)
	collatorMu.Unlock()
	if c != 0 {
		return c
	}
	// Equal collation of the stripped names means the same run structure,
	// so the two zero counts line up run for run.
	for i := 0; i < len(za) && i < len(zb); i++ {
		if za[i] != zb[i] {
			if za[i] < zb[i] {
				return -1
			}
			return 1
		}
	}
	switch {
	case len(za) < len(zb):
		return -1
	case len(za) > len(zb):
		return 1
	}
	return 0
}

// dropLeadingZeros returns the name with the leading zeros of every digit
// run removed — a run of zeros keeps its last digit, so "000" reads as "0" —
// and how many zeros each run lost, in the order the runs appear.
func dropLeadingZeros(s string) (string, []int) {
	var zeros []int
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); {
		if !isDigit(s[i]) {
			out = append(out, s[i])
			i++
			continue
		}
		run := i
		for run < len(s) && isDigit(s[run]) {
			run++
		}
		start := i
		for start < run-1 && s[start] == '0' {
			start++
		}
		zeros = append(zeros, start-i)
		out = append(out, s[start:run]...)
		i = run
	}
	return string(out), zeros
}

func isDigit(b byte) bool { return b >= '0' && b <= '9' }
