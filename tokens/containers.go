// Tonal containers: the ground a status role fills, and the mark read on it.
//
// A container is a role's own hue held at one measured chroma and realized
// at one measured depth — never a blend. Compositing the pinned base over
// the neutral Surface at 12% alpha interpolates in non-linear sRGB, which is
// neither hue-preserving nor chroma-preserving: the four status grounds come
// out at chroma 0.0155–0.0212, within a rounding error of the grey they are
// mixed into, and the red one's hue slides from 28.7° to 21.6°, seven degrees
// toward magenta — a dirty pink.
//
// Realizing the container instead fixes both halves at once. The tonal
// solver holds hue exactly by construction (it reduces chroma toward the
// neutral axis, never off the hue ray), so a container of a red role is
// red-family whatever the seed; and asking every role for the same chroma
// at the same depth means the four differ in hue and in nothing else, which
// is the only arrangement in which four status grounds read as four.
//
// The mark a role puts on a ground — an icon, a leading edge, a rule — is
// not text and is not chosen against a text floor. MarkOn picks it: the
// most chromatic rung of the role's own ramp that still reaches the floor
// asked for over that ground. Said the other way round, a mark should be as
// much its own colour as it can be while still reading, and which rung that
// is depends on the hue as much as on the ground — sRGB holds a saturated
// red only at mid depths and a saturated amber only at high ones, so a
// fixed rung would serve one hue at the cost of the others. It is one rule
// for two jobs: a status container's own mark takes it against the
// container at WCAG 1.4.11's 3:1 non-text floor, and a component marking
// some other ground — a transient chip's leading edge over the inverse
// surface, say — takes it against that ground at whatever floor that job
// owes.
//
// Body text on a container is a different pairing and does not use MarkOn.
// It reads in the neutral Text token, which measures 11.6:1 or better over
// every container this file derives.
//
// Two members answer the depth question, and only the depth question:
// Container at the fixed step the family was designed around, for a
// component whose container IS the elevation it claims, and ContainerOn
// against a named ground, for one placed at an arbitrary storey. Same
// construction, same chroma, same hue; the second only declines to land the
// fill on the surface it is filling.
//
// StatusContainer, OnStatusContainer and StatusContainerOn are Container,
// OnContainer and ContainerOn under the status family's own names. Two
// spellings, one derivation: the walk asks the role's ramp rather than a
// table, so it answers for the accent trio and the status four alike.
package tokens

import (
	stdcolor "image/color"

	"github.com/vibrantgio/theme/color"
)

// ContainerFloor is the least contrast a tonal container owes the ground it
// is placed on: enough that a reader sees a filled region there at all.
//
// A container's fill is the quietest a hue is spoken at — tinted toward the
// ground until it is a field rather than a mark — so its floor is a floor on
// being a field at all, and the ceiling is left to whatever draws it.
//
// It is not a WCAG criterion, because WCAG has none for this. 1.4.11's 3:1
// governs a mark that has to be resolved as a shape; a container carries no
// shape and no information of its own — it is a place, and what it owes is
// only that its edge be findable. Gating a container at 3:1 would make every
// tonal fill in the system as loud as the marks on it.
//
// 1.25:1 is measured rather than picked. Over the seed sweep — ten seeds,
// both schemes, both contrast variants, five roles, every storey — the fixed
// step-300 container's separation from the ground falls in bands with air
// between them: 1.00–1.01 where the container IS the surface, 1.16–1.21 where
// it is a shade the eye reads as the same surface, then 1.30–1.40 and 1.42
// upward, which is where the design's own pairings sit. The threshold goes in
// the empty stretch between 1.21 and 1.30, so it accepts every depth the
// design already ships and rejects only the collisions, and no seed sits near
// enough to it to flip on a rounding. Nothing about the number is a standard;
// it is where this palette's own answers separate.
const ContainerFloor = 1.25

// Container returns the role's tonal container: its ramp's
// containerStep rung with the chroma pulled down to containerChroma.
//
// It is defined for every role that has a ramp, status or accent, because
// the derivation asks the ramp and not a table. The depth and the hue both
// come off the rung the container is realized at, which is what lets a role
// whose hue varies with depth — warning, which bends toward orange as its
// tone deepens (see seed.go) — carry that hue into its container without
// this file knowing there is a bend.
//
// The chroma is the one quantity read off the ramp's mid-value step 500
// instead, and the asymmetry is the point: hue answers "which colour, at
// this depth", which is a property of the depth, while chroma is only being
// asked "has this role a colour at all", which is a property of the role
// and is best read where the gamut constrains it least. The answer is the
// lesser of the dial and what the role actually carries, which is what
// keeps a brandless palette brandless — a neutral seed's accent ramp
// carries no chroma, so its container carries none either, rather than
// inventing a hue the brand does not have.
//
// Moving the hue reading one rung costs the families that do not bend a
// byte and no more: over a 216-colour seed grid, 1172 of the 6048
// non-warning containers change at all and every one of them by 1/255 in a
// single channel, which is the whole difference between reading an angle
// off a rung at one chroma and off a rung at another. Nothing about their
// derivation changed; the eight bits landed one step over.
//
// RoleNeutral is accepted and yields the neutral ramp's own step, chroma 0.
func (t ColorTokens) Container(role Role) stdcolor.NRGBA {
	return t.containerAt(role, containerStep)
}

// StatusContainer is [ColorTokens.Container] under the status family's own
// name; see the file header.
func (t ColorTokens) StatusContainer(role Role) stdcolor.NRGBA {
	return t.Container(role)
}

// ContainerOn returns the role's tonal container for a component that stands
// on a named ground: [ColorTokens.Container]'s rung while that rung is visibly
// a different surface from the ground, and otherwise the first deeper rung
// that is.
//
// Depth is the one thing Container cannot answer alone. It is realized
// at a fixed step, which is right for a component that fills the page it is
// on — the depth then IS the elevation the container is claiming — and wrong
// for a small one placed at an arbitrary storey, because the elevation ladder
// walks through that fixed step. A dark scheme's level-2 surface and its
// step-300 container land within 1.01:1 of each other; a container drawn
// there is not subtle, it is absent, and the component wearing it silently
// loses whatever channel the fill was carrying.
//
// The walk deepens and never turns around: the first step at or past the
// reference one that clears [ContainerFloor]. It needs no direction because
// the ramps are paired scales — 100–300 are the tinted fills of whichever
// scheme is running and 700–900 the text over them — so a higher step is
// further from that scheme's own ground in the light ramp and in the dark
// one alike. A container therefore moves only as far as the ground forces it
// to, and every role on one ground moves together, which is what keeps four
// status fills reading as one set.
//
// Unlike [MarkOn]'s nearest-to-the-reference rule, which suits a mark that
// must land at one depth whichever side of the reference it comes from.
// Nearest-first here picks a fill from the wrong end: where a dark scheme's
// own surfaces occupy steps 200–400, the closest rung that clears the floor
// is step 900, and a caller that asked for a tint gets a white block.
//
// ContainerFloor, not GraphicFloor: a container is a region and not a mark
// on one, and 3:1 against the ground would make the four statuses read as
// four filled controls. What a fill owes is to be a different surface from
// the surface, which is a threshold about seeing a seam and not about
// resolving a shape.
//
// A ground no rung separates from yields the rung that separates most, so a
// caller always has a fill — an unseparated container is a defect the gates
// report, not a reason to paint nothing.
//
// Over the seed sweep the answer stays in a narrow band: worst 1.30:1,
// loudest 1.72:1. A container is never invisible and never a solid.
func (t ColorTokens) ContainerOn(role Role, ground stdcolor.NRGBA) stdcolor.NRGBA {
	best, bestAt := -1.0, containerStep
	for step := containerStep; step <= 900; step += 100 {
		fill := t.containerAt(role, step)
		got := color.ContrastRatio(fill, ground)
		if got >= ContainerFloor {
			return fill
		}
		if got > best {
			best, bestAt = got, step
		}
	}
	return t.containerAt(role, bestAt)
}

// StatusContainerOn is [ColorTokens.ContainerOn] under the status family's own
// name; see the file header.
func (t ColorTokens) StatusContainerOn(role Role, ground stdcolor.NRGBA) stdcolor.NRGBA {
	return t.ContainerOn(role, ground)
}

// containerAt realizes the role's container at one ramp step: the step's own
// tone and hue, at the chroma read off the ramp's mid-value step 500 and
// clamped to containerChroma. See [ColorTokens.Container] for why the two
// readings come off two different rungs.
func (t ColorTokens) containerAt(role Role, step int) stdcolor.NRGBA {
	r := t.rampFor(role) // validates role
	rung := r.Step(step)
	tone, _, _ := color.LabFromNRGBA(rung)
	_, _, hue := color.OKLChFromNRGBA(rung)
	_, chroma, _ := color.OKLChFromNRGBA(r.Step(500))
	if chroma > containerChroma {
		chroma = containerChroma
	}
	return color.NRGBAFromToneChromaHue(tone, chroma, hue)
}

// OnContainer returns the colour of the role's mark over its own
// container: MarkOn against Container(role) at graphicFloor, WCAG
// 1.4.11's 3:1 for a non-text graphic. Every light scheme lands on step 700
// and every dark scheme on step 500 — one depth per scheme for all four
// roles — and the worst pairing over the whole seed sweep measures 4.47:1.
//
// It is the floor a mark owes, so a caller setting a run of words on a
// container asks for the text floor instead: the role's ink derived against
// the fill ([ColorTokens.InkOn], or [ColorTokens.MarkOn] for RoleNeutral,
// at [TextFloor]).
func (t ColorTokens) OnContainer(role Role) stdcolor.NRGBA {
	return t.MarkOn(role, t.Container(role), graphicFloor)
}

// OnStatusContainer is [ColorTokens.OnContainer] under the status family's
// own name; see the file header.
func (t ColorTokens) OnStatusContainer(role Role) stdcolor.NRGBA {
	return t.OnContainer(role)
}

// MarkOn returns the colour the role marks ground with: the rung of the
// role's own ramp nearest the ramp's mid-value step 500 that reaches floor
// against ground.
//
// Step 500 is where a role is most itself — it is the ramp's mid-value
// reference — so the rule is "be that rung, or
// the nearest one to it that still reads". Walking outward from the middle
// rather than naming a rung is what keeps four hues in step. sRGB holds a
// saturated red only at mid depths and a saturated amber only at high ones,
// so a fixed rung serves one hue at the cost of the others; but the
// nearest-to-mid rung that clears a given ground is the same rung for all
// four, so every mark on one ground lands at one depth and reads at one
// weight, which is what makes a set of status marks a set.
//
// Reaching for the most chromatic rung instead does not hold that. Chroma
// peaks at a different depth for every hue — amber's peaks two rungs deeper
// on the dark scale than its siblings' do — so the amber mark comes back at
// 8.44:1 beside three siblings at 5.03:1 and pulls the eye across an alert
// stack for a reason no one reading it could infer. What a mark owes its
// ground is legibility, and past that, agreement with the other marks.
//
// A ground no rung can clear yields the rung that reads best rather than
// nothing, so a caller always has a colour: a mark too weak to meet its
// floor is a contrast defect the gates report, not a reason for a component
// to paint an unset colour.
func (t ColorTokens) MarkOn(role Role, ground stdcolor.NRGBA, floor float64) stdcolor.NRGBA {
	const mid = 4        // index of step 500, the ramp's mid-value reference
	r := t.rampFor(role) // validates role
	pick, dist := -1, 99
	fallback, fallbackAt := -1.0, 0
	for i, rung := range r {
		got := color.ContrastRatio(rung, ground)
		if got > fallback {
			fallback, fallbackAt = got, i
		}
		if got < floor {
			continue
		}
		d := i - mid
		if d < 0 {
			d = -d
		}
		if d < dist {
			pick, dist = i, d
		}
	}
	if pick < 0 {
		return r[fallbackAt]
	}
	return r[pick]
}
