package tokens

import (
	"math"
	"testing"
)

// TestDensityPicksMatchThePlatformScale pins the two control heights to the
// platform's scale as the reference now reads it, with the provenance
// recorded in density.go and in organisation/reference/macos/controls.md:
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

// TestFieldHeightsAreTheirMeasuredRows pins each field's Comfortable height
// to its row in organisation/reference/macos/controls.md. A field's height is
// measured per control and per variant, so each one is its own reading.
func TestFieldHeightsAreTheirMeasuredRows(t *testing.T) {
	for _, tc := range []struct {
		name string
		got  float32
		row  float32
		from string
	}{
		{"text field in a form", Comfortable.TextFieldHeight, 27, `save-dialog-{light,dark}.png, the "Tags:" field, y 243–269`},
		{"search field on the content", Comfortable.SearchFieldHeight, 26, "systemsettings-add-input-source-search-{light,dark}.png, y 347–372"},
		{"search field in a sidebar", Comfortable.SidebarSearchFieldHeight, 28, "system-settings-grouped-box-{light,dark}.png, y 61–88"},
	} {
		if tc.got != tc.row {
			t.Errorf("Comfortable %s = %v dp, want its measured row %v (%s)", tc.name, tc.got, tc.row, tc.from)
		}
	}
}

// TestCompactFieldHeightsAreDerived pins each Compact field height to the
// derivation the reference states: the Comfortable reading times the small
// control's 19 over the regular control's 24, rounded. No capture holds a
// small field.
func TestCompactFieldHeightsAreDerived(t *testing.T) {
	derive := func(comfortable float32) float32 {
		return float32(math.Round(float64(comfortable * CompactControlHeight / ComfortableControlHeight)))
	}
	for _, tc := range []struct {
		name             string
		got, comfortable float32
	}{
		{"text field", Compact.TextFieldHeight, Comfortable.TextFieldHeight},
		{"search field", Compact.SearchFieldHeight, Comfortable.SearchFieldHeight},
		{"sidebar search field", Compact.SidebarSearchFieldHeight, Comfortable.SidebarSearchFieldHeight},
	} {
		if want := derive(tc.comfortable); tc.got != want {
			t.Errorf("Compact %s = %v dp, want %v (%v × 19/24, rounded)", tc.name, tc.got, want, tc.comfortable)
		}
	}
	for name, d := range map[string]Density{"Comfortable": Comfortable, "Compact": Compact} {
		for field, h := range map[string]float32{"TextFieldHeight": d.TextFieldHeight, "SearchFieldHeight": d.SearchFieldHeight, "SidebarSearchFieldHeight": d.SidebarSearchFieldHeight} {
			if h <= d.ControlHeight {
				t.Errorf("%s.%s = %v, not above ControlHeight %v: the platform draws a field taller than a button", name, field, h, d.ControlHeight)
			}
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
		for field, h := range map[string]float32{"TextFieldHeight": d.TextFieldHeight, "SearchFieldHeight": d.SearchFieldHeight, "SidebarSearchFieldHeight": d.SidebarSearchFieldHeight} {
			if d.ToolbarControlHeight <= h {
				t.Errorf("%s.ToolbarControlHeight = %v, not above %s %v", name, d.ToolbarControlHeight, field, h)
			}
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
// strictly denser than Comfortable on every visual axis, and PaddingX is the
// pair the table states.
func TestDensitySettingsMatchTable(t *testing.T) {
	if Comfortable.ControlHeight != ComfortableControlHeight {
		t.Errorf("Comfortable.ControlHeight = %v, want %v", Comfortable.ControlHeight, ComfortableControlHeight)
	}
	if Compact.ControlHeight != CompactControlHeight {
		t.Errorf("Compact.ControlHeight = %v, want %v", Compact.ControlHeight, CompactControlHeight)
	}
	if Comfortable.TextFieldHeight != ComfortableTextFieldHeight || Comfortable.SearchFieldHeight != ComfortableSearchFieldHeight || Comfortable.SidebarSearchFieldHeight != ComfortableSidebarSearchFieldHeight {
		t.Errorf("Comfortable field heights = (%v, %v, %v), want the Comfortable consts", Comfortable.TextFieldHeight, Comfortable.SearchFieldHeight, Comfortable.SidebarSearchFieldHeight)
	}
	if Compact.TextFieldHeight != CompactTextFieldHeight || Compact.SearchFieldHeight != CompactSearchFieldHeight || Compact.SidebarSearchFieldHeight != CompactSidebarSearchFieldHeight {
		t.Errorf("Compact field heights = (%v, %v, %v), want the Compact consts", Compact.TextFieldHeight, Compact.SearchFieldHeight, Compact.SidebarSearchFieldHeight)
	}
	if Compact.ControlHeight >= Comfortable.ControlHeight || Compact.TextFieldHeight >= Comfortable.TextFieldHeight || Compact.SearchFieldHeight >= Comfortable.SearchFieldHeight || Compact.SidebarSearchFieldHeight >= Comfortable.SidebarSearchFieldHeight {
		t.Errorf("Compact heights must be < Comfortable heights: %+v against %+v", Compact, Comfortable)
	}
	if Compact.PaddingX >= Comfortable.PaddingX {
		t.Errorf("Compact PaddingX %v must be < Comfortable PaddingX %v", Compact.PaddingX, Comfortable.PaddingX)
	}
	if Comfortable.PaddingX != 8 {
		t.Errorf("Comfortable.PaddingX = %v, want 8 — the HIG's inset beside a regular push button's title", Comfortable.PaddingX)
	}
	if Compact.PaddingX != 7 {
		t.Errorf("Compact.PaddingX = %v, want 7 — the published 8 scaled by 19/22", Compact.PaddingX)
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

// TestControlHeightsAreMeasured pins each density's control and field
// heights to the heights the provenance table in density.go states. A control
// draws max(floor, line), and no padding adds to the line, so these floors are
// the heights the platform's controls draw at.
func TestControlHeightsAreMeasured(t *testing.T) {
	for _, tc := range []struct {
		name string
		got  float32
		want float32
	}{
		{"Comfortable control", Comfortable.ControlHeight, 24},
		{"Compact control", Compact.ControlHeight, 19},
		{"Comfortable text field", Comfortable.TextFieldHeight, 27},
		{"Compact text field", Compact.TextFieldHeight, 21},
		{"Comfortable search field", Comfortable.SearchFieldHeight, 26},
		{"Compact search field", Compact.SearchFieldHeight, 21},
		{"Comfortable sidebar search field", Comfortable.SidebarSearchFieldHeight, 28},
		{"Compact sidebar search field", Compact.SidebarSearchFieldHeight, 22},
	} {
		if tc.got != tc.want {
			t.Errorf("%s height = %v dp, want %v", tc.name, tc.got, tc.want)
		}
	}
}

// TestSymbolSizeIsItsOwnValue pins the size of a symbol in a control: 20 dp
// at Comfortable and 19 dp at Compact, a value of its own that no padding
// derives.
func TestSymbolSizeIsItsOwnValue(t *testing.T) {
	if Comfortable.SymbolSize != 20 {
		t.Errorf("Comfortable.SymbolSize = %v dp, want 20", Comfortable.SymbolSize)
	}
	if Compact.SymbolSize != 19 {
		t.Errorf("Compact.SymbolSize = %v dp, want 19", Compact.SymbolSize)
	}
}
