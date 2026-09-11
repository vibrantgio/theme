package system

import (
	"bufio"
	"fmt"
	"image/color"
	"math"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vibrantgio/theme/tokens"
)

// The catalogue the reader is pinned against is the copy tokens tests with;
// the organization's macOS reference holds the same file and the program
// that produced it.
const cataloguePath = "../tokens/testdata/nscolors.tsv"

// accentRows are the fields macOS paints with the user's accent colour. The
// catalogue was read on the platform's own blue, so they are comparable
// only on a machine whose accent is that blue — Multicolour, which is the
// absent key, or Blue.
var accentRows = map[string]bool{
	"ControlAccent":             true,
	"SelectedContentBackground": true,
	"SelectedTextBackground":    true,
	"SelectedControl":           true,
	"KeyboardFocusIndicator":    true,
}

// accentNames names an accent for a skip message.
var accentNames = map[Accent]string{
	AccentDefault:  "Multicolour",
	AccentRed:      "Red",
	AccentOrange:   "Orange",
	AccentYellow:   "Yellow",
	AccentGreen:    "Green",
	AccentBlue:     "Blue",
	AccentPurple:   "Purple",
	AccentPink:     "Pink",
	AccentGraphite: "Graphite",
}

// TestTheLivePlatformLightSetMatchesTheCatalogue pins the reader against the
// catalogue under the aqua appearance: every field carries what AppKit
// answers for its name, alpha included.
func TestTheLivePlatformLightSetMatchesTheCatalogue(t *testing.T) {
	liveSetMatchesCatalogue(t, false)
}

// TestTheLivePlatformDarkSetMatchesTheCatalogue is the same pin under
// darkAqua.
func TestTheLivePlatformDarkSetMatchesTheCatalogue(t *testing.T) {
	liveSetMatchesCatalogue(t, true)
}

// liveSetMatchesCatalogue compares one appearance's live set against the
// catalogue, field by field.
//
// The comparison is only meaningful on the system the catalogue was read
// on: a later macOS may move a value, and that is a new reading, not a
// failure — so the test skips, naming both versions. Alpha is compared to
// the precision the catalogue carries, two decimals, since the live value
// is the platform's full byte and the recorded one is that byte rounded
// through those two decimals.
func liveSetMatchesCatalogue(t *testing.T, dark bool) {
	t.Helper()
	rows, recordedOn := readCatalogue(t)
	recorded := tokens.PlatformLight
	if dark {
		recorded = tokens.PlatformDark
	}
	if running := productVersion(t); running != recordedOn {
		t.Skipf("the catalogue was read on macOS %s; this machine runs macOS %s", recordedOn, running)
	}
	accent := readAccent()
	onCatalogueAccent := accent == AccentDefault || accent == AccentBlue
	if !onCatalogueAccent {
		t.Logf("this machine's accent is %s, not the catalogue's blue; the accent rows are not compared", accentNames[accent])
	}

	set := reflect.ValueOf(liveSet(dark))
	typ := set.Type()
	for i := range typ.NumField() {
		field := typ.Field(i).Name
		got := set.Field(i).Interface().(color.NRGBA)
		if typ.Field(i).Tag.Get("appkit") == "-" {
			// A measured material: the platform gives it no NSColor
			// name, so the reader leaves the recorded value standing.
			if want := reflect.ValueOf(recorded).Field(i).Interface().(color.NRGBA); got != want {
				t.Errorf("%s = %s, want the measured %s: the live reader must not touch a field tagged `appkit:\"-\"`", field, hex(got), hex(want))
			}
			continue
		}
		if field == "FindHighlight" {
			// The one field that is not AppKit's answer: the reader does
			// not ask for findHighlightColor, so Mail's measured pair
			// stands.
			want := tokens.PlatformLight.FindHighlight
			if dark {
				want = tokens.PlatformDark.FindHighlight
			}
			if got != want {
				t.Errorf("FindHighlight = %s, want Mail's measured %s", hex(got), hex(want))
			}
			if got == rows[appKitName(field)].at(dark) {
				t.Error("FindHighlight took AppKit's findHighlightColor; it carries the find highlight as Mail paints it")
			}
			continue
		}
		if accentRows[field] && !onCatalogueAccent {
			continue
		}
		name := appKitName(field)
		row, ok := rows[name]
		if !ok {
			t.Errorf("%s: the catalogue holds no row named %q", field, name)
			continue
		}
		want := row.at(dark)
		if got.R != want.R || got.G != want.G || got.B != want.B || !sameAlpha(got.A, want.A) {
			t.Errorf("%s (%s) = %s, catalogue has %s", field, name, hex(got), hex(want))
		}
	}
}

// TestTheLiveSetIsReadOnTheAccentCadence pins that the two sets are read
// once and then served from cache until the interval the accent key is read
// on has elapsed — so a settings change arrives within a poll or two and no
// emission asks AppKit for the whole set.
func TestTheLiveSetIsReadOnTheAccentCadence(t *testing.T) {
	var clock time.Time
	reads := 0
	cache := &platformCache{
		interval: slowReadInterval,
		now:      func() time.Time { return clock },
		read: func(dark bool) tokens.PlatformColors {
			reads++
			if dark {
				return tokens.PlatformDark
			}
			return tokens.PlatformLight
		},
	}

	const polls = 60 // a one-second poll over a minute
	for i := range polls {
		light, dark := cache.sets()
		if light != tokens.PlatformLight || dark != tokens.PlatformDark {
			t.Fatalf("poll %d served a set the reader never returned", i)
		}
		clock = clock.Add(time.Second)
	}
	// One reading is both appearances: at 0 s, 10 s … 50 s, so 6 × 2.
	if reads > 12 {
		t.Errorf("AppKit was asked %d times over %d polls; want at most 12 (both appearances, once per %s)", reads, polls, slowReadInterval)
	}
	if reads < 2 {
		t.Errorf("AppKit was asked %d times; the first call must read both appearances", reads)
	}
}

// TestTheLiveSetRefreshesAfterTheInterval pins the other half of the
// cadence: a change the user makes is picked up once the interval has
// elapsed, not held forever.
func TestTheLiveSetRefreshesAfterTheInterval(t *testing.T) {
	var clock time.Time
	value := tokens.PlatformLight
	cache := &platformCache{
		interval: slowReadInterval,
		now:      func() time.Time { return clock },
		read:     func(bool) tokens.PlatformColors { return value },
	}

	if light, _ := cache.sets(); light != tokens.PlatformLight {
		t.Fatal("the first call did not read")
	}
	changed := tokens.PlatformLight
	changed.ControlAccent = tokens.PlatformLight.SystemPink
	value = changed
	clock = clock.Add(slowReadInterval - time.Second)
	if light, _ := cache.sets(); light != tokens.PlatformLight {
		t.Error("the set was re-read before the interval elapsed")
	}
	clock = clock.Add(2 * time.Second)
	if light, _ := cache.sets(); light != changed {
		t.Error("the set was not re-read after the interval elapsed")
	}
}

// measuredSection marks where the catalogue's AppKit rows stop and the
// fills read off the stored captures begin.
const measuredSection = "# measured materials"

// catalogueRow is one row of the catalogue: what AppKit answered for that
// name under each appearance.
type catalogueRow struct{ light, dark color.NRGBA }

func (r catalogueRow) at(dark bool) color.NRGBA {
	if dark {
		return r.dark
	}
	return r.light
}

// readCatalogue parses the tab-separated catalogue — "<appKitName>\t<light>\t<dark>",
// each value "#rrggbb" with an optional " a0.NN" — and returns the rows
// with the macOS version its header names.
func readCatalogue(t *testing.T) (rows map[string]catalogueRow, recordedOn string) {
	t.Helper()
	f, err := os.Open(cataloguePath)
	if err != nil {
		t.Fatalf("open catalogue: %v", err)
	}
	defer f.Close()
	rows = make(map[string]catalogueRow)
	sc := bufio.NewScanner(f)
	for line := 1; sc.Scan(); line++ {
		text := sc.Text()
		if strings.HasPrefix(text, measuredSection) {
			// The catalogue's second section: fills read off the stored
			// captures, which carry a fourth column and no AppKit name.
			// The live reader has nothing to ask for them, so the rows
			// stop here.
			break
		}
		if strings.HasPrefix(text, "#") {
			recordedOn = headerVersion(text)
			continue
		}
		if strings.TrimSpace(text) == "" {
			continue
		}
		cols := strings.Split(text, "\t")
		if len(cols) != 3 {
			t.Fatalf("catalogue line %d: want 3 tab-separated columns, got %d", line, len(cols))
		}
		light, err := parseCatalogueColor(cols[1])
		if err != nil {
			t.Fatalf("catalogue line %d, light: %v", line, err)
		}
		dark, err := parseCatalogueColor(cols[2])
		if err != nil {
			t.Fatalf("catalogue line %d, dark: %v", line, err)
		}
		rows[cols[0]] = catalogueRow{light: light, dark: dark}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("read catalogue: %v", err)
	}
	if recordedOn == "" {
		t.Fatal("the catalogue header names no macOS version")
	}
	return rows, recordedOn
}

// headerVersion lifts the product version out of the catalogue header,
// which opens "# macOS <version> (build <build>), read …".
func headerVersion(header string) string {
	_, rest, ok := strings.Cut(header, "macOS ")
	if !ok {
		return ""
	}
	return strings.Fields(rest)[0]
}

// productVersion is the macOS version this machine runs, read the way the
// catalogue's header records it.
func productVersion(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		t.Skipf("sw_vers unavailable, so the catalogue's version cannot be checked: %v", err)
	}
	return strings.TrimSpace(string(out))
}

func parseCatalogueColor(s string) (color.NRGBA, error) {
	fields := strings.Fields(s)
	if len(fields) == 0 || !strings.HasPrefix(fields[0], "#") || len(fields[0]) != 7 {
		return color.NRGBA{}, fmt.Errorf("want #rrggbb, got %q", s)
	}
	v, err := strconv.ParseUint(fields[0][1:], 16, 32)
	if err != nil {
		return color.NRGBA{}, fmt.Errorf("%q: %w", s, err)
	}
	c := color.NRGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xff}
	if len(fields) > 1 {
		if !strings.HasPrefix(fields[1], "a") {
			return color.NRGBA{}, fmt.Errorf("want an alpha of the form a0.NN, got %q", fields[1])
		}
		a, err := strconv.ParseFloat(fields[1][1:], 64)
		if err != nil {
			return color.NRGBA{}, fmt.Errorf("%q: %w", s, err)
		}
		c.A = uint8(math.Round(a * 255))
	}
	return c, nil
}

// sameAlpha compares a live alpha against a recorded one to the precision
// the catalogue carries: two decimals of coverage, which is up to half a
// 255th either side of the byte AppKit reports.
func sameAlpha(got, want uint8) bool {
	d := int(got) - int(want)
	return d >= -1 && d <= 1
}

func hex(c color.NRGBA) string {
	if c.A == 0xff {
		return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02x%02x%02x a%.2f", c.R, c.G, c.B, float64(c.A)/255)
}
