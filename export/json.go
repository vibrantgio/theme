package export

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/vibrantgio/theme/tokens"
)

// Parameters is theme.json's shape: what reproduces the theme. The theme
// colour rebuilds the platform's accent rows exactly — the round-trip test
// asserts it — and the platform's whole set, the scales, fonts, radius,
// density settings, shadow depths and motion set are recorded alongside so
// a reader (or a prototype) need not run the generator to know them.
type Parameters struct {
	// ThemeColor is the theme colour as lowercase #rrggbb: the colour the
	// platform's accent rows were rebuilt for, and the one colour a reader
	// needs to reproduce the set from the catalogue.
	//
	// The key is "themeColour", which is the key theme/brand's own file
	// carries, so an exported theme.json loads as a kept brand without
	// translation.
	ThemeColor string `json:"themeColour"`

	// Platform is the platform's colour set per appearance, keyed by the
	// same name the sheet's --platform-* custom properties carry.
	Platform ModePlatform `json:"platform"`

	// Fonts names the heading, body and mono faces. Heading and body are
	// Roboto until a heading face exists; mono is the code style's face.
	Fonts Fonts `json:"fonts"`

	// Radius is the base radius in dp — tokens.RadiusScale.Base, the sheet's
	// --radius-base.
	Radius float64 `json:"radius"`

	// Density records the theme's active setting by name plus both
	// published settings' metrics.
	Density DensityParams `json:"density"`

	// Elevation records the shadow depth each of the six levels casts, from
	// the backdrop up — the opt-in cue for floating transients. A level's
	// fill is the platform's own, under Platform.
	Elevation ElevationParams `json:"elevation"`

	// Motion records the captured MotionScale in full — duration stops,
	// easing beziers and spring presets — so the file reproduces it without
	// running the generator. Springs are Go-side physics with no CSS
	// counterpart; this is their only serialisation.
	Motion MotionParams `json:"motion"`
}

// DensityParams records the density model: the active setting's name
// ("comfortable" or "compact") and both settings' metrics. A control's
// pointer target is not a number here — it is the control's own height,
// which the metrics already carry.
type DensityParams struct {
	Setting     string         `json:"setting"`
	Comfortable DensityMetrics `json:"comfortable"`
	Compact     DensityMetrics `json:"compact"`
}

// DensityMetrics is one density setting's per-setting metrics in dp.
type DensityMetrics struct {
	ControlHeight float64 `json:"controlHeight"`
	ChipHeight    float64 `json:"chipHeight"`
	// ToolbarControlHeight is the height of a bordered control standing in a
	// toolbar band, which the platform draws taller than the control height
	// above — 36 against 24, measured.
	ToolbarControlHeight float64 `json:"toolbarControlHeight"`
	PaddingX             float64 `json:"paddingX"`
	PaddingY             float64 `json:"paddingY"`
}

// ElevationParams records the shadow depth of each of the six levels, from
// the backdrop up. A level's FILL is the platform's own name for what that
// region is, under Parameters.Platform; the shadow is the separate, opt-in
// cue a floating surface carries.
type ElevationParams struct {
	ShadowDp [6]float64 `json:"shadowDp"`
}

// MotionParams records the motion set.
type MotionParams struct {
	Durations DurationParams `json:"durations"`
	Easings   EasingParams   `json:"easings"`
	Springs   SpringParams   `json:"springs"`
}

// DurationParams carries the five duration stops in milliseconds.
type DurationParams struct {
	XFast  float64 `json:"xFast"`
	Fast   float64 `json:"fast"`
	Normal float64 `json:"normal"`
	Slow   float64 `json:"slow"`
	XSlow  float64 `json:"xSlow"`
}

// EasingParams carries each easing preset as its cubic-bezier control
// points [x1, y1, x2, y2] — the same four numbers the sheet's --ease-*
// variables carry inside cubic-bezier().
type EasingParams struct {
	Standard             [4]float64 `json:"standard"`
	StandardAccelerate   [4]float64 `json:"standardAccelerate"`
	StandardDecelerate   [4]float64 `json:"standardDecelerate"`
	Emphasized           [4]float64 `json:"emphasized"`
	EmphasizedAccelerate [4]float64 `json:"emphasizedAccelerate"`
	EmphasizedDecelerate [4]float64 `json:"emphasizedDecelerate"`
}

// SpringParams carries the three spring presets.
type SpringParams struct {
	Default SpringParam `json:"default"`
	Snappy  SpringParam `json:"snappy"`
	Gentle  SpringParam `json:"gentle"`
}

// SpringParam is one damped-oscillator preset. Damping is recorded at the
// shortest decimal that reproduces the float32 exactly (the critical
// presets are 2·√(k·m), an irrational number), so the file reproduces the
// Go value bit-for-bit — see f64.
type SpringParam struct {
	Mass      float64 `json:"mass"`
	Stiffness float64 `json:"stiffness"`
	Damping   float64 `json:"damping"`
}

// f64 widens a float32 token value for JSON without dragging float64
// conversion noise into the file: it goes through the shortest decimal
// representation of the float32 (0.2, not 0.20000000298023224), which
// converts back to the identical float32 — the round-trip test asserts
// exactly that — while staying readable.
func f64(v float32) float64 {
	f, err := strconv.ParseFloat(strconv.FormatFloat(float64(v), 'g', -1, 32), 64)
	if err != nil {
		panic("export: f64: " + err.Error()) // unreachable: FormatFloat output always parses
	}
	return f
}

// ModePlatform carries the platform's colour set for both appearances,
// each name mapped to its value as lowercase #rrggbb, or #rrggbbaa where
// the platform's answer is a colour at a coverage.
type ModePlatform struct {
	Light map[string]string `json:"light"`
	Dark  map[string]string `json:"dark"`
}

// Fonts names the typefaces. Mono is the code style's face — the sheet's
// --font-family-code.
type Fonts struct {
	Heading string `json:"heading"`
	Body    string `json:"body"`
	Mono    string `json:"mono"`
}

// densityMetricsOf reads one setting's metrics.
func densityMetricsOf(d tokens.Density) DensityMetrics {
	return DensityMetrics{
		ControlHeight:        f64(d.ControlHeight),
		ChipHeight:           f64(d.ChipHeight()),
		ToolbarControlHeight: f64(d.ToolbarControlHeight),
		PaddingX:             f64(d.PaddingX),
		PaddingY:             f64(d.PaddingY),
	}
}

// bezier4 flattens a Bezier into its cubic-bezier() control points.
func bezier4(bz tokens.Bezier) [4]float64 {
	return [4]float64{f64(bz.P1[0]), f64(bz.P1[1]), f64(bz.P2[0]), f64(bz.P2[1])}
}

// springOf flattens a Spring preset.
func springOf(sp tokens.Spring) SpringParam {
	return SpringParam{Mass: f64(sp.Mass), Stiffness: f64(sp.Stiffness), Damping: f64(sp.Damping)}
}

// durMs is a duration in milliseconds.
func durMs(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

// parameters assembles the Parameters for a snapshot.
func parameters(s Snapshot) Parameters {
	setting, _ := densitySetting(s.Density) // Capture already validated it
	var elev ElevationParams
	for i, level := range shadowLevels {
		elev.ShadowDp[i] = f64(level.dp(s.Elevation))
	}
	return Parameters{
		ThemeColor: hexRGB(s.PlatformLight.ControlAccent),
		Platform:   ModePlatform{Light: platformMap(s.PlatformLight), Dark: platformMap(s.PlatformDark)},
		Fonts:      Fonts{Heading: s.Typography.HeadlineLarge.Typeface, Body: s.Typography.BodyLarge.Typeface, Mono: s.Typography.Code.Typeface},
		Radius:     float64(s.Radius.Base),
		Density: DensityParams{
			Setting:     setting,
			Comfortable: densityMetricsOf(tokens.Comfortable),
			Compact:     densityMetricsOf(tokens.Compact),
		},
		Elevation: elev,
		Motion: MotionParams{
			Durations: DurationParams{
				XFast:  durMs(s.Motion.DurXFast),
				Fast:   durMs(s.Motion.DurFast),
				Normal: durMs(s.Motion.DurNormal),
				Slow:   durMs(s.Motion.DurSlow),
				XSlow:  durMs(s.Motion.DurXSlow),
			},
			Easings: EasingParams{
				Standard:             bezier4(s.Motion.EaseStandard),
				StandardAccelerate:   bezier4(s.Motion.EaseStandardAccelerate),
				StandardDecelerate:   bezier4(s.Motion.EaseStandardDecelerate),
				Emphasized:           bezier4(s.Motion.EaseEmphasized),
				EmphasizedAccelerate: bezier4(s.Motion.EaseEmphasizedAccelerate),
				EmphasizedDecelerate: bezier4(s.Motion.EaseEmphasizedDecelerate),
			},
			Springs: SpringParams{
				Default: springOf(s.Motion.SpringDefault),
				Snappy:  springOf(s.Motion.SpringSnappy),
				Gentle:  springOf(s.Motion.SpringGentle),
			},
		},
	}
}

// platformMap renders one appearance's set as name → value, under the same
// names the sheet's --platform-* custom properties carry.
func platformMap(p tokens.PlatformColors) map[string]string {
	m := make(map[string]string, len(platformNames))
	for _, n := range platformNames {
		m[n.name] = hexRGBA(n.pick(p))
	}
	return m
}

// themeJSON renders theme.json, indented and newline-terminated.
func themeJSON(s Snapshot) ([]byte, error) {
	js, err := json.MarshalIndent(parameters(s), "", "  ")
	if err != nil {
		return nil, err
	}
	return append(js, '\n'), nil
}
