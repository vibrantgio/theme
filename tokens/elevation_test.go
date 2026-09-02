package tokens_test

import (
	"image/color"
	"math"
	"testing"

	vgcolor "github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/tokens"
)

// levels is the elevation order, from the backdrop toward the reader. Every
// test below reads it in this direction and no other.
var levels = []struct {
	name  string
	level tokens.ElevationLevel
}{
	{"backdrop", tokens.LevelBackdrop},
	{"chrome", tokens.LevelChrome},
	{"content", tokens.Level0},
	{"raised", tokens.Level1},
	{"floating", tokens.Level2},
	{"top", tokens.Level3},
}

func lstar(c color.NRGBA) float64 {
	l, _, _ := vgcolor.LabFromNRGBA(c)
	return l
}

// TestElevationDpPreserved pins the dp shadow depths: the secondary cue,
// and what a depth effect and the token export read. Neither level under
// the content casts anything — the backdrop is the plane everything else
// stands on, and chrome lies flat on it.
func TestElevationDpPreserved(t *testing.T) {
	want := map[tokens.ElevationLevel]float32{
		tokens.LevelBackdrop: 0, tokens.LevelChrome: 0,
		tokens.Level0: 0, tokens.Level1: 1, tokens.Level2: 3, tokens.Level3: 6,
	}
	for level, dp := range want {
		if got := tokens.Elevation.Dp(level); got != dp {
			t.Errorf("Elevation.Dp(%d) = %v, want %v", level, got, dp)
		}
	}
	// The named fields carry the same dp values as the accessor.
	fields := map[tokens.ElevationLevel]float32{
		tokens.LevelBackdrop: tokens.Elevation.Backdrop,
		tokens.LevelChrome:   tokens.Elevation.Chrome,
		tokens.Level0:        tokens.Elevation.Level0,
		tokens.Level1:        tokens.Elevation.Level1,
		tokens.Level2:        tokens.Elevation.Level2,
		tokens.Level3:        tokens.Elevation.Level3,
	}
	for level, got := range fields {
		if want := tokens.Elevation.Dp(level); got != want {
			t.Errorf("Elevation field for level %d = %v, want %v", level, got, want)
		}
	}
}

// TestLightnessClimbsTowardTheViewer takes elevation's central invariant
// over the whole seed sweep in both schemes and both contrast variants:
// WALKING TOWARD THE VIEWER NEVER GETS DARKER. It is asserted over the
// population rather than over one palette, and with no mirror clause — the
// same sentence, the same direction, in every scheme.
func TestLightnessClimbsTowardTheViewer(t *testing.T) {
	var narrowest = 999.0
	var narrowestAt string
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		hcLight, hcDark := tokens.FromSeedHighContrast(seed)
		for _, scheme := range []struct {
			name string
			tok  tokens.ColorTokens
		}{{"light", light}, {"dark", dark}, {"light-hc", hcLight}, {"dark-hc", hcDark}} {
			for i := 1; i < len(levels); i++ {
				below := lstar(scheme.tok.SurfaceAt(levels[i-1].level))
				above := lstar(scheme.tok.SurfaceAt(levels[i].level))
				if above <= below {
					t.Fatalf("seed %v %s: %s L*%.2f is not lighter than %s L*%.2f — elevation inverts",
						seed, scheme.name, levels[i].name, above, levels[i-1].name, below)
				}
				if step := above - below; step < narrowest {
					narrowest, narrowestAt = step, scheme.name+" "+levels[i-1].name+"→"+levels[i].name
				}
			}
		}
	}
	t.Logf("over %d seeds × 4 schemes: narrowest level step %.2f L* (%s)",
		len(sweepSeeds()), narrowest, narrowestAt)
}

// TestTheLevelsKeepTheRestingPixels pins elevation's byte-for-byte
// identities: the light scheme's two levels under the content are Neutral
// 200 and 300 exactly, and the dark scheme's raised and floating levels
// land byte-for-byte on Neutral 200, 300 and 400 (the #222222/#2E2E2E
// family).
func TestTheLevelsKeepTheRestingPixels(t *testing.T) {
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		for level, step := range map[tokens.ElevationLevel]int{
			tokens.LevelChrome: 200, tokens.LevelBackdrop: 300,
		} {
			if got, want := light.SurfaceAt(level), light.Ramps.Neutral.Step(step); got != want {
				t.Fatalf("seed %v light: level %d = %v, want Neutral %d %v", seed, level, got, step, want)
			}
		}
		if got, want := light.SurfaceAt(tokens.Level0), light.Background; got != want {
			t.Fatalf("seed %v light: the content = %v, want the Background pin %v", seed, got, want)
		}
		for level, step := range map[tokens.ElevationLevel]int{
			tokens.Level1: 200, tokens.Level2: 300, tokens.Level3: 400,
		} {
			if got, want := dark.SurfaceAt(level), dark.Ramps.Neutral.Step(step); got != want {
				t.Fatalf("seed %v dark: level %d = %v, want Neutral %d %v", seed, level, got, step, want)
			}
		}
		if got, want := dark.SurfaceAt(tokens.Level0), dark.Background; got != want {
			t.Fatalf("seed %v dark: the content = %v, want the Background pin %v", seed, got, want)
		}
	}
	// The default palette's every level, written out, because these are
	// the values every screenshot in the plan is read against.
	want := map[string][6]color.NRGBA{
		"light": {
			{0xD4, 0xD4, 0xD4, 0xff}, {0xE8, 0xE8, 0xE8, 0xff}, {0xF6, 0xF6, 0xF6, 0xff},
			{0xF8, 0xF8, 0xF8, 0xff}, {0xFB, 0xFB, 0xFB, 0xff}, {0xFF, 0xFF, 0xFF, 0xff},
		},
		"dark": {
			{0x11, 0x11, 0x11, 0xff}, {0x15, 0x15, 0x15, 0xff}, {0x18, 0x18, 0x18, 0xff},
			{0x22, 0x22, 0x22, 0xff}, {0x2E, 0x2E, 0x2E, 0xff}, {0x47, 0x47, 0x47, 0xff},
		},
	}
	for _, scheme := range []struct {
		name string
		tok  tokens.ColorTokens
	}{{"light", tokens.DefaultLight}, {"dark", tokens.DefaultDark}} {
		for i, level := range levels {
			if got := scheme.tok.SurfaceAt(level.level); got != want[scheme.name][i] {
				t.Errorf("%s %s = %v, want %v", scheme.name, level.name, got, want[scheme.name][i])
			}
		}
	}
}

// TestTheChromeLevelTakesTheMeasuredStep pins the one level the ramp does
// not place in both schemes. Chrome's step under the content is a
// MEASUREMENT of the platform taken per scheme, and the two measurements
// differ by more than three times: about 4.9 L\* where the pin is the
// lightest surface the ramp carries — which is also the ramp's own first
// surface interval, so that half lands byte-for-byte on neutral 200 — and
// 1.48 L\* where the pin is the darkest, which is what Voice Memos (1.50),
// the reference chat application (1.71) and macOS Settings (3.81) measure
// and what a full band step of 4.98 overshot into pure black.
//
// The asymmetry is asserted rather than tolerated, because the temptation
// it guards against is deriving one number and mirroring it into the other
// scheme.
func TestTheChromeLevelTakesTheMeasuredStep(t *testing.T) {
	// One 8-bit level near either pin is well under a quarter of an L\*,
	// so the realized step cannot wander far from the measurement.
	const tolerance = 0.5
	deepest, deepestAt := 0.0, ""
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		hcLight, hcDark := tokens.FromSeedHighContrast(seed)
		for _, scheme := range []struct {
			name string
			tok  tokens.ColorTokens
			want float64 // the measured step for this scheme, in L*
		}{
			{"light", light, 0}, {"light-hc", hcLight, 0},
			{"dark", dark, 1.48}, {"dark-hc", hcDark, 1.48},
		} {
			pin := lstar(scheme.tok.Background)
			step := pin - lstar(scheme.tok.SurfaceAt(tokens.LevelChrome))
			want := scheme.want
			if want == 0 {
				// The light measurement IS the band's first surface
				// interval, so it is read off the band rather than named.
				band := scheme.tok.Ramps.Neutral
				want = lstar(band.Step(100)) - lstar(band.Step(200))
			}
			if step < want-tolerance || step > want+tolerance {
				t.Errorf("seed %v %s: chrome sits %.2f L* under the content, want the measured %.2f",
					seed, scheme.name, step, want)
			}
			if step > deepest {
				deepest, deepestAt = step, scheme.name
			}
		}
	}
	// The default dark chrome level, spelled out: the measured step realizes
	// #151515 under the #181818 content, and it is on no step of the ramp.
	chrome := tokens.DefaultDark.SurfaceAt(tokens.LevelChrome)
	if want := (color.NRGBA{0x15, 0x15, 0x15, 0xff}); chrome != want {
		t.Errorf("default dark chrome = %v, want the measured %v", chrome, want)
	}
	for step := 100; step <= 900; step += 100 {
		if rung := tokens.DefaultDark.Ramps.Neutral.Step(step); rung == chrome {
			t.Errorf("dark chrome %v landed on Neutral %d; the dark chrome step is measured, not a ramp step", chrome, step)
		}
	}
	t.Logf("over %d seeds × 4 schemes: deepest chrome step %.2f L* (%s)",
		len(sweepSeeds()), deepest, deepestAt)
}

// TestTheBackdropTakesTheDerivedStep pins the one level no stored platform
// capture can place: a macOS window paints its furniture edge to edge, so
// nothing measures a window plane beneath it. The backdrop's step is
// therefore the chrome step scaled by the surface band's own proportion —
// its second interval over its first — which is the same shape the levels
// above the pin take.
//
// Where the band descends from the pin that product is the band's second
// interval exactly, so the backdrop lands byte-for-byte on neutral 300; in
// the other direction it scales the measured whisper by the same
// proportion, and on the default dark palette realizes #111111 under the
// #151515 chrome. The derivation is asserted, not the number, because the
// number is not a measurement and must not be read as one.
func TestTheBackdropTakesTheDerivedStep(t *testing.T) {
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		hcLight, hcDark := tokens.FromSeedHighContrast(seed)
		for _, scheme := range []struct {
			name string
			tok  tokens.ColorTokens
		}{{"light", light}, {"light-hc", hcLight}, {"dark", dark}, {"dark-hc", hcDark}} {
			pin := lstar(scheme.tok.Background)
			band := scheme.tok.Ramps.Neutral
			first := math.Abs(lstar(band.Step(200)) - lstar(band.Step(100)))
			second := math.Abs(lstar(band.Step(300)) - lstar(band.Step(100)))
			chrome := pin - lstar(scheme.tok.SurfaceAt(tokens.LevelChrome))
			backdrop := pin - lstar(scheme.tok.SurfaceAt(tokens.LevelBackdrop))
			// One 8-bit level near either pin is well under a quarter of
			// an L*, so the realized step cannot wander far from the
			// derivation it renders.
			if want := chrome * second / first; math.Abs(backdrop-want) > 0.5 {
				t.Errorf("seed %v %s: the backdrop sits %.2f L* under the content, want the derived %.2f",
					seed, scheme.name, backdrop, want)
			}
		}
		if got, want := light.SurfaceAt(tokens.LevelBackdrop), light.Ramps.Neutral.Step(300); got != want {
			t.Errorf("seed %v light: the backdrop = %v, want Neutral 300 %v", seed, got, want)
		}
	}
	backdrop := tokens.DefaultDark.SurfaceAt(tokens.LevelBackdrop)
	if want := (color.NRGBA{0x11, 0x11, 0x11, 0xff}); backdrop != want {
		t.Errorf("default dark backdrop = %v, want the derived %v", backdrop, want)
	}
}

// TestTheHairlineCarriesTheWhisperStep is the light scheme's headroom
// strategy held to its own bargain. Above the paper the light scheme's
// steps are whispers — a fraction of an L* — so the derived hairline has
// to be what says where a raised surface is. This asserts that MarkOn's
// answer against every level clears WCAG 1.4.11's 3:1 in both schemes
// over the whole sweep, which is the condition under which the whisper is
// affordable, and it records how thin the fill step actually gets.
func TestTheHairlineCarriesTheWhisperStep(t *testing.T) {
	const graphicFloor = 3.0
	worst, worstAt := 99.0, ""
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		for _, scheme := range []struct {
			name string
			tok  tokens.ColorTokens
		}{{"light", light}, {"dark", dark}} {
			for _, level := range levels {
				fill := scheme.tok.SurfaceAt(level.level)
				got := contrastRatio(scheme.tok.MarkOn(tokens.RoleNeutral, fill, graphicFloor), fill)
				if got < graphicFloor {
					t.Errorf("seed %v %s %s: hairline reads %.2f:1 on the fill, floor %.1f:1",
						seed, scheme.name, level.name, got, graphicFloor)
				}
				if got < worst {
					worst, worstAt = got, scheme.name+" "+level.name
				}
			}
		}
	}
	// What the whisper measures on the default palette, so a reader can see
	// the size of the thing the hairline is covering for.
	for i := 3; i < len(levels); i++ {
		below := lstar(tokens.DefaultLight.SurfaceAt(levels[i-1].level))
		above := lstar(tokens.DefaultLight.SurfaceAt(levels[i].level))
		t.Logf("light %s→%s: %.2f L*", levels[i-1].name, levels[i].name, above-below)
	}
	t.Logf("over %d seeds: worst hairline on a level %.2f:1 (%s)", len(sweepSeeds()), worst, worstAt)
}

// TestRaisedWalksOneStepFromTheSurfaceItStandsOn asserts the step is taken
// from the level it is asked of and not from an absolute one: every level
// below the ceiling raises to the next one, so the same call answers
// differently for a thing on the window's content and the same thing on
// chrome. The walk is one-signed — the level it lands on is lighter in
// both schemes.
func TestRaisedWalksOneStepFromTheSurfaceItStandsOn(t *testing.T) {
	want := map[tokens.ElevationLevel]tokens.ElevationLevel{
		tokens.LevelBackdrop: tokens.LevelChrome,
		tokens.LevelChrome:   tokens.Level0,
		tokens.Level0:        tokens.Level1,
		tokens.Level1:        tokens.Level2,
		tokens.Level2:        tokens.Level3,
	}
	for ground, raised := range want {
		if got := ground.Raised(); got != raised {
			t.Errorf("ElevationLevel(%d).Raised() = %d, want %d", ground, got, raised)
		}
	}
	for _, scheme := range []struct {
		name string
		tok  tokens.ColorTokens
	}{{"light", tokens.DefaultLight}, {"dark", tokens.DefaultDark}} {
		for ground := range want {
			below := scheme.tok.SurfaceAt(ground)
			above := scheme.tok.SurfaceAt(ground.Raised())
			if above == below {
				t.Errorf("%s: level %d and the level raised from it both fill %v; a raised thing owes the surface it stands on one step",
					scheme.name, ground, above)
			}
			if lstar(above) <= lstar(below) {
				t.Errorf("%s: level %d raises from L*%.2f to L*%.2f — a step toward the viewer is lighter in both schemes",
					scheme.name, ground, lstar(below), lstar(above))
			}
		}
	}
}

// TestRaisedClampsAtTheCeiling pins the documented behaviour at the topmost
// level: the walk stops rather than naming a level the scale does not have.
// The scale has no level 4, so stepping past 3 here would hand every other
// accessor a level it panics on.
func TestRaisedClampsAtTheCeiling(t *testing.T) {
	if got := tokens.Level3.Raised(); got != tokens.Level3 {
		t.Errorf("Level3.Raised() = %d, want Level3 (%d): the levels end at 3", got, tokens.Level3)
	}
	if got, want := tokens.DefaultLight.SurfaceAt(tokens.Level3.Raised()), tokens.DefaultLight.SurfaceAt(tokens.Level3); got != want {
		t.Errorf("SurfaceAt(Level3.Raised()) = %v, want %v", got, want)
	}
}

// TestStateAtWalksFromTheLevel asserts the state walks compose on top of
// elevation from the level's OWN fill rather than from a ramp index the
// level does not have. The walk's direction is its own — toward the ramp's
// 900 end, darker in a light scheme and lighter in a dark one — because a
// state says something happened here, which is feedback and not depth.
func TestStateAtWalksFromTheLevel(t *testing.T) {
	for _, scheme := range []struct {
		name string
		tok  tokens.ColorTokens
		dark bool
	}{{"light", tokens.DefaultLight, false}, {"dark", tokens.DefaultDark, true}} {
		for _, level := range levels {
			fill := scheme.tok.SurfaceAt(level.level)
			if got := scheme.tok.StateAt(level.level, tokens.StateNormal); got != fill {
				t.Errorf("%s %s: normal = %v, want the level's own fill %v", scheme.name, level.name, got, fill)
			}
			hover := scheme.tok.StateAt(level.level, tokens.StateHover)
			if hover == fill {
				t.Errorf("%s %s: hover does not move off the fill %v", scheme.name, level.name, fill)
			}
			if lighter := lstar(hover) > lstar(fill); lighter != scheme.dark {
				t.Errorf("%s %s: hover walks the wrong way — fill L*%.2f, hover L*%.2f",
					scheme.name, level.name, lstar(fill), lstar(hover))
			}
			// The walk is PinnedStateColor's, taken from the level's own
			// fill, and it stops at the first depth clearing StateFloor:
			// never shallower than the one step that walk takes, and
			// never deeper than the floor requires.
			pinned := scheme.tok.PinnedStateColor(fill, tokens.StateHover)
			if lstar(hover) != lstar(pinned) {
				beyond := lstar(hover) > lstar(pinned)
				if beyond != scheme.dark {
					t.Errorf("%s %s: StateAt hover %v is shallower than the one-step walk %v",
						scheme.name, level.name, hover, pinned)
				}
			}
			if got := vgcolor.ContrastRatio(hover, fill); got < tokens.StateFloor {
				t.Errorf("%s %s: hover wash %v on the fill %v measures %.3f:1, under the %.2f:1 floor",
					scheme.name, level.name, hover, fill, got, tokens.StateFloor)
			}
			press := scheme.tok.StateAt(level.level, tokens.StatePressed)
			if lstar(press) == lstar(hover) || (lstar(press) > lstar(hover)) != scheme.dark {
				t.Errorf("%s %s: press %v does not lie beyond hover %v",
					scheme.name, level.name, press, hover)
			}
		}
	}
}

// TestSurfaceStepIsHalfTrue pins the deprecated accessor: it answers the
// ramp-index numbering, which is exactly the dark scheme's arrangement and
// stale in the light one, and it answers the two levels under the content
// with the "not a ramp step" sentinel it uses for the Background pin.
//
// The staleness is asserted rather than tolerated. A shim that is right in
// one scheme and wrong in the other is a thing to be moved off, and the
// assertion below says which half is which.
func TestSurfaceStepIsHalfTrue(t *testing.T) {
	want := map[tokens.ElevationLevel]int{
		tokens.LevelBackdrop: 0, tokens.LevelChrome: 0, tokens.Level0: 0,
		tokens.Level1: 200, tokens.Level2: 300, tokens.Level3: 400,
	}
	for level, step := range want {
		if got := tokens.Elevation.SurfaceStep(level); got != step {
			t.Errorf("Elevation.SurfaceStep(%d) = %d, want %d", level, got, step)
		}
		if step == 0 {
			continue
		}
		rung := tokens.DefaultDark.Ramps.Neutral.Step(step)
		if got := tokens.DefaultDark.SurfaceAt(level); got != rung {
			t.Errorf("dark: SurfaceStep(%d) names Neutral %d %v, but the level fills %v",
				level, step, rung, got)
		}
		if got := tokens.DefaultLight.SurfaceAt(level); got == tokens.DefaultLight.Ramps.Neutral.Step(step) {
			t.Errorf("light: level %d fills Neutral %d; the light scheme is meant to be OFF the ramp above the pin",
				level, step)
		}
	}
}

// TestElevationLevelPanics asserts out-of-vocabulary levels panic,
// matching Ramp.Step. 4 and 5 are in the list deliberately: the levels end
// at 3, so asking for a fifth is an error. −3 is the backdrop's counterpart
// at the other end: elevation has two levels below the content, not an
// open-ended basement.
// Raised answers to the same rule; its clamp at the ceiling is the one
// deliberate exception, and TestRaisedClampsAtTheCeiling pins it.
func TestElevationLevelPanics(t *testing.T) {
	for _, level := range []tokens.ElevationLevel{-5, -4, -3, 4, 5, 6} {
		for name, call := range map[string]func(){
			"SurfaceStep": func() { tokens.Elevation.SurfaceStep(level) },
			"Dp":          func() { tokens.Elevation.Dp(level) },
			"SurfaceAt":   func() { tokens.DefaultLight.SurfaceAt(level) },
			"StateAt":     func() { tokens.DefaultLight.StateAt(level, tokens.StateHover) },
			"Raised":      func() { level.Raised() },
		} {
			func() {
				defer func() {
					if recover() == nil {
						t.Errorf("%s(%d) did not panic", name, level)
					}
				}()
				call()
			}()
		}
	}
}
