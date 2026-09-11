package tokens

// Desktop density is the platform's control scale: Comfortable is the
// platform's regular control and Compact its small one. Nothing below is
// derived from a web or Material scale.
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
//
// What is measured and what is not. One Save panel, captured at 1x in both
// appearances, holds the regular push button, pop-up button, text field and
// checkbox this scale is named after, and every Comfortable number above is
// read off its pixels. No capture holds a control at the platform's small
// size, so Compact's height is still the published 19 and its field height
// is derived from the measured ratio; ADR-019 carries that gap row. The
// stored windows also measure the unified toolbar's controls at 36 px and
// Finder's info-pane text field at 33 px — different controls in different
// places, recorded beside these and correcting nothing here.
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
// ControlHeight is a floor for controls but a *pin* for stacked rows: list
// rows, table rows, header cells and sidebar items are ControlHeight tall
// exactly (see the row table below), so a change here re-pitches every dense
// list and table in the system.
//
// # Pointer targets: which WCAG level actually governs
//
// [MinHitTarget] is 44 dp, and 44 dp is not the AA requirement. Two success
// criteria are in play and they are a whole conformance level apart:
//
//	criterion                            level   threshold   applies to
//	---------                            -----   ---------   ----------
//	WCAG 2.5.5 Target Size (Enhanced)    AAA     44×44 CSS px  every pointer target
//	WCAG 2.5.8 Target Size (Minimum)     AA      24×24 CSS px  every pointer target
//
// (WCAG 2.2, https://www.w3.org/TR/WCAG22/#target-size-enhanced and
// #target-size-minimum. Both carry an inline/essential exception this system
// does not need to lean on.)
//
// The pointer area is 44 dp for standalone controls — button, checkbox,
// radio, text field, a picker's closed trigger — and deliberately not for
// stacked rows: list rows, table rows and header cells, and open picker
// option rows. Adjacent rows tile edge to edge, so slop granted to one row is
// stolen from its neighbour; the extension would not enlarge anything, it
// would only make the boundary lie about where it is. Rows rely on their full
// row width instead.
//
// So the stacked-row targets are as tall as the row is — ControlHeight is a
// floor, so a row whose content box is taller draws more than the token says:
//
//	row                              Comfortable   Compact   sizing
//	---                              -----------   -------   ------
//	list row                         24            19        pinned to ControlHeight
//	table body row and header cell   24            19        pinned to ControlHeight
//	sidebar item                     24            19        pinned to ControlHeight
//	picker option row                28            24        floor formula, BodyLarge
//
// At the platform's measured regular control a Comfortable pinned row is 24
// dp and so meets WCAG 2.5.8 Target Size (Minimum) exactly, the 24 dp
// criterion that governs at AA; a Compact pinned row is 19 and does not, and
// a row's pointer target is the row itself. That is what taking the
// platform's small control costs, and it is recorded here rather than left to
// be discovered; the standalone floor is untouched, so a button, a checkbox,
// a chip and a closed picker still extend to 44 dp at every density. An
// application that must claim 2.5.8 for its Compact rows sets a taller row of
// its own.

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
	// MinHitTarget is the pointer-target floor in dp for a *standalone*
	// control — one with space around it: button, checkbox, radio, text
	// field, a picker's closed trigger. Those extend their pointer area to at
	// least this on each axis, centred on the drawn control, whatever the
	// density.
	//
	// It is 44 dp, WCAG 2.5.5 Target Size (Enhanced), which is a AAA
	// criterion. It is not what stacked rows guarantee: list rows, table rows
	// and header cells, and open picker option rows are their own row height
	// because extending one row would steal its neighbour's slop. At the
	// platform's measured regular control a Comfortable row is 24 dp and
	// meets WCAG 2.5.8 Target Size (Minimum), the 24 dp criterion that
	// governs at AA; a Compact row is 19 dp and does not. See "Pointer
	// targets: which WCAG level actually governs" above.
	MinHitTarget float32 = 44
	// ChipDrop is how far under the control height the system's smallest
	// control is drawn, in dp. See [Density.ChipHeight]: it is the whole of
	// that relation, exported so a reader can see the two heights are one
	// number apart rather than two picks.
	ChipDrop float32 = 4
)

// Density is one density setting: the drawn control height and its inner
// padding, all in dp. It is a comparable value struct like the other token
// types. The standalone-control pointer-target floor is deliberately a method,
// not a field — see [Density.MinHitTarget] — so no Density value can carry a
// shrunken hit target.
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
	// PaddingX is the horizontal inner padding of a control in dp.
	PaddingX float32
	// PaddingY is the vertical inner padding of a control in dp.
	PaddingY float32
}

// MinHitTarget returns the standalone-control pointer-target floor in dp —
// the package const [MinHitTarget], 44 dp, WCAG 2.5.5 Target Size (Enhanced).
// It is a method rather than a struct field, so it is structurally identical
// across every density: Compact shrinks the drawn control, never the clickable
// area of a control that has room to grow into.
//
// It does not describe stacked rows. Read [MinHitTarget] before wiring this
// into anything that tiles.
func (Density) MinHitTarget() float32 { return MinHitTarget }

// ChipHeight returns the chip height in dp: [Density.ControlHeight] less
// [ChipDrop]. A chip is smaller than a button, and the relation says that once
// instead of pinning a second scale that can drift off the first —
// Comfortable lands on 20 and Compact on 15.
//
// It is a method for the same reason [Density.MinHitTarget] is one: no Density
// value can carry a chip height that has come loose from its control height.
// The pointer target is unaffected — a chip is a standalone control and
// extends to MinHitTarget like every other.
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
	Comfortable = Density{ControlHeight: ComfortableControlHeight, FieldHeight: ComfortableFieldHeight, PaddingX: 8, PaddingY: 2}
	// Compact is the dense mode: smaller drawn controls, same hit target.
	Compact = Density{ControlHeight: CompactControlHeight, FieldHeight: CompactFieldHeight, PaddingX: 7, PaddingY: 0}
)
