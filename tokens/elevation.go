// Tonal elevation (goal G-E2, task E2.1): elevation is a surface step.
//
// Per ADR-007, a raised surface separates from its ground primarily by
// colour — each elevation level fills with a step of the neutral ramp, one
// step deeper per storey — and only secondarily by a cast shadow. Because
// the light and dark ramps are paired scales (same step, same job), the
// same level reads as "raised" in both modes: a level-1 card is a light
// card on a lighter ground in light mode and a dark card on a darker
// ground in dark mode, with no mode-specific rule. The dp shadow survives
// as the secondary cue and is what effects/depth still renders.
package tokens

import (
	"fmt"
	"image/color"
)

// ElevationScale pairs, per level, the neutral-ramp step of the level's
// surface fill (the primary cue) with its shadow depth in dp (the
// secondary cue).
//
// The ladder is four storeys, 0 through 3. MD3 levels 4 and 5 survived
// through v0.1.x as clamps onto level 3's step, so that call sites written
// against the MD3 numbering kept compiling; v0.2.0 deleted them, since a
// desktop surface has no six-storey stack to describe. A call site that
// named Level4 or Level5 meant "as raised as it gets" and should name
// Level3.
//
// The LevelN fields carry the dp depths — effects/depth's lookup and
// theme/export's --shadow-* table read them — and the StepN fields the
// paired surface steps. Prefer the Dp and SurfaceStep accessors over field
// access in new code.
type ElevationScale struct {
	// Shadow depths in device-independent pixels, following Material
	// Design 3 elevation levels 0–3. The secondary cue.
	Level0 float32 // 0 dp
	Level1 float32 // 1 dp
	Level2 float32 // 3 dp
	Level3 float32 // 6 dp

	// Surface-fill steps on the neutral ramp. Step0 is not a ramp step:
	// its zero value marks the Background pin — a level-0 surface is the
	// app's bg pin sitting over the step-100 ground.
	Step0 int // 0 — sentinel: the Background pin, not a ramp step
	Step1 int // 200
	Step2 int // 300
	Step3 int // 400
}

// Elevation is the default scale instance.
var Elevation = ElevationScale{
	Level0: 0,
	Level1: 1,
	Level2: 3,
	Level3: 6,

	Step0: 0, // Background pin
	Step1: 200,
	Step2: 300,
	Step3: 400,
}

// ElevationLevel selects an entry on the [ElevationScale] by name.
// The dp and step values for a given level are read from the [Elevation]
// instance.
type ElevationLevel int

const (
	Level0 ElevationLevel = iota
	Level1
	Level2
	Level3
)

// Raised is the rung one step above level: the fill of something that
// stands on a surface at level rather than on the ground of the window.
//
// Rungs are walked from the surface a thing is lying on, never from an
// absolute step (ADR-021 R4). So a thing raised over a level-0 plane fills
// at level 1 and the same thing over a level-1 plane fills at level 2, and
// a caller that names its own ground gets the separation the grammar asks
// for wherever it is placed. Reaching for a fixed step instead holds only
// while every ground in a window happens to be the one that step was
// chosen against; move the ground one rung and the raised thing sits two
// rungs off it, which the grammar calls a mistake in either direction.
//
// Level3 is the ceiling. The ladder is four storeys and ends there, so
// raising Level3 returns Level3 rather than naming a storey the scale does
// not have: a caller already at the top asked for "as raised as it gets"
// and is given it, and the clamp is the one place this walk stops instead
// of stepping. A level the ladder has no rung for — a negative one, or the
// MD3 levels 4 and 5 that were deleted when the ladder settled at four —
// panics, matching [ElevationScale.SurfaceStep] and [Ramp.Step].
func (level ElevationLevel) Raised() ElevationLevel {
	switch level {
	case Level0, Level1, Level2:
		return level + 1
	case Level3:
		return Level3
	}
	panic(fmt.Sprintf("tokens: unknown ElevationLevel %d", level))
}

// Dp returns level's shadow depth in device-independent pixels. An
// out-of-vocabulary level panics, matching [Ramp.Step].
func (e ElevationScale) Dp(level ElevationLevel) float32 {
	switch level {
	case Level0:
		return e.Level0
	case Level1:
		return e.Level1
	case Level2:
		return e.Level2
	case Level3:
		return e.Level3
	}
	panic(fmt.Sprintf("tokens: unknown ElevationLevel %d", level))
}

// SurfaceStep returns the neutral-ramp step of level's surface fill, or 0
// for a level whose fill is the Background pin rather than a ramp step
// (level 0 on the default scale). An out-of-vocabulary level panics,
// matching [Ramp.Step].
func (e ElevationScale) SurfaceStep(level ElevationLevel) int {
	switch level {
	case Level0:
		return e.Step0
	case Level1:
		return e.Step1
	case Level2:
		return e.Step2
	case Level3:
		return e.Step3
	}
	panic(fmt.Sprintf("tokens: unknown ElevationLevel %d", level))
}

// SurfaceAt resolves the surface colour of an elevated component: the fill
// of the given elevation level on t, per the default [Elevation] scale's
// step mapping. Level 0 is the Background pin over the step-100 ground;
// levels 1–3 fill with Neutral steps 200, 300 and 400.
//
// D2.3's state walks compose on top with the level's step as the ground:
// hover on a level-1 surface is StateColor(RoleNeutral, 200, StateHover),
// i.e. Neutral step 300 in both modes, courtesy of the paired scales. A
// level-0 surface is the app background, which has no ramp ground; treat
// interactive regions on it as level-1 surfaces instead.
func (t ColorTokens) SurfaceAt(level ElevationLevel) color.NRGBA {
	step := Elevation.SurfaceStep(level) // validates level
	if step == 0 {
		return t.Background
	}
	return t.Ramps.Neutral.Step(step)
}
