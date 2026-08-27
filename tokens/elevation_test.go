package tokens

import "testing"

// TestElevationDpPreserved pins the dp shadow depths: the secondary cue
// E2.1 keeps byte-for-byte, and what effects/depth and theme/export read.
func TestElevationDpPreserved(t *testing.T) {
	want := map[ElevationLevel]float32{
		Level0: 0, Level1: 1, Level2: 3, Level3: 6,
	}
	for level, dp := range want {
		if got := Elevation.Dp(level); got != dp {
			t.Errorf("Elevation.Dp(%d) = %v, want %v", level, got, dp)
		}
	}
	// The named fields carry the same dp values as the accessor.
	fields := []float32{
		Elevation.Level0, Elevation.Level1, Elevation.Level2,
		Elevation.Level3,
	}
	for i, got := range fields {
		if want := Elevation.Dp(ElevationLevel(i)); got != want {
			t.Errorf("Elevation.Level%d field = %v, want %v", i, got, want)
		}
	}
}

// TestElevationSurfaceSteps asserts each level's fill sits on the neutral
// ramp — levels 1–3 on steps 200/300/400 — with level 0 the Background pin
// (sentinel step 0).
func TestElevationSurfaceSteps(t *testing.T) {
	want := map[ElevationLevel]int{
		Level0: 0, // Background pin, not a ramp step
		Level1: 200,
		Level2: 300,
		Level3: 400,
	}
	for level, step := range want {
		if got := Elevation.SurfaceStep(level); got != step {
			t.Errorf("Elevation.SurfaceStep(%d) = %d, want %d", level, got, step)
		}
	}
}

// TestSurfaceAt asserts the resolver in both modes: level 0 is the
// Background pin, levels 1–3 the neutral ramp steps 200/300/400.
func TestSurfaceAt(t *testing.T) {
	for _, mode := range []struct {
		name string
		tok  ColorTokens
	}{{"light", DefaultLight}, {"dark", DefaultDark}} {
		if got := mode.tok.SurfaceAt(Level0); got != mode.tok.Background {
			t.Errorf("%s: SurfaceAt(Level0) = %v, want Background pin %v", mode.name, got, mode.tok.Background)
		}
		steps := map[ElevationLevel]int{Level1: 200, Level2: 300, Level3: 400}
		for level, step := range steps {
			if got, want := mode.tok.SurfaceAt(level), mode.tok.Ramps.Neutral.Step(step); got != want {
				t.Errorf("%s: SurfaceAt(%d) = %v, want Neutral.Step(%d) %v", mode.name, level, got, step, want)
			}
		}
	}
}

// TestRaisedWalksOneRungFromTheLocalGround asserts the step is taken from
// the level it is asked of and not from an absolute rung: every level below
// the ceiling raises to the next one, so the same call answers differently
// for a thing on the window's paper and the same thing on furniture.
func TestRaisedWalksOneRungFromTheLocalGround(t *testing.T) {
	want := map[ElevationLevel]ElevationLevel{
		Level0: Level1,
		Level1: Level2,
		Level2: Level3,
	}
	for ground, raised := range want {
		if got := ground.Raised(); got != raised {
			t.Errorf("Level%d.Raised() = %d, want %d", ground, got, raised)
		}
	}
	// The fill a walked rung resolves to is the next step of the neutral
	// ramp, so what a raised thing paints is a rung off its own ground
	// rather than off the window's.
	for _, mode := range []struct {
		name string
		tok  ColorTokens
	}{{"light", DefaultLight}, {"dark", DefaultDark}} {
		for ground := range want {
			got := mode.tok.SurfaceAt(ground.Raised())
			if bg := mode.tok.SurfaceAt(ground); got == bg {
				t.Errorf("%s: a level-%d ground and its raised rung both fill %v; a raised thing owes its ground one rung", mode.name, ground, got)
			}
			if step := mode.tok.Ramps.Neutral.Step(Elevation.SurfaceStep(ground.Raised())); got != step {
				t.Errorf("%s: SurfaceAt(Level%d.Raised()) = %v, want %v", mode.name, ground, got, step)
			}
		}
	}
}

// TestRaisedClampsAtTheCeiling pins the documented behaviour at the top of
// the ladder: four storeys is all there is, so the walk stops rather than
// naming a fifth. Levels 4 and 5 were deleted when the ladder settled, and
// stepping past 3 here would hand every other accessor a level it panics
// on.
func TestRaisedClampsAtTheCeiling(t *testing.T) {
	if got := Level3.Raised(); got != Level3 {
		t.Errorf("Level3.Raised() = %d, want Level3 (%d): the ladder ends at 3", got, Level3)
	}
	if got, want := DefaultLight.SurfaceAt(Level3.Raised()), DefaultLight.SurfaceAt(Level3); got != want {
		t.Errorf("SurfaceAt(Level3.Raised()) = %v, want %v", got, want)
	}
}

// TestSurfaceStateComposition asserts D2.3's state walks compose on top of
// the elevation mapping: a state on an elevated surface uses the level's
// step as its ground, so hover on a level-1 surface is step 300 and
// pressed is step 400 — in both modes, courtesy of the paired scales.
func TestSurfaceStateComposition(t *testing.T) {
	for _, mode := range []struct {
		name string
		tok  ColorTokens
	}{{"light", DefaultLight}, {"dark", DefaultDark}} {
		ground := Elevation.SurfaceStep(Level1) // 200
		if got, want := mode.tok.StateColor(RoleNeutral, ground, StateHover), mode.tok.Ramps.Neutral.Step(300); got != want {
			t.Errorf("%s: hover on level-1 surface = %v, want Neutral.Step(300) %v", mode.name, got, want)
		}
		if got, want := mode.tok.StateColor(RoleNeutral, ground, StatePressed), mode.tok.Ramps.Neutral.Step(400); got != want {
			t.Errorf("%s: pressed on level-1 surface = %v, want Neutral.Step(400) %v", mode.name, got, want)
		}
	}
}

// TestElevationLevelPanics asserts out-of-vocabulary levels panic,
// matching Ramp.Step. 4 and 5 are in the list deliberately: through v0.1.x
// they were clamps onto level 3, and F3.3's sweep deleted them, so the
// ladder ends at 3 and asking for a fifth storey is now the error it always
// described. Raised answers to the same rule; its clamp at the ceiling is
// the one deliberate exception, and TestRaisedClampsAtTheCeiling pins it.
func TestElevationLevelPanics(t *testing.T) {
	for _, level := range []ElevationLevel{-1, 4, 5, 6} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("Elevation.SurfaceStep(%d) did not panic", level)
				}
			}()
			Elevation.SurfaceStep(level)
		}()
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("Elevation.Dp(%d) did not panic", level)
				}
			}()
			Elevation.Dp(level)
		}()
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("SurfaceAt(%d) did not panic", level)
				}
			}()
			DefaultLight.SurfaceAt(level)
		}()
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("ElevationLevel(%d).Raised() did not panic", level)
				}
			}()
			level.Raised()
		}()
	}
}
