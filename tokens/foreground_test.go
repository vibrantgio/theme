package tokens_test

import (
	"fmt"
	"image/color"
	"testing"

	vgcolor "github.com/vibrantgio/theme/color"
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
// link and owe their own surfaces the same floor, so the gate is read
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
						got := vgcolor.Magnitude(foreground, surface)
						if got < f.floor {
							t.Errorf("seed %v: %s %s: %s foreground %v on %v measures |Lc| %.2f, under the %.0f %s floor",
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
		t.Logf("over %d seeds: worst %s brand foreground |Lc| %.2f (%s)",
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
						if vgcolor.Magnitude(pin, surface) < f.floor {
							continue
						}
						if foreground := s.tok.ForegroundOnAtFloor(r.role, surface, f.floor); foreground != pin {
							t.Errorf("seed %v: %s %s: %s %s foreground moved from the pin %v to %v though the pin measured |Lc| %.2f",
								seed, s.name, st.name, r.name, f.name, pin, foreground, vgcolor.Magnitude(pin, surface))
						}
					}
				}
			}
		}
	}
}

// TestTheCanonicalSeedsBrandForegroundWalksForWordsAndStandsForMarks states
// what each floor costs the one palette every golden image in this design
// system is rendered from.
//
// At [tokens.GraphicFloor] every pinned base stands on every level, in both
// schemes and both derivations: a mark drawn in a role's colour is the
// role's own colour. At [tokens.TextFloor] the pin stands only where Lc 75
// is inside its reach — on the light scheme's white raised fills it is, on
// the content beneath them it is not, and in the dark scheme it never is —
// so a word in a brand colour is a step of that brand's ramp wherever the
// pin falls short, and the answer reaches the floor either way.
func TestTheCanonicalSeedsBrandForegroundWalksForWordsAndStandsForMarks(t *testing.T) {
	for _, s := range foregroundSchemes(tokens.DefaultSeed) {
		for _, st := range foregroundLevels {
			surface := s.tok.SurfaceAt(st.level)
			for _, r := range foregroundRoles {
				pin := r.pin(s.tok)
				if foreground := s.tok.ForegroundOnAtFloor(r.role, surface, tokens.GraphicFloor); foreground != pin {
					t.Errorf("%s %s: %s mark foreground is %v, not the pin %v",
						s.name, st.name, r.name, foreground, pin)
				}
				foreground := s.tok.ForegroundOnAtFloor(r.role, surface, tokens.TextFloor)
				if bare := vgcolor.Magnitude(pin, surface); (foreground == pin) != (bare >= tokens.TextFloor) {
					t.Errorf("%s %s: %s text foreground is %v while the pin %v measures |Lc| %.2f on %v",
						s.name, st.name, r.name, foreground, pin, bare, surface)
				}
				if got := vgcolor.Magnitude(foreground, surface); got < tokens.TextFloor {
					t.Errorf("%s %s: %s text foreground %v on %v measures |Lc| %.2f, under the %.0f text floor",
						s.name, st.name, r.name, foreground, surface, got, tokens.TextFloor)
				}
			}
		}
	}
}

// TestOnlyTheLightPrimaryPinEverWalksForAMark bounds what the mark floor can
// touch. Six of the seven pinned bases are realized at fixed perceptual
// depths, so their reading against the content is a property of the
// derivation rather than of the brand and they never need the walk at
// [tokens.GraphicFloor]. The light primary base is the brand colour itself at
// the brand's own depth, which is the one place a seed can put a mark too
// near its surface.
//
// [tokens.TextFloor] is a different matter and the counts record it: Lc 75
// over a level is outside most pins' reach, so a pin walks for words far
// more often than it walks for a mark.
func TestOnlyTheLightPrimaryPinEverWalksForAMark(t *testing.T) {
	walked := map[string]int{}
	stood := map[string]int{}
	worstPin := 999.0
	worstPinAt := ""
	for _, seed := range sweepSeeds() {
		for _, s := range foregroundSchemes(seed) {
			for _, st := range foregroundLevels {
				surface := s.tok.SurfaceAt(st.level)
				for _, r := range foregroundRoles {
					for _, f := range foregroundFloors {
						scheme := "dark"
						if s.light {
							scheme = "light"
						}
						if s.tok.ForegroundOnAtFloor(r.role, surface, f.floor) == r.pin(s.tok) {
							stood[f.name+" "+scheme]++
							continue
						}
						if f.floor == tokens.GraphicFloor && (r.role != tokens.RolePrimary || !s.light) {
							t.Errorf("seed %v: %s %s: the %s mark pin walked, but only the light primary pin follows the seed's own depth",
								seed, s.name, st.name, r.name)
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
		if got := vgcolor.Magnitude(light.Primary, surface); got < worstPin {
			worstPin, worstPinAt = got, hexOf(seed)
		}
	}
	t.Logf("over %d seeds, four palettes and four levels: %d light and %d dark mark walks, %d light and %d dark text walks (%d light and %d dark text pins stood); bare light pin over the content bottoms out at |Lc| %.2f (%s)",
		len(sweepSeeds()), walked["graphic light"], walked["graphic dark"],
		walked["text light"], walked["text dark"], stood["text light"], stood["text dark"], worstPin, worstPinAt)
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
	if bare := vgcolor.Magnitude(light.Primary, lightContent); bare >= tokens.TextFloor {
		t.Fatalf("the pastel seed's bare light pin now measures |Lc| %.2f over the content — this test no longer reads the shape it was written for", bare)
	}
	lightForeground := light.ForegroundOnAtFloor(tokens.RolePrimary, lightContent, tokens.TextFloor)
	if lightForeground == light.Primary {
		t.Errorf("light link foreground is still the bare pin %v", light.Primary)
	}
	if got := vgcolor.Magnitude(lightForeground, lightContent); got < tokens.TextFloor {
		t.Errorf("light link foreground %v on the content %v measures |Lc| %.2f, under %.0f",
			lightForeground, lightContent, got, tokens.TextFloor)
	}

	// The dark scheme was never the broken half — its pin is realized at a
	// fixed depth — but Lc 75 is further than that depth reaches, so it walks
	// one step and it is the walk that has to read.
	darkContent := dark.SurfaceAt(tokens.Level0)
	darkForeground := dark.ForegroundOnAtFloor(tokens.RolePrimary, darkContent, tokens.TextFloor)
	if got := vgcolor.Magnitude(darkForeground, darkContent); got < tokens.TextFloor {
		t.Errorf("dark link foreground %v on the content %v measures |Lc| %.2f, under %.0f",
			darkForeground, darkContent, got, tokens.TextFloor)
	}
	t.Logf("seed %s: light link %s on %s |Lc| %.2f (bare pin %s |Lc| %.2f); dark link %s on %s |Lc| %.2f",
		hexOf(seed), hexOf(lightForeground), hexOf(lightContent), vgcolor.Magnitude(lightForeground, lightContent),
		hexOf(light.Primary), vgcolor.Magnitude(light.Primary, lightContent),
		hexOf(darkForeground), hexOf(darkContent), vgcolor.Magnitude(darkForeground, darkContent))
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
			if ratio := vgcolor.Magnitude(got, fill); ratio < tokens.TextFloor {
				t.Errorf("%s role %d: foreground %s on its own container %s = |Lc| %.2f, under %.0f",
					sc.name, role, hexOf(got), hexOf(fill), ratio, tokens.TextFloor)
			}
		}
	}
}
