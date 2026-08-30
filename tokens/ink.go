// Brand ink: the colour a role reads in when it is drawn ON a surface
// rather than filling one.
//
// A pinned base and an ink are two jobs, and the derivation solves only the
// first. FromSeed pins each role's solid fill and then chooses
// the ink that reads over that fill (see onColour), so every pin carries a
// measured guarantee about text laid on TOP of it — and none at all about
// itself laid on the page. For six of the seven pinned roles that gap
// costs nothing, because their bases are realized at fixed perceptual
// depths (lightPinTone, statusPinTone, darkPinTone): whatever the seed, a
// secondary, tertiary or status pin measures 5.94:1 or better over the
// light paper and 10.99:1 or better over the dark one. The light primary
// base is the one exception in the whole palette, and it is an exception by
// design — it is the brand colour itself, at the brand's own CIELAB depth
// (see liftSeed), so whether it reads over the paper is a property of the
// seed and of nothing else.
//
// So the family this file gates has exactly one member that can fail, and
// naming it is worth more than counting it: the light primary pin, used as
// an ink. Over the seed sweep 280 of 414 light schemes put that pin under
// the 4.5:1 text floor against their own paper, bottoming out at 1.01:1,
// and 208 of 414 put it under the 3:1 graphic floor. The canonical seed
// #6750A4 sits at L* 51 and measures 5.94:1; a pastel accent of the kind a
// dark-scheme palette publishes sits near L* 73 and puts a 1.95:1 link on a
// near-white page.
//
// # Why the gate is here and not in the derivation
//
// Deepening a pale seed on the way in would fix the ink and break the
// palette. liftSeed is a projection — every colour it returns is a colour
// it leaves alone — and that is what lets a whole palette be rebuilt from
// the one colour it publishes, which is the guarantee a kept brand is
// stored under. A derivation that moved the brand's own colour to suit one
// consumer would no longer reproduce itself from its own output. The pin
// is therefore right as it stands, and it is the *use* of a fill colour as
// an ink that has to be measured.
//
// # Why the pin stands when it reads
//
// [ColorTokens.InkOn] answers the pin while the pin clears its floor and
// walks the role's ramp only when it does not. Always walking would be the
// simpler rule and it is the wrong one twice over: a brand that reads is
// entitled to be its own colour — the same reasoning onColour applies to the
// ink over a base — and a rule that moved a pairing already clearing its
// floor would move every downstream golden for nothing. The canonical seed
// clears in both schemes and both derivations, so this gate is a no-op on
// the palette every stored image is rendered from.
//
// What a walk answers is [ColorTokens.MarkOn]'s rung — the ramp step
// nearest the mid-value 500 that clears the floor over the ground. The
// ramp is realized at fixed depths, so the answer is unmistakably the brand
// hue and never too close to its ground: over the sweep, both schemes, both
// derivations and all four storeys a brand ink can be drawn on, the worst
// pairing any seed produces measures 4.50:1 where the floor is 4.5 and
// 3.01:1 where it is 3.
package tokens

import (
	stdcolor "image/color"

	"github.com/vibrantgio/theme/color"
)

// The two floors an ink is gated at, named for the consumers that draw
// with them. They are the derivation's own numbers under exported names
// rather than a second spelling of them: TextFloor is the ratio FromSeed
// holds every on-colour to, and GraphicFloor the one every status mark is
// chosen against.
const (
	// TextFloor is WCAG 1.4.3 AA for body text, 4.5:1 — what a link, a
	// label or any run of words owes the surface it is set on.
	TextFloor = onFloor
	// GraphicFloor is WCAG 1.4.11 non-text contrast, 3:1 — what a rule, a
	// bar, a tick or any other mark that carries meaning without being
	// text owes its ground.
	GraphicFloor = graphicFloor
)

// InkOn returns the colour role reads in when it is drawn on ground: the
// role's pinned base while that base clears floor against ground, and
// otherwise the rung [ColorTokens.MarkOn] answers.
//
// It is what a consumer wants wherever a brand colour is the ink rather
// than the fill — a link in a paragraph, a blockquote's bar, a task list's
// tick, an active tab's underline — and passing the ground rather than
// assuming one is the whole of it: the same role reads in different
// colours on the paper, on a card and on its own fill, and only the caller
// knows which it is drawing on. Pass [TextFloor] for words and
// [GraphicFloor] for a mark.
//
// RoleNeutral has no pinned base and panics, as it does everywhere else a
// pin is asked for. A neutral ink over the paper is the Text pin, which is
// derived against the Background pin already.
func (t ColorTokens) InkOn(role Role, ground stdcolor.NRGBA, floor float64) stdcolor.NRGBA {
	pin := t.pinFor(role) // validates role
	if color.ContrastRatio(pin, ground) >= floor {
		return pin
	}
	return t.MarkOn(role, ground, floor)
}
