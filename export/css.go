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

// hexRGBA formats a colour as lowercase #rrggbb, or #rrggbbaa where it
// carries a coverage: the platform's labels, seams, overlays and focus ring
// are a colour at a coverage over whatever is beneath them, and a sheet that
// dropped the coverage would state a colour the platform never paints.
func hexRGBA(c stdcolor.NRGBA) string {
	if c.A == 0xff {
		return hexRGB(c)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
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

// platformNames orders the platform's colour set under its CSS names, one
// entry per field of tokens.PlatformColors, in the struct's own order.
//
// The naming rule is mechanical and has no table of exceptions: the AppKit
// name the field carries, kebab-cased, under a --platform- prefix. So
// windowBackgroundColor is --platform-window-background and systemRed is
// --platform-system-red, exactly as the Go field drops AppKit's trailing
// "Color" and nothing else. The ten fields AppKit has no name for are spelled
// the same way from the name the token set gives them.
var platformNames = []struct {
	name string
	pick func(tokens.PlatformColors) stdcolor.NRGBA
}{
	{"window-background", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.WindowBackground }},
	{"under-page-background", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.UnderPageBackground }},
	{"control-background", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ControlBackground }},
	{"text-background", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.TextBackground }},
	{"selected-content-background", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SelectedContentBackground }},
	{"unemphasized-selected-content-background", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.UnemphasizedSelectedContentBackground }},
	{"selected-text-background", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SelectedTextBackground }},
	{"unemphasized-selected-text-background", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.UnemphasizedSelectedTextBackground }},
	{"find-highlight", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.FindHighlight }},
	{"separator", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Separator }},
	{"grid", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Grid }},
	{"label", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Label }},
	{"secondary-label", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SecondaryLabel }},
	{"tertiary-label", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.TertiaryLabel }},
	{"quaternary-label", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.QuaternaryLabel }},
	{"text", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Text }},
	{"placeholder-text", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.PlaceholderText }},
	{"selected-text", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SelectedText }},
	{"link", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Link }},
	{"header-text", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.HeaderText }},
	{"control", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Control }},
	{"control-text", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ControlText }},
	{"disabled-control-text", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.DisabledControlText }},
	{"selected-control", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SelectedControl }},
	{"selected-control-text", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SelectedControlText }},
	{"alternate-selected-control-text", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.AlternateSelectedControlText }},
	{"control-accent", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ControlAccent }},
	{"keyboard-focus-indicator", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.KeyboardFocusIndicator }},
	{"system-red", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemRed }},
	{"system-orange", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemOrange }},
	{"system-yellow", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemYellow }},
	{"system-green", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemGreen }},
	{"system-mint", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemMint }},
	{"system-teal", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemTeal }},
	{"system-cyan", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemCyan }},
	{"system-blue", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemBlue }},
	{"system-indigo", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemIndigo }},
	{"system-purple", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemPurple }},
	{"system-pink", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemPink }},
	{"system-brown", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemBrown }},
	{"system-gray", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SystemGray }},
	{"shadow", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Shadow }},
	{"highlight", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Highlight }},
	{"sidebar-material", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SidebarMaterial }},
	{"card-fill", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.CardFill }},
	{"push-button-fill", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.PushButtonFill }},
	{"hover-overlay", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.HoverOverlay }},
	{"press-overlay", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.PressOverlay }},
	{"floating-shadow", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.FloatingShadow }},
	{"field-edge", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.FieldEdge }},
	{"scrollbar-thumb", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ScrollbarThumb }},
	{"alternating-content-background", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.AlternatingContentBackground }},
	{"scrim", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Scrim }},
}

// platformVars is the platform's colour set as custom properties, one per
// field.
//
// A coverage is written out. Everywhere else in this sheet a colour is opaque
// and #rrggbb says all of it, but the platform's answer for a label, a seam,
// an overlay or the focus ring IS a coverage over whatever lies beneath, so
// those are emitted as #rrggbbaa and a rule that lays one over a fill gets
// what the platform gets. A browser composites #rrggbbaa in encoded sRGB,
// which is the space the platform composites in and the space theme/color's
// Flatten takes, so the two sides land on the same pixel.
func platformVars(p tokens.PlatformColors) []cssVar {
	vars := make([]cssVar, 0, len(platformNames))
	for _, n := range platformNames {
		vars = append(vars, cssVar{"--platform-" + n.name, hexRGBA(n.pick(p))})
	}
	return vars
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
	{"seam", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Seam }},
	// The inverse pair, emitted as first-class tokens for the same reason
	// the state walk below is: it resolves off the counterpart scheme's
	// neutral ramp, and a sheet holding only this scheme's ramps has no
	// var() arithmetic that could reach it.
	{"inverse-surface", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.InverseSurface }},
	{"on-inverse-surface", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.OnInverseSurface }},
	// The highlight: the fill marking content the reader was brought to. It
	// is emitted as a first-class token because it belongs to no ramp — it
	// is one reserved yellow laid over the surface at less than full
	// strength, which no var() reference over the ramps could reach — and
	// it is the fill laid over the surface these pages stand on, level 0,
	// which is the one answer a sheet has to give. It is not a status and
	// no status hue serves it; see the tokens package's highlight.go for
	// the colour and the coverage.
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
	// clears the text floor over it and the nearest step to the mid-value that
	// does otherwise (ForegroundOnAtFloor). Neutral has no pinned base — the
	// neutral ramp carries no pin — so it takes the walk directly (MarkOn), at
	// the same floor, and its fill comes back as depth alone.
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
	// the neutral step nearest step 500 that reaches graphicFloor against the
	// level-0 surface a control on the page is guaranteed against. Naming
	// step 500 in both schemes — which this sheet did, at every one of those
	// four sites — measures |Lc| 46.80 in the light scheme and 24.95 in the
	// dark, under the floor in the scheme most people read in at night. The
	// walk answers 500 in the light scheme and 600 in the dark and needs to
	// know nothing about either.
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
		// Both ceilings sit above anything a pairing can reach — |Lc| tops
		// out near 106 and the luminance ratio at 21 — so the first level
		// always lowers them.
		const maxLc, maxRatio = 110.0, 21.0
		worst, worstBorder := maxLc, maxRatio
		for _, lvl := range standableLevels {
			surface := t.SurfaceAt(lvl.level)
			if got := vgcolor.Magnitude(step, surface); got < worst {
				worst = got
			}
			border := t.MarkOn(tokens.RoleNeutral, surface, graphicFloor)
			if got := luminanceRatio(step, border); got < worstBorder {
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
	if fill.A == 0 || vgcolor.Magnitude(ring, fill) >= graphicFloor {
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

// The floors this sheet derives against are the theme's own, not a second
// spelling of them: onFloor is the text floor a mark on the inverse surface
// is chosen against — a toast's leading edge is the only thing that says
// which level the toast is, so it is held to the text floor rather than to
// what a non-text graphic owes — and graphicFloor is what a control's edge
// and its focus ring owe the surface they are drawn on, neither being
// decoration.
const (
	onFloor      = tokens.TextFloor
	graphicFloor = tokens.GraphicFloor
)

// luminanceRatio is the arithmetic [focusRingBorderSeparation] is measured
// in: (L1+0.05)/(L2+0.05) over the two relative luminances, lighter first.
// It is not a contrast measure — contrast is APCA
// ([vgcolor.Magnitude]) throughout — but the scale on which one grey stops
// being another grey.
func luminanceRatio(a, b stdcolor.NRGBA) float64 {
	la, lb := vgcolor.RelativeLuminance(a), vgcolor.RelativeLuminance(b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

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
	// The field's own height and the stacked row's pitch. The platform draws
	// a text field taller than a push button and a list row shorter than
	// both, all three measured, so each is stated rather than derived from
	// the control height.
	{"field-height", func(d tokens.Density) float32 { return d.FieldHeight }},
	{"row-height", func(d tokens.Density) float32 { return d.RowHeight }},
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
	// The platform's set stands beside the derived one in both blocks, under
	// its own prefix, so a rule may name either while the consumers convert.
	root := append(colorVars(s.Light), platformVars(s.PlatformLight)...)
	block(&b, ":root", append(root, scaleVars(s)...))
	b.WriteString("\n")
	block(&b, ".dark", append(colorVars(s.Dark), platformVars(s.PlatformDark)...))
	b.WriteString("\n")
	block(&b, ".compact", densityVars(tokens.Compact))
	b.WriteString("\n")
	b.WriteString(componentClasses)
	return b.String()
}

// componentClasses is the class layer: the component vocabulary the design surface composes screens from, defined
// entirely over the tokens above.
//
// It mirrors the Gio components, which are the sources of truth, and since
// Phase CE every colour in it is the platform's own name for what that
// element is on the platform — the same mapping components/button,
// components/input, components/badge and the patterns took. No role, no
// ramp step, no level and no derivation appears here.
//
// A coverage composites the same on both sides, which is what lets the
// sheet name a platform colour and stop. AppKit's labels, seams, overlays
// and focus ring are a colour at a coverage over whatever lies beneath; a
// browser composites #rrggbbaa in encoded sRGB, and so does the platform,
// and so does theme/color.Flatten, which is what the Gio side calls before
// it hands the rasterizer an opaque fill. So a rule may name the coverage
// directly and let the browser do the flattening — with one geometric
// proviso: the Gio side flattens each name onto the surface the component
// actually stands on (its Surface property), and CSS composites a
// translucent border or outline over whatever the painting order puts
// under it. Where the two agree — a hairline over the control's own fill —
// the rule lets the element's background paint under the border, as CSS
// does by default. Where the Gio side flattens onto the surface INSTEAD of
// the control's fill — the text field's focus ring, the checkbox's edge —
// the rule sets background-clip: padding-box, so the fill stops at the
// padding box and the edge composites over the page. That pair of clips is
// the whole of what replaced the old --surface-* inheritance chain: a
// control no longer asks which level hosts it, because the browser already
// knows what is under the pixel.
//
// The states are the platform's own answers. Nothing tints on hover but a
// toolbar button, which this sheet has no class for; a press lays
// --platform-press-overlay over whatever fill the variant carries, and
// over the page where it carries none, which is how a ghost gets a fill at
// all. Focus is --platform-keyboard-focus-indicator at --focus-ring-width,
// the same ring in every variant. Disabled is not a fade of the resting
// colours: a fill falls back to the push button's own, and every
// foreground becomes --platform-disabled-control-text.
//
// Every pointer/keyboard state rule also carries a forcing twin class
// (.is-hover, .is-active, .is-focus, .is-checked) grouped into the same
// rule. A static component page cannot press itself, and duplicating the
// declarations in the page would fork the resolution; a grouped selector
// emitted by this generator shares the exact declarations with the live
// pseudo-class, so a forced specimen provably renders as the live state.
// Disabled needs no twin: the pages force it with the native attribute.
const componentClasses = `/* ---- Component classes ----
   The class vocabulary, built only on the tokens above, and every colour in
   it the platform's own name. .btn mirrors components/button: filled by
   default, .tonal and .ghost the less pronounced variants, states resolved
   exactly as buttonColors resolves them. .input/.select/.checkbox/.radio
   mirror components/input and the trigger components/picker draws for it,
   .badge the inline annotation components/badge draws (the plain category
   label and the four statuses; the close mark is a Gio interaction and has
   no class here), .card the patterns/card box and .group the patterns/group
   hairline, .table the patterns/table grid, the navigation family —
   .navbar, .tabs, .sidebar (patterns) and .crumbs (components/breadcrumb) —
   and the overlay family — .scrim/.dialog (patterns/modal), .popover,
   .tooltip, .toast — the transient surfaces. Each state rule carries a
   forcing twin class (.is-hover, .is-active, .is-focus, .is-checked) so a
   static page can show the state with the very declarations the live
   pseudo-class applies. */

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
  background: var(--platform-control-accent);
  color: var(--platform-alternate-selected-control-text);
}

/* Tonal is the platform's ordinary push button: its own measured fill under
   the control text, inside a separator hairline. Both the label and the
   hairline are coverages the platform states over that fill, so the element's
   background paints under the border and the browser flattens them there. */
.btn.tonal {
  background: var(--platform-push-button-fill);
  border: 1px solid var(--platform-separator);
  padding: calc(var(--density-padding-y) - 1px) calc(var(--density-padding-x) - 1px);
  color: var(--platform-control-text);
}

/* Ghost is the borderless kind: no fill at rest, the control text over
   whatever it stands on. Held, the press overlay composites straight onto
   that surface, which is the only fill a ghost ever has. */
.btn.ghost {
  background: transparent;
  color: var(--platform-control-text);
}

/* Icon-only form (components/button drawIconButton): a square the density's
   control height on a side, the glyph inset by the density's vertical
   padding - content box ControlHeight - 2*PaddingY, icon.Size's rule. The
   variant reaches the colours and stops there: the square never shrinks. The
   glyph inherits the variant's text colour via currentColor. */
.btn.icon {
  width: var(--density-control-height);
  height: var(--density-control-height);
  min-height: var(--density-control-height);
  padding: var(--density-padding-y);
}
.btn.tonal.icon, .btn.icon:disabled {
  padding: calc(var(--density-padding-y) - 1px);
}
.btn.ghost.icon:disabled {
  padding: var(--density-padding-y);
}
.btn.icon svg {
  width: 100%;
  height: 100%;
  fill: currentColor;
}

/* No hover rule, in any variant, and the absence is measured: a Finder
   toolbar button tints under the pointer and a Save dialog's push button
   does not, so a push button on this platform answers the pointer only when
   it is held. Held, the platform lays its press overlay over the fill the
   variant carries - written as a one-colour gradient layer over the
   background colour, which is how CSS composites a coverage onto a fill in
   the same space Flatten does, and straight onto the page where the variant
   carries no fill, which is how a held ghost gets one at all. */
.btn:active, .btn.is-active {
  background-image: linear-gradient(var(--platform-press-overlay), var(--platform-press-overlay));
}

/* Keyboard focus: the platform's own indicator, inset in the control's
   outermost 2 dp so the button's box does not grow, and the same ring at the
   same width in every variant - keyboard visibility is not a prominence
   property. The coverage composites over whatever the ring lies on, which is
   the fill where the variant has one and the page where it has none: exactly
   what focus.Ring is handed on the Gio side. */
.btn:focus-visible, .btn.is-focus {
  outline: var(--focus-ring-width) solid var(--platform-keyboard-focus-indicator);
  outline-offset: calc(var(--focus-ring-width) / -2);
}

/* Disabled is the platform's own answer rather than a fade of the resting
   colours: a variant that carries a fill falls back to the push button's
   fill inside the separator hairline, and every foreground becomes the
   platform's disabled control text. The padding gives back the hairline's
   1px, as everywhere else in this sheet, so the drawn box does not grow. A
   ghost keeps its absence of fill: there is nothing to fall back to. */
.btn:disabled {
  cursor: default;
  background: var(--platform-push-button-fill);
  border: 1px solid var(--platform-separator);
  padding: calc(var(--density-padding-y) - 1px) calc(var(--density-padding-x) - 1px);
  color: var(--platform-disabled-control-text);
}
.btn.ghost:disabled {
  background: transparent;
  border: none;
  padding: var(--density-padding-y) var(--density-padding-x);
  color: var(--platform-disabled-control-text);
}

/* ---- Badge ----
   The inline annotation components/badge draws: the system's own word about
   a thing. A status is the platform's system colour for it, and the badge is
   that colour filled with its content in white — alternateSelectedControlText,
   which is white in both appearances — which is how the platform draws a
   count badge. The fill does not vary with the status beyond the hue and the
   white does not vary at all; never invert the pair, and never tint the
   system colour toward the surface. Neutral is systemGray, the platform
   giving no status colour for "no status".

   It is off the control metrics — no boundary, no minimum height, no vertical
   padding — so its whole height is the label role's line box. The side
   padding is the S2 stop and the corner the radius scale's Base stop,
   deliberately not the pill, which is a chip's shape.

   The label role is the comfortable density's. The sheet re-maps no type role
   under .compact, so the class states the comfortable role and a compact page
   sets its badges in the same one.

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
  background: var(--platform-system-gray);
  color: var(--platform-alternate-selected-control-text);
}
.badge.success { background: var(--platform-system-green); }
.badge.warning { background: var(--platform-system-orange); }
.badge.error { background: var(--platform-system-red); }
.badge.info { background: var(--platform-system-blue); }

/* ---- Form controls ----
   Native elements wearing components/input's resolution: the platform's text
   background under its text colour, the measured field edge as the resting
   hairline, the placeholder text for a prompt, the keyboard focus indicator
   focused, and the platform's disabled control text where a control cannot be
   used — the fill staying exactly where it was, because the platform fades
   the wording and leaves the control.

   background-clip is the padding box throughout this family, and that is the
   geometry rather than a style: components/input draws the edge as an outer
   shape with the fill inset inside it, so a focus ring is the platform's
   indicator over the SURFACE the control stands on, not over the control's
   own fill. Clipping the background to the padding box is how CSS composites
   the same pixel. */

/* Text field (components/input textfield.go). The height floor is
   Density.FieldHeight, not ControlHeight — the platform draws a field shorter
   than the button beside it — and the drawn height is max(that, line box +
   2*PaddingY). The horizontal inset is S3 (12 dp, static: it does not follow
   density) measured from the OUTER edge, so the padding gives back whatever
   the border occupies and the text lands where the Gio side puts it. */
.input {
  box-sizing: border-box;
  display: block;
  width: 100%;
  margin: 0;
  appearance: none;
  min-height: var(--density-field-height);
  padding: calc(var(--density-padding-y) - 1px) calc(var(--space-3) - 1px);
  border: 1px solid var(--platform-field-edge);
  border-radius: var(--radius-md);
  background: var(--platform-text-background);
  background-clip: padding-box;
  color: var(--platform-text);
  font-family: var(--font-family);
  font-size: var(--font-body-large-size);
  line-height: var(--font-body-large-line-height);
  font-weight: var(--font-body-large-weight);
  letter-spacing: var(--font-body-large-tracking);
}
.input::placeholder { color: var(--platform-placeholder-text); opacity: 1; }

/* Focus replaces the edge with the ring and draws it at focus.Width, the
   2 dp the Gio side draws, with the padding giving the two pixels back so the
   field's outer geometry and its text position do not move. */
.input:focus-visible, .input.is-focus {
  outline: none;
  border-width: var(--focus-ring-width);
  border-color: var(--platform-keyboard-focus-indicator);
  padding: calc(var(--density-padding-y) - var(--focus-ring-width)) calc(var(--space-3) - var(--focus-ring-width));
}
.input:disabled {
  color: var(--platform-disabled-control-text);
}
.input:disabled::placeholder {
  color: var(--platform-disabled-control-text);
}

/* Dropdown (components/input dropdown.go, drawn by components/picker's field
   trigger): a trigger is a BUTTON and not a field, so it takes the push
   button's own fill inside the separator hairline, the control height as its
   floor, and the control text — a prompt standing in for an unmade choice
   takes the placeholder instead. Its right side reserves S3 + the 16 dp
   chevron + S3, the same inset drawTrigger keeps clear of the label. The
   chevron is drawn by the .select-wrap wrapper — a native select cannot carry
   a generated child — as a border-built triangle 16 dp across and 8 dp tall
   (drawChevron's half/quarter geometry) in the platform's secondary label.
   The whole family's padding-box clip is overridden here: this edge is the
   separator over the trigger's own fill, so the fill paints under it. */
.select {
  min-height: var(--density-control-height);
  padding-right: calc(var(--space-3) * 2 + 16px - 1px);
  border-color: var(--platform-separator);
  background-color: var(--platform-push-button-fill);
  background-clip: border-box;
  color: var(--platform-control-text);
}
.select:focus-visible, .select.is-focus {
  border-color: var(--platform-keyboard-focus-indicator);
  padding-right: calc(var(--space-3) * 2 + 16px - var(--focus-ring-width));
}
.select:disabled { color: var(--platform-disabled-control-text); }
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
  border-top: 8px solid var(--platform-secondary-label);
  pointer-events: none;
}
.select-wrap:has(.select:disabled)::after {
  border-top-color: var(--platform-disabled-control-text);
}

/* Checkbox (components/input checkbox.go): a 16 dp glyph (checkboxBoxSize,
   measured off save-dialog-{light,dark}.png — a component constant, not a
   token; it does not follow density) over the platform's text background,
   inside the 2 dp field edge every control in this row wears. Checked, the
   box is the platform's accent under a check mark in
   alternateSelectedControlText, because a fill says a colour was applied and
   only the mark says what that means: a column of fills carries completion in
   hue alone, which is the one channel a reader may not have.

   The mark is drawn, not encoded. Gio strokes the icon set's centre line —
   (4.5,12) -> (9,16.5) -> (19.5,6) on the set's 24-unit grid, a 2-unit
   DIAGONAL band, round caps and joins — and at the 16 px glyph one grid unit
   is 2/3 px, so the band is 1.333 px wide (+/-0.667 either side of the centre)
   and the arms run from (3,8) to (6,11) to (13,4). Each arm is one background
   layer: a linear-gradient banding its own box perpendicular to the arm,
   45deg for the short "\" arm and 135deg for the long "/" one. Every stop is
   written from 50% because each box is sized to its arm — the segment grown
   by half a band along its own axis, which makes the arm the box's diagonal
   and the box's corners the round caps' own tips:
     short arm: 3.943 px square at 2.529,7.529   (3,8)->(6,11)
     long arm:  7.943 px square at 5.529,3.529   (6,11)->(13,4)
   CSS has no line cap, so the caps come out cut square inside those tips
   rather than rounded — the same trade the icon set's own SVG files make when
   they draw their caps as an explicit contour, and a sub-pixel one at this
   size. background-origin is the border box so the grid is the 16 px glyph
   the Gio side scales on, not the 12 px inside the edge. */
.checkbox, .radio {
  box-sizing: border-box;
  appearance: none;
  flex: none;
  width: 16px;  /* checkboxBoxSize / radioCircleSize: 16 dp, measured */
  height: 16px;
  margin: 0;
  border: 2px solid var(--platform-field-edge);
  background: var(--platform-text-background);
  background-clip: padding-box;
  cursor: pointer;
}
.checkbox:focus-visible, .checkbox.is-focus,
.radio:focus-visible, .radio.is-focus {
  outline: var(--focus-ring-width) solid var(--platform-keyboard-focus-indicator);
  outline-offset: 1px;
}
.checkbox {
  border-radius: var(--radius-sm);
}
.checkbox:checked, .checkbox.is-checked {
  border-color: var(--platform-control-accent);
  background-color: var(--platform-control-accent);
  background-clip: border-box;
  background-image:
    linear-gradient(45deg, transparent calc(50% - 0.667px), var(--platform-alternate-selected-control-text) calc(50% - 0.667px), var(--platform-alternate-selected-control-text) calc(50% + 0.667px), transparent calc(50% + 0.667px)),
    linear-gradient(135deg, transparent calc(50% - 0.667px), var(--platform-alternate-selected-control-text) calc(50% - 0.667px), var(--platform-alternate-selected-control-text) calc(50% + 0.667px), transparent calc(50% + 0.667px));
  background-origin: border-box;
  background-repeat: no-repeat;
  background-position: 2.529px 7.529px, 5.529px 3.529px;
  background-size: 3.943px 3.943px, 7.943px 7.943px;
}
/* Disabled: unchecked, the platform's disabled control text takes the edge
   and the fill stays; checked, the accent fill is the accent at the disabled
   fraction of its own coverage over the surface, with the mark following it. */
.checkbox:disabled, .radio:disabled {
  cursor: default;
  border-color: var(--platform-disabled-control-text);
}
.checkbox:checked:disabled, .checkbox.is-checked:disabled {
  border-color: color-mix(in srgb, var(--platform-control-accent) var(--state-disabled-opacity), transparent);
  background-color: color-mix(in srgb, var(--platform-control-accent) var(--state-disabled-opacity), transparent);
  background-image:
    linear-gradient(45deg, transparent calc(50% - 0.667px), var(--platform-disabled-control-text) calc(50% - 0.667px), var(--platform-disabled-control-text) calc(50% + 0.667px), transparent calc(50% + 0.667px)),
    linear-gradient(135deg, transparent calc(50% - 0.667px), var(--platform-disabled-control-text) calc(50% - 0.667px), var(--platform-disabled-control-text) calc(50% + 0.667px), transparent calc(50% + 0.667px));
}

/* Radio (components/input radio.go): the same 16 dp glyph as a circle.
   Selected is the platform's accent filling the whole circle with an 8 dp dot
   (radioDotSize) in alternateSelectedControlText at its centre — one fill and
   one mark, exactly the Gio nested ellipses, and no gap ring. */
.radio { border-radius: var(--radius-full); }
.radio:checked, .radio.is-checked {
  border-color: var(--platform-control-accent);
  background: radial-gradient(circle, var(--platform-alternate-selected-control-text) 4px, var(--platform-control-accent) 4px); /* 8 dp dot */
  background-clip: border-box;
}
.radio:checked:disabled, .radio.is-checked:disabled {
  border-color: color-mix(in srgb, var(--platform-control-accent) var(--state-disabled-opacity), transparent);
  background: radial-gradient(circle, var(--platform-disabled-control-text) 4px, color-mix(in srgb, var(--platform-control-accent) var(--state-disabled-opacity), transparent) 4px);
}

/* ---- Card and group ----
   The two are one ruling with two answers, and both are the platform's.
   patterns/card singles something out: the grouped box System Settings draws
   — a small step of fill from the surface it stands on, darker in the light
   appearance and lighter in the dark one, with NO hairline and no shadow. The
   step is the whole of what does the singling out. patterns/group divides the
   page: a separator hairline around related components, taking the fill of
   the surface the group is in, so it declares no background at all.

   Both carry radius Lg, an S4 inset and S3 gaps between the slots — exactly
   drawCard's rad.Lg / sp.S4 / sp.S3. The group's line lies inside its bounds
   where the CSS border does, so its padding gives the border's 1px back and
   the slots land where the Gio inset puts them; the card has no line to give
   back. Neither styles slot text of its own: the Gio card draws no text at
   all, and the group draws only its own label. */
.card {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: var(--space-4);
  border: none;
  border-radius: var(--radius-lg);
  background: var(--platform-card-fill);
  color: var(--platform-label);
}
.group {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: calc(var(--space-4) - 1px);
  border: 1px solid var(--platform-separator);
  border-radius: var(--radius-lg);
  color: var(--platform-label);
}
/* The group's own label: top-leading, inside the hairline, as the first row
   of the group's stack — the platform's idiom for a section header over a
   bordered container, and not the fieldset legend cut into the top line,
   which has no native counterpart and does not survive a Lg corner. The
   label-large role in the platform's secondary label: a group's own name is
   set under the content it names. */
.group-label {
  font-size: var(--font-label-large-size);
  line-height: var(--font-label-large-line-height);
  font-weight: var(--font-label-large-weight);
  letter-spacing: var(--font-label-large-tracking);
  color: var(--platform-secondary-label);
}

/* ---- Table ----
   patterns/table: the grid is printed on the platform's content fill, which
   is what a list, a table and a text view all stand on there, and the header
   band takes no fill of its own — on this platform a table's header is the
   content's fill under the header text, closed by a separator along its foot.
   Header and body rows are each exactly Density.RowHeight tall — the
   platform's stacked row, 20 dp comfortable, which is shorter than a control
   — so .compact re-pitches the whole grid. Every body row closes with the
   platform's grid colour drawn inside its own height; the header's foot is
   the separator over the content fill. Cells inset horizontally by S3
   (cellPadDp, 12 dp — static, the same inset rule the input uses); body text
   is body-medium in the platform's label (RenderTextCell), header text the
   platform's header text. Every second body row is laid on the platform's
   alternating content background, which is what Finder's list view draws.
   Sort is a 10 dp chevron in the header's own foreground at
   the header's right inset (drawSortChevron), drawn on the active column only
   — .sort-asc / .sort-desc on a .sortable header. */
.table {
  box-sizing: border-box;
  border-collapse: separate;
  border-spacing: 0;
  table-layout: fixed;
  width: 100%;
  background: var(--platform-control-background);
  color: var(--platform-label);
  font-family: var(--font-family);
  font-size: var(--font-body-medium-size);
  line-height: var(--font-body-medium-line-height);
  font-weight: var(--font-body-medium-weight);
  letter-spacing: var(--font-body-medium-tracking);
}
.table th, .table td {
  box-sizing: border-box;
  height: var(--density-row-height);
  padding: 0 var(--space-3);
  box-shadow: inset 0 -1px 0 var(--platform-grid);
  text-align: left;
  vertical-align: middle;
  white-space: nowrap;
  overflow: hidden;
}
/* The rule that closes a row is an inset shadow rather than a border,
   because a table cell's height is its content's floor: a 1 dp border would
   add a pixel to every row and walk the whole grid off the platform's pitch.
   Drawn inside the row's own height, which is where drawRow draws it. */
.table tbody tr:nth-child(even) td {
  background: var(--platform-alternating-content-background);
}
.table th {
  position: relative;
  background: var(--platform-control-background);
  box-shadow: inset 0 -1px 0 var(--platform-separator);
  color: var(--platform-header-text);
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
.table th.sort-asc::after { border-bottom: 5px solid var(--platform-header-text); }  /* 5 dp tall, apex up */
.table th.sort-desc::after { border-top: 5px solid var(--platform-header-text); }

/* ---- Navigation ----
   The four navigation patterns. All four stand on the chrome: the platform's
   sidebar material, which is a shade off the content in both appearances and
   is what a navbar, a tab strip and a rail are filled with. Two flush regions
   cannot be told apart by fill in the light appearance — the chrome material
   IS the content's white there — so the region above or leading draws the
   separator that says where it ends, inside its own bounds: the navbar's
   foot, the tab strip's foot, the rail's trailing edge. None of the four
   tints under the pointer: the platform tints neither a list row nor a
   sidebar row, and this is the measured absence rather than an omission. What
   they do draw is selection — the 2 dp underline on a link or a tab, the
   selected row's own fill on a rail — in the platform's selection colour,
   which follows the user's accent. */

/* Navbar (patterns/navbar navbar.go): a horizontal chrome bar — drawNavbar
   fills the material, insets PaddingY vertically and S4 horizontally, and
   patterns/shell pins the bar to ControlHeight + 2*PaddingY (28 dp
   comfortable, 19 compact). Slots run brand, centred links, actions; the
   links row centres in the space brand and actions leave over (that space
   halved), which margin-inline auto reproduces exactly. The foot hairline is
   an inset shadow rather than a border, because the Gio side paints it over
   the bar's own bottom pixel without insetting the slots above it. */
.navbar {
  box-sizing: border-box;
  display: flex;
  align-items: center;
  min-height: calc(var(--density-control-height) + 2 * var(--density-padding-y));
  padding: var(--density-padding-y) var(--space-4);
  background: var(--platform-sidebar-material);
  box-shadow: inset 0 -1px 0 var(--platform-separator);
  color: var(--platform-label);
}
.navbar-links {
  display: flex;
  align-items: center;
  gap: var(--space-2);  /* linksRow's HSpacer(sp.S2) between link cells */
  margin-inline: auto;  /* the leftover space, halved on either side */
}

/* A link cell (navbar.go linkWidget): label-large in the platform's label
   inside (S3, PaddingY) padding, with a 2 dp underline slot along the bottom
   edge that the Active link fills with the platform's selection colour — the
   underline runs the full cell width, padding included, exactly the
   image.Rect the Gio side fills. */
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
  color: var(--platform-label);
  border-bottom: 2px solid transparent;  /* underlineDp: the underline slot */
}
.navbar-link {
  padding: var(--density-padding-y) var(--space-3);
}
.navbar-link.selected, .tab.selected {
  border-bottom-color: var(--platform-selected-content-background);
}

/* Tabs (patterns/tabs tabs.go): the strip is a chrome row exactly
   ControlHeight tall (drawTabs pins stripH to the density), closed by the
   separator along its foot, each cell its label plus 2*S3 horizontal padding
   with the label centred in the height that remains above the 2 dp underline
   slot — which border-box centring reproduces. The selected cell fills the
   slot with the platform's selection colour; content panels below the strip
   are the caller's, and stand on the content fill. */
.tabs {
  box-sizing: border-box;
  display: flex;
  align-items: stretch;
  height: var(--density-control-height);
  background: var(--platform-sidebar-material);
  box-shadow: inset 0 -1px 0 var(--platform-separator);
}
.tab {
  height: 100%;
  padding: 0 var(--space-3);
  justify-content: center;
}

/* Sidebar (patterns/sidebar sidebar.go): a vertical chrome rail at the
   pattern's two contractual widths — 192 dp expanded, 48 dp collapsed
   (expandedDp/collapsedDp: component constants, deliberately not tokens and
   not density-responsive; a different rail copies the pattern) — closed by
   the separator along its trailing edge. The toggle row is ControlHeight;
   every item row is the sidebar's OWN row height, the pattern's RowHeight
   rather than the density scale's, because a chrome rail draws a taller row
   than a content list. The whole rail is one keyboard stop — the
   focus ring belongs to the rail, and the Arrow keys move .selected — so the
   ring rule below targets .sidebar itself, not the rows. */
.sidebar {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  width: 192px;  /* expandedDp */
  background: var(--platform-sidebar-material);
  box-shadow: inset -1px 0 0 var(--platform-separator);
  color: var(--platform-label);
  overflow: hidden;
}
.sidebar.collapsed { width: 48px; }  /* collapsedDp */

/* The collapse affordance (drawToggle): a full-width ControlHeight row with
   the icon set's sidebar mark centred in it at icon.Size — the control's
   inner content box, ControlHeight - 2*PaddingY — in the platform's secondary
   label. Pointer-only, never a Tab stop. */
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
  width: calc(var(--density-control-height) - 2 * var(--density-padding-y));
  height: calc(var(--density-control-height) - 2 * var(--density-padding-y));
  background: var(--platform-secondary-label);
}

/* An item row (drawItem): a 48 dp leading icon column (iconColDp) with the
   glyph centred in it, the label-large label starting at exactly the column
   edge, vertically centred, one line, clipped rather than wrapped — which is
   also what hides the labels at the collapsed width. Selected wears the
   platform's sidebar pill; nothing else moves, and nothing tints under the
   pointer.

   The row height, the pill's inset and its corner are MEASURED off the
   organization's macOS reference (patterns/sidebar RowHeight,
   SelectionInset, SelectionRadius): Finder's and Voice Memos' selected
   sidebar rows span 32 px at 1x, the pill is inset 10 from each edge of the
   rail, and a circular fit to its corner reads 8. */
.sidebar-item {
  box-sizing: border-box;
  position: relative;
  z-index: 0;
  flex: none;
  display: flex;
  align-items: center;
  height: 32px;  /* RowHeight */
  overflow: hidden;
  white-space: nowrap;
  cursor: pointer;
  font-family: var(--font-family);
  font-size: var(--font-label-large-size);
  line-height: var(--font-label-large-line-height);
  font-weight: var(--font-label-large-weight);
  letter-spacing: var(--font-label-large-tracking);
  color: var(--platform-label);
  text-decoration: none;
  user-select: none;
}
/* The pill is a layer behind the row's own content rather than the row's
   fill, because it is inset from the rail while the icon column and the
   label are not. */
.sidebar-item.selected {
  color: var(--platform-alternate-selected-control-text);
}
.sidebar-item.selected::before {
  content: "";
  position: absolute;
  inset: 0 10px;  /* SelectionInset */
  z-index: -1;
  border-radius: 8px;  /* SelectionRadius */
  background: var(--platform-control-accent);
}
.sidebar-item-icon {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;  /* iconColDp */
  height: 100%;
}

/* Breadcrumb (components/breadcrumb breadcrumb.go): a row of title-small
   segments with S2 gaps around 12 dp chevron separators (chevronDp). The last
   segment is the current location in the platform's label; every ancestor is
   a link and takes the platform's link colour. Colour follows position
   (labelColor), so :last-child carries it; .current forces it for a specimen. */
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
  color: var(--platform-link);
  text-decoration: none;
  white-space: nowrap;
}
.crumbs .crumb:last-child, .crumb.current { color: var(--platform-label); }

/* The chevron separator (chevronWidget/drawChevron): a right-pointing
   triangle half the 12 dp box wide and the full box tall, centred in it, in
   the platform's secondary label. */
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
  border-left: 6px solid var(--platform-secondary-label);  /* 6 dp deep, apex along +X */
}

/* The keyboard ring, identical to every other control's: per-cell for the
   navbar, tabs and breadcrumb (each cell is its own Clickable focus tag); on
   the rail itself for the sidebar, whose single stop is the item list. */
.navbar-link:focus-visible, .navbar-link.is-focus,
.tab:focus-visible, .tab.is-focus,
.crumb:focus-visible, .crumb.is-focus,
.sidebar:focus-visible, .sidebar.is-focus {
  outline: var(--focus-ring-width) solid var(--platform-keyboard-focus-indicator);
  outline-offset: calc(var(--focus-ring-width) / -2);
}

/* ---- Overlays ----
   The transient surfaces: the scrimmed dialog (patterns/modal), the
   unscrimmed popover (patterns/popover), the tooltip (components/tooltip) and
   the floating toast (components/toast, stood in a column by
   patterns/notifications). On this platform a floating surface is filled with
   the window's own background, and what says it floats is the measured
   shadow — not a lighter fill, and not a level. The popover and the tooltip
   carry the separator around theirs as well, because a still surface flush
   against the plane behind it needs the line; the dialog does not, the scrim
   being what parts it.

   The shadow is a measurement rather than a blur radius: it stands at the
   platform's floating shadow at the surface's edge and falls LINEARLY to
   nothing 24 px out (effects/depth, fitted to the stored sidebar-shadow
   captures). CSS has no linear falloff, so it is approximated by eight
   stacked spreads, each carrying an eighth of the coverage, which reproduces
   the ramp to within one level of the measured one. */
.dialog, .popover, .toast {
  --floating-shadow-step: color-mix(in srgb, var(--platform-floating-shadow) 12.5%, transparent);
  box-shadow:
    0 0 0 3px var(--floating-shadow-step),
    0 0 0 6px var(--floating-shadow-step),
    0 0 0 9px var(--floating-shadow-step),
    0 0 0 12px var(--floating-shadow-step),
    0 0 0 15px var(--floating-shadow-step),
    0 0 0 18px var(--floating-shadow-step),
    0 0 0 21px var(--floating-shadow-step),
    0 0 0 24px var(--floating-shadow-step);
}

/* Scrim (modal.go drawModal): the whole-plane dimmer under a dialog — the
   platform's own scrim coverage, identical in both appearances, because a
   scrim dims by reducing luminance and never flips with the scheme. The scrim
   centres the dialog, exactly as drawModal centres the surface in the window
   plane. Behaviour is part of the pattern: on a PANEL a scrim press invokes
   OnClose; on a DECISION the scrim is INERT — it absorbs presses and answers
   none of them, because dismissal is one of the decision's answers and a
   stray click must not give it. */
.scrim {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  background: var(--platform-scrim);
}

/* Dialog (modal.go drawModal): the centred surface — width 75% of the window
   plane clamped to 180-560 dp, height hugging its content between the 120 dp
   floor and the 560 dp cap (overflow clips), the window background under the
   platform's shadow and no hairline at all, radius Lg, an S5 inset and S3
   gaps between header, body and footer. G0A.2's two purposes share this one
   surface: a PANEL carries a ghost icon close (.btn.ghost.icon) in its header
   and no footer of its own; a DECISION carries no X anywhere and a
   .dialog-footer whose right-aligned actions end in the Return-bound
   default. */
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
  padding: var(--space-5);
  border: none;
  border-radius: var(--radius-lg);
  background: var(--platform-window-background);
  color: var(--platform-label);
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
   content plus an S3 inset (the padding gives back the hairline's 1px),
   clamped to the 48x24 dp minimum, the window background under the platform's
   shadow inside the platform's separator, radius Md. Positioning against the
   anchor is the page's: the pattern centres the surface on the anchor's
   midline, one S2 gap away on the Placement side. */
.popover {
  box-sizing: border-box;
  display: inline-block;
  min-width: 48px;
  min-height: 24px;
  padding: calc(var(--space-3) - 1px);
  border: 1px solid var(--platform-separator);
  border-radius: var(--radius-md);
  background: var(--platform-window-background);
  color: var(--platform-label);
}

/* The tail (popover.go drawTail): a triangle 12 dp across the base and 6 dp
   deep in the surface's own fill, bridging the gap with its tip at the
   anchor. The modifier names the popover's Placement — a .top popover sits
   above its anchor, so its tail points down. */
.popover-tail { width: 0; height: 0; }
.popover-tail.top {
  border-left: 6px solid transparent;   /* 12 dp base */
  border-right: 6px solid transparent;
  border-top: 6px solid var(--platform-window-background);
}
.popover-tail.bottom {
  border-left: 6px solid transparent;
  border-right: 6px solid transparent;
  border-bottom: 6px solid var(--platform-window-background);
}
.popover-tail.left {
  border-top: 6px solid transparent;
  border-bottom: 6px solid transparent;
  border-left: 6px solid var(--platform-window-background);
}
.popover-tail.right {
  border-top: 6px solid transparent;
  border-bottom: 6px solid transparent;
  border-right: 6px solid var(--platform-window-background);
}

/* Tooltip (components/tooltip tooltip.go drawSurface): the window background
   inside the platform's separator under a label in the platform's label
   colour, label-small, radius Sm, S2/S1 padding measured from the outer edge,
   clamped to the 24x16 dp minimum. It casts no shadow: the hairline is the
   whole of a still tooltip's edge, and what places it draws whatever shadow
   the placement owes. The same in both appearances. */
.tooltip {
  box-sizing: border-box;
  display: inline-flex;
  align-items: center;
  min-width: 24px;
  min-height: 16px;
  padding: calc(var(--space-1) - 1px) calc(var(--space-2) - 1px);
  border: 1px solid var(--platform-separator);
  border-radius: var(--radius-sm);
  background: var(--platform-window-background);
  color: var(--platform-label);
  white-space: nowrap;
  font-family: var(--font-family);
  font-size: var(--font-label-small-size);
  line-height: var(--font-label-small-line-height);
  font-weight: var(--font-label-small-weight);
  letter-spacing: var(--font-label-small-tracking);
}

/* Toast: one queued notification — 240 dp wide, a 36 dp legibility floor that
   deliberately does not follow density (a toast is not a control), radius Md,
   label-medium in the platform's label over the window background, and
   floating on the platform's measured shadow. The status is carried on the
   leading edge alone, one S2 wide, in the platform's system colour for it —
   painted as a two-stop gradient so the chip's own radius rounds it. A toast
   says one thing about one event, and the colour of that edge is the whole of
   what it says; the message stays in the reading colour every other surface
   uses. */
.toast {
  box-sizing: border-box;
  display: flex;
  align-items: center;
  width: 240px;
  min-height: 36px;
  padding: var(--space-2) var(--space-3);
  padding-left: calc(var(--space-2) + var(--space-3));
  border-radius: var(--radius-md);
  background: linear-gradient(to right, var(--platform-system-blue) 0 var(--space-2), var(--platform-window-background) var(--space-2));
  color: var(--platform-label);
  font-family: var(--font-family);
  font-size: var(--font-label-medium-size);
  line-height: var(--font-label-medium-line-height);
  font-weight: var(--font-label-medium-weight);
  letter-spacing: var(--font-label-medium-tracking);
}
.toast.success {
  background: linear-gradient(to right, var(--platform-system-green) 0 var(--space-2), var(--platform-window-background) var(--space-2));
}
.toast.warning {
  background: linear-gradient(to right, var(--platform-system-orange) 0 var(--space-2), var(--platform-window-background) var(--space-2));
}
.toast.error {
  background: linear-gradient(to right, var(--platform-system-red) 0 var(--space-2), var(--platform-window-background) var(--space-2));
}

/* The column (notifications.go paintColumn): a corner-anchored column with
   S2 gaps, inset S4 from the window plane's edges (the page anchors it);
   newest toast nearest the anchored edge. */
.toast-stack {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  width: 240px;
}
`
