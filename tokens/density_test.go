package tokens

import (
	"math"
	"testing"
)

// TestDensityPicksMatchThePlatformScale pins the two control heights to the
// platform's scale as the reference now reads it, with the provenance
// recorded in density.go and in .github/reference/macos/controls.md:
// Comfortable is the regular control measured off the stored Save panel at
// 24 dp, superseding the HIG's published 22 pt, and Compact is still the
// published small control at 19 dp because no capture holds one. Each is also
// the pointer target a control of that density offers.
func TestDensityPicksMatchThePlatformScale(t *testing.T) {
	const (
		measuredRegular  float32 = 24 // MEASURED: save-dialog-{light,dark}.png, push button and pop-up button
		publishedRegular float32 = 22 // PUBLISHED: the HIG's regular push button — superseded by the capture
		publishedSmall   float32 = 19 // PUBLISHED: the HIG's small push button — uncaptured
		publishedMini    float32 = 16 // PUBLISHED: the HIG's mini push button
	)
	if ComfortableControlHeight != measuredRegular {
		t.Errorf("ComfortableControlHeight = %v, want %v (the platform's regular control, measured)",
			ComfortableControlHeight, measuredRegular)
	}
	if ComfortableControlHeight == publishedRegular {
		t.Errorf("ComfortableControlHeight = %v, the published number: measured beats published here (owner ruling 2026-09-11)",
			ComfortableControlHeight)
	}
	if CompactControlHeight != publishedSmall {
		t.Errorf("CompactControlHeight = %v, want %v (the platform's small control, published)",
			CompactControlHeight, publishedSmall)
	}
	if CompactControlHeight >= ComfortableControlHeight {
		t.Errorf("CompactControlHeight (%v) must be < ComfortableControlHeight (%v)",
			CompactControlHeight, ComfortableControlHeight)
	}
	for name, v := range map[string]float32{
		"ComfortableControlHeight": ComfortableControlHeight,
		"CompactControlHeight":     CompactControlHeight,
	} {
		if v < publishedMini || v > measuredRegular {
			t.Errorf("%s = %v, want within the platform's control range [%v, %v]",
				name, v, publishedMini, measuredRegular)
		}
	}
}

// TestFieldHeightIsItsOwnMeasurement pins the text field's height to the
// platform's own reading of it. The Save panel measures the field at 27 px
// against the push button's 24, so a field that took the control height
// would be 3 dp short of the platform; Compact's 21 is that measured ratio
// applied to the uncaptured small control, and is what a capture of one
// replaces.
func TestFieldHeightIsItsOwnMeasurement(t *testing.T) {
	const (
		measuredField float32 = 27 // MEASURED: save-dialog-{light,dark}.png, the "Tags:" field
		derivedSmall  float32 = 21 // DERIVED: 27 × 19/24 = 21.4, rounded
	)
	if ComfortableFieldHeight != measuredField {
		t.Errorf("ComfortableFieldHeight = %v, want %v (the platform's measured text field)",
			ComfortableFieldHeight, measuredField)
	}
	if CompactFieldHeight != derivedSmall {
		t.Errorf("CompactFieldHeight = %v, want %v (27 × 19/24, rounded)",
			CompactFieldHeight, derivedSmall)
	}
	if got, want := float32(math.Round(float64(ComfortableFieldHeight*CompactControlHeight/ComfortableControlHeight))), CompactFieldHeight; got != want {
		t.Errorf("the Compact field height is %v, but the platform's ratio gives %v: the derivation in density.go is stale", want, got)
	}
	for name, d := range map[string]Density{
		"Comfortable": Comfortable,
		"Compact":     Compact,
	} {
		if d.FieldHeight <= d.ControlHeight {
			t.Errorf("%s.FieldHeight = %v, not above ControlHeight %v: the platform draws a field taller than a button, which is why this is a second number",
				name, d.FieldHeight, d.ControlHeight)
		}
	}
}

// TestToolbarControlIsItsOwnMeasurement pins the toolbar control's height to
// the platform's own reading of it. Every bordered control in the Finder
// toolbar captures measures 36 px — the frontmost untinted dark window's five
// controls between their rim rows, the inactive light window's fill, and the
// view pop-up in both frontmost captures — where the same window's dialog
// pop-up measures 24. The two are different controls in different places, so
// the toolbar's is a number of its own and neither reading corrects the other.
//
// Compact carries the same 36: no capture holds a toolbar drawn at the
// platform's small size, which is the one gap this number has.
func TestToolbarControlIsItsOwnMeasurement(t *testing.T) {
	const measuredToolbar float32 = 36 // MEASURED: the Finder toolbar captures, every bordered control
	if ComfortableToolbarControlHeight != measuredToolbar {
		t.Errorf("ComfortableToolbarControlHeight = %v, want %v (the platform's toolbar control, measured)",
			ComfortableToolbarControlHeight, measuredToolbar)
	}
	if CompactToolbarControlHeight != measuredToolbar {
		t.Errorf("CompactToolbarControlHeight = %v, want %v (carried until a capture holds a small toolbar)",
			CompactToolbarControlHeight, measuredToolbar)
	}
	for name, d := range map[string]Density{"Comfortable": Comfortable, "Compact": Compact} {
		if d.ToolbarControlHeight != measuredToolbar {
			t.Errorf("%s.ToolbarControlHeight = %v, want %v", name, d.ToolbarControlHeight, measuredToolbar)
		}
		if d.ToolbarControlHeight <= d.ControlHeight {
			t.Errorf("%s.ToolbarControlHeight = %v, not above ControlHeight %v: the platform draws a toolbar control taller than a dialog's, which is why this is a second number",
				name, d.ToolbarControlHeight, d.ControlHeight)
		}
		if d.ToolbarControlHeight <= d.FieldHeight {
			t.Errorf("%s.ToolbarControlHeight = %v, not above FieldHeight %v", name, d.ToolbarControlHeight, d.FieldHeight)
		}
	}
}

// TestStackedRowsAreTheirOwnTarget makes density.go's pointer-target section
// checkable instead of merely readable. The claim it pins: a stacked row —
// a list row, a table row, a header cell, a sidebar item — is its own target
// at RowHeight, the platform's measured list row, and never the control
// height, because adjacent rows would otherwise steal each other's pixels.
//
// The assertions fail in either direction on purpose: if a row height moves,
// density.go's recorded consequence is stale and has to be rewritten rather
// than quietly outgrown.
func TestStackedRowsAreTheirOwnTarget(t *testing.T) {
	if Comfortable.RowHeight != ComfortableRowHeight {
		t.Errorf("Comfortable.RowHeight = %v, want %v (Finder's list view, measured)",
			Comfortable.RowHeight, ComfortableRowHeight)
	}
	if Compact.RowHeight != CompactRowHeight {
		t.Errorf("Compact.RowHeight = %v, want %v", Compact.RowHeight, CompactRowHeight)
	}
	if Comfortable.RowHeight >= Comfortable.ControlHeight {
		t.Errorf("Comfortable row height %v is not under the control height %v: the platform draws a list row shorter than a button, which is why the row is a number of its own",
			Comfortable.RowHeight, Comfortable.ControlHeight)
	}
}

// TestDensitySettingsMatchTable pins the two settings to the provenance table
// in density.go: the heights are exactly the picked consts, Compact is
// strictly denser than Comfortable on every visual axis, and the paddings are
// the pair the table derives.
func TestDensitySettingsMatchTable(t *testing.T) {
	if Comfortable.ControlHeight != ComfortableControlHeight {
		t.Errorf("Comfortable.ControlHeight = %v, want %v", Comfortable.ControlHeight, ComfortableControlHeight)
	}
	if Compact.ControlHeight != CompactControlHeight {
		t.Errorf("Compact.ControlHeight = %v, want %v", Compact.ControlHeight, CompactControlHeight)
	}
	if Comfortable.FieldHeight != ComfortableFieldHeight {
		t.Errorf("Comfortable.FieldHeight = %v, want %v", Comfortable.FieldHeight, ComfortableFieldHeight)
	}
	if Compact.FieldHeight != CompactFieldHeight {
		t.Errorf("Compact.FieldHeight = %v, want %v", Compact.FieldHeight, CompactFieldHeight)
	}
	if Compact.ControlHeight >= Comfortable.ControlHeight || Compact.FieldHeight >= Comfortable.FieldHeight {
		t.Errorf("Compact heights (%v, %v) must be < Comfortable heights (%v, %v)",
			Compact.ControlHeight, Compact.FieldHeight, Comfortable.ControlHeight, Comfortable.FieldHeight)
	}
	if Compact.PaddingX >= Comfortable.PaddingX || Compact.PaddingY >= Comfortable.PaddingY {
		t.Errorf("Compact padding (%v, %v) must be < Comfortable padding (%v, %v)",
			Compact.PaddingX, Compact.PaddingY, Comfortable.PaddingX, Comfortable.PaddingY)
	}
	if Comfortable.PaddingX != 8 || Comfortable.PaddingY != 2 {
		t.Errorf("Comfortable padding = (%v, %v), want (8, 2) — the HIG's inset beside a regular push button's label, and what the measured 24 leaves around LabelLarge's line box",
			Comfortable.PaddingX, Comfortable.PaddingY)
	}
	if Compact.PaddingX != 7 || Compact.PaddingY != 0 {
		t.Errorf("Compact padding = (%v, %v), want (7, 0) — the published 8 scaled by 19/22, and no vertical room left at 19",
			Compact.PaddingX, Compact.PaddingY)
	}
}

// TestChipHeightRelation pins the chip height to the relation and not to two
// picked numbers: it is exactly ChipDrop under whatever control height the
// setting carries, which lands Comfortable on 20 and Compact on 15. Changing
// a control height moves the chip with it, and this test is what says so.
//
// The drawn chip is the shortest thing the system draws, and what it draws is
// what a pointer lands on: the chip's target is the chip.
func TestChipHeightRelation(t *testing.T) {
	for name, d := range map[string]Density{
		"Comfortable": Comfortable,
		"Compact":     Compact,
	} {
		if got, want := d.ChipHeight(), d.ControlHeight-ChipDrop; got != want {
			t.Errorf("%s.ChipHeight() = %v, want ControlHeight − ChipDrop = %v", name, got, want)
		}
		if d.ChipHeight() >= d.ControlHeight {
			t.Errorf("%s.ChipHeight() = %v, not under ControlHeight %v: a chip is smaller than a button",
				name, d.ChipHeight(), d.ControlHeight)
		}
	}
	if Comfortable.ChipHeight() != 20 || Compact.ChipHeight() != 15 {
		t.Errorf("chip heights = (%v, %v), want (20, 15)",
			Comfortable.ChipHeight(), Compact.ChipHeight())
	}
}

// TestControlFitsItsLabelLine states the fit rule against the typography the
// controls actually wear: a control draws max(floor, lineBox + 2×PaddingY),
// where the floor is ControlHeight for a button and FieldHeight for a text
// field, and lineBox is the role's line height, which density does not move.
// A button is set in LabelLarge (20 dp line box) and a text field in
// BodyLarge (24 dp), so this is what the two densities produce:
//
//   - Comfortable lands a button exactly on the platform's measured regular
//     control: 20 + 2×2 = 24, which is what PaddingY 2 is for.
//   - Compact cannot: LabelLarge's line box alone is 1 dp over the platform's
//     small control, so PaddingY is 0 — its minimum, never negative — and the
//     button draws 20 against a 19 dp floor.
//   - A Comfortable text field draws 28 against the platform's measured 27,
//     because BodyLarge's 24 dp line box plus the control's own padding is 28.
//
// Both overshoots close only by moving a typography role, which this scale
// does not do.
func TestControlFitsItsLabelLine(t *testing.T) {
	drawn := func(floor, lineBox, paddingY float32) float32 {
		h := lineBox + 2*paddingY
		if floor > h {
			return floor
		}
		return h
	}
	button := DefaultTypography.LabelLarge.LineHeight
	field := DefaultTypography.BodyLarge.LineHeight
	if button != 20 || field != 24 {
		t.Fatalf("line boxes = (LabelLarge %v, BodyLarge %v), want (20, 24): typography moved, so the density table above the constants is stale",
			button, field)
	}
	for name, d := range map[string]Density{
		"Comfortable": Comfortable,
		"Compact":     Compact,
	} {
		if d.PaddingY < 0 {
			t.Errorf("%s.PaddingY = %v: padding is never negative, the height is a floor instead", name, d.PaddingY)
		}
		// The fit rule itself, stated per density and per floor rather than
		// only through the worked rows below: a control is never shorter than
		// its floor, and never shorter than its own line box plus its padding.
		for what, pair := range map[string][2]float32{
			"button":     {d.ControlHeight, button},
			"text field": {d.FieldHeight, field},
		} {
			floor, lineBox := pair[0], pair[1]
			got := drawn(floor, lineBox, d.PaddingY)
			if got < floor {
				t.Errorf("%s %s draws %v, under its %v dp floor", name, what, got, floor)
			}
			if got < lineBox+2*d.PaddingY {
				t.Errorf("%s %s draws %v, under its content box %v + 2×%v", name, what, got, lineBox, d.PaddingY)
			}
		}
	}
	if got, want := button+2*Comfortable.PaddingY, ComfortableControlHeight; got != want {
		t.Errorf("Comfortable button = %v (LabelLarge %v + 2×%v), want exactly %v",
			got, button, Comfortable.PaddingY, want)
	}
	if Compact.PaddingY != 0 {
		t.Errorf("Compact.PaddingY = %v, want 0: LabelLarge's %v dp line box is already over the %v dp floor",
			Compact.PaddingY, button, CompactControlHeight)
	}
	if got, want := button-CompactControlHeight, float32(1); got != want {
		t.Errorf("Compact button overshoot = %v dp, want %v (LabelLarge %v over the %v dp floor)",
			got, want, button, CompactControlHeight)
	}
	if got, want := field+2*Comfortable.PaddingY-ComfortableFieldHeight, float32(1); got != want {
		t.Errorf("Comfortable text field overshoot = %v dp, want %v (BodyLarge %v + 2×%v over the measured %v)",
			got, want, field, Comfortable.PaddingY, ComfortableFieldHeight)
	}
	for _, tc := range []struct {
		name    string
		floor   float32
		lineBox float32
		pad     float32
		want    float32
	}{
		{"button, Comfortable", Comfortable.ControlHeight, button, Comfortable.PaddingY, 24},
		{"button, Compact", Compact.ControlHeight, button, Compact.PaddingY, 20},
		{"text field, Comfortable", Comfortable.FieldHeight, field, Comfortable.PaddingY, 28},
		{"text field, Compact", Compact.FieldHeight, field, Compact.PaddingY, 24},
	} {
		if got := drawn(tc.floor, tc.lineBox, tc.pad); got != tc.want {
			t.Errorf("%s draws %v dp, want %v — the table above the constants is stale", tc.name, got, tc.want)
		}
	}
}
