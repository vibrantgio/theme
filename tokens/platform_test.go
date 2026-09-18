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

	vgcolor "github.com/vibrantgio/theme/color"
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
// AppKit rows have no column for, and a shadow's row carries the geometry
// its coverage was fitted at beside each appearance's value.
type catalogueEntry struct {
	light, dark         color.NRGBA
	lightGeom, darkGeom shadowGeometry
	provenance          string
}

// shadowGeometry is the reach and the offset a row records beside a fitted
// coverage, in dp. spelt reports whether the row spelt them at all: a fill is
// a colour and nothing more.
type shadowGeometry struct {
	reach, offset float64
	spelt         bool
}

// measuredSection marks where the catalogue's AppKit rows stop and the
// measured materials begin.
const measuredSection = "# measured materials"

// readCatalogue parses the tab-separated catalogue into its two sections.
// The first is what AppKit answered: a leading "#" header naming the OS it
// was read on, then "<appKitName>\t<light>\t<dark>", each value "#rrggbb"
// with an optional alpha: " aN/255" in the AppKit rows, which carry the
// byte AppKit reported, and " a0.NNN" in the measured rows, which carry a
// fitted coverage. The [measuredSection] comment opens the
// second, whose rows carry a fourth column naming the capture the value was
// read from or the reason it is published rather than measured.
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
			lg, err := parseCatalogueGeometry(cols[1])
			if err != nil {
				t.Fatalf("catalogue line %d, light: %v", line, err)
			}
			dg, err := parseCatalogueGeometry(cols[2])
			if err != nil {
				t.Fatalf("catalogue line %d, dark: %v", line, err)
			}
			measured[cols[0]] = catalogueEntry{light: light, dark: dark, lightGeom: lg, darkGeom: dg, provenance: cols[3]}
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
			return color.NRGBA{}, fmt.Errorf("want an alpha of the form aN/255 or a0.NN, got %q", fields[1])
		}
		// Two forms, and the difference is what the row is. An AppKit row
		// carries the byte AppKit reported, "a216/255", and nothing rounds.
		// A measured row carries a fitted coverage, "a0.572", which is a
		// fraction because a fit is not a byte anyone read off the platform.
		if num, den, ok := strings.Cut(fields[1][1:], "/"); ok {
			n, err := strconv.ParseUint(num, 10, 8)
			if err != nil {
				return color.NRGBA{}, fmt.Errorf("%q: %w", s, err)
			}
			if den != "255" {
				return color.NRGBA{}, fmt.Errorf("%q: an alpha byte is written over 255, got %q", s, den)
			}
			c.A = uint8(n)
			return c, nil
		}
		a, err := strconv.ParseFloat(fields[1][1:], 64)
		if err != nil {
			return color.NRGBA{}, fmt.Errorf("%q: %w", s, err)
		}
		c.A = uint8(math.Round(a * 255))
	}
	return c, nil
}

// parseCatalogueGeometry reads the reach and the offset a shadow row spells
// after its coverage: "#000000 a0.035 reach 23 offset 9", both in dp. A row
// that spells neither is a fill and answers a zero value.
func parseCatalogueGeometry(s string) (shadowGeometry, error) {
	fields := strings.Fields(s)
	var g shadowGeometry
	for i := 2; i < len(fields); i += 2 {
		if i+1 >= len(fields) {
			return g, fmt.Errorf("%q: %q names no number", s, fields[i])
		}
		v, err := strconv.ParseFloat(fields[i+1], 64)
		if err != nil {
			return g, fmt.Errorf("%q: %w", s, err)
		}
		switch fields[i] {
		case "reach":
			g.reach, g.spelt = v, true
		case "offset":
			g.offset, g.spelt = v, true
		default:
			return g, fmt.Errorf("%q: unknown field %q, want reach or offset", s, fields[i])
		}
	}
	return g, nil
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
	"SidebarSelection",
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
			if !follows[field] {
				// Compared as whatever the field is: a shadow carries a
				// measured geometry beside its coverage, and the accent
				// moves neither.
				if g, w := got.Field(i).Interface(), want.Field(i).Interface(); g != w {
					t.Errorf("%s %s = %v, want the recorded %v: it does not follow the accent", set.name, field, g, w)
				}
				continue
			}
			g := got.Field(i).Interface().(color.NRGBA)
			w := want.Field(i).Interface().(color.NRGBA)
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

// TestMeasuredMaterialsMatchTheCatalogue pins every fill the platform gives
// no NSColor name against the catalogue's
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
		checkMeasuredRow(t, field, name, "light", light.Field(i).Interface(), row.light, row.lightGeom)
		checkMeasuredRow(t, field, name, "dark", dark.Field(i).Interface(), row.dark, row.darkGeom)
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

// checkMeasuredRow compares one field of one appearance against its
// catalogue row. A shadow answers a coverage and the geometry it was fitted
// at, and the row spells both; every other field is a colour and its row
// spells one.
func checkMeasuredRow(t *testing.T, field, name, appearance string, got any, want color.NRGBA, geom shadowGeometry) {
	t.Helper()
	switch v := got.(type) {
	case color.NRGBA:
		if geom.spelt {
			t.Errorf("%s: the catalogue's %s row spells a shadow geometry the field does not carry", field, name)
		}
		if v != want {
			t.Errorf("%s %s = %v, catalogue %s = %v", field, appearance, v, name, want)
		}
	case tokens.DropShadow:
		if !geom.spelt {
			t.Errorf("%s: the catalogue's %s row spells no reach or offset for a shadow", field, name)
		}
		if v.Peak != want {
			t.Errorf("%s %s peak = %v, catalogue %s = %v", field, appearance, v.Peak, name, want)
		}
		if float64(v.Reach) != geom.reach || float64(v.Offset) != geom.offset {
			t.Errorf("%s %s = reach %v offset %v, catalogue %s = reach %v offset %v",
				field, appearance, v.Reach, v.Offset, name, geom.reach, geom.offset)
		}
	default:
		t.Errorf("%s: a measured material is a colour or a shadow, got %T", field, got)
	}
}

// TestCardFillIsTheMeasuredGroupedBox pins the card's fill to the pixels of
// the System Settings grouped-box captures rather than to another token: the
// platform's box is a fill of its own in both schemes, one step off the
// plane it sits on, so it is neither the content's fill nor the chrome's.
func TestCardFillIsTheMeasuredGroupedBox(t *testing.T) {
	for _, c := range []struct {
		name string
		in   tokens.PlatformColors
		want color.NRGBA
	}{
		{"light", tokens.PlatformLight, color.NRGBA{R: 0xf7, G: 0xf7, B: 0xf7, A: 0xff}},
		{"dark", tokens.PlatformDark, color.NRGBA{R: 0x2a, G: 0x30, B: 0x34, A: 0xff}},
	} {
		if c.in.CardFill != c.want {
			t.Errorf("%s CardFill = %v, want the measured grouped box %v", c.name, c.in.CardFill, c.want)
		}
		if c.in.CardFill == c.in.ControlBackground {
			t.Errorf("%s CardFill = %v, the content's fill; the box was measured apart from it", c.name, c.in.CardFill)
		}
		// The light box and the light chrome material are the same
		// value on this platform: the grouped box reads #f7f7f7 over
		// System Settings' white plane, and Finder's sidebar reads
		// #f7f7f7 too. They were measured apart, off different windows,
		// and landed together; the dark pair does not.
		if c.name == "dark" && c.in.CardFill == c.in.SidebarMaterial {
			t.Errorf("%s CardFill = %v, the chrome material; the box was measured apart from it", c.name, c.in.CardFill)
		}
	}
}

// TestPushButtonFillIsTheMeasuredPushButton pins the push button's fill to
// the pixels of the Save dialog captures rather than to controlColor, which
// reports a different colour in both appearances: the platform draws the
// ordinary push button at a fill of its own, and a consumer that wore
// controlColor would paint white on the light sheet where the platform
// paints #ececec.
func TestPushButtonFillIsTheMeasuredPushButton(t *testing.T) {
	for _, c := range []struct {
		name string
		in   tokens.PlatformColors
		want color.NRGBA
	}{
		{"light", tokens.PlatformLight, color.NRGBA{R: 0xec, G: 0xec, B: 0xec, A: 0xff}},
		{"dark", tokens.PlatformDark, color.NRGBA{R: 0x33, G: 0x3a, B: 0x3f, A: 0xff}},
	} {
		if c.in.PushButtonFill != c.want {
			t.Errorf("%s PushButtonFill = %v, want the measured push button %v", c.name, c.in.PushButtonFill, c.want)
		}
		if c.in.PushButtonFill == c.in.Control {
			t.Errorf("%s PushButtonFill = %v, controlColor's own value; the button was measured apart from it", c.name, c.in.PushButtonFill)
		}
		if c.in.PushButtonFill.A != 0xff {
			t.Errorf("%s PushButtonFill = %v; the fill was read as a pixel, so it is opaque", c.name, c.in.PushButtonFill)
		}
	}
}

// TestTheStateOverlaysAreBlackOnLightAndWhiteOnDark pins the shape of the
// two overlays rather than their coverage, which the catalogue holds: each
// is the scheme's extreme at an alpha below 1, so it composites over
// whatever fill a control carries, and the press is the heavier of the two.
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

// TestAlphaNamesFlattenToTheCapturedBytes: an alpha name in this set is
// painted by flattening it over the surface it stands on in encoded sRGB,
// which is where the platform composites. These are the composites the Save
// dialog captures in the organization's macOS reference actually show, at
// the glyph cores of the sheet's own wording and of the Cancel button's
// title.
//
// Every row lands on the captured byte exactly. It did not while the
// catalogue recorded alpha to two decimals — Label carried round(0.85×255)
// = 217 where the platform's own coverage byte is 216, and the label rows
// landed one 255th light. The catalogue now carries the byte.
func TestAlphaNamesFlattenToTheCapturedBytes(t *testing.T) {
	for _, tc := range []struct {
		name          string
		fg, surface   color.NRGBA
		want, capture color.NRGBA
	}{
		{"Label on the light sheet", tokens.PlatformLight.Label, tokens.PlatformLight.WindowBackground,
			color.NRGBA{0x27, 0x27, 0x27, 0xff}, color.NRGBA{0x27, 0x27, 0x27, 0xff}},
		{"Label on the light push button", tokens.PlatformLight.Label, tokens.PlatformLight.PushButtonFill,
			color.NRGBA{0x24, 0x24, 0x24, 0xff}, color.NRGBA{0x24, 0x24, 0x24, 0xff}},
		{"Label on the dark push button", tokens.PlatformDark.Label, tokens.PlatformDark.PushButtonFill,
			color.NRGBA{0xe0, 0xe1, 0xe2, 0xff}, color.NRGBA{0xe0, 0xe1, 0xe2, 0xff}},
		{"SecondaryLabel on the dark sheet", tokens.PlatformDark.SecondaryLabel, color.NRGBA{0x23, 0x2a, 0x2f, 0xff},
			color.NRGBA{0x9c, 0x9f, 0xa1, 0xff}, color.NRGBA{0x9c, 0x9f, 0xa1, 0xff}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := vgcolor.Flatten(tc.fg, tc.surface)
			if got != tc.want {
				t.Errorf("Flatten(%v, %v) = %v, want %v", tc.fg, tc.surface, got, tc.want)
			}
			for _, d := range []int{int(got.R) - int(tc.capture.R), int(got.G) - int(tc.capture.G), int(got.B) - int(tc.capture.B)} {
				if d < -1 || d > 1 {
					t.Errorf("Flatten(%v, %v) = %v, more than one 255th from the captured %v", tc.fg, tc.surface, got, tc.capture)
				}
			}
		})
	}
}

// TestDisabledCoverageLandsOnTheCapturedSwitchedOffFill pins the disabled
// coverage against the pixels rather than against itself. MEASURED,
// save-dialog-{light,dark}.png: the "Options:" checkbox is switched off in
// both appearances and its box reads #f2f2f2 light and #2e3439 dark, on
// sheets of #ffffff and #232a2f; the "File Format:" pop-up seventeen rows
// above it is enabled on the same sheet and reads the push button's own
// #ececec and #333a3f. Light lands to the byte and dark within one 255th,
// which is the tolerance the hover overlay's dark reading carries.
func TestDisabledCoverageLandsOnTheCapturedSwitchedOffFill(t *testing.T) {
	for _, tc := range []struct {
		name  string
		in    tokens.PlatformColors
		sheet color.NRGBA
		want  color.NRGBA
		slack int
	}{
		{"light", tokens.PlatformLight, color.NRGBA{0xff, 0xff, 0xff, 0xff}, color.NRGBA{0xf2, 0xf2, 0xf2, 0xff}, 0},
		{"dark", tokens.PlatformDark, color.NRGBA{0x23, 0x2a, 0x2f, 0xff}, color.NRGBA{0x2e, 0x34, 0x39, 0xff}, 1},
	} {
		got := vgcolor.Flatten(vgcolor.Fade(tc.in.PushButtonFill, tokens.DisabledCoverage), tc.sheet)
		off := func(a, b uint8) int {
			if a > b {
				return int(a) - int(b)
			}
			return int(b) - int(a)
		}
		if off(got.R, tc.want.R) > tc.slack || off(got.G, tc.want.G) > tc.slack || off(got.B, tc.want.B) > tc.slack {
			t.Errorf("%s: the push button's fill at the disabled coverage lands on %v, want the capture's %v within %d",
				tc.name, got, tc.want, tc.slack)
		}
	}
}

// The disabled coverage is a fade and not an erasure: a control switched off
// still draws, and it still moves off the fill it had.
func TestDisabledCoverageIsAFade(t *testing.T) {
	if tokens.DisabledCoverage == 0 || tokens.DisabledCoverage == 0xff {
		t.Errorf("DisabledCoverage = %d; a switched-off control neither disappears nor draws at full strength", tokens.DisabledCoverage)
	}
}

// TestToolbarControlFillStandsOffTheChrome pins the toolbar control's fill to
// the pixels it was read at and holds the one property the reading was taken
// for: on the chrome material this set paints, the control is lighter than
// the band it stands on in both appearances, so it reads as a figure on the
// band rather than as part of it.
func TestToolbarControlFillStandsOffTheChrome(t *testing.T) {
	for _, c := range []struct {
		name string
		in   tokens.PlatformColors
		want color.NRGBA
	}{
		{"light", tokens.PlatformLight, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}},
		{"dark", tokens.PlatformDark, color.NRGBA{R: 0x26, G: 0x26, B: 0x26, A: 0xff}},
	} {
		if c.in.ToolbarControlFill != c.want {
			t.Errorf("%s ToolbarControlFill = %v, want the measured toolbar control %v",
				c.name, c.in.ToolbarControlFill, c.want)
		}
		if c.in.ToolbarControlFill == c.in.SidebarMaterial {
			t.Errorf("%s ToolbarControlFill = %v, the chrome material itself; a control that reads as its band is not a control",
				c.name, c.in.ToolbarControlFill)
		}
		fill, band := c.in.ToolbarControlFill, c.in.SidebarMaterial
		if fill.R <= band.R || fill.G <= band.G || fill.B <= band.B {
			t.Errorf("%s ToolbarControlFill = %v is not lighter than the chrome %v on every channel",
				c.name, fill, band)
		}
		if c.in.ToolbarControlFill == c.in.PushButtonFill {
			t.Errorf("%s ToolbarControlFill = %v, the push button's fill; the two controls were measured apart",
				c.name, c.in.ToolbarControlFill)
		}
	}
}

// TestToolbarSearchFillIsMeasuredApart pins the toolbar search recess to the
// pixels it was read at and holds it apart from the two fills it would
// otherwise be taken for: the sidebar recess, which the platform draws at a
// different value in the dark appearance, and the bordered toolbar control,
// which it draws at a different value in both.
func TestToolbarSearchFillIsMeasuredApart(t *testing.T) {
	for _, c := range []struct {
		name string
		in   tokens.PlatformColors
		want color.NRGBA
	}{
		{"light", tokens.PlatformLight, color.NRGBA{R: 0xe8, G: 0xe8, B: 0xe8, A: 0xff}},
		{"dark", tokens.PlatformDark, color.NRGBA{R: 0x36, G: 0x36, B: 0x36, A: 0xff}},
	} {
		if c.in.ToolbarSearchFill != c.want {
			t.Errorf("%s ToolbarSearchFill = %v, want the measured toolbar recess %v",
				c.name, c.in.ToolbarSearchFill, c.want)
		}
		if c.in.ToolbarSearchFill == c.in.ToolbarControlFill {
			t.Errorf("%s ToolbarSearchFill = %v, the bordered control's own fill; the two were measured apart",
				c.name, c.in.ToolbarSearchFill)
		}
		fill, band := c.in.ToolbarSearchFill, c.in.SidebarMaterial
		if fill == band {
			t.Errorf("%s ToolbarSearchFill = %v, the chrome material itself; a recess that reads as its band is not a recess",
				c.name, fill)
		}
	}
	// Light, the two recesses are one value; dark they are not, which is why
	// the toolbar's is a name of its own rather than the sidebar's reused.
	if tokens.PlatformLight.ToolbarSearchFill != tokens.PlatformLight.SidebarSearchFill {
		t.Errorf("light ToolbarSearchFill = %v against SidebarSearchFill %v; the two captures read one value",
			tokens.PlatformLight.ToolbarSearchFill, tokens.PlatformLight.SidebarSearchFill)
	}
	if tokens.PlatformDark.ToolbarSearchFill == tokens.PlatformDark.SidebarSearchFill {
		t.Errorf("dark ToolbarSearchFill = %v is the sidebar recess; the two were measured apart",
			tokens.PlatformDark.ToolbarSearchFill)
	}
}

// TestToolbarControlShadowIsTheMeasuredDarkening pins the toolbar control's
// drop shadow to the bytes it was read at: the peak is a coverage of black in
// both appearances, and laid over the band each capture holds it reproduces
// the darkening that capture shows.
func TestToolbarControlShadowIsTheMeasuredDarkening(t *testing.T) {
	for _, c := range []struct {
		name string
		in   tokens.PlatformColors
		want color.NRGBA
		// band is the capture's own toolbar band and darkened what the peak
		// lands on it. Dark, that is the byte the capture reads under the
		// control exactly; light, the capture reads 244 there against this
		// 246, the two 255ths the fitted ramp misses the darkest row by.
		band, darkened color.NRGBA
	}{
		{
			"light", tokens.PlatformLight,
			color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x09},
			color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			color.NRGBA{R: 0xf6, G: 0xf6, B: 0xf6, A: 0xff},
		},
		{
			"dark", tokens.PlatformDark,
			color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x06},
			color.NRGBA{R: 0x1e, G: 0x1e, B: 0x1e, A: 0xff},
			color.NRGBA{R: 0x1d, G: 0x1d, B: 0x1d, A: 0xff},
		},
	} {
		got := c.in.ToolbarControlShadow.Peak
		if got != c.want {
			t.Errorf("%s ToolbarControlShadow = %v, want the measured peak %v", c.name, got, c.want)
		}
		if got.R != 0 || got.G != 0 || got.B != 0 {
			t.Errorf("%s ToolbarControlShadow = %v, want a coverage of black: the platform darkens its band and never tints it", c.name, got)
		}
		if got.A == 0 || got.A == 0xff {
			t.Errorf("%s ToolbarControlShadow alpha = %d, want the platform's partial coverage", c.name, got.A)
		}
		if flat := vgcolor.Flatten(got, c.band); flat != c.darkened {
			t.Errorf("%s: the peak over the band gives %v, want %v", c.name, flat, c.darkened)
		}
	}
	// The light shadow is the deeper one: it is all that tells a #ffffff
	// control from a #ffffff band, where the dark control carries a fill and
	// a rim of its own and the platform leaves its band all but untouched.
	if tokens.PlatformLight.ToolbarControlShadow.Peak.A <= tokens.PlatformDark.ToolbarControlShadow.Peak.A {
		t.Errorf("the light shadow's coverage %d does not exceed the dark one's %d",
			tokens.PlatformLight.ToolbarControlShadow.Peak.A, tokens.PlatformDark.ToolbarControlShadow.Peak.A)
	}
	// It is not the floating surface's shadow: a control standing in a
	// toolbar is not a surface floating over the window.
	for _, c := range []tokens.PlatformColors{tokens.PlatformLight, tokens.PlatformDark} {
		if c.ToolbarControlShadow == c.FloatingShadow {
			t.Errorf("ToolbarControlShadow = %v, the floating surface's own; the two were measured apart", c.ToolbarControlShadow)
		}
	}
}
