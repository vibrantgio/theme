// The two variant rungs: the quiet neutral boundary and the muted neutral
// ink, both floored by construction rather than pinned to a step.
//
// A fixed rung cannot hold either job. The neutral ramp is a paired scale —
// the same step keeps the same job in both schemes — but a step is a depth,
// not a contrast, and the two grounds a neutral thing stands on without
// naming a storey ([ColorTokens.Surface] and [ColorTokens.Background]) sit at
// different depths in the light scheme and the dark one. Naming step 500 as
// "the outline" therefore states one colour and two measurements: over the two
// grounds it reads 3.07:1 in the dark scheme and 2.35:1 in the light, under
// the floor a boundary owes in the scheme most people read in. That is the
// failure this file exists to make impossible — the rung is chosen against the
// floor, so the token cannot vanish on a ground.
//
// Both members ask the same walk at two floors, which is what keeps them one
// idiom rather than two: the rung nearest the ramp's mid-value step 500 that
// clears the floor against BOTH neutral grounds. Nearest-to-the-middle is
// [ColorTokens.MarkOn]'s rule, and it lands on the quietest rung that clears
// in every scheme this palette produces, because the rungs below the middle
// are the tinted fills — too close to the ground to reach any floor over it —
// so the walk stops at the middle or at the first rung past it that clears.
//
// Both grounds rather than one because a neutral boundary and a neutral ink
// are drawn on the page and on the furniture alike, and a token that cleared
// only the easier of the two would be a promise kept on one of them.
//
// What the walk answers, over the seed sweep — 414 seeds, both schemes, both
// derivations:
//
//	token               scheme   step   ratio over the harder ground
//	-----               ------   ----   ----------------------------
//	OutlineVariant      light    600    3.55:1
//	OutlineVariant      dark     500    3.07:1
//	OnSurfaceVariant    light    700    5.46:1 (high contrast: 15.16:1)
//	OnSurfaceVariant    dark     600    5.72:1
//
// The rung does not move with the seed: the neutral ramp is realized on the
// shared lightness scale and both grounds come off it, so a brand changes the
// tint of these two colours and never their reading.
package tokens

import (
	stdcolor "image/color"

	"github.com/vibrantgio/theme/color"
)

// OutlineVariant returns the neutral boundary a resting edge is drawn in: the
// quiet line that says a region or a control is there without claiming to be
// its content.
//
// Floored at graphicFloor, WCAG 1.4.11's 3:1 — an edge that is the whole of
// what says which control this is carries meaning without being text, so it
// is not decoration and does not get a decorative floor. [ColorTokens.Divider]
// is the token for a separator that carries none.
func (t ColorTokens) OutlineVariant() stdcolor.NRGBA {
	return t.neutralVariant(graphicFloor)
}

// OnSurfaceVariant returns the muted neutral ink: the rung a secondary run of
// words is set in — quieter than [ColorTokens.Text], and still a colour text
// may legally be set in.
//
// Floored at onFloor, WCAG 1.4.3 AA's 4.5:1, because it is text. Muted is a
// property of the walk and not a second rule: the floor picks the quietest
// rung that reads, and Text is a pin derived against Background with far more
// room than the floor asks for, so the two are never one colour.
func (t ColorTokens) OnSurfaceVariant() stdcolor.NRGBA {
	return t.neutralVariant(onFloor)
}

// neutralVariant returns the rung of the neutral ramp nearest the ramp's
// mid-value step 500 that reaches floor against Surface and Background both.
//
// A palette whose ramp cleared neither ground yields the rung that comes
// closest, so a caller always has a colour: a boundary too weak to meet its
// floor is a contrast defect the gates report, not a reason to draw nothing.
func (t ColorTokens) neutralVariant(floor float64) stdcolor.NRGBA {
	const mid = 4 // index of step 500, the ramp's mid-value reference
	pick, dist := -1, len(t.Ramps.Neutral)
	widest, widestAt := -1.0, 0
	for i, rung := range t.Ramps.Neutral {
		worst := color.ContrastRatio(rung, t.Surface)
		if got := color.ContrastRatio(rung, t.Background); got < worst {
			worst = got
		}
		if worst > widest {
			widest, widestAt = worst, i
		}
		if worst < floor {
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
		return t.Ramps.Neutral[widestAt]
	}
	return t.Ramps.Neutral[pick]
}
