package color

import (
	stdcolor "image/color"
	"math"
	"testing"
)

// TestAPCAPublishedVectors verifies APCA against the published APCA-W3
// 0.0.98G-4g test vectors from the reference implementation's README
// (github.com/Myndex/apca-w3, "USAGE — the four-pair example set" and the
// values its live tool reports for these colours). The expectations are the
// reference implementation's full-precision outputs; agreement to 1e-6
// demonstrates the constants and curve are the published ones, not merely
// close.
func TestAPCAPublishedVectors(t *testing.T) {
	hex := func(r, g, b uint8) stdcolor.NRGBA { return stdcolor.NRGBA{r, g, b, 0xff} }
	cases := []struct {
		name     string
		text, bg stdcolor.NRGBA
		want     float64
	}{
		// The apca-w3 README's canonical example pairs, both polarities.
		{"#888 on #fff", hex(0x88, 0x88, 0x88), hex(0xff, 0xff, 0xff), 63.056469930209424},
		{"#fff on #888", hex(0xff, 0xff, 0xff), hex(0x88, 0x88, 0x88), -68.54146436644962},
		{"#000 on #aaa", hex(0x00, 0x00, 0x00), hex(0xaa, 0xaa, 0xaa), 58.146262578561334},
		{"#aaa on #000", hex(0xaa, 0xaa, 0xaa), hex(0x00, 0x00, 0x00), -56.24113336839742},
		{"#123 on #def", hex(0x11, 0x22, 0x33), hex(0xdd, 0xee, 0xff), 91.66830811481631},
		{"#def on #123", hex(0xdd, 0xee, 0xff), hex(0x11, 0x22, 0x33), -93.06770049484275},
		{"#123 on #444", hex(0x11, 0x22, 0x33), hex(0x44, 0x44, 0x44), 8.32326136957393},
	}
	for _, c := range cases {
		if got := APCA(c.text, c.bg); math.Abs(got-c.want) > 1e-6 {
			t.Errorf("APCA(%s) = %.12f, want %.12f", c.name, got, c.want)
		}
	}
}

// TestAPCAConventions verifies the formula's structural behaviour: the sign
// convention, the identical-pair zero, the low-contrast clip, and that
// alpha is ignored.
func TestAPCAConventions(t *testing.T) {
	white := stdcolor.NRGBA{0xff, 0xff, 0xff, 0xff}
	black := stdcolor.NRGBA{0x00, 0x00, 0x00, 0xff}
	if lc := APCA(black, white); lc <= 0 {
		t.Errorf("APCA(black on white) = %.2f, want positive (normal polarity)", lc)
	}
	if lc := APCA(white, black); lc >= 0 {
		t.Errorf("APCA(white on black) = %.2f, want negative (reverse polarity)", lc)
	}
	if lc := APCA(white, white); lc != 0 {
		t.Errorf("APCA(white on white) = %.2f, want 0 (identical pair)", lc)
	}
	// Adjacent greys sit under the low-contrast clip.
	a, b := stdcolor.NRGBA{0x80, 0x80, 0x80, 0xff}, stdcolor.NRGBA{0x84, 0x84, 0x84, 0xff}
	if lc := APCA(a, b); lc != 0 {
		t.Errorf("APCA(near-identical greys) = %.2f, want 0 (low-contrast clip)", lc)
	}
	// Alpha is ignored: NRGBA channels are non-premultiplied.
	translucent := stdcolor.NRGBA{0x00, 0x00, 0x00, 0x40}
	if got, want := APCA(translucent, white), APCA(black, white); got != want {
		t.Errorf("APCA with alpha 0x40 = %.4f, want alpha-independent %.4f", got, want)
	}
}

// TestMagnitudeDropsPolarity pins |Lc| against the same published pairs: the
// magnitude of a pairing is the same number whichever of the two is set on
// the other, up to the curve's own polarity difference.
func TestMagnitudeDropsPolarity(t *testing.T) {
	hex := func(r, g, b uint8) stdcolor.NRGBA { return stdcolor.NRGBA{r, g, b, 0xff} }
	cases := []struct {
		name     string
		text, bg stdcolor.NRGBA
		want     float64
	}{
		{"#888 on #fff", hex(0x88, 0x88, 0x88), hex(0xff, 0xff, 0xff), 63.056469930209424},
		{"#fff on #888", hex(0xff, 0xff, 0xff), hex(0x88, 0x88, 0x88), 68.54146436644962},
		{"#000 on #aaa", hex(0x00, 0x00, 0x00), hex(0xaa, 0xaa, 0xaa), 58.146262578561334},
		{"#aaa on #000", hex(0xaa, 0xaa, 0xaa), hex(0x00, 0x00, 0x00), 56.24113336839742},
	}
	for _, c := range cases {
		if got := Magnitude(c.text, c.bg); math.Abs(got-c.want) > 1e-6 {
			t.Errorf("Magnitude(%s) = %.12f, want %.12f", c.name, got, c.want)
		}
	}
}

// TestBestOnTakesTheGreaterMagnitude verifies the choice rule on the pairing
// that forced it: over the platform's own blue #007AFF, white reads at
// |Lc| 72.0 and black at 37.6, so the label is white — the reverse of what a
// luminance ratio ranks. Ties keep the first candidate offered.
func TestBestOnTakesTheGreaterMagnitude(t *testing.T) {
	white := stdcolor.NRGBA{0xff, 0xff, 0xff, 0xff}
	black := stdcolor.NRGBA{0x00, 0x00, 0x00, 0xff}
	blue := stdcolor.NRGBA{0x00, 0x7a, 0xff, 0xff}
	if got := BestOn(blue, white, black); got != white {
		t.Errorf("BestOn(#007AFF, white, black) = %v, want white (|Lc| %.1f over %.1f)",
			got, Magnitude(white, blue), Magnitude(black, blue))
	}
	if got := BestOn(blue, black, white); got != white {
		t.Errorf("BestOn is order-dependent: got %v, want white", got)
	}
	// A deep fill takes white; the order of the candidates does not decide it.
	deep := stdcolor.NRGBA{0x22, 0x00, 0x4e, 0xff}
	if got := BestOn(deep, black, white); got != white {
		t.Errorf("BestOn(#22004E, black, white) = %v, want white", got)
	}
	// Equals keep the first offered.
	if got := BestOn(blue, white, white); got != white {
		t.Errorf("BestOn with two equal candidates = %v, want the first", got)
	}
	if got := BestOn(blue); got != (stdcolor.NRGBA{}) {
		t.Errorf("BestOn with no candidates = %v, want the zero value", got)
	}
}
