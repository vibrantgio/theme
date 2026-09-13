package system

import (
	"os/exec"
	"testing"
)

// staticSource reports one Appearance forever.
type staticSource struct{ a Appearance }

func (s staticSource) Read() (Appearance, error) { return s.a, nil }

// TestTheMacOSMulticolourSettingAnswersSystemBlue: Multicolour is the
// absent AppleAccentColor key, which reports "no accent override" — and on
// macOS the colour a stream with nothing chosen falls back to is the one the
// system itself paints an application that has chosen none.
func TestTheMacOSMulticolourSettingAnswersSystemBlue(t *testing.T) {
	if _, err := exec.Command("defaults", "read", "-g", "AppleAccentColor").Output(); err == nil {
		t.Log("this Mac has a colour chosen; the absent-key half of the test is not exercised")
	} else if got := readAccent(); got != AccentDefault {
		t.Errorf("the absent AppleAccentColor key read as %d, want AccentDefault", got)
	}

	blue, ok := AccentBlue.Color()
	if !ok {
		t.Fatal("AccentBlue carries no colour")
	}
	got, ok := PlatformColor(Appearance{Accent: AccentDefault})
	if !ok {
		t.Fatal("Multicolour answered no colour at all")
	}
	if got != blue {
		t.Errorf("Multicolour answered %v, want systemBlue %v", got, blue)
	}
}

// TestAFailedReadOfTheMacOSAccentColourAnswersSystemBlue: with no
// `defaults` binary to run, the read cannot report a colour and folds to
// "no accent override" — which answers systemBlue like Multicolour, so a Mac
// whose read fails looks like a Mac on Multicolour.
func TestAFailedReadOfTheMacOSAccentColourAnswersSystemBlue(t *testing.T) {
	t.Setenv("PATH", t.TempDir()) // no `defaults` on the path

	if got := readAccent(); got != AccentDefault {
		t.Errorf("a failed read reported %d, want AccentDefault", got)
	}
	src := newDarwinSource()
	a, err := src.Read()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if a.Accent != AccentDefault || a.AccentColorSet {
		t.Fatalf("a failed read produced %+v, want the no-colour appearance", a)
	}

	blue, _ := AccentBlue.Color()
	if got, ok := PlatformColor(a); !ok || got != blue {
		t.Errorf("a failed read answered %v (ok=%v), want systemBlue %v", got, ok, blue)
	}
}

// TestThePlatformColorIsSystemBlue pins the colour itself, independently of
// the accent table: Apple's systemBlue in light appearance, sRGB.
func TestThePlatformColorIsSystemBlue(t *testing.T) {
	c, ok := platformColor()
	if !ok {
		t.Fatal("macOS reports no colour for an application that has chosen none")
	}
	if c.R != 0x00 || c.G != 0x7A || c.B != 0xFF || c.A != 0xFF {
		t.Errorf("the platform colour is %+v, want systemBlue 007AFF", c)
	}
	if blue, _ := AccentBlue.Color(); c != blue {
		t.Error("the platform colour is not the one AccentBlue carries")
	}
}
