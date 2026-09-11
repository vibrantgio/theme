package tokens_test

import (
	"bufio"
	"fmt"
	"image/color"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/vibrantgio/theme/tokens"
)

// The catalogue lives in the organization's macOS reference; testdata holds
// a copy so this module tests on its own, without a sibling checkout. The
// two files are byte-identical, header included.
const cataloguePath = "testdata/nscolors.tsv"

// mailFindHighlight is the find highlight as Mail paints it, measured off
// the stored captures of its find bar in both appearances. The catalogue's
// findHighlightColor row reports what AppKit answers for that name and is
// not what Mail puts on the page, so PlatformColors.FindHighlight carries
// the measurement and this test holds the two apart.
var mailFindHighlight = struct{ light, dark color.NRGBA }{
	light: color.NRGBA{R: 0xfa, G: 0xef, B: 0xbd, A: 0xff},
	dark:  color.NRGBA{R: 0x6e, G: 0x6e, B: 0x4d, A: 0xff},
}

// catalogueEntry is one row: the two appearances AppKit reported.
type catalogueEntry struct{ light, dark color.NRGBA }

// readCatalogue parses the tab-separated catalogue: a leading "#" header
// naming the OS it was read on, then "<appKitName>\t<light>\t<dark>", each
// value "#rrggbb" with an optional " a0.NN" alpha. Alpha becomes
// round(a*255), which is the quantization PlatformColors records.
func readCatalogue(t *testing.T) map[string]catalogueEntry {
	t.Helper()
	f, err := os.Open(cataloguePath)
	if err != nil {
		t.Fatalf("open catalogue: %v", err)
	}
	defer f.Close()
	rows := make(map[string]catalogueEntry)
	sc := bufio.NewScanner(f)
	header := false
	for line := 1; sc.Scan(); line++ {
		text := sc.Text()
		if strings.HasPrefix(text, "#") {
			header = true
			if !strings.Contains(text, "macOS") {
				t.Errorf("catalogue line %d: the header must name the macOS version it was read on, got %q", line, text)
			}
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
		rows[cols[0]] = catalogueEntry{light: light, dark: dark}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("read catalogue: %v", err)
	}
	if !header {
		t.Error("the catalogue carries no header naming the macOS version it was read on")
	}
	return rows
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

// appKitName is the inverse of the field-naming rule: AppKit's trailing
// "Color" is dropped and nothing else changes, so WindowBackground is
// windowBackgroundColor, and the system colours, which carry no such
// suffix, are systemRed and its fellows.
func appKitName(field string) string {
	name := strings.ToLower(field[:1]) + field[1:]
	if strings.HasPrefix(field, "System") {
		return name
	}
	return name + "Color"
}

// TestPlatformColorsMatchCatalogue pins both recorded sets against the
// catalogue read off AppKit, field by field and row by row: every field
// carries its name's recorded value, every recorded row has a field, and
// FindHighlight carries Mail's measurement rather than the AppKit answer.
func TestPlatformColorsMatchCatalogue(t *testing.T) {
	rows := readCatalogue(t)
	light := reflect.ValueOf(tokens.PlatformLight)
	dark := reflect.ValueOf(tokens.PlatformDark)
	typ := light.Type()

	seen := make(map[string]bool, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i).Name
		name := appKitName(field)
		seen[name] = true
		row, ok := rows[name]
		if !ok {
			t.Errorf("%s: the catalogue holds no row named %q", field, name)
			continue
		}
		gotLight := light.Field(i).Interface().(color.NRGBA)
		gotDark := dark.Field(i).Interface().(color.NRGBA)
		if field == "FindHighlight" {
			if gotLight != mailFindHighlight.light || gotDark != mailFindHighlight.dark {
				t.Errorf("FindHighlight = %v / %v, want Mail's measured %v / %v",
					gotLight, gotDark, mailFindHighlight.light, mailFindHighlight.dark)
			}
			if gotLight == row.light || gotDark == row.dark {
				t.Error("FindHighlight has taken the AppKit findHighlightColor answer; it carries Mail's measurement instead")
			}
			continue
		}
		if gotLight != row.light {
			t.Errorf("%s light = %v, catalogue %s = %v", field, gotLight, name, row.light)
		}
		if gotDark != row.dark {
			t.Errorf("%s dark = %v, catalogue %s = %v", field, gotDark, name, row.dark)
		}
	}
	for name := range rows {
		if !seen[name] {
			t.Errorf("the catalogue row %q has no field on PlatformColors", name)
		}
	}
}

// accentRows are the rows WithAccent rebuilds; every other field stands.
var accentRows = []string{
	"ControlAccent",
	"SelectedContentBackground",
	"SelectedTextBackground",
	"SelectedControl",
	"KeyboardFocusIndicator",
}

// TestWithAccentKeepsThePlatformBlue pins that the recorded sets already
// carry the platform's own accent: asking for it changes nothing, so a
// machine on the default accent gets the catalogue exactly.
func TestWithAccentKeepsThePlatformBlue(t *testing.T) {
	blue := tokens.PlatformLight.ControlAccent
	if got := tokens.PlatformLight.WithAccent(blue); got != tokens.PlatformLight {
		t.Error("the light set moved when asked for the accent it already carries")
	}
	if got := tokens.PlatformDark.WithAccent(blue); got != tokens.PlatformDark {
		t.Error("the dark set moved when asked for the accent it already carries")
	}
}

// TestWithAccentMovesOnlyTheAccentRows pins the rule's reach: the five rows
// the platform derives from the accent take the new colour's hue, and every
// other field — the inactive window's grey selection, the link, the system
// colours, the planes — is untouched.
func TestWithAccentMovesOnlyTheAccentRows(t *testing.T) {
	follows := make(map[string]bool, len(accentRows))
	for _, name := range accentRows {
		follows[name] = true
	}
	purple := color.NRGBA{R: 0xcb, G: 0x30, B: 0xe0, A: 0xff}
	for _, set := range []struct {
		name string
		in   tokens.PlatformColors
	}{{"light", tokens.PlatformLight}, {"dark", tokens.PlatformDark}} {
		got := reflect.ValueOf(set.in.WithAccent(purple))
		want := reflect.ValueOf(set.in)
		typ := want.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i).Name
			g := got.Field(i).Interface().(color.NRGBA)
			w := want.Field(i).Interface().(color.NRGBA)
			if !follows[field] {
				if g != w {
					t.Errorf("%s %s = %v, want the recorded %v: it does not follow the accent", set.name, field, g, w)
				}
				continue
			}
			if g == w {
				t.Errorf("%s %s did not follow the accent", set.name, field)
			}
			if g.A != w.A {
				t.Errorf("%s %s alpha = %d, want the recorded %d", set.name, field, g.A, w.A)
			}
		}
		if c := set.in.WithAccent(purple).ControlAccent; c.R != purple.R || c.G != purple.G || c.B != purple.B {
			t.Errorf("%s ControlAccent = %v, want the accent %v", set.name, c, purple)
		}
	}
}

// TestWithAccentDrainsAGreyAccent pins that a grey accent — the platform's
// graphite — leaves the following rows grey rather than blue: the rule
// carries the accent's saturation, not only its hue.
func TestWithAccentDrainsAGreyAccent(t *testing.T) {
	graphite := tokens.PlatformLight.SystemGray
	got := tokens.PlatformLight.WithAccent(graphite)
	for _, c := range []color.NRGBA{got.SelectedContentBackground, got.SelectedTextBackground, got.SelectedControl, got.KeyboardFocusIndicator} {
		span := int(max3(c.R, c.G, c.B)) - int(min3(c.R, c.G, c.B))
		if span > 8 {
			t.Errorf("%v spans %d across its channels; a graphite accent leaves these rows grey", c, span)
		}
	}
}

func max3(a, b, c uint8) uint8 {
	if b > a {
		a = b
	}
	if c > a {
		a = c
	}
	return a
}

func min3(a, b, c uint8) uint8 {
	if b < a {
		a = b
	}
	if c < a {
		a = c
	}
	return a
}
