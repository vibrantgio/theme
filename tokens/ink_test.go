package tokens_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/vibrantgio/theme/tokens"
)

// hexOf renders a colour the way the palette's own documents spell one, so
// a failing log line can be pasted into a contrast checker.
func hexOf(c color.NRGBA) string { return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B) }

// inkRoles are the seven roles that carry a pinned base, which is every
// role InkOn is defined for. RoleNeutral has no pin and is asserted to
// panic separately.
var inkRoles = []struct {
	name string
	role tokens.Role
	pin  func(tokens.ColorTokens) color.NRGBA
	ramp func(tokens.ColorTokens) tokens.Ramp
}{
	{"primary", tokens.RolePrimary, func(t tokens.ColorTokens) color.NRGBA { return t.Primary },
		func(t tokens.ColorTokens) tokens.Ramp { return t.Ramps.Primary }},
	{"secondary", tokens.RoleSecondary, func(t tokens.ColorTokens) color.NRGBA { return t.Secondary },
		func(t tokens.ColorTokens) tokens.Ramp { return t.Ramps.Secondary }},
	{"tertiary", tokens.RoleTertiary, func(t tokens.ColorTokens) color.NRGBA { return t.Tertiary },
		func(t tokens.ColorTokens) tokens.Ramp { return t.Ramps.Tertiary }},
	{"error", tokens.RoleError, func(t tokens.ColorTokens) color.NRGBA { return t.Error },
		func(t tokens.ColorTokens) tokens.Ramp { return t.Ramps.Error }},
	{"success", tokens.RoleSuccess, func(t tokens.ColorTokens) color.NRGBA { return t.Success },
		func(t tokens.ColorTokens) tokens.Ramp { return t.Ramps.Success }},
	{"warning", tokens.RoleWarning, func(t tokens.ColorTokens) color.NRGBA { return t.Warning },
		func(t tokens.ColorTokens) tokens.Ramp { return t.Ramps.Warning }},
	{"info", tokens.RoleInfo, func(t tokens.ColorTokens) color.NRGBA { return t.Info },
		func(t tokens.ColorTokens) tokens.Ramp { return t.Ramps.Info }},
}

// inkStoreys are the four storeys a brand ink can be drawn on: the paper a
// paragraph is set on and the three raised fills that host content above
// it. A link in a card and a link on the page are the same link and owe
// their own grounds the same ratio, so the gate is read against all four
// rather than against the page alone.
var inkStoreys = []struct {
	name  string
	level tokens.ElevationLevel
}{
	{"level 0 (paper)", tokens.Level0},
	{"level 1 (card)", tokens.Level1},
	{"level 2 (dialog)", tokens.Level2},
	{"level 3 (popover)", tokens.Level3},
}

// inkFloors are the two floors an ink is gated at: words and marks.
var inkFloors = []struct {
	name  string
	floor float64
}{
	{"text", tokens.TextFloor},
	{"graphic", tokens.GraphicFloor},
}

// inkSchemes yields every palette the sweep reads a seed as: both
// derivations, both schemes.
func inkSchemes(seed color.NRGBA) []struct {
	name  string
	tok   tokens.ColorTokens
	light bool
} {
	light, dark := tokens.FromSeed(seed)
	hcLight, hcDark := tokens.FromSeedHighContrast(seed)
	return []struct {
		name  string
		tok   tokens.ColorTokens
		light bool
	}{
		{"FromSeed light", light, true},
		{"FromSeed dark", dark, false},
		{"FromSeedHighContrast light", hcLight, true},
		{"FromSeedHighContrast dark", hcDark, false},
	}
}

// TestBrandInkClearsItsFloorForEverySeed is the gate the link ink defect
// asked for: whatever the seed, whatever the scheme, whatever the storey it
// is drawn on, a brand ink reaches the floor its job owes its ground.
//
// It is asserted for every pinned role and both floors rather than for the
// link alone, because the rule is one rule — a fill colour used as an ink —
// and the six roles that never fail are worth reading as measurements
// rather than assuming as facts.
func TestBrandInkClearsItsFloorForEverySeed(t *testing.T) {
	worst := map[string]float64{}
	worstAt := map[string]string{}
	for _, seed := range sweepSeeds() {
		for _, s := range inkSchemes(seed) {
			for _, st := range inkStoreys {
				ground := s.tok.SurfaceAt(st.level)
				for _, r := range inkRoles {
					for _, f := range inkFloors {
						ink := s.tok.InkOn(r.role, ground, f.floor)
						got := contrastRatio(ink, ground)
						if got < f.floor {
							t.Errorf("seed %v: %s %s: %s ink %v on %v measures %.2f:1, under the %.1f:1 %s floor",
								seed, s.name, st.name, r.name, ink, ground, got, f.floor, f.name)
						}
						// The answer is the brand's own colour or a rung of
						// the brand's own ramp, never something invented.
						if ink != r.pin(s.tok) {
							found := false
							for _, rung := range r.ramp(s.tok) {
								if rung == ink {
									found = true
									break
								}
							}
							if !found {
								t.Errorf("seed %v: %s %s: %s ink %v is neither the pin nor a rung of the role's ramp",
									seed, s.name, st.name, r.name, ink)
							}
						}
						key := f.name + " " + map[bool]string{true: "light", false: "dark"}[s.light]
						if w, ok := worst[key]; !ok || got < w {
							worst[key] = got
							worstAt[key] = r.name + " on " + st.name + " from seed " + hexOf(seed)
						}
					}
				}
			}
		}
	}
	for _, key := range []string{"text light", "text dark", "graphic light", "graphic dark"} {
		t.Logf("over %d seeds: worst %s brand ink %.2f:1 (%s)",
			len(sweepSeeds()), key, worst[key], worstAt[key])
	}
}

// TestBrandInkKeepsThePinThatReads is the other half of the rule, and it is
// what keeps this repair invisible to every stored image: a pin that clears
// its floor is returned untouched, so nothing moves for a palette that was
// already measuring.
func TestBrandInkKeepsThePinThatReads(t *testing.T) {
	for _, seed := range sweepSeeds() {
		for _, s := range inkSchemes(seed) {
			for _, st := range inkStoreys {
				ground := s.tok.SurfaceAt(st.level)
				for _, r := range inkRoles {
					for _, f := range inkFloors {
						pin := r.pin(s.tok)
						if contrastRatio(pin, ground) < f.floor {
							continue
						}
						if ink := s.tok.InkOn(r.role, ground, f.floor); ink != pin {
							t.Errorf("seed %v: %s %s: %s %s ink moved from the pin %v to %v though the pin measured %.2f:1",
								seed, s.name, st.name, r.name, f.name, pin, ink, contrastRatio(pin, ground))
						}
					}
				}
			}
		}
	}
}

// TestTheCanonicalSeedsBrandInkIsItsPin states the no-op for the one
// palette every golden image in this design system is rendered from. If
// this fails, a stored image somewhere else moved.
func TestTheCanonicalSeedsBrandInkIsItsPin(t *testing.T) {
	for _, s := range inkSchemes(tokens.DefaultSeed) {
		for _, st := range inkStoreys {
			ground := s.tok.SurfaceAt(st.level)
			for _, r := range inkRoles {
				for _, f := range inkFloors {
					if ink := s.tok.InkOn(r.role, ground, f.floor); ink != r.pin(s.tok) {
						t.Errorf("%s %s: %s %s ink is %v, not the pin %v — a stored image moved",
							s.name, st.name, r.name, f.name, ink, r.pin(s.tok))
					}
				}
			}
		}
	}
}

// TestOnlyTheLightPrimaryPinEverWalks records why the defect existed at all
// and bounds what the gate can touch. Six of the seven pinned bases are
// realized at fixed perceptual depths, so their contrast against the paper
// is a property of the derivation rather than of the brand and they never
// need the walk. The light primary base is the brand colour itself at the
// brand's own depth, which is the one place a seed can put an ink too near
// its ground — and the one place this gate ever answers with a rung.
func TestOnlyTheLightPrimaryPinEverWalks(t *testing.T) {
	walked := map[string]int{}
	worstPin := 99.0
	worstPinAt := ""
	for _, seed := range sweepSeeds() {
		for _, s := range inkSchemes(seed) {
			for _, st := range inkStoreys {
				ground := s.tok.SurfaceAt(st.level)
				for _, r := range inkRoles {
					for _, f := range inkFloors {
						if s.tok.InkOn(r.role, ground, f.floor) == r.pin(s.tok) {
							continue
						}
						scheme := "dark"
						if s.light {
							scheme = "light"
						}
						if r.role != tokens.RolePrimary || !s.light {
							t.Errorf("seed %v: %s %s: the %s %s pin walked, but only the light primary pin follows the seed's own depth",
								seed, s.name, st.name, r.name, f.name)
						}
						walked[f.name+" "+scheme]++
					}
				}
			}
		}
	}
	for _, seed := range sweepSeeds() {
		light, _ := tokens.FromSeed(seed)
		ground := light.SurfaceAt(tokens.Level0)
		if got := contrastRatio(light.Primary, ground); got < worstPin {
			worstPin, worstPinAt = got, hexOf(seed)
		}
	}
	t.Logf("over %d seeds, four palettes and four storeys: %d text walks and %d graphic walks, all of them the light primary pin; bare light pin over the paper bottoms out at %.2f:1 (%s)",
		len(sweepSeeds()), walked["text light"], walked["graphic light"], worstPin, worstPinAt)
}

// TestAPastelSeedGetsAReadableLinkInk is the regression this file was
// written for, read on the shape that produced it: an accent stated at a
// dark scheme's tone, used as a light scheme's seed. Its light primary pin
// lands a whisper off the paper, and before the gate that pin was the link
// colour a paragraph rendered with.
func TestAPastelSeedGetsAReadableLinkInk(t *testing.T) {
	seed := color.NRGBA{0x89, 0xb4, 0xfa, 0xff}
	light, dark := tokens.FromSeed(seed)

	lightPaper := light.SurfaceAt(tokens.Level0)
	if bare := contrastRatio(light.Primary, lightPaper); bare >= tokens.TextFloor {
		t.Fatalf("the pastel seed's bare light pin now measures %.2f:1 over the paper — this test no longer reads the shape it was written for", bare)
	}
	lightInk := light.InkOn(tokens.RolePrimary, lightPaper, tokens.TextFloor)
	if lightInk == light.Primary {
		t.Errorf("light link ink is still the bare pin %v", light.Primary)
	}
	if got := contrastRatio(lightInk, lightPaper); got < tokens.TextFloor {
		t.Errorf("light link ink %v on paper %v measures %.2f:1, under %.1f:1",
			lightInk, lightPaper, got, tokens.TextFloor)
	}

	// The dark scheme was never the broken half: its pin is realized at a
	// fixed depth, so it clears and is kept.
	darkPaper := dark.SurfaceAt(tokens.Level0)
	darkInk := dark.InkOn(tokens.RolePrimary, darkPaper, tokens.TextFloor)
	if darkInk != dark.Primary {
		t.Errorf("dark link ink walked to %v; the dark pin %v measures %.2f:1 and should stand",
			darkInk, dark.Primary, contrastRatio(dark.Primary, darkPaper))
	}
	if got := contrastRatio(darkInk, darkPaper); got < tokens.TextFloor {
		t.Errorf("dark link ink %v on paper %v measures %.2f:1, under %.1f:1",
			darkInk, darkPaper, got, tokens.TextFloor)
	}
	t.Logf("seed %s: light link %s on %s %.2f:1 (bare pin %s %.2f:1); dark link %s on %s %.2f:1",
		hexOf(seed), hexOf(lightInk), hexOf(lightPaper), contrastRatio(lightInk, lightPaper),
		hexOf(light.Primary), contrastRatio(light.Primary, lightPaper),
		hexOf(darkInk), hexOf(darkPaper), contrastRatio(darkInk, darkPaper))
}

// TestInkOnNeutralPanics: the neutral role carries surfaces and has no
// pinned fill, so asking it for a brand ink is a programming error — the
// same answer every other pin accessor gives.
func TestInkOnNeutralPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("InkOn(RoleNeutral, ...): expected panic")
		}
	}()
	tokens.DefaultLight.InkOn(tokens.RoleNeutral, tokens.DefaultLight.Background, tokens.TextFloor)
}
