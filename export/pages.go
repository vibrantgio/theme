package export

import (
	"fmt"
	"html"
	stdcolor "image/color"
	"strings"

	"github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/tokens"
)

// The foundation pages are static HTML that reads only from the emitted
// token sheet: every colour, size, radius, shadow and font value in a style
// position is a var(--...) reference into ../styles.css, so regenerating the
// sheet for another theme colour reflows every page with no page edit. The only
// literal token values in a page are annotation text — hexes, px numbers and
// contrast measurements printed for the reader — which the generator
// computes from the Snapshot at generation time. Chrome CSS (flex, grid,
// margins for the specimen scaffolding) uses tokens where a token fits and
// token-free units (rem, thin) elsewhere; it never carries a literal colour
// or px length, and the page test enforces exactly that.

// chromeCSS is the scaffolding shared by every foundation page. No literal
// hex colours and no literal px lengths: colours, sizes, spacing and radii
// all resolve through the token sheet.
const chromeCSS = `body {
  margin: 0;
  background: var(--platform-window-background);
  color: var(--platform-label);
  font-family: var(--font-family), system-ui, sans-serif;
  font-size: var(--font-body-medium-size);
  line-height: var(--font-body-medium-line-height);
  letter-spacing: var(--font-body-medium-tracking);
}
main {
  max-width: 64rem;
  margin: 0 auto;
  padding: var(--space-6) var(--space-6) var(--space-16);
}
.masthead {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: var(--space-4);
}
h1 {
  font-size: var(--font-headline-large-size);
  line-height: var(--font-headline-large-line-height);
  font-weight: var(--font-headline-large-weight);
  letter-spacing: var(--font-headline-large-tracking);
  margin: var(--space-6) 0 var(--space-2);
}
h2 {
  font-size: var(--font-title-large-size);
  line-height: var(--font-title-large-line-height);
  font-weight: var(--font-title-large-weight);
  letter-spacing: var(--font-title-large-tracking);
  margin: var(--space-10) 0 var(--space-3);
}
h3 {
  font-size: var(--font-title-medium-size);
  line-height: var(--font-title-medium-line-height);
  font-weight: var(--font-title-medium-weight);
  letter-spacing: var(--font-title-medium-tracking);
  margin: var(--space-6) 0 var(--space-2);
}
.intro {
  max-width: 48rem;
  color: var(--platform-secondary-label);
}
.annot {
  font-size: var(--font-label-small-size);
  line-height: var(--font-label-small-line-height);
  font-weight: var(--font-label-small-weight);
  letter-spacing: var(--font-label-small-tracking);
  color: var(--platform-secondary-label);
  margin: var(--space-1) 0 0;
}
.mode-toggle {
  font-family: var(--font-family), system-ui, sans-serif;
  font-size: var(--font-label-large-size);
  line-height: var(--font-label-large-line-height);
  font-weight: var(--font-label-large-weight);
  letter-spacing: var(--font-label-large-tracking);
  padding: var(--space-2) var(--space-4);
  color: var(--platform-label);
  background: var(--platform-card-fill);
  border: thin solid var(--platform-field-edge);
  border-radius: var(--radius-base);
  cursor: pointer;
}
.mode-toggle:hover {
  background: var(--platform-separator);
}
`

// page wraps a body in the shared skeleton: the sheet link (foundations/ is
// a subdirectory, so the sheet is ../styles.css), the chrome, a masthead
// with the light/dark toggle, and the toggle script that flips the .dark
// class on the root element — the sheet's .dark override block restyles
// everything var()-driven from there.
func page(title, heading, intro, style, body string) string {
	var b strings.Builder
	// The @dsCard marker must be the file's first line: claude.ai/design's
	// pane renders cards, not files, and it builds its card index from this
	// comment. A page without it uploads fine and then never appears in the
	// UI — the failure mode is silent, which is why the marker is emitted
	// here rather than remembered. The parser takes exactly this group-only
	// form; a first attempt with an extra name="…" attribute produced no
	// card, so the card's display name comes from the page, not the marker.
	b.WriteString("<!-- @dsCard group=\"Foundations\" -->\n")
	b.WriteString("<!doctype html>\n<html lang=\"en\">\n<head>\n")
	b.WriteString("<meta charset=\"utf-8\">\n")
	b.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	fmt.Fprintf(&b, "<title>%s</title>\n", html.EscapeString(title))
	b.WriteString("<link rel=\"stylesheet\" href=\"../styles.css\">\n")
	b.WriteString("<style>\n")
	b.WriteString(chromeCSS)
	b.WriteString(style)
	b.WriteString("</style>\n</head>\n<body>\n<main>\n")
	fmt.Fprintf(&b, "<div class=\"masthead\">\n<h1>%s</h1>\n", html.EscapeString(heading))
	b.WriteString("<button type=\"button\" class=\"mode-toggle\">Light / dark</button>\n</div>\n")
	fmt.Fprintf(&b, "<p class=\"intro\">%s</p>\n", intro)
	b.WriteString(body)
	b.WriteString("</main>\n<script>\n")
	b.WriteString("document.querySelector(\".mode-toggle\").addEventListener(\"click\", function () {\n")
	b.WriteString("  document.documentElement.classList.toggle(\"dark\");\n")
	b.WriteString("});\n")
	b.WriteString("</script>\n</body>\n</html>\n")
	return b.String()
}

// lcStr formats an APCA Lc measurement in the signed convention: positive
// dark-on-light, negative light-on-dark.
func lcStr(text, surface stdcolor.NRGBA) string {
	return fmt.Sprintf("%.1f", color.APCA(text, surface))
}

// modeHex annotates one token's value in both modes, labelled, because text
// cannot flip with a CSS class the way a painted swatch does.
func modeHex(light, dark stdcolor.NRGBA) string {
	return fmt.Sprintf("L %s · D %s", hexRGB(light), hexRGB(dark))
}

// contrastRow renders one measured text pair as a table row: the pair's APCA
// Lc in both modes, which is the whole of what the palette is gated on.
func contrastRow(b *strings.Builder, label string, lightText, lightSurface, darkText, darkSurface stdcolor.NRGBA) {
	fmt.Fprintf(b, "<tr><th scope=\"row\">%s</th><td>%s</td><td>%s</td></tr>\n",
		html.EscapeString(label),
		lcStr(lightText, lightSurface), lcStr(darkText, darkSurface))
}

// colorPageCSS is the colour page's specimen scaffolding.
const colorPageCSS = `.platform-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(11rem, 1fr));
  gap: var(--space-4);
  margin: var(--space-4) 0;
}
.platform-swatch {
  height: var(--space-12);
  border: thin solid var(--platform-separator);
  border-radius: var(--radius-sm);
}
.platform-name .annot {
  margin: var(--space-1) 0 0;
}
.ramp {
  display: grid;
  grid-template-columns: repeat(9, 1fr);
  gap: var(--space-2);
  margin: var(--space-4) 0;
}
.chip {
  height: var(--space-16);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  border: thin solid var(--platform-separator);
  font-size: var(--font-label-medium-size);
  font-weight: var(--font-label-medium-weight);
}
.pins {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-4);
  margin: var(--space-4) 0;
}
.pin {
  min-width: var(--space-24);
}
.pin .chip {
  padding: 0 var(--space-4);
}
.contrast {
  border-collapse: collapse;
  margin: var(--space-4) 0;
}
.contrast th, .contrast td {
  border: thin solid var(--platform-separator);
  padding: var(--space-1) var(--space-3);
  font-size: var(--font-body-small-size);
  line-height: var(--font-body-small-line-height);
  text-align: right;
}
.contrast th[scope="row"] {
  text-align: left;
  font-weight: var(--font-label-medium-weight);
}
.contrast thead th {
  background: var(--platform-card-fill);
  text-align: right;
}
`

// colorHTML renders foundations/color.html: the platform's colour set, one
// swatch per name, both appearances printed, and the measured contrast of
// the pairs the platform itself puts together.
//
// It is a listing and not a system. The values are read off the platform,
// so there is nothing to explain about how one was derived from another:
// what a page can usefully add is the name, the two values, and what each
// name reads at where the platform pairs it with a fill.
func colorHTML(s Snapshot) string {
	var b strings.Builder
	b.WriteString("<section>\n<h2>The platform\u2019s colour set</h2>\n")
	b.WriteString("<p class=\"intro\">AppKit\u2019s semantic colours under their own names, plus the fills the platform draws without naming one, measured. " +
		"Each swatch is painted through <code>var(--platform-&lt;name&gt;)</code>, so the toggle restyles it; the hexes beside it are printed for both schemes, labelled L and D. " +
		"A name that carries a coverage is printed and painted as <code>#rrggbbaa</code> &mdash; the platform\u2019s answer for a label, a seam, an overlay or the focus ring is a colour AT a coverage over whatever lies beneath, and the swatch shows it over the page it is on.</p>\n")
	b.WriteString("<div class=\"platform-grid\">\n")
	for _, n := range platformNames {
		fmt.Fprintf(&b, "<div class=\"platform-name\">\n<div class=\"platform-swatch\" style=\"background: var(--platform-%s)\"></div>\n", n.name)
		fmt.Fprintf(&b, "<p class=\"annot\"><code>--platform-%s</code><br>L %s &middot; D %s</p>\n</div>\n",
			n.name, hexRGBA(n.pick(s.PlatformLight)), hexRGBA(n.pick(s.PlatformDark)))
	}
	b.WriteString("</div>\n</section>\n")

	// The measured pairs: a foreground name over the fill the platform puts
	// it on, flattened first, because every one of these foregrounds is a
	// coverage and APCA can only be handed an opaque colour.
	b.WriteString("<section>\n<h2>Measured contrast</h2>\n")
	b.WriteString("<p class=\"intro\">APCA Lc is the one contrast measure (signed: negative means light-on-dark). Each foreground is flattened over the fill beneath it in encoded sRGB first, the way the platform composites it.</p>\n")
	b.WriteString("<table class=\"contrast\">\n<thead>\n")
	b.WriteString("<tr><th scope=\"col\">pair</th><th scope=\"col\">light Lc</th><th scope=\"col\">dark Lc</th></tr>\n")
	b.WriteString("</thead>\n<tbody>\n")
	for _, pair := range platformPairs {
		lightFill, darkFill := pair.fill(s.PlatformLight), pair.fill(s.PlatformDark)
		contrastRow(&b, pair.label,
			color.Flatten(pair.text(s.PlatformLight), lightFill), lightFill,
			color.Flatten(pair.text(s.PlatformDark), darkFill), darkFill)
	}
	b.WriteString("</tbody>\n</table>\n</section>\n")

	intro := "The platform\u2019s own colour set, one field per AppKit semantic name, with the values that name reports under each appearance. " +
		"Swatches are painted through the token sheet, so the toggle restyles them; " +
		"annotation values are printed for both schemes, labelled L and D."
	return page("Colour — Vibrant Gio foundations", "Colour", intro, colorPageCSS, b.String())
}

// platformPairs are the foreground/fill pairings the platform itself makes,
// which are the only pairings a listing of its names can honestly measure:
// a label on the plane it is set on, a control's text on the control's own
// fill, the text the platform names for a fill its accent paints.
var platformPairs = []struct {
	label      string
	text, fill func(tokens.PlatformColors) stdcolor.NRGBA
}{
	{"label on the window", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Label }, func(p tokens.PlatformColors) stdcolor.NRGBA { return p.WindowBackground }},
	{"secondary label on the window", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SecondaryLabel }, func(p tokens.PlatformColors) stdcolor.NRGBA { return p.WindowBackground }},
	{"label on the chrome material", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Label }, func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SidebarMaterial }},
	{"label on the card", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Label }, func(p tokens.PlatformColors) stdcolor.NRGBA { return p.CardFill }},
	{"control text on the push button", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ControlText }, func(p tokens.PlatformColors) stdcolor.NRGBA { return p.PushButtonFill }},
	{"alternate selected control text on the accent", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.AlternateSelectedControlText }, func(p tokens.PlatformColors) stdcolor.NRGBA { return p.ControlAccent }},
	{"alternate selected control text on the selected row", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.AlternateSelectedControlText }, func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SelectedContentBackground }},
	{"alternate selected control text on the sidebar pill", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.AlternateSelectedControlText }, func(p tokens.PlatformColors) stdcolor.NRGBA { return p.SidebarSelection }},
	{"link on the window", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.Link }, func(p tokens.PlatformColors) stdcolor.NRGBA { return p.WindowBackground }},
	{"disabled control text on the push button", func(p tokens.PlatformColors) stdcolor.NRGBA { return p.DisabledControlText }, func(p tokens.PlatformColors) stdcolor.NRGBA { return p.PushButtonFill }},
}

// typePageCSS is the type page's specimen scaffolding.
const typePageCSS = `.type-role {
  margin: var(--space-8) 0;
  border-bottom: thin solid var(--platform-separator);
  padding-bottom: var(--space-4);
}
.specimen {
  margin: 0;
}
`

// typeHTML renders foundations/type.html: the fifteen type roles plus the
// code style at their real size, weight, line height and tracking, each styled
// entirely through its --font-<role>-* vars and annotated with the numbers.
// The code specimen additionally names its own family through
// var(--font-family-code) — the one role that does not inherit the body face.
func typeHTML(s Snapshot) string {
	var b strings.Builder
	for _, role := range typeRoles {
		style := role.pick(s.Typography)
		family := ""
		specimen := "The five boxing wizards jump quickly"
		if role.name == "code" {
			family = "font-family: var(--font-family-code), monospace; "
			specimen = "if mono[0] != prose { align() }"
		}
		fmt.Fprintf(&b, "<section class=\"type-role\">\n")
		fmt.Fprintf(&b, "<p class=\"specimen\" style=\"%[2]sfont-size: var(--font-%[1]s-size); line-height: var(--font-%[1]s-line-height); font-weight: var(--font-%[1]s-weight); letter-spacing: var(--font-%[1]s-tracking)\">%[3]s</p>\n",
			role.name, family, html.EscapeString(specimen))
		fmt.Fprintf(&b, "<p class=\"annot\">%s &middot; %s / %s &middot; weight %d &middot; tracking %s</p>\n",
			role.name, px(style.Size), px(style.LineHeight), style.Weight, px(style.Tracking))
		b.WriteString("</section>\n")
	}
	intro := fmt.Sprintf("Every type role at its real size, weight, line height and tracking, styled through the "+
		"<code>--font-&lt;role&gt;-*</code> tokens. The face is the family the tokens name &mdash; %s, and %s for the "+
		"code style &mdash; via <code>var(--font-family)</code> and <code>var(--font-family-code)</code>; the browser "+
		"must have them installed, otherwise the system fallback face renders at the same metrics.",
		html.EscapeString(s.Typography.BodyLarge.Typeface), html.EscapeString(s.Typography.Code.Typeface))
	return page("Type — Vibrant Gio foundations", "Type", intro, typePageCSS, b.String())
}

// layoutPageCSS is the layout page's specimen scaffolding.
const layoutPageCSS = `.space-row {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  margin: var(--space-2) 0;
}
.space-row .annot {
  flex: 0 0 10rem;
  margin: 0;
}
.space-bar {
  height: var(--space-4);
  background: var(--platform-control-accent);
  border-radius: var(--radius-sm);
}
.specimen-grid {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-8);
  margin: var(--space-4) 0;
}
.radius-box {
  width: var(--space-24);
  height: var(--space-24);
  background: var(--platform-card-fill);
  border: thin solid var(--platform-field-edge);
  display: flex;
  align-items: center;
  justify-content: center;
}
.density-pair {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-8);
  margin: var(--space-4) 0;
}
.density-col {
  flex: 1 1 18rem;
  padding: var(--space-4);
  background: var(--platform-card-fill);
  border: thin solid var(--platform-separator);
  border-radius: var(--radius-md);
}
.control-row {
  display: flex;
  align-items: center;
  margin: var(--space-2) 0;
}
.control-bar {
  height: var(--density-control-height);
  padding: 0 var(--density-padding-x);
  display: inline-flex;
  align-items: center;
  background: var(--platform-control-accent);
  color: var(--platform-alternate-selected-control-text);
  border-radius: var(--radius-md);
  font-size: var(--font-label-large-size);
  font-weight: var(--font-label-large-weight);
}
.chip-bar {
  height: var(--density-chip-height);
  padding: 0 var(--space-3);
  margin-left: var(--space-2);
  display: inline-flex;
  align-items: center;
  background: var(--platform-card-fill);
  color: var(--platform-label);
  border: thin solid var(--platform-field-edge);
  border-radius: var(--radius-md);
  font-size: var(--font-label-large-size);
  font-weight: var(--font-label-large-weight);
}
.field-bar {
  height: var(--density-field-height);
  padding: 0 var(--density-padding-x);
  display: inline-flex;
  align-items: center;
  background: var(--platform-text-background);
  color: var(--platform-placeholder-text);
  border: thin solid var(--platform-field-edge);
  border-radius: var(--radius-md);
  font-size: var(--font-body-large-size);
}
.row-bar {
  height: var(--density-row-height);
  padding: 0 var(--space-2);
  display: flex;
  align-items: center;
  background: var(--platform-alternating-content-background);
  color: var(--platform-label);
  font-size: var(--font-body-large-size);
}
.toolbar-bar {
  height: var(--density-toolbar-control-height);
  padding: 0 var(--space-3);
  display: inline-flex;
  align-items: center;
  background: var(--platform-toolbar-control-fill);
  color: var(--platform-control-text);
  border-radius: calc(var(--density-toolbar-control-height) / 2);
  font-size: var(--font-label-large-size);
  font-weight: var(--font-label-large-weight);
}
.pad-box {
  display: inline-block;
  padding: var(--density-padding-y) var(--density-padding-x);
  background: var(--platform-card-fill);
  border: thin solid var(--platform-field-edge);
  border-radius: var(--radius-md);
  margin: var(--space-2) 0;
}
.elevation-row {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-8);
  padding: var(--space-6);
  background: var(--platform-window-background);
  border-radius: var(--radius-md);
  margin: var(--space-4) 0;
}
.surface-card {
  width: var(--space-24);
  height: var(--space-20);
  border-radius: var(--radius-md);
  border: thin solid var(--platform-separator);
  display: flex;
  align-items: center;
  justify-content: center;
}
`

// layoutHTML renders foundations/layout.html: the spacing scale as sized
// bars, the control metrics at both density settings side by side, the
// radius scale on sample boxes, and the shadow depth per level — the cue a
// floating transient carries over the platform fill its level is given —
// each specimen sized, padded, rounded, filled or shadowed by its var.
func layoutHTML(s Snapshot) string {
	var b strings.Builder

	b.WriteString("<section>\n<h2>Spacing</h2>\n")
	for _, key := range spaceKeys {
		fmt.Fprintf(&b, "<div class=\"space-row\">\n<p class=\"annot\"><code>--space-%s</code> &middot; %s</p>\n", key.name, px(key.pick(s.Spacing)))
		fmt.Fprintf(&b, "<div class=\"space-bar\" style=\"width: var(--space-%s)\"></div>\n</div>\n", key.name)
	}
	b.WriteString("</section>\n")

	b.WriteString("<section>\n<h2>Density</h2>\n")
	b.WriteString("<p class=\"intro\">Two published settings share one variable family: <code>:root</code> carries comfortable, " +
		"and a <code>.compact</code> class block overrides every per-setting metric of the <code>--density-*</code> family " +
		"the way <code>.dark</code> overrides the colours &mdash; the right column below is the same markup inside a <code>class=\"compact\"</code> wrapper. " +
		"A control's pointer target is the control: what it draws at <code>--density-control-height</code> is what a pointer lands on, " +
		"so compact shrinks the target with the pixels.</p>\n")
	b.WriteString("<div class=\"density-pair\">\n")
	for _, setting := range []struct {
		class string
		label string
		d     tokens.Density
	}{
		{"density-col", "comfortable (root)", tokens.Comfortable},
		{"density-col compact", "compact (.compact)", tokens.Compact},
	} {
		fmt.Fprintf(&b, "<div class=\"%s\">\n<h3>%s</h3>\n", setting.class, html.EscapeString(setting.label))
		b.WriteString("<div class=\"control-row\">\n<span class=\"control-bar\">Control</span>\n<span class=\"chip-bar\">Chip</span>\n</div>\n")
		fmt.Fprintf(&b, "<p class=\"annot\"><code>--density-control-height</code> &middot; %s &middot; <code>--density-chip-height</code> &middot; %s &mdash; each control's own pointer target</p>\n",
			px(setting.d.ControlHeight), px(setting.d.ChipHeight()))
		b.WriteString("<div><span class=\"field-bar\">Text field</span></div>\n")
		fmt.Fprintf(&b, "<p class=\"annot\"><code>--density-field-height</code> &middot; %s &mdash; the platform draws a field taller than a button</p>\n",
			px(setting.d.FieldHeight))
		b.WriteString("<div class=\"row-bar\">Stacked row</div>\n")
		fmt.Fprintf(&b, "<p class=\"annot\"><code>--density-row-height</code> &middot; %s &mdash; a pin, not a floor: rows tile</p>\n",
			px(setting.d.RowHeight))
		b.WriteString("<div><span class=\"toolbar-bar\">Toolbar control</span></div>\n")
		fmt.Fprintf(&b, "<p class=\"annot\"><code>--density-toolbar-control-height</code> &middot; %s &mdash; a control standing in a toolbar band is its own control, taller than the one in a dialog</p>\n",
			px(setting.d.ToolbarControlHeight))
		b.WriteString("<div class=\"pad-box\">padding</div>\n")
		fmt.Fprintf(&b, "<p class=\"annot\"><code>--density-padding-x</code> %s &middot; <code>--density-padding-y</code> %s</p>\n",
			px(setting.d.PaddingX), px(setting.d.PaddingY))
		b.WriteString("</div>\n")
	}
	b.WriteString("</div>\n</section>\n")

	b.WriteString("<section>\n<h2>Radius</h2>\n<div class=\"specimen-grid\">\n")
	for _, key := range radiusKeys {
		fmt.Fprintf(&b, "<div>\n<div class=\"radius-box\" style=\"border-radius: var(--radius-%s)\">Aa</div>\n", key.name)
		fmt.Fprintf(&b, "<p class=\"annot\"><code>--radius-%s</code> &middot; %s</p>\n</div>\n", key.name, px(key.pick(s.Radius)))
	}
	b.WriteString("</div>\n</section>\n")

	b.WriteString("<section>\n<h2>Shadow depth</h2>\n")
	b.WriteString("<p class=\"intro\">A level\u2019s FILL is the platform\u2019s own name for what that region is &mdash; " +
		"the window\u2019s plane, the chrome material, the content, the card. What the levels state here is the " +
		"shadow each casts, the opt-in cue a floating transient carries on top of its fill: " +
		"<code>--shadow-N</code>, for menus, dialogs and tooltips. Neither level under the content casts anything &mdash; " +
		"the backdrop is what everything stands on, and chrome lies flat on it. " +
		"The levels stop at 3 &mdash; desktop has no six-deep stack.</p>\n")
	b.WriteString("<div class=\"elevation-row\">\n")
	for i, level := range shadowLevels {
		wears := []string{"nothing: it shows where nothing stands", "navbars, toolbars, sidebars, panes",
			"the content", "cards, fences, fields", "dialogs, toasts", "menus, popovers"}[i]
		fmt.Fprintf(&b, "<div>\n<div class=\"surface-card\" style=\"background: var(--platform-card-fill); box-shadow: var(--shadow-%s)\">%s</div>\n", level.name, level.name)
		fmt.Fprintf(&b, "<p class=\"annot\"><code>--shadow-%s</code> &middot; depth %sdp &middot; %s</p>\n</div>\n", level.name, fnum(level.dp(s.Elevation)), wears)
	}
	b.WriteString("</div>\n</section>\n")

	intro := "The spacing scale as sized bars, the control metrics at both density settings, the radius scale on sample boxes " +
		"and the shadow depth each level casts &mdash; every bar width, control height, padding, corner radius and " +
		"shadow resolves through its token, so the sheet is the single source of these shapes."
	return page("Layout — Vibrant Gio foundations", "Layout", intro, layoutPageCSS, b.String())
}
