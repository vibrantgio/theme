// What a measured drop shadow is: its coverage and the geometry it is spread
// over.
package tokens

import "image/color"

// DropShadow is one reading of a shadow the platform draws: the peak coverage
// at the edge of the rectangle the ramp is spread from, how far that ramp
// carries past that rectangle, and how far below the shape the rectangle
// sits.
//
// The three are read off one capture together and are spent together — a
// coverage fitted at one reach is not the same coverage at another — which is
// why they travel as one value rather than as a colour here and two lengths
// in whatever draws it. A reader that has the shadow has the whole reading.
//
// Reach and Offset are in dp, as every length in this package is. The zero
// value draws nothing: a thing that casts no shadow answers it.
type DropShadow struct {
	// Peak is the coverage at the sunk rectangle's own edge. It is a
	// coverage and not a fill: what it lands as is this alpha over whatever
	// the shadow falls on.
	Peak color.NRGBA
	// Reach is how far the ramp carries past the rectangle, in dp, falling
	// linearly from Peak to nothing.
	Reach float32
	// Offset is how far below the shape that rectangle sits, in dp. Zero
	// centres the shadow on the shape; a positive value is what makes a
	// shadow heavier under a thing than over it.
	Offset float32
}

// WithPeak returns s drawn at another coverage — a surface fading in or out,
// or a control that is not offering itself — with the geometry it was
// measured at kept.
func (s DropShadow) WithPeak(c color.NRGBA) DropShadow {
	s.Peak = c
	return s
}

// ShadowCurveLen is the number of pixels of distance a [ShadowCurve] holds.
const ShadowCurveLen = 12

// ShadowCurve is a shadow read pixel by pixel: the coverage at each distance
// out from the edge of the rectangle it is cast from, where a [DropShadow]
// holds one peak and a linear ramp. The platform draws some shadows along a
// curve that no linear ramp lands, and one that differs per scheme.
//
// Alpha[d] is black at the coverage the shadow has d pixels out from that
// rectangle's own edge pixel, which is index 0. Entries past the shadow's
// end are zero. Offset is how far below the shape the rectangle sits, in dp.
//
// At a corner the shadow is the product of the two sides' curves:
// Alpha[dx]·Alpha[dy] / (2·Alpha[0]), which [ShadowCurve.CornerAlpha] gives.
//
// The zero value draws nothing.
type ShadowCurve struct {
	Alpha  [ShadowCurveLen]color.NRGBA
	Offset float32
}

// Reach is the number of pixels the shadow covers out from its rectangle,
// index 0 included: the index of the last non-zero entry plus one.
func (s ShadowCurve) Reach() int {
	for d := len(s.Alpha) - 1; d >= 0; d-- {
		if s.Alpha[d].A != 0 {
			return d + 1
		}
	}
	return 0
}

// CornerAlpha is the shadow's coverage, out of 255, at dx pixels out from
// the rectangle's side edge and dy pixels out from its top or bottom edge,
// both counted as Alpha counts. A distance past the curve has no shadow.
func (s ShadowCurve) CornerAlpha(dx, dy int) float64 {
	if dx < 0 || dy < 0 || dx >= len(s.Alpha) || dy >= len(s.Alpha) || s.Alpha[0].A == 0 {
		return 0
	}
	return float64(s.Alpha[dx].A) * float64(s.Alpha[dy].A) / (2 * float64(s.Alpha[0].A))
}
