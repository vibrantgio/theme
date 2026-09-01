package tokens

import "testing"

// TestDensityPicksWithinTableBounds pins the picked control heights to the
// measured table in density.go: Compact is denser than Comfortable, and both
// stay within [28, 44] — at or above macOS's large control (the densest
// desktop reference) and below the 44 dp touch-target floor.
func TestDensityPicksWithinTableBounds(t *testing.T) {
	if CompactControlHeight >= ComfortableControlHeight {
		t.Errorf("CompactControlHeight (%v) must be < ComfortableControlHeight (%v)",
			CompactControlHeight, ComfortableControlHeight)
	}
	for name, v := range map[string]float32{
		"ComfortableControlHeight": ComfortableControlHeight,
		"CompactControlHeight":     CompactControlHeight,
	} {
		if v < 28 || v > 44 {
			t.Errorf("%s = %v, want within [28, 44]", name, v)
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

// TestStackedRowsClearTheAAMinimumTarget makes density.go's pointer-target
// section checkable instead of merely readable. The claim it pins:
// stacked rows — list rows, table rows and header cells, sidebar items — are
// not extended to MinHitTarget, because adjacent rows would steal each other's
// slop. So the narrowest pointer target in the system is a Compact row at
// CompactControlHeight, and what that has to clear is WCAG 2.5.8 Target Size
// (Minimum) at 24 dp, the criterion that governs at AA — not MinHitTarget's
// 44 dp, which is 2.5.5 Target Size (Enhanced) and AAA.
//
// The third assertion is the one that would catch a well-meaning "fix": if
// Compact rows ever did reach 44, Compact would have stopped being compact and
// the doc would be describing something else.
func TestStackedRowsClearTheAAMinimumTarget(t *testing.T) {
	const (
		targetSizeMinimumAA   float32 = 24 // WCAG 2.5.8 Target Size (Minimum)
		targetSizeEnhancedAAA float32 = 44 // WCAG 2.5.5 Target Size (Enhanced)
	)
	if CompactControlHeight < targetSizeMinimumAA {
		t.Errorf("CompactControlHeight = %v, below WCAG 2.5.8's %v dp: the densest row in the system would fail AA",
			CompactControlHeight, targetSizeMinimumAA)
	}
	if MinHitTarget != targetSizeEnhancedAAA {
		t.Errorf("MinHitTarget = %v, want %v (WCAG 2.5.5 Target Size (Enhanced), AAA)",
			MinHitTarget, targetSizeEnhancedAAA)
	}
	if CompactControlHeight >= targetSizeEnhancedAAA {
		t.Errorf("CompactControlHeight = %v, at or above %v: rows are documented as clearing AA and not reaching AAA, so density.go needs rewriting before this passes",
			CompactControlHeight, targetSizeEnhancedAAA)
	}
}

// TestDensitySettingsMatchTable pins the two settings to the measured
// table in density.go: control heights are exactly the picked consts (so within the table's
// [28, 44] bounds already asserted above), Compact is strictly denser than
// Comfortable on every visual axis, and the paddings are the shadcn-derived
// pairs documented on the vars.
func TestDensitySettingsMatchTable(t *testing.T) {
	if Comfortable.ControlHeight != ComfortableControlHeight {
		t.Errorf("Comfortable.ControlHeight = %v, want %v", Comfortable.ControlHeight, ComfortableControlHeight)
	}
	if Compact.ControlHeight != CompactControlHeight {
		t.Errorf("Compact.ControlHeight = %v, want %v", Compact.ControlHeight, CompactControlHeight)
	}
	if Compact.ControlHeight >= Comfortable.ControlHeight {
		t.Errorf("Compact.ControlHeight (%v) must be < Comfortable.ControlHeight (%v)",
			Compact.ControlHeight, Comfortable.ControlHeight)
	}
	if Compact.PaddingX >= Comfortable.PaddingX || Compact.PaddingY >= Comfortable.PaddingY {
		t.Errorf("Compact padding (%v, %v) must be < Comfortable padding (%v, %v)",
			Compact.PaddingX, Compact.PaddingY, Comfortable.PaddingX, Comfortable.PaddingY)
	}
	if Comfortable.PaddingX != 16 || Comfortable.PaddingY != 8 {
		t.Errorf("Comfortable padding = (%v, %v), want (16, 8) — shadcn px-4 py-2",
			Comfortable.PaddingX, Comfortable.PaddingY)
	}
	if Compact.PaddingX != 12 || Compact.PaddingY != 6 {
		t.Errorf("Compact padding = (%v, %v), want (12, 6) — shadcn sm px-3, 2:1 vertical",
			Compact.PaddingX, Compact.PaddingY)
	}
}

// TestChipHeightRelation pins the chip height to the relation and not to two
// picked numbers: it is exactly ChipDrop under whatever control height the
// setting carries, which lands Comfortable on 32 and Compact on 24. Changing
// a control height moves the chip with it, and this test is what says so.
//
// The last assertion is the one that keeps the relation payable: the densest
// chip is the shortest thing the system draws, and 24 dp is exactly WCAG
// 2.5.8 Target Size (Minimum). A deeper drop would put the drawn control
// under the AA criterion even though the pointer target stays at MinHitTarget.
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
	if Comfortable.ChipHeight() != 32 || Compact.ChipHeight() != 24 {
		t.Errorf("chip heights = (%v, %v), want (32, 24)",
			Comfortable.ChipHeight(), Compact.ChipHeight())
	}
	const targetSizeMinimumAA float32 = 24 // WCAG 2.5.8 Target Size (Minimum)
	if Compact.ChipHeight() < targetSizeMinimumAA {
		t.Errorf("Compact.ChipHeight() = %v, under WCAG 2.5.8's %v dp: the shortest drawn control in the system would fail AA",
			Compact.ChipHeight(), targetSizeMinimumAA)
	}
}
