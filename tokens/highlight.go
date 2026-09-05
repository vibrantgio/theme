// The reserved highlighter: the fill that marks the content the user was
// brought to.
//
// It is not a colour role and reports no status, so its colour is reserved
// outside the role table — no status hue may serve it, and it does not
// rotate with the seed the accent family does. Every seed therefore
// derives the same two fills, #e7d700 light and #4e4800 dark, in the
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
// derivations, on every level, the realized fill stands at least 37.60°
// from every status container's hue and every status pin's, and 0.0397
// away from the nearest status container in OKLab. The two closest status
// containers, meanwhile, come to 30.18° and 0.0286 of each other, so the
// highlight is further from every status colour than two status colours
// are from each other.
//
// The chroma is the marker's own and not the container dial: a status fill
// is one of four that must read as four, so it is held to a shared dial,
// while a highlighter is one colour that must read as a marker — so the
// fill asks for more chroma than sRGB holds anywhere (highlightChroma) and
// takes what the gamut gives it at the depth it is realized at. That is
// 0.1845 at the light depth and 0.0850 at the dark one, against the 0.055
// the containers are dialled to. The two schemes cannot share a number:
// sRGB holds far less chroma in the dark half of the ramp than in the pale
// half, so one dial high enough for the light fill is unreachable in the
// dark and one low enough for the dark fill leaves the light fill a khaki.
//
// The depth is not the container step in both schemes, because what the
// depth has to deliver is a yellow. A light scheme's step-300 depth
// (L* 84.91) is where a marker sits anyway — pale enough to read as text on
// paper laid over with a pen, and it holds chroma to spare. A dark scheme's
// step-300 depth (L* 18.94) holds only 0.0650, which renders as an olive;
// its step-400 depth (L* 30.16) holds 0.0850 and renders as a yellow, so
// the fill is realized there. markerChroma is the line between those two
// readings, and the depth the fill takes is the shallowest step that holds
// it — a derivation, not a table: it answers 300 in a light scheme and 400
// in a dark one because that is where the gamut puts the yellow.
//
// The floor is ContainerFloor and not a text or graphic floor: a highlight
// is a field, not a mark on one, and what a field owes the surface it
// marks is that its edge be findable. So the fill moves off that depth
// wherever the surface it marks has walked into it — see
// [ColorTokens.HighlightOn] for the two directions it may move in and why:
// over the sweep the seam measures 1.317:1 at worst and 1.957:1 at best,
// and the neutral Text token clears TextFloor over every fill the walk
// returns, 5.157:1 at worst.
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

// highlightChroma is the chroma asked for at the reserved hue: more than
// sRGB holds at any depth, so the realization is clipped to the gamut and
// each scheme's fill carries as much yellow as its own depth can hold.
const highlightChroma = 0.4

// markerChroma is the least chroma the fill's depth must hold to read as
// the yellow a marker is: measured from renders of text over the fill, a
// dark scheme's step-300 depth holds 0.0650 and reads olive, its step-400
// depth holds 0.0850 and reads yellow.
const markerChroma = 0.08

// HighlightOn returns the reserved highlighter's fill for content standing
// on surface: [ColorTokens.Highlight]'s realization while that realization
// is visibly a different surface from the one it marks, and otherwise the
// first depth in its walk that is.
//
// The walk is needed for the reason [ColorTokens.ContainerOn]'s is: the
// fixed step the fill is realized at is one the elevation levels walk
// through, so a highlight drawn at it on a raised surface is not subtle but
// absent, and the marking silently stops marking. Content on the paper takes
// the resolved [ColorTokens.Highlight] field; content on any other level
// passes that level's fill here.
//
// Where a container's walk only deepens, this one may also step back toward
// the paper, because the two fills owe different things. A container is the
// surface its own content stands on and may be as deep as it likes; a
// highlight covers content the scheme's own body text is already set in, so
// it is bounded above by the depth at which that text stops clearing
// [TextFloor] — and in a dark scheme the level-3 surface stands at the very
// depth the fill deepens to, leaving nothing readable above it. So the
// candidates are the marker depth, then each deeper step whose fill the
// body text still clears, then the steps back toward the paper, which
// separate in the other direction; the first that clears [ContainerFloor]
// wins, and if none does, the one that separates most.
//
// A highlight marks content and reports no status, so nothing about the
// colour it answers with may be read as one; see the file header for the
// reservation and its measured distances.
func (t ColorTokens) HighlightOn(surface stdcolor.NRGBA) stdcolor.NRGBA {
	best, bestAt := -1.0, t.highlightStep()
	for _, step := range t.highlightWalk() {
		fill := t.highlightAt(step)
		got := color.ContrastRatio(fill, surface)
		if got >= ContainerFloor {
			return fill
		}
		if got > best {
			best, bestAt = got, step
		}
	}
	return t.highlightAt(bestAt)
}

// highlightWalk is the order the depths are tried in: the marker depth, the
// deeper steps whose fill the scheme's body text still clears TextFloor
// over, then the steps between the marker depth and containerStep, shallowest
// last. See [ColorTokens.HighlightOn].
func (t ColorTokens) highlightWalk() []int {
	from := t.highlightStep()
	steps := []int{}
	for step := from; step <= 900; step += 100 {
		if color.ContrastRatio(t.Text, t.highlightAt(step)) < TextFloor {
			break
		}
		steps = append(steps, step)
	}
	for step := from - 100; step >= containerStep; step -= 100 {
		steps = append(steps, step)
	}
	return steps
}

// highlightStep is the marker depth: the first step at or past containerStep
// whose realization holds markerChroma. See the file header for why the
// depth is asked to hold a chroma at all.
func (t ColorTokens) highlightStep() int {
	for step := containerStep; step <= 900; step += 100 {
		if _, chroma, _ := color.OKLChFromNRGBA(t.highlightAt(step)); chroma >= markerChroma {
			return step
		}
	}
	return containerStep
}

// highlightAt realizes the highlighter at one ramp step: the neutral step's
// own depth at the reserved hue, with all the chroma sRGB holds there. The depth
// comes off the neutral ramp because the fill carries no role — there is no
// ramp of its own to read a step from — and the neutral steps are the depths
// every scheme's own surfaces are placed against.
func (t ColorTokens) highlightAt(step int) stdcolor.NRGBA {
	tone, _, _ := color.LabFromNRGBA(t.Ramps.Neutral.Step(step))
	return color.NRGBAFromToneChromaHue(tone, highlightChroma, highlightHue)
}
