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

// fnum formats a float32 with no trailing zeros: 16 → "16", 0.15 → "0.15".
func fnum(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

// px formats a dp value as a CSS px length. Device-independent pixels map
// 1:1 onto CSS px, both being density-abstract logical pixels.
func px(v float32) string {
	return fnum(v) + "px"
}

// rampRoles orders ADR-007's colour roles under their CSS names.
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
	{"accent", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.Primary }},
	{"on-accent", func(t tokens.ColorTokens) stdcolor.NRGBA { return t.OnPrimary }},
	// The solid-fill state walk (ADR-007 / D2.3): hover one rung from the pin
	// toward the ramp's 900 end, pressed two — SolidStateColor, the exact
	// resolution components/button's filled register draws. They are emitted
	// as first-class tokens because a walked pin is off-ramp: no var()
	// arithmetic over the ramp steps could reproduce it, and ADR-007's whole
	// point is that states are real, addressable colours a sheet can emit.
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

// elevationLevels orders the elevation levels under their level numbers;
// each level's surface step and shadow dp are read off the snapshot's
// ElevationScale through its accessors.
var elevationLevels = []struct {
	name  string
	level tokens.ElevationLevel
}{
	{"0", tokens.Level0},
	{"1", tokens.Level1},
	{"2", tokens.Level2},
	{"3", tokens.Level3},
}

// densityMetrics orders the per-setting density metrics under their CSS
// names. The WCAG pointer-target floor is not here: it is not a per-setting
// metric — see densityVars.
var densityMetrics = []struct {
	name string
	pick func(tokens.Density) float32
}{
	{"control-height", func(d tokens.Density) float32 { return d.ControlHeight }},
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

// surfaceVarRef renders an elevation level's surface fill as a reference
// into the colour families — var(--color-bg) for the step-0 sentinel (the
// Background pin), var(--color-neutral-<step>) otherwise. Emitting a
// reference rather than a resolved hex keeps --elevation-* mode-invariant:
// the .dark block overrides the colour variables and every elevation
// surface flips with them, which is exactly the tonal model (an elevation
// level IS a neutral-ramp step).
func surfaceVarRef(step int) string {
	if step == 0 {
		return "var(--color-bg)"
	}
	return fmt.Sprintf("var(--color-neutral-%d)", step)
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
	return vars
}

// scaleVars renders the mode-invariant families: fonts, density
// (comfortable — the :root setting), spacing, radius, the tonal elevation
// surfaces (the default cue), the dp shadows (the opt-in cue for floating
// transients, per E2.2), and the motion set.
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
		vars = append(vars, cssVar{"--elevation-" + level.name, surfaceVarRef(s.Elevation.SurfaceStep(level.level))})
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
	// The interaction-state base the class layer builds on (G1.2). The focus
	// ring is Neutral step 500 (tokens.FocusRing) — a var() reference, like
	// the elevation surfaces, so .dark flips it with the ramp — at the 2 dp
	// stroke width components/button pins; the disabled fraction is
	// tokens.DisabledOpacity as a color-mix() percentage, because disabled is
	// an opacity, not a ramp step.
	vars = append(vars,
		cssVar{"--color-focus-ring", "var(--color-neutral-500)"},
		cssVar{"--focus-ring-width", px(focusRingWidthDp)},
		cssVar{"--state-disabled-opacity", fnum(tokens.DisabledOpacity*100) + "%"},
	)
	return vars
}

// focusRingWidthDp is the focus ring's stroke width — the 2 dp
// components/button draws (drawButton's gtx.Dp(2) stroke), identical in
// every emphasis register because keyboard visibility is not a loudness
// property.
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
func block(b *strings.Builder, selector string, vars []cssVar) {
	b.WriteString(selector)
	b.WriteString(" {\n")
	for _, v := range vars {
		fmt.Fprintf(b, "  %s: %s;\n", v.name, v.value)
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

// componentClasses is the class layer (G1.2, extended by G2.1–G2.3): the
// component vocabulary the design surface composes screens from, defined
// over the token variables above — not one literal colour, so a re-branded
// sheet re-brands the components with it, and .dark/.compact flip them like
// everything else. The only literal lengths are the same component
// constants the Gio side hardcodes as dp literals rather than tokens
// (checkbox/radio's 20 dp glyph and 10 dp dot, the dropdown's 16 dp
// chevron, the 1/2 dp input borders); each is commented at its source.
//
// It mirrors components/button and components/input, the sources of truth:
// .btn is the filled register by default, .tonal and .ghost the G0A.1
// emphasis modifiers, and every state resolves as ADR-007's ramp walks from
// exactly the rungs buttonColors picks (button.go: tonalGround 200 /
// tonalText 900, ghostGround 200 / ghostText 700 / ghostTextOnWash 900; the
// filled fill walks via SolidStateColor into --color-accent-hover/
// -pressed). Because each register's blocks override every state it treats,
// later register blocks never bleed a state from an earlier one; :disabled
// resolutions are per-register for the same reason. Selected resolves as
// tokens.StateColor resolves StateSelected — the two-step walk pressed
// takes. The form controls resolve as components/input does: Surface ground
// under body text, neutral 500 strong border, neutral 700 placeholder and
// glyph, focus promoting the border to the accent pin, disabled fading each
// colour to the disabled fraction of its alpha.
//
// Every pointer/keyboard state rule also carries a forcing twin class
// (.is-hover, .is-active, .is-focus, .is-checked) grouped into the same
// rule. A static component page cannot hover itself, and duplicating the
// declarations in the page would fork the resolution; a grouped selector
// emitted by this generator shares the exact declarations with the live
// pseudo-class, so a forced specimen provably renders as the live state.
// Disabled needs no twin: the pages force it with the native attribute.
const componentClasses = `/* ---- Component classes (G1.2, forms and tags G2.1, card and table G2.2,
   navigation G2.3) ----
   The class vocabulary, built only on the tokens above. .btn mirrors
   components/button: filled by default, .tonal and .ghost the emphasis
   modifiers, states resolved as ADR-007 ramp walks from the same rungs the
   Gio side draws. .input/.select/.checkbox/.radio mirror components/input,
   .tag the chip patterns/pricing and patterns/hero draw, .card the
   patterns/card surface, .table the patterns/table grid, and the navigation
   family — .navbar, .tabs, .sidebar, .crumbs — the four patterns of the same
   names. The focus ring is
   identical in every register: keyboard visibility is not an emphasis
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
   end (hover one rung, pressed and selected two) — the walked stops are
   tokens, not mixes. */
.btn:hover, .btn.is-hover { background: var(--color-accent-hover); }
.btn.selected { background: var(--color-accent-pressed); }
.btn:active, .btn.is-active { background: var(--color-accent-pressed); }

/* Keyboard focus keeps the resting fill and adds the ring — a stroke
   centred on the control's edge, as the Gio side draws it. */
.btn:focus-visible, .btn.is-focus,
.checkbox:focus-visible, .checkbox.is-focus,
.radio:focus-visible, .radio.is-focus {
  outline: var(--focus-ring-width) solid var(--color-focus-ring);
  outline-offset: calc(var(--focus-ring-width) / -2);
}

/* Disabled is an opacity, not a ramp step: each colour keeps its hue and
   fades to the disabled fraction of its alpha. */
.btn:disabled {
  cursor: default;
  background: color-mix(in srgb, var(--color-accent) var(--state-disabled-opacity), transparent);
  color: color-mix(in srgb, var(--color-on-accent) var(--state-disabled-opacity), transparent);
}

/* Tonal: a tinted fill off the role's own ramp — ground 200 under the
   ramp's 900 text; hover walks one step, pressed and selected two. */
.btn.tonal {
  background: var(--color-primary-200);
  color: var(--color-primary-900);
}
.btn.tonal:hover, .btn.tonal.is-hover { background: var(--color-primary-300); }
.btn.tonal.selected { background: var(--color-primary-400); }
.btn.tonal:active, .btn.tonal.is-active { background: var(--color-primary-400); }
.btn.tonal:disabled {
  background: color-mix(in srgb, var(--color-primary-200) var(--state-disabled-opacity), transparent);
  color: color-mix(in srgb, var(--color-primary-900) var(--state-disabled-opacity), transparent);
}

/* Ghost: no ground at rest — the neutral ramp's low-contrast text over
   whatever surface it sits on; under the pointer it performs that
   surface's own hover (300) and press (400), the text walking to 900 with
   the ground. No selected treatment: a ghost stays quiet. */
.btn.ghost {
  background: transparent;
  color: var(--color-neutral-700);
}
.btn.ghost:hover, .btn.ghost.is-hover {
  background: var(--color-neutral-300);
  color: var(--color-neutral-900);
}
.btn.ghost:active, .btn.ghost.is-active {
  background: var(--color-neutral-400);
  color: var(--color-neutral-900);
}
.btn.ghost:disabled {
  background: transparent;
  color: color-mix(in srgb, var(--color-neutral-700) var(--state-disabled-opacity), transparent);
}

/* Icon-only form (components/button drawIconButton): a square the
   density's control height on a side, the glyph inset by the density's
   vertical padding — content box ControlHeight − 2·PaddingY, icon.Size's
   rule (20 dp comfortable, 16 dp compact). Emphasis reaches the colours
   and stops there: the square never shrinks. The glyph inherits the
   register's text colour via currentColor. */
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

/* ---- Tag (G2.1) ----
   The chip the patterns draw: a Full-radius pill, S1/S2 padding, sized to
   its label-small text. Filled by default — patterns/pricing's "Popular"
   chip, the accent pin under its on-colour; .tonal is patterns/hero's
   eyebrow — the primary 200 tinted fill under the accent pin. Both call
   sites request SemiBold, which the pinned shaper resolves to the Medium
   face (the nearest registered weight), so the sheet says the label role's
   own weight rather than asking the browser to synthesize a 600. A tag is
   a label, not a control: it has no interaction states. */
.tag {
  box-sizing: border-box;
  display: inline-flex;
  align-items: center;
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-full);
  white-space: nowrap;
  font-family: var(--font-family);
  font-size: var(--font-label-small-size);
  line-height: var(--font-label-small-line-height);
  font-weight: var(--font-label-small-weight);
  letter-spacing: var(--font-label-small-tracking);
  background: var(--color-accent);
  color: var(--color-on-accent);
}
.tag.tonal {
  background: var(--color-primary-200);
  color: var(--color-accent);
}

/* ---- Form controls (G2.1) ----
   Native elements wearing components/input's resolution: Surface ground
   under body-large text, neutral 500 strong border, neutral 700
   placeholder, focus promoting the border to the accent pin, disabled
   fading every colour to the disabled fraction of its alpha. */

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
  border: 1px solid var(--color-neutral-500);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  color: var(--color-text);
  font-family: var(--font-family);
  font-size: var(--font-body-large-size);
  line-height: var(--font-body-large-line-height);
  font-weight: var(--font-body-large-weight);
  letter-spacing: var(--font-body-large-tracking);
}
.input::placeholder { color: var(--color-neutral-700); opacity: 1; }

/* Focus promotes the border to the accent pin and doubles it to the 2 dp
   the Gio side draws — the second pixel as an inset shadow, so the field's
   outer geometry and text position do not move (Gio thickens the border
   inward the same way). */
.input:focus-visible, .input.is-focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: inset 0 0 0 1px var(--color-accent);
}
.input:disabled {
  background: color-mix(in srgb, var(--color-surface) var(--state-disabled-opacity), transparent);
  color: color-mix(in srgb, var(--color-text) var(--state-disabled-opacity), transparent);
  border-color: color-mix(in srgb, var(--color-neutral-500) var(--state-disabled-opacity), transparent);
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
   — a component constant, not a token; it does not follow density) with a
   2 dp neutral 500 border over Surface when unchecked, a solid accent fill
   when checked — no checkmark glyph: the Gio side draws none. The focus
   ring is the shared rule above. */
.checkbox, .radio {
  box-sizing: border-box;
  appearance: none;
  flex: none;
  width: 20px;  /* checkboxBoxSize / radioCircleSize: 20 dp */
  height: 20px;
  margin: 0;
  border: 2px solid var(--color-neutral-500);
  background: var(--color-surface);
  cursor: pointer;
}
.checkbox { border-radius: var(--radius-sm); }
.checkbox:checked, .checkbox.is-checked {
  border-color: var(--color-accent);
  background: var(--color-accent);
}
.checkbox:disabled {
  cursor: default;
  border-color: color-mix(in srgb, var(--color-neutral-500) var(--state-disabled-opacity), transparent);
  background: color-mix(in srgb, var(--color-surface) var(--state-disabled-opacity), transparent);
}
.checkbox:checked:disabled, .checkbox.is-checked:disabled {
  border-color: color-mix(in srgb, var(--color-accent) var(--state-disabled-opacity), transparent);
  background: color-mix(in srgb, var(--color-accent) var(--state-disabled-opacity), transparent);
}

/* Radio (components/input radio.go): the same 20 dp glyph as a circle;
   selected keeps the Surface gap ring and fills a 10 dp accent dot
   (radioDotSize) — outer accent ring, surface, dot, exactly the Gio
   nested fills. */
.radio { border-radius: var(--radius-full); }
.radio:checked, .radio.is-checked {
  border-color: var(--color-accent);
  background: radial-gradient(circle, var(--color-accent) 5px, var(--color-surface) 5px); /* 10 dp dot */
}
.radio:disabled {
  cursor: default;
  border-color: color-mix(in srgb, var(--color-neutral-500) var(--state-disabled-opacity), transparent);
  background: color-mix(in srgb, var(--color-surface) var(--state-disabled-opacity), transparent);
}
.radio:checked:disabled, .radio.is-checked:disabled {
  border-color: color-mix(in srgb, var(--color-accent) var(--state-disabled-opacity), transparent);
  background: radial-gradient(circle, color-mix(in srgb, var(--color-accent) var(--state-disabled-opacity), transparent) 5px, color-mix(in srgb, var(--color-surface) var(--state-disabled-opacity), transparent) 5px);
}

/* ---- Card (G2.2) ----
   patterns/card: a rounded surface raised in place by tonal step alone
   (ADR-007, E2.2's verdict) — no cast shadow in either variant, because a
   card is raised, not floating; the dp shadows stay reserved for surfaces
   that can leave (menus, dialogs, toasts). The default outlined card fills
   at level 1 (--elevation-1, the neutral-200 storey) under a 1 dp neutral
   500 strong stroke; .elevated trades the stroke for one storey deeper
   (--elevation-2, neutral 300). Radius Lg, an S4 inset, S3 gaps between the
   slots — exactly drawCard's rad.Lg / sp.S4 / sp.S3. The Gio stroke is
   centred on the card's edge while the CSS border lies inside it, so the
   outlined padding gives back the border's 1px and the slots land where the
   Gio inset puts them. The card styles no slot text of its own: the Gio
   card draws no text, so slot typography belongs to the content. */
.card {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: calc(var(--space-4) - 1px);
  border: 1px solid var(--color-neutral-500);
  border-radius: var(--radius-lg);
  background: var(--elevation-1);
  color: var(--color-text);
}
.card.elevated {
  padding: var(--space-4);
  border: none;
  background: var(--elevation-2);
}

/* ---- Table (G2.2) ----
   patterns/table: the whole grid grounds on the Surface pin (drawTable);
   the header band fills neutral 300 under neutral 700 label-large text
   (drawHeaderRow / drawHeaderCell). Header and body rows are each exactly
   one control height tall — the E1.4 row rule (list.RowHeight), so .compact
   re-pitches the whole grid — and every row closes with a 1 dp Divider rule
   drawn inside its height. Cells inset horizontally by S3 (cellPadDp,
   12 dp — static, the E1.3 input rule); body text is body-medium at the
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
  background: var(--color-surface);
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
  background: var(--color-neutral-300);
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

/* ---- Navigation (G2.3) ----
   The four navigation patterns. All four rest their interactive cells on
   the Surface ground (neutral 200), so the pointer states are that ground's
   own ADR-007 walk — hover one rung to neutral 300, exactly the wash the
   ghost register performs on the same storey. The Gio side draws no hover
   (a native window has the pointer; a static page shows the resolution),
   but the rungs are the tokens' StateColor walk from ground 200, not a new
   mix. Selection is what the Gio side does draw: the 2 dp Primary underline
   (navbar.go / tabs.go underlineDp) or the sidebar's StateColor(RolePrimary,
   200, StateSelected) two-step walk to primary 400. */

/* Navbar (patterns/navbar navbar.go): a horizontal Surface bar —
   drawNavbar fills Surface, insets PaddingY vertically and S4 horizontally,
   and patterns/shell pins the bar to ControlHeight + 2*PaddingY (52 dp
   comfortable, 40 compact). Slots run brand, centred links, actions; the
   links row centres in the space brand and actions leave over (two equal
   flexed spacers), which margin-inline auto reproduces exactly, including
   the documented off-centre approximation when the end slots differ. */
.navbar {
  box-sizing: border-box;
  display: flex;
  align-items: center;
  min-height: calc(var(--density-control-height) + 2 * var(--density-padding-y));
  padding: var(--density-padding-y) var(--space-4);
  background: var(--color-surface);
  color: var(--color-text);
}
.navbar-links {
  display: flex;
  align-items: center;
  gap: var(--space-2);  /* linksRow's HSpacer(sp.S2) between link cells */
  margin-inline: auto;  /* the two equal flexed spacers */
}

/* A link cell (navbar.go linkWidget): label-large at the Text pin inside
   (S3, PaddingY) padding, with a 2 dp underline slot along the bottom edge
   that the Active link fills with the Primary pin — the underline runs the
   full cell width, padding included, exactly the image.Rect the Gio side
   fills. Hover is the Surface ground's one-rung walk. */
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

/* Tabs (patterns/tabs tabs.go): the strip is a Surface row of tab cells
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
  background: var(--color-surface);
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
   toggle row and every item row are exactly ControlHeight tall (the E1.4
   row rule, list.RowHeight), so .compact re-pitches the rail. The whole
   rail is one keyboard stop — the focus ring belongs to the rail, and the
   Arrow keys move .selected — so the ring rule below targets .sidebar
   itself, not the rows. */
.sidebar {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  width: 192px;  /* expandedDp */
  background: var(--color-surface);
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
   Surface ground to primary 400. */
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
   neutral 700 and, being links, hover to neutral 900 — the ghost register's
   text walk on the same ground. Colour follows position (labelColor), so
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
`
