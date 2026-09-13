package system

import (
	"image/color"
	"testing"
)

// TestWindowsSourceCarriesAccentColor verifies the source glue: a registry
// read that yields a colour lands in Appearance.AccentColor with the set
// flag raised, and Dark stays false (dark mode is not read on Windows yet).
func TestWindowsSourceCarriesAccentColor(t *testing.T) {
	want := color.NRGBA{R: 0x00, G: 0x78, B: 0xD7, A: 0xFF}
	src := &windowsSource{
		readAccentFn: func() (color.NRGBA, bool) { return want, true },
	}
	a, err := src.Read()
	if err != nil {
		t.Fatalf("Read(): %v", err)
	}
	if !a.AccentColorSet || a.AccentColor != want {
		t.Errorf("Read() AccentColor=%+v set=%v; want %+v, true", a.AccentColor, a.AccentColorSet, want)
	}
	if a.Dark {
		t.Error("Dark = true; Windows dark mode is not read yet")
	}
}

// TestWindowsSourceNoAccent verifies a failed registry read folds to the
// zero Appearance, per the package's error contract.
func TestWindowsSourceNoAccent(t *testing.T) {
	src := &windowsSource{
		readAccentFn: func() (color.NRGBA, bool) { return color.NRGBA{}, false },
	}
	a, err := src.Read()
	if err != nil {
		t.Fatalf("Read(): %v", err)
	}
	if a != (Appearance{}) {
		t.Errorf("no-accent Read() = %+v; want the zero Appearance", a)
	}
}

// TestReadAccentColorLive exercises the real registry read on a Windows
// machine. Any modern desktop Windows has the DWM key; if the value is
// absent the read must still fold cleanly rather than panic.
func TestReadAccentColorLive(t *testing.T) {
	c, ok := readAccentColor()
	t.Logf("live DWM AccentColor: %+v (set=%v)", c, ok)
	if ok && c.A != 0xFF {
		t.Errorf("live c alpha = %#02x; want opaque", c.A)
	}
}
