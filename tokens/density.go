package tokens

// Desktop density is the platform's control scale: Comfortable is the
// platform's regular control and Compact its small one. Every height below
// is read off a stored capture of the platform.
//
// Provenance, per number. The captures and their readings are in
// `organisation/reference/macos/controls.md`, under "What the density scale
// takes" and "What the small control captures measure":
//
//	number                      value   provenance
//	------                      -----   ----------
//	Comfortable control height  24 dp   MEASURED: the regular push button and pop-up button in save-dialog-{light,dark}.png, both appearances agreeing to the pixel
//	Compact control height      20 dp   MEASURED: the small pop-ups in finder-view-options-{light,dark}.png (Group By y 89–108, Sort By y 117–136 down x 150) and finder-search-criteria-{light,dark}.png (Kind y 7–26 down x 40); the small push buttons there agree (the − button y 7–26 down x 525, Use as Defaults y 592–611 down x 70)
//	Comfortable text field      27 dp   MEASURED: the text field in save-dialog-{light,dark}.png, y 243–269
//	Compact text field          21 dp   MEASURED: the small text field in finder-search-criteria-{light,dark}.png, its edge at y 38 and y 58 down x 340
//	Comfortable search field    26 dp   MEASURED: the search field on content in systemsettings-add-input-source-search-{light,dark}.png, y 347–372
//	Compact search field        22 dp   MEASURED: the small search field in automator-search-field-{light,dark}.png, its edge at y 69 and y 90 down x 250
//	Comfortable sidebar search  28 dp   MEASURED: the recess at the top of the sidebar in system-settings-grouped-box-{light,dark}.png, y 61–88
//	Compact sidebar search      28 dp   MEASURED: the recess in systemsettings-sidebar-search-small-{light,dark}.png, y 6–33 down x 120; the Medium crops give the same rows
//	Comfortable row height      20 dp   MEASURED: Finder's list view in finder-window-light.png, stripes on a 20 px pitch with no row seam (x=965, rows at y 85, 105, 125 … 284)
//	Compact row height          20 dp   MEASURED: Finder's list view at the Small sidebar icon size in finder-list-small-{light,dark}.png, stripes y 105–124, 125–144 … down x 580
//	Comfortable checkbox row    22 dp   MEASURED: the pitch between the two "Options:" checkboxes in save-dialog-{light,dark}.png, squares at y 372–387 and y 394–409
//	Compact checkbox row        20 dp   MEASURED: the pitch between the Show item info and Show icon preview checkboxes in finder-view-options-{light,dark}.png, boxes at y 338–351 and y 358–371 down x 27
//	Comfortable toolbar control 36 dp   MEASURED: every bordered control in the Finder toolbar captures — finder-window-untinted-dark.png (five controls, rim row y=46, fill y 47-80, rim row y=81), finder-window-untinted-light.png (fill y 34-69), finder-window-light.png (the view pop-up, y 8-43) and finder-window.png (the same pop-up, y 8-43); mail-window.png's search field and notes-toolbar.png agree at y 8-43
//	Compact toolbar control     36 dp   CARRIED: no stored capture holds a toolbar drawn at the platform's small size, and all five stored windows draw one toolbar control size
//	Comfortable PaddingX         8 dp   UNMEASURED: the platform's documented inset beside a regular push button's title. The captured dialog buttons sit at the platform's 74 px minimum width with their titles centred, so they measure centring, not an inset
//	Compact PaddingX            11 dp   MEASURED: the small pop-up's title inset, fill start to the first stem — Sort By x 81 to x 92 across y 126 in finder-view-options-{light,dark}.png, Kind x 7 to x 18 across y 17 in finder-search-criteria-{light,dark}.png
//	Comfortable SymbolSize      20 dp   CARRIED: the value a former padding derivation gave; no capture measures a symbol in a control yet
//	Compact SymbolSize          19 dp   CARRIED: the value a former padding derivation gave; no capture measures a symbol in a control yet
//	checkbox                    16 dp   MEASURED: 16 px square in save-dialog-{light,dark}.png. Not a token here: the checkbox keeps its own side length in components/input.
//
// The Sidebar icon size setting does not change a list's rows or the
// sidebar search field, so Compact's row and sidebar search heights equal
// Comfortable's. Density never changes a sidebar's rows: patterns/sidebar
// keeps them. The small push button centres its title in a width it sets
// itself, so the small pop-up gives the Compact inset. The stored windows
// measure Finder's info-pane text field at 33 px. That is a different
// control in a different place, and it corrects nothing here. Their unified
// toolbars measure their controls at 36 px, and that one IS a number here: a
// control standing in a toolbar is its own control, and
// [Density.ToolbarControlHeight] holds it.
//
// A control height is a floor, not a height. A control is as tall as its
// text's line needs, and never shorter than the density says:
//
//	height = max(floor, line)
//
// where the floor is [Density.ControlHeight] for a control and a field
// height for a field ([Density.TextFieldHeight], [Density.SearchFieldHeight],
// [Density.SidebarSearchFieldHeight]), and line is the text style's line
// height (see [TextStyle.LineHeight] and theme/typeset). The control centres
// the line in that height. The padding inside a control is what the height
// leaves around the line. A control never adds padding to the line to make
// its height. Typography does not move with density.
//
// A field's height is measured per control and per variant: a text field in
// a form (27), a search field on the content (26) and a search field in a
// sidebar (28) are three readings, each taller than a Comfortable button (24).
// A field in a toolbar is [Density.ToolbarControlHeight].
//
// A symbol in a control has its own size, [Density.SymbolSize]. No padding
// derives it.
//
// A stacked row is not a control, and takes a row height rather than the
// control height. Table rows and header cells are [Density.RowHeight] tall
// exactly (see the row table below), so a change there re-pitches every
// dense list and table in the system. A chrome rail's row is not one of
// them — the platform draws a sidebar row 12 px taller than a list row, and
// patterns/sidebar keeps that number. The platform draws a
// list row shorter than it draws a regular button — 20 px against 24 — which
// is why the row height is a number of its own and not the control's. A
// content list of one-line rows takes [Density.ContentListRowHeight] instead:
// the platform draws that kind of row 32 px tall.
//
// # Pointer targets
//
// A control's pointer target is the control: the box it draws at the heights
// above, with nothing added around it. The platform draws a 24 dp push
// button and 24 dp is what a pointer has to land on; a Compact button draws
// its 20 dp floor and offers 20. Density therefore moves
// the target with the pixels, and two controls standing side by side never
// claim the same pixels.
//
// A stacked row is its own target at its own height, and an option row at
// max(ControlHeight, its line):
//
//	row                              Comfortable   Compact   sizing
//	---                              -----------   -------   ------
//	content list row                 32            32        pinned to ContentListRowHeight
//	table body row and header cell   20            20        pinned to RowHeight
//	picker option row                24            20        max(ControlHeight, line)
//	toolbar control                  36            36        pinned to ToolbarControlHeight
//
// A checkbox is the one control drawn smaller than the box it stands in: the
// symbol keeps the platform's measured 16 dp square at every density, centred
// in a footprint of [Density.CheckboxRowHeight], and the footprint is the
// target. The platform's checkbox row is what a pointer lands on, never the
// symbol. The radio takes the same row: it is drawn at the checkbox's side
// length, and the two stand in one form.

const (
	// ComfortableControlHeight is the default desktop control-height floor in
	// dp: a Comfortable control is at least this tall, and taller when its
	// text's line needs it. It is the platform's regular control, measured
	// off the Save panel in organisation/reference/macos/save-dialog-light.png and
	// save-dialog-dark.png — a 24 px push button and pop-up button, both
	// appearances agreeing to the pixel.
	ComfortableControlHeight float32 = 24
	// CompactControlHeight is the dense-mode control-height floor in dp: the
	// platform's small control. MEASURED: the small pop-ups in
	// finder-view-options-{light,dark}.png run y 89–108 and y 117–136 down
	// x 150, 20 px each. The Kind pop-up and the − button in
	// finder-search-criteria-{light,dark}.png run y 7–26, also 20 px.
	CompactControlHeight float32 = 20
	// ComfortableTextFieldHeight is a text field's height floor in a form, in
	// dp. MEASURED: the Save panel's "Tags:" field, 27 px, y 243–269, edge
	// rows included, in save-dialog-{light,dark}.png. The same capture draws
	// the push button at 24.
	ComfortableTextFieldHeight float32 = 27
	// ComfortableSearchFieldHeight is a search field's height floor on the
	// content, in dp. MEASURED: the search field on System Settings' "add an
	// input source" sheet, 26 px, y 347–372, edge rows included, in
	// systemsettings-add-input-source-search-{light,dark}.png.
	ComfortableSearchFieldHeight float32 = 26
	// ComfortableSidebarSearchFieldHeight is a search field's height floor
	// in a sidebar, in dp. MEASURED: the recess at the top of System
	// Settings' sidebar, 28 px, y 61–88, which draws no edge, in
	// system-settings-grouped-box-{light,dark}.png.
	ComfortableSidebarSearchFieldHeight float32 = 28
	// CompactTextFieldHeight is the dense-mode text field's height floor in
	// dp. MEASURED: the small text field in
	// finder-search-criteria-{light,dark}.png has its edge at y 38 and y 58
	// down x 340, 21 px with the edge rows included.
	CompactTextFieldHeight float32 = 21
	// CompactSearchFieldHeight is the dense-mode search field's height floor
	// on the content in dp. MEASURED: the small search field in
	// automator-search-field-{light,dark}.png has its edge at y 69 and y 90
	// down x 250, 22 px with the edge rows included.
	CompactSearchFieldHeight float32 = 22
	// CompactSidebarSearchFieldHeight is the dense-mode search field's
	// height floor in a sidebar in dp. MEASURED: the recess in
	// systemsettings-sidebar-search-small-{light,dark}.png runs y 6–33 down
	// x 120, 28 px. The Medium crops give the same rows, so it equals
	// [ComfortableSidebarSearchFieldHeight].
	CompactSidebarSearchFieldHeight float32 = 28
	// ComfortableRowHeight is the height of a stacked row in dp — a list
	// row, a table row, a header cell — and a pin rather
	// than a floor: the row is drawn exactly this tall. It is the
	// platform's own list row, measured off Finder's list view in
	// organisation/reference/macos/finder-window-light.png, whose stripes
	// alternate on a 20 px pitch with no seam between them. A row is
	// shorter than a regular button on this platform, which is why it is a
	// number of its own.
	ComfortableRowHeight float32 = 20
	// CompactRowHeight is the dense-mode stacked row in dp. MEASURED:
	// Finder's list view at the Small sidebar icon size in
	// finder-list-small-{light,dark}.png, whose stripes run y 105–124,
	// 125–144 … 225–244 down x 580, a 20 px pitch. The Medium captures give
	// the same rows, so it equals [ComfortableRowHeight].
	CompactRowHeight float32 = 20
	// ContentListRowHeight is the height of a content list row in dp: one
	// entry of a content list of one-line rows, at both densities. It is a
	// pin, not a floor. MEASURED off Reminders' row, 32 px, y 169–200, in
	// reminders-row-selected-{light,dark}.png, read in "What a control
	// inside a row measures" in organisation/reference/macos/controls.md.
	// No capture holds this kind of row at the platform's small setting, so
	// the one measured row stands at Compact too.
	ContentListRowHeight float32 = 32
	// ComfortableCheckboxRowHeight is the height of a checkbox's row in dp:
	// the footprint the 16 dp symbol is centred in, and the pointer target.
	// It is a pin, not a floor — a checkbox draws no content box of its own.
	// MEASURED off the Save panel in organisation/reference/macos/save-dialog-light.png
	// and save-dialog-dark.png as the pitch between the two "Options:"
	// checkboxes, whose squares run y 372–387 and y 394–409: 22 px, both
	// appearances agreeing to the pixel. It is neither the control height nor
	// the stacked row's, which is why it is a number of its own.
	ComfortableCheckboxRowHeight float32 = 22
	// CompactCheckboxRowHeight is the dense-mode checkbox row in dp.
	// MEASURED: the pitch between the Show item info and Show icon preview
	// checkboxes in finder-view-options-{light,dark}.png, whose boxes run
	// y 338–351 and y 358–371 down x 27: 20 px. The small radio's pitch in
	// the same capture is 18 px.
	CompactCheckboxRowHeight float32 = 20
	// ComfortableToolbarControlHeight is the height of a bordered control
	// standing in a toolbar band, in dp: a search field, a capsule button, a
	// segmented control and a pop-up button alike. It is a floor in the same
	// sense as [ComfortableControlHeight].
	//
	// MEASURED off the Finder toolbar captures in
	// organisation/reference/macos, each read as a luminance run down a column
	// through the control's own middle:
	//
	//   - finder-window-untinted-dark.png, a frontmost window: all five
	//     bordered controls in its toolbar carry a rim row at y=46, the
	//     #262626 fill over y 47-80 and a rim row at y=81 — 36 px outer.
	//   - finder-window-untinted-light.png: the #f7f7f7 fill over y 34-69,
	//     36 px, no rim. That window is NOT frontmost, so its fills are the
	//     platform's inactive drawing; its extent is not.
	//   - finder-window-light.png, frontmost: the view pop-up's #ffffff over
	//     y 8-43 at x=715, told from the band it stands on by the drop shadow
	//     alone — 250 above it, 244 below.
	//   - finder-window.png, the same pop-up in the dark appearance: rim rows
	//     at y=8 and y=43 with the fill between them, 36 px.
	//
	// mail-window.png's toolbar search field (y 8-43) and notes-toolbar.png
	// (rim rows at y=8 and y=43) agree, so the number is the platform's and
	// not one application's.
	//
	// It is a number of its own because the platform draws it as one: 36
	// against the dialog pop-up's measured 24, in a 52 px unified toolbar
	// band that leaves 8 above and 8 below.
	ComfortableToolbarControlHeight float32 = 36
	// CompactToolbarControlHeight is the dense-mode toolbar control in dp. No
	// stored capture holds a toolbar drawn at the platform's small size, and
	// all five stored windows draw their toolbar controls at one height, so
	// this carries the measured 36 until a capture holds otherwise. A capture
	// of a small toolbar is on the reference's list.
	CompactToolbarControlHeight float32 = 36
	// ChipDrop is how far under the control height the system's smallest
	// control is drawn, in dp. See [Density.ChipHeight]: it is the whole of
	// that relation, exported so a reader can see the two heights are one
	// number apart rather than two picks.
	ChipDrop float32 = 4
)

// Density is one density setting: the drawn control heights, the horizontal
// padding inside a control and the size of a symbol in a control, all in dp. It is a comparable value struct like the other token
// types.
type Density struct {
	// ControlHeight is the minimum visual control height in dp
	// ([ComfortableControlHeight] or [CompactControlHeight]). It is a floor:
	// a control draws max(ControlHeight, line), so a line taller than this
	// makes the control taller.
	//
	// A field takes its own field height instead: on this platform a field
	// is taller than a button, and one number cannot be both.
	ControlHeight float32
	// TextFieldHeight is the minimum visual height of a text field in dp
	// ([ComfortableTextFieldHeight] or [CompactTextFieldHeight]), a floor in
	// the same sense as ControlHeight. A text field takes it wherever it
	// stands outside a toolbar.
	TextFieldHeight float32
	// SearchFieldHeight is the minimum visual height of a search field on
	// the content in dp ([ComfortableSearchFieldHeight] or
	// [CompactSearchFieldHeight]): the form variant.
	SearchFieldHeight float32
	// SidebarSearchFieldHeight is the minimum visual height of a search
	// field in a sidebar in dp ([ComfortableSidebarSearchFieldHeight] or
	// [CompactSidebarSearchFieldHeight]): the chrome variant's recess.
	SidebarSearchFieldHeight float32
	// RowHeight is the height of a stacked row in dp
	// ([ComfortableRowHeight] or [CompactRowHeight]) — a list row, a table
	// row, a header cell. A chrome rail's row is patterns/sidebar's own
	// number, not this one. Unlike ControlHeight it is a pin,
	// not a floor: rows tile, and a row that grew with its content would
	// cost the virtualised list the constant-time look-ahead that lets it
	// lay out only what is on screen.
	RowHeight float32
	// ContentListRowHeight is the height of a content list row in dp
	// ([ContentListRowHeight] at both settings): one entry of a content list
	// of one-line rows. It is a pin, as RowHeight is. It is separate because
	// the platform draws the content list row 32 px tall against the table
	// row's 20.
	ContentListRowHeight float32
	// CheckboxRowHeight is the height of a checkbox's row in dp
	// ([ComfortableCheckboxRowHeight] or [CompactCheckboxRowHeight]) — the
	// square footprint the checkbox's and the radio's 16 dp symbol is centred
	// in, and the pointer target both of them offer. Like RowHeight it is a
	// pin rather than a floor: the symbol does not grow, so there is no
	// content box to clear. It is separate because the platform draws it
	// separately — 22 px against the push button's 24 and the list row's 20
	// in the regular captures.
	CheckboxRowHeight float32
	// ToolbarControlHeight is the height floor of a bordered control standing
	// in a toolbar band in dp ([ComfortableToolbarControlHeight] or
	// [CompactToolbarControlHeight]) — the picker's chrome trigger, and
	// anything else this library draws bordered in a chrome region. It is
	// separate because the platform draws it separately: 36 px against the
	// dialog pop-up's 24 in the stored captures, which is the whole of why a
	// control that took the control height there would read as a dialog's
	// control standing in a toolbar.
	//
	// A field in a toolbar takes this height too: the search field's chrome
	// variant in a toolbar band measures 36 px, as every bordered control
	// there does. In a sidebar the same variant takes SidebarSearchFieldHeight.
	ToolbarControlHeight float32
	// PaddingX is the horizontal inner padding of a control in dp: the inset
	// between a control's fill and its title. The platform's small pop-up
	// insets its title 11 px; the provenance table at the top of this file
	// gives both values.
	PaddingX float32
	// SymbolSize is the side of a symbol drawn in a control, in dp: 20 at
	// Comfortable and 19 at Compact. It is a value of its own, and no padding
	// derives it. A former padding derivation gave these values. A symbol's
	// size in a control waits for a measurement, which replaces them.
	SymbolSize float32
}

// ChipHeight returns the chip height in dp: [Density.ControlHeight] less
// [ChipDrop]. A chip is smaller than a button, and the relation says that once
// instead of pinning a second scale that can drift off the first —
// Comfortable lands on 20 and Compact on 16.
//
// It is a method rather than a struct field so no Density value can carry a
// chip height that has come loose from its control height. The chip's pointer
// target is the chip: what it draws is what a pointer lands on.
func (d Density) ChipHeight() float32 { return d.ControlHeight - ChipDrop }

// PaddingX does not follow the control height. Comfortable's is the
// platform's documented inset beside a regular push button's title, and
// Compact's is the small pop-up's measured title inset. See the provenance
// table at the top of this file.
var (
	// Comfortable is the default desktop density.
	Comfortable = Density{ControlHeight: ComfortableControlHeight, TextFieldHeight: ComfortableTextFieldHeight, SearchFieldHeight: ComfortableSearchFieldHeight, SidebarSearchFieldHeight: ComfortableSidebarSearchFieldHeight, RowHeight: ComfortableRowHeight, ContentListRowHeight: ContentListRowHeight, CheckboxRowHeight: ComfortableCheckboxRowHeight, ToolbarControlHeight: ComfortableToolbarControlHeight, PaddingX: 8, SymbolSize: 20}
	// Compact is the dense mode: the platform's small controls.
	Compact = Density{ControlHeight: CompactControlHeight, TextFieldHeight: CompactTextFieldHeight, SearchFieldHeight: CompactSearchFieldHeight, SidebarSearchFieldHeight: CompactSidebarSearchFieldHeight, RowHeight: CompactRowHeight, ContentListRowHeight: ContentListRowHeight, CheckboxRowHeight: CompactCheckboxRowHeight, ToolbarControlHeight: CompactToolbarControlHeight, PaddingX: 11, SymbolSize: 19}
)
