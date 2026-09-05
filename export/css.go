package export

import (
	"fmt"
	stdcolor "image/color"
	"strconv"
	"strings"
	"time"

	vgcolor "github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/tokens"
)

// cssVar is one custom property, emitted in declaration order.
type cssVar struct {
	name, value string
}

// hexRGB formats a colour as lowercase #rrggbb. Every token colour is fully
// opaque, so alpha is never written.
func hexRGB(c stdcolor.NRGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// fnum formats a float32 with no trailing zeros: 16 → "16", 0.15 → "0.15".
func fnum(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

// px formats a dp value as a CSS px length. Device-independent pixels map
// 1:1 onto CSS px, both being density-abstract logical pixels.
func px(v float32) string {
	return fnum(v) + "px"
}

// rampRoles orders the colour roles under their CSS names.
var rampRoles = []struct {
	name string
	ramp func(tokens.RampSet) tokens.Ramp
}{
	{"neutral", func(r tokens.RampSet) tokens.Ramp { return r.Neutral }},
	{"primary", func(r tokens.RampSet) tokens.Ramp { return r.Primary }},
	{"secondary", func(r tokens.RampSet) tokens.Ramp { return r.Secondary }},
	{"tertiary", func(r tokens.RampSet) tokens.Ramp { return r.Tertiary }},
	{"error", func(r tokens.RampSet) tokens.Ramp { return r.Error }},
	{"success", func(r tokens.RampSet) tokens.Ramp { return r.Success }},
	{"warning", func(r tokens.RampSet) tokens.Ramp { return r.Warning }},
	{"info", func(r tokens.RampSet) tokens.Ramp { return r.Info }},
}

// pinRoles orders the pinned bases and the semantic layer under their CSS
// names; the doc comment on this package records the mapping.
var pinRoles = []struct {
	name string
	pick func(tokens.ColorTokens) stdcolor.NRGBA
}{
	{"bg", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Background }},
	{"surface", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Surface }},
	{"text", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Text }},
	{"divider", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Divider }},
	// The inverse pair, emitted as first-class tokens for the same reason
	// the state walk below is: it resolves off the counterpart scheme's
	// neutral ramp, and a sheet holding only this scheme's ramps has no
	// var() arithmetic that could reach it.
	{"inverse-surface", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.InverseSurface }},
	{"on-inverse-surface", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.OnInverseSurface }},
	// The reserved highlighter: the fill marking content the reader was
	// brought to. It is emitted as a first-class token because it belongs
	// to no ramp — its hue is reserved outside the roles, so no var()
	// reference over the ramps could reach it — and it is the fill resolved
	// for the surface these pages stand on, level 0, which is the one
	// answer a sheet has to give. It is not a status and no status hue
	// serves it; see the tokens package's highlight.go for the distances
	// that hold.
	{"highlight", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Highlight }},
	{"accent", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Primary }},
	{"on-accent", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.OnPrimary }},
	// The solid-fill state walk: hover one step from the pin toward the
	// ramp's 900 end, pressed two — SolidStateColor, the exact resolution a
	// filled button draws. They are emitted as first-class tokens because a
	// walked pin is off-ramp: no var() arithmetic over the ramp steps could
	// reproduce it, and a state is a real, addressable colour a sheet can
	// emit.
	{"accent-hover", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.SolidStateColor(tokens.RolePrimary, tokens.StateHover)
	}},
	{"accent-pressed", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.SolidStateColor(tokens.RolePrimary, tokens.StatePressed)
	}},
	{"secondary", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Secondary }},
	{"on-secondary", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.OnSecondary }},
	{"tertiary", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Tertiary }},
	{"on-tertiary", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.OnTertiary }},
	{"error", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Error }},
	{"on-error", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.OnError }},
	{"success", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Success }},
	{"on-success", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.OnSuccess }},
	{"warning", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Warning }},
	{"on-warning", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.OnWarning }},
	{"info", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Info }},
	{"on-info", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.OnInfo }},
	// The status containers and the marks read on them. They are emitted as
	// first-class tokens for the same reason the state walk above is: a
	// container is realized at a tone rather than mixed, so no var()
	// arithmetic over the ramp steps could reproduce one, and the mark is
	// the step the container's own contrast chose, which a sheet has no way
	// to measure.
	{"error-container", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.StatusContainer(tokens.RoleError)
	}},
	{"on-error-container", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.OnStatusContainer(tokens.RoleError)
	}},
	{"success-container", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.StatusContainer(tokens.RoleSuccess)
	}},
	{"on-success-container", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.OnStatusContainer(tokens.RoleSuccess)
	}},
	{"warning-container", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.StatusContainer(tokens.RoleWarning)
	}},
	{"on-warning-container", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.OnStatusContainer(tokens.RoleWarning)
	}},
	{"info-container", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.StatusContainer(tokens.RoleInfo)
	}},
	{"on-info-container", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.OnStatusContainer(tokens.RoleInfo)
	}},
	// Each status role's mark on the inverse surface: the step of that
	// role's ramp nearest its mid-value step that reads over the
	// counterpart scheme's card at the on-colour floor (MarkOn). It is a
	// token rather than a ramp reference because the two schemes do not
	// land on one step — a light scheme's marks come off step 500 and a
	// dark scheme's off step 400, its ramps having turned light by 500 —
	// so a sheet naming a step could not flip them with the scheme.
	{"error-on-inverse", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.MarkOn(tokens.RoleError, t.InverseSurface, onFloor)
	}},
	{"success-on-inverse", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.MarkOn(tokens.RoleSuccess, t.InverseSurface, onFloor)
	}},
	{"warning-on-inverse", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.MarkOn(tokens.RoleWarning, t.InverseSurface, onFloor)
	}},
	{"info-on-inverse", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.MarkOn(tokens.RoleInfo, t.InverseSurface, onFloor)
	}},
	// A badge draws one hue at two strengths, so each of the five variants
	// emits two tokens: the container fill it wears and the foreground read
	// on that fill. Both are resolved for the surface the sheet's pages stand
	// on, which is level 0 — a badge is small and its fill is
	// derived against whatever it is placed on, and a sheet has one answer to
	// give.
	//
	// The fill is a pale tint of the role's hue against that surface, at the
	// container chroma, at the depth that separates it from the surface
	// (StatusContainerOn). The foreground is then derived against the FILL
	// rather than against the surface: the role's pinned base while that base
	// clears the text floor over it and the nearest step to the mid-value
	// that does otherwise (InkOn). Neutral has no pinned base — the neutral
	// ramp carries no pin — so it takes the walk directly (MarkOn), at the
	// same floor, and its fill comes back as depth alone.
	//
	// The floor is the text floor for all five: a sign is the same utterance
	// at the same weight as a word. Never an inverted on-colour — a white
	// word on a saturated field is the variant interaction speaks in, and a
	// badge is read rather than used.
	{"badge-neutral-fill", badgeFill(tokens.RoleNeutral)},
	{"badge-neutral", badgeForeground(tokens.RoleNeutral)},
	{"badge-success-fill", badgeFill(tokens.RoleSuccess)},
	{"badge-success", badgeForeground(tokens.RoleSuccess)},
	{"badge-warning-fill", badgeFill(tokens.RoleWarning)},
	{"badge-warning", badgeForeground(tokens.RoleWarning)},
	{"badge-error-fill", badgeFill(tokens.RoleError)},
	{"badge-error", badgeForeground(tokens.RoleError)},
	{"badge-info-fill", badgeFill(tokens.RoleInfo)},
	{"badge-info", badgeForeground(tokens.RoleInfo)},
	// The tinted button is that same recipe under the accent role, so its
	// six tokens are the badge's two taken through the three states a
	// button answers a pointer with. They are tokens rather than ramp
	// references for the reason the containers above are: the fill is
	// realized at a tone against the surface it stands on, and the walk is
	// counted on the neutral scale from that realization, so no var()
	// arithmetic over the ramp steps reproduces either.
	//
	// Level 0, as the badge family is, and for the same reason: the Gio
	// side derives against the surface the control is placed on
	// (RenderState.Level) and a sheet has one answer to give.
	{"btn-tonal-fill", tonalFill(tokens.StateNormal)},
	{"btn-tonal", tonalForeground(tokens.StateNormal)},
	{"btn-tonal-fill-hover", tonalFill(tokens.StateHover)},
	{"btn-tonal-hover", tonalForeground(tokens.StateHover)},
	{"btn-tonal-fill-active", tonalFill(tokens.StatePressed)},
	{"btn-tonal-active", tonalForeground(tokens.StatePressed)},
	// The marks a control and a raised surface draw on themselves, each the
	// step its own ramp's MarkOn walk answers with at the graphic floor. All
	// are per-scheme tokens rather than named steps, because a named step is
	// a pairing and not a colour: the light and dark neutral ramps are
	// realized at the same perceptual depths from opposite ends, so one step
	// means two different contrasts against two surfaces that moved the whole
	// way.
	//
	// The two families ask the elevation levels different questions, and the
	// difference is the whole of why one is a set of four and the other a
	// single token. A resting edge asks which step of the neutral ramp reads
	// on the level the thing stands on, and each level may answer for
	// itself: an edge is the boundary of one surface, and two surfaces are
	// free to draw their own. A focus ring asks which step of the primary
	// ramp reads on EVERY level at once, because focus is one state and a
	// page that spelled it in two purples would be teaching two idioms.
	//
	// control-border is the level-0 answer for the row of controls that
	// says what it is with a line — the unchecked box, the unselected radio,
	// the text field, the dropdown trigger (components/input controlBorder):
	// the neutral step nearest step 500 that reaches 3:1 against the level-0
	// surface a control on the page is guaranteed against. Naming step 500 in
	// both schemes — which this sheet did, at every one of those four sites —
	// measured 6.63:1 in the dark and 2.67:1 in the light, under the floor in
	// the scheme most people read in. The walk answers 600 in the light
	// scheme and 500 in the dark and needs to know nothing about either.
	//
	// dialog-border and popover-border are the same walk taken against a
	// deeper level, and each serves both readings of "edge on that level":
	// the surface's own outline — a dialog's edge circles its level-2 fill,
	// a popover's its level 3, each pattern painting the fill it is measured
	// against — and the resting edge of any control standing on it, which is
	// the same line over the same surface and cannot sensibly be a second
	// colour. A checkbox in a dialog therefore takes dialog-border, not
	// control-border, and asks its own question rather than inheriting an
	// answer. A focus ring asks nothing of the level it was put on, which is
	// why no dialog-focus-ring stands beside dialog-border here.
	//
	// Level 1 has no member: a card is never outlined and draws no line of
	// its own, and the light scheme's level-0 step already clears level 1,
	// so a control standing on a card takes control-border unchanged.
	//
	// Whether the four answers differ is the derivation's to report, and
	// today they part in one scheme only. Because elevation lightens toward
	// the viewer in both schemes, a light window's hardest surface is its
	// CHROME level and a dark window's is its TOP level, and a step that
	// clears the hardest clears every other by more. In the light scheme one
	// neutral step therefore serves the whole window — 3.55:1 on chrome,
	// rising to 4.35:1 on a popover — and all four tokens repeat it. The dark
	// scheme's levels climb further from its chrome level: its level-0 step
	// reads 2.62:1 on a dialog's fill and 1.80:1 on a popover's, under the
	// floor a graphic owes the surface it stands on, so those two levels walk
	// on to a lighter step and the sheet states two edge colours where the
	// light one states one.
	//
	// focus-ring is the scheme's one ring, the colour every focused control
	// draws on every level: focusRing below, the step of the primary ramp
	// nearest its mid-value step that reaches the graphic floor against all
	// the levels a control can stand on at once. One token, not one per
	// surface, because a walk aimed at one surface answers that surface: two
	// controls whose fills lie three units apart on one level come back steps
	// 19 L* apart when the ramp carries a step between them, and two purples
	// for one state on one page is not an idiom. Asking every level is
	// affordable because a ring only ever lies on a level elevation carries,
	// and there are five of those rather than the whole scheme.
	//
	// The same pick answers the four edge tokens above as well, and owes them
	// a separation rather than a floor. A ring is a graphic on a surface, so
	// the levels are what the graphic floor is measured to; but the line the
	// ring replaces is the level's own resting edge, and a ring that matched
	// that edge in luminance would announce focus in hue alone —
	// focusRingBorderSeparation is what keeps the two apart in the channel a
	// forced-colors or greyscale display leaves standing.
	//
	// One surface belongs to no level: the accent fill a FILLED button's
	// ring lies on, because that ring is inset in the button's own
	// background rather than drawn at its boundary. It is the one place the
	// scheme's ring cannot be used — the ring is a step of the primary ramp
	// and so is that fill, and the two land on the same step, which is the
	// same colour twice. focusRingOn keeps the scheme's ring wherever it
	// reads on the fill and walks against the fill only where it cannot, so
	// the exception costs exactly the one surface that forces it.
	{"control-border", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.MarkOn(tokens.RoleNeutral, t.SurfaceAt(tokens.Level0), graphicFloor)
	}},
	{"dialog-border", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.MarkOn(tokens.RoleNeutral, t.SurfaceAt(tokens.Level2), graphicFloor)
	}},
	{"popover-border", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.MarkOn(tokens.RoleNeutral, t.SurfaceAt(tokens.Level3), graphicFloor)
	}},
	{"focus-ring", focusRing},
	{"focus-ring-on-accent", func(t tokens.ColorTokens) stdcolor.NRGBA {
		return focusRingOn(t, t.SolidStateColor(tokens.RolePrimary, tokens.StateFocus))
	}},
}

// primaryMidStep indexes the primary ramp's step 500, the mid-value step the
// ring's pick is aimed at. A tokens.Ramp is nine steps, 100 through 900.
const primaryMidStep = 4

// focusRingBorderSeparation is the least luminance separation the ring owes
// the neutral resting border a control on the same level draws — the line
// control-border, dialog-border and popover-border carry, and the
// line a focused field swaps for its ring. Colour is the ring's only channel,
// so a ring at the border's own luminance says nothing but hue, and hue is
// what Differentiate Without Color, forced-colors and a greyscale display take
// away.
//
// 1.25:1 is measured rather than picked. Over the seed sweep — 411 seeds, both
// schemes, both derivations, every level — the separations a step of the
// primary ramp can reach while still clearing graphicFloor fall in two bands
// with a wide empty stretch between them: 1.00–1.01, where the ring and the
// border are one grey, and 1.53 upward, where the ramp's next step is a
// different grey. The threshold goes in the empty stretch.
//
// It is components/internal/focus.BorderSeparation, restated here for the same
// reason the walk below is.
const focusRingBorderSeparation = 1.25

// focusRing is the colour every focused control draws its ring in, one per
// scheme: the step of the primary ramp nearest primaryMidStep that reaches
// graphicFloor against every elevation level, reaches
// focusRingBorderSeparation against every level's neutral resting border, and
// is not the accent fill itself. It is the
// derivation components/internal/focus draws by, restated here because the
// sheet is emitted a layer below the components and the two must land on the
// same hex.
//
// Every level rather than one, because the ring is one colour
// and a control may stand anywhere on it: a chip on a card and the button
// beside it are the same state and owe the reader the same pixel. The ramp
// is walked from its middle out, so where several steps clear every level the
// ring is the one nearest the depth the brand hue is most itself at, and the
// one furthest from both ends.
//
// The accent fill is excluded rather than measured, because what it owes the
// ring has no scale: it is what a checked box and a filled button paint at
// rest, and a dark scheme realizes it exactly on a step of this ramp. A ring
// drawn in it would announce focus in the colour the control was already
// speaking.
//
// A ring has to be drawn whatever it measures, so a palette no step satisfied
// all three on takes the step that comes closest against the levels rather
// than none.
func focusRing(t tokens.ColorTokens) stdcolor.NRGBA {
	pick, dist := -1, len(t.Ramps.Primary)
	widest, widestAt := -1.0, 0
	for i, step := range t.Ramps.Primary {
		// maxContrast is the ceiling of the WCAG ratio — black on white —
		// so the first level always lowers it.
		const maxContrast = 21.0
		worst, worstBorder := maxContrast, maxContrast
		for _, lvl := range standableLevels {
			surface := t.SurfaceAt(lvl.level)
			if got := vgcolor.ContrastRatio(step, surface); got < worst {
				worst = got
			}
			border := t.MarkOn(tokens.RoleNeutral, surface, graphicFloor)
			if got := vgcolor.ContrastRatio(step, border); got < worstBorder {
				worstBorder = got
			}
		}
		if worst > widest {
			widest, widestAt = worst, i
		}
		if worst < graphicFloor || worstBorder < focusRingBorderSeparation || step == t.Primary {
			continue
		}
		d := i - primaryMidStep
		if d < 0 {
			d = -d
		}
		if d < dist {
			pick, dist = i, d
		}
	}
	if pick < 0 {
		return t.Ramps.Primary[widestAt]
	}
	return t.Ramps.Primary[pick]
}

// focusRingOn is focusRing for a band lying inside a fill of the control's
// own, with that fill on both sides of it and no level anywhere near it —
// the filled button's inset ring, the only such band in the class layer. It
// answers the scheme's ring wherever that ring reads on the fill, and walks
// the primary ramp against the fill only where it cannot.
//
// A transparent fill is no fill: what a ghost button's ring lies on is the
// level showing through it, which the scheme's ring already answers.
func focusRingOn(t tokens.ColorTokens, fill stdcolor.NRGBA) stdcolor.NRGBA {
	ring := focusRing(t)
	if fill.A == 0 || vgcolor.ContrastRatio(ring, fill) >= graphicFloor {
		return ring
	}
	return t.MarkOn(tokens.RolePrimary, fill, graphicFloor)
}

// badgeFill and badgeForeground are the badge pair, written once per role
// rather than ten times. Splitting them is what keeps the two derivations
// honest: the foreground's surface is the fill, so a fill that moved without
// the foreground moving with it would emit a pairing nothing measured.
func badgeFill(role tokens.Role) func(tokens.ColorTokens) stdcolor.NRGBA {
	return func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.StatusContainerOn(role, t.SurfaceAt(tokens.Level0))
	}
}

func badgeForeground(role tokens.Role) func(tokens.ColorTokens) stdcolor.NRGBA {
	return func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.ForegroundOn(role, t.StatusContainerOn(role, t.SurfaceAt(tokens.Level0)))
	}
}

// tonalFill and tonalForeground are the SAME recipe under the accent role,
// which is what a tinted button wears: a tinted button and a status badge
// differ by no practical visual difference, so they speak one recipe and
// behaviour tells them apart. state is the walk the fill is under — normal,
// hover or pressed — and the foreground is derived against wherever that walk
// landed rather than against the resting fill.
func tonalFill(state tokens.State) func(tokens.ColorTokens) stdcolor.NRGBA {
	return func(t tokens.ColorTokens) stdcolor.NRGBA {
		rest := t.StatusContainerOn(tokens.RolePrimary, t.SurfaceAt(tokens.Level0))
		return t.PinnedStateColor(rest, state)
	}
}

func tonalForeground(state tokens.State) func(tokens.ColorTokens) stdcolor.NRGBA {
	return func(t tokens.ColorTokens) stdcolor.NRGBA {
		return t.ForegroundOn(tokens.RolePrimary, tonalFill(state)(t))
	}
}

// onFloor is WCAG AA for body text, the floor a mark on the inverse surface
// is chosen against: a toast's leading edge is the only thing that says
// which level the toast is, so it is held to the text floor rather than to
// the 3:1 a non-text graphic owes the surface it stands on. The floor does
// not bind — over the seed sweep, 3.0 picks the same steps — it states what
// the mark owes.
const onFloor = 4.5

// graphicFloor is WCAG 1.4.11's floor for a graphic that carries meaning
// without being text — 3:1 — the floor components/input measures a
// checkbox's edge to and components/internal/focus measures every focus
// ring to. A control's edge and its ring are the whole of what says which
// control it is and where the keyboard is, so neither is decoration.
const graphicFloor = 3.0

// typeRoles orders the fifteen MD3 type roles under their CSS names, plus
// code — the sixteenth style outside the MD3 grid, the mono face at
// body-medium's metrics (G-F0) — emitted last.
var typeRoles = []struct {
	name string
	pick func(tokens.Typography) tokens.TextStyle
}{
	{"display-large", func(t tokens.Typography) tokens.TextStyle { return t.DisplayLarge }},
	{"display-medium", func(t tokens.Typography) tokens.TextStyle { return t.DisplayMedium }},
	{"display-small", func(t tokens.Typography) tokens.TextStyle { return t.DisplaySmall }},
	{"headline-large", func(t tokens.Typography) tokens.TextStyle { return t.HeadlineLarge }},
	{"headline-medium", func(t tokens.Typography) tokens.TextStyle { return t.HeadlineMedium }},
	{"headline-small", func(t tokens.Typography) tokens.TextStyle { return t.HeadlineSmall }},
	{"title-large", func(t tokens.Typography) tokens.TextStyle { return t.TitleLarge }},
	{"title-medium", func(t tokens.Typography) tokens.TextStyle { return t.TitleMedium }},
	{"title-small", func(t tokens.Typography) tokens.TextStyle { return t.TitleSmall }},
	{"label-large", func(t tokens.Typography) tokens.TextStyle { return t.LabelLarge }},
	{"label-medium", func(t tokens.Typography) tokens.TextStyle { return t.LabelMedium }},
	{"label-small", func(t tokens.Typography) tokens.TextStyle { return t.LabelSmall }},
	{"body-large", func(t tokens.Typography) tokens.TextStyle { return t.BodyLarge }},
	{"body-medium", func(t tokens.Typography) tokens.TextStyle { return t.BodyMedium }},
	{"body-small", func(t tokens.Typography) tokens.TextStyle { return t.BodySmall }},
	{"code", func(t tokens.Typography) tokens.TextStyle { return t.Code }},
}

// spaceKeys orders the spacing stops under the Go scale's own key names.
var spaceKeys = []struct {
	name string
	pick func(tokens.SpacingScale) float32
}{
	{"0", func(s tokens.SpacingScale) float32 { return s.S0 }},
	{"1", func(s tokens.SpacingScale) float32 { return s.S1 }},
	{"2", func(s tokens.SpacingScale) float32 { return s.S2 }},
	{"3", func(s tokens.SpacingScale) float32 { return s.S3 }},
	{"4", func(s tokens.SpacingScale) float32 { return s.S4 }},
	{"5", func(s tokens.SpacingScale) float32 { return s.S5 }},
	{"6", func(s tokens.SpacingScale) float32 { return s.S6 }},
	{"8", func(s tokens.SpacingScale) float32 { return s.S8 }},
	{"10", func(s tokens.SpacingScale) float32 { return s.S10 }},
	{"12", func(s tokens.SpacingScale) float32 { return s.S12 }},
	{"16", func(s tokens.SpacingScale) float32 { return s.S16 }},
	{"20", func(s tokens.SpacingScale) float32 { return s.S20 }},
	{"24", func(s tokens.SpacingScale) float32 { return s.S24 }},
}

// radiusKeys orders the radius stops under Tailwind's names, which the Go
// field names mirror (Xl2 and Xl3 are Go spellings of 2xl and 3xl).
var radiusKeys = []struct {
	name string
	pick func(tokens.RadiusScale) float32
}{
	{"none", func(r tokens.RadiusScale) float32 { return r.None }},
	{"sm", func(r tokens.RadiusScale) float32 { return r.Sm }},
	{"base", func(r tokens.RadiusScale) float32 { return r.Base }},
	{"md", func(r tokens.RadiusScale) float32 { return r.Md }},
	{"lg", func(r tokens.RadiusScale) float32 { return r.Lg }},
	{"xl", func(r tokens.RadiusScale) float32 { return r.Xl }},
	{"2xl", func(r tokens.RadiusScale) float32 { return r.Xl2 }},
	{"3xl", func(r tokens.RadiusScale) float32 { return r.Xl3 }},
	{"full", func(r tokens.RadiusScale) float32 { return r.Full }},
}

// elevationLevels orders the levels from the backdrop up toward the reader:
// the two levels under the content under their own names, then the four
// numbered levels. Each level's fill is resolved per scheme through
// [tokens.ColorTokens.SurfaceAt] and its shadow dp read off the snapshot's
// ElevationScale.
//
// The backdrop and the chrome level are spelled out rather than numbered
// because the numbering is anchored on the content and they are below it:
// naming them "-2" and "-1" in a CSS variable would read as an arithmetic
// accident, and renumbering the four above them would rename every token to
// say the same thing.
var elevationLevels = []struct {
	name  string
	level tokens.ElevationLevel
}{
	{"backdrop", tokens.LevelBackdrop},
	{"chrome", tokens.LevelChrome},
	{"0", tokens.Level0},
	{"1", tokens.Level1},
	{"2", tokens.Level2},
	{"3", tokens.Level3},
}

// standableLevels is elevationLevels without the backdrop: the levels a
// control can be put on, and so the surfaces a derivation that answers "on
// every level" has to clear. Nothing is drawn at the backdrop — it shows
// wherever nothing stands — so a ring measured against it would be walked
// against a surface no ring ever lies on.
var standableLevels = elevationLevels[1:]

// densityMetrics orders the per-setting density metrics under their CSS
// names. The WCAG pointer-target floor is not here: it is not a per-setting
// metric — see densityVars.
var densityMetrics = []struct {
	name string
	pick func(tokens.Density) float32
}{
	{"control-height", func(d tokens.Density) float32 { return d.ControlHeight }},
	// The chip height is emitted rather than left to calc() over the control
	// height: the sheet states resolved values everywhere else, and a
	// var() subtraction would put the system's one statement of the relation
	// in a stylesheet instead of in the token layer.
	{"chip-height", func(d tokens.Density) float32 { return d.ChipHeight() }},
	{"padding-x", func(d tokens.Density) float32 { return d.PaddingX }},
	{"padding-y", func(d tokens.Density) float32 { return d.PaddingY }},
}

// easeRoles orders the MD3 easing presets under their CSS names.
var easeRoles = []struct {
	name string
	pick func(tokens.MotionScale) tokens.Bezier
}{
	{"standard", func(m tokens.MotionScale) tokens.Bezier { return m.EaseStandard }},
	{"standard-accelerate", func(m tokens.MotionScale) tokens.Bezier { return m.EaseStandardAccelerate }},
	{"standard-decelerate", func(m tokens.MotionScale) tokens.Bezier { return m.EaseStandardDecelerate }},
	{"emphasized", func(m tokens.MotionScale) tokens.Bezier { return m.EaseEmphasized }},
	{"emphasized-accelerate", func(m tokens.MotionScale) tokens.Bezier { return m.EaseEmphasizedAccelerate }},
	{"emphasized-decelerate", func(m tokens.MotionScale) tokens.Bezier { return m.EaseEmphasizedDecelerate }},
}

// durationStops orders the duration stops under their CSS names.
var durationStops = []struct {
	name string
	pick func(tokens.MotionScale) time.Duration
}{
	{"x-fast", func(m tokens.MotionScale) time.Duration { return m.DurXFast }},
	{"fast", func(m tokens.MotionScale) time.Duration { return m.DurFast }},
	{"normal", func(m tokens.MotionScale) time.Duration { return m.DurNormal }},
	{"slow", func(m tokens.MotionScale) time.Duration { return m.DurSlow }},
	{"x-slow", func(m tokens.MotionScale) time.Duration { return m.DurXSlow }},
}

// boxShadow approximates an elevation depth as a CSS box-shadow: y-offset
// the level's dp, blur twice it, no spread, black at 20%. Depth 0 casts no
// shadow at all, so it is "none" rather than an invisible shadow.
func boxShadow(dp float32) string {
	if dp == 0 {
		return "none"
	}
	return fmt.Sprintf("0 %s %s 0 rgba(0, 0, 0, 0.2)", px(dp), px(2*dp))
}

// cubicBezier renders a Bezier as the CSS cubic-bezier() function.
func cubicBezier(bz tokens.Bezier) string {
	return fmt.Sprintf("cubic-bezier(%s, %s, %s, %s)",
		fnum(bz.P1[0]), fnum(bz.P1[1]), fnum(bz.P2[0]), fnum(bz.P2[1]))
}

// ms formats a duration as CSS milliseconds.
func ms(d time.Duration) string {
	return strconv.FormatFloat(float64(d)/float64(time.Millisecond), 'f', -1, 64) + "ms"
}

// colorVars renders one colour scheme as its ramp and pin variables.
func colorVars(t tokens.ColorTokens) []cssVar {
	var vars []cssVar
	for _, role := range rampRoles {
		ramp := role.ramp(t.Ramps)
		for step := 100; step <= 900; step += 100 {
			vars = append(vars, cssVar{
				name:  fmt.Sprintf("--color-%s-%d", role.name, step),
				value: hexRGB(ramp.Step(step)),
			})
		}
	}
	for _, pin := range pinRoles {
		vars = append(vars, cssVar{"--color-" + pin.name, hexRGB(pin.pick(t))})
	}
	// The elevation levels' surface fills. They live with the colours
	// rather than with the mode-invariant scales because a level is not a
	// ramp step in both schemes: the levels are anchored on the Background
	// pin and placed in CIELAB L*, so the light scheme's levels above the
	// content are off the ramp and the dark scheme's backdrop is off it
	// below. No var() arithmetic over the ramp steps reaches those values,
	// so each scheme states its own, exactly as the walked pins and the
	// derived borders beside them do.
	//
	// --elevation-1 is not a table entry on the Go side: it is the raise
	// walked from the content. The sheet states it anyway because the walk
	// has no CSS arithmetic either, and the class layer expresses the walk
	// the way a cascade can — a host that raises what it holds redeclares
	// --surface-raised, so a control names var(--surface-raised,
	// var(--elevation-1)) once and lands one step above whatever it is
	// actually inside.
	for _, level := range elevationLevels {
		vars = append(vars, cssVar{"--elevation-" + level.name, hexRGB(t.SurfaceAt(level.level))})
	}
	// And the seam each level owes what stands on it: the hairline a raise
	// draws at its own edge where the scheme has no step left to tell it
	// with. It is `transparent` where the raise IS told by its fill, so a
	// rule can carry the border unconditionally and the geometry does not
	// move between the schemes — which is what the Gio side does too, its
	// stroke being centred on an edge the inset does not depend on.
	//
	// The backdrop has none: nothing stands on the backdrop.
	for _, level := range standableLevels {
		raise := t.RaisedOn(t.SurfaceAt(level.level))
		value := "transparent"
		if raise.Seamed {
			value = hexRGB(raise.Seam)
		}
		vars = append(vars, cssVar{"--elevation-" + level.name + "-seam", value})
	}
	// And the hairline two regions that SHARE a level's fill are parted by:
	// what a group draws at its own edge. It is the same derivation as the
	// seam above with both sides at one fill
	// (tokens.ColorTokens.SeamOn), and it differs from that seam in when it
	// is drawn: a raise owes its seam only where the fill cannot tell the
	// raise, so --elevation-N-seam is `transparent` in every other scheme,
	// while a group has no fill of its own and the line is the whole of what
	// says where it ends. It is never transparent.
	//
	// The backdrop has none: nothing is grouped on the bare window plane.
	for _, level := range standableLevels {
		vars = append(vars, cssVar{"--elevation-" + level.name + "-hairline",
			hexRGB(t.SeamOn(t.SurfaceAt(level.level)))})
	}
	// And each level's own interaction walk, for the same reason and one step
	// further: a ghost button paints no fill at rest and takes a state fill on
	// the surface it stands on under the pointer, so that fill is taken FROM
	// that level's fill (tokens.ColorTokens.StateAt, which is what
	// components/button's ghostWash performs). While a level was a ramp step
	// the sheet could name the step's neighbour and be done; a level off the
	// ramp has no neighbour to name, so the walk is written out per scheme
	// like the fill it starts from.
	for _, level := range elevationLevels {
		for _, st := range []struct {
			suffix string
			state  tokens.State
		}{
			{"-hover", tokens.StateHover},
			{"-active", tokens.StatePressed},
		} {
			vars = append(vars, cssVar{
				name:  "--elevation-" + level.name + st.suffix,
				value: hexRGB(t.StateAt(level.level, st.state)),
			})
		}
	}
	return vars
}

// scaleVars renders the mode-invariant families: fonts, density
// (comfortable — the :root setting), spacing, radius, the dp shadows (the
// opt-in cue for floating transients), and the motion set. The tonal
// surface fills the shadows layer over are NOT here: a level resolves per
// scheme, so --elevation-* sits with the colours.
func scaleVars(s Snapshot) []cssVar {
	vars := []cssVar{
		{"--font-family", strconv.Quote(s.Typography.BodyLarge.Typeface)},
		{"--font-family-code", strconv.Quote(s.Typography.Code.Typeface)},
	}
	for _, role := range typeRoles {
		style := role.pick(s.Typography)
		vars = append(vars,
			cssVar{"--font-" + role.name + "-size", px(style.Size)},
			cssVar{"--font-" + role.name + "-line-height", px(style.LineHeight)},
			cssVar{"--font-" + role.name + "-weight", strconv.Itoa(style.Weight)},
			cssVar{"--font-" + role.name + "-tracking", px(style.Tracking)},
		)
	}
	vars = append(vars, densityVars(tokens.Comfortable)...)
	vars = append(vars, cssVar{"--density-min-hit-target", px(tokens.Comfortable.MinHitTarget())})
	for _, key := range spaceKeys {
		vars = append(vars, cssVar{"--space-" + key.name, px(key.pick(s.Spacing))})
	}
	for _, key := range radiusKeys {
		vars = append(vars, cssVar{"--radius-" + key.name, px(key.pick(s.Radius))})
	}
	for _, level := range elevationLevels {
		vars = append(vars, cssVar{"--shadow-" + level.name, boxShadow(s.Elevation.Dp(level.level))})
	}
	for _, role := range easeRoles {
		vars = append(vars, cssVar{"--ease-" + role.name, cubicBezier(role.pick(s.Motion))})
	}
	for _, stop := range durationStops {
		vars = append(vars, cssVar{"--duration-" + stop.name, ms(stop.pick(s.Motion))})
	}
	// The interaction-state base the class layer builds on: the ring's 2 dp
	// stroke width, and the disabled fraction as
	// tokens.DisabledOpacity in color-mix() percent, because disabled is an
	// opacity and not a ramp step. Both are mode-invariant, which is why they
	// are here and the ring's COLOUR is not: --color-focus-ring is a measured
	// walk against surfaces that flip with the scheme, so it lives with the
	// colours (see pinRoles).
	vars = append(vars,
		cssVar{"--focus-ring-width", px(focusRingWidthDp)},
		cssVar{"--state-disabled-opacity", fnum(tokens.DisabledOpacity*100) + "%"},
	)
	// The scrim: a modal's dimmer over the whole window plane — black at alpha
	// 0x80, deliberately the same in both modes because a scrim dims by
	// reducing luminance, so it lives with the mode-invariant scales rather
	// than in the colour schemes. Like the shadows' fixed black, it is a
	// constant of the pattern, not a ramp resolution; emitting it as a token
	// keeps the class layer itself literal-free.
	vars = append(vars, cssVar{"--color-scrim", scrimRGBA})
	return vars
}

// scrimRGBA is the scrim colour — color.NRGBA{0, 0, 0, 0x80} — as the CSS
// colour that REPRODUCES it, which is not rgba(0,0,0,0.502): Gio
// composites the translucent black in linear RGB while a browser composites
// plain-alpha backgrounds in the sRGB space the pixels are stored in, so the
// literal alpha would dim roughly twice as hard as the pattern does
// (measured: 123 vs Gio's 181 over the light bg pin).
// The sRGB-equivalent alpha — the a solving srgb(bg)·(1−a) =
// srgb(linear(bg)·0.5) — is 0.267 over the light surfaces (bg ≈ 247), 0.28 at
// mid-grey and 0.30 near black: 0.28 is the compromise, within ±0.013 of
// exact across the whole tonal range (≤ ~3/255 per channel on any surface),
// and one value serves both modes exactly as the Gio constant does.
const scrimRGBA = "rgba(0, 0, 0, 0.28)"

// focusRingWidthDp is the focus ring's stroke width — the 2 dp
// components/button draws (drawButton's gtx.Dp(2) stroke), identical in
// every emphasis because keyboard visibility is not a matter of
// prominence.
const focusRingWidthDp = 2

// densityVars renders one density setting's per-setting metrics. The :root
// block carries tokens.Comfortable's; the .compact override block carries
// tokens.Compact's. --density-min-hit-target is deliberately not among
// them: the WCAG 2.5.5 pointer-target floor does not scale with density, so
// it is emitted once in :root and never overridden — the CSS mirror of
// Density.MinHitTarget being a method, not a field.
func densityVars(d tokens.Density) []cssVar {
	var vars []cssVar
	for _, m := range densityMetrics {
		vars = append(vars, cssVar{"--density-" + m.name, px(m.pick(d))})
	}
	return vars
}

// block renders one selector's declarations.
// kindOther is the Claude Design pane's kind marker for custom properties
// its token classifier cannot type on its own — the easing curves and
// durations. Without the marker the pane's self-check re-adds it by hand on
// every pass and the next regeneration wipes it again; emitting it here is
// the durable half of that handshake.
const kindOther = "/* @kind other */"

func block(b *strings.Builder, selector string, vars []cssVar) {
	b.WriteString(selector)
	b.WriteString(" {\n")
	for _, v := range vars {
		if strings.HasPrefix(v.name, "--ease-") || strings.HasPrefix(v.name, "--duration-") {
			fmt.Fprintf(b, "  %s: %s; %s\n", v.name, v.value, kindOther)
		} else {
			fmt.Fprintf(b, "  %s: %s;\n", v.name, v.value)
		}
	}
	b.WriteString("}\n")
}

// stylesCSS renders the full token sheet: the light scheme and every
// mode-invariant scale under :root, the paired dark colours under .dark,
// and the compact density metrics under .compact. The two class blocks are
// orthogonal switches — .dark flips the colours (and with them the
// var()-chained --elevation-* surfaces), .compact flips the per-setting
// density metrics — so a surface can be any of the four combinations.
func stylesCSS(s Snapshot) string {
	var b strings.Builder
	b.WriteString("/* Generated by theme/export (cmd/vg-tokens). Do not edit. */\n\n")
	// The faces behind --font-family and --font-family-code, self-hosted in
	// the bundle's fonts/ directory: the same Roboto regular and medium the
	// Gio applications embed, and the font repo's Roboto Mono. Without these
	// rules the design surface renders every specimen in a substitute face —
	// claude.ai/design flags the family as a missing brand font — and a
	// mirror scored against substitute-shaped text measures the machine, not
	// the mirror. Weights 400 and 500 are the only ones the token sheet uses.
	fontFace := func(family, weight, file string) {
		fmt.Fprintf(&b, "@font-face {\n  font-family: %s;\n  font-weight: %s;\n  font-style: normal;\n  src: url(\"fonts/%s\") format(\"truetype\");\n}\n", strconv.Quote(family), weight, file)
	}
	fontFace("Roboto", "400", "roboto-regular.ttf")
	fontFace("Roboto", "500", "roboto-medium.ttf")
	fontFace("Roboto Mono", "400", "robotomono-regular.ttf")
	b.WriteString("\n")
	block(&b, ":root", append(colorVars(s.Light), scaleVars(s)...))
	b.WriteString("\n")
	block(&b, ".dark", colorVars(s.Dark))
	b.WriteString("\n")
	block(&b, ".compact", densityVars(tokens.Compact))
	b.WriteString("\n")
	b.WriteString(componentClasses)
	return b.String()
}

// componentClasses is the class layer: the component vocabulary the design surface composes screens from, defined
// over the token variables above — not one literal colour, so a re-branded
// sheet re-brands the components with it, and .dark/.compact flip them like
// everything else. The only literal lengths are the same component
// constants the Gio side hardcodes as dp literals rather than tokens
// (checkbox/radio's 20 dp glyph and 10 dp dot, the dropdown's 16 dp
// chevron, the 1/2 dp input borders); each is commented at its source.
//
// It mirrors the button and input components, the sources of truth: .btn is
// the filled variant by default, .tonal and .ghost the emphasis
// modifiers, and every state resolves exactly as buttonColors resolves it
// (button.go: the filled fill walks via SolidStateColor into
// --color-accent-hover/-pressed; tonal takes the shared tint tokens above,
// which are the badge's own recipe under the accent role; ghost is
// ghostText 700 / ghostTextOnWash 900 over the local surface's own walk).
// A ghost's state fill derives from the local surface, so the raised hosts
// carry contextual overrides walking from their own level's step
// (ghostGroundStep: the level-2 dialog takes 400/500, the level-3 popover
// 500/600), matching RenderState.Level on the Gio side. Tonal derives
// against that same level on the Gio side and the sheet states level 0, as
// the badge tokens do. Because each variant's blocks override every state
// it treats, later variant blocks never bleed a state from an earlier one;
// :disabled resolutions are per-variant for the same reason. Selected
// resolves where pressed does — the two-step walk. The form controls
// resolve as components/input does: the raised level under body text, the
// ramp's own measured answer on the text field, the radio and the checkbox
// alike — all of it taken against the level the control stands on, which is
// what --surface-border and --surface-raised carry down from a raised host
// — neutral 700 placeholder and glyph, focus promoting the border to the
// scheme's ring, disabled fading each colour to the disabled fraction of
// its alpha. The checked checkbox carries the check mark the Gio side
// strokes, drawn out of the icon set's grid as two gradient bands rather
// than encoded as an image, so the layer stays literal-free.
//
// Every pointer/keyboard state rule also carries a forcing twin class
// (.is-hover, .is-active, .is-focus, .is-checked) grouped into the same
// rule. A static component page cannot hover itself, and duplicating the
// declarations in the page would fork the resolution; a grouped selector
// emitted by this generator shares the exact declarations with the live
// pseudo-class, so a forced specimen provably renders as the live state.
// Disabled needs no twin: the pages force it with the native attribute.
const componentClasses = `/* ---- Component classes ----
   The class vocabulary, built only on the tokens above. .btn mirrors
   components/button: filled by default, .tonal and .ghost the emphasis
   modifiers, states resolved exactly as the Gio side resolves them.
   .input/.select/.checkbox/.radio mirror components/input,
   .badge the inline annotation components/badge draws (the plain category
   label and the four status roles; the close mark is a Gio interaction and
   has no class here), .card the patterns/card surface and .group the
   patterns/group hairline, .table the patterns/table grid, the navigation
   family — .navbar, .tabs, .sidebar, .crumbs — the four patterns of the same
   names, and the overlay family — .scrim/.dialog (patterns/modal), .popover,
   .tooltip, .toast — the transient surfaces. The focus ring is the same
   ring in every variant — one width, one hue, one measured floor against
   whatever surface it circles: keyboard visibility is not an emphasis
   property. Each state rule carries a forcing twin class (.is-hover,
   .is-active, .is-focus, .is-checked) so a static page can show the state
   with the very declarations the live pseudo-class applies. */

.btn {
  box-sizing: border-box;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  appearance: none;
  border: none;
  margin: 0;
  text-decoration: none;
  user-select: none;
  white-space: nowrap;
  cursor: pointer;
  min-height: var(--density-control-height);
  padding: var(--density-padding-y) var(--density-padding-x);
  border-radius: var(--radius-md);
  font-family: var(--font-family);
  font-size: var(--font-label-large-size);
  line-height: var(--font-label-large-line-height);
  font-weight: var(--font-label-large-weight);
  letter-spacing: var(--font-label-large-tracking);
  background: var(--color-accent);
  color: var(--color-on-accent);
}

/* Filled states: the solid fill walks from the pin toward the ramp's 900
   end (hover one step, pressed and selected two) — the walked stops are
   tokens, not mixes. */
.btn:hover, .btn.is-hover { background: var(--color-accent-hover); }
.btn.selected { background: var(--color-accent-pressed); }
.btn:active, .btn.is-active { background: var(--color-accent-pressed); }

/* Keyboard focus keeps the resting fill and adds the ring — a stroke
   centred on the control's edge, as the Gio side draws it. One width, one
   hue, one value: var(--color-focus-ring) is the step of the primary ramp
   nearest its mid-value step that reaches 3:1 against every level at once
   and parts from every level's resting border in luminance, so a control
   wears the same ring wherever it was put, no rule here asks where that is,
   and focus survives a display that removes the colour. The one surface
   that belongs to no level is the accent fill a filled button's ring lies
   on — that ring is inset in the button's own background rather than drawn
   at its boundary, and no step that reads against the page reads against
   that fill — so that one takes
   var(--color-focus-ring-on-accent) wherever the button stands, and it is
   the same ring in the same place at the same width. */
.btn:focus-visible, .btn.is-focus {
  outline: var(--focus-ring-width) solid var(--color-focus-ring-on-accent);
  outline-offset: calc(var(--focus-ring-width) / -2);
}
.btn.tonal:focus-visible, .btn.tonal.is-focus,
.btn.ghost:focus-visible, .btn.ghost.is-focus,
.checkbox:focus-visible, .checkbox.is-focus,
.radio:focus-visible, .radio.is-focus {
  outline: var(--focus-ring-width) solid var(--color-focus-ring);
  outline-offset: calc(var(--focus-ring-width) / -2);
}

/* --surface-border and --surface-raised are the level-local pair, and the
   whole of how a raised host hands its level down to the controls inside
   it. The first is the neutral step that reads on this level; the second is
   the level a control that fills a box of its own is raised TO from here —
   the step above the host's own.
   Both inherit, so a surface declares them once — in its
   own rule, beside the --elevation-N it fills with — and every control below
   it re-derives, with no descendant selector anywhere in this sheet. It is
   the sheet's spelling of RenderState.Level on the Gio side: what a control
   resolves against is handed down, never looked up.

   Nothing declares them at the root. A control that no raised host contains
   stands on the paper, and the fallback inside each var() IS the paper's
   answer — so the default is written once, at the point of use, and cannot
   drift from the tokens the rules already name.

   A focus ring is not among them, and that absence is the point: the ring is
   one colour per scheme, measured against every level at once, so a raised
   host has nothing to hand down. A resting edge is the boundary of one
   surface and may differ per level — in the dark scheme the level-0 neutral
   step reads 2.62:1 over a level-2 fill and 1.80:1 over a level-3 one, under
   the 3:1 a graphic owes the surface it stands on, which is why a checkbox in
   a dialog wears the edge the dialog's own outline wears. The light scheme's
   levels climb less far from its backdrop and its level-0 step already clears
   every level, so the handed-down token repeats there and nothing moves — the
   derivation reporting that nothing needs to.

   The two join under different rules, and the difference is
   worth stating. --surface-border is a MEASUREMENT, so a host only declares it
   where the paper's answer stops clearing — which is why .card, at level 1,
   does not. --surface-raised is a LEVEL, and a
   level differs by construction: a control filling at its host's own step is
   invisible against it whatever the contrast table says. So every
   raised surface a control can be put inside declares it — .card,
   .dialog and .popover. A .group declares neither: it takes the fill of the
   surface it is in and raises nothing, so a control inside a group asks the
   same questions it would ask standing on that surface bare. The
   popover declares the ceiling: at the top of the scheme's own range the
   walk clamps rather than stepping, so a control in a popover fills flush
   with it and is read by its border and by --elevation-3-seam. */

/* Disabled is an opacity, not a ramp step: each colour keeps its hue and
   fades to the disabled fraction of its alpha. */
.btn:disabled {
  cursor: default;
  background: color-mix(in srgb, var(--color-accent) var(--state-disabled-opacity), transparent);
  color: color-mix(in srgb, var(--color-on-accent) var(--state-disabled-opacity), transparent);
}

/* Tonal: the accent's tint over the surface the button stands on, under the
   accent's own colour at the text floor — the same recipe .badge wears, one
   hue at two strengths, never an inverted on-colour. The foreground moves
   with the fill under the pointer: a colour held over a fill that walked two
   steps is measured against a surface no longer there. Selected resolves
   where pressed does, the two-step walk. */
.btn.tonal {
  background: var(--color-btn-tonal-fill);
  color: var(--color-btn-tonal);
}
.btn.tonal:hover, .btn.tonal.is-hover {
  background: var(--color-btn-tonal-fill-hover);
  color: var(--color-btn-tonal-hover);
}
.btn.tonal.selected,
.btn.tonal:active, .btn.tonal.is-active {
  background: var(--color-btn-tonal-fill-active);
  color: var(--color-btn-tonal-active);
}
.btn.tonal:disabled {
  background: color-mix(in srgb, var(--color-btn-tonal-fill) var(--state-disabled-opacity), transparent);
  color: color-mix(in srgb, var(--color-btn-tonal) var(--state-disabled-opacity), transparent);
}

/* Ghost: no fill at rest — the neutral ramp's low-contrast text over
   whatever surface it sits on; under the pointer it performs that
   surface's own hover and press walk, the text walking to 900 with the
   fill. No selected treatment: a ghost stays the least pronounced.

   The walk is named as a level's own state rather than as a ramp step,
   because a level is not a ramp step: the fills above the
   pin are off the ramp in the light scheme and the backdrop is off it in the
   dark one, so there is no index left to walk from. Each level's own
   -hover and -active pair is that walk taken
   from each level's own fill (components/button ghostWash, which is
   tokens.ColorTokens.StateAt). A ghost told nothing stands on the paper,
   so the base rule is level 0's. */
.btn.ghost {
  background: transparent;
  color: var(--color-neutral-700);
}
.btn.ghost:hover, .btn.ghost.is-hover {
  background: var(--elevation-0-hover);
  color: var(--color-neutral-900);
}
.btn.ghost:active, .btn.ghost.is-active {
  background: var(--elevation-0-active);
  color: var(--color-neutral-900);
}
.btn.ghost:disabled {
  background: transparent;
  color: color-mix(in srgb, var(--color-neutral-700) var(--state-disabled-opacity), transparent);
}

/* A ghost's state fill derives from the local surface it sits on, not the
   window's own: inside a host that is not the paper, the hover and press
   fills re-derive as that host surface's own walk (components/button
   buttonColors, walking from RenderState.Level). The card sits at level 1,
   the dialog at level 2, the popover at
   the deepest level 3; the text stays the ramp's 900 end, where the walk
   itself clamps. A group has no rule of its own: it raises nothing, so a
   ghost inside one takes the surface the group is in to that surface's
   own hover and press steps.

   Level 1 carries its own rule and the card is why: level 0 walks
   from the Background pin and level 1 from the level above it, which are
   two different fills and, in the dark scheme, two different state fills. */
.card .btn.ghost:hover, .card .btn.ghost.is-hover {
  background: var(--elevation-1-hover);
}
.card .btn.ghost:active, .card .btn.ghost.is-active {
  background: var(--elevation-1-active);
}
.dialog .btn.ghost:hover, .dialog .btn.ghost.is-hover {
  background: var(--elevation-2-hover);
}
.dialog .btn.ghost:active, .dialog .btn.ghost.is-active {
  background: var(--elevation-2-active);
}
.popover .btn.ghost:hover, .popover .btn.ghost.is-hover {
  background: var(--elevation-3-hover);
}
.popover .btn.ghost:active, .popover .btn.ghost.is-active {
  background: var(--elevation-3-active);
}

/* Icon-only form (components/button drawIconButton): a square the
   density's control height on a side, the glyph inset by the density's
   vertical padding — content box ControlHeight − 2·PaddingY, icon.Size's
   rule (20 dp comfortable, 16 dp compact). Emphasis reaches the colours
   and stops there: the square never shrinks. The glyph inherits the
   variant's text colour via currentColor. */
.btn.icon {
  width: var(--density-control-height);
  height: var(--density-control-height);
  min-height: var(--density-control-height);
  padding: var(--density-padding-y);
}
.btn.icon svg {
  width: 100%;
  height: 100%;
  fill: currentColor;
}

/* ---- Badge ----
   The inline annotation components/badge draws: the system's own word about
   a thing, sized to its type and coloured by the role it speaks in. It is
   off the control metrics — no boundary, no minimum height, no vertical
   padding — so its whole height is the label role's line box. A badge beside
   a control is a fraction of that control's height and is meant to be.

   It wears a tinted field, and the field is what says it is not a control.
   One hue at two strengths: a pale fill, the role's hue tinted toward the
   surface until it is a place rather than a mark, with the same hue at
   reading strength on top of it. A saturated fill under knocked-out white
   text is the variant interaction speaks in — .btn's — and a badge
   borrowing it would be claiming to do something. Never invert these; never
   put a status colour on inline style.

   Both halves are tokens because both are derived rather than named on a
   ramp. The fill is realized at a tone at the container chroma, at the depth
   that separates it from the surface the sheet's pages stand on, which no
   var() arithmetic over ramp steps could reproduce; the foreground is then
   the role's pinned base while that base clears the text floor OVER THAT
   FILL, and the nearest step to the mid-value that does otherwise. The pair
   is measured together and only means anything together.

   The side padding is the S2 stop and the vertical padding is none: the
   label role's line box carries its own leading, so the field already stands
   clear of the cap and the descender. The corner is the radius scale's Base
   stop and deliberately not the pill — the pill is .chip's shape, and a chip
   is the thing a badge must not be confused with.

   The label role is the comfortable density's, one step less pronounced
   than a chip's. The sheet re-maps no type role under .compact, so the class
   states the comfortable role and a compact page sets its badges in the
   same one.

   Five variants: the default is the plain category label,
   .success/.warning/.error/.info the four statuses. There is no emphasis
   axis — emphasis belongs where interaction does, and a badge is read rather
   than used.

   The close mark is a Gio interaction and has no class in this sheet, so a
   badge here carries no gap and no affordance — only the utterance. */
.badge {
  box-sizing: border-box;
  display: inline-flex;
  align-items: center;
  white-space: nowrap;
  padding: 0 var(--space-2);
  border-radius: var(--radius-base);
  font-family: var(--font-family);
  font-size: var(--font-label-medium-size);
  line-height: var(--font-label-medium-line-height);
  font-weight: var(--font-label-medium-weight);
  letter-spacing: var(--font-label-medium-tracking);
  background: var(--color-badge-neutral-fill);
  color: var(--color-badge-neutral);
}
.badge.success {
  background: var(--color-badge-success-fill);
  color: var(--color-badge-success);
}
.badge.warning {
  background: var(--color-badge-warning-fill);
  color: var(--color-badge-warning);
}
.badge.error {
  background: var(--color-badge-error-fill);
  color: var(--color-badge-error);
}
.badge.info {
  background: var(--color-badge-info-fill);
  color: var(--color-badge-info);
}

/* ---- Form controls ----
   Native elements wearing components/input's resolution: the raised level
   under body-large text, the neutral step that reads on the surface the
   control stands on for the border, neutral 700 placeholder, focus promoting
   the border to that surface's ring, disabled fading every colour to the
   disabled fraction of its alpha.

   Every fill below names var(--surface-raised, var(--elevation-1)) rather than
   a ramp step, and that is components/input's controlFill: a control that
   paints a box of its own is raised on whatever hosts it, so it fills one
   step nearer the viewer than its host. It used to be --color-surface, the
   neutral ramp's step 200, which lands on the raised level in the dark
   scheme by coincidence and on no level at all in the light one — a light
   field filled a whole band step BELOW the page it lies on. A surface nearer
   the viewer is never darker in either scheme, and on a desktop a text field
   is the lightest thing in the window, not the darkest. The 1 px border and
   the corner radius still carry the field's edge: a raise says a control is
   raised, not where it ends. */

/* Text field (components/input textfield.go). Height = ControlHeight as a
   floor, vertical inset PaddingY, horizontal inset S3 (12 dp — static, it
   does not follow density). The Gio border is drawn inside the field
   (nested fills), so the CSS padding gives back the 1px the border
   occupies and the outer geometry matches exactly. */
.input {
  box-sizing: border-box;
  display: block;
  width: 100%;
  margin: 0;
  appearance: none;
  min-height: var(--density-control-height);
  padding: calc(var(--density-padding-y) - 1px) calc(var(--space-3) - 1px);
  border: 1px solid var(--surface-border, var(--color-control-border));
  border-radius: var(--radius-md);
  background: var(--surface-raised, var(--elevation-1));
  color: var(--color-text);
  font-family: var(--font-family);
  font-size: var(--font-body-large-size);
  line-height: var(--font-body-large-line-height);
  font-weight: var(--font-body-large-weight);
  letter-spacing: var(--font-body-large-tracking);
}
.input::placeholder { color: var(--color-neutral-700); opacity: 1; }

/* Focus promotes the border to the ring and doubles it to the 2 dp the Gio
   side draws — the second pixel as an inset shadow, so the field's outer
   geometry and text position do not move (Gio thickens the border inward the
   same way). The colour is the ring rule above, not the accent pin: a
   promoted border IS the ring, and the accent pin is the seed a caller chose,
   which measures as low as 1.00:1 against the surface it would be drawn on.
   So the field takes the scheme's measured step like every other control —
   and that step is measured against the levels rather than the field's own
   fill because the band has the fill inside it and the level outside, and
   the level is the side every control on it shares. */
.input:focus-visible, .input.is-focus {
  outline: none;
  border-color: var(--color-focus-ring);
  box-shadow: inset 0 0 0 1px var(--color-focus-ring);
}
.input:disabled {
  background: color-mix(in srgb, var(--surface-raised, var(--elevation-1)) var(--state-disabled-opacity), transparent);
  color: color-mix(in srgb, var(--color-text) var(--state-disabled-opacity), transparent);
  border-color: color-mix(in srgb, var(--surface-border, var(--color-control-border)) var(--state-disabled-opacity), transparent);
}
.input:disabled::placeholder {
  color: color-mix(in srgb, var(--color-neutral-700) var(--state-disabled-opacity), transparent);
}

/* Dropdown (components/input dropdown.go): the trigger is a text field
   whose right side reserves S3 + 16 dp chevron + S3, the same inset
   drawTrigger keeps clear of the label. The chevron itself is drawn by the
   .select-wrap wrapper — a native select cannot carry a generated child —
   as a border-built triangle 16 dp wide and 8 dp tall (drawChevron's
   half/quarter geometry) in neutral 700, the low-contrast glyph step. */
.select { padding-right: calc(var(--space-3) * 2 + 16px - 1px); }
.select-wrap { position: relative; display: block; }
.select-wrap::after {
  content: "";
  position: absolute;
  right: var(--space-3);
  top: 50%;
  margin-top: -4px; /* half the 8px glyph height */
  width: 0;
  height: 0;
  border-left: 8px solid transparent;  /* 16 dp chevron width */
  border-right: 8px solid transparent;
  border-top: 8px solid var(--color-neutral-700);
  pointer-events: none;
}
.select-wrap:has(.select:disabled)::after {
  border-top-color: color-mix(in srgb, var(--color-neutral-700) var(--state-disabled-opacity), transparent);
}

/* Checkbox (components/input checkbox.go): a 20 dp glyph (checkboxBoxSize
   — a component constant, not a token; it does not follow density) over its
   own raised fill. Unchecked, its 2 dp edge is the neutral step the ramp answers
   with for a 3:1 graphic on the level the box stands on — the surface
   floor's --color-control-border (600 in the light scheme, 500 in the dark)
   unless a raised host has re-pointed --surface-border. The radio, the text
   field and the dropdown trigger wear that same edge, all four asking the
   ramp the one question rather than naming a step between them.
   Checked, the box is the
   accent fill under a check mark in the on-accent pin, because a fill says
   a colour was applied and only the mark says what that means: a column of
   fills carries completion in hue alone, which is the one channel a reader
   may not have.

   The mark is drawn, not encoded. Gio strokes the icon set's centre line —
   (4.5,12) → (9,16.5) → (19.5,6) on the set's 24-unit grid, a 2-unit
   DIAGONAL band, round caps and joins — and at the 20 px glyph one grid
   unit is 5/6 px, so the band is 1.667 px wide (±0.833 either side of the
   centre) and the arms run from (3.75,10) to (7.5,13.75) to (16.25,5).
   Each arm is one background layer: a linear-gradient banding its own box
   perpendicular to the arm, 45deg for the short "\" arm and 135deg for the
   long "/" one. Every stop is written from 50% because each box is sized
   to its arm — the segment grown by half a band along its own axis, which
   makes the arm the box's diagonal and the box's corners the round caps'
   own tips:
     short arm: 4.929 px square at 3.161,9.411   (3.75,10)→(7.5,13.75)
     long arm:  9.929 px square at 6.911,4.411   (7.5,13.75)→(16.25,5)
   CSS has no line cap, so the caps come out cut square inside those tips
   rather than rounded — the same trade the icon set's own SVG files make
   when they draw their caps as an explicit contour, and a sub-pixel one at
   this size. background-origin is the border box so the grid is the 20 px
   glyph the Gio side scales on, not the 16 px inside the edge. The focus
   ring is the shared rule above. */
.checkbox, .radio {
  box-sizing: border-box;
  appearance: none;
  flex: none;
  width: 20px;  /* checkboxBoxSize / radioCircleSize: 20 dp */
  height: 20px;
  margin: 0;
  border: 2px solid var(--surface-border, var(--color-control-border));
  background: var(--surface-raised, var(--elevation-1));
  cursor: pointer;
}
.checkbox {
  border-radius: var(--radius-sm);
}
.checkbox:checked, .checkbox.is-checked {
  border-color: var(--color-accent);
  background-color: var(--color-accent);
  background-image:
    linear-gradient(45deg, transparent calc(50% - 0.833px), var(--color-on-accent) calc(50% - 0.833px), var(--color-on-accent) calc(50% + 0.833px), transparent calc(50% + 0.833px)),
    linear-gradient(135deg, transparent calc(50% - 0.833px), var(--color-on-accent) calc(50% - 0.833px), var(--color-on-accent) calc(50% + 0.833px), transparent calc(50% + 0.833px));
  background-origin: border-box;
  background-repeat: no-repeat;
  background-position: 3.161px 9.411px, 6.911px 4.411px;
  background-size: 4.929px 4.929px, 9.929px 9.929px;
}
/* Disabled fades the box and its mark together, so the check stays the
   on-colour of the fill it is drawn on. background-color rather than the
   background shorthand: the shorthand would reset the checked rule's
   layer geometry and the disabled-checked box would lose its mark. */
.checkbox:disabled {
  cursor: default;
  border-color: color-mix(in srgb, var(--surface-border, var(--color-control-border)) var(--state-disabled-opacity), transparent);
  background-color: color-mix(in srgb, var(--surface-raised, var(--elevation-1)) var(--state-disabled-opacity), transparent);
}
.checkbox:checked:disabled, .checkbox.is-checked:disabled {
  border-color: color-mix(in srgb, var(--color-accent) var(--state-disabled-opacity), transparent);
  background-color: color-mix(in srgb, var(--color-accent) var(--state-disabled-opacity), transparent);
  background-image:
    linear-gradient(45deg, transparent calc(50% - 0.833px), color-mix(in srgb, var(--color-on-accent) var(--state-disabled-opacity), transparent) calc(50% - 0.833px), color-mix(in srgb, var(--color-on-accent) var(--state-disabled-opacity), transparent) calc(50% + 0.833px), transparent calc(50% + 0.833px)),
    linear-gradient(135deg, transparent calc(50% - 0.833px), color-mix(in srgb, var(--color-on-accent) var(--state-disabled-opacity), transparent) calc(50% - 0.833px), color-mix(in srgb, var(--color-on-accent) var(--state-disabled-opacity), transparent) calc(50% + 0.833px), transparent calc(50% + 0.833px));
}

/* Radio (components/input radio.go): the same 20 dp glyph as a circle;
   selected keeps the gap ring and fills a 10 dp accent dot (radioDotSize) —
   outer accent ring, the glyph's own raised fill, dot, exactly the Gio nested
   fills. The gap is the glyph's interior, so it takes --surface-raised in the
   chosen state as much as the resting one: the dot is drawn on the radio's
   own surface, not on the host's. */
.radio { border-radius: var(--radius-full); }
.radio:checked, .radio.is-checked {
  border-color: var(--color-accent);
  background: radial-gradient(circle, var(--color-accent) 5px, var(--surface-raised, var(--elevation-1)) 5px); /* 10 dp dot */
}
.radio:disabled {
  cursor: default;
  border-color: color-mix(in srgb, var(--surface-border, var(--color-control-border)) var(--state-disabled-opacity), transparent);
  background: color-mix(in srgb, var(--surface-raised, var(--elevation-1)) var(--state-disabled-opacity), transparent);
}
.radio:checked:disabled, .radio.is-checked:disabled {
  border-color: color-mix(in srgb, var(--color-accent) var(--state-disabled-opacity), transparent);
  background: radial-gradient(circle, color-mix(in srgb, var(--color-accent) var(--state-disabled-opacity), transparent) 5px, color-mix(in srgb, var(--surface-raised, var(--elevation-1)) var(--state-disabled-opacity), transparent) 5px);
}

/* ---- Card and group ----
   The two are one ruling with two answers. patterns/card singles something
   out: one rounded surface raised on the content by tonal step alone — no
   cast shadow, because a card is raised, not floating, and the dp shadows
   stay reserved for surfaces that can leave (menus, dialogs, toasts); no
   line of its own, because the raise is what does the singling out; and no
   role, because what a developer wants to say about a card is a .badge in
   its header. patterns/group divides the page: a hairline around related
   components at the level of the surface the group is in, taking that
   surface's own fill, so it declares no background at all and nothing is
   derived against it.

   A card fills at the raise walked from the content (--elevation-1) and
   carries the seam that raise owes — var(--surface-seam,
   var(--elevation-0-seam)), transparent in the scheme where the fill tells
   the raise on its own and a hairline where it does not. A card on the
   content is told by its fill in both schemes, so that border shows only on
   a card placed inside a host already at the top of its scheme, which is
   what --surface-seam is redeclared for on .dialog and .popover.

   A group's hairline is var(--surface-hairline,
   var(--elevation-0-hairline)): the seam of two regions sharing one fill,
   less pronounced than the 3:1 mark a graphic carrying meaning owes, and
   never transparent — it is the whole of what says where the group ends. A
   host at another level redeclares --surface-hairline, as .dialog and
   .popover do, so a group inside one is derived against that host's fill.

   Both carry radius Lg, an S4 inset and S3 gaps between the slots —
   exactly drawCard's rad.Lg / sp.S4 / sp.S3. The Gio line lies inside the
   bounds where the CSS border does, so the padding gives back the border's
   1px and the slots land where the Gio inset puts them. Neither styles
   slot text of its own: the Gio card draws no text at all, and the group
   draws only its own label, so slot typography belongs to the content. */
.card {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: calc(var(--space-4) - 1px);
  border: 1px solid var(--surface-seam, var(--elevation-0-seam));
  border-radius: var(--radius-lg);
  background: var(--elevation-1);
  color: var(--color-text);
  --surface-raised: var(--elevation-2);
}
.group {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: calc(var(--space-4) - 1px);
  border: 1px solid var(--surface-hairline, var(--elevation-0-hairline));
  border-radius: var(--radius-lg);
  color: var(--color-text);
}
/* The group's own label: top-leading, inside the hairline, as the first row
   of the group's stack — the platform's idiom for a section header over a
   bordered container, and not the fieldset legend cut into the top line,
   which has no native counterpart and does not survive a Lg corner. The
   label-large role in the ramp's low-contrast step: a group wears no role,
   so it is never the accent, and never the Text pin, which would give a
   section header the weight of the content it names. */
.group-label {
  font-size: var(--font-label-large-size);
  line-height: var(--font-label-large-line-height);
  font-weight: var(--font-label-large-weight);
  letter-spacing: var(--font-label-large-tracking);
  color: var(--color-neutral-700);
}

/* ---- Table ----
   patterns/table: the whole grid stands on the raised level (drawTable
   fills the surface prop, which defaults to level 1) and the header band on
   the raise walked from it (drawHeaderRow raises the grid's own fill, never
   an absolute step), under neutral 700 label-large text (drawHeaderCell).
   Both name --elevation-N rather than a ramp step, because a
   level is not a ramp step in both schemes and a table that named one
   would read as a mirror of itself between the two. Header and body rows are each exactly
   one control height tall — the row-height rule (list.RowHeight), so .compact
   re-pitches the whole grid — and every row closes with a 1 dp Divider rule
   drawn inside its height. Cells inset horizontally by S3 (cellPadDp,
   12 dp — static, the same inset rule the input uses); body text is body-medium at the
   Text pin (RenderTextCell). There is no zebra: rows separate by the rules
   alone. Sort is a 10x5 dp chevron in neutral 700 at the header's right
   inset (drawSortChevron), drawn on the active column only — .sort-asc /
   .sort-desc on a .sortable header. */
.table {
  box-sizing: border-box;
  border-collapse: separate;
  border-spacing: 0;
  table-layout: fixed;
  width: 100%;
  background: var(--elevation-1);
  color: var(--color-text);
  font-family: var(--font-family);
  font-size: var(--font-body-medium-size);
  line-height: var(--font-body-medium-line-height);
  font-weight: var(--font-body-medium-weight);
  letter-spacing: var(--font-body-medium-tracking);
}
.table th, .table td {
  box-sizing: border-box;
  height: var(--density-control-height);
  padding: 0 var(--space-3);
  border-bottom: 1px solid var(--color-divider);
  text-align: left;
  vertical-align: middle;
  white-space: nowrap;
  overflow: hidden;
}
.table th {
  position: relative;
  background: var(--elevation-2);
  color: var(--color-neutral-700);
  font-size: var(--font-label-large-size);
  line-height: var(--font-label-large-line-height);
  font-weight: var(--font-label-large-weight);
  letter-spacing: var(--font-label-large-tracking);
}
.table th.sortable { cursor: pointer; }
.table th.sort-asc::after, .table th.sort-desc::after {
  content: "";
  position: absolute;
  right: var(--space-3);
  top: 50%;
  transform: translateY(-50%);
  width: 0;
  height: 0;
  border-left: 5px solid transparent;  /* 10 dp chevron width */
  border-right: 5px solid transparent;
}
.table th.sort-asc::after { border-bottom: 5px solid var(--color-neutral-700); }  /* 5 dp tall, apex up */
.table th.sort-desc::after { border-top: 5px solid var(--color-neutral-700); }

/* ---- Navigation ----
   The four navigation patterns. All four rest their interactive cells on
   the Surface fill (neutral 200), so the pointer states are that fill's
   own one-step ramp walk — hover one step to neutral 300, exactly the state
   fill the ghost variant paints on the same level. The Gio side draws no hover
   (a native window has the pointer; a static page shows the resolution),
   but the steps are the tokens' StateColor walk from step 200, not a new
   mix. Selection is what the Gio side does draw: the 2 dp Primary underline
   (navbar.go / tabs.go underlineDp) or the sidebar's StateColor(RolePrimary,
   200, StateSelected) two-step walk to primary 400. */

/* Navbar (patterns/navbar navbar.go): a horizontal Surface bar —
   drawNavbar fills Surface, insets PaddingY vertically and S4 horizontally,
   and patterns/shell pins the bar to ControlHeight + 2*PaddingY (52 dp
   comfortable, 40 compact). Slots run brand, centred links, actions; the
   links row centres in the space brand and actions leave over (that space
   halved), which margin-inline auto reproduces exactly, including the
   documented off-centre approximation when the end slots differ. Each slot
   is drawn on the bar's own centre line, whatever height it comes back —
   align-items centre over a definite content box is that same line. */
.navbar {
  box-sizing: border-box;
  display: flex;
  align-items: center;
  min-height: calc(var(--density-control-height) + 2 * var(--density-padding-y));
  padding: var(--density-padding-y) var(--space-4);
  background: var(--elevation-chrome);
  color: var(--color-text);
}
.navbar-links {
  display: flex;
  align-items: center;
  gap: var(--space-2);  /* linksRow's HSpacer(sp.S2) between link cells */
  margin-inline: auto;  /* the leftover space, halved on either side */
}

/* A link cell (navbar.go linkWidget): label-large at the Text pin inside
   (S3, PaddingY) padding, with a 2 dp underline slot along the bottom edge
   that the Active link fills with the Primary pin — the underline runs the
   full cell width, padding included, exactly the image.Rect the Gio side
   fills. Hover is the Surface fill's one-step walk. */
.navbar-link, .tab {
  box-sizing: border-box;
  display: inline-flex;
  align-items: center;
  appearance: none;
  margin: 0;
  border: none;
  background: transparent;
  text-decoration: none;
  user-select: none;
  white-space: nowrap;
  cursor: pointer;
  font-family: var(--font-family);
  font-size: var(--font-label-large-size);
  line-height: var(--font-label-large-line-height);
  font-weight: var(--font-label-large-weight);
  letter-spacing: var(--font-label-large-tracking);
  color: var(--color-text);
  border-bottom: 2px solid transparent;  /* underlineDp: the underline slot */
}
.navbar-link {
  padding: var(--density-padding-y) var(--space-3);
}
.navbar-link:hover, .navbar-link.is-hover,
.tab:hover, .tab.is-hover {
  background: var(--color-neutral-300);
}
.navbar-link.selected, .tab.selected {
  border-bottom-color: var(--color-accent);
}

/* Tabs (patterns/tabs tabs.go): the strip is a raised row of tab cells
   exactly ControlHeight tall (drawTabs pins stripH to the density), each
   cell its label plus 2*S3 horizontal padding with the label centred in the
   height that remains above the 2 dp underline slot — which border-box
   centring reproduces. The selected cell fills the slot with the Primary
   pin; content panels below the strip are the caller's. */
.tabs {
  box-sizing: border-box;
  display: flex;
  align-items: stretch;
  height: var(--density-control-height);
  background: var(--elevation-1);
}
.tab {
  height: 100%;
  padding: 0 var(--space-3);
  justify-content: center;
}

/* Sidebar (patterns/sidebar sidebar.go): a vertical Surface rail at the
   pattern's two contractual widths — 192 dp expanded, 48 dp collapsed
   (expandedDp/collapsedDp: component constants, deliberately not tokens
   and not density-responsive; a different rail copies the pattern). The
   toggle row and every item row are exactly ControlHeight tall (the
   row-height rule, list.RowHeight), so .compact re-pitches the rail. The whole
   rail is one keyboard stop — the focus ring belongs to the rail, and the
   Arrow keys move .selected — so the ring rule below targets .sidebar
   itself, not the rows. */
.sidebar {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  width: 192px;  /* expandedDp */
  background: var(--elevation-chrome);
  color: var(--color-text);
  overflow: hidden;
}
.sidebar.collapsed { width: 48px; }  /* collapsedDp */

/* The collapse affordance (drawToggle): a full-width ControlHeight row with
   a 16 dp neutral-700 placeholder glyph centred in it — pointer-only, never
   a Tab stop. */
.sidebar-toggle {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  height: var(--density-control-height);
  cursor: pointer;
}
.sidebar-toggle::after {
  content: "";
  width: 16px;   /* drawToggle's 16 dp glyph */
  height: 16px;
  background: var(--color-neutral-700);
}

/* An item row (drawItem): a 48 dp leading icon column (iconColDp) with the
   glyph centred in it, the label-large label starting at exactly the column
   edge, vertically centred, one line, clipped rather than wrapped — which
   is also what hides the labels at the collapsed width. Selected is
   StateColor(RolePrimary, 200, StateSelected): the two-step walk past the
   Surface fill to primary 400. */
.sidebar-item {
  box-sizing: border-box;
  flex: none;
  display: flex;
  align-items: center;
  height: var(--density-control-height);
  overflow: hidden;
  white-space: nowrap;
  cursor: pointer;
  font-family: var(--font-family);
  font-size: var(--font-label-large-size);
  line-height: var(--font-label-large-line-height);
  font-weight: var(--font-label-large-weight);
  letter-spacing: var(--font-label-large-tracking);
  color: var(--color-text);
  text-decoration: none;
  user-select: none;
}
.sidebar-item:hover, .sidebar-item.is-hover { background: var(--color-neutral-300); }
.sidebar-item.selected { background: var(--color-primary-400); }
.sidebar-item-icon {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;  /* iconColDp */
  height: 100%;
}

/* Breadcrumb (patterns/breadcrumb breadcrumb.go): a row of title-small
   segments with S2 gaps around 12 dp chevron separators (chevronDp). The
   last segment is the current location at the Text pin; ancestors rest on
   neutral 700 and, being links, hover to neutral 900 — the ghost variant's
   text walk on the same surface. Colour follows position (labelColor), so
   :last-child carries it; .current forces it for a specimen. */
.crumbs {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-family: var(--font-family);
  font-size: var(--font-title-small-size);
  line-height: var(--font-title-small-line-height);
  font-weight: var(--font-title-small-weight);
  letter-spacing: var(--font-title-small-tracking);
}
.crumb {
  color: var(--color-neutral-700);
  text-decoration: none;
  white-space: nowrap;
}
.crumb:hover, .crumb.is-hover { color: var(--color-neutral-900); }
.crumbs .crumb:last-child, .crumb.current { color: var(--color-text); }

/* The chevron separator (chevronWidget/drawChevron): a right-pointing
   triangle half the 12 dp box wide and the full box tall, centred in it,
   in neutral 700 — the low-contrast glyph step. */
.crumb-sep {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 12px;   /* chevronDp */
  height: 12px;
}
.crumb-sep::before {
  content: "";
  width: 0;
  height: 0;
  border-top: 6px solid transparent;    /* 12 dp tall */
  border-bottom: 6px solid transparent;
  border-left: 6px solid var(--color-neutral-700);  /* 6 dp deep, apex along +X */
}

/* The keyboard ring, identical to every other control's: per-cell for the
   navbar, tabs and breadcrumb (each cell is its own Clickable focus tag);
   on the rail itself for the sidebar, whose single stop is the item list. */
.navbar-link:focus-visible, .navbar-link.is-focus,
.tab:focus-visible, .tab.is-focus,
.crumb:focus-visible, .crumb.is-focus,
.sidebar:focus-visible, .sidebar.is-focus {
  outline: var(--focus-ring-width) solid var(--color-focus-ring);
  outline-offset: calc(var(--focus-ring-width) / -2);
}

/* ---- Overlays ----
   The transient surfaces: the scrimmed dialog (patterns/modal), the
   unscrimmed popover (patterns/popover), the inverse-video tooltip
   (patterns/tooltip) and the floating toast (patterns/toast). The elevation
   grammar: a scrimmed modal sits at level 2 (the scrim, not the
   fill, isolates it); an unscrimmed, shadowless popover separates by fill
   alone and takes the deepest level 3; a toast takes no level at all — it
   inverts, and keeps the level-3 cast shadow to say it can leave; the tooltip
   takes no step at all — it inverts instead, because a bubble that small
   needs the stronger cue. */

/* Scrim (modal.go drawModal/scrimColor): the whole-plane dimmer under a
   dialog — --color-scrim, black at the fixed 50% alpha in both modes. The
   scrim centres the dialog, exactly as drawModal centres the surface in the
   window plane. Behaviour is part of the pattern: on a PANEL a scrim press
   invokes OnClose; on a DECISION the scrim is INERT — it absorbs presses and
   answers none of them, because dismissal is one of the decision's answers
   and a stray click must not give it. */
.scrim {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  background: var(--color-scrim);
}

/* Dialog (modal.go drawModal): the centred surface — width 75% of the
   window plane clamped to 180–560 dp, height hugging its content between the
   120 dp floor and the 560 dp cap (overflow clips), a level-2 fill under
   the 1 dp neutral 500 stroke, radius Lg, an S5 inset and S3 gaps between
   header, body and footer. The padding gives back the border's 1px, the
   card's trick, so content lands where the Gio inset puts it. G0A.2's two
   purposes share this one surface: a PANEL carries a ghost icon close
   (.btn.ghost.icon) in its header and no footer of its own; a DECISION
   carries no X anywhere and a .dialog-footer whose right-aligned actions
   end in the Return-bound default. */
.dialog {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  width: 75%;
  min-width: 180px;
  max-width: 560px;
  min-height: 120px;
  max-height: min(75%, 560px);
  overflow: hidden;
  padding: calc(var(--space-5) - 1px);
  border: 1px solid var(--color-dialog-border);
  border-radius: var(--radius-lg);
  background: var(--elevation-2);
  color: var(--color-text);
  --surface-border: var(--color-dialog-border);
  --surface-raised: var(--elevation-3);
  --surface-seam: var(--elevation-2-seam);
  --surface-hairline: var(--elevation-2-hairline);
}

/* The header row (modal.go headerWidget): the title-medium title on the
   left, the close affordance — when the purpose shows one — on the right,
   middle-aligned so the ghost button's square drives the row height. */
.dialog-header {
  display: flex;
  align-items: center;
}
.dialog-title {
  flex: 1;
  font-family: var(--font-family);
  font-size: var(--font-title-medium-size);
  line-height: var(--font-title-medium-line-height);
  font-weight: var(--font-title-medium-weight);
  letter-spacing: var(--font-title-medium-tracking);
}

/* The footer row (modal.go footerWidget): right-aligned actions with S2
   gaps. Each action is a bare component owning its own focus ring — the
   dialog wraps and decorates nothing. */
.dialog-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--space-2);
}

/* Popover (popover.go drawPopover): the unscrimmed anchored surface —
   content plus an S3 inset (the padding gives back the border's 1px),
   clamped to the 48x24 dp minimum, the deepest level-3 fill under the 1 dp
   neutral 500 stroke, radius Md. Positioning against the anchor is the
   page's: the pattern centres the surface on the anchor's midline, one S2
   gap away on the Placement side. */
.popover {
  box-sizing: border-box;
  display: inline-block;
  min-width: 48px;
  min-height: 24px;
  padding: calc(var(--space-3) - 1px);
  border: 1px solid var(--color-popover-border);
  border-radius: var(--radius-md);
  background: var(--elevation-3);
  color: var(--color-text);
  --surface-border: var(--color-popover-border);
  --surface-raised: var(--elevation-3);
  --surface-seam: var(--elevation-3-seam);
  --surface-hairline: var(--elevation-3-hairline);
}

/* The tail (popover.go drawTail): a triangle 12 dp across the base and
   6 dp deep in the surface's own fill, bridging the gap with its tip at
   the anchor. The modifier names the popover's Placement — a .top popover
   sits above its anchor, so its tail points down. */
.popover-tail { width: 0; height: 0; }
.popover-tail.top {
  border-left: 6px solid transparent;   /* 12 dp base */
  border-right: 6px solid transparent;
  border-top: 6px solid var(--elevation-3);  /* 6 dp deep, tip down */
}
.popover-tail.bottom {
  border-left: 6px solid transparent;
  border-right: 6px solid transparent;
  border-bottom: 6px solid var(--elevation-3);
}
.popover-tail.left {
  border-top: 6px solid transparent;
  border-bottom: 6px solid transparent;
  border-left: 6px solid var(--elevation-3);
}
.popover-tail.right {
  border-top: 6px solid transparent;
  border-bottom: 6px solid transparent;
  border-right: 6px solid var(--elevation-3);
}

/* Tooltip (tooltip.go drawSurface): the inverse-video bubble — the Text
   pin as the surface under a label in Surface, label-small, radius Sm, S2/S1
   padding, clamped to the 24x16 dp minimum. No elevation level and no
   shadow: inversion is the whole cue. */
.tooltip {
  box-sizing: border-box;
  display: inline-flex;
  align-items: center;
  min-width: 24px;
  min-height: 16px;
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-sm);
  background: var(--color-text);
  color: var(--color-surface);
  white-space: nowrap;
  font-family: var(--font-family);
  font-size: var(--font-label-small-size);
  line-height: var(--font-label-small-line-height);
  font-weight: var(--font-label-small-weight);
  letter-spacing: var(--font-label-small-tracking);
}

/* Toast: one queued notification — 240 dp wide, a 36 dp legibility floor
   that deliberately does not follow density (a toast is not a control),
   radius Md, label-medium on the inverse pair, and floating on the level-3
   cast shadow. It is the one surface built out of the counterpart scheme:
   dark on a light scheme, light on a dark one, so a message that can
   appear over any pane separates from all of them without claiming a
   level none of them can be under. The shadow stays for what it says
   rather than for the separation — this layer is temporary — and there is
   no outline, which on the old tinted level-2 base was the only thing
   giving the chip an edge.
   The level shows as a leading edge one S2 wide, painted as a two-stop
   gradient so the chip's own radius rounds it: the level's own mark on the
   inverse surface, the step of that level's ramp nearest its mid-value step
   that still reads over the chip. It was one S1, which is the width this
   desktop keeps for separators, pane strokes and insets — furniture it does
   not want looked at — and a mark identified by its colour cannot be drawn
   at furniture width. Two stops is as wide as the air above the message and
   two thirds of the air beside it, which is where the widening stops: an
   edge as wide as the gap it holds the text off by reads as a panel the
   message sits next to rather than as the chip's own edge.
   The mark arrives as a token rather than as a ramp reference because the
   two schemes do not land on one step — a light scheme's marks come off
   step 500 and a dark scheme's off step 400 — and because one step for both
   cost the light scheme its reds, the error edge coming out the pale salmon
   a red turns into when it is asked to sit as light as an orange wants to.
   Each level takes its own status ramp — info included, which reads off the
   info ramp rather than off the accent, so a themed brand cannot make an
   informational chip wear the colour of an alarming one. */
.toast {
  box-sizing: border-box;
  display: flex;
  align-items: center;
  width: 240px;
  min-height: 36px;
  padding: var(--space-2) var(--space-3);
  padding-left: calc(var(--space-2) + var(--space-3));
  border-radius: var(--radius-md);
  background: linear-gradient(to right, var(--color-info-on-inverse) 0 var(--space-2), var(--color-inverse-surface) var(--space-2));
  color: var(--color-on-inverse-surface);
  box-shadow: var(--shadow-3);
  font-family: var(--font-family);
  font-size: var(--font-label-medium-size);
  line-height: var(--font-label-medium-line-height);
  font-weight: var(--font-label-medium-weight);
  letter-spacing: var(--font-label-medium-tracking);
}
.toast.success {
  background: linear-gradient(to right, var(--color-success-on-inverse) 0 var(--space-2), var(--color-inverse-surface) var(--space-2));
}
.toast.warning {
  background: linear-gradient(to right, var(--color-warning-on-inverse) 0 var(--space-2), var(--color-inverse-surface) var(--space-2));
}
.toast.error {
  background: linear-gradient(to right, var(--color-error-on-inverse) 0 var(--space-2), var(--color-inverse-surface) var(--space-2));
}

/* The stack (toast.go paintStack): a corner-anchored column with S2 gaps,
   inset S4 from the window plane's edges (the page anchors it); newest toast
   nearest the anchored edge. */
.toast-stack {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  width: 240px;
}
`
