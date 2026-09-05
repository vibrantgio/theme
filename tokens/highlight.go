// The reserved highlighter: the fill that marks the content the user was
// brought to.
//
// It is not a colour role and reports no status, so its colour is reserved
// outside the role table — no status hue may serve it, and it does not
// rotate with the seed the accent family does. Every seed therefore
// derives the same two fills, #d9d6ad light and #322f09 dark, in the
// default derivation and the increased-contrast one alike.
//
// The hue is yellow, the colour a marker is: a reader who has asked to be
// shown where something is expects the mark a highlighter pen leaves, and
// no other hue is read that way. Yellow is free to be it because no status
// role stands there — the four anchors are error 28.7°, warning 64.05°,
// success 144.2° and info 248.8° (seed.go), and the run between the
// warning's orange and the success's green is 80.15° of hue with nothing
// in it. Within that run the reserved hue is the midpoint, 104.125°, the
// point furthest from either neighbour; Material Yellow 500 #FFEB3B, the
// canonical anchor of the yellow family, measures 102.50°, so the
// derivation lands 1.63° from the yellow a palette would have named.
//
// What the reservation is worth is a distance, and the distance is
// measured: over the shared seed sweep, in both schemes of both
// derivations, on every level, the realized fill stands at least 37.64°
// from every status container's hue and every status pin's, and 0.0359
// away from the nearest status container in OKLab. The two closest status
// containers, meanwhile, come to 30.18° and 0.0286 of each other, so the
// highlight is further from every status colour than two status colours
// are from each other.
//
// The chroma is containerChroma: the highlight is a fill of the same
// construction as a tonal container, so it differs from the four status
// fills in hue and in nothing else, which is the arrangement in which four
// status fills read as four (containers.go). sRGB holds the dial at this
// hue at both depths the fill is realized at — 0.1850 at the light
// step-300 depth (L* 85) and 0.0650 at the dark one (L* 19), against a
// dial of 0.055 — so the reservation never costs the fill its hue. The
// dark depth is the binding case for the dial across the whole palette:
// no hue the derivation asks for holds less there.
//
// The floor is ContainerFloor and not a text or graphic floor: a highlight
// is a field, not a mark on one, and what a field owes the surface it
// marks is that its edge be findable. So the fill deepens off the fixed
// step exactly as a container does, wherever the surface it marks has
// walked into it: over the sweep the seam measures 1.307:1 at worst and
// 1.705:1 at loudest, and the neutral Text token clears TextFloor over
// every fill the walk returns, 8.04:1 at worst.
//
// One strength. The fill is one colour per scheme; a second, stronger one
// would owe its own floor against the first and its own text floor over
// it, and this derivation carries neither.
package tokens

import (
	stdcolor "image/color"

	"github.com/vibrantgio/theme/color"
)

// highlightHue is the reserved highlighter's OKLCh hue: the midpoint of
// the run of the hue circle the yellows occupy, which the status anchors
// leave open between the warning and the success. The internal test
// recomputes it from the anchors themselves.
const highlightHue = 104.125

// HighlightOn returns the reserved highlighter's fill for content standing
// on surface: [ColorTokens.Highlight]'s realization while that realization
// is visibly a different surface from the one it marks, and otherwise the
// first deeper one that is.
//
// It is [ColorTokens.ContainerOn]'s walk at the reserved hue, and it is
// needed for the same reason: the fixed step the fill is realized at is one
// the elevation levels walk through, so a highlight drawn at it on a raised
// surface is not subtle but absent, and the marking silently stops marking.
// Content on the paper takes the resolved [ColorTokens.Highlight] field;
// content on any other level passes that level's fill here.
//
// A highlight marks content and reports no status, so nothing about the
// colour it answers with may be read as one; see the file header for the
// reservation and its measured distances.
func (t ColorTokens) HighlightOn(surface stdcolor.NRGBA) stdcolor.NRGBA {
	return fillToFloor(t.highlightAt, surface)
}

// highlightAt realizes the highlighter at one ramp step: the neutral step's
// own depth at the reserved hue and the container dial's chroma. The depth
// comes off the neutral ramp because the fill carries no role — there is no
// ramp of its own to read a step from — and the neutral steps are the depths
// every scheme's own surfaces are placed against.
func (t ColorTokens) highlightAt(step int) stdcolor.NRGBA {
	tone, _, _ := color.LabFromNRGBA(t.Ramps.Neutral.Step(step))
	return color.NRGBAFromToneChromaHue(tone, containerChroma, highlightHue)
}
