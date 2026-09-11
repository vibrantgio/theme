// The highlight: the fill that marks the content the user was brought to.
//
// One yellow, rgb(255, 208, 0), the same in both schemes, laid over the
// surface it marks at less than full strength so that surface shows through
// it. The colour and the coverage are read from Obsidian's stylesheet, which
// lays that yellow over the page at 40% and never overrides it per theme.
//
// The answer is an opaque colour, because the token contract is that no fill
// is translucent: the blend happens here and the caller paints what it is
// handed. It is blended in the eight-bit sRGB space ([color.Flatten]'s
// space) because the reference is what a browser puts on the screen and a
// browser composites there — over white the yellow lands on #ffec99 and over
// #1e1e1e on #786512, which is the reference as it reads.
//
// The text the highlight covers is not repainted, and the coverage is not
// tuned to a text floor: the marked words keep the colour they had, and what
// the composite costs them is measured rather than corrected for.
//
// The hue is yellow because that is the colour a marker leaves, and yellow is
// free to carry it because no status stands there: the four status anchors
// are error 28.7°, warning 64.05°, success 144.2° and info 248.8°
// (seed.go), and the run between the warning's orange and the success's green
// is 80.15° of hue with nothing in it. The marker's yellow falls inside that
// run; the internal test recomputes the run from the anchors themselves.
//
// One yellow, two strengths. A match and the current match differ in the
// coverage and in nothing else — the reader is stepping through matches, not
// looking at two kinds of thing.
package tokens

import (
	stdcolor "image/color"
	"math"
)

// markerYellow is the colour every highlight is laid on in: one yellow for
// both schemes, read from Obsidian's stylesheet.
var markerYellow = stdcolor.NRGBA{R: 0xff, G: 0xd0, B: 0x00, A: 0xff}

// The coverages the yellow is laid on at. highlightCoverage is Obsidian's
// own; currentMatchCoverage is the single dial that makes the current match
// the same yellow laid on more strongly.
const (
	highlightCoverage    = 0.40
	currentMatchCoverage = 0.70
)

// HighlightOn returns the highlight for content standing on surface: the
// marker yellow laid over that surface at 40%, opaque.
//
// It is [ColorTokens.ContainerOn]'s shape — the caller names what its content
// stands on and gets what to paint — because the surface shows through the
// yellow, so the fill is only right for the surface it was asked for. Content
// on level 0 may take the resolved [ColorTokens.Highlight] field instead.
//
// A highlight marks content and reports no status, so nothing about the
// colour may be read as one; see the file header for the reservation.
func (t ColorTokens) HighlightOn(surface stdcolor.NRGBA) stdcolor.NRGBA {
	return yellowOver(highlightCoverage, surface)
}

// CurrentMatchOn returns the fill for the current match among many — the one
// the reader is standing on — for content standing on surface: the same
// yellow [ColorTokens.HighlightOn] lays, over the same surface, at 70%
// instead of 40%. The coverage is the only difference between the two.
func (t ColorTokens) CurrentMatchOn(surface stdcolor.NRGBA) stdcolor.NRGBA {
	return yellowOver(currentMatchCoverage, surface)
}

// yellowOver lays the marker yellow over surface at coverage and returns the
// opaque colour that lands. surface is a surface — what is already on the
// screen — so its alpha is not read. See the file header for why the blend is
// taken in the eight-bit sRGB space.
func yellowOver(coverage float64, surface stdcolor.NRGBA) stdcolor.NRGBA {
	mix := func(src, dst uint8) uint8 {
		return uint8(math.Round(coverage*float64(src) + (1-coverage)*float64(dst)))
	}
	return stdcolor.NRGBA{
		R: mix(markerYellow.R, surface.R),
		G: mix(markerYellow.G, surface.G),
		B: mix(markerYellow.B, surface.B),
		A: 0xff,
	}
}
