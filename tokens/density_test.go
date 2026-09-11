package tokens

import "testing"

// TestDensityPicksMatchThePlatformScale pins the two control heights to the
// platform's published scale, recorded with its provenance in density.go and
// in .github/reference/macos/controls.md: Comfortable is the regular control
// (22 pt) and Compact the small one (19 pt), so both sit inside the HIG's
// published push-button range — mini 16 to regular 22 — and well under the
// 44 dp pointer-target floor, which is not a control height.
func TestDensityPicksMatchThePlatformScale(t *testing.T) {
	const (
		platformRegular float32 = 22 // HIG: regular push button, text field, pop-up button
		platformSmall   float32 = 19 // HIG: small push button
		platformMini    float32 = 16 // HIG: mini push button
	)
	if ComfortableControlHeight != platformRegular {
		t.Errorf("ComfortableControlHeight = %v, want %v (the platform's regular control)",
			ComfortableControlHeight, platformRegular)
	}
	if CompactControlHeight != platformSmall {
		t.Errorf("CompactControlHeight = %v, want %v (the platform's small control)",
			CompactControlHeight, platformSmall)
	}
	if CompactControlHeight >= ComfortableControlHeight {
		t.Errorf("CompactControlHeight (%v) must be < ComfortableControlHeight (%v)",
			CompactControlHeight, ComfortableControlHeight)
	}
	for name, v := range map[string]float32{
		"ComfortableControlHeight": ComfortableControlHeight,
		"CompactControlHeight":     CompactControlHeight,
	} {
		if v < platformMini || v > platformRegular {
			t.Errorf("%s = %v, want within the platform's published range [%v, %v]",
				name, v, platformMini, platformRegular)
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

// TestStackedRowsSitUnderTheAAMinimumTarget makes density.go's pointer-target
// section checkable instead of merely readable, and records what taking the
// platform's control heights costs. The claim it pins: stacked rows — list
// rows, table rows and header cells, sidebar items — are pinned to
// ControlHeight and are not extended to MinHitTarget, because adjacent rows
// would steal each other's slop. At the platform's heights that puts every
// pinned row under WCAG 2.5.8 Target Size (Minimum), the 24 dp criterion that
// governs at AA, while the standalone floor stays at 2.5.5's 44 dp for every
// control that has room around it.
//
// The assertions fail in either direction on purpose: if a row height ever
// reaches 24, density.go's recorded consequence is stale and has to be
// rewritten rather than quietly outgrown.
func TestStackedRowsSitUnderTheAAMinimumTarget(t *testing.T) {
	const (
		targetSizeMinimumAA   float32 = 24 // WCAG 2.5.8 Target Size (Minimum)
		targetSizeEnhancedAAA float32 = 44 // WCAG 2.5.5 Target Size (Enhanced)
	)
	for name, d := range map[string]Density{
		"Comfortable": Comfortable,
		"Compact":     Compact,
	} {
		if d.ControlHeight >= targetSizeMinimumAA {
			t.Errorf("%s row height = %v, at or above WCAG 2.5.8's %v dp: density.go records every pinned row as sitting under it, so that section needs rewriting before this passes",
				name, d.ControlHeight, targetSizeMinimumAA)
		}
		if d.MinHitTarget() != targetSizeEnhancedAAA {
			t.Errorf("%s.MinHitTarget() = %v, want %v (WCAG 2.5.5 Target Size (Enhanced), AAA): a standalone control keeps its floor whatever the row does",
				name, d.MinHitTarget(), targetSizeEnhancedAAA)
		}
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
	if Comfortable.PaddingX != 8 || Comfortable.PaddingY != 1 {
		t.Errorf("Comfortable padding = (%v, %v), want (8, 1) — the HIG's inset beside a regular push button's label, and what 22 leaves around LabelLarge's line box",
			Comfortable.PaddingX, Comfortable.PaddingY)
	}
	if Compact.PaddingX != 7 || Compact.PaddingY != 0 {
		t.Errorf("Compact padding = (%v, %v), want (7, 0) — 8 scaled by 19/22, and no vertical room left at 19",
			Compact.PaddingX, Compact.PaddingY)
	}
}

// TestChipHeightRelation pins the chip height to the relation and not to two
// picked numbers: it is exactly ChipDrop under whatever control height the
// setting carries, which lands Comfortable on 18 and Compact on 15. Changing
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
	if Comfortable.ChipHeight() != 18 || Compact.ChipHeight() != 15 {
		t.Errorf("chip heights = (%v, %v), want (18, 15)",
			Comfortable.ChipHeight(), Compact.ChipHeight())
	}
	const targetSizeMinimumAA float32 = 24 // WCAG 2.5.8 Target Size (Minimum)
	if Compact.MinHitTarget() < targetSizeMinimumAA {
		t.Errorf("Compact.MinHitTarget() = %v, under WCAG 2.5.8's %v dp: a chip is a standalone control and its pointer target is what the criterion measures",
			Compact.MinHitTarget(), targetSizeMinimumAA)
	}
}

// TestControlFitsItsLabelLine states the fit rule against the typography the
// controls actually wear: a control draws max(ControlHeight, lineBox +
// 2×PaddingY), where lineBox is the role's line height, which density does not
// move. A button is set in LabelLarge (20 dp line box) and a text field in
// BodyLarge (24 dp), so this is what the two densities produce:
//
//   - Comfortable lands a button exactly on the platform's regular control:
//     20 + 2×1 = 22, which is what PaddingY 1 is for.
//   - Compact cannot: LabelLarge's line box alone is 1 dp over the platform's
//     small control, so PaddingY is 0 — its minimum, never negative — and the
//     button draws 20 against a 19 dp floor. Closing that dp means moving a
//     typography role, which this scale does not do.
func TestControlFitsItsLabelLine(t *testing.T) {
	drawn := func(d Density, lineBox float32) float32 {
		h := lineBox + 2*d.PaddingY
		if d.ControlHeight > h {
			return d.ControlHeight
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
	for _, tc := range []struct {
		name    string
		d       Density
		lineBox float32
		want    float32
	}{
		{"button, Comfortable", Comfortable, button, 22},
		{"button, Compact", Compact, button, 20},
		{"text field, Comfortable", Comfortable, field, 26},
		{"text field, Compact", Compact, field, 24},
	} {
		if got := drawn(tc.d, tc.lineBox); got != tc.want {
			t.Errorf("%s draws %v dp, want %v — the table above the constants is stale", tc.name, got, tc.want)
		}
	}
}
