package color_test

import (
	stdcolor "image/color"
	"testing"

	"github.com/vibrantgio/theme/color"
)

func nrgba(hex uint32, a uint8) stdcolor.NRGBA {
	return stdcolor.NRGBA{R: uint8(hex >> 16), G: uint8(hex >> 8), B: uint8(hex), A: a}
}

// TestFlattenReproducesTheCapturedComposites pins the blend against pixels
// read off save-dialog-{light,dark}.png in the organization's macOS
// reference rather than computed here. Each row is a glyph core repeated
// across a run of the capture, surrounded by the flat plane it stands on.
//
// The alpha is the platform's own eight-bit coverage, which is 216 for
// labelColor and 140 for secondaryLabelColor: the catalogue records alpha to
// two decimals, so tokens.PlatformColors carries round(0.85×255) = 217 for
// the label and lands one 255th off each of these bytes — the gap the
// platform set's own doc comment records. The blend itself is exact.
func TestFlattenReproducesTheCapturedComposites(t *testing.T) {
	for _, tc := range []struct {
		name              string
		foreground        stdcolor.NRGBA
		surface, captured uint32
		where             string
	}{
		{"labelColor over the light sheet", nrgba(0x000000, 216), 0xffffff, 0x272727,
			"save-dialog-light.png, the sheet's own wording on its #ffffff plane"},
		{"labelColor over the light push button", nrgba(0x000000, 216), 0xececec, 0x242424,
			"save-dialog-light.png, the Cancel button's title on its #ececec fill"},
		{"labelColor over the dark push button", nrgba(0xffffff, 216), 0x333a3f, 0xe0e1e2,
			"save-dialog-dark.png, the Cancel button's title on its #333a3f fill"},
		{"secondaryLabelColor over the dark sheet", nrgba(0xffffff, 140), 0x232a2f, 0x9c9fa1,
			"save-dialog-dark.png, the sheet's secondary wording on its #232a2f plane"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := color.Flatten(tc.foreground, nrgba(tc.surface, 0xff))
			if want := nrgba(tc.captured, 0xff); got != want {
				t.Errorf("Flatten(%v, #%06X) = %v, want #%06X (%s)", tc.foreground, tc.surface, got, tc.captured, tc.where)
			}
		})
	}
}

// TestFlattenIsNotTheLinearBlend: the label the platform puts at #272727 on
// white is #6c6c6c when the same coverage is mixed in linear light, which is
// where Gio's rasterizer mixes a translucent fill it is handed. The two are
// not a shade apart, so handing the rasterizer the platform's alpha is not a
// rounding difference — it is a different colour.
func TestFlattenIsNotTheLinearBlend(t *testing.T) {
	white := nrgba(0xffffff, 0xff)
	got := color.Flatten(nrgba(0x000000, 216), white)
	if want := nrgba(0x272727, 0xff); got != want {
		t.Fatalf("Flatten(black@216, white) = %v, want %v", got, want)
	}
	if l := color.RelativeLuminance(got); l > color.RelativeLuminance(nrgba(0x6c6c6c, 0xff)) {
		t.Errorf("the flattened label is no darker than the linear blend's #6c6c6c")
	}
}

// TestFlattenAtTheEndsOfCoverage: no coverage is the surface and full
// coverage is the foreground, both exactly, so a caller can hand Flatten any
// alpha without special-casing either end.
func TestFlattenAtTheEndsOfCoverage(t *testing.T) {
	surface := nrgba(0x1E293B, 0xff)
	for _, foreground := range []uint32{0x000000, 0xFFFFFF, 0x5C5C5C, 0x3B82F6} {
		if got := color.Flatten(nrgba(foreground, 0), surface); got != surface {
			t.Errorf("Flatten(%06X@0, surface) = %v, want the surface %v", foreground, got, surface)
		}
		if got, want := color.Flatten(nrgba(foreground, 0xff), surface), nrgba(foreground, 0xff); got != want {
			t.Errorf("Flatten(%06X@255, surface) = %v, want the foreground %v", foreground, got, want)
		}
	}
}

// TestFlattenIsMonotonicInCoverage: raising coverage moves the composite
// toward the foreground and never away from it, on every channel.
func TestFlattenIsMonotonicInCoverage(t *testing.T) {
	for _, tc := range []struct{ foreground, surface uint32 }{
		{0x131313, 0xF6F6F6}, {0xEEEEEE, 0x181818}, {0x5C5C5C, 0xE8E8E8},
	} {
		surface := nrgba(tc.surface, 0xff)
		prev := color.RelativeLuminance(surface)
		toward := color.RelativeLuminance(nrgba(tc.foreground, 0xff)) - prev
		for a := 1; a <= 255; a++ {
			l := color.RelativeLuminance(color.Flatten(nrgba(tc.foreground, uint8(a)), surface))
			if (toward < 0 && l > prev) || (toward > 0 && l < prev) {
				t.Fatalf("foreground %06X over %06X: coverage %d moved the composite away from the foreground (%.6f from %.6f)",
					tc.foreground, tc.surface, a, l, prev)
			}
			prev = l
		}
	}
}
