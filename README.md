# theme

The theme runtime of [Vibrant Gio](https://github.com/vibrantgio), a design
system for native desktop applications on macOS, Windows and Linux, written in
pure Go on [Gio](https://gioui.org). theme is the layer that answers one
question — *what does this window look like right now* — and answers it as a
stream, so the answer can change while the application runs.

Following the operating system between light and dark is the kind of thing that
is easy to demo and tedious to actually do: something has to poll the OS,
notice a real change rather than re-emitting the same value, turn it into
design tokens, and get those tokens to every component on screen without the
application threading a parameter through its whole view tree. theme does
that in one line at startup. `system.LiveTheme` publishes the OS appearance as
an `rx.Observable[theme.Theme]`; `window.New` binds that observable to an
[mvu](https://github.com/vibrantgio/mvu) window and hands it to the builder
that constructs the layers. A component that takes the theme observable as
its first argument receives the appearance change with no application code
at all, which is why an application's bootstrap is the same two lines, and
none of them asks the OS about appearance a second time. The same stream
carries the OS accent colour — an accent change re-emits the theme just like
a dark-mode flip — and while the OS reports increased contrast, there is no
branch at all: the platform paints its own semantic colours differently
under "Increase Contrast", and a live read carries what it paints.

The module is small and testable without a display. Below `window`, two
packages import Gio themselves: `tokens`, whose `Typography` owns the
system's shapers and therefore imports Gio's text machinery, and `typeset`,
which lays text out with Gio's label. `system`, `a11y`, `preferences`,
`export`, `color`, `brand`, `imagecolor`, `system/naming` and
`system/openpanel` talk to the OS, the filesystem and the mathematics and
import no Gio package directly; those that use `tokens` reach Gio's text
packages through it.

`Typography` builds two shapers, and which one you take is a decision, not a
detail. `Shaper()` is what applications draw with: the embedded faces first,
then the platform's own fonts for anything they cannot serve, so text
resolves — all of it, including the arrows, box-drawing characters and
symbols Roboto and Roboto Mono do not carry. `DeterministicShaper()` is what
golden tests draw with: those faces and nothing else, system fonts off,
identical on every machine. They are cached apart, so neither can hand back
the other's.

Determinism belongs to the test that wants it. A test that needs a symbol
adds the face that carries it rather than reaching for the platform's:

```go
typ := tokens.DefaultTypography.WithFaces(notosansmono.FontFace())
shaper := typ.DeterministicShaper()
```

The same one-liner is how an application that cannot rely on system fonts — a
container, a kiosk — gets symbol coverage; see
[font](https://github.com/vibrantgio/font)'s `notosansmono`, which is optional
and not in `DefaultTypography.Faces`.

## Where it sits

The foundation tier of the stack: the module that owns the design values
everything above it is styled from. theme imports
[mvu](https://github.com/vibrantgio/mvu) and
[font](https://github.com/vibrantgio/font) — Roboto and Roboto Mono are the
default `Typography`'s faces — and nothing above it, which the layer check
records as no upward edge at all. The
[organization page](https://github.com/vibrantgio) has the full tier table.

```sh
go get github.com/vibrantgio/theme
```

## Packages

| Package | |
| --- | --- |
| `tokens` | The typed design values, all of them: `PlatformColors`, the platform's own colour set — one field per AppKit semantic colour name, with `PlatformLight` and `PlatformDark` carrying the values measured from the platform, `WithAccent` rebuilding the rows the platform derives from the accent, and `FloatingShadow` and `PaneShadow` as `DropShadow` readings; `Typography` — fifteen text roles plus `Code` and `DocumentHeadings`, the six-step heading scale prose surfaces set their headings in, carrying the face collection, `WithFaces` to widen it, `WithCodeFace` to name the code face, and two shapers cached apart: `Shaper()` with the system fallback for applications and `DeterministicShaper()` with the collection pinned for golden tests; `EmojiTypography()` is `DefaultTypography.WithEmoji()`, the typography a live theme stream emits; `Density` — `Comfortable` and `Compact`, whose minimum control heights are 24 dp and 19 dp: a control draws taller when its type role's line box and padding need it, so a Compact button draws 20 dp, and the text field, the list row, the checkbox row and the toolbar control each have a height of their own; `MotionScale` (duration stops, easings, spring presets, and `Reduced()` for the OS reduce-motion preference); and the 4-pt spacing and named radius scales. |
| `color` | The colour mathematics the theme measures and composites with: OKLab and OKLCh conversions and their inverses, gamut mapping that reduces chroma rather than clamping channels, the APCA contrast measure (`APCA`, `Magnitude`, `BestOn`), relative luminance, and the compositing helpers `Flatten`, `Fade` and `LinearCoverage`. No colour values live here. |
| `theme` | `Theme`: one `rx.Observable` per token category, so a reader subscribes to just the categories it reads. `Default()` and `AutoLightDark()` construct one — note `AutoLightDark()` reads the clock (hours 7–17 light), not the OS; `system.LiveTheme` is the real tracker. |
| `system` | The OS appearance — dark mode and accent colour — polled behind a `Source` interface and published as an observable that emits only on change. `Live` gives the raw `Appearance`; `LiveTheme` gives the `theme.Theme` a window wants, with the options `WithThemeColor` for the accent an application chooses, `WithA11ySource` for the accessibility preferences it reads and `WithTypography` for the type roles it emits (`tokens.EmojiTypography()` by default). `PlatformColor` is the colour a stream with no theme colour chosen rebuilds the accent rows for, for a given `Appearance`. Dark mode is read on macOS; the accent is read on all three platforms — macOS's accent choice, the Windows DWM registry value, GNOME's named accent and KDE's `kdeglobals` RGB. |
| `system/naming` | Orders names the way the platform's file browser orders them: case-insensitively, with a run of digits read as the number it spells. `Compare` and `Less`; on macOS the answer is the platform's own comparison. |
| `system/openpanel` | The platform's own chooser for a folder: `Available` says whether the platform has one, `ChooseDirectory` presents it and returns what was chosen. macOS with cgo has an open panel; elsewhere `Available` is false. |
| `a11y` | OS accessibility preferences — reduce motion, increased contrast, larger text — polled and published as an `rx.Observable[A11yPrefs]` that emits only on change. macOS and Windows report real preferences, Linux returns all-false. |
| `window` | Pairs an `mvu.Window` with the theme observable that scopes it, and hands that observable to the layer builder. Two windows built with two theme streams render in two different themes in the same process. |
| `typeset` | Lays a type role's text out in the line box that role names, rather than in the box its glyphs occupy: `Layout` wraps Gio's label and pads the result up to the line box, `Text`, `TextWeight` and `TextAligned` set a role's text in one call, `Label` and `Font` build the label and font from a `TextStyle`, and `LineWidth` measures a line. |
| `brand` | Keeps the brand colour a person chose — the theme colour, or that it follows the system, with a pair of code styles for light and dark, a code face and where the colour came from — as one JSON file, and hands it back as `LiveTheme` options (`Brand.Options`). `Load`, `Save`, `Kept` and their `From`/`To` forms read and write it. |
| `imagecolor` | Extracts theme-colour candidates from an image: colours the picture is made of, clustered in OKLab and ranked so a vivid minority leads a dull majority. `Extract`, `ExtractWith`, `ExtractPalette` and `ExtractPaletteWith` return `Candidate`s, never a palette. |
| `preferences` | Persists the user's explicit appearance choice — a theme name plus accessibility overrides — as JSON under the OS config directory, and reads it back at launch. |
| `export` | Serialises a `theme.Theme` emission into the project layout claude.ai/design reads — `theme.json`, `styles.css`, `readme.md` and the foundation pages. `cmd/vg-tokens` is the command-line front door. |
| `cmd/a11y-check` | Polls the OS accessibility preferences once a second and prints them at start and on each change, to check by hand that a toggle in System Settings arrives. |

## Usage

The whole bootstrap of a Vibrant Gio application, in its smallest complete
form. Two of these lines are theme:

```go
mvuWin := mvu.NewWindow(
	app.Title("My app"),
	app.Size(unit.Dp(650), unit.Dp(600)),
)
w := specwin.New(mvuWin, specsystem.LiveTheme(time.Second))

models, runner := mvu.Loop(mvuWin.Messages(), Init, Update)
defer func() { runner.Unsubscribe(); runner.Wait() }()

if err := w.Render(buildLayers(models)).Wait(); err != nil {
	fmt.Fprintln(os.Stderr, "my app:", err)
	os.Exit(1)
}
```

`mvu.Loop` returns the model observable already multicast, carrying the model
in force to every subscriber the moment it attaches, so the layer builder
subscribes it as often as its topology needs and nothing counts subscribers.

The argument is the poll interval. `Live`'s recommended interval is 100–250
ms; one second between an OS change and its emission is the outer budget.

Options on `LiveTheme` (and `FromSourceTheme`) brand the window without giving
up live OS tracking:

```go
// one theme colour; the rest of the set stays the platform's, dark mode
// stays live
specsystem.LiveTheme(time.Second, specsystem.WithThemeColor(brand))
```

Precedence, highest first: a chosen theme colour pins the accent rows — the
application chose its colour, the OS accent is ignored. With none, macOS
uses what AppKit reports and rebuilds the three measured accent fills — the
sidebar pill, the text on the grey pill and the default button — for the
live accent. The other desktops rebuild the recorded set's accent rows for
the colour they publish. Increased contrast needs no branch
of its own: the platform reports the values it paints under it. While the OS
reports reduced motion, `Motion` emits `MotionScale.Reduced()` — every
duration zero.

`Render` is where the theme becomes the application's. It calls the build
function with this window's own theme observable and renders the layers that
come back, so the observable is a parameter rather than a global — this is
`view.go` from the same app:

```go
func buildLayers(modelObs rx.Observable[Model]) func(th rx.Observable[theme.Theme]) []rx.Observable[layout.Widget] {
	return func(th rx.Observable[theme.Theme]) []rx.Observable[layout.Widget] {
		return []rx.Observable[layout.Widget]{
			BackdropLayer(th),
			ContentLayer(th, modelObs),
		}
	}
}
```

From there `th` goes straight into the components, which take it as their
first argument. A layer that needs the resolved values rather than the
`Theme` subscribes to the category it reads — each `LiveTheme` emission is a
static snapshot, every field an `rx.Of`, so the inner observable resolves
synchronously:

```go
themes := rx.SwitchMap(th, func(t theme.Theme) rx.Observable[themed] {
	return rx.Map(t.Platform, func(c tokens.PlatformColors) themed {
		return themed{components: t, palette: PaletteFrom(c)}
	})
})
```

To test any of this without an OS, implement `system.Source` and use
`FromSource` or `FromSourceTheme`; that is the whole test seam, and it is what
this module's own tests drive.

## For coding assistants

This README's own sections are the guide: [Where it sits](#where-it-sits) and
[Packages](#packages) place the module and its packages, [Usage](#usage) is
the whole bootstrap, and each package's own `doc.go` comment covers what the
table above only summarises.

Read this repository's [AGENTS.md](./AGENTS.md) before you write code against this module.

Build and test the root module with `go build ./...` and `go test ./...`.

This module stores no golden images.

## Status

Honest about what does not work yet:

- **Dark mode is only detected on macOS.** The Linux and Windows sources read
  the *accent* — GNOME's `gsettings` enum and KDE's `kdeglobals` RGB on Linux,
  the DWM registry value on Windows — but not dark mode: `Appearance.Dark`
  stays false there. Both files name the API a real implementation would
  use — an `org.freedesktop.appearance` portal read on Linux,
  `AppsUseLightTheme` plus a registry watch on Windows — and neither is
  written.
- **The theme snaps; nothing in this module cross-fades it.** `LiveTheme`
  emits the new colour set in one step. A cross-fade is the caller's to build
  by interpolating two `PlatformColors` values.
- **`preferences` persists a choice this module does not apply.** There is no
  mapping from the stored theme name to a `theme.Theme` — the string
  round-trips to disk and stops there, as do the stored a11y overrides.
  `Observe` is a live stream — it emits the persisted value and then
  re-emits on every in-process `Save` to the same path — but writes from
  other processes are unobserved.
- **The theme streams are shared.** One `LiveTheme` (or `Live`/`FromSource`)
  value runs one poll loop per OS source no matter how many layers
  subscribe: late subscribers replay the latest value, and the loops stop
  when the last subscriber unsubscribes. Sharing is per observable value —
  build the stream once and hand the same value around; each separate
  `LiveTheme` call costs its own loops.

## License

MIT — see [LICENSE](./LICENSE).
