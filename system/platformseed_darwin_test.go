package system

import (
	"os/exec"
	"testing"
	"time"

	"github.com/vibrantgio/theme/tokens"
)

// staticSource reports one Appearance forever.
type staticSource struct{ a Appearance }

func (s staticSource) Read() (Appearance, error) { return s.a, nil }

// lightColorsOf resolves one Appearance through the theme bridge with no
// palette option — the "nothing chosen" stream every unbranded application
// runs.
func lightColorsOf(t *testing.T, a Appearance) tokens.ColorTokens {
	t.Helper()
	th, err := FromSourceTheme(staticSource{a}, time.Hour).First()
	if err != nil {
		t.Fatalf("theme: %v", err)
	}
	colors, err := th.Color.First()
	if err != nil {
		t.Fatalf("colours: %v", err)
	}
	return colors
}

// TestTheMacOSMulticolourSettingDerivesFromSystemBlue: Multicolour is the
// absent AppleAccentColor key, which reports "no accent override" — and on
// macOS that derives from the colour the system itself paints an
// application that has chosen none, never from the package's own seed.
func TestTheMacOSMulticolourSettingDerivesFromSystemBlue(t *testing.T) {
	if _, err := exec.Command("defaults", "read", "-g", "AppleAccentColor").Output(); err == nil {
		t.Log("this Mac has a colour chosen; the absent-key half of the test is not exercised")
	} else if got := readAccent(); got != AccentDefault {
		t.Errorf("the absent AppleAccentColor key read as %d, want AccentDefault", got)
	}

	blue, ok := AccentBlue.Seed()
	if !ok {
		t.Fatal("AccentBlue carries no seed")
	}
	wantLight, _ := tokens.FromSeed(blue)
	got := lightColorsOf(t, Appearance{Accent: AccentDefault})
	if got != wantLight {
		t.Error("Multicolour did not derive from systemBlue")
	}
	if got == tokens.DefaultLight {
		t.Error("Multicolour fell through to the package's own default seed")
	}
}

// TestAFailedReadOfTheMacOSAccentColourDerivesFromSystemBlue: with no
// `defaults` binary to run, the read cannot report a colour and folds to
// "no accent override" — which derives from systemBlue like Multicolour,
// so a Mac whose read fails looks like a Mac on Multicolour, not like a
// branded application.
func TestAFailedReadOfTheMacOSAccentColourDerivesFromSystemBlue(t *testing.T) {
	t.Setenv("PATH", t.TempDir()) // no `defaults` on the path

	if got := readAccent(); got != AccentDefault {
		t.Errorf("a failed read reported %d, want AccentDefault", got)
	}
	src := newDarwinSource()
	a, err := src.Read()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if a.Accent != AccentDefault || a.AccentSeedSet {
		t.Fatalf("a failed read produced %+v, want the no-colour appearance", a)
	}

	blue, _ := AccentBlue.Seed()
	wantLight, _ := tokens.FromSeed(blue)
	if got := lightColorsOf(t, a); got != wantLight {
		t.Error("a failed read did not derive from systemBlue")
	}
}

// TestThePlatformSeedIsSystemBlue pins the colour itself, independently of
// the accent table: Apple's systemBlue in light appearance, sRGB.
func TestThePlatformSeedIsSystemBlue(t *testing.T) {
	seed, ok := platformSeed()
	if !ok {
		t.Fatal("macOS reports no colour for an application that has chosen none")
	}
	if seed.R != 0x00 || seed.G != 0x7A || seed.B != 0xFF || seed.A != 0xFF {
		t.Errorf("the platform colour is %+v, want systemBlue 007AFF", seed)
	}
	if blue, _ := AccentBlue.Seed(); seed != blue {
		t.Error("the platform colour is not the seed AccentBlue carries")
	}
}
