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
// published small control at 19 dp because no capture holds one. Both sit
// well under the 44 dp pointer-target floor, which is not a control height.
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
	if ComfortableControlHeight >= MinHitTarget {
		t.Errorf("ComfortableControlHeight (%v) should sit below the hit-target floor (%v): the floor, not the control, is what 44 dp is for",
			ComfortableControlHeight, MinHitTarget)
	}
	if MinHitTarget != 44 {
		t.Errorf("MinHitTarget = %v, want 44 (WCAG 2.5.5, components' current constant)", MinHitTarget)
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

// TestDensityHitTargetFloor asserts that both
// density settings satisfy the WCAG 2.5.5 pointer-target floor, and the floor
// is identical across settings — Compact shrinks the drawn control, never the
// clickable area.
func TestDensityHitTargetFloor(t *testing.T) {
	for name, d := range map[string]Density{
		"Comfortable": Comfortable,
		"Compact":     Compact,
	} {
		if d.MinHitTarget() < 44 {
			t.Errorf("%s.MinHitTarget() = %v, want >= 44 (WCAG 2.5.5)", name, d.MinHitTarget())
		}
	}
	if Comfortable.MinHitTarget() != Compact.MinHitTarget() {
		t.Errorf("hit target must not vary with density: Comfortable %v != Compact %v",
			Comfortable.MinHitTarget(), Compact.MinHitTarget())
	}
}

// TestStackedRowsAgainstTheAAMinimumTarget makes density.go's pointer-target
// section checkable instead of merely readable, and records what taking the
// platform's control heights costs. The claim it pins: stacked rows — list
// rows, table rows and header cells, sidebar items — are pinned to
// ControlHeight and are not extended to MinHitTarget, because adjacent rows
// would steal each other's slop. At the platform's measured regular control a
// Comfortable row is exactly WCAG 2.5.8 Target Size (Minimum), the 24 dp
// criterion that governs at AA; Compact's uncaptured small control leaves its
// rows under it. The standalone floor stays at 2.5.5's 44 dp for every
// control that has room around it.
//
// The assertions fail in either direction on purpose: if either row height
// moves across the criterion, density.go's recorded consequence is stale and
// has to be rewritten rather than quietly outgrown.
func TestStackedRowsAgainstTheAAMinimumTarget(t *testing.T) {
	const (
		targetSizeMinimumAA   float32 = 24 // WCAG 2.5.8 Target Size (Minimum)
		targetSizeEnhancedAAA float32 = 44 // WCAG 2.5.5 Target Size (Enhanced)
	)
	if Comfortable.ControlHeight != targetSizeMinimumAA {
		t.Errorf("Comfortable row height = %v, want exactly WCAG 2.5.8's %v dp: density.go records the measured regular control as landing on the criterion, so that section needs rewriting before this passes",
			Comfortable.ControlHeight, targetSizeMinimumAA)
	}
	if Compact.ControlHeight >= targetSizeMinimumAA {
		t.Errorf("Compact row height = %v, at or above WCAG 2.5.8's %v dp: density.go records the Compact row as sitting under it, so that section needs rewriting before this passes",
			Compact.ControlHeight, targetSizeMinimumAA)
	}
	for name, d := range map[string]Density{
		"Comfortable": Comfortable,
		"Compact":     Compact,
	} {
		if d.MinHitTarget() != targetSizeEnhancedAAA {
			t.Errorf("%s.MinHitTarget() = %v, want %v (WCAG 2.5.5 Target Size (Enhanced), AAA): a standalone control keeps its floor whatever the row does",
				name, d.MinHitTarget(), targetSizeEnhancedAAA)
		}
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
// The drawn chip is the shortest thing the system draws and is well under
// WCAG 2.5.8's 24 dp at both densities; what that criterion measures is the
// pointer target, and a chip is a standalone control, so it extends to
// MinHitTarget like every other — asserted per setting below.
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
		if d.MinHitTarget() != MinHitTarget {
			t.Errorf("%s.MinHitTarget() = %v, want %v: the drop is a drawn height and never a pointer target",
				name, d.MinHitTarget(), MinHitTarget)
		}
	}
	if Comfortable.ChipHeight() != 20 || Compact.ChipHeight() != 15 {
		t.Errorf("chip heights = (%v, %v), want (20, 15)",
			Comfortable.ChipHeight(), Compact.ChipHeight())
	}
	const targetSizeMinimumAA float32 = 24 // WCAG 2.5.8 Target Size (Minimum)
	if Compact.MinHitTarget() < targetSizeMinimumAA {
		t.Errorf("Compact.MinHitTarget() = %v, under WCAG 2.5.8's %v dp: a chip is a standalone control and its pointer target is what the criterion measures",
			Compact.MinHitTarget(), targetSizeMinimumAA)
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
