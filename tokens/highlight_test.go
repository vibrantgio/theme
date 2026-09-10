package tokens_test

import (
	stdcolor "image/color"
	"math"
	"testing"

	"github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/tokens"
)

// highlightLevels is the five levels a highlight can be drawn on: it marks
// content, and content stands on the content plane, on a card, in a dialog
// or in a popover as readily as in the window's own chrome. The backdrop
// is not among them — nothing stands on it, so no highlight is drawn there.
var highlightLevels = []tokens.ElevationLevel{
	tokens.LevelChrome, tokens.Level0, tokens.Level1, tokens.Level2, tokens.Level3,
}

// markerYellow is the one colour every highlight is laid on in, repeated here
// as the value the exported answers are read against.
var markerYellow = stdcolor.NRGBA{R: 0xff, G: 0xd0, B: 0x00, A: 0xff}

// oklabDistance is the Euclidean distance between two colours in OKLab — the
// perceptual distance the gates below read, in the space the whole derivation
// places its hues and chromas in.
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

// TestTheHighlightIsOneYellowOverAnySurface verifies that what the answers
// depend on is the surface and nothing else: over a given surface every seed,
// both derivations and both schemes lay the same yellow at the same coverage,
// and the fill that lands is the arithmetic blend of the two — 40% for a
// match, 70% for the current one. A highlight that followed the brand would
// be an accent under another name.
func TestTheHighlightIsOneYellowOverAnySurface(t *testing.T) {
	blend := func(coverage float64, surface stdcolor.NRGBA) stdcolor.NRGBA {
		mix := func(src, dst uint8) uint8 {
			return uint8(math.Round(coverage*float64(src) + (1-coverage)*float64(dst)))
		}
		return stdcolor.NRGBA{
			R: mix(markerYellow.R, surface.R),
			G: mix(markerYellow.G, surface.G),
			B: mix(markerYellow.B, surface.B),
			A: 0xff,
		}
	}
	for _, surface := range []stdcolor.NRGBA{
		{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		{R: 0x1e, G: 0x1e, B: 0x1e, A: 0xff},
		{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
	} {
		for _, seed := range sweepSeeds() {
			for _, s := range schemesOf(seed) {
				if got, want := s.tok.HighlightOn(surface), blend(0.40, surface); got != want {
					t.Fatalf("seed %v: %s HighlightOn(%v) = %v, want the yellow at 40%%, %v",
						seed, s.name, surface, got, want)
				}
				if got, want := s.tok.CurrentMatchOn(surface), blend(0.70, surface); got != want {
					t.Fatalf("seed %v: %s CurrentMatchOn(%v) = %v, want the yellow at 70%%, %v",
						seed, s.name, surface, got, want)
				}
			}
		}
	}
	// The reference the coverage is read from, as it lands: Obsidian's yellow
	// over its own white page and over its own dark one.
	for _, want := range []struct {
		surface, fill stdcolor.NRGBA
	}{
		{stdcolor.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, stdcolor.NRGBA{R: 0xff, G: 0xec, B: 0x99, A: 0xff}},
		{stdcolor.NRGBA{R: 0x1e, G: 0x1e, B: 0x1e, A: 0xff}, stdcolor.NRGBA{R: 0x78, G: 0x65, B: 0x12, A: 0xff}},
	} {
		light, _ := tokens.FromSeed(tokens.DefaultSeed)
		if got := light.HighlightOn(want.surface); got != want.fill {
			t.Errorf("the yellow over %v lands on %v, want %v", want.surface, got, want.fill)
		}
	}
}

// TestTheHighlightPinsItsValuesOnTheDefaultSeed pins the four fills a reader
// of the default palette actually sees — a match and the current match, over
// the content and over the chrome, in each scheme — and the [ColorTokens]
// field, which is the answer over the surface content stands on.
func TestTheHighlightPinsItsValuesOnTheDefaultSeed(t *testing.T) {
	light, dark := tokens.FromSeed(tokens.DefaultSeed)
	for _, want := range []struct {
		scheme      string
		tok         tokens.ColorTokens
		level       tokens.ElevationLevel
		match, curr stdcolor.NRGBA
	}{
		{"light", light, tokens.Level0,
			stdcolor.NRGBA{R: 0xf7, G: 0xe4, B: 0x91, A: 0xff}, stdcolor.NRGBA{R: 0xfb, G: 0xda, B: 0x48, A: 0xff}},
		{"light", light, tokens.LevelChrome,
			stdcolor.NRGBA{R: 0xee, G: 0xdb, B: 0x88, A: 0xff}, stdcolor.NRGBA{R: 0xf7, G: 0xd6, B: 0x44, A: 0xff}},
		{"dark", dark, tokens.Level0,
			stdcolor.NRGBA{R: 0x74, G: 0x62, B: 0x0e, A: 0xff}, stdcolor.NRGBA{R: 0xba, G: 0x99, B: 0x07, A: 0xff}},
		{"dark", dark, tokens.LevelChrome,
			stdcolor.NRGBA{R: 0x73, G: 0x60, B: 0x0d, A: 0xff}, stdcolor.NRGBA{R: 0xb9, G: 0x98, B: 0x06, A: 0xff}},
	} {
		surface := want.tok.SurfaceAt(want.level)
		if got := want.tok.HighlightOn(surface); got != want.match {
			t.Errorf("%s: the highlight on the level-%d surface %v = %v, want %v", want.scheme, want.level, surface, got, want.match)
		}
		if got := want.tok.CurrentMatchOn(surface); got != want.curr {
			t.Errorf("%s: the current match on the level-%d surface %v = %v, want %v", want.scheme, want.level, surface, got, want.curr)
		}
	}
	for _, s := range []struct {
		name string
		tok  tokens.ColorTokens
	}{{"light", light}, {"dark", dark}} {
		if got := s.tok.Highlight; got != s.tok.HighlightOn(s.tok.Background) {
			t.Errorf("%s: Highlight = %v but the answer over the content is %v; the field is the answer over the surface content stands on",
				s.name, got, s.tok.HighlightOn(s.tok.Background))
		}
	}
}

// TestTheCurrentMatchIsTheSameYellowLaidOnMoreStrongly is the gate the second
// answer exists for: the current match differs from a match in the coverage
// and in nothing else, so on every channel it stands nearer the marker yellow
// than the match does, and the step between them is one a reader can see.
func TestTheCurrentMatchIsTheSameYellowLaidOnMoreStrongly(t *testing.T) {
	// The sweep's own worst reading, less a rounding margin: the two fills
	// are always at least this far apart in OKLab.
	const stepBound = 0.05
	worst := 99.0
	for _, seed := range sweepSeeds() {
		for _, s := range schemesOf(seed) {
			for _, lv := range highlightLevels {
				surface := s.tok.SurfaceAt(lv)
				match, curr := s.tok.HighlightOn(surface), s.tok.CurrentMatchOn(surface)
				for _, ch := range []struct {
					name                string
					yellow, match, curr uint8
				}{
					{"red", markerYellow.R, match.R, curr.R},
					{"green", markerYellow.G, match.G, curr.G},
					{"blue", markerYellow.B, match.B, curr.B},
				} {
					near := func(v uint8) int {
						d := int(v) - int(ch.yellow)
						if d < 0 {
							return -d
						}
						return d
					}
					if near(ch.curr) > near(ch.match) {
						t.Errorf("seed %v: %s on level %d, the current match %v stands further from the yellow in %s than the match %v does",
							seed, s.name, lv, curr, ch.name, match)
					}
				}
				if got := oklabDistance(curr, match); got < stepBound {
					t.Errorf("seed %v: %s on level %d, the current match %v is %.4f from the match %v in OKLab — under the %.2f a reader can tell apart",
						seed, s.name, lv, curr, got, match, stepBound)
				} else if got < worst {
					worst = got
				}
			}
		}
	}
	t.Logf("over %d seeds, both derivations, both schemes, five levels: the current match stands %.4f from the match in OKLab at worst",
		len(sweepSeeds()), worst)
}

// TestBothHighlightAnswersAreOpaque holds the token contract: no fill this
// package answers with is translucent. The coverage is spent here, on the
// surface the caller named, and what the caller paints is a colour.
func TestBothHighlightAnswersAreOpaque(t *testing.T) {
	for _, seed := range sweepSeeds() {
		for _, s := range schemesOf(seed) {
			for _, lv := range highlightLevels {
				surface := s.tok.SurfaceAt(lv)
				for _, fill := range []stdcolor.NRGBA{s.tok.HighlightOn(surface), s.tok.CurrentMatchOn(surface), s.tok.Highlight} {
					if fill.A != 0xff {
						t.Fatalf("seed %v: %s on level %d answered %v, whose alpha is not opaque", seed, s.name, lv, fill)
					}
				}
			}
		}
	}
}

// TestTheHighlightSeparatesFromTheSurfaceItMarks verifies a highlight is
// visible as a field wherever content is marked, on every level in both
// schemes of both derivations, for every seed.
//
// The gate is a perceptual distance and not [tokens.ContainerFloor], because
// what separates this fill from its surface is mostly chroma and a contrast
// ratio reads luminance alone: the reference's own yellow over its own white
// page measures 1.186:1, and it is a highlight anyone can see.
func TestTheHighlightSeparatesFromTheSurfaceItMarks(t *testing.T) {
	// The sweep's own worst reading, less a rounding margin.
	const seamBound = 0.10
	worst, best := 99.0, 0.0
	for _, seed := range sweepSeeds() {
		for _, s := range schemesOf(seed) {
			for _, lv := range highlightLevels {
				surface := s.tok.SurfaceAt(lv)
				got := oklabDistance(s.tok.HighlightOn(surface), surface)
				if got < seamBound {
					t.Errorf("seed %v: %s highlight %v on the level-%d fill %v is %.4f away in OKLab — under the %.2f a field is findable by",
						seed, s.name, s.tok.HighlightOn(surface), lv, surface, got, seamBound)
				} else if got < worst {
					worst = got
				}
				if got > best {
					best = got
				}
			}
		}
	}
	t.Logf("over %d seeds, both derivations, both schemes, five levels: the highlight stands %.4f from the surface it marks at worst, %.4f at best",
		len(sweepSeeds()), worst, best)
}

// TestTheMarkedTextIsMeasuredNotCorrected records what a highlight costs the
// words it covers. The text is not repainted and the coverage is Obsidian's,
// not a number tuned to a floor, so this is a measurement with one gate under
// it: where content actually stands — the content plane and the window's
// chrome — the scheme's body foreground still reaches highlightReach over a
// match, a whisker under [tokens.TextFloor] in the dark scheme. Over the
// current match, and over a match on the raised levels, it falls further,
// and the readings are logged rather than corrected for.
func TestTheMarkedTextIsMeasuredNotCorrected(t *testing.T) {
	worst := 999.0
	for _, seed := range sweepSeeds() {
		for _, s := range schemesOf(seed) {
			for _, lv := range []tokens.ElevationLevel{tokens.LevelChrome, tokens.Level0} {
				fill := s.tok.HighlightOn(s.tok.SurfaceAt(lv))
				got := color.Magnitude(s.tok.Text, fill)
				if got < highlightReach {
					t.Errorf("seed %v: %s Text %v over the level-%d highlight %v measures |Lc| %.2f, under the %.0f the marker reaches",
						seed, s.name, s.tok.Text, lv, fill, got, highlightReach)
				} else if got < worst {
					worst = got
				}
			}
		}
	}
	t.Logf("over %d seeds, both derivations, both schemes: the body foreground over a match on the content or the chrome measures |Lc| %.2f at worst (floor %.0f)",
		len(sweepSeeds()), worst, highlightReach)
	light, dark := tokens.FromSeed(tokens.DefaultSeed)
	for _, s := range []struct {
		name string
		tok  tokens.ColorTokens
	}{{"light", light}, {"dark", dark}} {
		for _, lv := range highlightLevels {
			surface := s.tok.SurfaceAt(lv)
			t.Logf("default seed, %s, level %d on %v: text over the match |Lc| %.2f, over the current match |Lc| %.2f",
				s.name, lv, surface,
				color.Magnitude(s.tok.Text, s.tok.HighlightOn(surface)),
				color.Magnitude(s.tok.Text, s.tok.CurrentMatchOn(surface)))
		}
	}
}

// TestTheHighlightKeepsItsDistanceFromEveryStatus is the gate the reserved
// yellow exists for: a highlight reports no status, so no status may be read
// off it. It measures both answers against every status colour a reader could
// see beside them — each role's fixed container, the container resolved for
// the level the fill is on, and the role's pinned base — in OKLCh hue and in
// OKLab distance, over the whole seed sweep in both schemes of both
// derivations.
//
// The bounds are the sweep's own measurements less a rounding margin. The
// closest approach in hue is a dark scheme's current match against the
// warning's pin, the dark surface under the yellow carrying it that way; in
// OKLab, which reads the whole colour rather than its angle, nothing comes
// nearer than 0.0665, against the 0.0286 the two closest status containers
// come to each other.
func TestTheHighlightKeepsItsDistanceFromEveryStatus(t *testing.T) {
	const (
		hueBound = 26.0  // measured 26.81° over the sweep
		labBound = 0.066 // measured 0.0665 over the sweep
	)
	worstHue, worstHueAt := 999.0, ""
	worstLab, worstLabAt := 99.0, ""
	for _, seed := range sweepSeeds() {
		for _, s := range schemesOf(seed) {
			for _, lv := range highlightLevels {
				surface := s.tok.SurfaceAt(lv)
				for _, fill := range []stdcolor.NRGBA{s.tok.HighlightOn(surface), s.tok.CurrentMatchOn(surface)} {
					_, _, fillHue := color.OKLChFromNRGBA(fill)
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
							if got := hueGap(fillHue, hue); got < worstHue {
								worstHue = got
								worstHueAt = a.what
								if got < hueBound {
									t.Errorf("seed %v: %s highlight %v on level %d is %.2f° from the %s %s %v — under the %.1f° the reservation holds",
										seed, s.name, fill, lv, got, r.name, a.what, a.c, hueBound)
								}
							}
							if got := oklabDistance(fill, a.c); got < worstLab {
								worstLab = got
								worstLabAt = a.what
								if got < labBound {
									t.Errorf("seed %v: %s highlight %v on level %d is %.4f from the %s %s %v in OKLab — under the %.3f the reservation holds",
										seed, s.name, fill, lv, got, r.name, a.what, a.c, labBound)
								}
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
