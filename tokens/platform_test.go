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

// catalogueEntry is one row: the two appearances AppKit reported. A row of
// the measured-materials section carries its provenance as well, which the
// AppKit rows have no column for.
type catalogueEntry struct {
	light, dark color.NRGBA
	provenance  string
}

// measuredSection marks where the catalogue's AppKit rows stop and the
// measured materials begin.
const measuredSection = "# measured materials"

// readCatalogue parses the tab-separated catalogue into its two sections.
// The first is what AppKit answered: a leading "#" header naming the OS it
// was read on, then "<appKitName>\t<light>\t<dark>", each value "#rrggbb"
// with an optional " a0.NNN" alpha. The [measuredSection] comment opens the
// second, whose rows carry a fourth column naming the capture the value was
// read from or the reason it is published rather than measured. Alpha
// becomes round(a*255), which is the quantization PlatformColors records.
func readCatalogue(t *testing.T) (appKit, measured map[string]catalogueEntry) {
	t.Helper()
	f, err := os.Open(cataloguePath)
	if err != nil {
		t.Fatalf("open catalogue: %v", err)
	}
	defer f.Close()
	appKit = make(map[string]catalogueEntry)
	measured = make(map[string]catalogueEntry)
	sc := bufio.NewScanner(f)
	header := false
	inMeasured := false
	for line := 1; sc.Scan(); line++ {
		text := sc.Text()
		if strings.HasPrefix(text, measuredSection) {
			inMeasured = true
			continue
		}
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
		want := 3
		if inMeasured {
			want = 4
		}
		if len(cols) != want {
			t.Fatalf("catalogue line %d: want %d tab-separated columns, got %d", line, want, len(cols))
		}
		light, err := parseCatalogueColor(cols[1])
		if err != nil {
			t.Fatalf("catalogue line %d, light: %v", line, err)
		}
		dark, err := parseCatalogueColor(cols[2])
		if err != nil {
			t.Fatalf("catalogue line %d, dark: %v", line, err)
		}
		if inMeasured {
			measured[cols[0]] = catalogueEntry{light: light, dark: dark, provenance: cols[3]}
			continue
		}
		appKit[cols[0]] = catalogueEntry{light: light, dark: dark}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("read catalogue: %v", err)
	}
	if !header {
		t.Error("the catalogue carries no header naming the macOS version it was read on")
	}
	if !inMeasured {
		t.Errorf("the catalogue carries no %q section", measuredSection)
	}
	return appKit, measured
}

// measuredName is the naming rule for the measured materials, which have no
// AppKit name to invert: the field's own name with its first letter
// lowered, so SidebarMaterial is sidebarMaterial.
func measuredName(field string) string {
	return strings.ToLower(field[:1]) + field[1:]
}

// noAppKitName reports whether a field carries the `appkit:"-"` tag — the
// rule that says the platform gives this fill no NSColor name, so the live
// reader must not ask for one and the catalogue pins it from its measured
// materials instead.
func noAppKitName(f reflect.StructField) bool {
	return f.Tag.Get("appkit") == "-"
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
	rows, _ := readCatalogue(t)
	light := reflect.ValueOf(tokens.PlatformLight)
	dark := reflect.ValueOf(tokens.PlatformDark)
	typ := light.Type()

	seen := make(map[string]bool, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		if noAppKitName(typ.Field(i)) {
			continue
		}
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

// TestMeasuredMaterialsMatchTheCatalogue pins the five fills the platform
// gives no NSColor name — the chrome material, the card's fill, the hover
// and press overlays, the floating shadow — against the catalogue's
// measured-materials section in both schemes, and holds the two sections
// disjoint: a field tagged `appkit:"-"` has a measured row and no AppKit
// one, and every measured row has a field.
func TestMeasuredMaterialsMatchTheCatalogue(t *testing.T) {
	appKit, measured := readCatalogue(t)
	light := reflect.ValueOf(tokens.PlatformLight)
	dark := reflect.ValueOf(tokens.PlatformDark)
	typ := light.Type()

	seen := make(map[string]bool, len(measured))
	for i := 0; i < typ.NumField(); i++ {
		if !noAppKitName(typ.Field(i)) {
			continue
		}
		field := typ.Field(i).Name
		name := measuredName(field)
		seen[name] = true
		row, ok := measured[name]
		if !ok {
			t.Errorf("%s: the catalogue's measured materials hold no row named %q", field, name)
			continue
		}
		if row.provenance == "" {
			t.Errorf("%s: the row %q names no provenance", field, name)
		}
		if _, clash := appKit[appKitName(field)]; clash {
			t.Errorf("%s carries `appkit:\"-\"` but the catalogue answers for %q; it belongs in the AppKit rows", field, appKitName(field))
		}
		if got := light.Field(i).Interface().(color.NRGBA); got != row.light {
			t.Errorf("%s light = %v, catalogue %s = %v", field, got, name, row.light)
		}
		if got := dark.Field(i).Interface().(color.NRGBA); got != row.dark {
			t.Errorf("%s dark = %v, catalogue %s = %v", field, got, name, row.dark)
		}
	}
	for name := range measured {
		if !seen[name] {
			t.Errorf("the measured row %q has no field on PlatformColors carrying `appkit:\"-\"`", name)
		}
	}
	if len(seen) == 0 {
		t.Error("no field carries `appkit:\"-\"`; the measured materials have lost their rule")
	}
}

// TestCardFillStandsInForTheContentsFill pins the stand-in the card's fill
// is until the System Settings grouped-box capture lands: it is the
// content's fill exactly, in both schemes, and nothing has quietly invented
// a number for it. When that capture lands and CardFill is read off its
// pixels, this test goes with the stand-in.
func TestCardFillStandsInForTheContentsFill(t *testing.T) {
	if tokens.PlatformLight.CardFill != tokens.PlatformLight.ControlBackground {
		t.Errorf("light CardFill = %v, want the content's fill %v", tokens.PlatformLight.CardFill, tokens.PlatformLight.ControlBackground)
	}
	if tokens.PlatformDark.CardFill != tokens.PlatformDark.ControlBackground {
		t.Errorf("dark CardFill = %v, want the content's fill %v", tokens.PlatformDark.CardFill, tokens.PlatformDark.ControlBackground)
	}
}

// TestTheStateOverlaysAreBlackOnLightAndWhiteOnDark pins the shape of the
// two overlays rather than their coverage, which is published and not
// measured: each is the scheme's extreme at an alpha below 1, so it
// composites over whatever fill a control carries, and the press is the
// heavier of the two.
func TestTheStateOverlaysAreBlackOnLightAndWhiteOnDark(t *testing.T) {
	for _, set := range []struct {
		name  string
		in    tokens.PlatformColors
		level uint8
	}{{"light", tokens.PlatformLight, 0x00}, {"dark", tokens.PlatformDark, 0xff}} {
		for _, o := range []struct {
			field string
			c     color.NRGBA
		}{{"HoverOverlay", set.in.HoverOverlay}, {"PressOverlay", set.in.PressOverlay}} {
			if o.c.R != set.level || o.c.G != set.level || o.c.B != set.level {
				t.Errorf("%s %s = %v, want %#02x on every channel", set.name, o.field, o.c, set.level)
			}
			if o.c.A == 0 || o.c.A == 0xff {
				t.Errorf("%s %s alpha = %d; an overlay composites, so it is neither absent nor opaque", set.name, o.field, o.c.A)
			}
		}
		if set.in.PressOverlay.A <= set.in.HoverOverlay.A {
			t.Errorf("%s PressOverlay covers %d and HoverOverlay %d; a press is the heavier of the two", set.name, set.in.PressOverlay.A, set.in.HoverOverlay.A)
		}
	}
}
