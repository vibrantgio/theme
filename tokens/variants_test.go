package tokens_test

import (
	"fmt"
	stdcolor "image/color"
	"testing"

	"github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/tokens"
)

// variantSurfaces are the two neutral surfaces both variant steps are chosen
// against: the surface a chrome region wears and the content surface beside
// it. A token that cleared only the easier of the two would be a promise
// kept on one of them.
func variantSurfaces(t tokens.ColorTokens) []struct {
	name    string
	surface stdcolor.NRGBA
} {
	return []struct {
		name    string
		surface stdcolor.NRGBA
	}{
		{"Surface", t.Surface},
		{"Background", t.Background},
	}
}

// TestVariantStepsClearTheirFloors is the gate the two variant steps exist
// for: whatever the seed, the scheme or the derivation, the boundary step
// reaches the 3:1 a mark owes its surface and the muted foreground reaches the 4.5:1
// a run of words owes it — against BOTH neutral surfaces, not the easier one.
//
// It is the failure M3's fixed outline role has and this palette must not:
// a step named once states one colour and as many measurements as it has
// surfaces.
func TestVariantStepsClearTheirFloors(t *testing.T) {
	worstOutline, worstOutlineAt := 99.0, ""
	worstForeground, worstForegroundAt := 99.0, ""
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		hcLight, hcDark := tokens.FromSeedHighContrast(seed)
		for _, s := range []struct {
			name string
			tok  tokens.ColorTokens
		}{
			{"FromSeed light", light}, {"FromSeed dark", dark},
			{"FromSeedHighContrast light", hcLight}, {"FromSeedHighContrast dark", hcDark},
		} {
			for _, g := range variantSurfaces(s.tok) {
				if got := color.Magnitude(s.tok.OutlineVariant(), g.surface); got < tokens.GraphicFloor {
					t.Errorf("seed %v: %s OutlineVariant on %s measures |Lc| %.2f, under the %.0f boundary floor",
						seed, s.name, g.name, got, tokens.GraphicFloor)
				} else if got < worstOutline {
					worstOutline, worstOutlineAt = got, fmt.Sprintf("%s on %s", s.name, g.name)
				}
				if got := color.Magnitude(s.tok.OnSurfaceVariant(), g.surface); got < tokens.TextFloor {
					t.Errorf("seed %v: %s OnSurfaceVariant on %s measures |Lc| %.2f, under the %.0f text floor",
						seed, s.name, g.name, got, tokens.TextFloor)
				} else if got < worstForeground {
					worstForeground, worstForegroundAt = got, fmt.Sprintf("%s on %s", s.name, g.name)
				}
			}
		}
	}
	t.Logf("over %d seeds: worst OutlineVariant |Lc| %.2f (%s), worst OnSurfaceVariant |Lc| %.2f (%s)",
		len(sweepSeeds()), worstOutline, worstOutlineAt, worstForeground, worstForegroundAt)
}

// TestVariantStepsStayMuted holds the other half of each step's contract.
// Clearing a floor is easy at the ramp's far end — step 900 clears
// everything — and a token that went there would be the most pronounced step
// under a muted name. Both are the LEAST PRONOUNCED step that clears, so each
// must stay under the token that speaks at full strength: the boundary under
// the muted foreground, and the muted foreground under Text.
func TestVariantStepsStayMuted(t *testing.T) {
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		hcLight, hcDark := tokens.FromSeedHighContrast(seed)
		for _, s := range []struct {
			name string
			tok  tokens.ColorTokens
		}{
			{"FromSeed light", light}, {"FromSeed dark", dark},
			{"FromSeedHighContrast light", hcLight}, {"FromSeedHighContrast dark", hcDark},
		} {
			for _, g := range variantSurfaces(s.tok) {
				outline := color.Magnitude(s.tok.OutlineVariant(), g.surface)
				foreground := color.Magnitude(s.tok.OnSurfaceVariant(), g.surface)
				text := color.Magnitude(s.tok.Text, g.surface)
				if outline > foreground {
					t.Errorf("seed %v: %s OutlineVariant on %s measures |Lc| %.2f, more pronounced than OnSurfaceVariant's %.2f",
						seed, s.name, g.name, outline, foreground)
				}
				if foreground >= text {
					t.Errorf("seed %v: %s OnSurfaceVariant on %s measures |Lc| %.2f, at or past Text's %.2f — a muted foreground that is not muted",
						seed, s.name, g.name, foreground, text)
				}
			}
		}
	}
}

// TestVariantStepsBeatTheFixedStep is the comparison the derivation is worth
// having: naming one neutral step and calling it the outline — M3's fixed
// role, and what this sheet's borders once did — reads under the boundary
// floor in the light scheme, and the derived token does not.
//
// Step 500 falls short of the boundary floor over the harder of the two
// light surfaces and clears it over the harder dark one, so the fixed
// naming fails in the scheme most people read in and passes in the other. The
// test asserts the shape rather than the digits — that the fixed step falls
// short somewhere the derived step does not — so a re-derived ramp moves the
// numbers without turning this red for the wrong reason.
func TestVariantStepsBeatTheFixedStep(t *testing.T) {
	fixedFails := false
	for _, s := range []struct {
		name string
		tok  tokens.ColorTokens
	}{
		{"DefaultLight", tokens.DefaultLight},
		{"DefaultDark", tokens.DefaultDark},
	} {
		fixed := s.tok.Ramps.Neutral.Step(500)
		for _, g := range variantSurfaces(s.tok) {
			got := color.Magnitude(fixed, g.surface)
			t.Logf("%s: fixed step 500 on %s measures |Lc| %.2f; OutlineVariant |Lc| %.2f",
				s.name, g.name, got, color.Magnitude(s.tok.OutlineVariant(), g.surface))
			if got < tokens.GraphicFloor {
				fixedFails = true
			}
		}
	}
	if !fixedFails {
		t.Error("the fixed step-500 outline clears the boundary floor on every surface of both default schemes; the derived step's whole reason is that it does not, so variants.go needs rewriting before this passes")
	}
}

// TestContainerSpellingsAgree holds the container family to one derivation.
// StatusContainer and OnStatusContainer are the status family's names for
// Container and OnContainer; two spellings that answered differently for any
// role would be two recipes wearing one word.
//
// The accent trio is asserted alongside the status four because the general
// spelling exists for it: the walk asks the role's ramp and not a table, so a
// role with a ramp has a container whether or not it is a status.
func TestContainerSpellingsAgree(t *testing.T) {
	for _, s := range []struct {
		name string
		tok  tokens.ColorTokens
	}{
		{"DefaultLight", tokens.DefaultLight},
		{"DefaultDark", tokens.DefaultDark},
	} {
		for _, r := range []struct {
			name string
			role tokens.Role
		}{
			{"Neutral", tokens.RoleNeutral},
			{"Primary", tokens.RolePrimary},
			{"Secondary", tokens.RoleSecondary},
			{"Tertiary", tokens.RoleTertiary},
			{"Error", tokens.RoleError},
			{"Success", tokens.RoleSuccess},
			{"Warning", tokens.RoleWarning},
			{"Info", tokens.RoleInfo},
		} {
			if got, want := s.tok.StatusContainer(r.role), s.tok.Container(r.role); got != want {
				t.Errorf("%s %s: StatusContainer = %v, Container = %v", s.name, r.name, got, want)
			}
			if got, want := s.tok.OnStatusContainer(r.role), s.tok.OnContainer(r.role); got != want {
				t.Errorf("%s %s: OnStatusContainer = %v, OnContainer = %v", s.name, r.name, got, want)
			}
			for _, lv := range []tokens.ElevationLevel{tokens.LevelBackdrop, tokens.LevelChrome, tokens.Level0, tokens.Level1, tokens.Level2, tokens.Level3} {
				surface := s.tok.SurfaceAt(lv)
				if got, want := s.tok.StatusContainerOn(r.role, surface), s.tok.ContainerOn(r.role, surface); got != want {
					t.Errorf("%s %s on %v: StatusContainerOn = %v, ContainerOn = %v", s.name, r.name, lv, got, want)
				}
			}
			if got := color.Magnitude(s.tok.OnContainer(r.role), s.tok.Container(r.role)); got < tokens.GraphicFloor {
				t.Errorf("%s %s: OnContainer on Container measures |Lc| %.2f, under the %.0f mark floor",
					s.name, r.name, got, tokens.GraphicFloor)
			}
		}
	}
}
