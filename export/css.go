package export

import (
	"fmt"
	stdcolor "image/color"
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
	// The focus ring's 2 dp stroke width, mode-invariant, which is why it
	// is here and the ring's COLOUR is not: the ring wears
	// --platform-keyboard-focus-indicator, which flips with the appearance.
	vars = append(vars, cssVar{"--focus-ring-width", px(focusRingWidthDp)})
	return vars
}

// focusRingWidthDp is the focus ring's stroke width — the 2 dp
// components/button draws (drawButton's gtx.Dp(2) stroke), identical in
// every emphasis because keyboard visibility is not a matter of
// prominence.
const focusRingWidthDp = 2

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
	block(&b, ":root", append(platformVars(s.PlatformLight), scaleVars(s)...))
	b.WriteString("\n")
	block(&b, ".dark", platformVars(s.PlatformDark))
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
// all. Focus is --platform-keyboard-focus-indicator at --focus-ring-width,
// the same ring in every variant. Disabled falls a fill back to the push
// button's own and takes every foreground to
// --platform-disabled-control-text.
//
// This sheet emits no hover rule and no fade on the disabled fill, and the
// components draw both: the library lays --platform-hover-overlay over the
// fill a control carries and fades a switched-off control toward the
// surface it stands on. The sheet has not been brought in step, and neither
// has the trigger this sheet still draws with a hairline and a solid
// triangle; a reader comparing the two is reading a sheet behind the
// components, not two answers deliberately kept apart.
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

/* No hover rule in any variant, and the absence is a sheet behind the
   components rather than a measurement: no stored capture holds a push
   button under the pointer, so nothing measures one as exempt, and the
   library lays the platform's hover overlay on every variant.
   Held, the platform lays its press overlay over the fill the
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

/* Disabled falls a variant that carries a fill back to the push button's
   fill inside the separator hairline and takes every foreground to the
   platform's disabled control text. The library fades that fill and that
   hairline toward the surface the control stands on as well, at the
   coverage the Save dialog's switched-off checkbox measures; this sheet
   does not, and is behind it. The padding gives back the hairline's
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
   and the fill stays; checked, the accent drains — the fill is the
   platform's disabled control text over the surface, with the mark in the
   platform's own control text so it still reads. */
.checkbox:disabled, .radio:disabled {
  cursor: default;
  border-color: var(--platform-disabled-control-text);
}
.checkbox:checked:disabled, .checkbox.is-checked:disabled {
  border-color: var(--platform-disabled-control-text);
  background-color: var(--platform-disabled-control-text);
  background-image:
    linear-gradient(45deg, transparent calc(50% - 0.667px), var(--platform-control-text) calc(50% - 0.667px), var(--platform-control-text) calc(50% + 0.667px), transparent calc(50% + 0.667px)),
    linear-gradient(135deg, transparent calc(50% - 0.667px), var(--platform-control-text) calc(50% - 0.667px), var(--platform-control-text) calc(50% + 0.667px), transparent calc(50% + 0.667px));
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
  border-color: var(--platform-disabled-control-text);
  background: radial-gradient(circle, var(--platform-control-text) 4px, var(--platform-disabled-control-text) 4px);
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
  background: var(--platform-sidebar-selection);
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
