package system_test

import (
	"testing"
	"time"

	"github.com/vibrantgio/theme/system"
	"github.com/vibrantgio/theme/tokens"
)

// platformFor is what the bridge is expected to emit for an appearance: the
// recorded set for its side, with the accent-following rows taken from the
// colour PlatformColor resolves.
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

func firstPlatform(t *testing.T, a system.Appearance) tokens.PlatformColors {
	t.Helper()
	src := &fakeSource{vals: []system.Appearance{a}}
	themes, err := collect(system.FromSourceTheme(src, time.Hour).Take(1))
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
	cases := []struct {
		name string
		app  system.Appearance
	}{
		{"light, no accent reported", system.Appearance{}},
		{"dark, no accent reported", system.Appearance{Dark: true}},
		{"light, a named accent", system.Appearance{Accent: system.AccentPurple}},
		{"dark, a named accent", system.Appearance{Dark: true, Accent: system.AccentGraphite}},
		{"light, a raw accent colour", system.Appearance{AccentSeed: tokens.PlatformLight.SystemPink, AccentSeedSet: true}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := firstPlatform(t, tc.app), platformFor(tc.app); got != want {
				t.Errorf("platform set = %+v, want %+v", got, want)
			}
		})
	}
}

// TestThePlatformSetIgnoresAPinnedBrand pins that the platform set is the
// platform's answer and not the application's: an explicit palette option
// steers Color and leaves the platform set where the OS put it.
func TestThePlatformSetIgnoresAPinnedBrand(t *testing.T) {
	src := &fakeSource{vals: []system.Appearance{{}}}
	themes, err := collect(system.FromSourceTheme(src, time.Hour, system.WithSeed(tokens.PlatformLight.SystemGreen)).Take(1))
	if err != nil {
		t.Fatalf("theme observe: %v", err)
	}
	sets, err := collect(themes[0].Platform)
	if err != nil {
		t.Fatalf("platform observe: %v", err)
	}
	if len(sets) != 1 || sets[0] != platformFor(system.Appearance{}) {
		t.Errorf("a pinned brand moved the platform set; got %+v", sets)
	}
}
