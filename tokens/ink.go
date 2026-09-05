// The brand foreground: the colour a role reads in when it is drawn ON a
// surface rather than filling one.
//
// A pinned base and a foreground are two jobs, and the derivation solves only
// the first. FromSeed pins each role's solid fill and then chooses the
// foreground that reads over that fill (see onColour), so every pin carries a
// measured guarantee about text laid on TOP of it — and none at all about
// itself laid on the page. For six of the seven pinned roles that gap costs
// nothing, because their bases are realized at fixed perceptual depths
// (lightPinTone, statusPinTone, darkPinTone): whatever the seed, a secondary,
// tertiary or status pin measures 5.94:1 or better over the light paper and
// 10.99:1 or better over the dark one. The light primary base is the one
// exception in the whole palette, and it is an exception by design — it is the
// brand colour itself, at the brand's own CIELAB depth (see liftSeed), so
// whether it reads over the paper is a property of the seed and of nothing
// else.
//
// So the family this file gates has exactly one member that can fail, and
// naming it is worth more than counting it: the light primary pin, used as a
// foreground. Over the seed sweep 280 of 414 light schemes put that pin under
// the 4.5:1 text floor against their own paper, bottoming out at 1.01:1, and
// 208 of 414 put it under the 3:1 graphic floor. The canonical seed #6750A4
// sits at L* 51 and measures 5.94:1; a pastel accent of the kind a dark-scheme
// palette publishes sits near L* 73 and puts a 1.95:1 link on a near-white
// page.
//
// # Why the gate is here and not in the derivation
//
// Deepening a pale seed on the way in would fix the foreground and break the
// palette. liftSeed is a projection — every colour it returns is a colour
// it leaves alone — and that is what lets a whole palette be rebuilt from
// the one colour it publishes, which is the guarantee a kept brand is
// stored under. A derivation that moved the brand's own colour to suit one
// consumer would no longer reproduce itself from its own output. The pin
// is therefore right as it stands, and it is the *use* of a fill colour as
// a foreground that has to be measured.
//
// # Why the pin stands when it reads
//
// [ColorTokens.InkOn] answers the pin while the pin clears its floor and walks
// the role's ramp only when it does not. Always walking would be the simpler
// rule and it is the wrong one twice over: a brand that reads is entitled to
// be its own colour — the same reasoning onColour applies to the foreground
// over a base — and a rule that moved a pairing already clearing its floor
// would move every downstream golden for nothing. The canonical seed clears in
// both schemes and both derivations, so this gate is a no-op on the palette
// every stored image is rendered from.
//
// What a walk answers is [ColorTokens.MarkOn]'s step — the ramp step nearest
// the mid-value 500 that clears the floor over the surface. The ramp is
// realized at fixed depths, so the answer is unmistakably the brand hue and
// never too close to the surface beneath: over the sweep, both schemes, both
// derivations and all four levels a brand foreground can be drawn on, the
// worst pairing any seed produces measures 4.50:1 where the floor is 4.5 and
// 3.01:1 where it is 3.
package tokens

import (
	stdcolor "image/color"

	"github.com/vibrantgio/theme/color"
)

// The two floors a foreground is gated at, named for the consumers that draw
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
	// text owes the surface it is drawn on.
	GraphicFloor = graphicFloor
)

// InkOn returns the colour role reads in when it is drawn on the given
// surface: the role's pinned base while that base clears floor against it,
// and otherwise the step [ColorTokens.MarkOn] answers.
//
// It is what a consumer wants wherever a brand colour is the foreground
// rather than the fill — a link in a paragraph, a blockquote's bar, a task
// list's tick, an active tab's underline — and passing that surface rather
// than assuming one is the whole of it: the same role reads in different
// colours on the paper, on a card and on its own fill, and only the caller
// knows which it is drawing on. Pass [TextFloor] for words and
// [GraphicFloor] for a mark.
//
// RoleNeutral has no pinned base and panics, as it does everywhere else a pin
// is asked for. A neutral foreground over the paper is the Text pin, which is
// derived against the Background pin already.
func (t ColorTokens) InkOn(role Role, ground stdcolor.NRGBA, floor float64) stdcolor.NRGBA {
	pin := t.pinFor(role) // validates role
	if color.ContrastRatio(pin, ground) >= floor {
		return pin
	}
	return t.MarkOn(role, ground, floor)
}

// ForegroundOn returns the colour a role's content — a word, a count, or a
// sign standing in for one — reads in over surface: [ColorTokens.InkOn] at
// [TextFloor] for the roles that carry a pinned base, and
// [ColorTokens.MarkOn]'s walk for RoleNeutral, which carries none.
//
// It is the foreground half of the one tonal recipe, and it exists as a
// function because more than one component draws content in a role's own
// hue over a fill of that same hue: a tinted button and a status badge are
// the same tint, and it is behaviour rather than colour that tells them
// apart. Two spellings of one derivation is how they drift.
//
// The floor is the text one whatever the content is. A component that says
// its word as a sign is making the same utterance at the same weight, so
// deriving the sign at [GraphicFloor] would make one component read at two
// strengths depending on which of its faces it wore.
//
// surface is whatever is ACTUALLY behind the content, which for a fill that
// walks under the pointer is the walked fill and not the resting one: a
// foreground held over a fill that moved is derived against a surface that
// is no longer there.
func (t ColorTokens) ForegroundOn(role Role, surface stdcolor.NRGBA) stdcolor.NRGBA {
	if role == RoleNeutral {
		// InkOn asks a role for its pinned base and neutral has none; the
		// walk is the whole derivation for it.
		return t.MarkOn(role, surface, TextFloor)
	}
	return t.InkOn(role, surface, TextFloor)
}
