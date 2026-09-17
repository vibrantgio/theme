package export

import (
	"fmt"
	stdcolor "image/color"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/theme"
)

// pageFiles are the foundation pages Write must emit, relative to the
// project root.
var pageFiles = []string{
	filepath.Join("foundations", "color.html"),
	filepath.Join("foundations", "type.html"),
	filepath.Join("foundations", "layout.html"),
}

// writeProject writes the default theme's full project into a temp dir and
// returns the snapshot, the parsed sheet and the page sources by name.
func writeProject(t *testing.T) (Snapshot, map[string]map[string]string, map[string]string) {
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
	pages := map[string]string{}
	for _, name := range append([]string{"readme.md"}, pageFiles...) {
		src, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("Write must emit %s: %v", name, err)
		}
		if len(src) == 0 {
			t.Fatalf("%s is empty", name)
		}
		pages[name] = string(src)
	}
	return snap, parseSheet(t, string(css)), pages
}

var varRefRE = regexp.MustCompile(`var\((--[a-zA-Z0-9-]+)\)`)

// TestPageVarClosure asserts every var() reference in every page resolves in
// the emitted sheet's :root block — a page can never name a token the sheet
// does not declare — and that every referenced colour variable is also
// overridden in .dark, so the toggle restyles all of them.
func TestPageVarClosure(t *testing.T) {
	_, sheet, pages := writeProject(t)
	root, dark := sheet[":root"], sheet[".dark"]
	for _, name := range pageFiles {
		refs := varRefRE.FindAllStringSubmatch(pages[name], -1)
		if len(refs) == 0 {
			t.Errorf("%s references no token variables at all", name)
		}
		for _, ref := range refs {
			v := ref[1]
			if _, ok := root[v]; !ok {
				t.Errorf("%s references %s, which styles.css :root does not declare", name, v)
			}
			if strings.HasPrefix(v, "--color-") {
				if _, ok := dark[v]; !ok {
					t.Errorf("%s references colour %s, which the .dark block does not override", name, v)
				}
			}
		}
	}
}

var (
	styleBlockRE = regexp.MustCompile(`(?s)<style>(.*?)</style>`)
	styleAttrRE  = regexp.MustCompile(`style="([^"]*)"`)
	hexLitRE     = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	pxLitRE      = regexp.MustCompile(`\d(?:\.\d+)?px\b`)
)

// TestPagesNoHardCodedTokenValues enforces the task's rule: the pages read
// only from the token sheet. In every style context — <style> blocks and
// style attributes — no literal hex colour and no literal px length may
// appear; token values reach styling exclusively through var() references.
// (Hexes and px numbers in annotation text are reader-facing data, not
// styling, and are out of scope by construction.)
func TestPagesNoHardCodedTokenValues(t *testing.T) {
	_, _, pages := writeProject(t)
	for _, name := range pageFiles {
		src := pages[name]
		var contexts []string
		for _, m := range styleBlockRE.FindAllStringSubmatch(src, -1) {
			contexts = append(contexts, m[1])
		}
		attrs := styleAttrRE.FindAllStringSubmatch(src, -1)
		if len(attrs) == 0 {
			t.Errorf("%s has no style attributes; the specimens are expected to be var()-styled inline", name)
		}
		for _, m := range attrs {
			contexts = append(contexts, m[1])
		}
		for _, css := range contexts {
			if hit := hexLitRE.FindString(css); hit != "" {
				t.Errorf("%s: hard-coded colour %q in a style context", name, hit)
			}
			if hit := pxLitRE.FindString(css); hit != "" {
				t.Errorf("%s: hard-coded length %q in a style context", name, hit)
			}
		}
	}
}

// TestPagesDarkToggle asserts each page links the shared sheet relatively
// and carries the light/dark toggle script flipping .dark on the root
// element — the mechanism the dark render confirmation rests on.
func TestPagesDarkToggle(t *testing.T) {
	_, _, pages := writeProject(t)
	for _, name := range pageFiles {
		src := pages[name]
		if !strings.Contains(src, `<link rel="stylesheet" href="../styles.css">`) {
			t.Errorf("%s does not link ../styles.css", name)
		}
		if !strings.Contains(src, `classList.toggle("dark")`) {
			t.Errorf("%s has no .dark toggle script", name)
		}
	}
}

// wantRow renders a contrast table row the way the colour page must,
// written out independently so the page and the test cannot drift together:
// APCA Lc, signed, one decimal, light then dark.
func wantRow(label string, lt, lg, dt, dg stdcolor.NRGBA) string {
	return fmt.Sprintf(`<tr><th scope="row">%s</th><td>%.1f</td><td>%.1f</td></tr>`,
		label, color.APCA(lt, lg), color.APCA(dt, dg))
}

// TestColorPageAnnotatesContrast asserts the colour page carries the
// measured APCA Lc, in both appearances, for every pairing the platform
// itself makes — each foreground flattened over the fill beneath it first,
// which is the only form APCA can be handed.
func TestColorPageAnnotatesContrast(t *testing.T) {
	snap, _, pages := writeProject(t)
	src := pages[filepath.Join("foundations", "color.html")]

	for _, pair := range platformPairs {
		lightFill, darkFill := pair.fill(snap.PlatformLight), pair.fill(snap.PlatformDark)
		row := wantRow(pair.label,
			color.Flatten(pair.text(snap.PlatformLight), lightFill), lightFill,
			color.Flatten(pair.text(snap.PlatformDark), darkFill), darkFill)
		if !strings.Contains(src, row) {
			t.Errorf("color.html lacks the measured row for %q:\n%s", pair.label, row)
		}
	}
}

// TestColorPageAnnotatesBothModeValues spot-checks that swatch annotations
// carry both appearances' hexes, labelled, since text cannot flip with the
// class.
func TestColorPageAnnotatesBothModeValues(t *testing.T) {
	snap, _, pages := writeProject(t)
	src := pages[filepath.Join("foundations", "color.html")]
	for _, n := range platformNames {
		want := fmt.Sprintf("L %s &middot; D %s", hexRGBA(n.pick(snap.PlatformLight)), hexRGBA(n.pick(snap.PlatformDark)))
		if !strings.Contains(src, want) {
			t.Errorf("color.html lacks the dual-appearance annotation for %s: %q", n.name, want)
		}
	}
}

// TestReadmeNamesFamilies asserts readme.md names every token family the
// sheet emits, deriving the expected mentions from the same tables the
// emitter renders from.
func TestReadmeNamesFamilies(t *testing.T) {
	snap, sheet, pages := writeProject(t)
	readme := pages["readme.md"]

	var want []string
	for _, n := range platformNames {
		want = append(want, "--platform-"+n.name)
	}
	want = append(want, "--font-family", "-size", "-line-height", "-weight", "-tracking")
	for _, role := range typeRoles {
		want = append(want, role.name)
	}
	for _, key := range spaceKeys {
		want = append(want, "--space-"+key.name)
	}
	for _, key := range radiusKeys {
		want = append(want, "--radius-"+key.name)
	}
	for _, level := range shadowLevels {
		want = append(want, "--shadow-"+level.name)
	}
	for _, m := range densityMetrics {
		want = append(want, "--density-"+m.name)
	}
	for _, role := range easeRoles {
		want = append(want, "--ease-"+role.name)
	}
	for _, stop := range durationStops {
		want = append(want, "--duration-"+stop.name)
	}
	want = append(want,
		wantHex(snap.PlatformLight.ControlAccent),
		"styles.css", "theme.json",
		"foundations/color.html", "foundations/type.html", "foundations/layout.html",
	)
	for _, w := range want {
		if !strings.Contains(readme, w) {
			t.Errorf("readme.md does not mention %q", w)
		}
	}

	// Paranoia in the other direction: every variable the sheet actually
	// emits must be documented, so a new family cannot ship unnamed. Each
	// variable maps to the string the readme must contain for it.
	for name := range sheet[":root"] {
		mention := name // --platform-*, --space-*, --radius-*, --shadow-*: listed in full
		if role, metric, ok := fontMetric(name); ok {
			if !strings.Contains(readme, metric) {
				t.Errorf("readme.md does not mention the %q metric suffix for %s", metric, name)
			}
			mention = role
		}
		if !strings.Contains(readme, mention) {
			t.Errorf("sheet variable %s is undocumented: readme.md lacks %q", name, mention)
		}
	}
}

// TestLayoutPageDensityAndShadow asserts the layout page's contract:
// the control metrics render at BOTH density settings — the compact column
// is the same markup inside a .compact wrapper, exercising the sheet's
// override block — and the shadow section shows each level's depth.
func TestLayoutPageDensityAndShadow(t *testing.T) {
	_, _, pages := writeProject(t)
	src := pages[filepath.Join("foundations", "layout.html")]

	if !strings.Contains(src, `class="density-col compact"`) {
		t.Error("layout.html has no .compact density column; both settings must render side by side")
	}
	for _, m := range densityMetrics {
		if !strings.Contains(src, "var(--density-"+m.name+")") {
			t.Errorf("layout.html does not style through var(--density-%s)", m.name)
		}
	}
	for _, level := range shadowLevels {
		if !strings.Contains(src, fmt.Sprintf("box-shadow: var(--shadow-%s)", level.name)) {
			t.Errorf("layout.html does not show the shadow var(--shadow-%s)", level.name)
		}
	}
}

// fontMetric splits a --font-<role>-<metric> variable into the role name
// and its metric suffix; --font-family is not a per-role metric.
func fontMetric(name string) (role, metric string, ok bool) {
	base, found := strings.CutPrefix(name, "--font-")
	if !found || name == "--font-family" {
		return "", "", false
	}
	for _, suf := range []string{"-size", "-line-height", "-weight", "-tracking"} {
		if r, has := strings.CutSuffix(base, suf); has {
			return r, suf, true
		}
	}
	return "", "", false
}
