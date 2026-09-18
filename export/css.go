package export

import (
	"fmt"
	stdcolor "image/color"
	"math"
	"strconv"
	"strings"
	"time"

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
//
// A test walks tokens.PlatformColors by reflection and fails if a field has
// no entry here: a name the set carries and the sheet drops is a colour a
// rule can reference and never resolve.
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
	{"sidebar-selection", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SidebarSelection }},
	{"sidebar-count", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SidebarCount }},
	{"card-fill", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.CardFill }},
	{"push-button-fill", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.PushButtonFill }},
	{"hover-overlay", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.HoverOverlay }},
	{"press-overlay", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.PressOverlay }},
	{"floating-shadow", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.FloatingShadow }},
	{"field-edge", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.FieldEdge }},
	{"scrollbar-thumb", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ScrollbarThumb }},
	{"alternating-content-background", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.AlternatingContentBackground }},
	{"scrim", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Scrim }},
	{"sidebar-search-fill", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SidebarSearchFill }},
	{"toolbar-control-fill", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ToolbarControlFill }},
	{"toolbar-control-rim", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ToolbarControlRim }},
	{"toolbar-search-fill", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ToolbarSearchFill }},
	{"toolbar-search-rim", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ToolbarSearchRim }},
	{"toolbar-control-shadow", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ToolbarControlShadow }},
	{"toolbar-checked-overlay", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ToolbarCheckedOverlay }},
	{"pane-rim", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.PaneRim }},
	{"pane-shadow", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.PaneShadow }},
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

// The bordered toolbar control's drop shadow is the one thing this sheet
// draws whose GEOMETRY flips with the appearance rather than only its colour,
// so its reach and offset are emitted per appearance beside the platform's
// colour set. Its peak coverage is
// --platform-toolbar-control-shadow, a platform name like every other.
//
// MEASURED, components/internal/control's fit, restated here because the
// module graph runs the other way (components imports theme): light,
// finder-window-light.png, 23 px of reach with the shadow's rectangle sunk
// 9 px, which is what tells a #ffffff control from a #ffffff band; dark,
// finder-window-untinted-dark.png and notes-toolbar.png, 2 px of reach with
// the rectangle sunk 6 px, the band one 255th deep over seven rows. Which
// reading answers is the platform's own behaviour and not an appearance this
// code tests for: where the platform gives the control a rim the control is
// told from its band by that rim and its fill, and the shadow is a hint sunk
// under it; where it gives none the shadow is the whole of the step.
const (
	toolbarShadowReachLightDp  = 23
	toolbarShadowOffsetLightDp = 9
	toolbarShadowReachDarkDp   = 2
	toolbarShadowOffsetDarkDp  = 6
)

// toolbarShadowVars renders that geometry for one appearance, told apart by
// whether the platform draws the control a rim there.
func toolbarShadowVars(p tokens.PlatformColors) []cssVar {
	reach, offset := float32(toolbarShadowReachLightDp), float32(toolbarShadowOffsetLightDp)
	if p.ToolbarControlRim.A != 0 {
		reach, offset = toolbarShadowReachDarkDp, toolbarShadowOffsetDarkDp
	}
	return []cssVar{
		{"--toolbar-control-shadow-reach", px(reach)},
		{"--toolbar-control-shadow-offset", px(offset)},
	}
}

// typeRoles orders the fifteen type roles under their CSS names, plus
// code — the sixteenth style outside the type scale, the mono face at
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

// shadowLevels orders the six levels from the backdrop up toward the
// reader, with the shadow depth each casts. A level's FILL is the
// platform's own name for what that region is; the shadow is the separate,
// opt-in cue a floating surface carries, and this is the only thing the
// sheet states per level.
//
// The backdrop and the chrome level are spelled out rather than numbered
// because the numbering is anchored on the content and they are below it:
// naming them "-2" and "-1" in a CSS variable would read as an arithmetic
// accident, and renumbering the four above them would rename every token to
// say the same thing.
var shadowLevels = []struct {
	name string
	dp   func(tokens.ElevationScale) float32
}{
	{"backdrop", func(e tokens.ElevationScale) float32 { return e.Backdrop }},
	{"chrome", func(e tokens.ElevationScale) float32 { return e.Chrome }},
	{"0", func(e tokens.ElevationScale) float32 { return e.Level0 }},
	{"1", func(e tokens.ElevationScale) float32 { return e.Level1 }},
	{"2", func(e tokens.ElevationScale) float32 { return e.Level2 }},
	{"3", func(e tokens.ElevationScale) float32 { return e.Level3 }},
}

// densityMetrics orders the per-setting density metrics under their CSS
// names.
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
	// The checkbox's row: the square footprint the 16 dp glyph is centred in
	// and the pointer target both the checkbox and the radio offer. Measured
	// at 22 against the push button's 24 and the list row's 20, so it is a
	// number of its own and not either of theirs.
	{"checkbox-row-height", func(d tokens.Density) float32 { return d.CheckboxRowHeight }},
	// The toolbar control's own height. A bordered control standing in a
	// toolbar band is 36 on this platform against the dialog control's 24,
	// measured, so the sheet states it rather than leaving a consumer to
	// reach for the control height there.
	{"toolbar-control-height", func(d tokens.Density) float32 { return d.ToolbarControlHeight }},
	{"padding-x", func(d tokens.Density) float32 { return d.PaddingX }},
	{"padding-y", func(d tokens.Density) float32 { return d.PaddingY }},
}

// easeRoles orders the easing presets under their CSS names.
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

// scaleVars renders the mode-invariant families: fonts, density
// (comfortable — the :root setting), spacing, radius, the dp shadows (the
// opt-in cue for floating transients), and the motion set.
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
	for _, key := range spaceKeys {
		vars = append(vars, cssVar{"--space-" + key.name, px(key.pick(s.Spacing))})
	}
	for _, key := range radiusKeys {
		vars = append(vars, cssVar{"--radius-" + key.name, px(key.pick(s.Radius))})
	}
	for _, level := range shadowLevels {
		vars = append(vars, cssVar{"--shadow-" + level.name, boxShadow(level.dp(s.Elevation))})
	}
	for _, role := range easeRoles {
		vars = append(vars, cssVar{"--ease-" + role.name, cubicBezier(role.pick(s.Motion))})
	}
	for _, stop := range durationStops {
		vars = append(vars, cssVar{"--duration-" + stop.name, ms(stop.pick(s.Motion))})
	}
	// The focus halo's band width, mode-invariant, which is why it is here
	// and the halo's COLOUR is not: the halo wears
	// --platform-keyboard-focus-indicator, which flips with the appearance.
	// The platform's measured disabled coverage sits beside it: a switched-off
	// control is its own fill at that coverage over the surface it stands on,
	// which is a number rather than a colour and belongs to no appearance.
	vars = append(vars,
		cssVar{"--focus-halo-width", px(focusHaloWidthDp)},
		cssVar{"--disabled-coverage", coverage(tokens.DisabledCoverage)},
	)
	return vars
}

// focusHaloWidthDp is the focus halo's band width — the 4 dp every control in
// this library draws (components/internal/focus.Width), half of it past the
// control's own box and half over the control's outermost band. It is
// identical in every variant because keyboard visibility is not a matter of
// prominence, and identical at every density because a halo is a keyboard
// affordance rather than an ornament.
//
// MEASURED, save-dialog-{light,dark}.png, the focused "Save As:" field: four
// px on every side of a box running x 264-495, hard-edged, straddling the box.
const focusHaloWidthDp = 4

// coverage formats a 0-255 coverage as a CSS percentage, which is the form
// color-mix() takes it in: a rule states the platform's own colour at the
// platform's own coverage and lets the browser composite it over whatever is
// beneath, exactly as theme/color.Fade and Flatten do on the Gio side.
func coverage(a uint8) string {
	return strconv.FormatFloat(math.Round(float64(a)/255*1e4)/100, 'f', -1, 64) + "%"
}

// densityVars renders one density setting's per-setting metrics. The :root
// block carries tokens.Comfortable's; the .compact override block carries
// tokens.Compact's. A control's pointer target is not among them because it
// is not a metric of its own: a control's target is the box it draws, which
// --density-control-height already gives.
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

// stylesCSS renders the full token sheet: the platform's light set and
// every mode-invariant scale under :root, the platform's dark set under
// .dark, and the compact density metrics under .compact. The two class
// blocks are orthogonal switches — .dark flips the colours, .compact flips
// the per-setting density metrics — so a surface can be any of the four
// combinations.
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
	light := append(platformVars(s.PlatformLight), toolbarShadowVars(s.PlatformLight)...)
	block(&b, ":root", append(light, scaleVars(s)...))
	b.WriteString("\n")
	block(&b, ".dark", append(platformVars(s.PlatformDark), toolbarShadowVars(s.PlatformDark)...))
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
// components/input, components/badge and the patterns took. Nothing is
// derived here: every value is a platform name.
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
// The states are the platform's own answers. A press lays
// --platform-press-overlay over whatever fill the variant carries, and
// over the page where it carries none, which is how a ghost gets a fill at
// all. Focus is the halo: --platform-keyboard-focus-indicator at
// --focus-halo-width on the control's own outline, the same band in every
// variant. Disabled fades a fill to the push button's own at
// --disabled-coverage over the surface the control stands on and takes every
// foreground to --platform-disabled-control-text.
//
// The pointer overlays and the switched-off fade are the platform's answers
// too, and this sheet states them the way the components draw them: hover
// lays --platform-hover-overlay over whatever fill the variant carries, a
// press lays --platform-press-overlay there instead, and a switched-off
// control is the push button's own fill at --disabled-coverage over the
// surface it stands on. A coverage over a fill is what color-mix() renders
// and what theme/color.Fade computes, so the two sides land on the same
// pixel.
//
// No rule in this layer spends a type role's tracking. The library's typeset
// lays a label out at the role's size, weight and line height and spends no
// letter spacing at all, so a sheet that spent the role's tracking token
// would set every label a fraction wider than the component beside it.
// Whether the token should ever be spent is a typography question and is not
// answered here; what is answered is that the sheet draws what the library
// draws.
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
  letter-spacing: 0;
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

/* The bordered toolbar control (components/button's chrome variant, drawn
   through internal/toolbarface): a capsule at the toolbar band's own measured
   height - 36 against a dialog control's 24 - cornered at half of it, filled
   with the platform's measured toolbar control fill, rimmed with the
   platform's measured value for that rim (which answers no colour in the
   light appearance, where the platform draws none) and casting the measured
   drop shadow that tells a light control from a light band. The rim is an
   inset ring rather than a border so the control's box does not grow, exactly
   as toolbarface lays its band ON the shape's outline.

   The shadow's ramp is the linear falloff effects/depth draws, approximated
   by eight stacked spreads each carrying an eighth of the peak, sunk the
   measured offset below the control and reaching the measured reach out (the
   two --toolbar-control-shadow-* lengths). Both flip with the appearance,
   which is why they are emitted per appearance beside the colour set: this is
   the one thing here whose geometry the platform draws differently light and
   dark. */
.btn.chrome {
  min-height: var(--density-toolbar-control-height);
  height: var(--density-toolbar-control-height);
  border: none;
  border-radius: calc(var(--density-toolbar-control-height) / 2);
  padding: 0 var(--space-3);
  background: var(--platform-toolbar-control-fill);
  color: var(--platform-control-text);
  --toolbar-shadow-step: color-mix(in srgb, var(--platform-toolbar-control-shadow) 12.5%, transparent);
  box-shadow:
    0 var(--toolbar-control-shadow-offset) 0 calc(var(--toolbar-control-shadow-reach) / 8) var(--toolbar-shadow-step),
    0 var(--toolbar-control-shadow-offset) 0 calc(var(--toolbar-control-shadow-reach) / 4) var(--toolbar-shadow-step),
    0 var(--toolbar-control-shadow-offset) 0 calc(var(--toolbar-control-shadow-reach) * 3 / 8) var(--toolbar-shadow-step),
    0 var(--toolbar-control-shadow-offset) 0 calc(var(--toolbar-control-shadow-reach) / 2) var(--toolbar-shadow-step),
    0 var(--toolbar-control-shadow-offset) 0 calc(var(--toolbar-control-shadow-reach) * 5 / 8) var(--toolbar-shadow-step),
    0 var(--toolbar-control-shadow-offset) 0 calc(var(--toolbar-control-shadow-reach) * 3 / 4) var(--toolbar-shadow-step),
    0 var(--toolbar-control-shadow-offset) 0 calc(var(--toolbar-control-shadow-reach) * 7 / 8) var(--toolbar-shadow-step),
    0 var(--toolbar-control-shadow-offset) 0 var(--toolbar-control-shadow-reach) var(--toolbar-shadow-step),
    inset 0 0 0 1px var(--platform-toolbar-control-rim);
}
/* A chrome control whose whole label is a symbol: the measured 24 dp mark
   box with the measured 7 dp of clear room a side, so the capsule is 38
   across at the band's own height. */
.btn.chrome.icon {
  width: calc(24px + 2 * 7px);  /* ChromeMarkDp + 2 x ChromeMarkSideDp */
  height: var(--density-toolbar-control-height);
  padding: 0;
}
.btn.chrome.icon svg {
  width: 24px;  /* ChromeMarkDp */
  height: 24px;
  margin: auto;
}

/* Under the pointer the platform lays its hover overlay over whatever fill
   the variant carries, and held its press overlay over that same fill - each
   written as a one-colour gradient layer over the background colour, which is
   how CSS composites a coverage onto a fill in the same space Flatten does,
   and straight onto the page where the variant carries no fill, which is how
   a ghost gets one at all. A press wins over a hover: the two are one answer
   and not two laid on each other, so the held rule stands after the hovered
   one and replaces its layer. */
.btn:hover:not(:disabled), .btn.is-hover {
  background-image: linear-gradient(var(--platform-hover-overlay), var(--platform-hover-overlay));
}
.btn:active:not(:disabled), .btn.is-active {
  background-image: linear-gradient(var(--platform-press-overlay), var(--platform-press-overlay));
}

/* Keyboard focus: the platform's halo, a --focus-halo-width band lying on the
   control's own outline with half of it past the box and half over the box's
   outermost band, so the button's box does not grow and the fill and edge
   under it stay where they are. The same band at the same width in every
   variant - keyboard visibility is not a prominence property. The indicator
   carries a coverage, so the half past the box composites over the page and
   the half on the box over the control's own fill, which is exactly the pair
   of colours focus.Halo is handed on the Gio side. */
.btn:focus-visible, .btn.is-focus {
  outline: var(--focus-halo-width) solid var(--platform-keyboard-focus-indicator);
  outline-offset: calc(var(--focus-halo-width) / -2);
}

/* Disabled is the platform's fade: a variant that carries a fill falls back
   to the push button's own fill at --disabled-coverage over the surface the
   control stands on, its hairline is the separator at that same coverage, and
   every foreground becomes the platform's disabled control text. The fill is
   clipped to the padding box so the faded fill and the faded hairline each
   composite over the page rather than over one another, which is where
   control.Faded lands them. The padding gives back the hairline's 1px, as
   everywhere else in this sheet, so the drawn box does not grow. A ghost
   keeps its absence of fill: there is nothing to fall back to. */
.btn:disabled {
  cursor: default;
  background: color-mix(in srgb, var(--platform-push-button-fill) var(--disabled-coverage), transparent);
  background-clip: padding-box;
  border: 1px solid color-mix(in srgb, var(--platform-separator) var(--disabled-coverage), transparent);
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
  letter-spacing: 0;
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
  letter-spacing: 0;
}
.input::placeholder { color: var(--platform-placeholder-text); opacity: 1; }

/* Focus adds the halo and moves nothing: the field keeps its own edge, its
   own fill and its own insets, and wears the --focus-halo-width band on the
   box it already draws, half past it and half over its outermost band. That
   is the one focus idiom every control in this library wears (focus.Halo),
   and the indicator's coverage composites over the page outside the box and
   over the field's own edge and fill inside it. */
.input:focus-visible, .input.is-focus {
  outline: var(--focus-halo-width) solid var(--platform-keyboard-focus-indicator);
  outline-offset: calc(var(--focus-halo-width) / -2);
}
.input:disabled {
  color: var(--platform-disabled-control-text);
}
.input:disabled::placeholder {
  color: var(--platform-disabled-control-text);
}

/* The pop-up trigger (components/picker's form trigger, which
   components/input's dropdown forwards to): a BUTTON and not a field, and
   drawn as the platform's pop-up button rather than as the field beside it.

   MEASURED, save-dialog-{light,dark}.png, the "File Format:" pop-up: the
   control is 24 px tall, the same height the push button beside it draws and
   not the text field's 27; a run down its middle gives the push button's own
   fill from its first row to its last with no darker column at either end, so
   it draws NO edge and its fill meets the surface directly; its label's
   origin is 11 columns in from that fill's edge (the measured twelve less the
   one column of bearing the capture's own first letter carries), five deeper
   than the field's; and the mark's last column stands 9 clear of the trailing
   edge. Because there is no edge column, both insets are spent from the
   fill's own edge and nothing gives a hairline back.

   The pointer states are the pop-up's own. MEASURED,
   control-hover-{light,dark}.png - the Finder toolbar's view pop-up, the
   control drawing this very mark, under the pointer: the platform's hover
   overlay over the fill it stands on. Held, the press overlay over the same
   fill. Switched off, the fill is the push button's at --disabled-coverage
   over the surface the trigger stands on, and every foreground over it
   becomes the platform's disabled control text.

   The trailing room the label is kept clear of is the S3 gap plus the mark's
   own 8 px plus the measured 9 px of clearance - the gap is the trigger's,
   what stops a long value running into the mark, and not one of its two
   ends. */
.select {
  min-height: var(--density-control-height);
  border: none;
  border-radius: var(--radius-md);
  padding: 0 calc(var(--space-3) + 8px + 9px) 0 11px;  /* S3 gap + MarkWDp + PopupMarkTrailDp; PopupLeadDp */
  background-color: var(--platform-push-button-fill);
  background-clip: border-box;
  color: var(--platform-control-text);
}
.select:hover:not(:disabled), .select.is-hover {
  background-image: linear-gradient(var(--platform-hover-overlay), var(--platform-hover-overlay));
}
.select:active:not(:disabled), .select.is-active {
  background-image: linear-gradient(var(--platform-press-overlay), var(--platform-press-overlay));
}
/* Focus is the halo on the edgeless shape the trigger already draws, and
   adds nothing else: the insets and the height do not move. */
.select:focus-visible, .select.is-focus {
  border: none;
  padding: 0 calc(var(--space-3) + 8px + 9px) 0 11px;
  outline: var(--focus-halo-width) solid var(--platform-keyboard-focus-indicator);
  outline-offset: calc(var(--focus-halo-width) / -2);
}
.select:disabled {
  background-color: color-mix(in srgb, var(--platform-push-button-fill) var(--disabled-coverage), transparent);
  background-image: none;
  color: var(--platform-disabled-control-text);
}
.select-wrap { position: relative; display: block; }

/* The mark is the platform's pop-up mark: two chevrons stacked point to
   point, the upper pointing up and the lower down, saying the control holds
   one of several values and never which way its menu goes. A native select
   cannot carry a generated child, so the .select-wrap wrapper draws it.

   MEASURED, save-dialog-{light,dark}.png, both appearances agreeing to the
   pixel: the pair spans 8 columns and 11 rows - each chevron 5 rows with one
   clear row between them - and its last column stands 9 clear of the fill's
   trailing edge. The size is FIXED: the Finder toolbar draws the same 8 by 11
   in a control 36 px tall, so the mark does not scale with what it stands in.
   The arm's weight is 1.5 px perpendicular, fitted to the capture's coverage.

   It is a masked SVG rather than a border-built triangle because the platform
   draws two strokes and not a solid wedge: the mask carries the two chevrons
   at exactly the geometry control.DrawMark strokes, and what shows through is
   the platform's own name for a control's own marks, --platform-control-text
   (MEASURED: the pair's covered pixels read the control text's coverage over
   the trigger's fill to the byte; the secondary label's would land a hundred
   levels away). */
.select-wrap::after {
  content: "";
  position: absolute;
  right: 9px;  /* PopupMarkTrailDp */
  top: 50%;
  margin-top: -5.5px;  /* half MarkHDp */
  width: 8px;   /* MarkWDp */
  height: 11px; /* MarkHDp: two 5-row chevrons and the clear row between */
  background: var(--platform-control-text);
  /* The two chevrons as control.DrawMark strokes them: each arm a V from
     (0,5) to (4,0) to (8,5) closed back along the base, the closing side
     offset by the 1.5 px arm's own half-width resolved along the base
     (1.921 across, 2.401 down at the apex, which is (stroke/2) x sqrt(1+k2)/k
     for the centreline's k = 2h/w = 1.25). The lower chevron is the same
     figure inverted one clear row below it. */
  -webkit-mask: url("data:image/svg+xml,%3Csvg%20xmlns=%22http://www.w3.org/2000/svg%22%20width=%228%22%20height=%2211%22%20viewBox=%220%200%208%2011%22%3E%3Cpath%20d=%22M0%205L4%200L8%205L6.079%205L4%202.401L1.921%205ZM0%206L4%2011L8%206L6.079%206L4%208.599L1.921%206Z%22/%3E%3C/svg%3E") center / 8px 11px no-repeat;
  mask: url("data:image/svg+xml,%3Csvg%20xmlns=%22http://www.w3.org/2000/svg%22%20width=%228%22%20height=%2211%22%20viewBox=%220%200%208%2011%22%3E%3Cpath%20d=%22M0%205L4%200L8%205L6.079%205L4%202.401L1.921%205ZM0%206L4%2011L8%206L6.079%206L4%208.599L1.921%206Z%22/%3E%3C/svg%3E") center / 8px 11px no-repeat;
  pointer-events: none;
}
.select-wrap:has(.select:disabled)::after {
  background: var(--platform-disabled-control-text);
}

/* Checkbox (components/input checkbox.go): the 16 dp glyph the Save dialog
   measures, centred in the density's checkbox row - the square footprint the
   platform gives a pointer, 22 dp comfortable against the push button's 24
   and the list row's 20. The glyph does not follow density and the footprint
   does, so the slack around it is written as the margin that centres it;
   nothing else in this family moves with the density.

   The corner is the measured 5 dp: a circular fit to the per-row coverage of
   the switched-off boxes' corner ramp in save-dialog-{light,dark}.png answers
   5.04 and 5.34 against the same fit's habit of sitting a fifth over
   everywhere in that reference. The edge is the measured 1 dp of the
   platform's field hairline - the Save dialog's "Tags:" field is the sheet's
   one unfocused enabled control that draws an edge at all, and it draws a
   single pixel of it.

   Checked, the box is the platform's accent under a check mark in
   alternateSelectedControlText, because a fill says a colour was applied and
   only the mark says what that means: a column of fills carries completion in
   hue alone, which is the one channel a reader may not have.

   The mark is drawn, not encoded. Gio strokes the icon set's centre line -
   (4.5,12) -> (9,16.5) -> (19.5,6) on the set's 24-unit grid, a 2-unit
   DIAGONAL band, round caps and joins - and at the 16 px glyph one grid unit
   is 2/3 px, so the band is 1.333 px wide (+/-0.667 either side of the centre)
   and the arms run from (3,8) to (6,11) to (13,4). Each arm is one background
   layer: a linear-gradient banding its own box perpendicular to the arm,
   45deg for the short "\" arm and 135deg for the long "/" one. Every stop is
   written from 50% because each box is sized to its arm - the segment grown
   by half a band along its own axis, which makes the arm the box's diagonal
   and the box's corners the round caps' own tips:
     short arm: 3.943 px square at 2.529,7.529   (3,8)->(6,11)
     long arm:  7.943 px square at 5.529,3.529   (6,11)->(13,4)
   CSS has no line cap, so the caps come out cut square inside those tips
   rather than rounded - the same trade the icon set's own SVG files make when
   they draw their caps as an explicit contour, and a sub-pixel one at this
   size. background-origin is the border box so the grid is the 16 px glyph
   the Gio side scales on, not the 15 px inside the edge. */
.checkbox, .radio {
  box-sizing: border-box;
  appearance: none;
  flex: none;
  width: 16px;  /* checkboxBoxSize / radioCircleSize: 16 dp, measured */
  height: 16px;
  /* The slack the density's footprint holds around the glyph, which is what
     centres it in the row and what a pointer is given. */
  margin: calc((var(--density-checkbox-row-height) - 16px) / 2);
  border: 1px solid var(--platform-field-edge);  /* controlEdgeWidth: 1 dp, measured */
  background: var(--platform-text-background);
  background-clip: padding-box;
  cursor: pointer;
}
/* Focus is the same halo every other control wears, on the glyph's own
   outline, riding in the slack the footprint holds around it - so nothing
   about the control moves when it takes the keyboard. */
.checkbox:focus-visible, .checkbox.is-focus,
.radio:focus-visible, .radio.is-focus {
  outline: var(--focus-halo-width) solid var(--platform-keyboard-focus-indicator);
  outline-offset: calc(var(--focus-halo-width) / -2);
}
.checkbox {
  border-radius: 5px;  /* checkboxCornerRadius: measured */
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
/* Switched off, the glyph is ONE FILL AND NO EDGE at all. MEASURED,
   save-dialog-{light,dark}.png: both switched-off "Options:" checkboxes read
   the push button's own fill faded to --disabled-coverage over the sheet, to
   the byte light and one 255th over on dark green and blue, and neither box
   draws an edge column in either appearance - its rim is a one-pixel antialiased ramp from this fill
   to the sheet. The border is kept at its width in transparent so the glyph's
   drawn box does not move. Checked, the mark takes the colour the
   switched-off label beside it takes, the platform's tertiary label, no
   stored capture holding a switched-off checked box. */
.checkbox:disabled, .radio:disabled {
  cursor: default;
  border-color: transparent;
  background-color: color-mix(in srgb, var(--platform-push-button-fill) var(--disabled-coverage), transparent);
  background-clip: border-box;
}
.checkbox:checked:disabled, .checkbox.is-checked:disabled {
  background-image:
    linear-gradient(45deg, transparent calc(50% - 0.667px), var(--platform-tertiary-label) calc(50% - 0.667px), var(--platform-tertiary-label) calc(50% + 0.667px), transparent calc(50% + 0.667px)),
    linear-gradient(135deg, transparent calc(50% - 0.667px), var(--platform-tertiary-label) calc(50% - 0.667px), var(--platform-tertiary-label) calc(50% + 0.667px), transparent calc(50% + 0.667px));
}

/* Radio (components/input radio.go): the same 16 dp glyph as a circle.
   Selected is the platform's accent filling the whole disc with the dot in
   alternateSelectedControlText at its centre - one fill and one mark, exactly
   the Gio nested ellipses, and no gap ring.

   The dot is 5 dp across, MEASURED off System Settings' selected radio in
   both appearances: a least-squares circle fitted to the white dot's
   sub-pixel edges reads 5.00 px across inside a 16 px disc - five sixteenths
   of the glyph, not the half a fill drawn to the glyph's own ratio would
   give. */
.radio { border-radius: var(--radius-full); }
.radio:checked, .radio.is-checked {
  border-color: var(--platform-control-accent);
  background: radial-gradient(circle, var(--platform-alternate-selected-control-text) 2.5px, var(--platform-control-accent) 2.5px); /* radioDotSize: 5 dp, measured */
  background-clip: border-box;
}
.radio:checked:disabled, .radio.is-checked:disabled {
  background: radial-gradient(circle, var(--platform-tertiary-label) 2.5px, color-mix(in srgb, var(--platform-push-button-fill) var(--disabled-coverage), transparent) 2.5px);
  background-clip: border-box;
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
  letter-spacing: 0;
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
  letter-spacing: 0;
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
  letter-spacing: 0;
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
  letter-spacing: 0;
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
   inner content box, ControlHeight - 2*PaddingY — in the platform's
   secondary label. Pointer-only, never a Tab stop. */
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

/* An item row (drawItem): a symbol, a label and, at the trailing end, a count
   when the entry has one. The symbol stands in a 24 dp square set 17 in from
   the rail's leading edge, the label-large label starts at 48 in, and the
   count's trailing edge is 17 in from the rail's trailing edge; each part is
   vertically centred, one line, clipped rather than wrapped — which is also
   what hides the label and the count at the collapsed width. Selected wears
   the platform's sidebar pill, and on it the label and the count both take
   the foreground the platform pairs with that fill; nothing else moves, and
   nothing tints under the pointer.

   The row height, the pill's inset and its corner, the three columns and the
   heading's block are MEASURED off the organization's macOS reference
   (patterns/sidebar RowHeight, SelectionInset, SelectionRadius, SymbolBox,
   SymbolInset, LabelInset, CountInset, SectionHeight, SectionInset,
   SectionBaseline): Finder's and Voice Memos' selected sidebar rows span 32 px
   at 1x, the pill is inset 10 from each edge of the rail, and a circular fit
   to its corner reads 8; Voice Memos' folder mark is centred on 29 in from the
   panel's edge, its names start at 48 and its counts end 17 in from the
   trailing edge. */
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
  letter-spacing: 0;
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
  background: var(--platform-sidebar-selection);
}
.sidebar-item-icon {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;  /* SymbolBox */
  height: 100%;
  margin-left: 17px;  /* SymbolInset */
  margin-right: 7px;  /* LabelInset less SymbolInset and SymbolBox */
}
/* Collapsed a row is its symbol and nothing else: drawItem returns after the
   symbol, so the label, the count and the section heading are not drawn at
   all, and the symbol box is centred in the rail rather than set at the
   leading inset. */
.sidebar.collapsed .sidebar-item-icon {
  margin-left: calc((48px - 24px) / 2);  /* collapsedDp less SymbolBox, halved */
  margin-right: 0;
}
.sidebar.collapsed .sidebar-item-label,
.sidebar.collapsed .sidebar-item-count,
.sidebar.collapsed .sidebar-section {
  display: none;
}
/* A row with no symbol still starts its label at the same column, so the
   names of a list whose entries differ still line up. */
.sidebar-item-label {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}
.sidebar-item-label:first-child { margin-left: 48px; }  /* LabelInset */
.sidebar-item-count {
  flex: none;
  margin-left: var(--space-2);
  margin-right: 17px;  /* CountInset */
  color: var(--platform-sidebar-count);
}
.sidebar-item.selected .sidebar-item-count {
  color: var(--platform-alternate-selected-control-text);
}

/* A section's heading: a small label in the platform's secondary label,
   17 in from the rail's leading edge, its baseline 30 down a block of 42,
   and parted from the rows around it by that air alone — the platform draws
   no line there and neither does this.

   A sheet cannot place a baseline, so what is written is where the library's
   own placement puts the line box: 30 less the role's ascent inside its line
   box, which for the shipped face at 11 dp on a 16 dp line is 18. */
.sidebar-section {
  box-sizing: border-box;
  flex: none;
  display: flex;
  align-items: flex-start;
  height: 42px;  /* SectionHeight */
  padding: 18px 17px 0;  /* SectionBaseline less the role's ascent; SectionInset */
  font-family: var(--font-family);
  font-size: var(--font-label-small-size);
  line-height: var(--font-label-small-line-height);
  font-weight: var(--font-label-small-weight);
  letter-spacing: 0;
  color: var(--platform-secondary-label);
  user-select: none;
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
  letter-spacing: 0;
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

/* The keyboard halo, identical to every other control's: per-cell for the
   navbar, tabs and breadcrumb (each cell is its own Clickable focus tag); on
   the rail itself for the sidebar, whose single stop is the item list. */
.navbar-link:focus-visible, .navbar-link.is-focus,
.tab:focus-visible, .tab.is-focus,
.crumb:focus-visible, .crumb.is-focus,
.sidebar:focus-visible, .sidebar.is-focus {
  outline: var(--focus-halo-width) solid var(--platform-keyboard-focus-indicator);
  outline-offset: calc(var(--focus-halo-width) / -2);
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
  letter-spacing: 0;
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

/* Menu (components/picker menu.go, the surface its trigger opens): the
   platform's menu as the Level entry gives every floating surface - the
   window's own background under the measured floating shadow, which is the
   whole of what tells the menu's fill from the window's behind it, inside the
   platform's separator laid ON the plane's outline so the edge costs the
   plane no height.

   The plane's corner and the rows' pill are NOT MEASURED: no stored capture
   holds an open menu at all, so each is the corner and the inset of the one
   pill the reference does measure - the sidebar's selected row, inset 10 from
   its rail and cornered at 8 - and the capture that would settle a menu's own
   is on the reference's list.

   A row is the control height tall with the density's vertical padding, its
   mark box standing 10 in from the leading edge at the 16 dp size a mark
   beside a line of body text is drawn at, its label starting a text lead
   after that box, and 16 clear at the trailing end. The current row wears the
   pill and takes the foreground the platform pairs with that fill; nothing
   else moves.

   An open menu stands OVER its trigger with the current row on the trigger's
   own label rather than dropping below it, which is placement and the page's
   to arrange - the class carries the surface and the rows. */
.menu {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  min-width: 48px;
  overflow: hidden;
  border-radius: 8px;  /* planeRadius: the sidebar pill's own corner */
  background: var(--platform-window-background);
  color: var(--platform-label);
  font-family: var(--font-family);
  font-size: var(--font-body-large-size);
  line-height: var(--font-body-large-line-height);
  font-weight: var(--font-body-large-weight);
  letter-spacing: 0;
  --floating-shadow-step: color-mix(in srgb, var(--platform-floating-shadow) 12.5%, transparent);
  box-shadow:
    0 0 0 3px var(--floating-shadow-step),
    0 0 0 6px var(--floating-shadow-step),
    0 0 0 9px var(--floating-shadow-step),
    0 0 0 12px var(--floating-shadow-step),
    0 0 0 15px var(--floating-shadow-step),
    0 0 0 18px var(--floating-shadow-step),
    0 0 0 21px var(--floating-shadow-step),
    0 0 0 24px var(--floating-shadow-step),
    inset 0 0 0 1px var(--platform-separator);  /* edgeDp, laid on the outline */
}
.menu-item {
  box-sizing: border-box;
  position: relative;
  z-index: 0;
  display: flex;
  align-items: center;
  min-height: var(--density-control-height);
  padding: var(--density-padding-y) 16px var(--density-padding-y) 32px;  /* selectionInset + TextTrail; selectionInset + markBox + TextLead */
  white-space: nowrap;
  cursor: pointer;
  user-select: none;
}
/* The pill is a layer behind the row's own content rather than the row's
   fill, because it is inset from the plane while the mark column is not. */
.menu-item.selected { color: var(--platform-alternate-selected-control-text); }
.menu-item.selected::before {
  content: "";
  position: absolute;
  inset: 0 10px;  /* selectionInsetDp */
  z-index: -1;
  border-radius: 8px;  /* selectionRadiusDp */
  background: var(--platform-sidebar-selection);
}
/* The check beside the current item, drawn in the row's own foreground out of
   the same two gradient bands the checkbox's mark is drawn from, in a 16 dp
   box standing 10 in from the plane's leading edge. */
.menu-item.selected::after {
  content: "";
  position: absolute;
  left: 10px;  /* selectionInsetDp */
  top: 50%;
  margin-top: -8px;
  width: 16px;  /* markBoxDp */
  height: 16px;
  background-image:
    linear-gradient(45deg, transparent calc(50% - 0.667px), currentColor calc(50% - 0.667px), currentColor calc(50% + 0.667px), transparent calc(50% + 0.667px)),
    linear-gradient(135deg, transparent calc(50% - 0.667px), currentColor calc(50% - 0.667px), currentColor calc(50% + 0.667px), transparent calc(50% + 0.667px));
  background-repeat: no-repeat;
  background-position: 2.529px 7.529px, 5.529px 3.529px;
  background-size: 3.943px 3.943px, 7.943px 7.943px;
}

/* Pane (patterns/pane): the platform's sidebar is an inset rounded panel
   standing inside the window, not a flush column parted by a seam. MEASURED
   off the stored Voice Memos captures: the panel is inset 8 from the window's
   leading, top and bottom edges, cornered at 18 - concentric with the
   window's own 26 one margin out - wears a 1 px rim just inside its edge in
   the platform's measured value for it, and casts a shadow onto what lies
   beside it from a rectangle sunk 9 below the panel, carrying 24 out. The
   rim is an inset ring so the panel's box does not grow, and the shadow is
   the same eight-step approximation of effects/depth's linear falloff the
   floating surfaces take. The content beside it is the window's own plane:
   the rim and the shadow are the boundary, and no seam is drawn. */
.pane {
  box-sizing: border-box;
  overflow: hidden;
  margin: 8px;  /* MarginDp */
  border-radius: 18px;  /* RadiusDp: the window's 26 less one margin */
  background: var(--platform-sidebar-material);
  --pane-shadow-step: color-mix(in srgb, var(--platform-pane-shadow) 12.5%, transparent);
  box-shadow:
    0 9px 0 3px var(--pane-shadow-step),
    0 9px 0 6px var(--pane-shadow-step),
    0 9px 0 9px var(--pane-shadow-step),
    0 9px 0 12px var(--pane-shadow-step),
    0 9px 0 15px var(--pane-shadow-step),
    0 9px 0 18px var(--pane-shadow-step),
    0 9px 0 21px var(--pane-shadow-step),
    0 9px 0 24px var(--pane-shadow-step),  /* ShadowSinkDp, ShadowReachDp */
    inset 0 0 0 1px var(--platform-pane-rim);  /* RimDp */
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
  letter-spacing: 0;
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
  letter-spacing: 0;
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
