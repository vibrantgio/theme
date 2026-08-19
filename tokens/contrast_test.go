package tokens_test

import (
	"fmt"
	stdcolor "image/color"
	"math"
	"testing"

	"github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/tokens"
)

// TestAPCAContrastGate is ADR-007's contrast gate over the default palette,
// in both modes: in every role ramp, step 900 must reach |Lc| ≥ 90 and step
// 700 |Lc| ≥ 60 over the step-100 and step-200 grounds, and each pinned
// base's on-colour |Lc| ≥ 60 over the base.
//
// Reading: ADR-007's sentence — "in both ramps, step 900 must reach Lc 90
// and step 700 Lc 60 over the step-100 and step-200 grounds" — is read with
// the grounds taken from the SAME role's ramp, because the ADR assigns
// 700–900 the job "text over tinted fills and pressed states" and the
// tinted fills 100–300 come from the ramp being read. Since every ramp
// shares one lightness scale, the neutral-grounds reading differs only by
// hue-induced luminance wiggle; the same-role reading covers neutral anyway
// (neutral is one of the seven gated ramps).
//
// WCAG 2 ratios for the same pairs are logged alongside — conformance
// claims cite them per ADR-007 — but they do not gate: only APCA failures
// fail this test, so a WCAG regression shows up in the log, never as a
// verdict.
func TestAPCAContrastGate(t *testing.T) {
	for _, s := range []struct {
		name string
		tok  tokens.ColorTokens
	}{
		{"DefaultLight", tokens.DefaultLight},
		{"DefaultDark", tokens.DefaultDark},
	} {
		t.Run(s.name, func(t *testing.T) {
			for _, r := range namedRamps(s.tok) {
				for _, tc := range []struct {
					textStep int
					minLc    float64
				}{
					{900, 90},
					{700, 60},
				} {
					for _, groundStep := range []int{100, 200} {
						text := r.ramp.Step(tc.textStep)
						ground := r.ramp.Step(groundStep)
						lc := color.APCA(text, ground)
						wcag := color.ContrastRatio(text, ground)
						t.Logf("%s %d on %d: Lc %.2f (gate ≥ %.0f), WCAG %.2f:1 (AA %.1f:1: %s, cited not gating)",
							r.name, tc.textStep, groundStep, lc, tc.minLc, wcag, wcagAA, wcagVerdict(wcag))
						if math.Abs(lc) < tc.minLc {
							t.Errorf("%s: step %d on step-%d ground: |Lc| %.2f < %.0f",
								r.name, tc.textStep, groundStep, math.Abs(lc), tc.minLc)
						}
					}
				}
			}
			for _, p := range []struct {
				name     string
				base, on stdcolor.NRGBA
			}{
				{"Primary", s.tok.Primary, s.tok.OnPrimary},
				{"Secondary", s.tok.Secondary, s.tok.OnSecondary},
				{"Tertiary", s.tok.Tertiary, s.tok.OnTertiary},
				{"Error", s.tok.Error, s.tok.OnError},
				{"Success", s.tok.Success, s.tok.OnSuccess},
				{"Warning", s.tok.Warning, s.tok.OnWarning},
			} {
				lc := color.APCA(p.on, p.base)
				wcag := color.ContrastRatio(p.on, p.base)
				t.Logf("pin %s: on-colour Lc %.2f (gate ≥ 60), WCAG %.2f:1 (AA %.1f:1: %s, cited not gating)",
					p.name, lc, wcag, wcagAA, wcagVerdict(wcag))
				if math.Abs(lc) < 60 {
					t.Errorf("pin %s: on-colour |Lc| %.2f < 60", p.name, math.Abs(lc))
				}
			}
		})
	}
}

// wcagVerdict renders a WCAG AA pass/fail for the gate test's log lines —
// reported, never gated on.
func wcagVerdict(ratio float64) string {
	if ratio >= wcagAA {
		return "pass"
	}
	return "fail"
}

// TestAPCAContrastGateHighContrast is the E3.3 gate over the high-contrast
// variant of the default seed, with the variant's floors above the
// defaults': in every role ramp, step 900 must reach |Lc| ≥ 90 as before
// AND step 700 must now also reach |Lc| ≥ 90 (the default gate asks 60)
// over the step-100 and step-200 grounds, and each pinned base's on-colour
// |Lc| ≥ 75 (the default asks 60). Per ADR-007's arrangement WCAG ratios
// are reported alongside — here against AAA (7:1), the level a
// high-contrast conformance claim would cite — but never gated on: only
// APCA failures fail this test.
//
// Measured margins at recording time: light min 700 Lc 90.7, 900 Lc 92.3,
// pins Lc 85.7; dark min 700 Lc 93.1, 900 Lc 104.4, pins Lc 76.3. F4.6's
// Success and Warning ramps join the gate without moving any of those
// minima — each clears its floor above the worst case already recorded.
func TestAPCAContrastGateHighContrast(t *testing.T) {
	hcLight, hcDark := tokens.FromSeedHighContrast(tokens.DefaultSeed)
	for _, s := range []struct {
		name string
		tok  tokens.ColorTokens
	}{
		{"HighContrastLight", hcLight},
		{"HighContrastDark", hcDark},
	} {
		t.Run(s.name, func(t *testing.T) {
			for _, r := range namedRamps(s.tok) {
				for _, textStep := range []int{900, 700} {
					for _, groundStep := range []int{100, 200} {
						text := r.ramp.Step(textStep)
						ground := r.ramp.Step(groundStep)
						lc := color.APCA(text, ground)
						wcag := color.ContrastRatio(text, ground)
						t.Logf("%s %d on %d: Lc %.2f (gate ≥ 90), WCAG %.2f:1 (AAA %.1f:1: %s, cited not gating)",
							r.name, textStep, groundStep, lc, wcag, wcagAAA, wcagAAAVerdict(wcag))
						if math.Abs(lc) < 90 {
							t.Errorf("%s: step %d on step-%d ground: |Lc| %.2f < 90",
								r.name, textStep, groundStep, math.Abs(lc))
						}
					}
				}
			}
			for _, p := range []struct {
				name     string
				base, on stdcolor.NRGBA
			}{
				{"Primary", s.tok.Primary, s.tok.OnPrimary},
				{"Secondary", s.tok.Secondary, s.tok.OnSecondary},
				{"Tertiary", s.tok.Tertiary, s.tok.OnTertiary},
				{"Error", s.tok.Error, s.tok.OnError},
				{"Success", s.tok.Success, s.tok.OnSuccess},
				{"Warning", s.tok.Warning, s.tok.OnWarning},
			} {
				lc := color.APCA(p.on, p.base)
				wcag := color.ContrastRatio(p.on, p.base)
				t.Logf("pin %s: on-colour Lc %.2f (gate ≥ 75), WCAG %.2f:1 (AAA %.1f:1: %s, cited not gating)",
					p.name, lc, wcag, wcagAAA, wcagAAAVerdict(wcag))
				if math.Abs(lc) < 75 {
					t.Errorf("pin %s: on-colour |Lc| %.2f < 75", p.name, math.Abs(lc))
				}
			}
		})
	}
}

// TestInverseSurfaceBodyTextContrast gates the inverse pair in every
// scheme the seed pipeline derives — both default schemes and both
// high-contrast ones — at WCAG AA for body text (4.5:1). The inverse
// surface is what a transient message stands on, so its on-colour is read
// as running text and the body-text ratio is the bar that matters; unlike
// the ramp gates above this one is WCAG-gated rather than WCAG-reported,
// because 4.5:1 is the number the role's contract is written in. APCA is
// logged alongside, the mirror of the arrangement the other gates use.
//
// The pair is also, by derivation, the counterpart scheme's own Surface
// and Text — so a failure here is a failure of that scheme's reading pair,
// not of a separate approximation. The test asserts the derivation too:
// measuring a pair that had quietly stopped being the counterpart's would
// prove nothing about the chip a light scheme actually paints.
func TestInverseSurfaceBodyTextContrast(t *testing.T) {
	hcLight, hcDark := tokens.FromSeedHighContrast(tokens.DefaultSeed)
	for _, s := range []struct {
		name             string
		tok, counterpart tokens.ColorTokens
	}{
		{"DefaultLight", tokens.DefaultLight, tokens.DefaultDark},
		{"DefaultDark", tokens.DefaultDark, tokens.DefaultLight},
		{"HighContrastLight", hcLight, hcDark},
		{"HighContrastDark", hcDark, hcLight},
	} {
		t.Run(s.name, func(t *testing.T) {
			wcag := color.ContrastRatio(s.tok.OnInverseSurface, s.tok.InverseSurface)
			lc := color.APCA(s.tok.OnInverseSurface, s.tok.InverseSurface)
			t.Logf("inverse pair %v on %v: WCAG %.2f:1 (gate ≥ %.1f:1), Lc %.2f (reported)",
				s.tok.OnInverseSurface, s.tok.InverseSurface, wcag, wcagAA, lc)
			if wcag < wcagAA {
				t.Errorf("inverse pair: WCAG %.2f:1 < %.1f:1 — body text on the inverse surface is unreadable",
					wcag, wcagAA)
			}
			if got, want := s.tok.InverseSurface, s.counterpart.Surface; got != want {
				t.Errorf("InverseSurface = %v, want the counterpart scheme's Surface %v", got, want)
			}
			if got, want := s.tok.OnInverseSurface, s.counterpart.Text; got != want {
				t.Errorf("OnInverseSurface = %v, want the counterpart scheme's Text %v", got, want)
			}
		})
	}
}

// accentPairs returns one scheme's pinned accent pairings: each base with
// the ink the derivation chose to read it in.
func accentPairs(t tokens.ColorTokens) []struct {
	name     string
	base, on stdcolor.NRGBA
} {
	return []struct {
		name     string
		base, on stdcolor.NRGBA
	}{
		{"Primary", t.Primary, t.OnPrimary},
		{"Secondary", t.Secondary, t.OnSecondary},
		{"Tertiary", t.Tertiary, t.OnTertiary},
		{"Error", t.Error, t.OnError},
		{"Success", t.Success, t.OnSuccess},
		{"Warning", t.Warning, t.OnWarning},
	}
}

// TestAccentOnColoursClearTheFloorForEverySeed is the whole-population gate
// on the on-colour rule: over the shared seed sweep, in both schemes of both
// derivations, every pinned accent pairing reaches WCAG AA for body text.
//
// It is the property the rule exists for. The bases are pinned to depths
// their usual ink clears — except the light primary base, which is the brand
// colour itself and can land anywhere on the axis, which is how a light
// brand colour used to come back under white text at 2.1:1. The ink is
// chosen by measurement now, and because the two candidates are the ends of
// the tonal axis, the better of them clears 4.5:1 over any colour whatever:
// no seed can produce a pairing this gate has to fail.
//
// Three further properties are asserted alongside the number, because a
// number alone would not notice them going:
//
//   - The ink is always one of the two ends on offer, and where the
//     preferred one falls short the chosen one reads at least as well. A
//     rule that flipped an ink into a worse pairing would still clear the
//     floor most of the time.
//   - Nothing moves for a base whose usual ink already clears the floor.
//     This is what keeps every downstream golden in the design system on
//     the canonical seed exactly where it was.
//   - The increased-contrast variant never reads below the default's, which
//     is what its stricter floor is worth on a pairing whose two candidates
//     are already the ends of the axis.
func TestAccentOnColoursClearTheFloorForEverySeed(t *testing.T) {
	worstLight, worstDark := 99.0, 99.0
	flips := 0
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		hcLight, hcDark := tokens.FromSeedHighContrast(seed)
		for _, s := range []struct {
			name string
			tok  tokens.ColorTokens
			// variant is the increased-contrast scheme this one is asked to
			// be no better than; for the variant's own rows it is itself.
			variant tokens.ColorTokens
			light   bool
			// deflt marks the rows FromSeed derived, the only ones the
			// no-op guarantee is written about: the variant's floor is its
			// own and higher, so it is entitled to question an ink FromSeed
			// is satisfied with.
			deflt bool
		}{
			{"FromSeed light", light, hcLight, true, true},
			{"FromSeed dark", dark, hcDark, false, true},
			{"FromSeedHighContrast light", hcLight, hcLight, true, false},
			{"FromSeedHighContrast dark", hcDark, hcDark, false, false},
		} {
			for i, p := range accentPairs(s.tok) {
				got := color.ContrastRatio(p.on, p.base)
				if got < wcagAA {
					t.Errorf("seed %v: %s %s: %v on %v measures %.2f:1, under the %.1f:1 floor",
						seed, s.name, p.name, p.on, p.base, got, wcagAA)
				}
				if s.light {
					if got < worstLight {
						worstLight = got
					}
					// The light scheme's two candidates are the ends of the
					// axis, so the ink is one of them, and White is what a
					// pairing that already clears the floor keeps.
					if p.on != tokens.White && p.on != tokens.Black {
						t.Errorf("seed %v: %s %s: ink %v is neither end of the tonal axis",
							seed, s.name, p.name, p.on)
					}
					if p.on == tokens.Black {
						flips++
						white := color.ContrastRatio(tokens.White, p.base)
						if got < white {
							t.Errorf("seed %v: %s %s: ink flipped to Black at %.2f:1, worse than White's %.2f:1",
								seed, s.name, p.name, got, white)
						}
						if s.deflt && white >= wcagAA {
							t.Errorf("seed %v: %s %s: ink flipped to Black though White measured %.2f:1 — a passing pairing moved",
								seed, s.name, p.name, white)
						}
					}
				} else if got < worstDark {
					worstDark = got
				}
				// The variant asks more of the same pairing and can never
				// answer with less.
				v := accentPairs(s.variant)[i]
				if hc := color.ContrastRatio(v.on, v.base); hc < got-1e-9 {
					t.Errorf("seed %v: %s %s: the high-contrast variant measures %.2f:1, under the default's %.2f:1",
						seed, s.name, p.name, hc, got)
				}
			}
		}
	}
	t.Logf("over %d seeds: worst light accent pairing %.2f:1, worst dark %.2f:1; %d light inks flipped to Black",
		len(sweepSeeds()), worstLight, worstDark, flips)
}

// TestContainerAndFillPairingsClearTheFloorForEverySeed gates the pairings
// the accents' containers and fills are made of, which the on-colour rule
// deliberately leaves alone: the ramps' own text steps over their own
// tinted grounds, in every role ramp of every scheme the pipeline derives.
//
// They need no rule of their own and the gate says why: a ramp step is
// realized at a fixed CIELAB depth, and depth is what luminance is, so a
// step's contrast against another step of the same ramp is the same
// measurement for every seed there is. The sweep is here to hold that claim
// rather than to look for a seed that breaks it.
func TestContainerAndFillPairingsClearTheFloorForEverySeed(t *testing.T) {
	worst, worstAt := 99.0, ""
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
			for _, r := range namedRamps(s.tok) {
				for _, text := range []int{700, 900} {
					for _, ground := range []int{100, 200} {
						got := color.ContrastRatio(r.ramp.Step(text), r.ramp.Step(ground))
						if got < wcagAA {
							t.Errorf("seed %v: %s %s step %d on step %d measures %.2f:1, under the %.1f:1 floor",
								seed, s.name, r.name, text, ground, got, wcagAA)
						}
						if got < worst {
							worst, worstAt = got, fmt.Sprintf("%s %s %d on %d", s.name, r.name, text, ground)
						}
					}
				}
			}
		}
	}
	t.Logf("over %d seeds: worst container/fill pairing %.2f:1 (%s)", len(sweepSeeds()), worst, worstAt)
}

// wcagAAA is WCAG 2's AAA normal-text ratio, reported (never gated on) by
// the high-contrast gate's log lines.
const wcagAAA = 7.0

// wcagAAAVerdict renders a WCAG AAA pass/fail for those log lines.
func wcagAAAVerdict(ratio float64) string {
	if ratio >= wcagAAA {
		return "pass"
	}
	return "fail"
}
