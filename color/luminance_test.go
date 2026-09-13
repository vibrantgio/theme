package color_test

import (
	stdcolor "image/color"
	"math"
	"testing"

	"github.com/vibrantgio/theme/color"
)

// TestRelativeLuminanceAnchors pins relative luminance at the values its own
// numbers force: black 0, white 1, and each pure primary exactly its
// luminance coefficient (channel 1.0 linearizes to 1.0 in both piecewise
// branches' meeting of the formula).
func TestRelativeLuminanceAnchors(t *testing.T) {
	for _, c := range []struct {
		name string
		in   stdcolor.NRGBA
		want float64
	}{
		{"black", stdcolor.NRGBA{A: 255}, 0},
		{"white", stdcolor.NRGBA{R: 255, G: 255, B: 255, A: 255}, 1},
		{"red", stdcolor.NRGBA{R: 255, A: 255}, 0.2126},
		{"green", stdcolor.NRGBA{G: 255, A: 255}, 0.7152},
		{"blue", stdcolor.NRGBA{B: 255, A: 255}, 0.0722},
	} {
		if got := color.RelativeLuminance(c.in); math.Abs(got-c.want) > 1e-12 {
			t.Errorf("RelativeLuminance(%s) = %.13f, want %g", c.name, got, c.want)
		}
	}
	// Alpha is documented as ignored.
	opaque := stdcolor.NRGBA{R: 0x67, G: 0x50, B: 0xa4, A: 0xff}
	translucent := stdcolor.NRGBA{R: 0x67, G: 0x50, B: 0xa4, A: 0x40}
	if color.RelativeLuminance(opaque) != color.RelativeLuminance(translucent) {
		t.Error("RelativeLuminance depends on alpha; it must ignore it")
	}
}
