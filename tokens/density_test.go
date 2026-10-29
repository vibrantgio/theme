package tokens

import "testing"

// TestDensityPicksMatchThePlatformScale pins the two control heights to the
// platform's scale as the reference reads it, with the provenance recorded
// in density.go and in organisation/reference/macos/controls.md: Comfortable
// is the regular control measured off the stored Save panel at 24 dp, and
// Compact is the small control measured off Finder's View Options window and
// search criteria rows at 20 dp. Each is also the pointer target a control
// of that density offers.
func TestDensityPicksMatchThePlatformScale(t *testing.T) {
	const (
		measuredRegular float32 = 24 // MEASURED: save-dialog-{light,dark}.png, push button and pop-up button
		measuredSmall   float32 = 20 // MEASURED: finder-view-options-{light,dark}.png and finder-search-criteria-{light,dark}.png, pop-ups and push buttons
	)
	if ComfortableControlHeight != measuredRegular {
		t.Errorf("ComfortableControlHeight = %v, want %v (the platform's regular control, measured)",
			ComfortableControlHeight, measuredRegular)
	}
	if CompactControlHeight != measuredSmall {
		t.Errorf("CompactControlHeight = %v, want %v (the platform's small control, measured)",
			CompactControlHeight, measuredSmall)
	}
	if CompactControlHeight >= ComfortableControlHeight {
		t.Errorf("CompactControlHeight (%v) must be < ComfortableControlHeight (%v)",
			CompactControlHeight, ComfortableControlHeight)
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

// TestCompactValuesAreTheirMeasuredRows pins each Compact value to its
// MEASURED row in organisation/reference/macos/controls.md, under "What the
// density scale takes". Each one is its own reading of a small control, and
// none is computed from a Comfortable value.
func TestCompactValuesAreTheirMeasuredRows(t *testing.T) {
	for _, tc := range []struct {
		name string
		got  float32
		row  float32
		from string
	}{
		{"control height", Compact.ControlHeight, 20, "finder-view-options-{light,dark}.png and finder-search-criteria-{light,dark}.png, the small pop-ups and push buttons"},
		{"text field", Compact.TextFieldHeight, 21, "finder-search-criteria-{light,dark}.png, the edge at y 38 and y 58"},
		{"search field", Compact.SearchFieldHeight, 22, "automator-search-field-{light,dark}.png, the edge at y 69 and y 90"},
		{"sidebar search field", Compact.SidebarSearchFieldHeight, 28, "systemsettings-sidebar-search-small-{light,dark}.png, y 6–33"},
		{"row height", Compact.RowHeight, 20, "finder-list-small-{light,dark}.png, a 20 px stripe pitch"},
		{"checkbox row", Compact.CheckboxRowHeight, 20, "finder-view-options-{light,dark}.png, boxes at y 338–351 and y 358–371"},
		{"PaddingX", Compact.PaddingX, 11, "finder-view-options-{light,dark}.png and finder-search-criteria-{light,dark}.png, the small pop-up's fill start to the first stem"},
	} {
		if tc.got != tc.row {
			t.Errorf("Compact %s = %v dp, want its measured row %v (%s)", tc.name, tc.got, tc.row, tc.from)
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
// in density.go: the heights are exactly the consts, no Compact height
// exceeds its Comfortable height, and PaddingX and ChipPaddingX are the pairs the
// table states.
// The sidebar search field is the same at both, as the platform draws it.
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
	if Compact.ControlHeight >= Comfortable.ControlHeight || Compact.TextFieldHeight >= Comfortable.TextFieldHeight || Compact.SearchFieldHeight >= Comfortable.SearchFieldHeight || Compact.SidebarSearchFieldHeight > Comfortable.SidebarSearchFieldHeight {
		t.Errorf("Compact heights must not exceed Comfortable heights: %+v against %+v", Compact, Comfortable)
	}
	if Comfortable.PaddingX != 13 {
		t.Errorf("Comfortable.PaddingX = %v, want 13 — the regular push button's and pop-up's measured title inset in save-panel-untinted-{light,dark}.png and save-panel-footer-untinted-{light,dark}.png", Comfortable.PaddingX)
	}
	if Compact.PaddingX != 11 {
		t.Errorf("Compact.PaddingX = %v, want 11 — the small pop-up's measured title inset", Compact.PaddingX)
	}
	if Comfortable.ChipPaddingX != 8 || Compact.ChipPaddingX != 11 {
		t.Errorf("ChipPaddingX = (%v, %v), want the carried (8, 11)", Comfortable.ChipPaddingX, Compact.ChipPaddingX)
	}
}

// TestChipHeightRelation pins the chip height to the relation and not to two
// picked numbers: it is exactly ChipDrop under whatever control height the
// setting carries, which lands Comfortable on 20 and Compact on 16. Changing
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
	if Comfortable.ChipHeight() != 20 || Compact.ChipHeight() != 16 {
		t.Errorf("chip heights = (%v, %v), want (20, 16)",
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
		{"Compact control", Compact.ControlHeight, 20},
		{"Comfortable text field", Comfortable.TextFieldHeight, 27},
		{"Compact text field", Compact.TextFieldHeight, 21},
		{"Comfortable search field", Comfortable.SearchFieldHeight, 26},
		{"Compact search field", Compact.SearchFieldHeight, 22},
		{"Comfortable sidebar search field", Comfortable.SidebarSearchFieldHeight, 28},
		{"Compact sidebar search field", Compact.SidebarSearchFieldHeight, 28},
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
