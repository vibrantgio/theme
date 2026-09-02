package tokens_test

import (
	"fmt"
	stdcolor "image/color"
	"testing"

	"github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/tokens"
)

// variantGrounds are the two neutral grounds both variant rungs are chosen
// against: the surface a piece of furniture wears and the paper the window
// stands on. A token that cleared only the easier of the two would be a
// promise kept on one of them.
func variantGrounds(t tokens.ColorTokens) []struct {
	name   string
	ground stdcolor.NRGBA
} {
	return []struct {
		name   string
		ground stdcolor.NRGBA
	}{
		{"Surface", t.Surface},
		{"Background", t.Background},
	}
}

// TestVariantRungsClearTheirFloors is the gate the two variant rungs exist
// for: whatever the seed, the scheme or the derivation, the boundary rung
// reaches the 3:1 a mark owes its ground and the muted ink reaches the 4.5:1
// a run of words owes it — against BOTH neutral grounds, not the easier one.
//
// It is the failure M3's fixed outline role has and this palette must not:
// a rung named once states one colour and as many measurements as it has
// grounds.
func TestVariantRungsClearTheirFloors(t *testing.T) {
	worstOutline, worstOutlineAt := 99.0, ""
	worstInk, worstInkAt := 99.0, ""
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
			for _, g := range variantGrounds(s.tok) {
				if got := color.ContrastRatio(s.tok.OutlineVariant(), g.ground); got < tokens.GraphicFloor {
					t.Errorf("seed %v: %s OutlineVariant on %s measures %.2f:1, under the %.1f:1 boundary floor",
						seed, s.name, g.name, got, tokens.GraphicFloor)
				} else if got < worstOutline {
					worstOutline, worstOutlineAt = got, fmt.Sprintf("%s on %s", s.name, g.name)
				}
				if got := color.ContrastRatio(s.tok.OnSurfaceVariant(), g.ground); got < tokens.TextFloor {
					t.Errorf("seed %v: %s OnSurfaceVariant on %s measures %.2f:1, under the %.1f:1 text floor",
						seed, s.name, g.name, got, tokens.TextFloor)
				} else if got < worstInk {
					worstInk, worstInkAt = got, fmt.Sprintf("%s on %s", s.name, g.name)
				}
			}
		}
	}
	t.Logf("over %d seeds: worst OutlineVariant %.2f:1 (%s), worst OnSurfaceVariant %.2f:1 (%s)",
		len(sweepSeeds()), worstOutline, worstOutlineAt, worstInk, worstInkAt)
}

// TestVariantRungsStayQuiet holds the other half of each rung's contract.
// Clearing a floor is easy at the ramp's far end — step 900 clears
// everything — and a token that went there would be the loud rung under a
// quiet name. Both are the QUIETEST rung that clears, so each must stay under
// the token that speaks at full strength: the boundary under the muted ink,
// and the muted ink under Text.
func TestVariantRungsStayQuiet(t *testing.T) {
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
			for _, g := range variantGrounds(s.tok) {
				outline := color.ContrastRatio(s.tok.OutlineVariant(), g.ground)
				ink := color.ContrastRatio(s.tok.OnSurfaceVariant(), g.ground)
				text := color.ContrastRatio(s.tok.Text, g.ground)
				if outline > ink {
					t.Errorf("seed %v: %s OutlineVariant on %s measures %.2f:1, louder than OnSurfaceVariant's %.2f:1",
						seed, s.name, g.name, outline, ink)
				}
				if ink >= text {
					t.Errorf("seed %v: %s OnSurfaceVariant on %s measures %.2f:1, at or past Text's %.2f:1 — a muted ink that is not muted",
						seed, s.name, g.name, ink, text)
				}
			}
		}
	}
}

// TestVariantRungsBeatTheFixedStep is the comparison the derivation is worth
// having: naming one neutral rung and calling it the outline — M3's fixed
// role, and what this sheet's borders once did — reads under the boundary
// floor in the light scheme, and the derived token does not.
//
// The numbers are the file header's: step 500 measures 2.35:1 over the harder
// of the two light grounds and 3.07:1 over the harder dark one, so the fixed
// naming fails in the scheme most people read in and passes in the other. The
// test asserts the shape rather than the digits — that the fixed rung falls
// short somewhere the derived rung does not — so a re-derived ramp moves the
// numbers without turning this red for the wrong reason.
func TestVariantRungsBeatTheFixedStep(t *testing.T) {
	fixedFails := false
	for _, s := range []struct {
		name string
		tok  tokens.ColorTokens
	}{
		{"DefaultLight", tokens.DefaultLight},
		{"DefaultDark", tokens.DefaultDark},
	} {
		fixed := s.tok.Ramps.Neutral.Step(500)
		for _, g := range variantGrounds(s.tok) {
			got := color.ContrastRatio(fixed, g.ground)
			t.Logf("%s: fixed step 500 on %s measures %.2f:1; OutlineVariant %.2f:1",
				s.name, g.name, got, color.ContrastRatio(s.tok.OutlineVariant(), g.ground))
			if got < tokens.GraphicFloor {
				fixedFails = true
			}
		}
	}
	if !fixedFails {
		t.Error("the fixed step-500 outline clears the boundary floor on every ground of both default schemes; the derived rung's whole reason is that it does not, so variants.go needs rewriting before this passes")
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
				ground := s.tok.SurfaceAt(lv)
				if got, want := s.tok.StatusContainerOn(r.role, ground), s.tok.ContainerOn(r.role, ground); got != want {
					t.Errorf("%s %s on %v: StatusContainerOn = %v, ContainerOn = %v", s.name, r.name, lv, got, want)
				}
			}
			if got := color.ContrastRatio(s.tok.OnContainer(r.role), s.tok.Container(r.role)); got < tokens.GraphicFloor {
				t.Errorf("%s %s: OnContainer on Container measures %.2f:1, under the %.1f:1 mark floor",
					s.name, r.name, got, tokens.GraphicFloor)
			}
		}
	}
}
