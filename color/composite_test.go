package color_test

import (
	stdcolor "image/color"
	"math"
	"testing"

	"github.com/vibrantgio/theme/color"
)

func nrgba(hex uint32, a uint8) stdcolor.NRGBA {
	return stdcolor.NRGBA{R: uint8(hex >> 16), G: uint8(hex >> 8), B: uint8(hex), A: a}
}

// TestOverBlendsInLinearLight pins the blend against a pixel that was read
// off a recorded golden image rather than computed here: the design
// system's overlay scrollbar thumb — neutral step 700, #5C5C5C, at coverage
// 100 — renders as #C1C1C1 over the light Surface #E8E8E8 in
// components/scrollbar's stored render. That is what Gio's rasterizer
// actually wrote, so it is what Over has to answer; the eight-bit average
// of the same two colours is #B1B1B1, sixteen levels away.
func TestOverBlendsInLinearLight(t *testing.T) {
	got := color.Over(nrgba(0x5C5C5C, 100), nrgba(0xE8E8E8, 0xff))
	if want := nrgba(0xC1C1C1, 0xff); got != want {
		t.Errorf("Over(#5C5C5C@100, #E8E8E8) = %v, want %v (the recorded golden pixel)", got, want)
	}
	// The same ink over the light page, the pairing Phase AS was opened by.
	page := nrgba(0xF6F6F6, 0xff)
	comp := color.Over(nrgba(0x5C5C5C, 100), page)
	if want := nrgba(0xCCCCCC, 0xff); comp != want {
		t.Errorf("Over(#5C5C5C@100, #F6F6F6) = %v, want %v", comp, want)
	}
	if got := color.ContrastRatio(comp, page); math.Abs(got-1.49) > 0.005 {
		t.Errorf("the composited thumb measures %.2f:1 against the page, want 1.49:1", got)
	}
}

// TestOverAtTheEndsOfCoverage: no coverage is the ground and full coverage
// is the ink, both exactly, so a caller can hand Over any alpha without
// special-casing either end.
func TestOverAtTheEndsOfCoverage(t *testing.T) {
	ground := nrgba(0x1E293B, 0xff)
	for _, ink := range []uint32{0x000000, 0xFFFFFF, 0x5C5C5C, 0x3B82F6} {
		if got := color.Over(nrgba(ink, 0), ground); got != ground {
			t.Errorf("Over(%06X@0, ground) = %v, want the ground %v", ink, got, ground)
		}
		if got, want := color.Over(nrgba(ink, 0xff), ground), nrgba(ink, 0xff); got != want {
			t.Errorf("Over(%06X@255, ground) = %v, want the ink %v", ink, got, want)
		}
	}
}

// TestOverIsMonotonicInCoverage: raising coverage moves the composite
// toward the ink and never away from it. The derivations that solve for a
// coverage — components/scrollbar's thumb above all — walk alpha upward and
// stop at the first value that clears a floor, which is only the least such
// value if the walk is monotonic.
func TestOverIsMonotonicInCoverage(t *testing.T) {
	for _, tc := range []struct{ ink, ground uint32 }{
		{0x131313, 0xF6F6F6}, {0xEEEEEE, 0x181818}, {0x5C5C5C, 0xE8E8E8},
	} {
		ground := nrgba(tc.ground, 0xff)
		prev := color.RelativeLuminance(ground)
		toward := color.RelativeLuminance(nrgba(tc.ink, 0xff)) - prev
		for a := 1; a <= 255; a++ {
			l := color.RelativeLuminance(color.Over(nrgba(tc.ink, uint8(a)), ground))
			if (toward < 0 && l > prev) || (toward > 0 && l < prev) {
				t.Fatalf("ink %06X over %06X: coverage %d moved the composite away from the ink (%.6f from %.6f)",
					tc.ink, tc.ground, a, l, prev)
			}
			prev = l
		}
	}
}
