// The raise: what a thing standing on a surface fills with, and the seam it
// draws where the fill cannot say it on its own.
//
// A raise is WALKED from the surface a thing stands on, never read off a
// table of levels. Whatever a thing stands on, raised means one step lighter
// than that: a card on the content, a field on that card, a card on a modal.
// The walk is what makes the grammar hold wherever a thing is placed —
// reaching for a fixed level instead holds only while every surface in a
// window happens to be the one that level was chosen against, and moving the
// host one level then puts the raised thing two levels off it.
//
// One step is the surface band's own first interval (elevation.go), and the
// scheme's ceiling stops the walk: the band's own top where the band climbs
// away from its 100 stop, and the tonal axis where it does not. So the walk
// runs out — at the top of the dark band, and at white in the light scheme,
// where the content pin stands one band step under the axis and its first
// raise IS white.
//
// Where it runs out the raise is still told, by a seam at its own edge
// instead of by its fill: a raise never vanishes. [Raise.Seamed] is that
// obligation, and [Raise.Seam] the hairline it is discharged with, drawn
// once, by the raised surface at its own edge — an inset object needs no
// seam on the side the surface beneath it shows around.
//
// A caller that already draws an edge of its own has discharged the seam
// with it: a text field's resting border is a 3:1 mark around exactly this
// pairing, and a second hairline inside it would be two lines saying one
// thing.
package tokens

import (
	"image/color"
	"math"

	vgcolor "github.com/vibrantgio/theme/color"
)

// RaiseFloor is the least contrast a raise owes the surface beneath it for
// its FILL to be doing the telling: below it the step is there but nobody
// can find it, and the raise owes a seam instead.
//
// 1.09:1 is measured rather than picked. Over the seed sweep — both schemes,
// both contrast variants, every named level walked and every walk taken
// again off its own answer, eight deep — a raise's separation from the
// surface beneath takes four values and no others: 1.0000 where the walk
// clamps on a ceiling it has already reached, 1.0483 and 1.0816 where the
// ceiling is nearer than a step away and the walk lands a partial one, and
// 1.1079 upward where a FULL band step separates the two. The threshold goes
// in the empty stretch between 1.0816 and 1.1079, so a raise is told by its
// fill when it got a whole step and by its seam when it did not, and no
// palette sits near enough to it to flip on a rounding.
//
// A partial step is rejected on the platform's own evidence rather than on
// taste: a macOS dark window parts its content from its furniture by 1.48
// L*, and that boundary is drawn with a seam. A 1.99 L* raise is the same
// order of whisper and owes the same hairline.
//
// It is far under the 1.25:1 that [ContainerFloor], [StateFloor] and the
// focus ring's border separation share, and deliberately: those three ask
// when a fill laid ON a surface stops being that surface, and are set where
// this design's own pairings already sit. A raise is the surface itself
// moving, and every raise the elevation carries is quieter than any of the
// state washes drawn on it. So this floor is measured off the levels' own
// population and not borrowed from theirs.
const RaiseFloor = 1.09

// SeamRatio is how far a seam stands from the fills it parts: 1.51:1, and it
// is a MEASUREMENT of the platform rather than a floor anything has to
// clear. The platform draws this hairline and draws it quietly — Voice Memos
// outlines its inset panel at #3A3A3A on a #1B1B1B panel, 1.514:1.
//
// It is deliberately NOT the 3:1 graphic floor a mark derives to: a 3:1 mark
// carries meaning by itself and owes WCAG 1.4.11, while a seam only says
// where one region ends and the next begins, read alongside the corner
// radius and the shadow that say the same thing. At 3:1 every card in a
// light window would wear a line louder than anything the platform draws.
const SeamRatio = 1.51

// Raise is what a thing standing on a surface is drawn with: the surface one
// step up, and the seam that tells the raise where that step cannot.
//
// Seam is filled in whether or not it is owed, so a caller that wants the
// hairline for a reason of its own has it; Seamed is the obligation.
type Raise struct {
	// Fill is the surface one step above the one the thing stands on,
	// clamped at the scheme's ceiling.
	Fill color.NRGBA
	// Seam is the hairline the raise draws at its own edge: findable
	// against both fills, in either scheme.
	Seam color.NRGBA
	// Seamed reports that Fill alone does not clear [RaiseFloor] against
	// the surface beneath, so the raise is owed its Seam.
	Seamed bool
}

// RaisedOn returns the raise of a thing standing on surface: the surface one
// step up, whether that step can be seen, and the seam it owes if it cannot.
//
// It is [ColorTokens.ContainerOn]'s and [ColorTokens.HighlightOn]'s shape —
// the caller names what it stands on and gets what to paint — and it is the
// only way to ask for a raised fill. surface is a surface of this scheme:
// a level's own fill, or another raise's.
//
// The fill is realized at the Background pin's hue and chroma, the way every
// level is, so a raise carries whatever tint the palette carries and none of
// its own; a depth that coincides with one of the band's own steps answers
// with that step verbatim, which is what makes the dark scheme's first raise
// Neutral 200 byte-for-byte rather than approximately.
func (t ColorTokens) RaisedOn(surface color.NRGBA) Raise {
	band, tone := t.surfaceBand()
	standing, _, _ := vgcolor.LabFromNRGBA(surface)
	target := math.Min(standing+raiseStep(tone), ceiling(tone))
	if target < standing {
		target = standing // the surface is already past the ceiling
	}
	fill := t.realizeSurface(target, band, tone)
	return Raise{
		Fill:   fill,
		Seam:   t.seamBetween(surface, fill),
		Seamed: vgcolor.ContrastRatio(fill, surface) < RaiseFloor,
	}
}

// seamBetween derives the hairline parting two flush surfaces: [SeamRatio]
// away from both of them, toward the scheme's own foreground, realized at
// the raised fill's own hue and chroma.
//
// Two things are derived and neither names a scheme. The DISTANCE is solved
// in the luminance a contrast ratio is taken in, against whichever of the
// two fills is the binding one — the further target of the two — so the
// hairline is findable against both and not only against the one it happens
// to be drawn over. The DIRECTION is toward the scheme's own foreground: a
// dark scheme's seam is lighter than the surfaces it parts, as the platform
// draws it, and a light scheme's is darker, which is the only direction a
// light seam has room in when the surface above it is already white.
func (t ColorTokens) seamBetween(below, above color.NRGBA) color.NRGBA {
	yb := vgcolor.RelativeLuminance(below)
	ya := vgcolor.RelativeLuminance(above)
	toward := 1.0 // the scheme's foreground is lighter than its content
	target := math.Max(SeamRatio*(yb+0.05), SeamRatio*(ya+0.05)) - 0.05
	if fg, bg := lightnessOf(t.Text), lightnessOf(t.Background); fg < bg {
		toward = -1
		target = math.Min((yb+0.05)/SeamRatio, (ya+0.05)/SeamRatio) - 0.05
	}
	_, chroma, hue := vgcolor.OKLChFromNRGBA(above)
	tone := toneOfLuminance(math.Min(math.Max(target, 0), 1))
	seam := vgcolor.NRGBAFromToneChromaHue(tone, chroma, hue)
	// The realization is eight-bit, so it can land a hair inside the target
	// it was solved for. Findable is the promise, so the derivation confirms
	// against both fills and steps once more toward the foreground rather
	// than shipping a seam that misses by a rounding.
	for i := 0; i < 8; i++ {
		if vgcolor.ContrastRatio(seam, below) >= SeamRatio && vgcolor.ContrastRatio(seam, above) >= SeamRatio {
			break
		}
		if tone += toward * 0.5; tone < 0 || tone > axisTop {
			break
		}
		seam = vgcolor.NRGBAFromToneChromaHue(tone, chroma, hue)
	}
	return seam
}

// lightnessOf is a colour's CIELAB L*, which is what "toward the
// foreground" compares: a seam's direction is a question about lightness and
// nothing else.
func lightnessOf(c color.NRGBA) float64 {
	l, _, _ := vgcolor.LabFromNRGBA(c)
	return l
}

// toneOfLuminance is the CIELAB lightness of a relative luminance — the
// inverse of the Y a WCAG contrast ratio is taken on. A distance stated as a
// ratio is solved in Y; the toolkit realizes a colour from a tone, a chroma
// and a hue; this is the one step between them.
func toneOfLuminance(y float64) float64 {
	if y <= 216.0/24389.0 {
		return y * 24389.0 / 27.0
	}
	return 116*math.Cbrt(y) - 16
}
