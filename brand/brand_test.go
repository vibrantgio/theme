package brand_test

import (
	"encoding/json"
	"image/color"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vibrantgio/theme/brand"
	"github.com/vibrantgio/theme/export"
	"github.com/vibrantgio/theme/system"
	"github.com/vibrantgio/theme/tokens"
)

var harbourRed = color.NRGBA{R: 0xe8, G: 0x11, B: 0x2d, A: 0xff}

// file names a file inside a fresh temporary directory. Nothing here ever
// touches the real config directory: the package's own path is asserted
// once, by name, in TestThePathIsOneSharedFileUnderTheConfigDir.
func file(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "theme.json")
}

// TestAKeptSeedRegeneratesBothSchemesExactly is the whole promise of
// keeping one colour instead of a palette: what comes back off disk derives
// the same two schemes, field for field, as what went in.
func TestAKeptSeedRegeneratesBothSchemesExactly(t *testing.T) {
	path := file(t)
	if err := brand.SaveTo(path, brand.Brand{Seed: harbourRed, Source: "harbour.jpg"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, ok, err := brand.LoadFrom(path)
	if err != nil || !ok {
		t.Fatalf("load: got (%v, %v), want a brand and no error", ok, err)
	}
	if got.Seed != harbourRed {
		t.Fatalf("seed came back as %v, want %v", got.Seed, harbourRed)
	}
	wantLight, wantDark := tokens.FromSeed(harbourRed)
	gotLight, gotDark := got.Colors()
	if gotLight != wantLight {
		t.Error("the light scheme regenerated from the kept seed is not the one the seed derives")
	}
	if gotDark != wantDark {
		t.Error("the dark scheme regenerated from the kept seed is not the one the seed derives")
	}
}

// TestProvenanceSurvivesTheRoundTrip: the file has to be able to say what it
// is months later, so where the colour came from and when it was kept come
// back too.
func TestProvenanceSurvivesTheRoundTrip(t *testing.T) {
	path := file(t)
	kept := time.Date(2026, 8, 19, 11, 4, 31, 0, time.UTC)
	if err := brand.SaveTo(path, brand.Brand{Seed: harbourRed, Source: "harbour.jpg", Saved: kept}); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, _, err := brand.LoadFrom(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.Source != "harbour.jpg" {
		t.Errorf("source came back as %q, want %q", got.Source, "harbour.jpg")
	}
	if !got.Saved.Equal(kept) {
		t.Errorf("saved-at came back as %v, want %v", got.Saved, kept)
	}
}

// TestTheFileSpellsItsSeedTheWayTheExportDoes pins the one thing another
// reader has to agree with: the key and the spelling of the colour.
func TestTheFileSpellsItsSeedTheWayTheExportDoes(t *testing.T) {
	path := file(t)
	if err := brand.SaveTo(path, brand.Brand{Seed: harbourRed}); err != nil {
		t.Fatalf("save: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var f struct {
		Seed string `json:"seed"`
	}
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatalf("the file is not JSON: %v", err)
	}
	if f.Seed != "#e8112d" {
		t.Errorf("seed written as %q, want lowercase #rrggbb %q", f.Seed, "#e8112d")
	}
}

// TestAnExportedThemeJSONIsAKeptBrand: the exported project's theme.json
// names its seed under the same key in the same spelling, so it loads here
// without translation — which is why this file has that name and no other
// format was minted for it.
func TestAnExportedThemeJSONIsAKeptBrand(t *testing.T) {
	th, err := system.FromSourceTheme(fixed{}, time.Hour, system.WithSeed(harbourRed)).First()
	if err != nil {
		t.Fatalf("theme: %v", err)
	}
	snapshot, err := export.Capture(th)
	if err != nil {
		t.Fatalf("capture: %v", err)
	}
	dir := t.TempDir()
	if err := export.Write(dir, snapshot); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, ok, err := brand.LoadFrom(filepath.Join(dir, "theme.json"))
	if err != nil || !ok {
		t.Fatalf("load: got (%v, %v), want a brand and no error", ok, err)
	}
	if got.Seed != harbourRed {
		t.Errorf("the exported theme.json read back as %v, want %v", got.Seed, harbourRed)
	}
}

// TestNothingKeptIsExactlyTheDefaultPalette: an application that has never
// been given a brand must behave as though this package were not there.
func TestNothingKeptIsExactlyTheDefaultPalette(t *testing.T) {
	path := filepath.Join(t.TempDir(), "absent.json")
	b, ok, err := brand.LoadFrom(path)
	if ok || err != nil {
		t.Fatalf("a missing file read as (%v, %v), want no brand and no error", ok, err)
	}
	if b.Chosen() {
		t.Error("a missing file produced a brand")
	}
	assertDefaults(t, brand.KeptFrom(path))
}

// TestAFileThatWillNotParseIsNotACrash: every shape of damage a hand-edited
// or half-written file can have leaves the application on the defaults, and
// says why to a caller that asks.
func TestAFileThatWillNotParseIsNotACrash(t *testing.T) {
	for _, tc := range []struct{ name, body string }{
		{"empty", ""},
		{"truncated", `{"seed": "#e811`},
		{"not json at all", "harbour red, I think\n"},
		{"no seed", `{"source": "harbour.jpg"}`},
		{"seed is a name", `{"seed": "red"}`},
		{"seed is short", `{"seed": "#e12"}`},
		{"seed is not hex", `{"seed": "#zzzzzz"}`},
		{"seed is a number", `{"seed": 15208749}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := file(t)
			if err := os.WriteFile(path, []byte(tc.body), 0o644); err != nil {
				t.Fatalf("write: %v", err)
			}
			if _, ok, err := brand.LoadFrom(path); ok || err == nil {
				t.Errorf("a damaged file read as (%v, %v), want no brand and an error saying so", ok, err)
			}
			assertDefaults(t, brand.KeptFrom(path))
		})
	}
}

// TestABadTimestampCostsTheProvenanceAndNotTheBrand: the colour is what the
// file is for, and it parsed.
func TestABadTimestampCostsTheProvenanceAndNotTheBrand(t *testing.T) {
	path := file(t)
	if err := os.WriteFile(path, []byte(`{"seed": "#e8112d", "saved": "last tuesday"}`), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, ok, err := brand.LoadFrom(path)
	if err != nil || !ok {
		t.Fatalf("load: got (%v, %v), want the brand and no error", ok, err)
	}
	if got.Seed != harbourRed {
		t.Errorf("seed came back as %v, want %v", got.Seed, harbourRed)
	}
	if !got.Saved.IsZero() {
		t.Errorf("an unreadable timestamp came back as %v, want the zero time", got.Saved)
	}
}

// TestTheKeptBrandPinsTheStreamsPaletteOnBothSides is the adoption seam: the
// options put the kept seed on a live theme stream, and the OS keeps
// deciding which side of the pair shows.
func TestTheKeptBrandPinsTheStreamsPaletteOnBothSides(t *testing.T) {
	path := file(t)
	if err := brand.SaveTo(path, brand.Brand{Seed: harbourRed}); err != nil {
		t.Fatalf("save: %v", err)
	}
	opts := brand.KeptFrom(path).Options()
	light, dark := tokens.FromSeed(harbourRed)
	for _, tc := range []struct {
		name string
		app  system.Appearance
		want tokens.ColorTokens
	}{
		{"the desktop is light", system.Appearance{}, light},
		{"the desktop is dark", system.Appearance{Dark: true}, dark},
		// An accent the OS reports loses to a brand a person chose.
		{"the desktop has an accent", system.Appearance{AccentSeed: color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xff}, AccentSeedSet: true}, light},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := system.FromSourceTheme(fixed{tc.app}, time.Hour, opts...).First()
			if err != nil {
				t.Fatalf("theme: %v", err)
			}
			colors, err := got.Color.First()
			if err != nil {
				t.Fatalf("colours: %v", err)
			}
			if colors != tc.want {
				t.Error("the stream did not emit the kept brand's scheme")
			}
		})
	}
}

// TestWithNothingKeptTheStreamIsTheOneItAlwaysWas: splatting an empty option
// set is the same call as making none.
func TestWithNothingKeptTheStreamIsTheOneItAlwaysWas(t *testing.T) {
	opts := brand.KeptFrom(filepath.Join(t.TempDir(), "absent.json")).Options()
	if opts != nil {
		t.Fatalf("nothing kept produced %d options, want none", len(opts))
	}
	got, err := system.FromSourceTheme(fixed{}, time.Hour, opts...).First()
	if err != nil {
		t.Fatalf("theme: %v", err)
	}
	colors, err := got.Color.First()
	if err != nil {
		t.Fatalf("colours: %v", err)
	}
	if colors != tokens.DefaultLight {
		t.Error("with nothing kept the stream did not emit the default palette")
	}
}

// TestSavingNoColourIsRefused: a file the reader would take for an absent
// one is a slower way of deleting it, so it is never written.
func TestSavingNoColourIsRefused(t *testing.T) {
	path := file(t)
	if err := brand.SaveTo(path, brand.Brand{Source: "harbour.jpg"}); err == nil {
		t.Error("a brand with no colour was saved")
	}
	if _, err := os.Stat(path); err == nil {
		t.Error("a file was written for a brand with no colour")
	}
}

// TestSaveFillsInWhenItWasKept: a caller that does not care about the clock
// still writes an honest file.
func TestSaveFillsInWhenItWasKept(t *testing.T) {
	path := file(t)
	before := time.Now().Add(-time.Second)
	if err := brand.SaveTo(path, brand.Brand{Seed: harbourRed}); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, _, err := brand.LoadFrom(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.Saved.Before(before) {
		t.Errorf("saved-at came back as %v, want a time from this test run", got.Saved)
	}
}

// TestKeepingAgainReplacesTheBrand: choosing a second colour is not
// appending to a list.
func TestKeepingAgainReplacesTheBrand(t *testing.T) {
	path := file(t)
	teal := color.NRGBA{R: 0x00, G: 0x7a, B: 0x7a, A: 0xff}
	if err := brand.SaveTo(path, brand.Brand{Seed: harbourRed, Source: "harbour.jpg"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := brand.SaveTo(path, brand.Brand{Seed: teal, Source: "lagoon.png"}); err != nil {
		t.Fatalf("save again: %v", err)
	}
	got := brand.KeptFrom(path)
	if got.Seed != teal || got.Source != "lagoon.png" {
		t.Errorf("the second keep read back as %v from %q, want %v from %q", got.Seed, got.Source, teal, "lagoon.png")
	}
}

// TestThePathIsOneSharedFileUnderTheConfigDir: the path is per user and not
// per application, because everything the person opens is meant to wear the
// same brand.
func TestThePathIsOneSharedFileUnderTheConfigDir(t *testing.T) {
	path, err := brand.Path()
	if err != nil {
		t.Skipf("this machine has no user config directory: %v", err)
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		t.Skipf("this machine has no user config directory: %v", err)
	}
	want := filepath.Join(dir, "vibrantgio", "theme.json")
	if path != want {
		t.Errorf("the kept brand lives at %s, want %s", path, want)
	}
}

// TestTheChosenBaseSurvivesTheRoundTrip: the second choice the file carries
// comes back as it went in, under its own key, beside the seed.
func TestTheChosenBaseSurvivesTheRoundTrip(t *testing.T) {
	path := file(t)
	if err := brand.SaveTo(path, brand.Brand{Seed: harbourRed, Base: "catppuccin-latte"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	got := brand.KeptFrom(path)
	if got.Base != "catppuccin-latte" {
		t.Errorf("the base came back as %q, want catppuccin-latte", got.Base)
	}
	var raw map[string]any
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("the file is not JSON: %v", err)
	}
	if raw["base"] != "catppuccin-latte" {
		t.Errorf("the file spells the base as %v under \"base\", want catppuccin-latte", raw["base"])
	}
}

// TestAFileWithNoBaseIsStillAKeptBrand: every theme.json written before the
// field existed has no base in it, and reading one has to be uneventful —
// the brand loads, the seed is intact, and the base comes back empty, which
// is the value that means "the reader's own default applies".
func TestAFileWithNoBaseIsStillAKeptBrand(t *testing.T) {
	path := file(t)
	const older = `{"seed":"#e8112d","source":"harbour.jpg"}`
	if err := os.WriteFile(path, []byte(older), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, ok, err := brand.LoadFrom(path)
	if err != nil || !ok {
		t.Fatalf("load: got (%v, %v), want a brand and no error", ok, err)
	}
	if got.Seed != harbourRed {
		t.Errorf("seed came back as %v, want %v", got.Seed, harbourRed)
	}
	if got.Base != "" {
		t.Errorf("a file with no base loaded base %q, want none", got.Base)
	}
}

// TestNoBaseIsNotWritten: an unchosen base leaves no key behind, so the file
// says what was chosen and nothing else.
func TestNoBaseIsNotWritten(t *testing.T) {
	path := file(t)
	if err := brand.SaveTo(path, brand.Brand{Seed: harbourRed}); err != nil {
		t.Fatalf("save: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("the file is not JSON: %v", err)
	}
	if _, ok := raw["base"]; ok {
		t.Errorf("a brand with no base wrote %q into the file", "base")
	}
}

// TestTheStylesFolderSitsBesideTheFile: styles a person adds are shared the
// way the brand is, in one folder under the same directory, so a style added
// once is offered by everything that looks for one.
func TestTheStylesFolderSitsBesideTheFile(t *testing.T) {
	dir, err := brand.StylesDir()
	if err != nil {
		t.Skipf("this machine has no user config directory: %v", err)
	}
	path, err := brand.Path()
	if err != nil {
		t.Skipf("this machine has no user config directory: %v", err)
	}
	if want := filepath.Join(filepath.Dir(path), "styles"); dir != want {
		t.Errorf("the styles folder is %s, want %s", dir, want)
	}
}

// assertDefaults asserts that a brand carrying nothing leaves an
// application exactly where it was.
func assertDefaults(t *testing.T, b brand.Brand) {
	t.Helper()
	if b.Chosen() {
		t.Error("a brand was reported where there is none")
	}
	if opts := b.Options(); opts != nil {
		t.Errorf("no brand produced %d stream options, want none", len(opts))
	}
	light, dark := b.Colors()
	if light != tokens.DefaultLight || dark != tokens.DefaultDark {
		t.Error("no brand produced something other than the default palette")
	}
}

// fixed is an appearance source that always reports the same thing, so a
// test's emissions cannot depend on the desktop it runs on.
type fixed struct{ a system.Appearance }

func (f fixed) Read() (system.Appearance, error) { return f.a, nil }
