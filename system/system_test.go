package system_test

import (
	"context"
	"image/color"
	"runtime"
	"testing"
	"time"

	"github.com/reactivego/rx"
	"github.com/vibrantgio/theme/a11y"
	"github.com/vibrantgio/theme/system"
	"github.com/vibrantgio/theme/tokens"
)

// fakeSource returns successive values from vals on each Read call,
// repeating the last value once the slice is exhausted. Mirrors the
// pattern used in theme/a11y/preferences_test.go.
type fakeSource struct {
	vals []system.Appearance
	n    int
}

func (f *fakeSource) Read() (system.Appearance, error) {
	v := f.vals[f.n]
	if f.n < len(f.vals)-1 {
		f.n++
	}
	return v, nil
}

func collect[T any](obs rx.Observable[T]) ([]T, error) {
	var out []T
	err := obs.Subscribe(context.Background(), func(v T, _ error, done bool) {
		if !done {
			out = append(out, v)
		}
	}).Wait()
	return out, err
}

func TestFromSourceEmitsInitialValue(t *testing.T) {
	want := system.Appearance{Dark: true, Accent: system.AccentBlue}
	src := &fakeSource{vals: []system.Appearance{want}}

	got, err := collect(system.FromSource(src, time.Hour).Take(1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 emission, got %d", len(got))
	}
	if got[0] != want {
		t.Errorf("got %+v, want %+v", got[0], want)
	}
}

func TestFromSourceEmitsOnDarkChange(t *testing.T) {
	light := system.Appearance{Dark: false}
	dark := system.Appearance{Dark: true}
	src := &fakeSource{vals: []system.Appearance{light, dark}}

	got, err := collect(system.FromSource(src, time.Millisecond).Take(2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 emissions, got %d", len(got))
	}
	if got[0] != light {
		t.Errorf("first: got %+v, want %+v", got[0], light)
	}
	if got[1] != dark {
		t.Errorf("second: got %+v, want %+v", got[1], dark)
	}
}

func TestFromSourceDeduplicates(t *testing.T) {
	a := system.Appearance{Dark: false}
	b := system.Appearance{Dark: true}
	src := &fakeSource{vals: []system.Appearance{a, a, b}}

	got, err := collect(system.FromSource(src, time.Millisecond).Take(2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 emissions (a then b), got %d", len(got))
	}
	if got[0] != a || got[1] != b {
		t.Errorf("got %+v then %+v, want %+v then %+v", got[0], got[1], a, b)
	}
}

func TestFromSourceEmitsOnAccentChange(t *testing.T) {
	a := system.Appearance{Accent: system.AccentBlue}
	b := system.Appearance{Accent: system.AccentRed}
	src := &fakeSource{vals: []system.Appearance{a, b}}

	got, err := collect(system.FromSource(src, time.Millisecond).Take(2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 emissions, got %d", len(got))
	}
	if got[0].Accent != system.AccentBlue || got[1].Accent != system.AccentRed {
		t.Errorf("accent transitions wrong: %+v then %+v", got[0], got[1])
	}
}

func TestFromSourceThemeBridgesDarkToTheDarkSet(t *testing.T) {
	src := &fakeSource{vals: []system.Appearance{{Dark: true}}}

	themes, err := collect(system.FromSourceTheme(src, time.Hour).Take(1))
	if err != nil {
		t.Fatalf("theme observe: %v", err)
	}
	if len(themes) != 1 {
		t.Fatalf("expected 1 theme, got %d", len(themes))
	}
	colors, err := collect(themes[0].Platform)
	if err != nil {
		t.Fatalf("color observe: %v", err)
	}
	wantDark := firstPlatform(t, system.Appearance{Dark: true})
	if len(colors) != 1 || colors[0] != wantDark {
		t.Errorf("dark appearance must yield the platform's dark side; got %+v", colors)
	}
}

func TestFromSourceThemeBridgesLightToTheLightSet(t *testing.T) {
	src := &fakeSource{vals: []system.Appearance{{Dark: false}}}

	themes, err := collect(system.FromSourceTheme(src, time.Hour).Take(1))
	if err != nil {
		t.Fatalf("theme observe: %v", err)
	}
	colors, err := collect(themes[0].Platform)
	if err != nil {
		t.Fatalf("color observe: %v", err)
	}
	wantLight := firstPlatform(t, system.Appearance{})
	if len(colors) != 1 || colors[0] != wantLight {
		t.Errorf("light appearance must yield the platform's light side; got %+v", colors)
	}
}

func TestFromSourceThemeReemitsOnChange(t *testing.T) {
	light := system.Appearance{Dark: false}
	dark := system.Appearance{Dark: true}
	src := &fakeSource{vals: []system.Appearance{light, dark}}

	themes, err := collect(system.FromSourceTheme(src, time.Millisecond).Take(2))
	if err != nil {
		t.Fatalf("theme observe: %v", err)
	}
	if len(themes) != 2 {
		t.Fatalf("expected 2 themes, got %d", len(themes))
	}
	wantLight := firstPlatform(t, system.Appearance{})
	wantDark := firstPlatform(t, system.Appearance{Dark: true})
	for i, want := range []tokens.PlatformColors{wantLight, wantDark} {
		colors, err := collect(themes[i].Platform)
		if err != nil {
			t.Fatalf("theme[%d] color observe: %v", i, err)
		}
		if len(colors) != 1 || colors[0] != want {
			t.Errorf("theme[%d] colors: got %+v, want %+v", i, colors, want)
		}
	}
}

// customThemeColor is a colour distinct from the platform's own accent, so
// any leak of the platform's answer into a pinned stream is detectable.
var customThemeColor = color.NRGBA{R: 0x00, G: 0x6E, B: 0x2E, A: 0xff}

// TestFromSourceThemeWithThemeColorRebuildsTheAccentRows: a pinned theme
// colour is what the emitted set's accent rows are rebuilt for, in both
// appearances, and nothing else in the set moves.
func TestFromSourceThemeWithThemeColorRebuildsTheAccentRows(t *testing.T) {
	for _, a := range []system.Appearance{{}, {Dark: true}} {
		src := &fakeSource{vals: []system.Appearance{a}}
		themes, err := collect(system.FromSourceTheme(src, time.Hour, system.WithThemeColor(customThemeColor)).Take(1))
		if err != nil || len(themes) != 1 {
			t.Fatalf("theme observe: err=%v len=%d", err, len(themes))
		}
		colors, err := collect(themes[0].Platform)
		if err != nil || len(colors) != 1 {
			t.Fatalf("color observe: err=%v len=%d", err, len(colors))
		}
		if want := firstPlatform(t, a).WithAccent(customThemeColor); colors[0] != want {
			t.Errorf("%+v: the emitted set is not the platform's with the accent rows rebuilt", a)
		}
	}
}

// TestFromSourceThemeWithThemeColorBeatsTheOSAccent: the application chose
// its colour, so whatever the desktop reports is ignored.
func TestFromSourceThemeWithThemeColorBeatsTheOSAccent(t *testing.T) {
	for _, a := range []system.Appearance{
		{Accent: system.AccentPurple},
		{AccentColor: rawAccent, AccentColorSet: true},
	} {
		src := &fakeSource{vals: []system.Appearance{a}}
		themes, err := collect(system.FromSourceTheme(src, time.Hour, system.WithThemeColor(customThemeColor)).Take(1))
		if err != nil || len(themes) != 1 {
			t.Fatalf("theme observe: err=%v len=%d", err, len(themes))
		}
		colors, err := collect(themes[0].Platform)
		if err != nil || len(colors) != 1 {
			t.Fatalf("color observe: err=%v len=%d", err, len(colors))
		}
		got := colors[0].ControlAccent
		if got.R != customThemeColor.R || got.G != customThemeColor.G || got.B != customThemeColor.B {
			t.Errorf("%+v: the accent is %v, want the chosen %v", a, got, customThemeColor)
		}
	}
}

// rawAccent is an arbitrary colour of the kind the Windows registry or a
// KDE kdeglobals delivers — deliberately none of the enum accent colours.
var rawAccent = color.NRGBA{R: 0x00, G: 0x78, B: 0xD7, A: 0xFF} // Windows default blue

func TestFromSourceEmitsOnAccentColorChange(t *testing.T) {
	a := system.Appearance{AccentColor: rawAccent, AccentColorSet: true}
	b := system.Appearance{AccentColor: color.NRGBA{R: 0xE6, G: 0x2D, B: 0x42, A: 0xFF}, AccentColorSet: true}
	src := &fakeSource{vals: []system.Appearance{a, a, b}}

	got, err := collect(system.FromSource(src, time.Millisecond).Take(2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != a || got[1] != b {
		t.Errorf("colour transitions wrong: got %+v", got)
	}
}

func TestFromSourceThemeWithTypographyEmitsIt(t *testing.T) {
	src := &fakeSource{vals: []system.Appearance{{}}}
	want := tokens.CodeFace("JetBrains Mono")

	themes, err := collect(system.FromSourceTheme(src, time.Hour, system.WithTypography(want)).Take(1))
	if err != nil {
		t.Fatalf("theme observe: %v", err)
	}
	if len(themes) != 1 {
		t.Fatalf("expected 1 theme, got %d", len(themes))
	}
	got, err := collect(themes[0].Typography)
	if err != nil {
		t.Fatalf("typography observe: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 typography, got %d", len(got))
	}
	if got[0].Code.Typeface != "JetBrains Mono" {
		t.Errorf("Code.Typeface = %q, want JetBrains Mono", got[0].Code.Typeface)
	}
	if got[0].Shaper() != want.Shaper() {
		t.Error("the stream built its own shaper instead of emitting the one WithTypography was given")
	}
}

func TestFromSourceThemeWithoutTypographyIsTheDefault(t *testing.T) {
	src := &fakeSource{vals: []system.Appearance{{}}}
	themes, err := collect(system.FromSourceTheme(src, time.Hour).Take(1))
	if err != nil {
		t.Fatalf("theme observe: %v", err)
	}
	got, err := collect(themes[0].Typography)
	if err != nil {
		t.Fatalf("typography observe: %v", err)
	}
	live := tokens.EmojiTypography()
	if len(got) != 1 || got[0].Shaper() != live.Shaper() {
		t.Error("with no WithTypography the stream did not emit EmojiTypography")
	}
}

// --- accessibility preferences composed into the theme ---

// fakeA11ySource returns successive values from vals on each Read call,
// repeating the last value once the slice is exhausted — the a11y twin of
// fakeSource above.
type fakeA11ySource struct {
	vals []a11y.A11yPrefs
	n    int
}

func (f *fakeA11ySource) Read() (a11y.A11yPrefs, error) {
	v := f.vals[f.n]
	if f.n < len(f.vals)-1 {
		f.n++
	}
	return v, nil
}

// TestFromSourceThemeReduceMotionSnaps is the snap test: under an OS
// reduce-motion preference the emitted motion scale has every duration at
// zero, so an animated component that derives its frame count from the
// scale — round(d·fps) = 0 frames for every stop — is at its target on the
// first frame it draws.
func TestFromSourceThemeReduceMotionSnaps(t *testing.T) {
	appearance := &fakeSource{vals: []system.Appearance{{}}}
	prefs := &fakeA11ySource{vals: []a11y.A11yPrefs{{ReduceMotion: true}}}

	themes, err := collect(system.FromSourceTheme(appearance, time.Hour, system.WithA11ySource(prefs)).Take(1))
	if err != nil {
		t.Fatalf("theme observe: %v", err)
	}
	if len(themes) != 1 {
		t.Fatalf("expected 1 theme, got %d", len(themes))
	}
	motions, err := collect(themes[0].Motion)
	if err != nil {
		t.Fatalf("motion observe: %v", err)
	}
	if len(motions) != 1 {
		t.Fatalf("expected 1 motion emission, got %d", len(motions))
	}
	m := motions[0]
	for _, d := range []struct {
		name string
		v    time.Duration
	}{
		{"DurXFast", m.DurXFast}, {"DurFast", m.DurFast}, {"DurNormal", m.DurNormal},
		{"DurSlow", m.DurSlow}, {"DurXSlow", m.DurXSlow},
	} {
		if d.v != 0 {
			t.Errorf("ReduceMotion: %s = %v, want 0 (animations must reach their target in one frame)", d.name, d.v)
		}
	}
	if m != tokens.Motion.Reduced() {
		t.Errorf("ReduceMotion: motion scale is not tokens.Motion.Reduced():\ngot %+v", m)
	}
}

func TestFromSourceThemeReduceMotionOffKeepsMotion(t *testing.T) {
	appearance := &fakeSource{vals: []system.Appearance{{}}}
	prefs := &fakeA11ySource{vals: []a11y.A11yPrefs{{}}}

	themes, err := collect(system.FromSourceTheme(appearance, time.Hour, system.WithA11ySource(prefs)).Take(1))
	if err != nil {
		t.Fatalf("theme observe: %v", err)
	}
	motions, err := collect(themes[0].Motion)
	if err != nil {
		t.Fatalf("motion observe: %v", err)
	}
	if len(motions) != 1 || motions[0] != tokens.Motion {
		t.Errorf("all-off prefs: motion scale must be tokens.Motion unchanged; got %+v", motions)
	}
}

func TestFromSourceThemeDefaultA11yIsHermetic(t *testing.T) {
	// Without WithA11ySource, FromSourceTheme must not read the machine's
	// real accessibility preferences: the default is a constant all-off
	// stream, so this asserts the full default emission regardless of host.
	appearance := &fakeSource{vals: []system.Appearance{{}}}

	themes, err := collect(system.FromSourceTheme(appearance, time.Hour).Take(1))
	if err != nil {
		t.Fatalf("theme observe: %v", err)
	}
	motions, err := collect(themes[0].Motion)
	if err != nil {
		t.Fatalf("motion observe: %v", err)
	}
	if len(motions) != 1 || motions[0] != tokens.Motion {
		t.Errorf("default a11y stream must be all-off: got motion %+v", motions)
	}
	colors, err := collect(themes[0].Platform)
	if err != nil {
		t.Fatalf("color observe: %v", err)
	}
	wantLight := firstPlatform(t, system.Appearance{})
	if len(colors) != 1 || colors[0] != wantLight {
		t.Errorf("default a11y stream must be all-off: got colors %+v", colors)
	}
}

func TestFromSourceThemeReduceMotionToggleReemits(t *testing.T) {
	// A preference change alone — the appearance never changes — re-emits
	// the theme, normal motion first, reduced second.
	appearance := &fakeSource{vals: []system.Appearance{{}}}
	prefs := &fakeA11ySource{vals: []a11y.A11yPrefs{{}, {ReduceMotion: true}}}

	themes, err := collect(system.FromSourceTheme(appearance, time.Millisecond, system.WithA11ySource(prefs)).Take(2))
	if err != nil {
		t.Fatalf("theme observe: %v", err)
	}
	if len(themes) != 2 {
		t.Fatalf("expected 2 themes (one per preference value), got %d", len(themes))
	}
	for i, want := range []tokens.MotionScale{tokens.Motion, tokens.Motion.Reduced()} {
		motions, err := collect(themes[i].Motion)
		if err != nil {
			t.Fatalf("theme[%d] motion observe: %v", i, err)
		}
		if len(motions) != 1 || motions[0] != want {
			t.Errorf("theme[%d] motion: got %+v, want %+v", i, motions, want)
		}
	}
}

// TestPlatformColorWalksTheFallthrough walks the accessor down the same
// fallthrough the stream applies, and pins the per-platform end of it: an
// application that offers this colour as a choice and a stream with nothing
// chosen must never disagree about what it is.
func TestPlatformColorWalksTheFallthrough(t *testing.T) {
	desktop := color.NRGBA{R: 0x35, G: 0x84, B: 0xE4, A: 0xFF}
	for _, tc := range []struct {
		name string
		app  system.Appearance
		want color.NRGBA
		ok   bool
	}{
		{"a desktop reporting a colour of its own", system.Appearance{AccentColor: desktop, AccentColorSet: true}, desktop, true},
		{"a macOS accent colour chosen by name", system.Appearance{Accent: system.AccentPink}, colorOf(t, system.AccentPink), true},
		// The raw colour wins where a source ever reports both, which is the
		// order the stream resolves them in.
		{"both", system.Appearance{AccentColor: desktop, AccentColorSet: true, Accent: system.AccentPink}, desktop, true},
		// Multicolour, an unsupported desktop, a failed read: the platform's
		// own colour, where it has one.
		{"nothing chosen", system.Appearance{}, platformOwnColor(t), runtime.GOOS == "darwin"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := system.PlatformColor(tc.app)
			if ok != tc.ok {
				t.Fatalf("PlatformColor reported ok=%v on %s, want %v", ok, runtime.GOOS, tc.ok)
			}
			if ok && got != tc.want {
				t.Errorf("PlatformColor reported %v, want %v", got, tc.want)
			}
		})
	}
}

// platformOwnColor is the colour this platform paints an application that
// has chosen none, or the zero colour where it paints none.
func platformOwnColor(t *testing.T) color.NRGBA {
	t.Helper()
	c, _ := system.PlatformColor(system.Appearance{})
	return c
}

// colorOf is the colour one named macOS accent carries.
func colorOf(t *testing.T, a system.Accent) color.NRGBA {
	t.Helper()
	c, ok := a.Color()
	if !ok {
		t.Fatalf("%v carries no colour", a)
	}
	return c
}
