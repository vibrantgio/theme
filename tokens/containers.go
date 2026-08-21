// Tonal containers: the ground a status role fills, and the mark read on it.
//
// A container is a role's own hue held at one measured chroma and realized
// at one measured depth — never a blend. Blending was how the tinted banner
// used to get its ground: the pinned base composited over the neutral
// Surface at 12% alpha, in non-linear sRGB. That is neither hue-preserving
// nor chroma-preserving, and it showed. The four status grounds came out at
// chroma 0.0155–0.0212 — within a rounding error of the grey they were
// mixed into, so no two of them could be told apart — and the red one's hue
// slid from 28.7° to 21.6° on the way, seven degrees toward magenta. A red
// that has given up seven-eighths of its chroma and seven degrees of its
// hue is what a container reported as dirty pink actually was.
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
package tokens

import (
	stdcolor "image/color"

	"github.com/vibrantgio/theme/color"
)

// StatusContainer returns the role's tonal container: its ramp's
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
func (t ColorTokens) StatusContainer(role Role) stdcolor.NRGBA {
	r := t.rampFor(role) // validates role
	rung := r.Step(containerStep)
	tone, _, _ := color.LabFromNRGBA(rung)
	_, _, hue := color.OKLChFromNRGBA(rung)
	_, chroma, _ := color.OKLChFromNRGBA(r.Step(500))
	if chroma > containerChroma {
		chroma = containerChroma
	}
	return color.NRGBAFromToneChromaHue(tone, chroma, hue)
}

// OnStatusContainer returns the colour of the role's mark over its own
// container: MarkOn against StatusContainer(role) at graphicFloor, WCAG
// 1.4.11's 3:1 for a non-text graphic. Every light scheme lands on step 700
// and every dark scheme on step 500 — one depth per scheme for all four
// roles — and the worst pairing over the whole seed sweep measures 4.47:1.
func (t ColorTokens) OnStatusContainer(role Role) stdcolor.NRGBA {
	return t.MarkOn(role, t.StatusContainer(role), graphicFloor)
}

// MarkOn returns the colour the role marks ground with: the rung of the
// role's own ramp nearest the ramp's mid-value step 500 that reaches floor
// against ground.
//
// Step 500 is where a role is most itself — ADR-007 gives it the job of
// being the ramp's mid-value reference — so the rule is "be that rung, or
// the nearest one to it that still reads". Walking outward from the middle
// rather than naming a rung is what keeps four hues in step. sRGB holds a
// saturated red only at mid depths and a saturated amber only at high ones,
// so a fixed rung serves one hue at the cost of the others; but the
// nearest-to-mid rung that clears a given ground is the same rung for all
// four, so every mark on one ground lands at one depth and reads at one
// weight, which is what makes a set of status marks a set.
//
// Reaching for the most chromatic rung instead was tried and does not hold
// that. Chroma peaks at a different depth for every hue — amber's peaks two
// rungs deeper on the dark scale than its siblings' do — so the amber mark
// came back at 8.44:1 beside three siblings at 5.03:1 and pulled the eye
// across an alert stack for a reason no one reading it could infer. What a
// mark owes its ground is legibility, and past that, agreement with the
// other marks; being the most colourful it could have been is neither.
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
