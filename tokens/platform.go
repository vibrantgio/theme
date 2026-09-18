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
// The catalogue records alpha as the byte AppKit reports, re-read
// 2026-09-13: every alpha the platform names is exactly n/255, so a field
// stores that n and nothing rounds. Label is 216, not the 217 a two-decimal
// 0.85 would round to. On macOS the live reader in theme/system asks AppKit
// for every one of these names itself; these recorded sets are what the
// other platforms — and every test — read.
//
// Twelve fields are not AppKit's. The platform paints a sidebar, a grouped
// box, a selected sidebar row, a search field's recess on a sidebar, a push
// button, a hovered and a pressed
// control, the shadow under a
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
//	SidebarMaterial  #f7f7f7        #1c1c1c        measured: the chrome band, wallpaper tinting off
//	SidebarSelection #178bfb        #1994fc        measured: the selected sidebar row's pill
//	CardFill         #f7f7f7        #2a3034        measured: the System Settings grouped box
//	PushButtonFill   #ececec        #333a3f        measured: the Save dialog's push button at rest
//	HoverOverlay     #000000 a0.051 #ffffff a0.094 measured: a toolbar button under the pointer
//	PressOverlay     #000000 a0.098 #ffffff a0.098 measured: a push button held down
//	FloatingShadow   #000000 a0.075 #000000 a0.075 measured: the sidebar shadow's peak coverage, 24 px of reach
//	FieldEdge        #f3f3f3        #2c3338        measured: the unfocused text field's hairline in the Save dialog
//	ScrollbarThumb   #000000 a0.572 #ffffff a0.572 measured: the overlay scrollbar's knob over its track
//	AlternatingContentBackground  #f4f5f5  #ffffff a0.05  measured: the second of alternatingContentBackgroundColors, and Finder's list stripes
//	Scrim            #000000 a0.20  #000000 a0.26  measured: the dim under a Save sheet in save-dialog-{light,dark}.png
//	SidebarSearchFill  #e8e8e8      #2f3234        measured: the recess a search field is on a sidebar
//	ToolbarControlFill #ffffff      #262626        measured: a bordered control in a frontmost window's toolbar
//	ToolbarSearchFill  #e8e8e8      #363636        measured: the recess a search field is in a toolbar
//	ToolbarControlShadow  #000000 a0.035  #000000 a0.024  measured: the drop shadow a bordered toolbar control casts, 23 px of reach 9 px below the control
//
// The chrome material and the sidebar's pill were read with "Tint window
// background with wallpaper colour" off, so they carry the platform's own
// shade and not the desktop picture's. The grouped box, the push button,
// the field edge, the scrim and the sidebar search field's dark recess were
// read with it on, which is what carries them off neutral grey in the dark
// scheme.
package tokens

import (
	"image/color"
	"math"
)

// DisabledCoverage is how much of its own paint a control keeps while it is
// switched off: the fraction of 255 its fill, its edge and its marks are
// drawn at over the surface it stands on. The control fades toward that
// surface; it is not tinted by an overlay of a fixed colour, which no single
// coverage could reproduce in both appearances.
//
// MEASURED, save-dialog-light.png and save-dialog-dark.png: the "Options:"
// checkbox is switched off in both, its 16 px box reading #f2f2f2 light and
// #2e3439 dark, on sheets of #ffffff and #232a2f — and seventeen rows above
// it the "File Format:" pop-up, enabled on the same sheet, reads the push
// button's own #ececec and #333a3f. 170 of 255 puts the enabled fill on the
// disabled reading exactly light and exactly on dark red, one 255th over on
// dark green and blue, which is the tolerance the hover overlay's dark
// reading carries in the same reference.
//
// The dialog holds no switched-off push button and no switched-off edged
// control, so this one coverage carries both the fill and the edge; a
// capture of either is on the capture list. Text is not faded by it: a
// switched-off control's wording is DisabledControlText, which is the
// platform's own reduced coverage and reads at it in the same capture.
const DisabledCoverage uint8 = 170

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
	// navbars, inspectors, status bars. #f7f7f7 light and #1c1c1c dark,
	// flat-region samples of Finder's sidebar in
	// finder-window-untinted-light.png and finder-window-untinted-dark.png,
	// read with window-background wallpaper tinting off so the value is the
	// platform's own and not the desktop picture's. The content beside it
	// reads #ffffff and #1e1e1e in the same captures, so the chrome is a
	// shade under the content in both appearances. The platform paints the
	// dark band #1b1b1b through #1d1d1d down its height; #1c1c1c is what it
	// holds over the band's flat middle, and this set paints it flat.
	SidebarMaterial color.NRGBA `appkit:"-"`

	// SidebarSelection is the pill the platform lays under the selected row
	// of a sidebar: #178bfb light and #1994fc dark, flat-region samples of
	// Voice Memos' selected row in voicememos-sidebar-light.png (x 74-273,
	// y 363-394) and voicememos-sidebar-dark.png, both 32 tall and cornered
	// at 8, with a white label.
	//
	// It is deliberately neither ControlAccent nor
	// SelectedContentBackground. The pill follows the user's accent, so
	// [PlatformColors.WithAccent] moves it; but the platform lifts it above
	// the accent's own #007aff over the chrome material, and a content
	// list's selected row wears a different colour again. So the lift is
	// recorded as the pixel in both appearances.
	SidebarSelection color.NRGBA `appkit:"-"`

	// CardFill is the fill of the platform's box — a card, a grouped
	// box, a filled inset: #f7f7f7 over a #ffffff plane light, #2a3034
	// over a #232a2e plane dark, flat-region samples of the Appearance
	// pane's boxes in system-settings-grouped-box-light.png and
	// -dark.png. The box carries no hairline and no shadow: its edge is
	// a 2–3 px antialiased blend straight from the plane to the fill, and
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
	// closer to an even lightening than to a white composite.
	//
	// The overlay is the platform's one hover answer and a caller lays it
	// on whatever fill the control carries. The reference holds no capture
	// of a push button or a list row under the pointer and not held, so
	// neither is measured as exempt; both captures are on the capture
	// list. What the pressed capture does settle is that the two states do
	// not stack: the held push button reads #d5d5d5 light and #474d52
	// dark, which is PressOverlay straight over PushButtonFill, with no
	// hover under it.
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
	// track. The coverage is 0.572 — 146/255 — read off
	// textedit-scrollbar.png: the knob reads #9d9fa1 over a #1a2124 track,
	// which white at that coverage reproduces within one 255th on every
	// channel when flattened in encoded sRGB, the space the platform
	// composites in. The 0.337 this row carried before was fitted through a
	// blend in linear light. No stored capture
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

	// SidebarSearchFill is the recess a search field standing on chrome is
	// drawn as: #e8e8e8 light and #2f3234 dark, measured off the field at
	// the top of System Settings' sidebar in
	// system-settings-grouped-box-light.png and -dark.png. The recess
	// carries no edge, its ends are fully rounded, and it measures 28 px
	// tall with 8 px of sidebar on either side of it.
	//
	// It is a fill and not a coverage over what it stands on: the same
	// #e8e8e8 stands on System Settings' #fafaf9 sidebar and on
	// Voice Memos' #ffffff toolbar band, and the direction does not survive the
	// scheme — light the recess is 18 levels darker than the material,
	// dark it is 19 levels lighter.
	//
	// The dark value was read with wallpaper tinting on, as CardFill's
	// was, over a #1c2124 sidebar; no stored capture holds a sidebar
	// search field in the dark appearance untinted.
	SidebarSearchFill color.NRGBA `appkit:"-"`

	// ToolbarControlFill is what a control standing in a toolbar is filled
	// with: #ffffff light and #262626 dark.
	//
	// The dark value is a flat-region sample of every bordered control in
	// the toolbar of finder-window-untinted-dark.png — that window is
	// frontmost, its traffic lights saturated — over a band reading
	// #1e1e1e, eight levels lighter than what it stands on. The tinted
	// capture of the same window agrees in direction: #242d32 on a #232a2e
	// band in finder-window.png. The light value is the same control in
	// finder-window-light.png, where the band beneath it is the content's
	// own #ffffff and the control is told from it by its shadow alone, so
	// what the capture fixes is the fill's value and not its step; on the
	// chrome material this set paints it is eight levels lighter, the dark
	// appearance's step to the level.
	//
	// The light reading is the frontmost window's. finder-window-untinted-light.png
	// is an inactive window — no traffic light in it carries a hue — and its
	// controls are faded: they read #f7f7f7 on a #ffffff band with their
	// glyphs at TertiaryLabel, which is the platform's inactive drawing and
	// not a control's own fill.
	ToolbarControlFill color.NRGBA `appkit:"-"`

	// ToolbarSearchFill is the recess a search field standing in a TOOLBAR
	// is drawn as: #e8e8e8 light and #363636 dark. It is a second material
	// from [PlatformColors.SidebarSearchFill] because the platform draws the
	// two recesses apart — one fill on a sidebar, another in a toolbar — and
	// from [PlatformColors.ToolbarControlFill] because a search field is not
	// the bordered control beside it.
	//
	// MEASURED, voicememos-sidebar-light.png and voicememos-window.png, one
	// capture per appearance of a frontmost Voice Memos window whose toolbar
	// carries a search field at its trailing end. Light: the field spans
	// y 46-81 and its interior is flat #e8e8e8, the same value the sidebar
	// recess carries. Dark: the field spans y 8-43 at x 643-967, its interior
	// #363636 over a #1e1e1e band, where the sidebar recess reads #2f3234 and
	// a bordered toolbar control reads #262626. Neither stored capture of the
	// sidebar recess corrects this one and neither is corrected by it, so both
	// are recorded.
	//
	// Both Voice Memos captures are untinted — every channel of the fill and
	// of the band it stands on is equal — so the dark value carries no
	// wallpaper cast, unlike SidebarSearchFill's.
	//
	// Two other stored windows draw a toolbar search field at their own
	// bordered control's fill instead — Finder's #262626 in
	// finder-window-untinted-dark.png and Mail's #242d32 in mail-window.png —
	// and neither light toolbar in those two holds a recess to read at all.
	// The application this library reads the toolbar search field from is
	// Voice Memos, which is the one the Language names.
	ToolbarSearchFill color.NRGBA `appkit:"-"`

	// ToolbarControlShadow is the peak coverage of the drop shadow a
	// BORDERED TOOLBAR CONTROL casts on the band it stands in: black at
	// 0.035 light and at 0.024 dark. It is the shadow that tells a light
	// control from its band — the platform draws no edge there and the fill
	// is the band's own white — and this field carries the peak the way
	// [PlatformColors.FloatingShadow] does, the caller spreading it over the
	// measured reach.
	//
	// The reach and the offset are ONE geometry for both appearances and are
	// spent by the drawing: 23 px of reach with the shadow's rectangle sunk
	// 9 px below the control, which is what makes it heavier under the
	// control than over it. They are not this field's to carry because they
	// are lengths and not colours.
	//
	// MEASURED, finder-window-light.png, the view pop-up at x 694-742,
	// y 8-43, on a band flat at #ffffff: the row under the control reads 244
	// and the row over it 250, the columns beside it 248, and the band
	// recovers to #ffffff 39 rows below the control, 31 columns beside it and
	// — extrapolated, the window's own top edge cutting the reading off at 8
	// rows — about 15 rows above it. A linear ramp at this coverage over that
	// reach, sunk by that offset, lands every one of the 77 sampled pixels
	// within two 255ths and most within one.
	//
	// MEASURED, finder-window-untinted-dark.png: the same shadow in the dark
	// appearance is one 255th deep and seven rows tall. Under each of the
	// five bordered controls the #1e1e1e band reads #1d1d1d over y 82-88 and
	// nothing above or beside them; finder-window.png and mail-window.png
	// agree on their tinted #232a2e band (#222a2d under the control), and
	// notes-toolbar.png agrees on its own. Black at 0.024 reproduces all four
	// of those bytes exactly. Spread over the one geometry it darkens a wider
	// halo than the platform's seven rows, which is a miss of one 255th on a
	// band this dark; the dark control is told from its band by its fill and
	// its rim, not by this.
	ToolbarControlShadow color.NRGBA `appkit:"-"`
}

// PlatformLight and PlatformDark are the recorded sets, aqua and darkAqua,
// with the platform's own accent — systemBlue — in the accent rows. Treat
// them as read-only, like every package-level value here: take a copy, or
// [PlatformColors.WithAccent], rather than assigning into one.
var (
	PlatformLight = PlatformColors{
		WindowBackground:    color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		UnderPageBackground: color.NRGBA{R: 0x96, G: 0x96, B: 0x96, A: 0xe5},
		ControlBackground:   color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		TextBackground:      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},

		SelectedContentBackground:             color.NRGBA{R: 0x00, G: 0x64, B: 0xe1, A: 0xff},
		UnemphasizedSelectedContentBackground: color.NRGBA{R: 0xdc, G: 0xdc, B: 0xdc, A: 0xff},
		SelectedTextBackground:                color.NRGBA{R: 0xb3, G: 0xd7, B: 0xff, A: 0xff},
		UnemphasizedSelectedTextBackground:    color.NRGBA{R: 0xdc, G: 0xdc, B: 0xdc, A: 0xff},

		FindHighlight: color.NRGBA{R: 0xfa, G: 0xef, B: 0xbd, A: 0xff},

		Separator: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x19},
		Grid:      color.NRGBA{R: 0xe6, G: 0xe6, B: 0xe6, A: 0xff},

		Label:           color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xd8},
		SecondaryLabel:  color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x7f},
		TertiaryLabel:   color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x42},
		QuaternaryLabel: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x19},

		Text:            color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		PlaceholderText: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x7f},
		SelectedText:    color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		Link:            color.NRGBA{R: 0x00, G: 0x68, B: 0xda, A: 0xff},
		HeaderText:      color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xd8},

		Control:                      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		ControlText:                  color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xd8},
		DisabledControlText:          color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x3f},
		SelectedControl:              color.NRGBA{R: 0xb3, G: 0xd7, B: 0xff, A: 0xff},
		SelectedControlText:          color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xd8},
		AlternateSelectedControlText: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},

		ControlAccent:          color.NRGBA{R: 0x00, G: 0x7a, B: 0xff, A: 0xff},
		KeyboardFocusIndicator: color.NRGBA{R: 0x00, G: 0x67, B: 0xf4, A: 0x7f},

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

		SidebarMaterial:  color.NRGBA{R: 0xf7, G: 0xf7, B: 0xf7, A: 0xff},
		SidebarSelection: color.NRGBA{R: 0x17, G: 0x8b, B: 0xfb, A: 0xff},
		CardFill:         color.NRGBA{R: 0xf7, G: 0xf7, B: 0xf7, A: 0xff},
		PushButtonFill:   color.NRGBA{R: 0xec, G: 0xec, B: 0xec, A: 0xff},
		HoverOverlay:     color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x0d},
		PressOverlay:     color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x19},
		FloatingShadow:   color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x13},
		FieldEdge:        color.NRGBA{R: 0xf3, G: 0xf3, B: 0xf3, A: 0xff},
		ScrollbarThumb:   color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x92},

		AlternatingContentBackground: color.NRGBA{R: 0xf4, G: 0xf5, B: 0xf5, A: 0xff},
		Scrim:                        color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x33},
		SidebarSearchFill:            color.NRGBA{R: 0xe8, G: 0xe8, B: 0xe8, A: 0xff},
		ToolbarControlFill:           color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		ToolbarSearchFill:            color.NRGBA{R: 0xe8, G: 0xe8, B: 0xe8, A: 0xff},
		ToolbarControlShadow:         color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x09},
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

		Separator: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x19},
		Grid:      color.NRGBA{R: 0x1a, G: 0x1a, B: 0x1a, A: 0xff},

		Label:           color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xd8},
		SecondaryLabel:  color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x8c},
		TertiaryLabel:   color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x3f},
		QuaternaryLabel: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x19},

		Text:            color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		PlaceholderText: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x8c},
		SelectedText:    color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		Link:            color.NRGBA{R: 0x41, G: 0x9c, B: 0xff, A: 0xff},
		HeaderText:      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},

		Control:                      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x3f},
		ControlText:                  color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xd8},
		DisabledControlText:          color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x3f},
		SelectedControl:              color.NRGBA{R: 0x3f, G: 0x63, B: 0x8b, A: 0xff},
		SelectedControlText:          color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xd8},
		AlternateSelectedControlText: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},

		ControlAccent:          color.NRGBA{R: 0x00, G: 0x7a, B: 0xff, A: 0xff},
		KeyboardFocusIndicator: color.NRGBA{R: 0x1a, G: 0xa9, B: 0xff, A: 0x7f},

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

		SidebarMaterial:  color.NRGBA{R: 0x1c, G: 0x1c, B: 0x1c, A: 0xff},
		SidebarSelection: color.NRGBA{R: 0x19, G: 0x94, B: 0xfc, A: 0xff},
		CardFill:         color.NRGBA{R: 0x2a, G: 0x30, B: 0x34, A: 0xff},
		PushButtonFill:   color.NRGBA{R: 0x33, G: 0x3a, B: 0x3f, A: 0xff},
		HoverOverlay:     color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x18},
		PressOverlay:     color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x19},
		FloatingShadow:   color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x13},
		FieldEdge:        color.NRGBA{R: 0x2c, G: 0x33, B: 0x38, A: 0xff},
		ScrollbarThumb:   color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x92},

		AlternatingContentBackground: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x0d},
		Scrim:                        color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x42},
		SidebarSearchFill:            color.NRGBA{R: 0x2f, G: 0x32, B: 0x34, A: 0xff},
		ToolbarControlFill:           color.NRGBA{R: 0x26, G: 0x26, B: 0x26, A: 0xff},
		ToolbarSearchFill:            color.NRGBA{R: 0x36, G: 0x36, B: 0x36, A: 0xff},
		ToolbarControlShadow:         color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x06},
	}
)

// WithAccent returns the set with the six rows the platform derives from
// the accent colour rebuilt for accent. Every other field is untouched: the
// unemphasized selection pair is the inactive window's grey, Link is a fixed
// blue, and the system colours are a fixed catalogue — none of them moves
// when the user changes the accent.
//
// The six rows and their rule:
//
//   - ControlAccent is the accent itself, at the recorded row's alpha.
//   - SelectedContentBackground, SelectedTextBackground, SelectedControl,
//     KeyboardFocusIndicator and SidebarSelection keep the accent's hue and
//     saturation and
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
	p.SidebarSelection = atAccentHue(hue, sat, p.SidebarSelection)
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
