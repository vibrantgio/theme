package tokens

// Desktop density is the platform's control scale: Comfortable is the
// platform's regular control and Compact its small one. Nothing below is
// derived from a web or Material scale.
//
// Provenance, per number — the reference and its gap are
// `.github/reference/macos/controls.md`, indexed by ADR-019:
//
//	number                      value   provenance
//	------                      -----   ----------
//	Comfortable control height  22 dp   published: the macOS HIG's regular push button, text field and pop-up button
//	Compact control height      19 dp   published: the macOS HIG's small push button
//	Comfortable PaddingX         8 dp   published: the HIG's horizontal inset beside a regular push button's label
//	Compact PaddingX             7 dp   derived: 8 × 19/22 rounded — no small-size inset is published
//	Comfortable PaddingY         1 dp   derived: (22 − LabelLarge's 20 dp line box) / 2, which lands a button exactly on 22
//	Compact PaddingY             0 dp   derived: LabelLarge's line box is already over 19, so there is nothing to pad with
//
// The stored captures hold no push button, text field or pop-up at regular
// size: what they hold is the unified toolbar's controls, 36 px across five
// applications, and Finder's info-pane text field at 33 px. Those are toolbar
// and pane controls, not the regular control this scale is named after. Both
// readings are recorded beside the published numbers, with a gap row asking
// for the one capture that would replace "published" with "measured" here.
//
// A control height is a floor, not a height. A control is as tall as its
// content box needs, and never shorter than the density says:
//
//	height = max(ControlHeight, lineBox + 2×PaddingY)
//
// where lineBox is the type role's line height (see [TextStyle.LineHeight]
// and theme/typeset). Typography does not move with density, so at the
// platform's heights the two terms are close and either can win:
//
//	control                  role         line box   + 2×PaddingY   ControlHeight   drawn
//	-------                  ----         --------   ------------   -------------   -----
//	button, Comfortable      LabelLarge   20         22             22              22
//	button, Compact          LabelLarge   20         20             19              20
//	text field, Comfortable  BodyLarge    24         26             22              26
//	text field, Compact      BodyLarge    24         24             19              24
//
// Comfortable's padding is what makes a button exact: 20 + 2×1 = 22. Compact
// has no padding left to spend — LabelLarge's 20 dp line box is 1 dp over the
// 19 dp floor before any padding — so a Compact button draws 20 dp and the
// floor is not the answer. That 1 dp closes only by moving a typography role,
// which this scale does not do.
//
// The consequence worth stating: controls in different type roles come out at
// different heights, and a Comfortable text field (26) is taller than a
// Comfortable button (22) because BodyLarge is a larger role than LabelLarge.
// Both are honest readings of the tokens. A design that wants them equal
// changes the roles or the padding, not the measurement.
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
//	list row                         22            19        pinned to ControlHeight
//	table body row and header cell   22            19        pinned to ControlHeight
//	sidebar item                     22            19        pinned to ControlHeight
//	picker option row                26            24        floor formula, BodyLarge
//
// At the platform's control heights no pinned row reaches WCAG 2.5.8's 24 dp:
// a Comfortable row is 22 and a Compact row 19, and a row's pointer target is
// the row itself. That is what taking the platform's scale costs, and it is
// recorded here rather than left to be discovered; the standalone floor is
// untouched, so a
// button, a checkbox, a chip and a closed picker still extend to 44 dp at
// every density. An application that must claim 2.5.8 for its rows sets a
// taller row of its own.

const (
	// ComfortableControlHeight is the default desktop control-height floor in
	// dp: a Comfortable control is at least this tall, and taller when its
	// content box needs it. It is the platform's regular control — the HIG's
	// 22 pt push button, text field and pop-up button.
	ComfortableControlHeight float32 = 22
	// CompactControlHeight is the dense-mode control-height floor in dp: the
	// platform's small control, the HIG's 19 pt push button. A control set in
	// LabelLarge or BodyLarge draws over it — see the table above — so for
	// those it is never the answer, and for a stacked row it is a pin rather
	// than a floor.
	CompactControlHeight float32 = 19
	// MinHitTarget is the pointer-target floor in dp for a *standalone*
	// control — one with space around it: button, checkbox, radio, text
	// field, a picker's closed trigger. Those extend their pointer area to at
	// least this on each axis, centred on the drawn control, whatever the
	// density.
	//
	// It is 44 dp, WCAG 2.5.5 Target Size (Enhanced), which is a AAA
	// criterion. It is not what stacked rows guarantee: list rows, table rows
	// and header cells, and open picker option rows are their own row height
	// (19 dp at Compact) because extending one row would steal its
	// neighbour's slop, and at the platform's control heights that is under
	// WCAG 2.5.8 Target Size (Minimum), the 24 dp criterion that governs at
	// AA. See "Pointer targets: which WCAG level actually governs" above.
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
	ControlHeight float32
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
// Comfortable lands on 18 and Compact on 15.
//
// It is a method for the same reason [Density.MinHitTarget] is one: no Density
// value can carry a chip height that has come loose from its control height.
// The pointer target is unaffected — a chip is a standalone control and
// extends to MinHitTarget like every other.
func (d Density) ChipHeight() float32 { return d.ControlHeight - ChipDrop }

// The padding comes from the platform beside the control height: PaddingX is
// the HIG's inset beside a push button's label, PaddingY is what is left of
// the control height once LabelLarge's 20 dp line box has taken its share —
// 1 dp at Comfortable, nothing at Compact, whose line box is already over the
// height. See the provenance table at the top of this file.
var (
	// Comfortable is the default desktop density.
	Comfortable = Density{ControlHeight: ComfortableControlHeight, PaddingX: 8, PaddingY: 1}
	// Compact is the dense mode: smaller drawn controls, same hit target.
	Compact = Density{ControlHeight: CompactControlHeight, PaddingX: 7, PaddingY: 0}
)
