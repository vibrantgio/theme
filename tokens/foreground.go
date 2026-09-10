// The brand foreground: the colour a role reads in when it is drawn ON a
// surface rather than filling one.
//
// A pinned base and a foreground are two jobs, and the derivation solves only
// the first. FromSeed pins each role's solid fill and then chooses the
// foreground that reads over that fill, so every pin carries a measured
// guarantee about text laid on TOP of it — and none at all about itself laid
// on the page. Under APCA's own body-text level, Lc 75, that gap is wide:
// the light pins are realized at L* 40 or 39 and the dark ones at L* 82, and
// none of those depths reaches Lc 75 over the content plane. So the pin
// stands as a foreground only where the surface happens to be far enough
// away — the light scheme's white raised fills — and walks everywhere else.
//
// Over the 414-seed sweep 346 of the light schemes put the primary pin under
// [TextFloor] against their own content and 175 put it under [GraphicFloor],
// bottoming out at |Lc| 0.00 where the seed IS the content's own colour. The
// canonical seed #6750A4 realizes at L* 41 and measures |Lc| 72.71 over the
// content; a pastel accent of the kind a dark-scheme palette publishes sits
// near L* 73 and puts an |Lc| 32.73 link on a near-white page.
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
// [ColorTokens.ForegroundOnAtFloor] answers the pin while the pin clears
// its floor and walks the role's ramp only when it does not. Always walking
// would be the simpler rule and it is the wrong one twice over: a brand
// that reads is entitled to be its own colour — the same reasoning
// [color.BestOn] applies to the foreground over a base — and a rule that
// moved a pairing already clearing its floor would move every downstream
// golden for nothing.
//
// What a walk answers is [ColorTokens.MarkOn]'s step — the ramp step nearest
// the mid-value 500 that clears the floor over the surface. The ramp is
// realized at fixed depths, so the answer is unmistakably the brand hue and
// never too close to the surface beneath: over the sweep, both schemes, both
// derivations and all four levels a brand foreground can be drawn on, the
// worst pairing any seed produces measures |Lc| 75.15 where the floor is 75
// and 45.15 where it is 45.
package tokens

import (
	stdcolor "image/color"

	"github.com/vibrantgio/theme/color"
)

// The floors a foreground is gated at, named for the consumers that draw
// with them. Each is a floor on |Lc| — [color.Magnitude] — and each is
// APCA's own published level for the job, so a consumer that imports one
// and a derivation that gates on it are reading the same number:
// [color.APCA] is the system's one contrast measure and these are the
// levels it is read against.
const (
	// TextFloor is APCA's minimum for body text, Lc 75 — what a link, a
	// label or any run of words owes the surface it is set on.
	TextFloor = onFloor
	// GraphicFloor is APCA's level for a mark that carries meaning without
	// being text, Lc 45 — what a rule, a bar, a tick or any other such
	// mark owes the surface it is drawn on.
	GraphicFloor = graphicFloor
	// IncreasedTextFloor is APCA's preferred level for body text, Lc 90:
	// the floor the increased-contrast variant asks of the same pairing.
	IncreasedTextFloor = hcOnFloor
	// IncreasedGraphicFloor is the increased step of the mark level,
	// Lc 60 — a mark under the increased-contrast variant.
	IncreasedGraphicFloor = hcGraphicFloor
)

// ForegroundOnAtFloor returns the colour role reads in when it is drawn on
// the given surface: the role's pinned base while that base clears floor
// against it, and otherwise the step [ColorTokens.MarkOn] answers.
//
// It is what a consumer wants wherever a brand colour is the foreground
// rather than the fill — a link in a paragraph, a blockquote's bar, a task
// list's tick, an active tab's underline — and passing that surface rather
// than assuming one is the whole of it: the same role reads in different
// colours on the content, on a card and on its own fill, and only the caller
// knows which it is drawing on. Pass [TextFloor] for words and
// [GraphicFloor] for a mark.
//
// RoleNeutral has no pinned base and panics, as it does everywhere else a pin
// is asked for. A neutral foreground over the content is the Text pin, which is
// derived against the Background pin already.
func (t ColorTokens) ForegroundOnAtFloor(role Role, surface stdcolor.NRGBA, floor float64) stdcolor.NRGBA {
	pin := t.pinFor(role) // validates role
	if color.Magnitude(pin, surface) >= floor {
		return pin
	}
	return t.MarkOn(role, surface, floor)
}

// ForegroundOn returns the colour a role's content — a word, a count, or a
// sign standing in for one — reads in over surface:
// [ColorTokens.ForegroundOnAtFloor] at [TextFloor] for the roles that carry
// a pinned base, and [ColorTokens.MarkOn]'s walk for RoleNeutral, which
// carries none.
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
		// ForegroundOnAtFloor asks a role for its pinned base and neutral has
		// none; the walk is the whole derivation for it.
		return t.MarkOn(role, surface, TextFloor)
	}
	return t.ForegroundOnAtFloor(role, surface, TextFloor)
}
