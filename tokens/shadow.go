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
