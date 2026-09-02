// The reserved highlighter: the wash that marks the content the user was
// brought to.
//
// It is not a colour role and reports no status, so its colour is reserved
// outside the role table — no status hue may serve it, and it does not
// rotate with the seed the accent family does. Every seed therefore
// derives the same two washes, #e6cbee light and #3b2641 dark, in the
// default derivation and the increased-contrast one alike.
//
// The hue is the midpoint of the widest arc the four status anchors leave
// open on the OKLCh hue circle. Sorted, the anchors are error 28.7°,
// warning 84.9°, success 144.2° and info 248.8° (seed.go); the runs
// between them measure 56.2°, 59.3°, 104.6° and 139.9°, and the widest is
// info's to error's, whose midpoint is 318.75°. Warning is the collision
// the reservation exists to avoid — a yellow wash is read as a warning,
// and warning's anchor is yellow — and warning's bend toward orange with
// depth, as far as 51.9°, only widens the run this hue sits in.
//
// What the reservation is worth is a distance, and the distance is
// measured: over the shared seed sweep, in both schemes of both
// derivations, on every level, the realized wash stands at least 64.01°
// from every status container's hue and every status pin's, and 0.0582
// away from the nearest status container in OKLab. The two closest status
// containers, meanwhile, come to 48.33° and 0.0453 of each other, so the
// highlight is a third further again from every status colour than two
// status colours are from each other.
//
// The chroma is containerChroma: the highlight is a wash of the same
// construction as a tonal container, so it differs from the four status
// washes in hue and in nothing else, which is the arrangement in which
// four status grounds read as four (containers.go). sRGB holds the dial at
// this hue with headroom at both depths the wash is realized at — 0.0935
// at the light step-300 depth (L* 85) and 0.1575 at the dark one (L* 19),
// against a dial of 0.055 — so the reservation never costs the wash its
// hue.
//
// The floor is ContainerFloor and not a text or graphic floor: a highlight
// is a field, not a mark on one, and what a field owes the surface it
// marks is that its edge be findable. So the wash deepens off the fixed
// step exactly as a container does, wherever the surface it marks has
// walked into it: over the sweep the seam measures 1.302:1 at worst and
// 1.713:1 at loudest, and the neutral Text token clears TextFloor over
// every wash the walk returns, 8.00:1 at worst.
//
// One strength. The wash is one colour per scheme; a second, stronger one
// would owe its own floor against the first and its own text floor over
// it, and this derivation carries neither.
package tokens

import (
	stdcolor "image/color"

	"github.com/vibrantgio/theme/color"
)

// highlightHue is the reserved highlighter's OKLCh hue: the midpoint of
// the widest arc the four status anchors leave open, info's to error's.
// The internal test recomputes it from the anchors themselves.
const highlightHue = 318.75

// HighlightOn returns the reserved highlighter's wash for content standing
// on surface: [ColorTokens.Highlight]'s realization while that realization
// is visibly a different surface from the one it marks, and otherwise the
// first deeper one that is.
//
// It is [ColorTokens.ContainerOn]'s walk at the reserved hue, and it is
// needed for the same reason: the fixed step the wash is realized at is
// one the elevation ladder walks through, so a highlight drawn at it on a
// raised surface is not subtle but absent, and the marking silently stops
// marking. Content on the paper takes the resolved [ColorTokens.Highlight]
// field; content on any other level passes that level's fill here.
//
// A highlight marks content and reports no status, so nothing about the
// colour it answers with may be read as one; see the file header for the
// reservation and its measured distances.
func (t ColorTokens) HighlightOn(surface stdcolor.NRGBA) stdcolor.NRGBA {
	return fillToFloor(t.highlightAt, surface)
}

// highlightAt realizes the highlighter at one ramp step: the neutral rung's
// own depth at the reserved hue and the container dial's chroma. The depth
// comes off the neutral ramp because the wash carries no role — there is no
// ramp of its own to read a rung from — and the neutral rungs are the
// depths every scheme's own surfaces are placed against.
func (t ColorTokens) highlightAt(step int) stdcolor.NRGBA {
	tone, _, _ := color.LabFromNRGBA(t.Ramps.Neutral.Step(step))
	return color.NRGBAFromToneChromaHue(tone, containerChroma, highlightHue)
}
