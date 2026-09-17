package export

import (
	"encoding/json"
	"fmt"
	stdcolor "image/color"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/reactivego/rx"
	"github.com/vibrantgio/theme/theme"
	"github.com/vibrantgio/theme/tokens"
)

// parseSheet is a tolerant reader of the emitted styles.css: it returns the
// custom properties per selector block, ignoring anything that is not a
// block opener, a block closer or a --declaration.
func parseSheet(t *testing.T, src string) map[string]map[string]string {
	t.Helper()
	blocks := map[string]map[string]string{}
	var cur map[string]string
	for _, line := range strings.Split(src, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasSuffix(line, "{"):
			sel := strings.TrimSpace(strings.TrimSuffix(line, "{"))
			if strings.HasPrefix(sel, "@") {
				// At-rules (@font-face) carry no custom properties and may
				// legitimately repeat — one block per face. Read through
				// them without recording a selector block.
				cur = map[string]string{}
				continue
			}
			if _, dup := blocks[sel]; dup {
				t.Fatalf("styles.css: duplicate block %q", sel)
			}
			cur = map[string]string{}
			blocks[sel] = cur
		case line == "}":
			cur = nil
		case strings.HasPrefix(line, "--"):
			if cur == nil {
				t.Fatalf("styles.css: declaration outside a block: %q", line)
			}
			// Declarations may carry a trailing kind annotation for the
			// Claude Design pane's token classifier ("; /* @kind other */");
			// strip it so the value checks see the bare declaration.
			if i := strings.Index(line, "/*"); i >= 0 && strings.HasSuffix(line, "*/") {
				line = strings.TrimSpace(line[:i])
			}
			name, val, ok := strings.Cut(line, ":")
			if !ok || !strings.HasSuffix(val, ";") {
				t.Fatalf("styles.css: malformed declaration: %q", line)
			}
			name = strings.TrimSpace(name)
			if _, dup := cur[name]; dup {
				t.Fatalf("styles.css: duplicate declaration %q", name)
			}
			cur[name] = strings.TrimSpace(strings.TrimSuffix(val, ";"))
		}
	}
	return blocks
}

// wantHex formats a colour the way the sheet must: lowercase #rrggbb.
// Deliberately written out rather than shared with the implementation, so
// the test and the serialiser cannot drift together.
func wantHex(c stdcolor.NRGBA) string {
	const digits = "0123456789abcdef"
	return string([]byte{'#',
		digits[c.R>>4], digits[c.R&0xf],
		digits[c.G>>4], digits[c.G&0xf],
		digits[c.B>>4], digits[c.B&0xf]})
}

// wantPx parses a px length back to its float32 value.
func wantPx(t *testing.T, name, v string) float32 {
	t.Helper()
	num, ok := strings.CutSuffix(v, "px")
	if !ok {
		t.Fatalf("%s: value %q is not a px length", name, v)
	}
	f, err := strconv.ParseFloat(num, 32)
	if err != nil {
		t.Fatalf("%s: value %q: %v", name, v, err)
	}
	return float32(f)
}

func writeDefault(t *testing.T) (Snapshot, map[string]map[string]string, []byte) {
	t.Helper()
	snap, err := Capture(theme.Default())
	if err != nil {
		t.Fatalf("Capture(theme.Default()): %v", err)
	}
	dir := t.TempDir()
	if err := Write(dir, snap); err != nil {
		t.Fatalf("Write: %v", err)
	}
	css, err := os.ReadFile(filepath.Join(dir, "styles.css"))
	if err != nil {
		t.Fatal(err)
	}
	js, err := os.ReadFile(filepath.Join(dir, "theme.json"))
	if err != nil {
		t.Fatal(err)
	}
	return snap, parseSheet(t, string(css)), js
}

// TestRoundTripColors parses the emitted CSS back and asserts every colour
// variable in both blocks equals the Go token it came from.
func TestRoundTripColors(t *testing.T) {
	snap, sheet, _ := writeDefault(t)
	root, dark := sheet[":root"], sheet[".dark"]
	if root == nil || dark == nil {
		t.Fatalf("styles.css must carry a :root and a .dark block; got %v", len(sheet))
	}

	schemes := []struct {
		vars     map[string]string
		platform tokens.PlatformColors
	}{{root, snap.PlatformLight}, {dark, snap.PlatformDark}}

	for _, scheme := range schemes {
		// The platform's own set, one property per field of
		// tokens.PlatformColors, its coverage written out where it carries
		// one — a label or a seam is a colour AT a coverage and a sheet
		// that dropped it would state a colour the platform never paints.
		for _, n := range platformNames {
			name := "--platform-" + n.name
			if got, want := scheme.vars[name], hexRGBA(n.pick(scheme.platform)); got != want {
				t.Errorf("%s = %q, want %q", name, got, want)
			}
		}
	}

	// The dark block carries exactly the overrides that resolve against an
	// appearance — every variable it declares must exist in :root, and
	// nothing but the platform's own set may differ between the two.
	for name := range dark {
		if _, ok := root[name]; !ok {
			t.Errorf(".dark declares %s which :root does not", name)
		}
		if !strings.HasPrefix(name, "--platform-") {
			t.Errorf(".dark declares non-scheme variable %s", name)
		}
	}
	if want := len(platformNames); len(dark) != want {
		t.Errorf(".dark declares %d variables, want %d", len(dark), want)
	}
}

// TestRoundTripScales asserts the font, spacing, radius and shadow
// variables all parse back to the Go values they came from.
func TestRoundTripScales(t *testing.T) {
	snap, sheet, _ := writeDefault(t)
	root := sheet[":root"]

	if got, err := strconv.Unquote(root["--font-family"]); err != nil || got != snap.Typography.BodyLarge.Typeface {
		t.Errorf("--font-family = %q (%v), want quoted %q", root["--font-family"], err, snap.Typography.BodyLarge.Typeface)
	}
	if got, err := strconv.Unquote(root["--font-family-code"]); err != nil || got != snap.Typography.Code.Typeface {
		t.Errorf("--font-family-code = %q (%v), want quoted %q", root["--font-family-code"], err, snap.Typography.Code.Typeface)
	}
	for _, role := range typeRoles {
		style := role.pick(snap.Typography)
		base := "--font-" + role.name
		if got := wantPx(t, base+"-size", root[base+"-size"]); got != style.Size {
			t.Errorf("%s-size = %v, want %v", base, got, style.Size)
		}
		if got := wantPx(t, base+"-line-height", root[base+"-line-height"]); got != style.LineHeight {
			t.Errorf("%s-line-height = %v, want %v", base, got, style.LineHeight)
		}
		if got, err := strconv.Atoi(root[base+"-weight"]); err != nil || got != style.Weight {
			t.Errorf("%s-weight = %q, want %d", base, root[base+"-weight"], style.Weight)
		}
		if got := wantPx(t, base+"-tracking", root[base+"-tracking"]); got != style.Tracking {
			t.Errorf("%s-tracking = %v, want %v", base, got, style.Tracking)
		}
	}

	for _, key := range spaceKeys {
		name := "--space-" + key.name
		if got := wantPx(t, name, root[name]); got != key.pick(snap.Spacing) {
			t.Errorf("%s = %v, want %v", name, got, key.pick(snap.Spacing))
		}
	}
	for _, key := range radiusKeys {
		name := "--radius-" + key.name
		if got := wantPx(t, name, root[name]); got != key.pick(snap.Radius) {
			t.Errorf("%s = %v, want %v", name, got, key.pick(snap.Radius))
		}
	}

	for _, level := range shadowLevels {
		name, dp := "--shadow-"+level.name, level.dp(snap.Elevation)
		v := root[name]
		if dp == 0 {
			if v != "none" {
				t.Errorf("%s = %q, want \"none\" at depth 0", name, v)
			}
			continue
		}
		mid, ok := strings.CutPrefix(v, "0 ")
		if !ok {
			t.Errorf("%s = %q: want a y-offset shadow with no x-offset", name, v)
			continue
		}
		mid, ok = strings.CutSuffix(mid, " 0 rgba(0, 0, 0, 0.2)")
		if !ok {
			t.Errorf("%s = %q: want no spread and black at 20%%", name, v)
			continue
		}
		lengths := strings.Fields(mid)
		if len(lengths) != 2 {
			t.Errorf("%s = %q: want a y-offset and a blur", name, v)
			continue
		}
		y, blur := wantPx(t, name, lengths[0]), wantPx(t, name, lengths[1])
		if y != dp || blur != 2*dp {
			t.Errorf("%s = %q: y %v blur %v, want dp %v and 2dp", name, v, y, blur, dp)
		}
	}
}

// TestRoundTripDensity asserts the density variables parse back to the Go
// settings: :root carries tokens.Comfortable plus the invariant hit-target
// floor, and the .compact block overrides exactly the three per-setting
// metrics with tokens.Compact's — never the hit target.
func TestRoundTripDensity(t *testing.T) {
	_, sheet, _ := writeDefault(t)
	root, compact := sheet[":root"], sheet[".compact"]
	if compact == nil {
		t.Fatalf("styles.css must carry a .compact block")
	}

	for _, m := range densityMetrics {
		name := "--density-" + m.name
		if got := wantPx(t, name, root[name]); got != m.pick(tokens.Comfortable) {
			t.Errorf(":root %s = %v, want comfortable %v", name, got, m.pick(tokens.Comfortable))
		}
		if got := wantPx(t, name, compact[name]); got != m.pick(tokens.Compact) {
			t.Errorf(".compact %s = %v, want compact %v", name, got, m.pick(tokens.Compact))
		}
	}
	name := "--density-min-hit-target"
	if got := wantPx(t, name, root[name]); got != tokens.MinHitTarget {
		t.Errorf("%s = %v, want %v", name, got, tokens.MinHitTarget)
	}

	// The compact block carries exactly the per-setting overrides: every
	// variable it declares exists in :root, is a --density-* metric, and the
	// hit-target floor is not among them.
	for n := range compact {
		if _, ok := root[n]; !ok {
			t.Errorf(".compact declares %s which :root does not", n)
		}
		if !strings.HasPrefix(n, "--density-") {
			t.Errorf(".compact declares non-density variable %s", n)
		}
	}
	if _, ok := compact[name]; ok {
		t.Errorf(".compact overrides %s; the WCAG 2.5.5 pointer-target floor must not scale with density", name)
	}
	if want := len(densityMetrics); len(compact) != want {
		t.Errorf(".compact declares %d variables, want %d", len(compact), want)
	}
}

// TestRoundTripMotion asserts the easing variables parse back structurally
// — cubic-bezier() with the Go control points — and the duration stops
// numerically in ms.
func TestRoundTripMotion(t *testing.T) {
	snap, sheet, _ := writeDefault(t)
	root := sheet[":root"]

	for _, role := range easeRoles {
		name := "--ease-" + role.name
		// Parse into float32: the emitted decimals are shortest float32
		// representations, so the round-trip contract is that they parse
		// back to the exact float32 control points.
		var p [4]float32
		if _, err := fmt.Sscanf(root[name], "cubic-bezier(%f, %f, %f, %f)", &p[0], &p[1], &p[2], &p[3]); err != nil {
			t.Errorf("%s = %q: not a cubic-bezier(): %v", name, root[name], err)
			continue
		}
		bz := role.pick(snap.Motion)
		if want := [4]float32{bz.P1[0], bz.P1[1], bz.P2[0], bz.P2[1]}; p != want {
			t.Errorf("%s = %v, want control points %v", name, p, want)
		}
	}

	for _, stop := range durationStops {
		name := "--duration-" + stop.name
		num, ok := strings.CutSuffix(root[name], "ms")
		if !ok {
			t.Errorf("%s = %q: not a ms duration", name, root[name])
			continue
		}
		f, err := strconv.ParseFloat(num, 64)
		if err != nil {
			t.Errorf("%s = %q: %v", name, root[name], err)
			continue
		}
		if want := float64(stop.pick(snap.Motion)) / 1e6; f != want {
			t.Errorf("%s = %v ms, want %v ms", name, f, want)
		}
	}
}

func TestRoundTripButtonClasses(t *testing.T) {
	snap, sheet, _ := writeDefault(t)
	root := sheet[":root"]

	if got := wantPx(t, "--focus-ring-width", root["--focus-ring-width"]); got != 2 {
		t.Errorf("--focus-ring-width = %v, want the 2 dp stroke components/button draws", got)
	}
	// The class layer itself: the platform's names only, and no literal
	// colour anywhere.
	src := stylesCSS(snap)
	idx := strings.Index(src, ".btn")
	if idx < 0 {
		t.Fatal("styles.css has no .btn class layer")
	}
	classes := src[idx:]
	if strings.Contains(classes, "#") {
		t.Error("the class layer contains a literal colour; every value must be a token reference")
	}
	// Not one name from the derived set survives in the layer: those
	// variables are still emitted above it for the pages that have not
	// converted, and a class reaching back up to them is the drift this
	// assertion exists to catch.
	for _, gone := range []string{"var(--color-", "var(--elevation-", "var(--shadow-", "var(--surface-"} {
		if strings.Contains(classes, gone) {
			t.Errorf("the class layer still names %s; every colour in it is the platform's", gone)
		}
	}
	for _, frag := range []string{
		// Structure from the density, radius and label-large tokens.
		"min-height: var(--density-control-height);",
		"padding: var(--density-padding-y) var(--density-padding-x);",
		"border-radius: var(--radius-md);",
		"font-size: var(--font-label-large-size);",
		// Filled is the platform's default action.
		"background: var(--platform-control-accent);",
		"color: var(--platform-alternate-selected-control-text);",
		// Tonal is the platform's ordinary push button, inside the seam.
		"background: var(--platform-push-button-fill);",
		"border: 1px solid var(--platform-separator);",
		"color: var(--platform-control-text);",
		// Ghost carries no fill at all.
		".btn.ghost {",
		"background: transparent;",
		// No hover rule in any variant: a push button does not tint under
		// the pointer on this platform. Held, the press overlay goes over
		// the fill as a one-colour gradient layer.
		".btn:active, .btn.is-active {",
		"background-image: linear-gradient(var(--platform-press-overlay), var(--platform-press-overlay));",
		// One ring, one width, every variant, and its forcing twins: a
		// static page shows a state through a class grouped into the same
		// rule as the live pseudo-class, never through duplicated
		// declarations.
		"outline: var(--focus-ring-width) solid var(--platform-keyboard-focus-indicator);",
		".btn:focus-visible, .btn.is-focus {",
		".checkbox:focus-visible, .checkbox.is-focus,",
		".radio:focus-visible, .radio.is-focus {",
		// Disabled is the platform's pair, not a fade of the resting colours.
		"color: var(--platform-disabled-control-text);",
		// Icon-only: a control-height square, glyph inset by PaddingY.
		"width: var(--density-control-height);",
		"padding: var(--density-padding-y);",
		// Badge: the platform's system colour for the status under white.
		".badge {",
		"font-size: var(--font-label-medium-size);",
		"background: var(--platform-system-gray);",
		".badge.success { background: var(--platform-system-green); }",
		".badge.warning { background: var(--platform-system-orange); }",
		".badge.error { background: var(--platform-system-red); }",
		".badge.info { background: var(--platform-system-blue); }",
		// Forms: the platform's text background under its text, the
		// measured field edge at rest, the field's own height floor, and
		// the padding-box clip that puts a focus ring over the surface
		// rather than over the field's own fill.
		"border: 1px solid var(--platform-field-edge);",
		"background: var(--platform-text-background);",
		"background-clip: padding-box;",
		"min-height: var(--density-field-height);",
		"font-size: var(--font-body-large-size);",
		".input::placeholder { color: var(--platform-placeholder-text); opacity: 1; }",
		"border-color: var(--platform-keyboard-focus-indicator);",
		// The dropdown trigger is a button, not a field.
		".select {",
		"background-color: var(--platform-push-button-fill);",
		"border-top: 8px solid var(--platform-secondary-label);",
		// Checkbox/radio: the 16 dp measured glyph inside the 2 dp field
		// edge; checked is the accent under a mark drawn out of gradients
		// rather than encoded as an image, so the no-literal guard above
		// still holds over the whole layer.
		"border: 2px solid var(--platform-field-edge);",
		".checkbox:checked, .checkbox.is-checked {",
		"background-position: 2.529px 7.529px, 5.529px 3.529px;",
		"background-size: 3.943px 3.943px, 7.943px 7.943px;",
		"radial-gradient(circle, var(--platform-alternate-selected-control-text) 4px, var(--platform-control-accent) 4px)",
		// Card and group: the grouped box with no hairline, and the
		// hairline with no box.
		"background: var(--platform-card-fill);",
		".group {",
		"padding: calc(var(--space-4) - 1px);",
		"border-radius: var(--radius-lg);",
		"color: var(--platform-secondary-label);",
		// Table: the content fill, the platform's row pitch, the grid
		// closing each row and the separator closing the header band.
		"background: var(--platform-control-background);",
		"height: var(--density-row-height);",
		"box-shadow: inset 0 -1px 0 var(--platform-grid);",
		".table tbody tr:nth-child(even) td {",
		"background: var(--platform-alternating-content-background);",
		"padding: 0 var(--space-3);",
		".table th.sort-asc::after { border-bottom: 5px solid var(--platform-header-text); }",
		".table th.sort-desc::after { border-top: 5px solid var(--platform-header-text); }",
		// Navigation: the chrome material under the platform's label, the
		// separator where two flush regions meet, the tab strip's selected
		// underline in the platform's selection colour, and the sidebar's
		// own measured pill in the accent.
		"min-height: calc(var(--density-control-height) + 2 * var(--density-padding-y));",
		"background: var(--platform-sidebar-material);",
		"box-shadow: inset 0 -1px 0 var(--platform-separator);",
		"box-shadow: inset -1px 0 0 var(--platform-separator);",
		".navbar-link.selected, .tab.selected {",
		"border-bottom-color: var(--platform-selected-content-background);",
		".sidebar.collapsed { width: 48px; }",
		"height: 32px;  /* RowHeight */",
		".sidebar-item.selected::before {",
		"inset: 0 10px;  /* SelectionInset */",
		"border-radius: 8px;  /* SelectionRadius */",
		"background: var(--platform-sidebar-selection);",
		// Breadcrumb: the ancestors are links, the current segment the label.
		"font-size: var(--font-title-small-size);",
		"color: var(--platform-link);",
		".crumbs .crumb:last-child, .crumb.current { color: var(--platform-label); }",
		"border-left: 6px solid var(--platform-secondary-label);",
		// Overlays: the platform's scrim, the window background under the
		// measured shadow, and the separator where a still surface needs a
		// line.
		"background: var(--platform-scrim);",
		"min-width: 180px;",
		"max-width: 560px;",
		"min-height: 120px;",
		"padding: var(--space-5);",
		"font-size: var(--font-title-medium-size);",
		".dialog-footer {",
		"justify-content: flex-end;",
		"background: var(--platform-window-background);",
		"padding: calc(var(--space-3) - 1px);",
		"border-top: 6px solid var(--platform-window-background);",
		"--floating-shadow-step: color-mix(in srgb, var(--platform-floating-shadow) 12.5%, transparent);",
		"0 0 0 24px var(--floating-shadow-step);",
		// Tooltip: the window background inside the separator, S2/S1 padding
		// measured from the outer edge.
		"border-radius: var(--radius-sm);",
		"padding: calc(var(--space-1) - 1px) calc(var(--space-2) - 1px);",
		// Toast: the window background under the platform's label, the
		// status carried on a leading edge one S2 wide in the system colour
		// for it.
		"background: linear-gradient(to right, var(--platform-system-blue) 0 var(--space-2), var(--platform-window-background) var(--space-2));",
		"min-height: 36px;",
		"padding-left: calc(var(--space-2) + var(--space-3));",
		"background: linear-gradient(to right, var(--platform-system-green) 0 var(--space-2), var(--platform-window-background) var(--space-2));",
		"background: linear-gradient(to right, var(--platform-system-orange) 0 var(--space-2), var(--platform-window-background) var(--space-2));",
		"background: linear-gradient(to right, var(--platform-system-red) 0 var(--space-2), var(--platform-window-background) var(--space-2));",
	} {
		if !strings.Contains(classes, frag) {
			t.Errorf("class layer lacks %q", frag)
		}
	}

}

// TestThemeJSONReproduces asserts theme.json's reproducibility claim: every
// recorded parameter matches the tokens and the sheet, value for value, so
// the file alone rebuilds what was exported.
func TestThemeJSONReproduces(t *testing.T) {
	snap, sheet, js := writeDefault(t)
	var p Parameters
	if err := json.Unmarshal(js, &p); err != nil {
		t.Fatalf("theme.json: %v", err)
	}

	// The theme colour, under the key theme/brand's own file carries, so
	// an exported theme.json loads there without translation.
	if got, want := p.ThemeColor, wantHex(snap.PlatformLight.ControlAccent); got != want {
		t.Errorf("the theme colour is %q, want %q", got, want)
	}

	// The platform's set, name for name, against the sheet's own blocks.
	root, darkVars := sheet[":root"], sheet[".dark"]
	for _, mode := range []struct {
		name string
		set  map[string]string
		vars map[string]string
	}{{":root", p.Platform.Light, root}, {".dark", p.Platform.Dark, darkVars}} {
		if len(mode.set) != len(platformNames) {
			t.Errorf("theme.json platform.%s carries %d names, want %d", mode.name, len(mode.set), len(platformNames))
		}
		for _, n := range platformNames {
			if got, want := mode.set[n.name], mode.vars["--platform-"+n.name]; got != want {
				t.Errorf("theme.json platform.%s[%s] = %q, sheet says %q", mode.name, n.name, got, want)
			}
		}
	}

	if p.Fonts.Heading != "Roboto" || p.Fonts.Body != "Roboto" || p.Fonts.Mono != "Roboto Mono" {
		t.Errorf("fonts = %+v, want Roboto/Roboto/Roboto Mono", p.Fonts)
	}
	if p.Radius != float64(snap.Radius.Base) {
		t.Errorf("radius = %v, want the base radius %v", p.Radius, snap.Radius.Base)
	}
	// Density: the active setting by name, both published settings' metrics,
	// and the invariant floor.
	if p.Density.Setting != "comfortable" {
		t.Errorf("density.setting = %q, want %q (theme.Default() emits tokens.Comfortable)", p.Density.Setting, "comfortable")
	}
	wantMetrics := func(label string, got DensityMetrics, d tokens.Density) {
		want := DensityMetrics{
			ControlHeight: float64(d.ControlHeight),
			ChipHeight:    float64(d.ChipHeight()),
			PaddingX:      float64(d.PaddingX),
			PaddingY:      float64(d.PaddingY),
		}
		if got != want {
			t.Errorf("density.%s = %+v, want %+v", label, got, want)
		}
	}
	wantMetrics("comfortable", p.Density.Comfortable, tokens.Comfortable)
	wantMetrics("compact", p.Density.Compact, tokens.Compact)
	if p.Density.MinHitTarget != float64(tokens.MinHitTarget) {
		t.Errorf("density.minHitTarget = %v, want %v", p.Density.MinHitTarget, tokens.MinHitTarget)
	}

	// The shadow depth each level casts, off the captured scale.
	for i, level := range shadowLevels {
		if got, want := p.Elevation.ShadowDp[i], float64(level.dp(snap.Elevation)); got != want {
			t.Errorf("elevation.shadowDp[%d] = %v, want %v", i, got, want)
		}
	}

	// Motion: durations in ms, easing control points and springs must all
	// reproduce the captured scale exactly (float64(float32) round-trips
	// through JSON unchanged).
	durs := []struct {
		name string
		got  float64
		want time.Duration
	}{
		{"xFast", p.Motion.Durations.XFast, snap.Motion.DurXFast},
		{"fast", p.Motion.Durations.Fast, snap.Motion.DurFast},
		{"normal", p.Motion.Durations.Normal, snap.Motion.DurNormal},
		{"slow", p.Motion.Durations.Slow, snap.Motion.DurSlow},
		{"xSlow", p.Motion.Durations.XSlow, snap.Motion.DurXSlow},
	}
	for _, d := range durs {
		if d.got != float64(d.want)/1e6 {
			t.Errorf("motion.durations.%s = %v ms, want %v ms", d.name, d.got, float64(d.want)/1e6)
		}
	}
	eases := []struct {
		name string
		got  [4]float64
		want tokens.Bezier
	}{
		{"standard", p.Motion.Easings.Standard, snap.Motion.EaseStandard},
		{"standardAccelerate", p.Motion.Easings.StandardAccelerate, snap.Motion.EaseStandardAccelerate},
		{"standardDecelerate", p.Motion.Easings.StandardDecelerate, snap.Motion.EaseStandardDecelerate},
		{"emphasized", p.Motion.Easings.Emphasized, snap.Motion.EaseEmphasized},
		{"emphasizedAccelerate", p.Motion.Easings.EmphasizedAccelerate, snap.Motion.EaseEmphasizedAccelerate},
		{"emphasizedDecelerate", p.Motion.Easings.EmphasizedDecelerate, snap.Motion.EaseEmphasizedDecelerate},
	}
	for _, e := range eases {
		// The file records shortest-float32 decimals; converting the parsed
		// float64s back to float32 must land on the exact control points.
		got := [4]float32{float32(e.got[0]), float32(e.got[1]), float32(e.got[2]), float32(e.got[3])}
		if want := [4]float32{e.want.P1[0], e.want.P1[1], e.want.P2[0], e.want.P2[1]}; got != want {
			t.Errorf("motion.easings.%s = %v, want %v", e.name, e.got, want)
		}
	}
	springs := []struct {
		name string
		got  SpringParam
		want tokens.Spring
	}{
		{"default", p.Motion.Springs.Default, snap.Motion.SpringDefault},
		{"snappy", p.Motion.Springs.Snappy, snap.Motion.SpringSnappy},
		{"gentle", p.Motion.Springs.Gentle, snap.Motion.SpringGentle},
	}
	for _, sp := range springs {
		// Same shortest-float32 contract as the easings: the recorded
		// decimals must reproduce the Go float32s bit-for-bit.
		got := tokens.Spring{
			Mass:      float32(sp.got.Mass),
			Stiffness: float32(sp.got.Stiffness),
			Damping:   float32(sp.got.Damping),
		}
		if got != sp.want {
			t.Errorf("motion.springs.%s = %+v, want %+v", sp.name, sp.got, sp.want)
		}
	}
}

// TestCaptureRejectsIrreproducible asserts Capture refuses inputs
// theme.json could not honestly reproduce.
func TestCaptureRejectsIrreproducible(t *testing.T) {
	if _, err := Capture(theme.Theme{}); err == nil {
		t.Error("Capture of a zero Theme (nil observables) must error")
	}

	th := theme.Default()
	th.Density = rx.Of(tokens.Density{ControlHeight: 30, PaddingX: 10, PaddingY: 5})
	if _, err := Capture(th); err == nil {
		t.Error("Capture of a non-preset density must error: theme.json records density as a named setting")
	}
}

// TestCaptureAChosenThemeColour asserts a set whose accent rows were rebuilt
// for a chosen colour captures with that colour, and that the dark
// counterpart carries it too.
func TestCaptureAChosenThemeColour(t *testing.T) {
	chosen := stdcolor.NRGBA{R: 0x00, G: 0x68, B: 0x74, A: 0xff}
	th := theme.Default()
	th.Platform = rx.Of(tokens.PlatformLight.WithAccent(chosen))
	snap, err := Capture(th)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if got := snap.PlatformLight.ControlAccent; got != chosen {
		t.Errorf("PlatformLight.ControlAccent = %v, want %v", got, chosen)
	}
	if want := tokens.PlatformDark.WithAccent(chosen); snap.PlatformDark != want {
		t.Error("the dark counterpart does not carry the chosen theme colour")
	}
}

// TestPlatformNamesCoverEveryField walks tokens.PlatformColors by reflection
// and fails if a field has no entry in platformNames. A name the set carries
// and the sheet drops is a colour a class-layer rule can reference and never
// resolve — which is how the sidebar's pill lost its fill once.
func TestPlatformNamesCoverEveryField(t *testing.T) {
	byName := make(map[string]bool, len(platformNames))
	for _, n := range platformNames {
		byName[n.name] = true
	}
	typ := reflect.TypeOf(tokens.PlatformColors{})
	if typ.NumField() != len(platformNames) {
		t.Errorf("PlatformColors has %d fields and platformNames %d entries", typ.NumField(), len(platformNames))
	}
	for i := range typ.NumField() {
		if want := kebab(typ.Field(i).Name); !byName[want] {
			t.Errorf("PlatformColors.%s has no --platform-%s entry", typ.Field(i).Name, want)
		}
	}
}

// kebab is the naming rule platformNames follows: the Go field name with a
// hyphen before each interior capital, lowercased.
func kebab(field string) string {
	var b strings.Builder
	for i, r := range field {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('-')
		}
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
		}
		b.WriteRune(r)
	}
	return b.String()
}
