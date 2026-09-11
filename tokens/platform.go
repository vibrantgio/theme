// The platform's colour set: AppKit's semantic colours, one field per name.
//
// The values were read off AppKit on macOS 26.5.2 (build 25F84) on
// 2026-09-10 by a command-line program — no application was launched — under
// the aqua and darkAqua appearances, converted to sRGB, alpha kept. That
// catalogue and the program that produced it are stored in the
// organization's macOS reference; a copy of the catalogue is this package's
// testdata, and a test pins both recorded sets against it, field by field.
//
// Alpha is part of the platform's answer. A label, a seam, a disabled
// control's text is black or white at a coverage and composites over
// whatever lies beneath it, which is how one recorded value reads correctly
// on every fill. Nothing here is pre-composited: a caller that paints one of
// these over something else composites it there.
//
// The catalogue records alpha to two decimals, so each field stores
// round(a × 255) of what was written down — at most one 255th away from what
// AppKit reported. On macOS the live reader in theme/system removes that
// gap: it asks AppKit for every one of these names itself, and these
// recorded sets are what the other platforms — and every test — read.
//
// Ten fields are not AppKit's. The platform paints a sidebar, a grouped
// box, a push button, a hovered and a pressed control, the shadow under a
// floating pane, a text field's hairline, an overlay scrollbar's knob, a
// list's alternating row and the dim a sheet lays over the window it
// interrupts without giving any of them an NSColor name, so those fills
// were read off the stored captures in the organization's macOS
// reference — the alternating row off the array AppKit answers with
// instead of a name, and confirmed against the capture. Each
// carries the struct tag `appkit:"-"`, which is the whole of the rule the
// live reader and the tests use: a field so tagged has no name to ask
// AppKit for and is pinned against the catalogue's "measured materials"
// section instead of its AppKit rows.
//
//	field            light          dark           provenance
//	-----            -----          ----           ----------
//	SidebarMaterial  #ffffff        #232a2e        measured: the chrome band and the list below it
//	CardFill         #f7f7f7        #2a3034        measured: the System Settings grouped box
//	PushButtonFill   #ececec        #333a3f        measured: the Save dialog's push button at rest
//	HoverOverlay     #000000 a0.051 #ffffff a0.094 measured: a toolbar button under the pointer
//	PressOverlay     #000000 a0.098 #ffffff a0.098 measured: a push button held down
//	FloatingShadow   #000000 a0.075 #000000 a0.075 measured: the sidebar shadow's peak coverage, 24 px of reach
//	FieldEdge        #f3f3f3        #2c3338        measured: the unfocused text field's hairline in the Save dialog
//	ScrollbarThumb   #000000 a0.337 #ffffff a0.337 measured: the overlay scrollbar's knob over its track
//	AlternatingContentBackground  #f4f5f5  #ffffff a0.05  measured: the second of alternatingContentBackgroundColors, and Finder's list stripes
//	Scrim            #000000 a0.20  #000000 a0.26  measured: the dim under a Save sheet in save-dialog-{light,dark}.png
//
// Every one of these was read on a desktop whose "Tint window background
// with wallpaper colour" is on, which is what carries the chrome material
// and the grouped box off neutral grey in the dark scheme.
package tokens

import (
	"image/color"
	"math"
)

// PlatformColors is the platform's colour set: one field per AppKit
// semantic colour name, in Go casing, with the value that name reports for
// one appearance. [PlatformLight] and [PlatformDark] are the two recorded
// sets; [PlatformColors.WithAccent] substitutes the rows the platform
// derives from the accent colour.
//
// Field names drop AppKit's trailing "Color" and nothing else, so
// windowBackgroundColor is WindowBackground and systemRed is SystemRed. The
// one field whose value is not AppKit's own is FindHighlight; see its
// comment.
type PlatformColors struct {
	// The planes. WindowBackground is the window's own; ControlBackground
	// and TextBackground are the content, the lists, the tables and the
	// fields; UnderPageBackground is the backdrop showing around a pane.
	WindowBackground    color.NRGBA
	UnderPageBackground color.NRGBA
	ControlBackground   color.NRGBA
	TextBackground      color.NRGBA

	// Selection. The emphasized rows are the active window's and follow
	// the accent; the unemphasized pair is the inactive window's grey and
	// does not.
	SelectedContentBackground             color.NRGBA
	UnemphasizedSelectedContentBackground color.NRGBA
	SelectedTextBackground                color.NRGBA
	UnemphasizedSelectedTextBackground    color.NRGBA

	// FindHighlight is the find highlight as Mail paints it, measured off
	// the stored captures of Mail's find bar in both appearances. It is
	// deliberately not AppKit's findHighlightColor, which reports a
	// saturated yellow the platform's own applications do not paint. The
	// text the highlight covers keeps its colour.
	FindHighlight color.NRGBA

	// Separator is every seam, laid over whatever is beneath it. Grid is
	// the line a table rules its cells with.
	Separator color.NRGBA
	Grid      color.NRGBA

	// Text, at four strengths, over whatever it stands on.
	Label           color.NRGBA
	SecondaryLabel  color.NRGBA
	TertiaryLabel   color.NRGBA
	QuaternaryLabel color.NRGBA

	// The text-system colours: a text view's own, its placeholder, the
	// text under a selection, a link, and a header row's.
	Text            color.NRGBA
	PlaceholderText color.NRGBA
	SelectedText    color.NRGBA
	Link            color.NRGBA
	HeaderText      color.NRGBA

	// Controls: an ordinary control's fill and text, the disabled text,
	// the selected control's fill and text, and the text that reads on a
	// fill the accent paints.
	Control                      color.NRGBA
	ControlText                  color.NRGBA
	DisabledControlText          color.NRGBA
	SelectedControl              color.NRGBA
	SelectedControlText          color.NRGBA
	AlternateSelectedControlText color.NRGBA

	// The accent itself, and the ring a focused control wears.
	ControlAccent          color.NRGBA
	KeyboardFocusIndicator color.NRGBA

	// The system colours, available by name. Error, Success, Warning and
	// Info are SystemRed, SystemGreen, SystemOrange and SystemBlue.
	SystemRed    color.NRGBA
	SystemOrange color.NRGBA
	SystemYellow color.NRGBA
	SystemGreen  color.NRGBA
	SystemMint   color.NRGBA
	SystemTeal   color.NRGBA
	SystemCyan   color.NRGBA
	SystemBlue   color.NRGBA
	SystemIndigo color.NRGBA
	SystemPurple color.NRGBA
	SystemPink   color.NRGBA
	SystemBrown  color.NRGBA
	SystemGray   color.NRGBA

	// Shadow is what a shadow is drawn in; Highlight is the light edge a
	// bezel catches.
	Shadow    color.NRGBA
	Highlight color.NRGBA

	// SidebarMaterial is the chrome regions' fill: sidebars, toolbars,
	// navbars, inspectors, status bars. Dark is #232a2e, a flat-region
	// sample of mail-window.png — the toolbar band at y 1–31 and the
	// flush mailbox list below the band's hairline carry the same value,
	// and finder-window.png carries it again as the window's own plane
	// behind the floating pane. Light is #ffffff, the same regions read
	// off mail-window-light.png and finder-window-light.png; System
	// Settings' own plane agrees in both schemes. So on macOS 26 the
	// light chrome is the content's fill exactly, and a consumer that
	// needs the chrome to read apart from the content in the light
	// scheme cannot get that separation from this value.
	SidebarMaterial color.NRGBA `appkit:"-"`

	// CardFill is the fill of the platform's box — a card, a grouped
	// box, a filled inset: #f7f7f7 over a #ffffff plane light, #2a3034
	// over a #232a2e plane dark, flat-region samples of the Appearance
	// pane's boxes in system-settings-grouped-box-light.png and
	// -dark.png. The box carries no hairline and no shadow: its edge is
	// a 2–3 px antialiased ramp straight from the plane to the fill, and
	// the step across it is about 3% of the way to black in light and to
	// white in dark.
	CardFill color.NRGBA `appkit:"-"`

	// PushButtonFill is what an ordinary push button is actually filled
	// with: #ececec light and #333a3f dark, flat-region samples of the
	// "Cancel" button in save-dialog-light.png and save-dialog-dark.png
	// (x 362-430, y 504-521), the button the sheet does not fill with the
	// accent.
	//
	// It is not Control. AppKit's controlColor reports #ffffff light and
	// white at a quarter dark, which is the bezel's own backing rather
	// than the fill the platform draws: on the light sheet a push button
	// reads eight percent off the white behind it, and on the dark one it
	// reads lighter than the sheet by a fixed pixel rather than by a
	// coverage. So the fill is recorded as the pixel and Control is left
	// answering for what AppKit says it answers for.
	PushButtonFill color.NRGBA `appkit:"-"`

	// HoverOverlay and PressOverlay are what a control lays over its own
	// fill while the pointer is on it and while it is held: black in
	// light, white in dark, at the coverage that reproduces the captured
	// state over the captured resting fill. Hover is a Finder toolbar
	// button under the pointer (control-hover-light.png: #ffffff to
	// #f2f2f2; -dark.png: #242d32 to #384146); press is a Save dialog's
	// push button held down (control-pressed-light.png: #ececec to
	// #d5d5d5; -dark.png: #333a3f to #474d52). The light readings are
	// exact on every channel; the dark hover is exact on green and
	// within one 255th on red and blue, the platform's dark tint being
	// closer to an even lightening than to a white composite. A push
	// button does not tint under the pointer at all on macOS 26 and a
	// list row does not either, so a caller applies HoverOverlay only
	// where the platform does.
	HoverOverlay color.NRGBA `appkit:"-"`
	PressOverlay color.NRGBA `appkit:"-"`

	// FloatingShadow is the shadow a floating surface casts, black at
	// its peak coverage. Measured off finder-sidebar-shadow.png and
	// reminders-sidebar-shadow.png, which agree: outward from the pane's
	// 1 px edge stroke the window's #232a2e plane reads #20272b, then
	// recovers to #232a2e over 24 px. Black at 0.075 reproduces that
	// darkest pixel on every channel, so the shadow is 0.075 at the edge
	// falling to nothing 24 px out; this field carries the peak and the
	// caller spreads it.
	FloatingShadow color.NRGBA `appkit:"-"`

	// FieldEdge is the hairline a text field draws around itself,
	// unfocused: #f3f3f3 light and #2c3338 dark, the single border row of
	// the "Tags:" field in save-dialog-light.png and save-dialog-dark.png
	// (y 243 and y 269), read over that sheet's own fill — #ffffff light
	// and #232a2f dark, since the field's interior is the sheet's there.
	// The light value is Separator laid over that white exactly, to the
	// byte; the dark one is not, and nothing near it — Separator over the
	// dark sheet reads #606264, six times the step the platform draws — so
	// the edge is recorded as the pixel rather than as a coverage.
	FieldEdge color.NRGBA `appkit:"-"`

	// ScrollbarThumb is the overlay scrollbar's thumb: Label's own black
	// or white at the coverage the platform's knob measures, laid over the
	// track. The coverage is 0.337, read off textedit-scrollbar.png — the
	// knob reads #9d9fa1 over a #1a2124 track, which white at 0.337
	// reproduces within one 255th on every channel. No stored capture
	// holds a light-appearance overlay scrollbar (the platform hides the
	// overlay bar unless it is being operated, and the light window
	// captures caught none), so the light row carries the dark row's
	// coverage under Label's black until one does.
	ScrollbarThumb color.NRGBA `appkit:"-"`

	// AlternatingContentBackground is the fill a list lays under every
	// second row where the platform stripes one: #f4f5f5 light and white
	// at 0.05 dark, the second entry of AppKit's
	// alternatingContentBackgroundColors, read 2026-09-11 by the same
	// command-line program that read the catalogue. It carries no NSColor
	// name of its own — AppKit answers for the pair as an array — so it is
	// recorded here rather than asked for by name. The light value is what
	// Finder's list view draws to the byte: the stripes in
	// finder-window-light.png alternate #ffffff and #f4f5f5 on a 20 px
	// pitch. No stored capture holds a dark list view, so the dark row is
	// the array's answer alone, and it is a coverage rather than a pixel.
	AlternatingContentBackground color.NRGBA `appkit:"-"`

	// Scrim is the dim a modal lays over everything it interrupts: black
	// at 0.20 light and at 0.26 dark, measured off the window standing
	// behind the sheet in save-dialog-light.png and save-dialog-dark.png.
	// The light window's plane reads #cccccc against the #ffffff the sheet
	// itself carries, which is black at 0.20 exactly; the dark window's
	// chrome band reads #1a1f22 against the measured #232a2e chrome
	// material, which black at 0.26 reproduces on every channel (0.255 to
	// 0.275 all do; 0.25 does not).
	//
	// It is the one alpha in this set a caller may hand the rasterizer as
	// it stands: a scrim covers whatever the window happens to be showing,
	// so there is no one surface to flatten it onto.
	Scrim color.NRGBA `appkit:"-"`
}

// PlatformLight and PlatformDark are the recorded sets, aqua and darkAqua,
// with the platform's own accent — systemBlue — in the accent rows. Treat
// them as read-only, like every package-level value here: take a copy, or
// [PlatformColors.WithAccent], rather than assigning into one.
var (
	PlatformLight = PlatformColors{
		WindowBackground:    color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		UnderPageBackground: color.NRGBA{R: 0x96, G: 0x96, B: 0x96, A: 0xe6},
		ControlBackground:   color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		TextBackground:      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},

		SelectedContentBackground:             color.NRGBA{R: 0x00, G: 0x64, B: 0xe1, A: 0xff},
		UnemphasizedSelectedContentBackground: color.NRGBA{R: 0xdc, G: 0xdc, B: 0xdc, A: 0xff},
		SelectedTextBackground:                color.NRGBA{R: 0xb3, G: 0xd7, B: 0xff, A: 0xff},
		UnemphasizedSelectedTextBackground:    color.NRGBA{R: 0xdc, G: 0xdc, B: 0xdc, A: 0xff},

		FindHighlight: color.NRGBA{R: 0xfa, G: 0xef, B: 0xbd, A: 0xff},

		Separator: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x1a},
		Grid:      color.NRGBA{R: 0xe6, G: 0xe6, B: 0xe6, A: 0xff},

		Label:           color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xd9},
		SecondaryLabel:  color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x80},
		TertiaryLabel:   color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x42},
		QuaternaryLabel: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x1a},

		Text:            color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		PlaceholderText: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x80},
		SelectedText:    color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		Link:            color.NRGBA{R: 0x00, G: 0x68, B: 0xda, A: 0xff},
		HeaderText:      color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xd9},

		Control:                      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		ControlText:                  color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xd9},
		DisabledControlText:          color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x40},
		SelectedControl:              color.NRGBA{R: 0xb3, G: 0xd7, B: 0xff, A: 0xff},
		SelectedControlText:          color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xd9},
		AlternateSelectedControlText: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},

		ControlAccent:          color.NRGBA{R: 0x00, G: 0x7a, B: 0xff, A: 0xff},
		KeyboardFocusIndicator: color.NRGBA{R: 0x00, G: 0x67, B: 0xf4, A: 0x80},

		SystemRed:    color.NRGBA{R: 0xff, G: 0x38, B: 0x3c, A: 0xff},
		SystemOrange: color.NRGBA{R: 0xff, G: 0x8d, B: 0x28, A: 0xff},
		SystemYellow: color.NRGBA{R: 0xff, G: 0xcc, B: 0x00, A: 0xff},
		SystemGreen:  color.NRGBA{R: 0x34, G: 0xc7, B: 0x59, A: 0xff},
		SystemMint:   color.NRGBA{R: 0x00, G: 0xc8, B: 0xb3, A: 0xff},
		SystemTeal:   color.NRGBA{R: 0x00, G: 0xc3, B: 0xd0, A: 0xff},
		SystemCyan:   color.NRGBA{R: 0x00, G: 0xc0, B: 0xe8, A: 0xff},
		SystemBlue:   color.NRGBA{R: 0x00, G: 0x88, B: 0xff, A: 0xff},
		SystemIndigo: color.NRGBA{R: 0x61, G: 0x55, B: 0xf5, A: 0xff},
		SystemPurple: color.NRGBA{R: 0xcb, G: 0x30, B: 0xe0, A: 0xff},
		SystemPink:   color.NRGBA{R: 0xff, G: 0x2d, B: 0x55, A: 0xff},
		SystemBrown:  color.NRGBA{R: 0xac, G: 0x7f, B: 0x5e, A: 0xff},
		SystemGray:   color.NRGBA{R: 0x8e, G: 0x8e, B: 0x93, A: 0xff},

		Shadow:    color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		Highlight: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},

		SidebarMaterial: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		CardFill:        color.NRGBA{R: 0xf7, G: 0xf7, B: 0xf7, A: 0xff},
		PushButtonFill:  color.NRGBA{R: 0xec, G: 0xec, B: 0xec, A: 0xff},
		HoverOverlay:    color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x0d},
		PressOverlay:    color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x19},
		FloatingShadow:  color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x13},
		FieldEdge:       color.NRGBA{R: 0xf3, G: 0xf3, B: 0xf3, A: 0xff},
		ScrollbarThumb:  color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x56},

		AlternatingContentBackground: color.NRGBA{R: 0xf4, G: 0xf5, B: 0xf5, A: 0xff},
		Scrim:                        color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x33},
	}

	PlatformDark = PlatformColors{
		WindowBackground:    color.NRGBA{R: 0x1e, G: 0x1e, B: 0x1e, A: 0xff},
		UnderPageBackground: color.NRGBA{R: 0x28, G: 0x28, B: 0x28, A: 0xff},
		ControlBackground:   color.NRGBA{R: 0x1e, G: 0x1e, B: 0x1e, A: 0xff},
		TextBackground:      color.NRGBA{R: 0x1e, G: 0x1e, B: 0x1e, A: 0xff},

		SelectedContentBackground:             color.NRGBA{R: 0x00, G: 0x59, B: 0xd1, A: 0xff},
		UnemphasizedSelectedContentBackground: color.NRGBA{R: 0x46, G: 0x46, B: 0x46, A: 0xff},
		SelectedTextBackground:                color.NRGBA{R: 0x3f, G: 0x63, B: 0x8b, A: 0xff},
		UnemphasizedSelectedTextBackground:    color.NRGBA{R: 0x46, G: 0x46, B: 0x46, A: 0xff},

		FindHighlight: color.NRGBA{R: 0x6e, G: 0x6e, B: 0x4d, A: 0xff},

		Separator: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x1a},
		Grid:      color.NRGBA{R: 0x1a, G: 0x1a, B: 0x1a, A: 0xff},

		Label:           color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xd9},
		SecondaryLabel:  color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x8c},
		TertiaryLabel:   color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x40},
		QuaternaryLabel: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x1a},

		Text:            color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		PlaceholderText: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x8c},
		SelectedText:    color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		Link:            color.NRGBA{R: 0x41, G: 0x9c, B: 0xff, A: 0xff},
		HeaderText:      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},

		Control:                      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x40},
		ControlText:                  color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xd9},
		DisabledControlText:          color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x40},
		SelectedControl:              color.NRGBA{R: 0x3f, G: 0x63, B: 0x8b, A: 0xff},
		SelectedControlText:          color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xd9},
		AlternateSelectedControlText: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},

		ControlAccent:          color.NRGBA{R: 0x00, G: 0x7a, B: 0xff, A: 0xff},
		KeyboardFocusIndicator: color.NRGBA{R: 0x1a, G: 0xa9, B: 0xff, A: 0x80},

		SystemRed:    color.NRGBA{R: 0xff, G: 0x42, B: 0x45, A: 0xff},
		SystemOrange: color.NRGBA{R: 0xff, G: 0x92, B: 0x30, A: 0xff},
		SystemYellow: color.NRGBA{R: 0xff, G: 0xd6, B: 0x00, A: 0xff},
		SystemGreen:  color.NRGBA{R: 0x30, G: 0xd1, B: 0x58, A: 0xff},
		SystemMint:   color.NRGBA{R: 0x00, G: 0xda, B: 0xc3, A: 0xff},
		SystemTeal:   color.NRGBA{R: 0x00, G: 0xd2, B: 0xe0, A: 0xff},
		SystemCyan:   color.NRGBA{R: 0x3c, G: 0xd3, B: 0xfe, A: 0xff},
		SystemBlue:   color.NRGBA{R: 0x00, G: 0x91, B: 0xff, A: 0xff},
		SystemIndigo: color.NRGBA{R: 0x6d, G: 0x7c, B: 0xff, A: 0xff},
		SystemPurple: color.NRGBA{R: 0xdb, G: 0x34, B: 0xf2, A: 0xff},
		SystemPink:   color.NRGBA{R: 0xff, G: 0x37, B: 0x5f, A: 0xff},
		SystemBrown:  color.NRGBA{R: 0xb7, G: 0x8a, B: 0x66, A: 0xff},
		SystemGray:   color.NRGBA{R: 0x98, G: 0x98, B: 0x9d, A: 0xff},

		Shadow:    color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		Highlight: color.NRGBA{R: 0xb4, G: 0xb4, B: 0xb4, A: 0xff},

		SidebarMaterial: color.NRGBA{R: 0x23, G: 0x2a, B: 0x2e, A: 0xff},
		CardFill:        color.NRGBA{R: 0x2a, G: 0x30, B: 0x34, A: 0xff},
		PushButtonFill:  color.NRGBA{R: 0x33, G: 0x3a, B: 0x3f, A: 0xff},
		HoverOverlay:    color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x18},
		PressOverlay:    color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x19},
		FloatingShadow:  color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x13},
		FieldEdge:       color.NRGBA{R: 0x2c, G: 0x33, B: 0x38, A: 0xff},
		ScrollbarThumb:  color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x56},

		AlternatingContentBackground: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x0d},
		Scrim:                        color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x42},
	}
)

// WithAccent returns the set with the five rows the platform derives from
// the accent colour rebuilt for accent. Every other field is untouched: the
// unemphasized selection pair is the inactive window's grey, Link is a fixed
// blue, and the system colours are a fixed catalogue — none of them moves
// when the user changes the accent.
//
// The five rows and their rule:
//
//   - ControlAccent is the accent itself, at the recorded row's alpha.
//   - SelectedContentBackground, SelectedTextBackground, SelectedControl
//     and KeyboardFocusIndicator keep the accent's hue and saturation and
//     take the recorded row's HSL lightness and alpha. In the catalogue all
//     four are the platform's blue at a lightness of their own — the
//     emphasized selection darker than the accent, the text selection and
//     the selected control much lighter, the focus ring a shade darker at
//     half coverage — and that lightness is what carries over.
//
// The rule is an APPROXIMATION of what AppKit does, not a reproduction of
// it: the platform derives these in its own colour space, and rebuilding
// the catalogue's blue rows through this rule lands within a few units per
// channel rather than on the byte. It is the rule on Windows and Linux,
// whose desktops publish an accent colour but no colour set of their own;
// macOS does not use it at all, because theme/system's live reader asks
// AppKit for each of these names directly and gets the platform's own
// answer. Where accent is the platform's own blue —
// the accent already in the recorded rows — the set is returned unchanged,
// so a machine on the default accent gets the catalogue exactly.
func (p PlatformColors) WithAccent(accent color.NRGBA) PlatformColors {
	if accent.R == p.ControlAccent.R && accent.G == p.ControlAccent.G && accent.B == p.ControlAccent.B {
		return p
	}
	hue, sat, _ := hsl(accent)
	p.ControlAccent = color.NRGBA{R: accent.R, G: accent.G, B: accent.B, A: p.ControlAccent.A}
	p.SelectedContentBackground = atAccentHue(hue, sat, p.SelectedContentBackground)
	p.SelectedTextBackground = atAccentHue(hue, sat, p.SelectedTextBackground)
	p.SelectedControl = atAccentHue(hue, sat, p.SelectedControl)
	p.KeyboardFocusIndicator = atAccentHue(hue, sat, p.KeyboardFocusIndicator)
	return p
}

// atAccentHue rebuilds one recorded row at the accent's hue and saturation,
// keeping the row's own lightness and alpha.
func atAccentHue(hue, sat float64, recorded color.NRGBA) color.NRGBA {
	_, _, light := hsl(recorded)
	out := fromHSL(hue, sat, light)
	out.A = recorded.A
	return out
}

// hsl reads a colour's hue in degrees, saturation and lightness, ignoring
// alpha. Saturation and lightness run 0..1.
func hsl(c color.NRGBA) (hue, sat, light float64) {
	r, g, b := float64(c.R)/255, float64(c.G)/255, float64(c.B)/255
	high := math.Max(r, math.Max(g, b))
	low := math.Min(r, math.Min(g, b))
	light = (high + low) / 2
	span := high - low
	if span == 0 {
		return 0, 0, light
	}
	sat = span / (1 - math.Abs(2*light-1))
	if sat > 1 {
		sat = 1
	}
	switch high {
	case r:
		hue = math.Mod((g-b)/span, 6)
	case g:
		hue = (b-r)/span + 2
	default:
		hue = (r-g)/span + 4
	}
	hue *= 60
	if hue < 0 {
		hue += 360
	}
	return hue, sat, light
}

// fromHSL is hsl's inverse, returning an opaque colour.
func fromHSL(hue, sat, light float64) color.NRGBA {
	span := (1 - math.Abs(2*light-1)) * sat
	sector := hue / 60
	mid := span * (1 - math.Abs(math.Mod(sector, 2)-1))
	var r, g, b float64
	switch {
	case sector < 1:
		r, g, b = span, mid, 0
	case sector < 2:
		r, g, b = mid, span, 0
	case sector < 3:
		r, g, b = 0, span, mid
	case sector < 4:
		r, g, b = 0, mid, span
	case sector < 5:
		r, g, b = mid, 0, span
	default:
		r, g, b = span, 0, mid
	}
	base := light - span/2
	return color.NRGBA{R: eightBit(r + base), G: eightBit(g + base), B: eightBit(b + base), A: 0xff}
}

// eightBit quantizes a 0..1 channel to a byte, clamping out-of-range input.
func eightBit(v float64) uint8 {
	n := math.Round(v * 255)
	if n < 0 {
		n = 0
	}
	if n > 255 {
		n = 255
	}
	return uint8(n)
}
