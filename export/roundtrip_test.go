package export

import (
	"encoding/json"
	"fmt"
	stdcolor "image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/reactivego/rx"
	"github.com/vibrantgio/theme/color"
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
		tokens   tokens.ColorTokens
		platform tokens.PlatformColors
	}{{root, snap.Light, snap.PlatformLight}, {dark, snap.Dark, snap.PlatformDark}}

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
		for _, role := range rampRoles {
			ramp := role.ramp(scheme.tokens.Ramps)
			for step := 100; step <= 900; step += 100 {
				name := fmt.Sprintf("--color-%s-%d", role.name, step)
				if got, want := scheme.vars[name], wantHex(ramp.Step(step)); got != want {
					t.Errorf("%s = %q, want %q", name, got, want)
				}
			}
		}
		for _, pin := range pinRoles {
			name := "--color-" + pin.name
			if got, want := scheme.vars[name], wantHex(pin.pick(scheme.tokens)); got != want {
				t.Errorf("%s = %q, want %q", name, got, want)
			}
		}
	}

	// The dark block carries exactly the overrides that resolve against a
	// scheme — every variable it declares must exist in :root, and nothing
	// but a colour or an elevation level may differ per mode. A level is
	// placed against the Background pin rather than named as a ramp step,
	// so it resolves per scheme like the walked pins do and cannot be a
	// var() reference the .dark block flips underneath.
	for name := range dark {
		if _, ok := root[name]; !ok {
			t.Errorf(".dark declares %s which :root does not", name)
		}
		if !strings.HasPrefix(name, "--color-") && !strings.HasPrefix(name, "--elevation-") && !strings.HasPrefix(name, "--platform-") {
			t.Errorf(".dark declares non-scheme variable %s", name)
		}
	}
	// Three families per level — the fill, and the hover and press walks
	// taken from it — plus two hairlines every level a thing can stand on
	// carries: the seam it owes what stands on it, and the line two regions
	// sharing its own fill are parted by. All resolve per scheme for the
	// same reason.
	if want := len(rampRoles)*9 + len(pinRoles) + 3*len(elevationLevels) + 2*len(standableLevels) + len(platformNames); len(dark) != want {
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

	for _, level := range elevationLevels {
		name, dp := "--shadow-"+level.name, snap.Elevation.Dp(level.level)
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

// TestRoundTripElevationSurfaces asserts every --elevation-* variable is a
// literal that equals, in its own block, exactly the colour SurfaceAt
// returns for that scheme — the sheet's default elevation cue cannot drift
// from the Go resolver.
//
// They cannot be var() references into the neutral ramp: a level is placed
// against the Background pin in CIELAB L*, so the light scheme's levels
// above the content and the dark scheme's floor are not ramp steps at all and
// no var() chain reaches them. Each block states its own five.
//
// The scale's direction is asserted here too: read down the levels and
// the fill gets lighter, in the :root block and in the .dark one, with no
// mirror clause between them.
func TestRoundTripElevationSurfaces(t *testing.T) {
	snap, sheet, _ := writeDefault(t)
	root, dark := sheet[":root"], sheet[".dark"]

	for _, mode := range []struct {
		name   string
		scheme tokens.ColorTokens
		vars   map[string]string
	}{{":root", snap.Light, root}, {".dark", snap.Dark, dark}} {
		var last float64 = -1
		for _, level := range elevationLevels {
			name := "--elevation-" + level.name
			got, ok := mode.vars[name]
			if !ok {
				t.Fatalf("%s does not declare %s; every scheme states its own scale", mode.name, name)
			}
			fill := mode.scheme.SurfaceAt(level.level)
			if want := wantHex(fill); got != want {
				t.Errorf("%s %s = %q, want SurfaceAt = %q", mode.name, name, got, want)
			}
			if l, _, _ := color.LabFromNRGBA(fill); l < last {
				t.Errorf("%s %s is L*%.2f, under the level below it (L*%.2f)",
					mode.name, name, l, last)
			} else {
				last = l
			}
		}
		// And the seam each standable level owes what stands on it:
		// transparent where the raise is told by its own fill, the derived
		// hairline where it is not. The backdrop declares none — nothing
		// stands on the backdrop.
		if _, ok := mode.vars["--elevation-backdrop-seam"]; ok {
			t.Errorf("%s declares a seam for the backdrop; nothing stands on it", mode.name)
		}
		for _, level := range standableLevels {
			name := "--elevation-" + level.name + "-seam"
			got, ok := mode.vars[name]
			if !ok {
				t.Fatalf("%s does not declare %s", mode.name, name)
			}
			raise := mode.scheme.RaisedOn(mode.scheme.SurfaceAt(level.level))
			want := "transparent"
			if raise.Seamed {
				want = wantHex(raise.Seam)
			}
			if got != want {
				t.Errorf("%s %s = %q, want %q", mode.name, name, got, want)
			}
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
		t.Errorf(".compact overrides %s; the WCAG floor must not scale with density", name)
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

// TestRoundTripButtonClasses asserts two things that must not drift. The
// emitted derived variables still resolve to what their own rules say — the
// walked solid-fill stops equal SolidStateColor's per mode, the ring and the
// control edge are the steps their ramps measure, the disabled fraction is
// DisabledOpacity — because the pages that have not converted still read
// them. And the class layer below them names the PLATFORM's colours and
// nothing else: every rule is a --platform- name, the mapping the Gio
// components take, with not one literal colour and not one reference back up
// into the derived set.
// stepDistance is how far a ramp index sits from step 500, the mid-value
// depth the ring's pick is aimed at. The test measures it for itself rather
// than importing the emitter's constant, so a drift in the aim is a
// failure here rather than a silent agreement.
func stepDistance(i int) int {
	const mid = 4 // steps run 100…900
	if i < mid {
		return mid - i
	}
	return i - mid
}

func TestRoundTripButtonClasses(t *testing.T) {
	snap, sheet, _ := writeDefault(t)
	root, dark := sheet[":root"], sheet[".dark"]

	// The solid-fill state walk, per mode. Written against SolidStateColor
	// directly, not through pinRoles, so the emitter cannot drift with its
	// own table.
	for i, mode := range []struct {
		vars map[string]string
		tok  tokens.ColorTokens
	}{{root, snap.Light}, {dark, snap.Dark}} {
		if got, want := mode.vars["--color-accent-hover"], wantHex(mode.tok.SolidStateColor(tokens.RolePrimary, tokens.StateHover)); got != want {
			t.Errorf("--color-accent-hover (mode %d) = %q, want SolidStateColor hover %q", i, got, want)
		}
		if got, want := mode.vars["--color-accent-pressed"], wantHex(mode.tok.SolidStateColor(tokens.RolePrimary, tokens.StatePressed)); got != want {
			t.Errorf("--color-accent-pressed (mode %d) = %q, want SolidStateColor pressed %q", i, got, want)
		}
		// The ring: one per mode, restated here from the rule rather than
		// called out of the emitter, so the sheet and its generator cannot
		// agree on a wrong answer. The rule is the step of the primary ramp
		// nearest step 500 that reaches [tokens.GraphicFloor] against EVERY
		// level a control can stand on — every level but the backdrop, which
		// nothing is drawn at — 1.25:1 in luminance against every one of
		// those levels' neutral resting borders, and is not the accent fill
		// — the property the Gio side
		// derives by, and the reason no per-level ring token exists to pin.
		// The second floor is what keeps focus from being spelled in hue
		// alone: the resting border is the line a focused field swaps for its
		// ring. The exclusion keeps it from being spelled in the colour a
		// checked box already paints.
		var wantRing stdcolor.NRGBA
		wantAt, clearing := -1, 0
		for r, step := range mode.tok.Ramps.Primary {
			clears := true
			for _, level := range standableLevels {
				surface := mode.tok.SurfaceAt(level.level)
				border := mode.tok.MarkOn(tokens.RoleNeutral, surface, tokens.GraphicFloor)
				if color.Magnitude(step, surface) < tokens.GraphicFloor ||
					luminanceRatio(step, border) < 1.25 {
					clears = false
					break
				}
			}
			if !clears || step == mode.tok.Primary {
				continue
			}
			clearing++
			// Steps run 100…900, so index 4 is step 500. Nearest to it wins;
			// walking upward, a tie keeps the lower step, as the sheet does.
			if wantAt < 0 || stepDistance(r) < stepDistance(wantAt) {
				wantRing, wantAt = step, r
			}
		}
		if clearing == 0 {
			t.Fatalf("mode %d: no step of the primary ramp clears both the ring's floors on every level a control stands on — the sheet's ring rule has nothing to pick", i)
		}
		if got, want := mode.vars["--color-focus-ring"], wantHex(wantRing); got != want {
			t.Errorf("--color-focus-ring (mode %d) = %q, want the primary step nearest step 500 that reads on every level and parts from every resting border %q", i, got, want)
		}
		// The one exception, and the only surface that belongs to no level:
		// the fill a filled button insets its ring in. The scheme's ring
		// serves wherever it reads on that fill; where it cannot — a solid
		// primary fill being a step of the ring's own ramp — the ramp is
		// walked against the fill instead.
		fill := mode.tok.SolidStateColor(tokens.RolePrimary, tokens.StateFocus)
		onAccent := wantRing
		if color.Magnitude(wantRing, fill) < tokens.GraphicFloor {
			onAccent = mode.tok.MarkOn(tokens.RolePrimary, fill, tokens.GraphicFloor)
		}
		if got, want := mode.vars["--color-focus-ring-on-accent"], wantHex(onAccent); got != want {
			t.Errorf("--color-focus-ring-on-accent (mode %d) = %q, want the ring the filled button's own fill can carry %q", i, got, want)
		}
		// The level-varying ring tokens are gone, and their absence is
		// pinned: a ring that depended on the surface is the divergence this
		// sheet exists not to reintroduce.
		for _, gone := range []string{"--color-dialog-focus-ring", "--color-popover-focus-ring"} {
			if got, ok := mode.vars[gone]; ok {
				t.Errorf("%s (mode %d) = %q, want no such token: the ring does not vary with the level", gone, i, got)
			}
		}
		// The control row's resting edge, likewise: components/input's
		// controlBorder is MarkOn against the level-0 surface, and the two
		// outline tokens are the same walk two and three levels up — the
		// fills patterns/modal and patterns/popover paint and measure their
		// own edges against, and the edge any control standing on those
		// levels wears. Level 1 has none: a card draws no line of its own,
		// and a control on a card takes control-border unchanged.
		if got, want := mode.vars["--color-control-border"], wantHex(mode.tok.MarkOn(tokens.RoleNeutral, mode.tok.SurfaceAt(tokens.Level0), tokens.GraphicFloor)); got != want {
			t.Errorf("--color-control-border (mode %d) = %q, want the neutral step that reads on the window surface %q", i, got, want)
		}
		for _, edge := range []struct {
			name  string
			level tokens.ElevationLevel
			what  string
		}{
			{"--color-dialog-border", tokens.Level2, "the dialog's level-2 fill"},
			{"--color-popover-border", tokens.Level3, "the popover's level-3 fill"},
		} {
			want := wantHex(mode.tok.MarkOn(tokens.RoleNeutral, mode.tok.SurfaceAt(edge.level), tokens.GraphicFloor))
			if got := mode.vars[edge.name]; got != want {
				t.Errorf("%s (mode %d) = %q, want the neutral step that reads on %s %q", edge.name, i, got, edge.what, want)
			}
		}
	}
	if got := wantPx(t, "--focus-ring-width", root["--focus-ring-width"]); got != 2 {
		t.Errorf("--focus-ring-width = %v, want the 2 dp stroke components/button draws", got)
	}
	if got, want := root["--state-disabled-opacity"], fmt.Sprintf("%v%%", tokens.DisabledOpacity*100); got != want {
		t.Errorf("--state-disabled-opacity = %q, want %q", got, want)
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
		"background: var(--platform-control-accent);",
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

	// The scrim token: the modal scrim's black at alpha 0x80 in
	// the alpha that reproduces it under sRGB compositing (Gio composites in
	// linear RGB — see scrimRGBA's derivation), mode-invariant like the
	// shadows' fixed black, so it lives in :root and .dark never overrides
	// it.
	if got := root["--color-scrim"]; got != scrimRGBA {
		t.Errorf("--color-scrim = %q, want scrimRGBA %q", got, scrimRGBA)
	}
	if _, ok := dark["--color-scrim"]; ok {
		t.Error("--color-scrim is overridden in .dark; the scrim is mode-invariant")
	}
}

// TestThemeJSONReproduces asserts theme.json's reproducibility claim: the
// one colour it records alone regenerates the exported palette through
// FromSeed, and every recorded parameter matches the tokens and the sheet.
// That colour is the light scheme's primary base — the brand seed with the
// palette's accent chroma on it — and FromSeed reproduces itself from it.
func TestThemeJSONReproduces(t *testing.T) {
	snap, sheet, js := writeDefault(t)
	var p Parameters
	if err := json.Unmarshal(js, &p); err != nil {
		t.Fatalf("theme.json: %v", err)
	}

	var r, g, b uint8
	if _, err := fmt.Sscanf(p.Seed, "#%02x%02x%02x", &r, &g, &b); err != nil {
		t.Fatalf("theme.json seed %q: %v", p.Seed, err)
	}
	seed := stdcolor.NRGBA{R: r, G: g, B: b, A: 0xff}
	if seed != tokens.DefaultLight.Primary {
		t.Errorf("seed = %q, want the default palette's primary base %s", p.Seed, wantHex(tokens.DefaultLight.Primary))
	}
	light, dark := tokens.FromSeed(seed)
	if light != snap.Light || dark != snap.Dark {
		t.Errorf("FromSeed(theme.json seed) does not reproduce the exported palette")
	}

	_, chroma, hue := color.OKLChFromNRGBA(seed)
	if diff := p.Hue - hue; diff < -0.005 || diff > 0.005 {
		t.Errorf("hue = %v, want %v within 0.005", p.Hue, hue)
	}
	if diff := p.Sat - chroma; diff < -0.00005 || diff > 0.00005 {
		t.Errorf("sat = %v, want %v within 0.00005", p.Sat, chroma)
	}

	root, darkVars := sheet[":root"], sheet[".dark"]
	for _, mode := range []struct {
		pins Pins
		vars map[string]string
	}{{p.Pins.Light, root}, {p.Pins.Dark, darkVars}} {
		checks := []struct{ name, got string }{
			{"--color-bg", mode.pins.Bg},
			{"--color-text", mode.pins.Text},
			{"--color-accent", mode.pins.Accent},
			{"--color-secondary", mode.pins.Secondary},
			{"--color-tertiary", mode.pins.Tertiary},
			{"--color-error", mode.pins.Error},
			{"--color-success", mode.pins.Success},
			{"--color-warning", mode.pins.Warning},
		}
		for _, c := range checks {
			if c.got != mode.vars[c.name] {
				t.Errorf("theme.json pin %s = %q, sheet says %q", c.name, c.got, mode.vars[c.name])
			}
		}
	}

	if p.Fonts.Heading != "Roboto" || p.Fonts.Body != "Roboto" || p.Fonts.Mono != "Roboto Mono" {
		t.Errorf("fonts = %+v, want Roboto/Roboto/Roboto Mono", p.Fonts)
	}
	if p.Radius != float64(snap.Radius.Base) {
		t.Errorf("radius = %v, want the base radius %v", p.Radius, snap.Radius.Base)
	}
	if want := [9]int{97, 92, 85, 74, 63, 51, 39, 28, 6}; p.Scale.Light != want {
		t.Errorf("scale.light = %v, want ADR-007's shared scale %v", p.Scale.Light, want)
	}
	if want := [9]int{8, 13, 19, 30, 46, 64, 82, 86, 94}; p.Scale.Dark != want {
		t.Errorf("scale.dark = %v, want the paired dark scale %v", p.Scale.Dark, want)
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

	// Elevation: the level fill per scheme and the shadow dp per level,
	// off the captured snapshot through the same resolver the sheet uses.
	for i, level := range elevationLevels {
		if got, want := p.Elevation.Surfaces.Light[i], hexRGB(snap.Light.SurfaceAt(level.level)); got != want {
			t.Errorf("elevation.surfaces.light[%d] = %q, want %q", i, got, want)
		}
		if got, want := p.Elevation.Surfaces.Dark[i], hexRGB(snap.Dark.SurfaceAt(level.level)); got != want {
			t.Errorf("elevation.surfaces.dark[%d] = %q, want %q", i, got, want)
		}
		if got, want := p.Elevation.ShadowDp[i], float64(snap.Elevation.Dp(level.level)); got != want {
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
	th.Color = rx.Of(tokens.DefaultDark) // a dark scheme: its Primary pin is not the seed
	if _, err := Capture(th); err == nil {
		t.Error("Capture of a dark colour emission must error: FromSeed(pin) cannot reproduce it")
	}

	th = theme.Default()
	th.Density = rx.Of(tokens.Density{ControlHeight: 30, PaddingX: 10, PaddingY: 5})
	if _, err := Capture(th); err == nil {
		t.Error("Capture of a non-preset density must error: theme.json records density as a named setting")
	}
}

// TestCaptureCustomSeed asserts a re-branded light scheme captures with its
// own seed recovered.
func TestCaptureCustomSeed(t *testing.T) {
	seed := stdcolor.NRGBA{R: 0x00, G: 0x68, B: 0x74, A: 0xff}
	light, dark := tokens.FromSeed(seed)
	th := theme.Default()
	th.Color = rx.Of(light)
	snap, err := Capture(th)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if snap.Seed != seed {
		t.Errorf("Seed = %v, want %v", snap.Seed, seed)
	}
	if snap.Dark != dark {
		t.Errorf("Dark scheme is not FromSeed(seed)'s pair")
	}
}
