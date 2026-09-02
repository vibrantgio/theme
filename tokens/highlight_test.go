package tokens_test

import (
	stdcolor "image/color"
	"math"
	"testing"

	"github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/tokens"
)

// highlightLevels is the five levels a highlight can be drawn on: it marks
// content, and content stands on the paper, on a card, in a dialog or in a
// popover as readily as on the window's own furniture.
var highlightLevels = []tokens.ElevationLevel{
	tokens.LevelBackdrop, tokens.Level0, tokens.Level1, tokens.Level2, tokens.Level3,
}

// oklabDistance is the Euclidean distance between two colours in OKLab —
// the perceptual distance the gates below read, in the space the whole
// derivation places its hues and chromas in. Hue angle alone cannot judge a
// wash (two hues 90° apart at no chroma are one grey), so the reservation
// is measured both ways.
func oklabDistance(a, b stdcolor.NRGBA) float64 {
	l1, a1, b1 := color.OKLabFromNRGBA(a)
	l2, a2, b2 := color.OKLabFromNRGBA(b)
	return math.Sqrt((l1-l2)*(l1-l2) + (a1-a2)*(a1-a2) + (b1-b2)*(b1-b2))
}

// schemesOf returns the four schemes a whole-population gate reads: both
// schemes of both derivations, from one seed.
func schemesOf(seed stdcolor.NRGBA) []struct {
	name string
	tok  tokens.ColorTokens
} {
	light, dark := tokens.FromSeed(seed)
	hcLight, hcDark := tokens.FromSeedHighContrast(seed)
	return []struct {
		name string
		tok  tokens.ColorTokens
	}{
		{"FromSeed light", light}, {"FromSeed dark", dark},
		{"FromSeedHighContrast light", hcLight}, {"FromSeedHighContrast dark", hcDark},
	}
}

// TestTheHighlightIsReservedAgainstTheBrand verifies the highlighter does
// not rotate with the seed: it is reserved outside the role table, so every
// seed derives the same two washes in both derivations. A highlighter that
// followed the brand would be an accent under another name.
func TestTheHighlightIsReservedAgainstTheBrand(t *testing.T) {
	lightWash := stdcolor.NRGBA{0xe6, 0xcb, 0xee, 0xff}
	darkWash := stdcolor.NRGBA{0x3b, 0x26, 0x41, 0xff}
	for _, seed := range sweepSeeds() {
		for i, s := range schemesOf(seed) {
			want := lightWash
			if i%2 == 1 {
				want = darkWash
			}
			if s.tok.Highlight != want {
				t.Fatalf("seed %v: %s Highlight = %v, want the reserved %v — the highlighter does not rotate with the brand",
					seed, s.name, s.tok.Highlight, want)
			}
		}
	}
}

// TestTheHighlightClearsTheSeamFloorOnEveryLevel verifies the wash is
// visible as a field wherever content is marked: it clears ContainerFloor
// against every level's fill, in both schemes of both derivations, for
// every seed — and stays under the ratio at which a tint stops being the
// ground of something and becomes a fill in its own right.
func TestTheHighlightClearsTheSeamFloorOnEveryLevel(t *testing.T) {
	// The threshold TestContainersSeparateFromEveryLevelItStandsOn reads a
	// container against: past this a wash is a control's fill, not a mark
	// on content.
	const solid = 2.5
	worst, loudest := 99.0, 0.0
	for _, seed := range sweepSeeds() {
		for _, s := range schemesOf(seed) {
			for _, lv := range highlightLevels {
				surface := s.tok.SurfaceAt(lv)
				wash := s.tok.HighlightOn(surface)
				got := color.ContrastRatio(wash, surface)
				if got < tokens.ContainerFloor {
					t.Errorf("seed %v: %s highlight %v on the level-%d fill %v measures %.3f:1, under the %.2f:1 seam floor",
						seed, s.name, wash, lv, surface, got, tokens.ContainerFloor)
				} else if got < worst {
					worst = got
				}
				if got > solid {
					t.Errorf("seed %v: %s highlight %v on the level-%d fill %v measures %.3f:1 — that is a fill, not a mark on content",
						seed, s.name, wash, lv, surface, got)
				} else if got > loudest {
					loudest = got
				}
			}
		}
	}
	t.Logf("over %d seeds, both derivations, both schemes, five levels: worst seam %.3f:1 (floor %.2f), loudest %.3f:1",
		len(sweepSeeds()), worst, tokens.ContainerFloor, loudest)
}

// TestContentInkClearsItsFloorOverTheHighlight verifies a highlight never
// costs the content it marks its legibility: the scheme's own body ink
// clears TextFloor over every wash the walk returns, on every level, in
// both schemes of both derivations. A highlight is applied to content, so
// the words it covers are the whole point of it.
func TestContentInkClearsItsFloorOverTheHighlight(t *testing.T) {
	worst := 99.0
	for _, seed := range sweepSeeds() {
		for _, s := range schemesOf(seed) {
			for _, lv := range highlightLevels {
				wash := s.tok.HighlightOn(s.tok.SurfaceAt(lv))
				got := color.ContrastRatio(s.tok.Text, wash)
				if got < tokens.TextFloor {
					t.Errorf("seed %v: %s Text %v over the level-%d highlight %v measures %.3f:1, under the %.1f:1 text floor",
						seed, s.name, s.tok.Text, lv, wash, got, tokens.TextFloor)
				} else if got < worst {
					worst = got
				}
			}
		}
	}
	t.Logf("over %d seeds, both derivations, both schemes, five levels: worst body ink over a highlight %.3f:1 (floor %.1f)",
		len(sweepSeeds()), worst, tokens.TextFloor)
}

// TestTheHighlightKeepsItsDistanceFromEveryStatus is the gate the
// reservation exists for: a highlight reports no status, so no status may
// be read off it. It measures the realized wash against every status
// colour a reader could see beside it — each role's fixed container, the
// container resolved for the level the wash is on, and the role's pinned
// base — in OKLCh hue and in OKLab distance, over the whole seed sweep in
// both schemes of both derivations.
//
// The bounds are the sweep's own measurements less a rounding margin, and
// they are read against what the palette already asks a reader to tell
// apart: the two closest status containers come to 19.27° and 0.0183 of
// each other (the bent warning beside the error), while the highlight
// stands 65.92° and 0.0606 from the nearest of them. Nothing on this
// palette is as far from a status colour as the highlight is.
func TestTheHighlightKeepsItsDistanceFromEveryStatus(t *testing.T) {
	const (
		hueBound = 60.0  // measured 65.92° over the sweep
		labBound = 0.055 // measured 0.0606 over the sweep
	)
	worstHue, worstHueAt := 999.0, ""
	worstLab, worstLabAt := 99.0, ""
	for _, seed := range sweepSeeds() {
		for _, s := range schemesOf(seed) {
			for _, lv := range highlightLevels {
				surface := s.tok.SurfaceAt(lv)
				wash := s.tok.HighlightOn(surface)
				_, _, washHue := color.OKLChFromNRGBA(wash)
				for _, r := range statusRoles {
					against := []struct {
						what string
						c    stdcolor.NRGBA
					}{
						{"container", s.tok.StatusContainer(r.role)},
						{"container on this level", s.tok.StatusContainerOn(r.role, surface)},
						{"pin", statusPin(s.tok, r.role)},
					}
					for _, a := range against {
						_, _, hue := color.OKLChFromNRGBA(a.c)
						if got := hueGap(washHue, hue); got < worstHue {
							worstHue = got
							worstHueAt = a.what
							if got < hueBound {
								t.Errorf("seed %v: %s highlight %v on level %d is %.2f° from the %s %s %v — under the %.1f° the reservation holds",
									seed, s.name, wash, lv, got, r.name, a.what, a.c, hueBound)
							}
						}
						if got := oklabDistance(wash, a.c); got < worstLab {
							worstLab = got
							worstLabAt = a.what
							if got < labBound {
								t.Errorf("seed %v: %s highlight %v on level %d is %.4f from the %s %s %v in OKLab — under the %.3f the reservation holds",
									seed, s.name, wash, lv, got, r.name, a.what, a.c, labBound)
							}
						}
					}
				}
			}
		}
	}
	t.Logf("over %d seeds, both derivations, both schemes, five levels: the highlight's closest approach to a status colour is %.2f° (%s) and %.4f in OKLab (%s)",
		len(sweepSeeds()), worstHue, worstHueAt, worstLab, worstLabAt)
}

// statusPin answers a status role's pinned base, which ColorTokens carries
// under the role's own name rather than behind an accessor.
func statusPin(t tokens.ColorTokens, role tokens.Role) stdcolor.NRGBA {
	switch role {
	case tokens.RoleError:
		return t.Error
	case tokens.RoleSuccess:
		return t.Success
	case tokens.RoleWarning:
		return t.Warning
	case tokens.RoleInfo:
		return t.Info
	}
	panic("highlight_test: statusPin: not a status role")
}

// TestHighlightOnHoldsTheResolvedWashWhereItAlreadyWorks pins the
// relationship between the field and the walk: ColorTokens.Highlight is
// HighlightOn against the paper, and the walk moves off that realization
// only where the level has walked into it — so content on the paper and
// content on a card are not marked in two different washes side by side
// for no reason.
func TestHighlightOnHoldsTheResolvedWashWhereItAlreadyWorks(t *testing.T) {
	moved := 0
	for _, seed := range sweepSeeds() {
		for _, s := range schemesOf(seed) {
			if got := s.tok.HighlightOn(s.tok.Background); got != s.tok.Highlight {
				t.Fatalf("seed %v: %s HighlightOn(Background) = %v but Highlight = %v; the field is the walk against the paper",
					seed, s.name, got, s.tok.Highlight)
			}
			for _, lv := range highlightLevels {
				surface := s.tok.SurfaceAt(lv)
				if s.tok.HighlightOn(surface) == s.tok.Highlight {
					continue
				}
				moved++
				if color.ContrastRatio(s.tok.Highlight, surface) >= tokens.ContainerFloor {
					t.Errorf("seed %v: %s highlight walked off %v on the level-%d fill %v, which it already cleared at %.3f:1",
						seed, s.name, s.tok.Highlight, lv, surface, color.ContrastRatio(s.tok.Highlight, surface))
				}
			}
		}
	}
	t.Logf("over %d seeds, both derivations, both schemes, five levels: the wash deepens off its resolved realization in %d pairings",
		len(sweepSeeds()), moved)
}
