// Tonal containers: the wash a status role fills the surface it stands on
// with, and the mark read on it.
//
// A container is a role's own hue held at one measured chroma and realized
// at one measured depth — never a blend. Compositing the pinned base over
// the neutral Surface at 12% alpha interpolates in non-linear sRGB, which is
// neither hue-preserving nor chroma-preserving: the four status washes come
// out at chroma 0.0155–0.0212, within a rounding error of the grey they are
// mixed into, and the red one's hue slides from 28.7° to 21.6°, seven degrees
// toward magenta — a dirty pink.
//
// Realizing the container instead fixes both halves at once. The tonal
// solver holds hue exactly by construction (it reduces chroma toward the
// neutral axis, never off the hue ray), so a container of a red role is
// red-family whatever the seed; and asking every role for the same chroma
// at the same depth means the four differ in hue and in nothing else, which
// is the only arrangement in which four status washes read as four.
//
// A fill's hue is its role's, not its depth's, so it is read at one fixed
// step — the ramp's pale tint depth, the third step counted from the ramp's
// pale end — and not at whatever depth the fill is realized at. Reading it
// off the realized step instead would let the eight-bit realization and the
// gamut mapping at that depth answer a question about which colour the role
// is: at the container dial a step's hue reads up to 1.42° off the angle it
// was asked for, and the deepest steps of a hue sRGB starves read further
// off still.
//
// That also leaves the walk below with one job. A hue read off the realized
// step rotates whenever the walk deepens the fill, which would put a warning
// badge on the content and one on the window's furniture beside it in two
// different hues for a reason no reader could infer.
//
// The mark a role puts on a ground — an icon, a leading edge, a rule — is
// not text and is not chosen against a text floor. MarkOn picks it: the
// most chromatic rung of the role's own ramp that still reaches the floor
// asked for over that ground. Said the other way round, a mark should be as
// much its own colour as it can be while still reading, and which step that
// is depends on the hue as much as on the surface — sRGB holds a saturated
// red only at mid depths and a saturated orange only at high ones, so a
// fixed step would serve one hue at the cost of the others. It is one rule
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
// against a named surface, for one placed at an arbitrary level. Same
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

// ContainerSeparation is the least perceptual distance two status fills
// drawn on one surface may come to, in OKLab. The four are one set, and a
// set of four is only four if a reader can tell its members apart.
//
// 0.028 is measured rather than picked. The error and the warning are the
// closest pair the anchors leave — an orange beside a red, 35.35° apart
// before any tint — and a seed lying between them tints them toward each
// other, so over the seed sweep, both derivations, both schemes, every
// level, the pair closes to 30.18° at the dial and 0.0286 in OKLab. The
// threshold goes just under that, so it accepts what the derivation can
// produce and rejects a set that has closed further.
//
// Two readings bracket it. Two fills 31.4° apart at this dial read as two
// colours; two 19.3° apart read as one brown twice, which is where a hue
// read off the realized rung put a warning that rotated with depth. The
// derivation ships the first kind.
//
// Hue angle alone cannot state it — two hues 90° apart at no chroma are one
// grey — so the distance is taken in the space the derivation places its
// hues and chromas in.
const ContainerSeparation = 0.028

// Container returns the role's tonal container: its ramp's containerStep
// rung's tone, at the hue the role wears at its ramp's pale tint depth,
// with the chroma pulled down to containerChroma.
//
// It is defined for every role that has a ramp, status or accent, because
// the derivation asks the ramp and not a table.
//
// Three readings off three rungs, and the asymmetry is the point. Tone
// answers "how deep", which is the container's own question and so is read
// at the rung it is realized at. Hue answers "which colour", which is a
// property of the role rather than of the depth — see the file header for
// what reading it off the realized rung cost the status set — so it is read
// at the pale tint depth. Chroma is only being asked "has this role a colour
// at all", which is best read where the gamut constrains it least, so it
// comes off the ramp's mid-value step 500; the answer is the lesser of the
// dial and what the role actually carries, which is what keeps a brandless
// palette brandless — a neutral seed's accent ramp carries no chroma, so its
// container carries none either, rather than inventing a hue the brand does
// not have.
//
// Reading the hue one step over costs a container an eight-bit step and no
// more: over the seed sweep every container that changed at all changed by
// at most two units summed across its channels, which is the whole
// difference between reading an angle off a step at one chroma and off a
// step at another.
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
// on a named surface: [ColorTokens.Container]'s rung while that rung is
// visibly a different surface from that one, and otherwise the first deeper
// rung that is.
//
// Depth is the one thing Container cannot answer alone. It is realized
// at a fixed step, which is right for a component that fills the page it is
// on — the depth then IS the elevation the container is claiming — and wrong
// for a small one placed at an arbitrary level, because the elevation ladder
// walks through that fixed step. A dark scheme's level-2 surface and its
// step-300 container land within 1.01:1 of each other; a container drawn
// there is not subtle, it is absent, and the component wearing it silently
// loses whatever channel the fill was carrying.
//
// The walk deepens and never turns around: the first step at or past the
// reference one that clears [ContainerFloor]. It needs no direction because
// the ramps are paired scales — 100–300 are the tinted fills of whichever
// scheme is running and 700–900 the text over them — so a higher step is
// further from that scheme's own surfaces in the light ramp and in the dark
// one alike. A container therefore moves only as far as the surface forces
// it to, and every role on one surface moves together, which is what keeps
// four status fills reading as one set.
//
// The walk moves depth and nothing else: hue comes off the pale tint depth
// whatever rung the walk lands on (see the file header), so a role's wash is
// one hue wherever it is drawn.
//
// Unlike [MarkOn]'s nearest-to-the-reference rule, which suits a mark that
// must land at one depth whichever side of the reference it comes from.
// Nearest-first here picks a fill from the wrong end: where a dark scheme's
// own surfaces occupy steps 200–400, the closest rung that clears the floor
// is step 900, and a caller that asked for a tint gets a white block.
//
// ContainerFloor, not GraphicFloor: a container is a region and not a mark
// on one, and 3:1 would make the four statuses read as four filled controls.
// What a fill owes is to be a different surface from the one it stands on,
// which is a threshold about seeing a seam and not about resolving a shape.
//
// A surface no rung separates from yields the rung that separates most, so a
// caller always has a fill — an unseparated container is a defect the gates
// report, not a reason to paint nothing.
//
// Over the seed sweep the answer stays in a narrow band: worst 1.30:1,
// loudest 1.72:1. A container is never invisible and never a solid.
func (t ColorTokens) ContainerOn(role Role, surface stdcolor.NRGBA) stdcolor.NRGBA {
	return fillToFloor(func(step int) stdcolor.NRGBA { return t.containerAt(role, step) }, surface)
}

// fillToFloor is the walk both ground-aware tonal fills take: the first
// step at or past containerStep whose realization clears [ContainerFloor]
// against surface, and otherwise the step that separates most. at realizes
// one step of the fill being walked — a role's container, or the reserved
// highlighter (highlight.go). See [ColorTokens.ContainerOn] for why the
// walk only ever deepens and why an unseparated fill still returns a
// colour.
func fillToFloor(at func(step int) stdcolor.NRGBA, surface stdcolor.NRGBA) stdcolor.NRGBA {
	best, bestAt := -1.0, containerStep
	for step := containerStep; step <= 900; step += 100 {
		fill := at(step)
		got := color.ContrastRatio(fill, surface)
		if got >= ContainerFloor {
			return fill
		}
		if got > best {
			best, bestAt = got, step
		}
	}
	return at(bestAt)
}

// StatusContainerOn is [ColorTokens.ContainerOn] under the status family's own
// name; see the file header.
func (t ColorTokens) StatusContainerOn(role Role, surface stdcolor.NRGBA) stdcolor.NRGBA {
	return t.ContainerOn(role, surface)
}

// containerAt realizes the role's container at one ramp step: the step's own
// tone, the hue read at the ramp's pale tint depth, and the chroma read off
// the ramp's mid-value step 500 and clamped to containerChroma. See
// [ColorTokens.Container] for why the three readings come off three
// different rungs.
func (t ColorTokens) containerAt(role Role, step int) stdcolor.NRGBA {
	r := t.rampFor(role) // validates role
	tone, _, _ := color.LabFromNRGBA(r.Step(step))
	_, _, hue := color.OKLChFromNRGBA(r.Step(containerHueStep(r)))
	_, chroma, _ := color.OKLChFromNRGBA(r.Step(500))
	if chroma > containerChroma {
		chroma = containerChroma
	}
	return color.NRGBAFromToneChromaHue(tone, chroma, hue)
}

// containerHueStep is the step a container reads its hue at: the third rung
// from the ramp's pale end, which is step 300 in a light scheme's ramp and
// step 700 in a dark scheme's. The two are one tone by construction — the
// ramps are paired scales — and it is the tone at which the one family whose
// hue moves with depth still carries its anchor, so a wash realized anywhere
// on the ramp wears the hue of its role and not of its depth.
//
// The end is read off the scale rather than named, the way every other
// scheme-agnostic rule in the derivation reads it: nothing here knows which
// scheme is running.
func containerHueStep(r Ramp) int {
	pale, _, _ := color.LabFromNRGBA(r.Step(100))
	deep, _, _ := color.LabFromNRGBA(r.Step(900))
	if pale >= deep {
		return 300
	}
	return 700
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
// saturated red only at mid depths and a saturated orange only at high ones,
// so a fixed step serves one hue at the cost of the others; but the
// nearest-to-mid step that clears a given ground is the same step for all
// four, so every mark on one ground lands at one depth and reads at one
// weight, which is what makes a set of status marks a set.
//
// Reaching for the most chromatic step instead does not hold that. Chroma
// peaks at a different depth for every hue — the green's peaks two steps
// deeper on the dark scale than its siblings' do — so the green mark comes
// back at 9.45:1 beside three siblings at 4.87:1 and pulls the eye across an
// alert stack for a reason no one reading it could infer. What a mark owes its
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
