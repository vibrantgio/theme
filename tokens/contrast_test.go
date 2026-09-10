package tokens_test

import (
	"fmt"
	stdcolor "image/color"
	"math"
	"testing"

	"github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/tokens"
)

// TestAPCAContrastGate is the contrast gate over the default palette,
// in both modes: in every role ramp, step 900 must reach |Lc| ≥ 90 and step
// 700 |Lc| ≥ 60 over the step-100 and step-200 surfaces, and each pinned
// base's on-colour |Lc| ≥ 60 over the base.
//
// Reading: the surfaces are taken from the SAME role's ramp, because steps
// 700–900 carry the job "text over tinted fills and pressed states" and the
// tinted fills 100–300 come from the ramp being read. Since every ramp
// shares one lightness scale, the neutral-surfaces reading differs only by
// hue-induced luminance wiggle; the same-role reading covers neutral anyway
// (neutral is one of the seven gated ramps).
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
					for _, surfaceStep := range []int{100, 200} {
						text := r.ramp.Step(tc.textStep)
						surface := r.ramp.Step(surfaceStep)
						lc := color.APCA(text, surface)
						t.Logf("%s %d on %d: Lc %.2f (gate ≥ %.0f)",
							r.name, tc.textStep, surfaceStep, lc, tc.minLc)
						if math.Abs(lc) < tc.minLc {
							t.Errorf("%s: step %d on step-%d surface: |Lc| %.2f < %.0f",
								r.name, tc.textStep, surfaceStep, math.Abs(lc), tc.minLc)
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
				{"Info", s.tok.Info, s.tok.OnInfo},
			} {
				lc := color.APCA(p.on, p.base)
				t.Logf("pin %s: on-colour Lc %.2f (gate ≥ 60)", p.name, lc)
				if math.Abs(lc) < 60 {
					t.Errorf("pin %s: on-colour |Lc| %.2f < 60", p.name, math.Abs(lc))
				}
			}
		})
	}
}

// TestAPCAContrastGateHighContrast is the gate over the high-contrast
// variant of the default seed, with the variant's floors above the
// defaults': in every role ramp, step 900 must reach |Lc| ≥ 90 as before
// AND step 700 must now also reach |Lc| ≥ 90 (the default gate asks 60)
// over the step-100 and step-200 surfaces, and each pinned base's on-colour
// |Lc| ≥ 75 (the default asks 60).
//
// Measured margins at recording time: light min 700 Lc 90.7, 900 Lc 92.3,
// pins Lc 85.7; dark min 700 Lc 93.1, 900 Lc 104.4, pins Lc 76.3. The
// Success and Warning ramps clear their floors above those minima.
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
					for _, surfaceStep := range []int{100, 200} {
						text := r.ramp.Step(textStep)
						surface := r.ramp.Step(surfaceStep)
						lc := color.APCA(text, surface)
						t.Logf("%s %d on %d: Lc %.2f (gate ≥ 90)", r.name, textStep, surfaceStep, lc)
						if math.Abs(lc) < 90 {
							t.Errorf("%s: step %d on step-%d surface: |Lc| %.2f < 90",
								r.name, textStep, surfaceStep, math.Abs(lc))
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
				{"Info", s.tok.Info, s.tok.OnInfo},
			} {
				lc := color.APCA(p.on, p.base)
				t.Logf("pin %s: on-colour Lc %.2f (gate ≥ 75)", p.name, lc)
				if math.Abs(lc) < 75 {
					t.Errorf("pin %s: on-colour |Lc| %.2f < 75", p.name, math.Abs(lc))
				}
			}
		})
	}
}

// TestInverseSurfaceBodyTextContrast gates the inverse pair in every
// scheme the seed pipeline derives — both default schemes and both
// high-contrast ones — at [tokens.TextFloor]. The inverse surface is what a
// transient message stands on, so its on-colour is read as running text and
// the body-text floor is the bar that matters.
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
			lc := color.APCA(s.tok.OnInverseSurface, s.tok.InverseSurface)
			t.Logf("inverse pair %v on %v: Lc %.2f (gate |Lc| ≥ %.0f)",
				s.tok.OnInverseSurface, s.tok.InverseSurface, lc, tokens.TextFloor)
			if math.Abs(lc) < tokens.TextFloor {
				t.Errorf("inverse pair: |Lc| %.2f < %.0f — body text on the inverse surface is unreadable",
					math.Abs(lc), tokens.TextFloor)
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
// the foreground the derivation chose to read it in.
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
		{"Info", t.Info, t.OnInfo},
	}
}

// TestAccentOnColoursClearTheFloorForEverySeed is the whole-population gate
// on the on-colour rule: over the shared seed sweep, in both schemes of both
// derivations, every pinned accent pairing reaches the |Lc| the two
// candidates make available.
//
// It is the property the rule exists for. The bases are pinned to depths a
// foreground reads over — except the light primary base, which is the brand
// colour itself and can land anywhere on the axis. The foreground is chosen
// by measurement, and because the light scheme's two candidates are the ends
// of the tonal axis, the better of them reaches |Lc| 54.55 over any sRGB
// colour whatever (measured over a dense grid of the whole cube, the
// binding case a mid green #8EAB87): no seed can produce a pairing under
// axisEndGuarantee.
//
// Three further properties are asserted alongside the number, because a
// number alone would not notice them going:
//
//   - The foreground is always one of the two ends on offer, and it is the
//     one that reads furthest. A rule that flipped a foreground into a
//     worse pairing would still look plausible on most seeds.
//   - The increased-contrast variant never reads below the default's on the
//     same pairing.
//   - The light scheme's flip to Black happens where and only where Black
//     reads further than White, which is the whole of the rule.
func TestAccentOnColoursClearTheFloorForEverySeed(t *testing.T) {
	// axisEndGuarantee is the least |Lc| the better end of the tonal axis
	// reaches over any sRGB colour; see the doc comment for the measurement.
	const axisEndGuarantee = 54.55
	worstLight, worstDark := 999.0, 999.0
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
			// own and higher, so it is entitled to question a foreground FromSeed
			// is satisfied with.
			deflt bool
		}{
			{"FromSeed light", light, hcLight, true, true},
			{"FromSeed dark", dark, hcDark, false, true},
			{"FromSeedHighContrast light", hcLight, hcLight, true, false},
			{"FromSeedHighContrast dark", hcDark, hcDark, false, false},
		} {
			for i, p := range accentPairs(s.tok) {
				got := color.Magnitude(p.on, p.base)
				if s.light {
					if got < axisEndGuarantee {
						t.Errorf("seed %v: %s %s: %v on %v measures |Lc| %.2f, under the %.2f the axis ends guarantee",
							seed, s.name, p.name, p.on, p.base, got, axisEndGuarantee)
					}
					if got < worstLight {
						worstLight = got
					}
					// The light scheme's two candidates are the ends of the
					// axis, so the foreground is one of them, and it is the
					// one that reads furthest.
					if p.on != tokens.White && p.on != tokens.Black {
						t.Errorf("seed %v: %s %s: foreground %v is neither end of the tonal axis",
							seed, s.name, p.name, p.on)
					}
					white, black := color.Magnitude(tokens.White, p.base), color.Magnitude(tokens.Black, p.base)
					if p.on == tokens.Black {
						flips++
						if black < white {
							t.Errorf("seed %v: %s %s: foreground flipped to Black at |Lc| %.2f, under White's %.2f",
								seed, s.name, p.name, black, white)
						}
					} else if white < black {
						t.Errorf("seed %v: %s %s: foreground kept White at |Lc| %.2f though Black reads %.2f",
							seed, s.name, p.name, white, black)
					}
				} else if got < worstDark {
					worstDark = got
				}
				// The variant asks more of the same pairing and can never
				// answer with less.
				v := accentPairs(s.variant)[i]
				if hc := color.Magnitude(v.on, v.base); hc < got-1e-9 {
					t.Errorf("seed %v: %s %s: the high-contrast variant measures |Lc| %.2f, under the default's %.2f",
						seed, s.name, p.name, hc, got)
				}
			}
		}
	}
	t.Logf("over %d seeds: worst light accent pairing |Lc| %.2f, worst dark |Lc| %.2f; %d light foregrounds took Black",
		len(sweepSeeds()), worstLight, worstDark, flips)
}

// TestContainerAndFillPairingsClearTheFloorForEverySeed gates the pairings
// the accents' containers and fills are made of, which the on-colour rule
// deliberately leaves alone: the ramps' own text steps over their own
// tinted surfaces, in every role ramp of every scheme the pipeline derives.
//
// They need no rule of their own and the gate says why: a ramp step is
// realized at a fixed CIELAB depth, and depth is what luminance is, so a
// step's contrast against another step of the same ramp is the same
// measurement for every seed there is. The sweep is here to hold that claim
// rather than to look for a seed that breaks it.
func TestContainerAndFillPairingsClearTheFloorForEverySeed(t *testing.T) {
	// The ramps' own text steps are contracted at |Lc| 60 over their own
	// tinted surfaces, the level TestAPCAContrastGate holds step 700 to;
	// the sweep's worst is 66.56.
	const rampTextFloor = 60.0
	worst, worstAt := 999.0, ""
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
					for _, surface := range []int{100, 200} {
						got := color.Magnitude(r.ramp.Step(text), r.ramp.Step(surface))
						if got < rampTextFloor {
							t.Errorf("seed %v: %s %s step %d on step %d measures |Lc| %.2f, under the %.0f floor",
								seed, s.name, r.name, text, surface, got, rampTextFloor)
						}
						if got < worst {
							worst, worstAt = got, fmt.Sprintf("%s %s %d on %d", s.name, r.name, text, surface)
						}
					}
				}
			}
		}
	}
	t.Logf("over %d seeds: worst container/fill pairing |Lc| %.2f (%s)", len(sweepSeeds()), worst, worstAt)
}

// statusRoles is the four semantic status roles with the fixed anchor hue
// each is derived from, in the order the derivation builds them. Each anchor
// is one hue at every depth: a semantic colour that changed family as it
// deepened would be a different colour at every size it was drawn.
var statusRoles = []struct {
	name   string
	role   tokens.Role
	anchor float64 // the fixed OKLCh anchor hue, before any seed tint
}{
	{"Error", tokens.RoleError, 28.7},
	{"Success", tokens.RoleSuccess, 144.2},
	{"Warning", tokens.RoleWarning, 64.05},
	{"Info", tokens.RoleInfo, 248.8},
}

// statusTintBound is the derivation's own bound on how far a seed may rotate
// a status anchor, restated here so the gates read against a number of their
// own rather than against the one the code under test used.
const statusTintBound = 1.75

// hueReadChroma is the realized OKLCh chroma below which a step's hue cannot
// be read — by these gates or by anyone. A hue is an angle on a circle whose
// radius is the chroma, so eight-bit quantization costs about 0.12/C degrees
// of hue: measured over the four anchors across every step of all four
// lightness scales, a step at chroma 0.0119 reads 3.5° off the hue it was
// asked for and one at 0.0131 reads 4.8° off, while every step at or above
// 0.045 reads within 1.7°. Those pale and near-black steps are exactly the
// ones with no hue to confuse — a colour that has given up its chroma has
// given up its family too — so the per-step gates below skip them and report
// how many they skipped rather than pretending to measure them.
const hueReadChroma = 0.045

// stepSlack is the hue-reading slack the per-step gates allow above that
// threshold: 1.7° measured, 2.0° allowed. What the sweep then reports is a
// worst drift of 3.20° against a 3.75° allowance (1.75° of tint plus this)
// — the whole of that margin being the eight bits, not the derivation.
const stepSlack = 2.0

// paleTintStep is the step a container reads its hue at, restated here so
// the gate reads it off the ramp rather than off the code under test: the
// third step counted from the ramp's pale end.
func paleTintStep(r tokens.Ramp) int {
	pale, _, _ := color.LabFromNRGBA(r.Step(100))
	deep, _, _ := color.LabFromNRGBA(r.Step(900))
	if pale >= deep {
		return 300
	}
	return 700
}

// roleRamp returns one role's ramp out of a scheme.
func roleRamp(t tokens.ColorTokens, role tokens.Role) tokens.Ramp {
	switch role {
	case tokens.RoleError:
		return t.Ramps.Error
	case tokens.RoleSuccess:
		return t.Ramps.Success
	case tokens.RoleWarning:
		return t.Ramps.Warning
	case tokens.RoleInfo:
		return t.Ramps.Info
	case tokens.RoleSecondary:
		return t.Ramps.Secondary
	case tokens.RoleTertiary:
		return t.Ramps.Tertiary
	case tokens.RoleNeutral:
		return t.Ramps.Neutral
	default:
		return t.Ramps.Primary
	}
}

// stepTone reads back the CIELAB depth a realized step sits at, which is
// what the warning family's hue rule keys on.
func stepTone(c stdcolor.NRGBA) int {
	l, _, _ := color.LabFromNRGBA(c)
	return int(math.Round(l))
}

// hueGap is the angular distance between two OKLCh hues in degrees, in
// [0,180] — the shorter way round the circle.
func hueGap(a, b float64) float64 {
	d := math.Abs(math.Mod(a-b+540, 360) - 180)
	return d
}

// roleHue reads a role's realized hue off its ramp's mid-value step, which
// is the step furthest from both ends of the lightness scale and so the one
// least disturbed by gamut mapping or 8-bit quantization.
func roleHue(t tokens.ColorTokens, role tokens.Role) float64 {
	var r tokens.Ramp
	switch role {
	case tokens.RoleError:
		r = t.Ramps.Error
	case tokens.RoleSuccess:
		r = t.Ramps.Success
	case tokens.RoleWarning:
		r = t.Ramps.Warning
	case tokens.RoleInfo:
		r = t.Ramps.Info
	default:
		r = t.Ramps.Primary
	}
	_, _, h := color.OKLChFromNRGBA(r.Step(500))
	return h
}

// The slack the gates below allow a hue measurement, which is entirely the
// eight-bit realization and none of it the derivation. How much a byte costs
// depends on how much chroma the colour has to spend it on: a ramp's
// mid-value step carries most of its role's chroma and pins its hue to
// within a degree, while a container carries the 0.055 dial and pins its hue
// to within about a degree and a half — 1.42° measured over the sweep,
// against the step it is realized from. Both are far below the 3° tint bound
// these gates are actually about, and both are far below anything an eye
// resolves at these chromas.
const (
	quantizationSlack = 1.0 // a hue read off a ramp's mid-value step
	containerSlack    = 2.0 // a hue read off a container, at the container dial
)

// statusSeparation is the floor between two status families' realized hues.
//
// The anchors sit 35.35°, 80.2°, 104.6° and 139.9° apart, so the error and
// the warning are the closest pair the set has: an orange next to a red. A
// seed lying between the two tints them toward each other by the tint bound
// each, leaving 31.85° at the anchors, and the steps realize that to 30.2° at
// the worst chroma a container carries. The gate asks 29°, the narrowest the
// derivation can produce less the hue-reading slack two steps carry between
// them.
//
// 30° is not a gate-shaped compromise but a number rendered and judged, the
// two families side by side at every depth a warning is painted at: at L* 39
// the warning realizes #874e00 against the error's #b1241c, at L* 28 #633800
// against #8b0002, at L* 63 #dd8300 against #fa6b5b. The orange reads as an
// orange at every one of them.
const statusSeparation = 29.0

// TestStatusAnchorsHoldTheirFamiliesForEverySeed is the whole-population
// gate on the status anchors and the bound on the seed's tint, over the
// shared seed sweep, in both schemes of both derivations, at every step of
// every ramp rather than at one reference step.
//
// Three properties, and they are the reasons the anchors exist:
//
//   - Every status step stays within statusTintBound of its own family's
//     anchor, so no brand can rotate a semantic colour out of its family. An
//     error is red under every seed there is, and a warning is orange.
//   - No two status families come within statusSeparation of each other at
//     the same depth.
//   - The error role is never further from true red than the accent is. A
//     red-heavy brand pulls the error onto the red anchor rather than
//     pulling the accent past it, which is what stops a themed accent from
//     out-reddening the colour that means "this went wrong".
//
// Steps under hueReadChroma are skipped rather than measured, for the reason
// recorded on that constant; the count is logged so a derivation that
// quietly drained a family would show up as a skip count that moved.
func TestStatusAnchorsHoldTheirFamiliesForEverySeed(t *testing.T) {
	const redAnchor = 28.7
	worstTint, worstSep, worstRed := 0.0, 999.0, 999.0
	worstRedAt, worstSepAt := "", ""
	measured, skipped := 0, 0
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
			for step := 0; step < 9; step++ {
				// Each family's realized hue at this one depth, or "no hue
				// to read" for a step with too little chroma to carry one.
				hues := make([]float64, len(statusRoles))
				legible := make([]bool, len(statusRoles))
				for i, r := range statusRoles {
					stepColor := roleRamp(s.tok, r.role)[step]
					_, chroma, hue := color.OKLChFromNRGBA(stepColor)
					if chroma < hueReadChroma {
						skipped++
						continue
					}
					measured++
					hues[i], legible[i] = hue, true
					tone := stepTone(stepColor)
					if drift := hueGap(hue, r.anchor); drift > statusTintBound+stepSlack {
						t.Errorf("seed %v: %s %s step %d (L* %d) sits at hue %.1f°, %.1f° off its %.1f° anchor — past the %.1f° tint bound",
							seed, s.name, r.name, (step+1)*100, tone, hue, drift, r.anchor, statusTintBound)
					} else if drift > worstTint {
						worstTint = drift
					}
				}
				for i := range hues {
					for j := i + 1; j < len(hues); j++ {
						if !legible[i] || !legible[j] {
							continue
						}
						gap := hueGap(hues[i], hues[j])
						if gap < statusSeparation {
							t.Errorf("seed %v: %s %s and %s are %.1f° apart at step %d, under the %.1f° floor",
								seed, s.name, statusRoles[i].name, statusRoles[j].name, gap, (step+1)*100, statusSeparation)
						}
						if gap < worstSep {
							worstSep, worstSepAt = gap, fmt.Sprintf("%s and %s", statusRoles[i].name, statusRoles[j].name)
						}
					}
				}
			}
			// The accent-versus-error gate. Redness is read as distance to
			// the fixed red anchor: the error role must be at least as close
			// to it as the accent is. A seed inside the error's own tint
			// window puts the two on the same hue ray, where the two
			// realizations differ by a fraction of a degree and the sign of
			// the margin is eight-bit noise — hence the slack, which is what
			// the narrowest margins logged below are made of. Both are read
			// off their ramps' mid-value step, the step furthest from either
			// end of the lightness scale and so the least disturbed by gamut
			// mapping.
			errRed := hueGap(roleHue(s.tok, tokens.RoleError), redAnchor)
			accentRed := hueGap(roleHue(s.tok, tokens.RolePrimary), redAnchor)
			if errRed > accentRed+quantizationSlack {
				t.Errorf("seed %v: %s: the accent sits %.1f° from true red and the error %.1f° — the accent reads redder than the error",
					seed, s.name, accentRed, errRed)
			}
			if margin := accentRed - errRed; margin < worstRed {
				worstRed, worstRedAt = margin, fmt.Sprintf("seed %v, %s", seed, s.name)
			}
		}
	}
	t.Logf("over %d seeds: %d steps measured, %d skipped under chroma %.3f; worst tint drift %.2f° (bound %.1f°), closest pair %.1f° apart (%s, floor %.1f°), narrowest accent-versus-error margin %.2f° (%s)",
		len(sweepSeeds()), measured, skipped, hueReadChroma, worstTint, statusTintBound,
		worstSep, worstSepAt, statusSeparation, worstRed, worstRedAt)
}

// TestStatusContainersKeepTheirParentsHue is the whole-population gate on
// the tonal container derivation, over the shared seed sweep, in both
// schemes of both derivations. It is the gate the blended treatment could
// not have passed: its fills measured chroma 0.0155–0.0212 against a dial of
// 0.055, and its red fill's hue landed 7° off its parent's.
//
// Four properties per container: it carries its parent role's hue, it
// carries the container dial's chroma, the neutral body-text token reaches
// [tokens.TextFloor] over it, and the mark the derivation chose for it
// reaches [tokens.GraphicFloor] over it.
//
// "Its parent's hue" is read off the step the container reads its own hue
// from — the ramp's pale tint depth, step 300 counted from the pale end —
// because a container's hue is its role's and not its depth's (containers.go).
// Measuring it against the step it is realized at instead would gate the
// bend, which is a rule for marks and not for containers: the whole point of
// the derivation is that a role's container is one hue at every depth it is
// drawn at.
func TestStatusContainersKeepTheirParentsHue(t *testing.T) {
	const dial, dialSlack = 0.055, 0.004
	worstText, worstMark, worstChroma := 999.0, 999.0, 99.0
	worstHue := 0.0
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
			for _, r := range statusRoles {
				container := s.tok.StatusContainer(r.role)
				mark := s.tok.OnStatusContainer(r.role)
				_, chroma, hue := color.OKLChFromNRGBA(container)
				ramp := roleRamp(s.tok, r.role)
				parent := ramp.Step(300)
				_, _, parentHue := color.OKLChFromNRGBA(ramp.Step(paleTintStep(ramp)))
				if drift := hueGap(hue, parentHue); drift > containerSlack {
					t.Errorf("seed %v: %s %s container %v sits at hue %.1f°, %.1f° off the pale tint depth's %.1f° — a container has left its family",
						seed, s.name, r.name, container, hue, drift, parentHue)
				} else if drift > worstHue {
					worstHue = drift
				}
				if got, want := stepTone(container), stepTone(parent); got != want {
					t.Errorf("seed %v: %s %s container %v realizes at L* %d, not its step-300 L* %d",
						seed, s.name, r.name, container, got, want)
				}
				if math.Abs(chroma-dial) > dialSlack {
					t.Errorf("seed %v: %s %s container %v measures chroma %.4f, want the %.3f dial",
						seed, s.name, r.name, container, chroma, dial)
				}
				if chroma < worstChroma {
					worstChroma = chroma
				}
				if got := color.Magnitude(s.tok.Text, container); got < tokens.TextFloor {
					t.Errorf("seed %v: %s %s container: body text measures |Lc| %.2f, under the %.0f floor",
						seed, s.name, r.name, got, tokens.TextFloor)
				} else if got < worstText {
					worstText = got
				}
				if got := color.Magnitude(mark, container); got < tokens.GraphicFloor {
					t.Errorf("seed %v: %s %s container: the mark %v measures |Lc| %.2f, under the %.0f mark floor",
						seed, s.name, r.name, mark, got, tokens.GraphicFloor)
				} else if got < worstMark {
					worstMark = got
				}
			}
		}
	}
	t.Logf("over %d seeds: worst container chroma %.4f (dial %.3f), worst hue drift from the pale tint depth %.2f° (slack %.1f°), worst body text on a container |Lc| %.2f, worst mark on a container |Lc| %.2f",
		len(sweepSeeds()), worstChroma, dial, worstHue, containerSlack, worstText, worstMark)
}

// TestStatusPinsAreTheirRampsSeventhStep holds the status roles' pin rule: a
// status role's pinned base and its ramp's step 700 are one colour, in both
// schemes, for every seed. Realizing them at two depths a single tone apart
// — the pin at MD3's tone 40 and the step at the scale's 39 — puts five of
// six light role fills 3/255 per channel beside the step they claim to be,
// which is a palette shipping two Errors that differ by one percent.
//
// It is asserted of FromSeed only. The increased-contrast variant deepens
// its text steps and leaves its fills where they are, exactly as it leaves
// the light pins' White on-colours alone, so its 700 step is deliberately
// not its pin.
func TestStatusPinsAreTheirRampsSeventhStep(t *testing.T) {
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		for _, s := range []struct {
			name string
			tok  tokens.ColorTokens
		}{{"light", light}, {"dark", dark}} {
			for _, p := range []struct {
				name string
				pin  stdcolor.NRGBA
				ramp tokens.Ramp
			}{
				{"Error", s.tok.Error, s.tok.Ramps.Error},
				{"Success", s.tok.Success, s.tok.Ramps.Success},
				{"Warning", s.tok.Warning, s.tok.Ramps.Warning},
				{"Info", s.tok.Info, s.tok.Ramps.Info},
			} {
				if got, want := p.pin, p.ramp.Step(700); got != want {
					t.Errorf("seed %v: %s %s pin = %v, want its ramp's step 700 %v",
						seed, s.name, p.name, got, want)
				}
			}
		}
	}
}

// TestContainersSeparateFromEveryLevelItStandsOn is the gate on the
// surface-aware member of the container family: whatever level a component
// places its tonal container on, the fill is a visibly different surface from
// that surface — and is still only a tint, never a solid.
//
// The fixed-depth container cannot pass this and is not asked to: it is
// realized at one step, and the elevation levels walk THROUGH that step, so
// a dark scheme's level-2 surface and its own step-300 container are the same
// colour to within 1.01:1. That is the whole reason StatusContainerOn exists,
// and the two bounds below are what it buys.
//
// RoleNeutral is swept alongside the four statuses because the derivation
// asks the ramp rather than a table: a neutral container is the neutral
// ramp's own depth at chroma 0, and a component labelling a plain category
// needs it to separate from the page exactly as a status one does.
func TestContainersSeparateFromEveryLevelItStandsOn(t *testing.T) {
	// A container is a tint. Past this it stops being the surface of
	// something and starts being a fill in its own right, which is the
	// register a control occupies.
	const solid = 2.5
	roles := []struct {
		name string
		role tokens.Role
	}{{"Neutral", tokens.RoleNeutral}}
	for _, r := range statusRoles {
		roles = append(roles, struct {
			name string
			role tokens.Role
		}{r.name, r.role})
	}
	levels := []tokens.ElevationLevel{
		tokens.LevelBackdrop, tokens.LevelChrome, tokens.Level0, tokens.Level1, tokens.Level2, tokens.Level3,
	}
	// strongest and faintest are contrast readings, not prominence.
	worst, strongest := 99.0, 0.0
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
			for _, r := range roles {
				for _, lv := range levels {
					surface := s.tok.SurfaceAt(lv)
					fill := s.tok.StatusContainerOn(r.role, surface)
					got := contrastRatio(fill, surface)
					if got < tokens.ContainerFloor {
						t.Errorf("seed %v: %s %s container %v on the level-%d fill %v measures %.3f:1, under the %.2f:1 seam floor",
							seed, s.name, r.name, fill, lv, surface, got, tokens.ContainerFloor)
					} else if got < worst {
						worst = got
					}
					if got > solid {
						t.Errorf("seed %v: %s %s container %v on the level-%d fill %v measures %.3f:1 — that is a fill, not a tint",
							seed, s.name, r.name, fill, lv, surface, got)
					} else if got > strongest {
						strongest = got
					}
				}
			}
		}
	}
	t.Logf("over %d seeds, both derivations, both schemes, six levels: worst seam %.3f:1 (floor %.2f), best %.3f:1",
		len(sweepSeeds()), worst, tokens.ContainerFloor, strongest)
}

// TestStatusFillsKeepTheirHuesApartOnEveryLevel is the whole-population
// gate on the status set as a set. Three floors are read together because a
// container owes all three at once and trading one away for another is the
// defect they exist to catch: it must be a visible field on the surface it
// stands on, it must not be so pronounced that it reads as a control's fill,
// its content must be legible over it, and it must be tellable apart from
// the other three.
//
// The separation floor is the one this gate was written for. It is measured,
// not picked: [tokens.ContainerSeparation] is the light scheme's own closest
// pairing, and the dark scheme measured 0.0183 against it while a container
// read its hue off the step it was realized at — the bent warning beside the
// error, two browns.
func TestStatusFillsKeepTheirHuesApartOnEveryLevel(t *testing.T) {
	// The threshold TestContainersSeparateFromEveryLevelItStandsOn reads a
	// container against: past this a container is a control's fill, not the
	// surface something else stands on.
	const solid = 2.5
	// containerForegroundReach is what a role's own ramp actually reaches
	// over that role's container, which is under [tokens.TextFloor]: a
	// container is realized at the step-300 depth or deeper, and the ramp's
	// text steps reach |Lc| 66.6 at best over the shallower step-100
	// surface. The walk answers with the step that reads furthest and the
	// sweep's worst is 62.54.
	const containerForegroundReach = 62.0
	levels := []tokens.ElevationLevel{
		tokens.LevelBackdrop, tokens.LevelChrome, tokens.Level0, tokens.Level1, tokens.Level2, tokens.Level3,
	}
	// strongestSeam is a contrast reading, not prominence.
	worstSeam, strongestSeam, worstText := 99.0, 0.0, 999.0
	worstSep, worstSepAt := 99.0, ""
	worstHue := 999.0
	for _, seed := range sweepSeeds() {
		for _, s := range schemesOf(seed) {
			for _, lv := range levels {
				surface := s.tok.SurfaceAt(lv)
				fills := make([]stdcolor.NRGBA, len(statusRoles))
				for i, r := range statusRoles {
					fills[i] = s.tok.StatusContainerOn(r.role, surface)
					got := contrastRatio(fills[i], surface)
					if got < tokens.ContainerFloor {
						t.Errorf("seed %v: %s %s container %v on the level-%d surface %v measures %.3f:1, under the %.2f:1 seam floor",
							seed, s.name, r.name, fills[i], lv, surface, got, tokens.ContainerFloor)
					} else if got < worstSeam {
						worstSeam = got
					}
					if got > solid {
						t.Errorf("seed %v: %s %s container %v on the level-%d surface %v measures %.3f:1 — that is a control's fill, not a container",
							seed, s.name, r.name, fills[i], lv, surface, got)
					} else if got > strongestSeam {
						strongestSeam = got
					}
					fg := s.tok.ForegroundOn(r.role, fills[i])
					if got := color.Magnitude(fg, fills[i]); got < containerForegroundReach {
						t.Errorf("seed %v: %s %s foreground %v over its own container %v measures |Lc| %.2f, under the %.0f the role's ramp reaches over its own container",
							seed, s.name, r.name, fg, fills[i], got, containerForegroundReach)
					} else if got < worstText {
						worstText = got
					}
				}
				for i := range fills {
					for j := i + 1; j < len(fills); j++ {
						got := oklabDistance(fills[i], fills[j])
						_, _, hi := color.OKLChFromNRGBA(fills[i])
						_, _, hj := color.OKLChFromNRGBA(fills[j])
						if got < tokens.ContainerSeparation {
							t.Errorf("seed %v: %s on the level-%d surface the %s container %v and the %s container %v are %.4f apart in OKLab (%.2f° of hue), under the %.3f the set owes",
								seed, s.name, lv, statusRoles[i].name, fills[i], statusRoles[j].name, fills[j], got, hueGap(hi, hj), tokens.ContainerSeparation)
						} else if got < worstSep {
							worstSep = got
							worstSepAt = fmt.Sprintf("%s beside %s, %s level %d", statusRoles[i].name, statusRoles[j].name, s.name, lv)
						}
						if h := hueGap(hi, hj); h < worstHue {
							worstHue = h
						}
					}
				}
			}
		}
	}
	t.Logf("over %d seeds, both derivations, both schemes, five levels: closest two containers %.4f in OKLab (floor %.3f, %s) at %.2f° of hue; seam worst %.3f:1 best %.3f:1; worst foreground over a container |Lc| %.2f",
		len(sweepSeeds()), worstSep, tokens.ContainerSeparation, worstSepAt, worstHue, worstSeam, strongestSeam, worstText)
}

// TestTheSurfaceAwareContainerHoldsTheFixedOneWhereItAlreadyWorks pins the
// relationship between the two members: StatusContainerOn moves off
// StatusContainer's depth only where the fixed depth has collided with the
// level. Everywhere else the two are one colour, so a component that names
// its level and one that fills the page do not draw the same role in two
// different tints beside each other.
func TestTheSurfaceAwareContainerHoldsTheFixedOneWhereItAlreadyWorks(t *testing.T) {
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		for _, s := range []struct {
			name string
			tok  tokens.ColorTokens
		}{{"light", light}, {"dark", dark}} {
			for _, r := range statusRoles {
				for _, lv := range []tokens.ElevationLevel{
					tokens.LevelBackdrop, tokens.LevelChrome, tokens.Level0, tokens.Level1, tokens.Level2, tokens.Level3,
				} {
					surface := s.tok.SurfaceAt(lv)
					fixed := s.tok.StatusContainer(r.role)
					if contrastRatio(fixed, surface) < tokens.ContainerFloor {
						continue // the collision case: moving is the point
					}
					if got := s.tok.StatusContainerOn(r.role, surface); got != fixed {
						t.Errorf("seed %v: %s %s on the level-%d fill: the surface-aware container %v left the fixed one %v while the fixed one still cleared the seam floor",
							seed, s.name, r.name, lv, got, fixed)
					}
				}
			}
		}
	}
}

// The four properties the tone curve is asked to hold. All four are gated
// below over every ramp of both schemes, for every seed in the sweep;
// the bracketing measurements are that sweep's own extremes.
const (
	// gapCeiling is the widest a ramp may leave two neighbouring steps in
	// CIELAB L*. The binding case is the light scale's own extreme, where
	// the 900 stop drops to L* 6 to clear the Lc ≥ 90 text gate: 22.33 L*
	// over the sweep. A dark scale read from the same curve tops out at
	// 18.2, between its 600 and 700 steps.
	gapCeiling = 23.0
	// adjacencyFloor is the least two neighbouring steps may measure
	// against each other, which is what makes them two steps rather than
	// one. The binding case is step 100 against step 200 — the window floor
	// under the content, 5 L* apart in both schemes by measurement — which
	// bottoms out at 1.1062:1 over the sweep, in the dark scheme.
	adjacencyFloor = 1.10
	// The two bands a ramp has to put a step in, measured in |Lc| against
	// its own step-100 surface: [tokens.GraphicFloor] for a mark that is not
	// text and [tokens.TextFloor] for one that is. Each band is closed at
	// the next level up, so clearing it takes a step of about the right
	// weight rather than one pronounced enough to clear everything. The
	// boundary tone is light's 600 and dark's 500, measuring |Lc| 46.09 to
	// 50.06 over the sweep; the text tone is light's 800 and dark's 700,
	// |Lc| 75.00 to 85.51.
	boundaryBand = tokens.GraphicFloor
	textBand     = tokens.TextFloor
	strongBand   = tokens.IncreasedTextFloor
)

// TestRampsCoverTheirRange gates the tone curve itself, in every ramp of
// both schemes and for every seed: no two steps further apart than
// gapCeiling, none closer than adjacencyFloor, and a step in each of the
// two contrast bands over the ramp's own surface. Together they say a ramp
// is a progression covering its range rather than two clusters with a hole
// between them — a scale with no step in the mark band has no boundary tone
// to draw an outline in, and one with no step in the text band has no text
// tone.
//
// The increased-contrast variant is deliberately outside this gate: it
// vacates the upper middle by design, deepening its 700 stop to the default
// scale's 900 depth so 700 text clears Lc ≥ 90, which leaves both its
// schemes a gap the default scale would fail on. What the variant owes is
// gated in TestAPCAContrastGateHighContrast.
func TestRampsCoverTheirRange(t *testing.T) {
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		for _, s := range []struct {
			name string
			tok  tokens.ColorTokens
		}{{"light", light}, {"dark", dark}} {
			for _, r := range namedRamps(s.tok) {
				surface := r.ramp.Step(100)
				var boundary, text int
				for step := 200; step <= 900; step += 100 {
					below, at := r.ramp.Step(step-100), r.ramp.Step(step)
					belowL, _, _ := color.LabFromNRGBA(below)
					atL, _, _ := color.LabFromNRGBA(at)
					if gap := math.Abs(atL - belowL); gap > gapCeiling {
						t.Errorf("seed %v: %s %s steps %d–%d are %.1f L* apart, over the %.0f ceiling",
							seed, s.name, r.name, step-100, step, gap, gapCeiling)
					}
					if got := contrastRatio(below, at); got < adjacencyFloor {
						t.Errorf("seed %v: %s %s steps %d–%d measure %.4f:1 against each other, under the %.2f floor",
							seed, s.name, r.name, step-100, step, got, adjacencyFloor)
					}
					switch got := color.Magnitude(at, surface); {
					case got >= boundaryBand && got < textBand && boundary == 0:
						boundary = step
					case got >= textBand && got < strongBand && text == 0:
						text = step
					}
				}
				if boundary == 0 {
					t.Errorf("seed %v: %s %s has no step between |Lc| %.0f and %.0f over its own surface — no boundary tone",
						seed, s.name, r.name, boundaryBand, textBand)
				}
				if text == 0 {
					t.Errorf("seed %v: %s %s has no step between |Lc| %.0f and %.0f over its own surface — no text tone",
						seed, s.name, r.name, textBand, strongBand)
				}
				if seed == tokens.DefaultSeed {
					t.Logf("default seed, %s %s: boundary tone step %d, text tone step %d",
						s.name, r.name, boundary, text)
				}
			}
		}
	}
}

// TestFillsClearThePerceptibilityFloor gates the state fill a control with no
// fill of its own paints on the surface it stands on: over the seed sweep,
// both derivations, both schemes and every level, hover and press each
// separate from that surface by at least tokens.StateFloor, and press lies
// beyond hover.
//
// Both colours in this pairing come off the one neutral scale, so the walk
// is the only one in the package whose own step can be too small to see —
// before the floor, the dark scheme's content hovered at 1.12:1, which is a
// signal that has stopped signalling.
func TestFillsClearThePerceptibilityFloor(t *testing.T) {
	// strongest and faintest are contrast readings, not prominence.
	worst, worstAt, strongest := 99.0, "", 0.0
	worstStep, worstStepAt := 99.0, ""
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
			for _, lv := range levels {
				surface := s.tok.SurfaceAt(lv.level)
				hover := s.tok.StateAt(lv.level, tokens.StateHover)
				press := s.tok.StateAt(lv.level, tokens.StatePressed)
				for _, w := range []struct {
					name string
					fill stdcolor.NRGBA
				}{{"hover", hover}, {"press", press}} {
					got := contrastRatio(w.fill, surface)
					where := fmt.Sprintf("seed %v %s %s %s", seed, s.name, lv.name, w.name)
					if got < tokens.StateFloor {
						t.Errorf("%s: state fill %v on the surface %v measures %.3f:1, under the %.2f:1 floor",
							where, w.fill, surface, got, tokens.StateFloor)
					} else if got < worst {
						worst, worstAt = got, where
					}
					if got > strongest {
						strongest = got
					}
				}
				step := contrastRatio(press, hover)
				if lstar(press) == lstar(hover) {
					t.Errorf("seed %v %s %s: press %v does not lie beyond hover %v",
						seed, s.name, lv.name, press, hover)
				} else if step < worstStep {
					worstStep, worstStepAt = step, fmt.Sprintf("seed %v %s %s", seed, s.name, lv.name)
				}
			}
		}
	}
	t.Logf("over %d seeds, both derivations, both schemes, five levels: worst state fill %.3f:1 (floor %.2f, %s), best %.3f:1; worst press-over-hover %.3f:1 (%s)",
		len(sweepSeeds()), worst, tokens.StateFloor, worstAt, strongest, worstStep, worstStepAt)
}
