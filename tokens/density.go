package tokens

// Desktop density is the platform's control scale: Comfortable is the
// platform's regular control and Compact its small one. Nothing below is
// derived from a web scale.
//
// Measured beats published. Where a stored capture and the platform's
// published guidelines disagree, the capture wins (owner ruling,
// 2026-09-11): the platform is what the application is judged against, and
// the published table is a document about it.
//
// Provenance, per number — the reference and its remaining gap are
// `.github/reference/macos/controls.md`, indexed by ADR-019:
//
//	number                      value   provenance
//	------                      -----   ----------
//	Comfortable control height  24 dp   MEASURED: the regular push button and pop-up button in save-dialog-light.png and save-dialog-dark.png, both appearances agreeing to the pixel; the HIG publishes 22 pt, superseded
//	Compact control height      19 dp   PUBLISHED: the HIG's small push button — no small control is captured yet
//	Comfortable field height    27 dp   MEASURED: the text field in the same pair; the HIG publishes 22 pt, superseded
//	Compact field height        21 dp   DERIVED: 27 × 19/24 = 21.4, rounded — the platform's field-to-control ratio applied to the small control, until a small field is captured
//	Comfortable PaddingX         8 dp   PUBLISHED: the HIG's horizontal inset beside a regular push button's label — the captured buttons both sit at the platform's 74 px minimum width with their labels centred, so what they measure is centring, not an inset
//	Compact PaddingX             7 dp   DERIVED: 8 × 19/22, rounded — the published small-to-regular ratio, both operands published, since neither the inset nor the small control is captured
//	Comfortable PaddingY         2 dp   DERIVED: (24 − LabelLarge's 20 dp line box) / 2, which lands a button exactly on 24
//	Compact PaddingY             0 dp   DERIVED: LabelLarge's line box is already over 19, so there is nothing to pad with
//	checkbox                    16 dp   MEASURED: 16 px square in the same pair. Not a token here: the checkbox carries its own side length in components/input, which is where this number lands when a consumer takes it.
//	Comfortable row height      20 dp   MEASURED: Finder's list view in finder-window-light.png — the stripes alternate on a 20 px pitch with no row seam between them (x=965, rows at y 85, 105, 125 … 284)
//	Compact row height          19 dp   PUBLISHED: the platform's small control, carried until a capture holds a list drawn dense
//	Comfortable checkbox row    22 dp   MEASURED: the pitch between the two "Options:" checkboxes in the same pair — squares at y 372–387 and y 394–409, so 394 − 372
//	Compact checkbox row        17 dp   DERIVED: 22 × 19/24 = 17.4, rounded — the same ratio the compact field height takes, until a small checkbox is captured — no stored capture does
//	Comfortable toolbar control 36 dp   MEASURED: every bordered control in the Finder toolbar captures — finder-window-untinted-dark.png (five controls, rim row y=46, fill y 47-80, rim row y=81), finder-window-untinted-light.png (fill y 34-69), finder-window-light.png (the view pop-up, y 8-43) and finder-window.png (the same pop-up, y 8-43); mail-window.png's search field and notes-toolbar.png agree at y 8-43
//	Compact toolbar control     36 dp   CARRIED: no stored capture holds a toolbar drawn at the platform's small size, and all five stored windows draw one toolbar control size, so Compact carries the measured 36 until a capture holds otherwise
//
// What is measured and what is not. One Save panel, captured at 1x in both
// appearances, holds the regular push button, pop-up button, text field and
// checkbox this scale is named after, and every Comfortable number above is
// read off its pixels. No capture holds a control at the platform's small
// size, so Compact's height is still the published 19 and its field height
// is derived from the measured ratio; ADR-019 carries that gap row. The
// stored windows measure Finder's info-pane text field at 33 px — a different
// control in a different place, recorded beside these and correcting nothing
// here. Their unified toolbars measure their controls at 36 px, and that one
// IS a number here: a control standing in a toolbar is its own control, and
// [Density.ToolbarControlHeight] carries it.
//
// A control height is a floor, not a height. A control is as tall as its
// content box needs, and never shorter than the density says:
//
//	height = max(floor, lineBox + 2×PaddingY)
//
// where the floor is [Density.ControlHeight] for a control and
// [Density.FieldHeight] for a text field, and lineBox is the type role's
// line height (see [TextStyle.LineHeight] and theme/typeset). Typography
// does not move with density, so at the platform's heights the two terms are
// close and either can win:
//
//	control                  role         line box   + 2×PaddingY   floor              drawn
//	-------                  ----         --------   ------------   -----              -----
//	button, Comfortable      LabelLarge   20         24             24 ControlHeight   24
//	button, Compact          LabelLarge   20         20             19 ControlHeight   20
//	text field, Comfortable  BodyLarge    24         28             27 FieldHeight     28
//	text field, Compact      BodyLarge    24         24             21 FieldHeight     24
//
// Comfortable's padding is what makes a button exact: 20 + 2×2 = 24. Two
// overshoots are recorded rather than hidden, and both close only by moving
// a typography role, which this scale does not do:
//
//   - a Compact button draws 20 dp against a 19 dp floor, because
//     LabelLarge's line box is 1 dp over the platform's small control before
//     any padding, so PaddingY is 0 and there is nothing left to spend.
//   - a Comfortable text field draws 28 dp against the platform's measured
//     27, because BodyLarge's 24 dp line box plus the control's own 2 dp
//     padding is 28. The measured height is the floor; the content box wins
//     by a dp.
//
// The consequence worth stating: controls in different type roles come out at
// different heights, and a Comfortable text field (28) is taller than a
// Comfortable button (24) because BodyLarge is a larger role than LabelLarge.
// Both are honest readings of the tokens. A design that wants them equal
// changes the roles or the padding, not the measurement. The platform draws
// the same difference, 27 against 24.
//
// A stacked row is not a control, and takes [Density.RowHeight] rather than
// the control height: list rows, table rows and header cells are RowHeight
// tall exactly (see the row table below), so a change there re-pitches every
// dense list and table in the system. A chrome rail's row is not one of
// them — the platform draws a sidebar row 12 px taller than a list row, and
// patterns/sidebar carries that number. The platform draws a
// list row shorter than it draws a button — 20 px against 24 — which is why
// the row height is a number of its own and not the control's.
//
// # Pointer targets
//
// A control's pointer target is the control: the box it draws at the heights
// above, with nothing added around it. The platform draws a 24 dp push
// button and 24 dp is what a pointer has to land on; a Compact control draws
// 19 and offers 19. Density therefore moves the target with the pixels, and
// two controls standing side by side never claim the same pixels.
//
// A stacked row is its own target at its own height, and an option row at
// the height its content box gives it:
//
//	row                              Comfortable   Compact   sizing
//	---                              -----------   -------   ------
//	list row                         20            19        pinned to RowHeight
//	table body row and header cell   20            19        pinned to RowHeight
//	picker option row                28            24        floor formula, BodyLarge
//	toolbar control                  36            36        pinned to ToolbarControlHeight
//
// A checkbox is the one control drawn smaller than the box it stands in: the
// glyph keeps the platform's measured 16 dp square at every density, centred
// in a footprint of [Density.CheckboxRowHeight], and the footprint is the
// target. The platform's checkbox row is what a pointer lands on, never the
// glyph. The radio takes the same row: it is drawn at the checkbox's side
// length, and the two stand in one form.

const (
	// ComfortableControlHeight is the default desktop control-height floor in
	// dp: a Comfortable control is at least this tall, and taller when its
	// content box needs it. It is the platform's regular control, measured
	// off the Save panel in .github/reference/macos/save-dialog-light.png and
	// save-dialog-dark.png — a 24 px push button and pop-up button, both
	// appearances agreeing to the pixel. The HIG publishes 22 pt; the capture
	// supersedes it.
	ComfortableControlHeight float32 = 24
	// CompactControlHeight is the dense-mode control-height floor in dp: the
	// platform's small control, the HIG's published 19 pt push button. It is
	// the one height here no capture holds. A control set in LabelLarge or
	// BodyLarge draws over it — see the table above — so for those it is
	// never the answer, and for a stacked row it is a pin rather than a
	// floor.
	CompactControlHeight float32 = 19
	// ComfortableFieldHeight is the text field's own height floor in dp,
	// measured at 27 px off the same Save panel. A field is not a button's
	// height on this platform and does not take one: the HIG publishes 22 pt
	// for both, and the capture draws 27 against 24.
	ComfortableFieldHeight float32 = 27
	// CompactFieldHeight is the dense-mode text field's height floor in dp,
	// derived at 27 × 19/24 = 21.4, rounded: the platform's measured
	// field-to-control ratio applied to the small control, since no capture
	// holds a small field. It is the number a capture of one replaces.
	CompactFieldHeight float32 = 21
	// ComfortableRowHeight is the height of a stacked row in dp — a list
	// row, a table row, a header cell — and a pin rather
	// than a floor: the row is drawn exactly this tall. It is the
	// platform's own list row, measured off Finder's list view in
	// .github/reference/macos/finder-window-light.png, whose stripes
	// alternate on a 20 px pitch with no seam between them. A row is
	// shorter than a button on this platform, which is why it is a number
	// of its own.
	ComfortableRowHeight float32 = 20
	// CompactRowHeight is the dense-mode stacked row in dp. No stored
	// capture holds a list drawn dense, so it carries the platform's small
	// control height until one does; ADR-019 has the gap row.
	CompactRowHeight float32 = 19
	// ComfortableCheckboxRowHeight is the height of a checkbox's row in dp:
	// the footprint the 16 dp glyph is centred in, and the pointer target.
	// It is a pin, not a floor — a checkbox draws no content box of its own.
	// MEASURED off the Save panel in .github/reference/macos/save-dialog-light.png
	// and save-dialog-dark.png as the pitch between the two "Options:"
	// checkboxes, whose squares run y 372–387 and y 394–409: 22 px, both
	// appearances agreeing to the pixel. It is neither the control height nor
	// the stacked row's, which is why it is a number of its own.
	ComfortableCheckboxRowHeight float32 = 22
	// CompactCheckboxRowHeight is the dense-mode checkbox row in dp, derived
	// at 22 × 19/24 = 17.4, rounded: the platform's regular-to-small control
	// ratio applied to the measured row, the same derivation
	// [CompactFieldHeight] takes, since no capture holds a small checkbox.
	CompactCheckboxRowHeight float32 = 17
	// ComfortableToolbarControlHeight is the height of a bordered control
	// standing in a toolbar band, in dp: a search field, a capsule button, a
	// segmented control and a pop-up button alike. It is a floor in the same
	// sense as [ComfortableControlHeight].
	//
	// MEASURED off the Finder toolbar captures in
	// .github/reference/macos, each read as a luminance run down a column
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

// Density is one density setting: the drawn control height and its inner
// padding, all in dp. It is a comparable value struct like the other token
// types.
type Density struct {
	// ControlHeight is the minimum visual control height in dp
	// ([ComfortableControlHeight] or [CompactControlHeight]). It is a floor:
	// a control draws max(ControlHeight, contentHeight + 2×PaddingY), so a
	// content box taller than this makes the control taller. See the table
	// above the constants for which controls clear it and by how much.
	//
	// A text field takes [Density.FieldHeight] instead: on this platform a
	// field is taller than a button, and one number cannot be both.
	ControlHeight float32
	// FieldHeight is the minimum visual height of a text field in dp
	// ([ComfortableFieldHeight] or [CompactFieldHeight]), a floor in the same
	// sense as ControlHeight. It is separate because the platform draws it
	// separately — 27 px against the push button's 24 in the stored capture —
	// so a field that took the control height would be 3 dp short of the
	// platform at Comfortable.
	FieldHeight float32
	// RowHeight is the height of a stacked row in dp
	// ([ComfortableRowHeight] or [CompactRowHeight]) — a list row, a table
	// row, a header cell. A chrome rail's row is patterns/sidebar's own
	// number, not this one. Unlike ControlHeight it is a pin,
	// not a floor: rows tile, and a row that grew with its content would
	// cost the virtualised list the constant-time look-ahead that lets it
	// lay out only what is on screen.
	RowHeight float32
	// CheckboxRowHeight is the height of a checkbox's row in dp
	// ([ComfortableCheckboxRowHeight] or [CompactCheckboxRowHeight]) — the
	// square footprint the checkbox's and the radio's 16 dp glyph is centred
	// in, and the pointer target both of them offer. Like RowHeight it is a
	// pin rather than a floor: the glyph does not grow, so there is no
	// content box to clear. It is separate because the platform draws it
	// separately — 22 px against the push button's 24 and the list row's 20
	// in the stored captures.
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
	// A recess is not one of these. The search field's chrome variant is the
	// platform's flat recess and takes FieldHeight, not this.
	ToolbarControlHeight float32
	// PaddingX is the horizontal inner padding of a control in dp.
	PaddingX float32
	// PaddingY is the vertical inner padding of a control in dp.
	PaddingY float32
}

// ChipHeight returns the chip height in dp: [Density.ControlHeight] less
// [ChipDrop]. A chip is smaller than a button, and the relation says that once
// instead of pinning a second scale that can drift off the first —
// Comfortable lands on 20 and Compact on 15.
//
// It is a method rather than a struct field so no Density value can carry a
// chip height that has come loose from its control height. The chip's pointer
// target is the chip: what it draws is what a pointer lands on.
func (d Density) ChipHeight() float32 { return d.ControlHeight - ChipDrop }

// The padding comes from the platform beside the control height: PaddingX is
// the HIG's inset beside a push button's label, PaddingY is what is left of
// the control height once LabelLarge's 20 dp line box has taken its share —
// 2 dp at Comfortable, nothing at Compact, whose line box is already over the
// height. PaddingX did not move with the measured heights: both captured
// buttons sit at the platform's minimum push-button width with their labels
// centred, so the capture measures centring and not an inset. See the
// provenance table at the top of this file.
var (
	// Comfortable is the default desktop density.
	Comfortable = Density{ControlHeight: ComfortableControlHeight, FieldHeight: ComfortableFieldHeight, RowHeight: ComfortableRowHeight, CheckboxRowHeight: ComfortableCheckboxRowHeight, ToolbarControlHeight: ComfortableToolbarControlHeight, PaddingX: 8, PaddingY: 2}
	// Compact is the dense mode: smaller drawn controls, tighter padding.
	Compact = Density{ControlHeight: CompactControlHeight, FieldHeight: CompactFieldHeight, RowHeight: CompactRowHeight, CheckboxRowHeight: CompactCheckboxRowHeight, ToolbarControlHeight: CompactToolbarControlHeight, PaddingX: 7, PaddingY: 0}
)
