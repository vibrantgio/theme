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

// foregroundRoles are the seven roles that carry a pinned base, which is
// every role ForegroundOnAtFloor is defined for. RoleNeutral has no pin and
// is asserted to panic separately.
var foregroundRoles = []struct {
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

// foregroundLevels are the four levels a brand foreground can be drawn on:
// the content a paragraph is set on and the three raised fills that host
// content above it. A link in a card and a link on the page are the same
// link and owe their own surfaces the same ratio, so the gate is read
// against all four rather than against the page alone.
var foregroundLevels = []struct {
	name  string
	level tokens.ElevationLevel
}{
	{"level 0 (the content)", tokens.Level0},
	{"level 1 (card)", tokens.Level1},
	{"level 2 (dialog)", tokens.Level2},
	{"level 3 (popover)", tokens.Level3},
}

// foregroundFloors are the two floors a foreground is gated at: words and
// marks.
var foregroundFloors = []struct {
	name  string
	floor float64
}{
	{"text", tokens.TextFloor},
	{"graphic", tokens.GraphicFloor},
}

// foregroundSchemes yields every palette the sweep reads a seed as: both
// derivations, both schemes.
func foregroundSchemes(seed color.NRGBA) []struct {
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

// TestBrandForegroundClearsItsFloorForEverySeed is the gate the link
// foreground defect asked for: whatever the seed, whatever the scheme,
// whatever the level it is drawn on, a brand foreground reaches the floor
// its job owes its surface.
//
// It is asserted for every pinned role and both floors rather than for the
// link alone, because the rule is one rule — a fill colour used as a foreground —
// and the six roles that never fail are worth reading as measurements
// rather than assuming as facts.
func TestBrandForegroundClearsItsFloorForEverySeed(t *testing.T) {
	worst := map[string]float64{}
	worstAt := map[string]string{}
	for _, seed := range sweepSeeds() {
		for _, s := range foregroundSchemes(seed) {
			for _, st := range foregroundLevels {
				surface := s.tok.SurfaceAt(st.level)
				for _, r := range foregroundRoles {
					for _, f := range foregroundFloors {
						foreground := s.tok.ForegroundOnAtFloor(r.role, surface, f.floor)
						got := contrastRatio(foreground, surface)
						if got < f.floor {
							t.Errorf("seed %v: %s %s: %s foreground %v on %v measures %.2f:1, under the %.1f:1 %s floor",
								seed, s.name, st.name, r.name, foreground, surface, got, f.floor, f.name)
						}
						// The answer is the brand's own colour or a step of
						// the brand's own ramp, never something invented.
						if foreground != r.pin(s.tok) {
							found := false
							for _, step := range r.ramp(s.tok) {
								if step == foreground {
									found = true
									break
								}
							}
							if !found {
								t.Errorf("seed %v: %s %s: %s foreground %v is neither the pin nor a step of the role's ramp",
									seed, s.name, st.name, r.name, foreground)
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
		t.Logf("over %d seeds: worst %s brand foreground %.2f:1 (%s)",
			len(sweepSeeds()), key, worst[key], worstAt[key])
	}
}

// TestBrandForegroundKeepsThePinThatReads is the other half of the rule,
// and it is what keeps this repair invisible to every stored image: a pin
// that clears its floor is returned untouched, so nothing moves for a
// palette that was already measuring.
func TestBrandForegroundKeepsThePinThatReads(t *testing.T) {
	for _, seed := range sweepSeeds() {
		for _, s := range foregroundSchemes(seed) {
			for _, st := range foregroundLevels {
				surface := s.tok.SurfaceAt(st.level)
				for _, r := range foregroundRoles {
					for _, f := range foregroundFloors {
						pin := r.pin(s.tok)
						if contrastRatio(pin, surface) < f.floor {
							continue
						}
						if foreground := s.tok.ForegroundOnAtFloor(r.role, surface, f.floor); foreground != pin {
							t.Errorf("seed %v: %s %s: %s %s foreground moved from the pin %v to %v though the pin measured %.2f:1",
								seed, s.name, st.name, r.name, f.name, pin, foreground, contrastRatio(pin, surface))
						}
					}
				}
			}
		}
	}
}

// TestTheCanonicalSeedsBrandForegroundIsItsPin states the no-op for the one
// palette every golden image in this design system is rendered from. If
// this fails, a stored image somewhere else moved.
func TestTheCanonicalSeedsBrandForegroundIsItsPin(t *testing.T) {
	for _, s := range foregroundSchemes(tokens.DefaultSeed) {
		for _, st := range foregroundLevels {
			surface := s.tok.SurfaceAt(st.level)
			for _, r := range foregroundRoles {
				for _, f := range foregroundFloors {
					if foreground := s.tok.ForegroundOnAtFloor(r.role, surface, f.floor); foreground != r.pin(s.tok) {
						t.Errorf("%s %s: %s %s foreground is %v, not the pin %v — a stored image moved",
							s.name, st.name, r.name, f.name, foreground, r.pin(s.tok))
					}
				}
			}
		}
	}
}

// TestOnlyTheLightPrimaryPinEverWalks records why the defect existed at all
// and bounds what the gate can touch. Six of the seven pinned bases are
// realized at fixed perceptual depths, so their contrast against the content
// is a property of the derivation rather than of the brand and they never
// need the walk. The light primary base is the brand colour itself at the
// brand's own depth, which is the one place a seed can put a foreground too near
// its surface — and the one place this gate ever answers with a step.
func TestOnlyTheLightPrimaryPinEverWalks(t *testing.T) {
	walked := map[string]int{}
	worstPin := 99.0
	worstPinAt := ""
	for _, seed := range sweepSeeds() {
		for _, s := range foregroundSchemes(seed) {
			for _, st := range foregroundLevels {
				surface := s.tok.SurfaceAt(st.level)
				for _, r := range foregroundRoles {
					for _, f := range foregroundFloors {
						if s.tok.ForegroundOnAtFloor(r.role, surface, f.floor) == r.pin(s.tok) {
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
		surface := light.SurfaceAt(tokens.Level0)
		if got := contrastRatio(light.Primary, surface); got < worstPin {
			worstPin, worstPinAt = got, hexOf(seed)
		}
	}
	t.Logf("over %d seeds, four palettes and four levels: %d text walks and %d graphic walks, all of them the light primary pin; bare light pin over the content bottoms out at %.2f:1 (%s)",
		len(sweepSeeds()), walked["text light"], walked["graphic light"], worstPin, worstPinAt)
}

// TestAPastelSeedGetsAReadableLinkForeground is the regression this file was
// written for, read on the shape that produced it: an accent stated at a
// dark scheme's tone, used as a light scheme's seed. Its light primary pin
// lands a whisper off the content, and before the gate that pin was the link
// colour a paragraph rendered with.
func TestAPastelSeedGetsAReadableLinkForeground(t *testing.T) {
	seed := color.NRGBA{0x89, 0xb4, 0xfa, 0xff}
	light, dark := tokens.FromSeed(seed)

	lightContent := light.SurfaceAt(tokens.Level0)
	if bare := contrastRatio(light.Primary, lightContent); bare >= tokens.TextFloor {
		t.Fatalf("the pastel seed's bare light pin now measures %.2f:1 over the content — this test no longer reads the shape it was written for", bare)
	}
	lightForeground := light.ForegroundOnAtFloor(tokens.RolePrimary, lightContent, tokens.TextFloor)
	if lightForeground == light.Primary {
		t.Errorf("light link foreground is still the bare pin %v", light.Primary)
	}
	if got := contrastRatio(lightForeground, lightContent); got < tokens.TextFloor {
		t.Errorf("light link foreground %v on the content %v measures %.2f:1, under %.1f:1",
			lightForeground, lightContent, got, tokens.TextFloor)
	}

	// The dark scheme was never the broken half: its pin is realized at a
	// fixed depth, so it clears and is kept.
	darkContent := dark.SurfaceAt(tokens.Level0)
	darkForeground := dark.ForegroundOnAtFloor(tokens.RolePrimary, darkContent, tokens.TextFloor)
	if darkForeground != dark.Primary {
		t.Errorf("dark link foreground walked to %v; the dark pin %v measures %.2f:1 and should stand",
			darkForeground, dark.Primary, contrastRatio(dark.Primary, darkContent))
	}
	if got := contrastRatio(darkForeground, darkContent); got < tokens.TextFloor {
		t.Errorf("dark link foreground %v on the content %v measures %.2f:1, under %.1f:1",
			darkForeground, darkContent, got, tokens.TextFloor)
	}
	t.Logf("seed %s: light link %s on %s %.2f:1 (bare pin %s %.2f:1); dark link %s on %s %.2f:1",
		hexOf(seed), hexOf(lightForeground), hexOf(lightContent), contrastRatio(lightForeground, lightContent),
		hexOf(light.Primary), contrastRatio(light.Primary, lightContent),
		hexOf(darkForeground), hexOf(darkContent), contrastRatio(darkForeground, darkContent))
}

// TestForegroundOnAtFloorNeutralPanics: the neutral role carries surfaces
// and has no pinned fill, so asking it for a brand foreground is a
// programming error — the same answer every other pin accessor gives.
func TestForegroundOnAtFloorNeutralPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("ForegroundOnAtFloor(RoleNeutral, ...): expected panic")
		}
	}()
	tokens.DefaultLight.ForegroundOnAtFloor(tokens.RoleNeutral, tokens.DefaultLight.Background, tokens.TextFloor)
}

// TestForegroundOnAnswersEveryRole: the shared foreground derivation is
// [ColorTokens.ForegroundOnAtFloor] at the text floor for the roles that
// carry a pinned base, and the walk for RoleNeutral, which carries none —
// so unlike ForegroundOnAtFloor it answers every role rather than refusing
// one. A component drawing a role's word over a fill of that role's hue
// calls this, and only this; a second spelling is how two components stop
// matching.
func TestForegroundOnAnswersEveryRole(t *testing.T) {
	for _, sc := range []struct {
		name string
		c    tokens.ColorTokens
	}{{"light", tokens.DefaultLight}, {"dark", tokens.DefaultDark}} {
		c := sc.c
		for _, role := range []tokens.Role{
			tokens.RoleNeutral, tokens.RolePrimary, tokens.RoleSecondary,
			tokens.RoleTertiary, tokens.RoleError, tokens.RoleSuccess,
			tokens.RoleWarning, tokens.RoleInfo,
		} {
			fill := c.ContainerOn(role, c.SurfaceAt(tokens.Level0))
			got := c.ForegroundOn(role, fill)

			want := c.MarkOn(role, fill, tokens.TextFloor)
			if role != tokens.RoleNeutral {
				want = c.ForegroundOnAtFloor(role, fill, tokens.TextFloor)
			}
			if got != want {
				t.Errorf("%s role %d: ForegroundOn = %s, want %s",
					sc.name, role, hexOf(got), hexOf(want))
			}
			if ratio := contrastRatio(got, fill); ratio < tokens.TextFloor {
				t.Errorf("%s role %d: foreground %s on its own container %s = %.2f:1, under %.1f:1",
					sc.name, role, hexOf(got), hexOf(fill), ratio, tokens.TextFloor)
			}
		}
	}
}
