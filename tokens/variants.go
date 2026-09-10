// The two variant steps: the faint neutral boundary and the muted neutral
// foreground, both floored by construction rather than pinned to a step.
//
// A fixed step cannot hold either job. The neutral ramp is a paired scale —
// the same step keeps the same job in both schemes — but a step is a depth,
// not a contrast, and the two surfaces a neutral thing stands on without
// naming a level ([ColorTokens.Surface] and [ColorTokens.Background]) sit at
// different depths in the light scheme and the dark one. Naming step 500 as
// "the outline" therefore states one colour and two measurements: over the two
// surfaces it reads |Lc| 41.53 in the light scheme and 23.64 in the dark, under
// the floor a boundary owes in both. That is the
// failure this file exists to make impossible — the step is chosen against the
// floor, so the token cannot vanish on a surface.
//
// Both members ask the same walk at two floors, which is what keeps them one
// idiom rather than two: the step nearest the ramp's mid-value step 500 that
// clears the floor against BOTH neutral surfaces. Nearest-to-the-middle is
// [ColorTokens.MarkOn]'s rule, and it lands on the least pronounced step that
// clears in every scheme this palette produces, because the steps below the
// middle are the tinted fills — too close to the surface to reach any floor
// over it — so the walk stops at the middle or at the first step past it that
// clears.
//
// Both surfaces rather than one because a neutral boundary and a neutral
// foreground are drawn on the page and on the chrome alike, and a token
// that cleared only the easier of the two would be a promise kept on one of
// them.
//
// What the walk answers, over the seed sweep — 414 seeds, both schemes, both
// derivations:
//
//	token               scheme   step   |Lc| over the harder surface
//	-----               ------   ----   ---------------------------
//	OutlineVariant      light    600    56.58
//	OutlineVariant      dark     600    45.81
//	OnSurfaceVariant    light    800    79.76 (increased contrast: 700, 91.57)
//	OnSurfaceVariant    dark     800    79.90 (increased contrast: 700, 94.34)
//
// The step does not move with the seed: the neutral ramp is realized on the
// shared lightness scale and both surfaces come off it, so a brand changes the
// tint of these two colours and never their reading.
package tokens

import (
	stdcolor "image/color"

	"github.com/vibrantgio/theme/color"
)

// OutlineVariant returns the neutral boundary a resting edge is drawn in: the
// faint line that says a region or a control is there without claiming to be
// its content.
//
// Floored at graphicFloor, APCA's Lc 45 for a mark — an edge that is the whole of
// what says which control this is carries meaning without being text, so it
// is not decoration and does not get a decorative floor. [ColorTokens.Seam]
// is the token for a separator that carries none.
func (t ColorTokens) OutlineVariant() stdcolor.NRGBA {
	return t.neutralVariant(graphicFloor)
}

// OnSurfaceVariant returns the muted neutral foreground: the step a secondary
// run of words is set in — less pronounced than [ColorTokens.Text], and still
// a colour text may legally be set in.
//
// Floored at onFloor, APCA's Lc 75 for body text, because it is text. Muted is a
// property of the walk and not a second rule: the floor picks the least
// pronounced step that reads, and Text is a pin derived against Background
// with far more room than the floor asks for, so the two are never one colour.
func (t ColorTokens) OnSurfaceVariant() stdcolor.NRGBA {
	return t.neutralVariant(onFloor)
}

// neutralVariant returns the step of the neutral ramp nearest the ramp's
// mid-value step 500 that reaches floor against Surface and Background both.
//
// A palette whose ramp cleared neither surface yields the step that comes
// closest, so a caller always has a colour: a boundary too weak to meet its
// floor is a contrast defect the gates report, not a reason to draw nothing.
func (t ColorTokens) neutralVariant(floor float64) stdcolor.NRGBA {
	const mid = 4 // index of step 500, the ramp's mid-value reference
	pick, dist := -1, len(t.Ramps.Neutral)
	widest, widestAt := -1.0, 0
	for i, step := range t.Ramps.Neutral {
		worst := color.Magnitude(step, t.Surface)
		if got := color.Magnitude(step, t.Background); got < worst {
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
