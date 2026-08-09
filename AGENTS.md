# AGENTS.md — spectrum

The token and theme layer of the Vibrant Gio design system, and the surface
every layer above it styles against. `tokens` holds the typed values —
ADR-007's colour ramps derived from a brand seed, the MD3 type roles, and
the 4-pt spacing, radius, elevation and motion scales; `theme` publishes
them as one `rx.Observable` per category, so a consumer rebuilds only for
what it reads; `color` is the perceptual mathematics those palettes are
derived with, CIELAB tone and OKLCh hue and chroma, gamut mapped rather
than clamped. Around that core: `system` polls the OS light/dark preference
and accent colour, `a11y` polls its reduce-motion, contrast and text-size
preferences, `window` pairs a theme observable with an mvu window,
`preferences` persists the user's explicit choice and republishes it as an
`mvu/stream.Value` — ADR-008's third destination, one current-value stream
per path, conflating rather than queueing so a stalled observer can never
wedge a save, `typeset` lays a type role's text out in the line box the
role names rather than the one its glyphs happen to ink, and `export` —
with `cmd/vg-tokens` in front of it — writes a theme out as the project
layout `claude.ai/design` consumes. Interpolating between two themes is not
here; it is a layer up, in `pulse/transition`.

**Layer.** Tier 1 of ADR-001's stack, `mvu → spectrum → prism → pulse →
cadence → markdown`. The token, theme and `a11y` contract the rest of the
system styles against lives here rather than in prism: goals G-B3 and E3.2
moved it down, which is what makes a tier-1 spectrum possible at all. Its
root module imports `font` and `mvu`. Imported by `cadence`, `markdown`,
`prism` and `pulse`. Outside the tier table, also by the demo modules
`mvu/example` and `prism/gallery` and all seven workbench applications.
Both directions are measured rather than typed — `scripts/check-layers.sh
--edges` reports the graph and `scripts/sync-agents.sh` renders these
sentences from it — so correcting them here changes nothing.

**Read the canonical guide before you write code against this module.** It is
the organization's one agent guide — the module inventory with current tags,
the application skeleton, the MVU loop and rx semantics, typography, and the
pitfalls that are not guessable. It lives exactly once, in `vibrantgio/.github`,
and this file links it rather than copying it:

    https://raw.githubusercontent.com/vibrantgio/.github/master/llms.txt

**Module.** `github.com/vibrantgio/spectrum`, one module at the repository
root.

**Build and test.** From the repository root:

    go build ./... && go test ./...

**Two shapers, and the choice is yours to make.** `tokens.Typography` builds
both, cached in separate fields so neither can hand back the other's:

- `Shaper()` — the system fallback is **on**. This is what applications and
  library components draw with. The embedded faces answer first; the platform's
  fonts answer for everything they lack, which is every arrow, box-drawing
  character and dingbat, because Roboto and Roboto Mono carry none of them.
  Never disable system fonts here to make an output stable — that is the F4.2
  defect, and it ships tofu to every user.
- `DeterministicShaper()` — system fonts **off**, the collection pinned. This
  is what a golden or pixel test draws with, and the reason it exists: a test
  that says which faces it wants cannot drift when the default changes.

Widen the collection with `WithFaces`, which copies and clears both caches:
`tokens.DefaultTypography.WithFaces(notosansmono.FontFace())`. That is how a
test that legitimately draws an arrow stays deterministic, and how an
application that cannot rely on system fonts — a container, a kiosk — ships its
own symbol coverage. The face is optional and is not in
`DefaultTypography.Faces`; see ADR-003.

**Line height means the line box, and `typeset` is how.**
`tokens.TextStyle.LineHeight` is the CSS thing — the height of one line box,
leading split evenly around the ink — and `gioui.org/widget.Label` does not
deliver it. Gio gives the first line its own ascent plus descent and spends the
line height only on the gap to the next, so a `MaxLines: 1` label measures the
same at every line height there is. `spectrum/typeset` wraps `widget.Label` and
adds the missing leading; every component in the org that draws a role's text
lays it out through `typeset.Layout`. `spectrum/export` writes the same number
into `--font-<role>-line-height`, so the two surfaces state one fact.

The consequence for sizing: `Density.ControlHeight` is a **floor**, not a
height. A control draws `max(ControlHeight, lineBox + 2×PaddingY)`, so a
Comfortable text field in BodyLarge is 40 dp against a 36 dp floor while a
Comfortable button in LabelLarge is exactly 36.

**Golden images.** None, and the absence above is that fact rather than an
omission: `sync-agents.sh` renders a Golden images section only for a clone
that has a `testdata/golden/` directory, and `find . -type d -name golden`
here finds none. Spectrum stores no rendered output — it computes colour, type
and layout values and asserts on numbers. The harness the repositories that do
render share is `prism/golden`.
