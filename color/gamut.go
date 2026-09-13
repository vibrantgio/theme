// Gamut mapping for the OKLab inverse: chroma reduction at constant OKLab
// lightness and constant hue.
//
// Near the lightness extremes a requested (L, a, b) lies outside sRGB, and
// clamping R, G and B independently shifts hue and chroma badly. So chroma
// walks toward the neutral axis instead, keeping lightness and hue exact,
// until the colour fits: the neutral axis is in gamut for every L in [0,1],
// so the bisection converges on the gamut boundary from inside. Hue is held
// exactly by construction — every candidate's a,b sit on the requested hue
// ray, so the only hue movement in a mapped result is 8-bit quantization
// noise.
package color

import "math"

// gamutEps is the linear-sRGB tolerance inside which a channel counts as
// in gamut. Two kinds of residue must fall inside it: the boundary
// overshoot the bisection leaves, and the round-trip residue of the
// conversion chains, measured at up to ~7e-7 over the sRGB cube. It must
// also stay below half the smallest 8-bit quantization step (~1.5e-4 in
// linear near black), so that snapping the residue onto the boundary picks
// the same byte true gamut mapping would. 1e-4 sits safely between the
// two.
const gamutEps = 1e-4

// inSRGBGamut reports whether a linear sRGB triple is inside the unit
// cube, within gamutEps of numeric residue per channel.
func inSRGBGamut(red, green, blue float64) bool {
	const lo, hi = -gamutEps, 1 + gamutEps
	return red >= lo && red <= hi &&
		green >= lo && green <= hi &&
		blue >= lo && blue <= hi
}

// quantizeLinear converts a linear sRGB triple to 8-bit channels. The
// clamp only snaps the ≤ gamutEps boundary residue the bisection leaves —
// it is not a per-channel gamut clamp; callers guarantee the triple is in
// gamut within gamutEps.
func quantizeLinear(red, green, blue float64) (R, G, B uint8) {
	q := func(v float64) uint8 {
		if v < 0 {
			v = 0
		} else if v > 1 {
			v = 1
		}
		return uint8(math.Round(255.0 * SRGBFromLinear(v)))
	}
	return q(red), q(green), q(blue)
}

// mapOKLabChroma scales a,b down toward the neutral axis at constant L
// until the linear sRGB triple fits the cube.
func mapOKLabChroma(L, a, b float64) (red, green, blue float64) {
	red, green, blue = LinearRGBFromOKLab(L, 0, 0)
	lo, hi := 0.0, 1.0
	for i := 0; i < 48; i++ {
		mid := (lo + hi) / 2
		if r, g, bl := LinearRGBFromOKLab(L, mid*a, mid*b); inSRGBGamut(r, g, bl) {
			red, green, blue = r, g, bl
			lo = mid
		} else {
			hi = mid
		}
	}
	return red, green, blue
}
