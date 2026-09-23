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
// Twenty-two fields are not AppKit's. The platform paints a sidebar, a
// grouped box, a selected sidebar row, the count at a sidebar row's trailing
// end, a search field's recess on a sidebar, a push
// button, a hovered and a pressed
// control, the shadow under a
// floating pane, a text field's hairline, an overlay scrollbar's knob, a
// list's alternating row, the rim and the shadow of an inset sidebar panel,
// the several fills a toolbar band's controls take
// and the dim a sheet lays over the window it
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
//	SidebarSelection #178bfb        #1994fc        measured: the selected sidebar row's pill, its list holding the keyboard
//	SidebarSelectionUnemphasized       #000000 a0.043  #ffffff a0.063  measured: that pill while its list does not hold the keyboard, a coverage over the rail
//	SidebarSelectionUnemphasizedLabel  #0072f7  #148fff  measured: the label, symbol and count on that grey pill
//	SidebarCount     #6d6d6d        #a4a4a4        measured: the count at the trailing end of a sidebar row
//	SidebarSymbol    #000000        #ffffff        measured: the symbol at the leading end of a sidebar row
//	CardFill         #f7f7f7        #2a3034        measured: the System Settings grouped box
//	PushButtonFill   #ececec        #333a3f        measured: the Save dialog's push button at rest
//	DefaultButtonFill #157efb       #157efb        measured: the Save dialog's default push button at rest
//	HoverOverlay     #000000 a0.051 #ffffff a0.094 measured: a toolbar button under the pointer
//	PressOverlay     #000000 a0.098 #ffffff a0.098 measured: a push button held down
//	FloatingShadow   #000000 a0.075 #000000 a0.075 measured: the sidebar shadow, 24 px of reach centred on the surface
//	FieldEdge        #f3f3f3        #2c3338        measured: the unfocused text field's hairline in the Save dialog
//	ScrollbarThumb   #000000 a0.572 #ffffff a0.572 measured: the overlay scrollbar's knob over its track
//	AlternatingContentBackground  #f4f5f5  #ffffff a0.05  measured: the second of alternatingContentBackgroundColors, and Finder's list stripes
//	Scrim            #000000 a0.20  #000000 a0.26  measured: the dim under a Save sheet in save-dialog-{light,dark}.png
//	SidebarSearchFill  #e8e8e8      #2f3234        measured: the recess a search field is on a sidebar
//	ToolbarControlFill #ffffff      #262626        measured: a bordered control in a frontmost window's toolbar
//	ToolbarControlRim  #000000 a0.00 #404040       measured: the rim that control wears in a dark toolbar; none in light
//	ToolbarSearchFill  #e8e8e8      #363636        measured: the recess a search field is in a toolbar
//	ToolbarSearchRim   #000000 a0.00 #4d4d4d       measured: the rim that recess wears in a dark toolbar; none in light
//	ToolbarLabel     #4d4d4d        #e9e9e9        measured: what a toolbar's own title and glyphs are drawn in
//	ToolbarControlSeam #f2f2f2      #3a3a3a        measured: the line dividing one bordered toolbar control into segments
//	ToolbarControlShadow  #000000 a0.035  #000000 a0.024  measured: the drop shadow a bordered toolbar control casts, 23 px of reach sunk 9 light and 2 sunk 6 dark
//	ToolbarCheckedOverlay #000000 a0.102  #ffffff a0.161  measured: the chosen segment of a segmented toolbar control, over the control's own fill
//	PaneRim          #ffffff        #3a3a3a        measured: the 1 px rim the inset sidebar panel wears on every side
//	PaneShadow       #000000 a0.051 #000000 a0.051 measured: the shadow that panel casts, 24 px of reach with its rectangle sunk 9
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
	// the accent; the unemphasized rows are the inactive window's grey and
	// do not. UnemphasizedSelectedText is the colour a run reads in over
	// that grey: AppKit answers the opaque black and white
	// SelectedText does, and not the label's 216 of 255.
	SelectedContentBackground             color.NRGBA
	UnemphasizedSelectedContentBackground color.NRGBA
	SelectedTextBackground                color.NRGBA
	UnemphasizedSelectedTextBackground    color.NRGBA
	UnemphasizedSelectedText              color.NRGBA

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

	// SidebarSelectionUnemphasized is that same pill while its list does
	// not hold the keyboard: black at 11/255 light and white at 16/255
	// dark, laid over whatever the rail is painting. It is a COVERAGE and
	// not a value, which the capture says outright: in
	// finder-sidebar-unfocused-light.png the rail alternates #fafafa and
	// #fafaf9 down its own dither, and the pill over it alternates #efefef
	// and #efefee in the same columns — 123 columns of one pair and 108 of
	// the other, read across the pill's middle row. The dark capture does
	// the same, #1c1c1c to #2a2a2a in 164 columns and #1b1b1b to #292929 in
	// 59. A fill that carries the rail's own dither through it is
	// translucent; an opaque paint would flatten it.
	//
	// The pill's box is x 18-307, y 212-243 in both captures — 290 by 32,
	// inset 10 from the panel's rim and cornered at 8, the emphasized
	// pill's own geometry unchanged between the two states.
	//
	// Flattened onto SidebarMaterial with [github.com/vibrantgio/theme/color.Flatten]
	// the coverage lands #ececec light and #2a2a2a dark, an 11 of 255 step
	// under the rail in light and 14 in dark — the step the captures hold,
	// which recording the pixel instead would have shrunk to 8 in light,
	// because the capture's rail carries wallpaper tinting and this set's
	// material is the untinted reading.
	//
	// It is not UnemphasizedSelectedContentBackground, which reports
	// #dcdcdc light and #464646 dark. It is black and white rather than a
	// colour, so WithAccent leaves it where it stands. See
	// reference/macos/controls.md, "What the unfocused sidebar pill
	// measures".
	SidebarSelectionUnemphasized color.NRGBA `appkit:"-"`

	// SidebarSelectionUnemphasizedLabel is what the label, the symbol and
	// the count on that grey pill are drawn in: #0072f7 light and #148fff
	// dark, the plateau the row's name holds over 47 and 43 pixels in
	// finder-sidebar-unfocused-{light,dark}.png, and the plateau the
	// document mark beside it holds over 23 pixels in each. A vector mark
	// takes no stem darkening, so the two rasterizations agreeing on one
	// value is the drawn colour and not a stem's shortfall.
	//
	// It is the accent as the platform's vibrancy lands it on the pill, and
	// no name in the set is it in both appearances: ControlAccent's #007aff
	// is 8 of 255 off in light and 21 in dark, and SidebarSelection's own
	// lift is 5 off in dark and 25 in light. The label follows the user's
	// accent, so [PlatformColors.WithAccent] moves it, as it moves the
	// emphasized pill.
	SidebarSelectionUnemphasizedLabel color.NRGBA `appkit:"-"`

	// SidebarCount is what the count at the trailing end of a sidebar row
	// is drawn in: #6d6d6d light and #a4a4a4 dark, the plateau every count
	// holds in voicememos-multi-folder-2026-09-18.png and
	// voicememos-sidebar-dark.png over the panel's own #f9f9f9 and #1c1c1c.
	//
	// It is not SecondaryLabel, which is what the section label above those
	// rows IS: that label plateaus at #7d7d7d and #999999 on the same fills,
	// black at 0.5 and white at 0.55 to the byte, and the count plateaus 16
	// of 255 darker in light and 11 lighter in dark. Both are plateaux over
	// whole pixels, so the gap is the drawn colour and not a stem's
	// shortfall. No field of this set flattens within a level of either
	// reading — the nearest is ScrollbarThumb at 3 and 6, which is a
	// scrollbar's knob — so the count is recorded as a measured value of the
	// sidebar, as the pill is. See reference/macos/controls.md, "What a
	// sidebar row measures".
	SidebarCount color.NRGBA `appkit:"-"`

	// SidebarSymbol is what the symbol at the leading end of a sidebar row
	// is drawn in: #000000 light and #ffffff dark, the plateau the folder
	// mark holds over 45 pixels in voicememos-multi-folder-2026-09-18.png
	// and over 30 in voicememos-sidebar-light.png, and #ffffff in
	// voicememos-sidebar-dark.png.
	//
	// It is not Label, which is what the name beside it is: that label
	// plateaus at #262626 and #dcdcdc on the same fills, black and white at
	// 216 of 255, so the symbol stands 38 of 255 stronger in light and 35 in
	// dark. A vector mark takes no stem darkening, so both readings are the
	// drawn colour. No field of this set answers for black or white outright,
	// so the symbol is recorded as a measured value of the sidebar, as the
	// pill and the count are. On the pill it wears the pill's own foreground
	// instead. See reference/macos/controls.md, "What a sidebar row
	// measures".
	SidebarSymbol color.NRGBA `appkit:"-"`

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

	// DefaultButtonFill is what the DEFAULT push button — the one a sheet
	// binds Return to — is actually filled with: #157efb in BOTH
	// appearances, flat-region samples of the "Save" button in
	// save-dialog-light.png and save-dialog-dark.png (x 445-460 and
	// x 496-512, y 505-520, #157efb over every one of those 528 pixels in
	// each capture). Its label is #ffffff on both sheets, which is
	// AlternateSelectedControlText to the byte.
	//
	// It is not ControlAccent. AppKit's controlAccentColor reports #007aff
	// in both appearances; the platform draws the default button's bezel
	// over it and what reaches the screen is eleven of 255 off that on the
	// red channel and four on the green. The lift is the same kind
	// SidebarSelection carries and is recorded the same way — as the pixel,
	// in both appearances — so the library paints what the platform draws
	// rather than the published name. The bezel's own gradient is not
	// recorded; a consumer paints this flat.
	//
	// It follows the user's accent, as SidebarSelection does:
	// [PlatformColors.WithAccent] moves it.
	DefaultButtonFill color.NRGBA `appkit:"-"`

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
	// falling to nothing 24 px out, centred on the surface: the captures
	// are of a vertical edge and nothing measured lights this one from
	// above, so its rectangle is the surface's own.
	//
	// It does not vary with what is floating — the platform draws one
	// shadow — which is why one reading serves a dialog, a menu, a popover
	// and a toast alike.
	FloatingShadow DropShadow `appkit:"-"`

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

	// ToolbarControlRim is the hairline a BORDERED TOOLBAR CONTROL wears
	// round its own edge: #404040 in the dark appearance and no colour at
	// all in the light one, where the platform draws none. It is a value of
	// its own beside [PlatformColors.ToolbarControlFill] because no alpha
	// name lands it — the platform's separator over that fill falls five of
	// 255 short — and apart from [PlatformColors.ToolbarSearchRim], which is
	// what the toolbar's search recess wears: a different pixel over a
	// different fill.
	//
	// MEASURED, finder-window-untinted-dark.png, a frontmost window: the
	// rows immediately above and below a bordered control's #262626 fill
	// read 64 at every column of its flat middle, and the columns at either
	// end read 61-62 through the corner's antialiasing. It is lighter than
	// both the fill and the #1e1e1e band — the highlight every bordered
	// control in a dark toolbar band wears. separatorColor over that fill
	// gives #3b3b3b and over the band #323232, so neither name lands it.
	//
	// MEASURED, finder-window-light.png and finder-window-untinted-light.png:
	// a light band steps straight to the fill in one row with no stroke row
	// on any side, so the light value is black at zero coverage — no colour,
	// and a caller draws nothing where it answers one. What falls outside the
	// control there is its drop shadow, which darkens away from the control
	// and never sits against it.
	ToolbarControlRim color.NRGBA `appkit:"-"`

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

	// ToolbarSearchRim is the hairline the toolbar search recess wears round
	// its own edge: #4d4d4d in the dark appearance and no colour at all in
	// the light one, where the platform draws none. It is a value of its own
	// beside [PlatformColors.ToolbarSearchFill] because no alpha name lands
	// it — the platform's separator over that fill falls three of 255 short
	// — and apart from the rim a BORDERED toolbar control wears, which is a
	// different pixel over a different fill.
	//
	// MEASURED, voicememos-window.png, the search field at the trailing end
	// of a frontmost dark Voice Memos toolbar, x 643-967, y 8-43: the rows at
	// y=8 and y=43 read #4d4d4d flat over x 669-941 and fall away through the
	// corners' antialiasing, and the columns at x=643 and x=967 read #4b4b4b
	// and #4a4a4a at the control's own middle row. It runs the whole way
	// round, so it is the control's own edge and not the band's seam, and it
	// is lighter than both its #363636 fill and the #1e1e1e band — the
	// highlight every bordered control in a dark toolbar band wears.
	// Separator over that fill gives #4a4a4a, three of 255 short of the
	// pixel, which is why the name is carried here rather than flattened.
	//
	// MEASURED, voicememos-sidebar-light.png and finder-window-light.png: a
	// light toolbar band steps straight to the fill with no stroke row on any
	// side, so the light value is black at zero coverage — no colour, and a
	// caller draws nothing where it answers one.
	ToolbarSearchRim color.NRGBA `appkit:"-"`

	// ToolbarLabel is what a TOOLBAR draws its own words and its own glyphs
	// in: #4d4d4d light and #e9e9e9 dark. It is the band's title standing
	// bare, the wording inside a bordered control, and every mark a chrome
	// control carries — one foreground for everything a toolbar says.
	//
	// It is a value of its own beside [PlatformColors.ControlText] because
	// no name in the catalogue flattens there. Light, over the band's
	// #ffffff: ControlText's 216 of 255 gives #272727, SecondaryLabel's 127
	// gives #808080 and TertiaryLabel's 66 gives #bdbdbd, against the
	// measured #4d4d4d, which is black at 178 of 255 and no name's coverage.
	// Dark, over the #1e1e1e band: ControlText gives #dcdcdc against the
	// measured #e9e9e9, white at 230 of 255.
	//
	// A FORM control is not this: the Save dialog's pop-up draws its mark and
	// its label at ControlText exactly, which is why that name stays where a
	// dialog's controls are drawn and this one answers for the toolbar.
	//
	// MEASURED at 1x, finder-window-light.png, a frontmost light Finder
	// window whose band is #ffffff: the title "Applications" standing bare in
	// the band (x 413-499) plateaus at #4d4d4d over 131 pixels, the group
	// pull-down's grid glyph (x 769-786, y 17-34) over 24 and the search
	// capsule's magnifier (x 965-980, y 18-34) over 15. A word and a glyph
	// hold the same plateau, so it is the drawn colour and not the partial
	// coverage a thin stroke reaches.
	//
	// MEASURED at 1x, finder-window-untinted-dark.png (the window's own
	// origin at x=56, y=38 in that capture): the pull-down's glyph plateaus
	// at #e9e9e9 over its own fill and the window's title at #e8e8e8 over the
	// band, one 255th below it — text against a vector mark, which takes no
	// stem darkening. The glyph's reading is the one carried.
	ToolbarLabel color.NRGBA `appkit:"-"`

	// ToolbarControlSeam is the line that divides ONE bordered toolbar
	// control into segments: #f2f2f2 light over that control's #ffffff fill
	// and #3a3a3a dark over its #262626. The segments are divisions of one
	// capsule and not controls side by side, so the seam is the control's own
	// and stops short of its top and its foot.
	//
	// It is a value of its own beside [PlatformColors.ToolbarControlRim], as
	// that rim is, because the platform's separator over the fill lands
	// neither: light it gives #e6e6e6, twelve of 255 short of the pixel, and
	// dark #3b3b3b, one over it.
	//
	// MEASURED at 1x, finder-window-light.png and
	// finder-window-untinted-dark.png, the back/forward pair each window
	// keeps at the leading end of its content column: the control spans
	// 73 × 36 px and its seam is one column wide and twenty rows tall, eight
	// rows clear of the control's top and eight of its foot.
	ToolbarControlSeam color.NRGBA `appkit:"-"`

	// ToolbarControlShadow is the peak coverage of the drop shadow a
	// BORDERED TOOLBAR CONTROL casts on the band it stands in: black at
	// 0.035 light and at 0.024 dark. It is the shadow that tells a light
	// control from its band — the platform draws no edge there and the fill
	// is the band's own white — and this field carries the peak the way
	// [PlatformColors.FloatingShadow] does, the caller spreading it over the
	// measured reach.
	//
	// The reach and the offset are this field's own, per appearance: light
	// the ramp carries 23 px from a rectangle sunk 9 px below the control,
	// dark 2 px from one sunk 6. That is what makes the shadow heavier under
	// the control than over it, and it is carried here because a coverage
	// fitted at one reach is not the same coverage at another — see
	// [DropShadow].
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
	// of those bytes exactly; the dark control is told from its band by its
	// fill and its rim rather than by this, which is why so little of it
	// shows.
	//
	// MEASURED, dark, finder-window-untinted-dark.png and notes-toolbar.png,
	// re-read 2026-09-18: fitted over 10,575 band pixels around the Finder
	// search field (x 1155-1379, y 46-81) and the Notes compose control
	// (x 8-44, y 8-43) at this coverage, the best whole-pixel pair is 2 px of
	// reach with the rectangle sunk 6 px — 312 of 31,725 channel samples off
	// by one 255th and none by more except at the control's own antialiased
	// corner. The light pair fitted to the same samples is eleven times
	// worse, which is why the geometry is per appearance rather than one
	// shape drawn twice.
	ToolbarControlShadow DropShadow `appkit:"-"`

	// ToolbarCheckedOverlay is what a bordered toolbar control lays over its
	// own fill while it records a yes: black at 0.102 light and white at
	// 0.161 dark, the coverage the platform draws the chosen segment of a
	// segmented control at. The patch it fills is not the control's whole
	// box — see components/internal/controlface, which draws it inset and
	// cornered as measured.
	//
	// MEASURED, finder-window-untinted-dark.png, the four-segment view
	// control at x 796-943, y 46-81 with its list segment chosen: the patch
	// spans x 836-867, y 51-76 and reads #494949 flat against the control's
	// own #262626, which is white at 41 of 255 over it to the byte.
	//
	// MEASURED, finder-window-untinted-light.png, the same control in that
	// window: the patch spans x 814-845, y 39-64 — the same 32 by 26 — and
	// reads #dedede against the control's #f7f7f7, which is black at 26 of
	// 255 over it to the byte.
	// That window is NOT frontmost, so both its fill and its patch are the
	// platform's faded drawing and the PIXEL is not an active control's; the
	// COVERAGE between the two is what this field carries, and over the
	// frontmost light control's own #ffffff it lands #e5e5e5. No stored
	// capture holds a frontmost light window with a chosen segment in it,
	// which is on the capture list.
	//
	// It is a coverage rather than a fill because the fill beneath it moves:
	// a checked control still tints under the pointer, and the patch is laid
	// over whatever that leaves.
	ToolbarCheckedOverlay color.NRGBA `appkit:"-"`

	// PaneRim is the 1 px rim the platform draws round an inset sidebar
	// panel, on every side: opaque #ffffff light and #3a3a3a dark. It is the
	// boundary between the panel and what stands around it — the window's
	// own plane on three sides, the content on the fourth — which is why a
	// panel set into the window needs no seam.
	//
	// MEASURED, voicememos-multi-folder-2026-09-18.png at 1x, the panel at
	// x 64-283, y 46-786 in a window at x 56-1031, y 38-794: the rim reads
	// #ffffff flat down the column at x=64, along the rows at y=46 and
	// y=786 and down x=283, with the panel's own #f9f9f9 inside it and the
	// window's plane or the content outside. It is opaque and not a
	// coverage: white at any coverage over a fill that light lands short of
	// 255, and the pixel is 255 on every channel.
	//
	// MEASURED, finder-window-untinted-dark.png, the panel at x 64-373,
	// y 46-1024 in a window at x 56-1442, y 38-1032: the rim reads #3a3a3a
	// flat down x=64 and along y=46 and y=1024, over the panel's own
	// #1c1c1c and the plane beside it. Down the trailing edge at x=373 it
	// reads #404040 through the band and #434343 below it, the panel's own
	// sidebar material lifting toward that edge; the three sides over the
	// plane are what this field carries. Separator over the dark fill gives #3b3b3b,
	// one of 255 off the three flat sides and nine off the trailing one, and
	// over the light fill #e6e6e6, nowhere near white — so the value is
	// carried rather than flattened.
	PaneRim color.NRGBA `appkit:"-"`

	// PaneShadow is the peak coverage of the shadow an inset sidebar panel
	// casts on what lies around it: black at 0.051 — 13 of 255 — in both
	// appearances, as [PlatformColors.FloatingShadow] is one value in both.
	// It is spread over 24 px of reach with its rectangle sunk 9 px below
	// the panel, which is what makes it heavier under the panel than over
	// it.
	//
	// MEASURED, voicememos-multi-folder-2026-09-18.png: the panel at
	// x 64-283, y 46-786 stands on a white plane and a white content column.
	// Beside its trailing rim the content reads 244 and recovers to #ffffff
	// 33 columns out; the 8 px of plane at its leading edge reads 239 at the
	// rim and 245 at the window's edge; the 8 px above it reads 247 at the
	// rim and 251 at the window's edge. A linear ramp of black at 13 of 255
	// over 24 px of reach, its rectangle sunk 9 px, lands those 9,096
	// sampled pixels at an rms of 1.19 of 255 and a worst miss of 3.4.
	//
	// The 8 px of plane BELOW the panel reads 227 at the rim and 234 at the
	// window's edge, 8 to 11 of 255 deeper than one sunk rectangle puts it:
	// the platform's own shadow is blurred and lit from above, and a single
	// rectangle with one peak cannot be both that deep below and that light
	// beside. The same limit is recorded for
	// [PlatformColors.ToolbarControlShadow], fitted the same way.
	//
	// MEASURED, finder-window-untinted-dark.png: the same shadow is one to
	// two of 255 deep — the plane beside the panel's leading rim reads 27
	// and recovers to 28 within 6 columns, the content beside its trailing
	// rim 29 against its own #1e1e1e — which is what this coverage lands on
	// a plane that dark, so one value serves both appearances.
	PaneShadow DropShadow `appkit:"-"`
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
		UnemphasizedSelectedText:              color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},

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

		SidebarMaterial:                   color.NRGBA{R: 0xf7, G: 0xf7, B: 0xf7, A: 0xff},
		SidebarSelection:                  color.NRGBA{R: 0x17, G: 0x8b, B: 0xfb, A: 0xff},
		SidebarSelectionUnemphasized:      color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x0b},
		SidebarSelectionUnemphasizedLabel: color.NRGBA{R: 0x00, G: 0x72, B: 0xf7, A: 0xff},
		SidebarCount:                      color.NRGBA{R: 0x6d, G: 0x6d, B: 0x6d, A: 0xff},
		SidebarSymbol:                     color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		CardFill:                          color.NRGBA{R: 0xf7, G: 0xf7, B: 0xf7, A: 0xff},
		PushButtonFill:                    color.NRGBA{R: 0xec, G: 0xec, B: 0xec, A: 0xff},
		DefaultButtonFill:                 color.NRGBA{R: 0x15, G: 0x7e, B: 0xfb, A: 0xff},
		HoverOverlay:                      color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x0d},
		PressOverlay:                      color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x19},
		FloatingShadow:                    DropShadow{Peak: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x13}, Reach: 24},
		FieldEdge:                         color.NRGBA{R: 0xf3, G: 0xf3, B: 0xf3, A: 0xff},
		ScrollbarThumb:                    color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x92},

		AlternatingContentBackground: color.NRGBA{R: 0xf4, G: 0xf5, B: 0xf5, A: 0xff},
		Scrim:                        color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x33},
		SidebarSearchFill:            color.NRGBA{R: 0xe8, G: 0xe8, B: 0xe8, A: 0xff},
		ToolbarControlFill:           color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		ToolbarControlRim:            color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00},
		ToolbarSearchFill:            color.NRGBA{R: 0xe8, G: 0xe8, B: 0xe8, A: 0xff},
		ToolbarSearchRim:             color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00},
		ToolbarLabel:                 color.NRGBA{R: 0x4d, G: 0x4d, B: 0x4d, A: 0xff},
		ToolbarControlSeam:           color.NRGBA{R: 0xf2, G: 0xf2, B: 0xf2, A: 0xff},
		ToolbarControlShadow:         DropShadow{Peak: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x09}, Reach: 23, Offset: 9},
		ToolbarCheckedOverlay:        color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x1a},
		PaneRim:                      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		PaneShadow:                   DropShadow{Peak: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x0d}, Reach: 24, Offset: 9},
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
		UnemphasizedSelectedText:              color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},

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

		SidebarMaterial:                   color.NRGBA{R: 0x1c, G: 0x1c, B: 0x1c, A: 0xff},
		SidebarSelection:                  color.NRGBA{R: 0x19, G: 0x94, B: 0xfc, A: 0xff},
		SidebarSelectionUnemphasized:      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x10},
		SidebarSelectionUnemphasizedLabel: color.NRGBA{R: 0x14, G: 0x8f, B: 0xff, A: 0xff},
		SidebarCount:                      color.NRGBA{R: 0xa4, G: 0xa4, B: 0xa4, A: 0xff},
		SidebarSymbol:                     color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		CardFill:                          color.NRGBA{R: 0x2a, G: 0x30, B: 0x34, A: 0xff},
		PushButtonFill:                    color.NRGBA{R: 0x33, G: 0x3a, B: 0x3f, A: 0xff},
		DefaultButtonFill:                 color.NRGBA{R: 0x15, G: 0x7e, B: 0xfb, A: 0xff},
		HoverOverlay:                      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x18},
		PressOverlay:                      color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x19},
		FloatingShadow:                    DropShadow{Peak: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x13}, Reach: 24},
		FieldEdge:                         color.NRGBA{R: 0x2c, G: 0x33, B: 0x38, A: 0xff},
		ScrollbarThumb:                    color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x92},

		AlternatingContentBackground: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x0d},
		Scrim:                        color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x42},
		SidebarSearchFill:            color.NRGBA{R: 0x2f, G: 0x32, B: 0x34, A: 0xff},
		ToolbarControlFill:           color.NRGBA{R: 0x26, G: 0x26, B: 0x26, A: 0xff},
		ToolbarControlRim:            color.NRGBA{R: 0x40, G: 0x40, B: 0x40, A: 0xff},
		ToolbarSearchFill:            color.NRGBA{R: 0x36, G: 0x36, B: 0x36, A: 0xff},
		ToolbarSearchRim:             color.NRGBA{R: 0x4d, G: 0x4d, B: 0x4d, A: 0xff},
		ToolbarLabel:                 color.NRGBA{R: 0xe9, G: 0xe9, B: 0xe9, A: 0xff},
		ToolbarControlSeam:           color.NRGBA{R: 0x3a, G: 0x3a, B: 0x3a, A: 0xff},
		ToolbarControlShadow:         DropShadow{Peak: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x06}, Reach: 2, Offset: 6},
		ToolbarCheckedOverlay:        color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x29},
		PaneRim:                      color.NRGBA{R: 0x3a, G: 0x3a, B: 0x3a, A: 0xff},
		PaneShadow:                   DropShadow{Peak: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x0d}, Reach: 24, Offset: 9},
	}
)

// WithAccent returns the set with the seven rows the platform derives from
// the accent colour rebuilt for accent. Every other field is untouched: the
// unemphasized selection pair is the inactive window's grey, Link is a fixed
// blue, and the system colours are a fixed catalogue — none of them moves
// when the user changes the accent.
//
// The seven rows and their rule:
//
//   - ControlAccent is the accent itself, at the recorded row's alpha.
//   - SelectedContentBackground, SelectedTextBackground, SelectedControl,
//     KeyboardFocusIndicator, SidebarSelection and DefaultButtonFill keep
//     the accent's hue and saturation and
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
	p.SidebarSelectionUnemphasizedLabel = atAccentHue(hue, sat, p.SidebarSelectionUnemphasizedLabel)
	p.DefaultButtonFill = atAccentHue(hue, sat, p.DefaultButtonFill)
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
