// Package export serialises a theme.Theme emission into the project layout
// claude.ai/design consumes: theme.json, the machine-readable generative
// parameters; styles.css, the token sheet; readme.md, the project's front
// door; and the foundation pages under foundations/. Capture collects the
// first emission of each Theme observable into a Snapshot; Write renders
// the whole tree into a target directory. cmd/vg-tokens is the command-line
// front door.
//
// # The token sheet
//
// styles.css carries one :root block (the light scheme plus every
// mode-invariant scale, comfortable density), one .dark override block (the
// paired dark colours only) and one .compact override block (the compact
// density metrics only). The reference project's token families are
// recorded but not its dark-mode selector, so the sheet uses class
// overrides and .compact follows the same pattern; the two switches are
// orthogonal.
//
// Colour variables are the platform's set and nothing else:
//
//   - --platform-<name>, one per field of tokens.PlatformColors, named as
//     AppKit names the colour with the humps cut at hyphens:
//     --platform-window-background, --platform-control-accent,
//     --platform-secondary-label, --platform-sidebar-material,
//     --platform-sidebar-selection. :root carries the light appearance and
//     the .dark block the dark one. A value the platform gives at a
//     coverage is emitted with its alpha, because compositing it over
//     whatever is beneath is how the platform gets one value to read right
//     on every fill.
//   - --focus-ring-width, the ring's stroke, mode-invariant.
//
// The remaining families, all emitted in :root only because they do not
// change with the scheme:
//
//   - --font-family and --font-family-code (the code style's mono family),
//     plus --font-<role>-size, -line-height, -weight and -tracking per type
//     role (display-large … body-small, and code — the mono style outside
//     the type scale, at body-medium's metrics): sizes, line heights and
//     tracking in px, weights as CSS numeric weights.
//   - --density-control-height, --density-padding-x and --density-padding-y
//     from tokens.Density: :root carries tokens.Comfortable, the .compact
//     block overrides with tokens.Compact. --density-min-hit-target is the
//     WCAG 2.5.5 pointer-target floor, emitted once and never overridden —
//     density scales the drawn control, never the clickable area.
//   - --space-<key> from tokens.SpacingScale, keys as the Go scale names
//     them (0, 1, 2, … 24), in px.
//   - --radius-<key> from tokens.RadiusScale in Tailwind naming (none, sm,
//     base, md, lg, xl, 2xl, 3xl, full), in px; Base is also theme.json's
//     base radius parameter.
//   - --shadow-<level> (backdrop, chrome, 0, 1, 2, 3): CSS box-shadow
//     approximations of the dp depths, the cue a floating transient (a
//     menu, a dialog, a tooltip) carries over its platform fill. What each
//     level is filled with is a --platform-* name, so there is no fill
//     variable of its own. The backdrop, the chrome level and level 0 cast
//     nothing. Each level's dp depth d becomes
//     "0 <d>px <2d>px 0 rgba(0, 0, 0, 0.2)" — y-offset the depth, blur
//     twice it, no spread, black at 20% — and a zero depth is "none".
//   - --ease-<name> from tokens.MotionScale: the six easing presets as
//     cubic-bezier() strings (standard and emphasized families, each plain
//     / -accelerate / -decelerate).
//   - --duration-<stop> (x-fast, fast, normal, slow, x-slow): the five
//     duration stops in ms. The spring presets have no CSS
//     counterpart — springs are Go-side physics — and are serialised only
//     in theme.json's motion parameters.
//
// # The foundation pages
//
// foundations/color.html, type.html and layout.html render the scales at
// real sizes, and readme.md is the file a human or an agent reads first.
// The pages are static HTML that reads only from the emitted sheet: every
// styled colour, size, radius, shadow and font value is a var() reference
// into ../styles.css, so regenerating the sheet for another theme colour
// reflows every page. Literal token values appear only as annotation text — hexes,
// px numbers, and the measured APCA Lc of each text pair —
// printed for both modes (labelled L and D) because text cannot flip with a
// class the way painted specimens do. Each page carries a light/dark toggle
// flipping the .dark class on the root element. The page test enforces the
// no-hard-coded-values rule: no hex colours or px lengths in any style
// context, and every referenced variable declared by the sheet.
//
// # The generative parameters
//
// theme.json records what reproduces the theme: the theme colour as
// lowercase #rrggbb, the platform's whole set per appearance keyed by the
// sheet's names, the heading, body and mono faces, the base radius, the
// density model (the active setting by name, both settings' metrics and the
// invariant hit-target minimum), the shadow depth per level and the motion
// set (durations in ms, easing control points, spring presets). The
// platform's own reading of the theme colour is what the accent rows carry,
// so the file rebuilds the set exactly — the round-trip test asserts it.
// The per-role --font-*-size tokens come from Typography, the theme's only
// type source.
package export
