package system_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/vibrantgio/theme/system"
	"github.com/vibrantgio/theme/tokens"
)

// platformFor is what the bridge is expected to emit for an appearance
// where the platform publishes an accent colour but no colour set of its
// own — Windows and Linux: the recorded set for its side, with the
// accent-following rows taken from the colour PlatformColor resolves. macOS
// reads the whole set off AppKit instead, so the tests that use this skip
// there.
func platformFor(a system.Appearance) tokens.PlatformColors {
	p := tokens.PlatformLight
	if a.Dark {
		p = tokens.PlatformDark
	}
	if accent, ok := system.PlatformColor(a); ok {
		p = p.WithAccent(accent)
	}
	return p
}

func firstPlatform(t *testing.T, a system.Appearance, opts ...system.Option) tokens.PlatformColors {
	t.Helper()
	src := &fakeSource{vals: []system.Appearance{a}}
	themes, err := collect(system.FromSourceTheme(src, time.Hour, opts...).Take(1))
	if err != nil {
		t.Fatalf("theme observe: %v", err)
	}
	if len(themes) != 1 {
		t.Fatalf("expected 1 theme, got %d", len(themes))
	}
	if themes[0].Platform == nil {
		t.Fatal("the emitted theme carries no platform colour set")
	}
	sets, err := collect(themes[0].Platform)
	if err != nil {
		t.Fatalf("platform observe: %v", err)
	}
	if len(sets) != 1 {
		t.Fatalf("expected 1 platform set, got %d", len(sets))
	}
	return sets[0]
}

// TestThemeCarriesThePlatformSetForTheAppearance pins that the bridge emits
// the platform's own colour set beside the palette: the recorded side for
// the appearance, and the accent rows following the reported accent.
func TestThemeCarriesThePlatformSetForTheAppearance(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("macOS reads the whole set off AppKit; see TestThePlatformSetOnMacOSIsTheLiveSet")
	}
	cases := []struct {
		name string
		app  system.Appearance
	}{
		{"light, no accent reported", system.Appearance{}},
		{"dark, no accent reported", system.Appearance{Dark: true}},
		{"light, a named accent", system.Appearance{Accent: system.AccentPurple}},
		{"dark, a named accent", system.Appearance{Dark: true, Accent: system.AccentGraphite}},
		{"light, a raw accent colour", system.Appearance{AccentColor: tokens.PlatformLight.SystemPink, AccentColorSet: true}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := firstPlatform(t, tc.app), platformFor(tc.app); got != want {
				t.Errorf("platform set = %+v, want %+v", got, want)
			}
		})
	}
}

// TestThePlatformSetOnMacOSIsTheLiveSet pins the macOS half: the set comes
// from AppKit, so the appearance chooses a side and nothing else about the
// Appearance reaches it — a stubbed accent is the same setting read a
// second way and does not move the platform's own answer.
func TestThePlatformSetOnMacOSIsTheLiveSet(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("only macOS reads the set off the platform")
	}
	light := firstPlatform(t, system.Appearance{})
	if dark := firstPlatform(t, system.Appearance{Dark: true}); dark == light {
		t.Error("the two appearances answered the same set")
	}
	if got := firstPlatform(t, system.Appearance{Accent: system.AccentPurple}); got != light {
		t.Error("a stubbed accent moved the set AppKit answers for the light appearance")
	}
	if got := firstPlatform(t, system.Appearance{AccentColor: tokens.PlatformLight.SystemPink, AccentColorSet: true}); got != light {
		t.Error("a stubbed accent colour moved the set AppKit answers for the light appearance")
	}
}

// TestAPinnedThemeColourMovesTheAccentRowsAlone pins the reach of a chosen
// theme colour: it rebuilds the rows the platform derives from the accent
// and leaves every other name where the OS put it.
func TestAPinnedThemeColourMovesTheAccentRowsAlone(t *testing.T) {
	green := tokens.PlatformLight.SystemGreen
	for _, a := range []system.Appearance{{}, {Dark: true}} {
		branded := firstPlatform(t, a, system.WithThemeColor(green))
		if want := firstPlatform(t, a).WithAccent(green); branded != want {
			t.Errorf("a pinned theme colour did not rebuild the accent rows for %+v", a)
		}
		if got := branded.ControlAccent; got.R != green.R || got.G != green.G || got.B != green.B {
			t.Errorf("the accent is %v for %+v, want the chosen %v", got, a, green)
		}
		if branded.WindowBackground != firstPlatform(t, a).WindowBackground {
			t.Errorf("a pinned theme colour moved the window's plane for %+v", a)
		}
	}
}
