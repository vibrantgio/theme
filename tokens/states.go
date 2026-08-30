// Interaction states as step walks.
//
// There are no MD3 alpha state layers — no translucent 8%/12%/16% overlays of
// the role colour. A state is a real, addressable colour a token sheet can
// emit. Two regimes:
//
//   - Tinted (ground) surfaces walk the ramp by index: hover is one step
//     past the component's ground, pressed and selected two (a card on 200
//     hovers on 300 and presses on 400). Because the light and dark ramps
//     are paired scales — same step, same job — the SAME index walk holds
//     in both modes: past 200 lies 300, a darker hover on a light ground
//     and a lighter one on a dark ground, with no mode-specific rule.
//
//   - Solid fills — the pinned role bases — walk one or two steps from the
//     pin toward 900. Pins are off-ramp by design (the seed sits deep, so
//     bases are pins), so a pin cannot walk by ramp index. The rule, chosen
//     to keep "same step, same job" across modes: the pin's CIELAB L* is
//     located on its
//     role ramp's own measured L* ladder as a fractional step index; each
//     state moves that index toward the 900 end by the tinted regime's walk
//     (hover one rung, pressed and selected two), the target L* is read off
//     the ladder by linear interpolation, and the colour is realized at the
//     pin's own OKLCh hue and chroma by the tonal solver. The dark pin sits
//     at the dark scale's step-700 depth by construction, so its hover lands
//     at step-800 depth and its pressed clamps at
//     the step-900 depth — exactly the walk a tinted surface performs —
//     while the light seed pin (≈ step-700 depth for the default seed)
//     darkens toward the 800/900 depths. In both modes the walk heads
//     toward the 900 end of the paired scale: darker in light mode, lighter
//     in dark mode.
//
//   - A pin the scheme does not carry at all — a colour a caller pins on one
//     control, for which no role and no ramp exists — walks by the very same
//     rule (PinnedStateColor). Having no role ramp to ladder on, it ladders
//     on the neutral one: every ramp in a scheme sweeps the same lightness
//     scale and only the neutral sweeps it at zero chroma, so the neutral
//     ramp is that scale itself rather than one role's tinted copy of it.
//

// Clamping: a walk past the ramp end clamps to the 900 stop — ground 800
// pressed resolves to step 900, ground 900 hover stays at 900, and a solid
// walk never passes the 900 rung's depth. Clamping is friendlier than
// erroring: a component on the deepest ground still gets a stable colour.
//
// The remaining states are not walks: disabled is an opacity
// (DisabledOpacity, MD3's 38%) applied to the state's normal colour; focus
// keeps the surface colour and adds the focus ring (FocusRing, Neutral
// step 500, the strong-border focusable edge); dragged resolves exactly as
// pressed.
package tokens

import (
	"fmt"
	stdcolor "image/color"
	"math"

	"github.com/vibrantgio/theme/color"
)

// State enumerates the interaction states a component surface can be in.
// Hover, Pressed, Selected and Dragged resolve as step walks; Disabled is
// an opacity; Focus keeps the surface colour (the ring is FocusRing).
type State int

const (
	StateNormal State = iota
	StateHover
	StatePressed
	StateSelected
	StateDisabled
	StateFocus
	StateDragged
)

// Role selects one of the colour-role ramps of a RampSet, and for the
// accent and status roles also the matching pinned base on ColorTokens.
// Neutral has a ramp but no pinned solid fill.
type Role int

const (
	RoleNeutral Role = iota
	RolePrimary
	RoleSecondary
	RoleTertiary
	RoleError
	RoleSuccess
	RoleWarning
	// RoleInfo is last because the values above it are what callers have
	// compiled against; a fourth status role is an addition, not a
	// renumbering.
	RoleInfo
)

// DisabledOpacity is the fraction of full alpha a disabled element keeps —
// MD3's 38% disabled-content opacity (alpha 0x61 of 0xff). Disabled is an
// opacity, not a ramp step: a disabled surface keeps its colour and fades.
const DisabledOpacity = 0.38

// Disabled returns c with its alpha scaled by DisabledOpacity. Colours are
// non-premultiplied NRGBA, so scaling alpha alone is the whole rule; an
// opaque colour comes back at alpha 0x61.
func Disabled(c stdcolor.NRGBA) stdcolor.NRGBA {
	c.A = uint8(math.Round(float64(c.A) * DisabledOpacity))
	return c
}

// FocusRing returns the focus-ring colour: Neutral step 500, the
// strong-border focusable edge.
func (t ColorTokens) FocusRing() stdcolor.NRGBA {
	return t.Ramps.Neutral.Step(500)
}

// StateColor resolves a tinted surface: the colour of a component whose
// normal ground is the given ramp step (100–900 in hundreds), under the
// given state. Hover walks one step toward 900, pressed, selected and
// dragged two; walks past the ramp end clamp to step 900. Normal and focus
// return the ground itself (draw FocusRing for the ring); disabled returns
// the ground at DisabledOpacity. An out-of-vocabulary ground, role or
// state panics, matching Ramp.Step.
func (t ColorTokens) StateColor(role Role, ground int, state State) stdcolor.NRGBA {
	r := t.rampFor(role)
	base := r.Step(ground) // validates ground
	switch state {
	case StateNormal, StateFocus:
		return base
	case StateDisabled:
		return Disabled(base)
	}
	step := ground + 100*stateWalk(state)
	if step > 900 {
		step = 900 // clamp at the ramp end
	}
	return r.Step(step)
}

// SolidStateColor resolves a solid fill: the role's pinned base under the
// given state, walking from the pin toward the 900 end of its ramp per the
// rule in the package-file header. Normal and focus return the pin itself;
// disabled returns the pin at DisabledOpacity. RoleNeutral panics: neutral
// has no pinned solid fill.
func (t ColorTokens) SolidStateColor(role Role, state State) stdcolor.NRGBA {
	pin := t.pinFor(role)
	switch state {
	case StateNormal, StateFocus:
		return pin
	case StateDisabled:
		return Disabled(pin)
	}
	return solidWalk(pin, t.rampFor(role), stateWalk(state))
}

// PinnedStateColor resolves a solid fill from a pin the scheme does not
// carry: a colour the caller fixes on one control — a status the palette has
// no role for, a shade a platform prescribes and expects to survive a scheme
// change — under the given state. The walk is SolidStateColor's: normal and
// focus return the pin itself, disabled returns it at DisabledOpacity, and
// hover and press move it one rung and two toward the 900 end at its own hue
// and chroma.
//
// What differs is only the ladder those rungs are counted on. A role's pin
// ladders on its role's ramp; a caller's pin belongs to no role, so it
// ladders on the neutral ramp — the scheme's shared lightness scale carried
// at zero chroma, which every coloured ramp is a tinted sweep of. The
// substitution is small by construction: across a seed sweep the role
// ladders sit within an L* of the neutral one (TestRoleLaddersTrackTheNeutral
// Ladder), all of them being the same scale realized at different chromas.
//
// The pin is expected opaque, as the scheme's own pins are: the walked
// states are realized opaque whatever alpha the pin carried, and only
// disabled scales alpha.
//
// The walk inherits the scale's own coarseness, which is a feature of the
// ladder rather than of the pin: where two rungs sit far apart — the light
// scale's 28 and 6, say — a press from between them lands far away, exactly
// as the scheme's own primary pin does there.
func (t ColorTokens) PinnedStateColor(pin stdcolor.NRGBA, state State) stdcolor.NRGBA {
	switch state {
	case StateNormal, StateFocus:
		return pin
	case StateDisabled:
		return Disabled(pin)
	}
	return solidWalk(pin, t.Ramps.Neutral, stateWalk(state))
}

// stateWalk returns how many ramp steps past the ground a state sits:
// hover one, pressed, selected and dragged two, everything else zero.
func stateWalk(state State) int {
	switch state {
	case StateNormal, StateDisabled, StateFocus:
		return 0
	case StateHover:
		return 1
	case StatePressed, StateSelected, StateDragged:
		return 2
	}
	panic(fmt.Sprintf("tokens: unknown State %d", state))
}

func (t ColorTokens) rampFor(role Role) Ramp {
	switch role {
	case RoleNeutral:
		return t.Ramps.Neutral
	case RolePrimary:
		return t.Ramps.Primary
	case RoleSecondary:
		return t.Ramps.Secondary
	case RoleTertiary:
		return t.Ramps.Tertiary
	case RoleError:
		return t.Ramps.Error
	case RoleSuccess:
		return t.Ramps.Success
	case RoleWarning:
		return t.Ramps.Warning
	case RoleInfo:
		return t.Ramps.Info
	}
	panic(fmt.Sprintf("tokens: unknown Role %d", role))
}

func (t ColorTokens) pinFor(role Role) stdcolor.NRGBA {
	switch role {
	case RolePrimary:
		return t.Primary
	case RoleSecondary:
		return t.Secondary
	case RoleTertiary:
		return t.Tertiary
	case RoleError:
		return t.Error
	case RoleSuccess:
		return t.Success
	case RoleWarning:
		return t.Warning
	case RoleInfo:
		return t.Info
	}
	panic(fmt.Sprintf("tokens: Role %d has no pinned solid fill", role))
}

// solidWalk moves the pin n rungs toward the 900 end of r's measured L*
// ladder and realizes the target depth at the pin's own hue and chroma.
func solidWalk(pin stdcolor.NRGBA, r Ramp, n int) stdcolor.NRGBA {
	var ladder [9]float64
	for i, c := range r {
		ladder[i], _, _ = color.LabFromNRGBA(c)
	}
	pinL, _, _ := color.LabFromNRGBA(pin)
	idx := ladderIndex(ladder, pinL) + float64(n)
	if idx > 8 {
		idx = 8 // clamp at the 900 rung
	}
	targetL := ladderAt(ladder, idx)
	_, chroma, hue := color.OKLChFromNRGBA(pin)
	return color.NRGBAFromToneChromaHue(targetL, chroma, hue)
}

// ladderIndex locates L on a monotonic (ascending or descending) ladder as
// a fractional index in [0,8], clamping beyond either end.
func ladderIndex(ladder [9]float64, L float64) float64 {
	descending := ladder[0] > ladder[8]
	past := func(a, b float64) bool { // b lies at or past a, toward index 8
		if descending {
			return b <= a
		}
		return b >= a
	}
	if past(L, ladder[0]) {
		return 0
	}
	for i := 0; i < 8; i++ {
		if past(L, ladder[i+1]) {
			return float64(i) + (L-ladder[i])/(ladder[i+1]-ladder[i])
		}
	}
	return 8
}

// ladderAt reads the ladder's L* at a fractional index by linear
// interpolation between adjacent rungs.
func ladderAt(ladder [9]float64, idx float64) float64 {
	i := int(idx)
	if i >= 8 {
		return ladder[8]
	}
	return ladder[i] + (idx-float64(i))*(ladder[i+1]-ladder[i])
}
