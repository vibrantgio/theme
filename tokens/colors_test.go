package tokens_test

import (
	"image/color"
	"math"
	"math/rand"
	"testing"

	"github.com/vibrantgio/theme/tokens"
)

// contrastPair is a named foreground/background pair for WCAG AA verification.
type contrastPair struct {
	name string
	bg   color.NRGBA
	fg   color.NRGBA
}

// tokenPairs returns the foreground/background pairs every scheme must
// carry at WCAG AA. The pinned roles pair with their "On" colours; the
// surface pairs are named by the ramp steps they resolve from.
func tokenPairs(t tokens.ColorTokens) []contrastPair {
	n := t.Ramps.Neutral
	return []contrastPair{
		{"Background/Text", t.Background, t.Text},
		{"InverseSurface/OnInverseSurface", t.InverseSurface, t.OnInverseSurface},
		{"Surface/Neutral.Step(900)", t.Surface, n.Step(900)},
		{"Neutral.Step(300)/Neutral.Step(700)", n.Step(300), n.Step(700)},
		{"Primary/OnPrimary", t.Primary, t.OnPrimary},
		{"Secondary/OnSecondary", t.Secondary, t.OnSecondary},
		{"Tertiary/OnTertiary", t.Tertiary, t.OnTertiary},
		{"Error/OnError", t.Error, t.OnError},
		{"Success/OnSuccess", t.Success, t.OnSuccess},
		{"Warning/OnWarning", t.Warning, t.OnWarning},
		{"Info/OnInfo", t.Info, t.OnInfo},
	}
}

func linearize(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func relativeLuminance(c color.NRGBA) float64 {
	r := linearize(float64(c.R) / 255.0)
	g := linearize(float64(c.G) / 255.0)
	b := linearize(float64(c.B) / 255.0)
	return 0.2126*r + 0.7152*g + 0.0722*b
}

func contrastRatio(c1, c2 color.NRGBA) float64 {
	l1 := relativeLuminance(c1)
	l2 := relativeLuminance(c2)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

const wcagAA = 4.5

// sweepSeeds is the seed sweep the derivation's whole-population properties
// are asserted over: fourteen chosen colours — the default seed, the nine
// macOS system accents, both ends of the tonal axis and three pastels — and
// four hundred random ones from a fixed source, so the population is wide
// and the run is the same one every time. It is deliberately shared: a
// property that has to hold for every seed should be read against the same
// seeds as every other, or two gates disagree about what "every seed"
// meant.
//
// The three pastels are named rather than left to the random draw because
// of the shape they have, not the colours they are. A palette published for
// a dark scheme states its accents at a high tone — around L* 73 to L* 85 —
// and people seed a brand with one, which is a seed that reads perfectly in
// the dark scheme and lands the light scheme's primary pin a whisper off
// the paper. That shape is what put a 1.95:1 link on a light page (see
// ink.go); the random draw covers it thinly and by accident, and a
// regression that only a randomly drawn seed catches is one a future change
// to the draw can lose. So the shape is in the matrix by name: a blue, a
// mauve and a green at L* 72.8, 74.0 and 84.8.
func sweepSeeds() []color.NRGBA {
	rng := rand.New(rand.NewSource(20260818))
	seeds := []color.NRGBA{
		tokens.DefaultSeed,
		{0xff, 0x3b, 0x30, 0xff}, {0xff, 0x95, 0x00, 0xff}, {0xff, 0xcc, 0x00, 0xff},
		{0x28, 0xcd, 0x41, 0xff}, {0x00, 0x7a, 0xff, 0xff}, {0xaf, 0x52, 0xde, 0xff},
		{0xff, 0x2d, 0x55, 0xff}, {0x8e, 0x8e, 0x93, 0xff}, {0x00, 0x00, 0x00, 0xff},
		{0xff, 0xff, 0xff, 0xff},
		{0x89, 0xb4, 0xfa, 0xff}, {0xcb, 0xa6, 0xf7, 0xff}, {0xa6, 0xe3, 0xa1, 0xff},
	}
	for i := 0; i < 400; i++ {
		seeds = append(seeds, color.NRGBA{
			uint8(rng.Intn(256)), uint8(rng.Intn(256)), uint8(rng.Intn(256)), 0xff})
	}
	return seeds
}

// TestRampStepAddressing verifies Step's 100–900 addressing over the
// backing array and that out-of-vocabulary steps panic.
func TestRampStepAddressing(t *testing.T) {
	var r tokens.Ramp
	for i := range r {
		r[i] = color.NRGBA{R: uint8(i + 1), A: 0xff}
	}
	for n := 100; n <= 900; n += 100 {
		want := color.NRGBA{R: uint8(n / 100), A: 0xff}
		if got := r.Step(n); got != want {
			t.Errorf("Step(%d) = %v, want %v", n, got, want)
		}
	}
	for _, bad := range []int{0, 50, 99, 150, 901, 1000, -100} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("Step(%d): expected panic", bad)
				}
			}()
			r.Step(bad)
		}()
	}
}

// TestSemanticLayerResolvesFromRamps verifies the semantic layer resolves
// from ramp steps in both default schemes: Surface, Divider and the inverse
// pair. The inverse pair is the one resolution that reads across the
// scheme boundary: it is the counterpart scheme's own Surface and Text, so
// each scheme's inverse chip is built out of the other one.
func TestSemanticLayerResolvesFromRamps(t *testing.T) {
	for _, s := range []struct {
		name        string
		tok         tokens.ColorTokens
		counterpart tokens.ColorTokens
	}{
		{"DefaultLight", tokens.DefaultLight, tokens.DefaultDark},
		{"DefaultDark", tokens.DefaultDark, tokens.DefaultLight},
	} {
		n := s.tok.Ramps.Neutral
		checks := []struct {
			name string
			got  color.NRGBA
			want color.NRGBA
		}{
			{"Surface = Neutral.Step(200)", s.tok.Surface, n.Step(200)},
			{"Divider = Neutral.Step(300)", s.tok.Divider, n.Step(300)},
			{"InverseSurface = the counterpart scheme's Surface", s.tok.InverseSurface, s.counterpart.Surface},
			{"OnInverseSurface = the counterpart scheme's Text", s.tok.OnInverseSurface, s.counterpart.Text},
		}
		for _, c := range checks {
			if c.got != c.want {
				t.Errorf("%s: %s: got %v, want %v", s.name, c.name, c.got, c.want)
			}
		}
	}
}

// TestAllRampStepsPopulated verifies every step of every role ramp, every
// pin and every semantic colour is an opaque colour in both default schemes.
func TestAllRampStepsPopulated(t *testing.T) {
	for _, s := range []struct {
		name string
		tok  tokens.ColorTokens
	}{
		{"DefaultLight", tokens.DefaultLight},
		{"DefaultDark", tokens.DefaultDark},
	} {
		for _, r := range namedRamps(s.tok) {
			for step := 100; step <= 900; step += 100 {
				if c := r.ramp.Step(step); c.A != 0xff {
					t.Errorf("%s: %s.Step(%d) = %v, want an opaque colour", s.name, r.name, step, c)
				}
			}
		}
		for name, c := range map[string]color.NRGBA{
			"Primary": s.tok.Primary, "OnPrimary": s.tok.OnPrimary,
			"Secondary": s.tok.Secondary, "OnSecondary": s.tok.OnSecondary,
			"Tertiary": s.tok.Tertiary, "OnTertiary": s.tok.OnTertiary,
			"Error": s.tok.Error, "OnError": s.tok.OnError,
			"Success": s.tok.Success, "OnSuccess": s.tok.OnSuccess,
			"Warning": s.tok.Warning, "OnWarning": s.tok.OnWarning,
			"Background": s.tok.Background, "Text": s.tok.Text,
			"InverseSurface": s.tok.InverseSurface, "OnInverseSurface": s.tok.OnInverseSurface,
		} {
			if c.A != 0xff {
				t.Errorf("%s: %s = %v, want an opaque colour", s.name, name, c)
			}
		}
	}
}

type namedRamp struct {
	name string
	ramp tokens.Ramp
}

func namedRamps(t tokens.ColorTokens) []namedRamp {
	return []namedRamp{
		{"Neutral", t.Ramps.Neutral},
		{"Primary", t.Ramps.Primary},
		{"Secondary", t.Ramps.Secondary},
		{"Tertiary", t.Ramps.Tertiary},
		{"Error", t.Ramps.Error},
		{"Success", t.Ramps.Success},
		{"Warning", t.Ramps.Warning},
		{"Info", t.Ramps.Info},
	}
}

// TestRampLightnessDirection verifies the paired scales run the documented
// way in every ramp: in light mode step 100 is the lightest and luminance
// strictly falls toward 900; in dark mode step 100 is the darkest ground
// and luminance strictly rises toward 900 — same step, same job.
func TestRampLightnessDirection(t *testing.T) {
	for _, s := range []struct {
		name       string
		tok        tokens.ColorTokens
		descending bool
	}{
		{"DefaultLight", tokens.DefaultLight, true},
		{"DefaultDark", tokens.DefaultDark, false},
	} {
		for _, r := range namedRamps(s.tok) {
			for step := 200; step <= 900; step += 100 {
				prev := relativeLuminance(r.ramp.Step(step - 100))
				cur := relativeLuminance(r.ramp.Step(step))
				if s.descending && cur >= prev {
					t.Errorf("%s: %s.Step(%d) luminance %.4f not below step %d's %.4f",
						s.name, r.name, step, cur, step-100, prev)
				}
				if !s.descending && cur <= prev {
					t.Errorf("%s: %s.Step(%d) luminance %.4f not above step %d's %.4f",
						s.name, r.name, step, cur, step-100, prev)
				}
			}
		}
	}
}

// TestFromSeedPinsTheLiftedSeed verifies the light primary base is the seed
// at its own hue and CIELAB depth with the accent dial on its chroma —
// never read off a ramp step — with the alpha forced opaque. A seed the
// dial leaves alone (already past it, or grey) comes back byte-for-byte,
// which is what keeps a palette seeded from a desktop accent colour
// matching that accent.
func TestFromSeedPinsTheLiftedSeed(t *testing.T) {
	for _, tc := range []struct {
		name string
		seed color.NRGBA
		want color.NRGBA
	}{
		{"the default seed, lifted", color.NRGBA{0x67, 0x50, 0xA4, 0xff}, color.NRGBA{0x72, 0x3a, 0xd4, 0xff}},
		{"already past the dial", color.NRGBA{0xff, 0x3b, 0x30, 0xff}, color.NRGBA{0xff, 0x3b, 0x30, 0xff}},
		{"grey stays grey", color.NRGBA{0x8e, 0x8e, 0x93, 0xff}, color.NRGBA{0x8e, 0x8e, 0x93, 0xff}},
		{"alpha forced opaque", color.NRGBA{0xff, 0x3b, 0x30, 0x80}, color.NRGBA{0xff, 0x3b, 0x30, 0xff}},
	} {
		light, _ := tokens.FromSeed(tc.seed)
		if light.Primary != tc.want {
			t.Errorf("%s: FromSeed(%v): light Primary = %v, want %v",
				tc.name, tc.seed, light.Primary, tc.want)
		}
	}
	if p := tokens.DefaultLight.Primary; p != (color.NRGBA{0x72, 0x3a, 0xd4, 0xff}) {
		t.Errorf("DefaultLight.Primary = %v, want the lifted default seed #723AD4", p)
	}
}

// TestFromSeedReproducesItselfFromItsBase verifies the derivation is a
// projection: deriving from a pair's own light primary base gives that pair
// back byte-for-byte, in the default derivation and the high-contrast one.
// It is the property a serialized theme rests on — a theme names one colour
// and is rebuilt from it — so it is asserted over a wide sweep of seeds
// rather than the handful the rest of the package uses.
func TestFromSeedReproducesItselfFromItsBase(t *testing.T) {
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		if l, d := tokens.FromSeed(light.Primary); l != light || d != dark {
			t.Fatalf("FromSeed(%v): re-deriving from the light Primary base %v does not reproduce the pair",
				seed, light.Primary)
		}
		hcLight, hcDark := tokens.FromSeedHighContrast(seed)
		if l, d := tokens.FromSeedHighContrast(light.Primary); l != hcLight || d != hcDark {
			t.Fatalf("FromSeedHighContrast(%v): re-deriving from the light Primary base %v does not reproduce the pair",
				seed, light.Primary)
		}
	}
}

// defaultGolden returns the recorded palette FromSeed derives from
// DefaultSeed (#6750A4): every ramp step, pin and semantic colour, byte for
// byte. It is the regression anchor for the whole colour engine — any
// change to the derivation (scales, chromas, hues, pins) must show up here
// as an explicit, reviewed palette change. Landmarks that tie the recording
// to the measured evidence: the dark Neutral steps 100–400, 700 and 900
// hold the measured dark column at its measured depths with the seed's tint
// taken out (#181818, #222222, #2e2e2e, #474747, #cccccc, #eeeeee — the
// source read #18171c, #222126, #2e2e33, #47464c, #cccbd2, #eeedf4, its
// neutrals carrying the brand). Steps 500, 600 and 800 are the tone curve's,
// not the column's: the source has one grey between its deepest surface and
// its pin (#9e9da4, L* 65) and the ramp needs three. Two further recorded
// values sit off the raw measurements: the light 900 stops are at L* 6
// rather than the measured L* 18 (with Text following the neutral 900) to
// clear the Lc ≥ 90 gate over the step-200 grounds, and the dark pins are
// at L* 82 (step-700 depth, byte-identical to the ramp's 700 stop) rather
// than L* 65, because no on-colour reaches Lc 60 over an L* 65 mid-tone.
func defaultGolden() (light, dark tokens.ColorTokens) {
	hex := func(r, g, b uint8) color.NRGBA { return color.NRGBA{r, g, b, 0xff} }
	light = tokens.ColorTokens{
		Ramps: tokens.RampSet{
			Neutral: tokens.Ramp{
				hex(0xf6, 0xf6, 0xf6),
				hex(0xe8, 0xe8, 0xe8),
				hex(0xd4, 0xd4, 0xd4),
				hex(0xb6, 0xb6, 0xb6),
				hex(0x98, 0x98, 0x98),
				hex(0x79, 0x79, 0x79),
				hex(0x5c, 0x5c, 0x5c),
				hex(0x42, 0x42, 0x42),
				hex(0x13, 0x13, 0x13),
			},
			Primary: tokens.Ramp{
				hex(0xf7, 0xf5, 0xff),
				hex(0xea, 0xe5, 0xff),
				hex(0xd8, 0xce, 0xff),
				hex(0xbd, 0xaa, 0xff),
				hex(0xa5, 0x84, 0xff),
				hex(0x8c, 0x59, 0xf4),
				hex(0x6f, 0x36, 0xd1),
				hex(0x57, 0x0a, 0xb2),
				hex(0x1d, 0x00, 0x44),
			},
			Secondary: tokens.Ramp{
				hex(0xf7, 0xf5, 0xff),
				hex(0xea, 0xe5, 0xff),
				hex(0xd8, 0xce, 0xff),
				hex(0xba, 0xaf, 0xe3),
				hex(0x9c, 0x92, 0xc4),
				hex(0x7d, 0x73, 0xa3),
				hex(0x60, 0x55, 0x83),
				hex(0x46, 0x3b, 0x67),
				hex(0x18, 0x0b, 0x32),
			},
			Tertiary: tokens.Ramp{
				hex(0xff, 0xf4, 0xf8),
				hex(0xff, 0xe0, 0xeb),
				hex(0xff, 0xc4, 0xdb),
				hex(0xf2, 0x9d, 0xc0),
				hex(0xd2, 0x7f, 0xa2),
				hex(0xb0, 0x60, 0x83),
				hex(0x8e, 0x42, 0x65),
				hex(0x70, 0x28, 0x4a),
				hex(0x2f, 0x00, 0x19),
			},
			Error: tokens.Ramp{
				hex(0xff, 0xf4, 0xf3),
				hex(0xff, 0xe2, 0xde),
				hex(0xff, 0xc7, 0xc1),
				hex(0xff, 0x9b, 0x92),
				hex(0xfa, 0x6a, 0x63),
				hex(0xd5, 0x48, 0x44),
				hex(0xb1, 0x22, 0x26),
				hex(0x8a, 0x00, 0x10),
				hex(0x32, 0x00, 0x02),
			},
			Success: tokens.Ramp{
				hex(0xde, 0xff, 0xe0),
				hex(0x96, 0xff, 0xa4),
				hex(0x7f, 0xeb, 0x8f),
				hex(0x5f, 0xcc, 0x71),
				hex(0x3e, 0xad, 0x54),
				hex(0x0c, 0x8d, 0x35),
				hex(0x00, 0x6b, 0x24),
				hex(0x00, 0x4e, 0x18),
				hex(0x00, 0x18, 0x03),
			},
			Warning: tokens.Ramp{
				hex(0xff, 0xf5, 0xe4),
				hex(0xff, 0xe5, 0xb7),
				hex(0xff, 0xcd, 0x6e),
				hex(0xff, 0xa0, 0x31),
				hex(0xed, 0x78, 0x19),
				hex(0xc1, 0x5d, 0x00),
				hex(0x94, 0x46, 0x00),
				hex(0x6d, 0x31, 0x00),
				hex(0x25, 0x0c, 0x00),
			},
			Info: tokens.Ramp{
				hex(0xf0, 0xf7, 0xff),
				hex(0xd8, 0xea, 0xff),
				hex(0xb5, 0xd8, 0xff),
				hex(0x7c, 0xba, 0xff),
				hex(0x3a, 0x9c, 0xfd),
				hex(0x09, 0x7b, 0xd9),
				hex(0x00, 0x5d, 0xa8),
				hex(0x00, 0x43, 0x7c),
				hex(0x00, 0x14, 0x2c),
			},
		},
		Primary:     hex(0x72, 0x3a, 0xd4),
		OnPrimary:   hex(0xff, 0xff, 0xff),
		Secondary:   hex(0x62, 0x58, 0x86),
		OnSecondary: hex(0xff, 0xff, 0xff),
		Tertiary:    hex(0x91, 0x45, 0x67),
		OnTertiary:  hex(0xff, 0xff, 0xff),
		Error:       hex(0xb1, 0x22, 0x26),
		OnError:     hex(0xff, 0xff, 0xff),
		Success:     hex(0x00, 0x6b, 0x24),
		OnSuccess:   hex(0xff, 0xff, 0xff),
		Warning:     hex(0x94, 0x46, 0x00),
		OnWarning:   hex(0xff, 0xff, 0xff),
		Info:        hex(0x00, 0x5d, 0xa8),
		OnInfo:      hex(0xff, 0xff, 0xff),
		Background:  hex(0xf6, 0xf6, 0xf6),
		Text:        hex(0x13, 0x13, 0x13),
	}
	dark = tokens.ColorTokens{
		Ramps: tokens.RampSet{
			Neutral: tokens.Ramp{
				hex(0x18, 0x18, 0x18),
				hex(0x22, 0x22, 0x22),
				hex(0x2e, 0x2e, 0x2e),
				hex(0x47, 0x47, 0x47),
				hex(0x6d, 0x6d, 0x6d),
				hex(0x9b, 0x9b, 0x9b),
				hex(0xcc, 0xcc, 0xcc),
				hex(0xd7, 0xd7, 0xd7),
				hex(0xee, 0xee, 0xee),
			},
			Primary: tokens.Ramp{
				hex(0x22, 0x00, 0x4e),
				hex(0x2f, 0x00, 0x66),
				hex(0x3f, 0x00, 0x85),
				hex(0x5c, 0x15, 0xb7),
				hex(0x80, 0x4b, 0xe5),
				hex(0xa7, 0x87, 0xff),
				hex(0xd0, 0xc4, 0xff),
				hex(0xda, 0xd2, 0xff),
				hex(0xef, 0xec, 0xff),
			},
			Secondary: tokens.Ramp{
				hex(0x1c, 0x10, 0x37),
				hex(0x26, 0x1a, 0x43),
				hex(0x32, 0x27, 0x51),
				hex(0x4a, 0x40, 0x6c),
				hex(0x71, 0x66, 0x96),
				hex(0x9f, 0x95, 0xc7),
				hex(0xd0, 0xc5, 0xfa),
				hex(0xda, 0xd2, 0xff),
				hex(0xef, 0xec, 0xff),
			},
			Tertiary: tokens.Ramp{
				hex(0x37, 0x00, 0x1d),
				hex(0x49, 0x00, 0x29),
				hex(0x59, 0x10, 0x36),
				hex(0x76, 0x2c, 0x4f),
				hex(0xa2, 0x54, 0x76),
				hex(0xd5, 0x82, 0xa5),
				hex(0xff, 0xb8, 0xd4),
				hex(0xff, 0xc8, 0xdd),
				hex(0xff, 0xe8, 0xf0),
			},
			Error: tokens.Ramp{
				hex(0x3a, 0x00, 0x03),
				hex(0x4d, 0x00, 0x05),
				hex(0x64, 0x00, 0x08),
				hex(0x93, 0x00, 0x12),
				hex(0xc6, 0x39, 0x38),
				hex(0xfd, 0x6d, 0x65),
				hex(0xff, 0xbb, 0xb4),
				hex(0xff, 0xcb, 0xc5),
				hex(0xff, 0xe9, 0xe6),
			},
			Success: tokens.Ramp{
				hex(0x00, 0x1d, 0x05),
				hex(0x00, 0x28, 0x08),
				hex(0x00, 0x37, 0x0e),
				hex(0x00, 0x53, 0x1a),
				hex(0x00, 0x7f, 0x2c),
				hex(0x41, 0xb0, 0x57),
				hex(0x77, 0xe3, 0x87),
				hex(0x82, 0xee, 0x92),
				hex(0xb7, 0xff, 0xbe),
			},
			Warning: tokens.Ramp{
				hex(0x2c, 0x10, 0x00),
				hex(0x3b, 0x18, 0x00),
				hex(0x4e, 0x22, 0x00),
				hex(0x74, 0x35, 0x00),
				hex(0xae, 0x53, 0x00),
				hex(0xf0, 0x7b, 0x1e),
				hex(0xff, 0xc2, 0x42),
				hex(0xff, 0xd0, 0x79),
				hex(0xff, 0xec, 0xc9),
			},
			Info: tokens.Ramp{
				hex(0x00, 0x18, 0x33),
				hex(0x00, 0x22, 0x44),
				hex(0x00, 0x2f, 0x5a),
				hex(0x00, 0x48, 0x84),
				hex(0x00, 0x6e, 0xc6),
				hex(0x3e, 0x9e, 0xff),
				hex(0xa6, 0xd0, 0xff),
				hex(0xba, 0xda, 0xff),
				hex(0xe2, 0xef, 0xff),
			},
		},
		Primary:     hex(0xd0, 0xc4, 0xff),
		OnPrimary:   hex(0x22, 0x00, 0x4e),
		Secondary:   hex(0xd0, 0xc5, 0xfa),
		OnSecondary: hex(0x1c, 0x10, 0x37),
		Tertiary:    hex(0xff, 0xb8, 0xd4),
		OnTertiary:  hex(0x37, 0x00, 0x1d),
		Error:       hex(0xff, 0xbb, 0xb4),
		OnError:     hex(0x3a, 0x00, 0x03),
		Success:     hex(0x77, 0xe3, 0x87),
		OnSuccess:   hex(0x00, 0x1d, 0x05),
		Warning:     hex(0xff, 0xc2, 0x42),
		OnWarning:   hex(0x2c, 0x10, 0x00),
		Info:        hex(0xa6, 0xd0, 0xff),
		OnInfo:      hex(0x00, 0x18, 0x33),
		Background:  hex(0x18, 0x18, 0x18),
		Text:        hex(0xee, 0xee, 0xee),
	}
	// Surface, Divider and the inverse pair are recorded as the resolutions
	// they are — the first two off this scheme's neutral ramp, the inverse
	// pair off the counterpart scheme's, which is what makes a light
	// scheme's inverse chip dark and a dark scheme's light.
	fill := func(t, counterpart tokens.ColorTokens) tokens.ColorTokens {
		n, o := t.Ramps.Neutral, counterpart.Ramps.Neutral
		t.Surface = n.Step(200)
		t.Divider = n.Step(300)
		t.InverseSurface = o.Step(200)
		t.OnInverseSurface = o.Step(900)
		return t
	}
	return fill(light, dark), fill(dark, light)
}

// hcGolden returns the recorded palette FromSeedHighContrast derives from
// DefaultSeed: the high-contrast variant, byte for byte. Landmarks
// that tie the recording to the variant's three widenings: steps 100–600 of
// every ramp are byte-identical to defaultGolden's (the grounds stay), the
// 700 stops sit at the default scale's 900 depth (light #131313 neutral 700
// = defaultGolden's neutral 900; dark #eeeeee mirrors it) with 800/900
// sliding to the axis ends (light 900 pure black, dark 900 pure white, Text
// following in both modes), Divider is the step-500 strong border rather
// than 300, and the dark pins keep their L* 82 bases while their on-colours
// drop to pure black (tone 0). The light pins are byte-identical to
// defaultGolden's — White already clears the raised Lc ≥ 75 floor — so the
// HC light Primary is the same lifted seed.
func hcGolden() (light, dark tokens.ColorTokens) {
	hex := func(r, g, b uint8) color.NRGBA { return color.NRGBA{r, g, b, 0xff} }
	light = tokens.ColorTokens{
		Ramps: tokens.RampSet{
			Neutral: tokens.Ramp{
				hex(0xf6, 0xf6, 0xf6),
				hex(0xe8, 0xe8, 0xe8),
				hex(0xd4, 0xd4, 0xd4),
				hex(0xb6, 0xb6, 0xb6),
				hex(0x98, 0x98, 0x98),
				hex(0x79, 0x79, 0x79),
				hex(0x13, 0x13, 0x13),
				hex(0x0b, 0x0b, 0x0b),
				hex(0x00, 0x00, 0x00),
			},
			Primary: tokens.Ramp{
				hex(0xf7, 0xf5, 0xff),
				hex(0xea, 0xe5, 0xff),
				hex(0xd8, 0xce, 0xff),
				hex(0xbd, 0xaa, 0xff),
				hex(0xa5, 0x84, 0xff),
				hex(0x8c, 0x59, 0xf4),
				hex(0x1d, 0x00, 0x44),
				hex(0x12, 0x00, 0x30),
				hex(0x00, 0x00, 0x00),
			},
			Secondary: tokens.Ramp{
				hex(0xf7, 0xf5, 0xff),
				hex(0xea, 0xe5, 0xff),
				hex(0xd8, 0xce, 0xff),
				hex(0xba, 0xaf, 0xe3),
				hex(0x9c, 0x92, 0xc4),
				hex(0x7d, 0x73, 0xa3),
				hex(0x18, 0x0b, 0x32),
				hex(0x10, 0x03, 0x28),
				hex(0x00, 0x00, 0x00),
			},
			Tertiary: tokens.Ramp{
				hex(0xff, 0xf4, 0xf8),
				hex(0xff, 0xe0, 0xeb),
				hex(0xff, 0xc4, 0xdb),
				hex(0xf2, 0x9d, 0xc0),
				hex(0xd2, 0x7f, 0xa2),
				hex(0xb0, 0x60, 0x83),
				hex(0x2f, 0x00, 0x19),
				hex(0x20, 0x00, 0x0f),
				hex(0x00, 0x00, 0x00),
			},
			Error: tokens.Ramp{
				hex(0xff, 0xf4, 0xf3),
				hex(0xff, 0xe2, 0xde),
				hex(0xff, 0xc7, 0xc1),
				hex(0xff, 0x9b, 0x92),
				hex(0xfa, 0x6a, 0x63),
				hex(0xd5, 0x48, 0x44),
				hex(0x32, 0x00, 0x02),
				hex(0x22, 0x00, 0x01),
				hex(0x00, 0x00, 0x00),
			},
			Success: tokens.Ramp{
				hex(0xde, 0xff, 0xe0),
				hex(0x96, 0xff, 0xa4),
				hex(0x7f, 0xeb, 0x8f),
				hex(0x5f, 0xcc, 0x71),
				hex(0x3e, 0xad, 0x54),
				hex(0x0c, 0x8d, 0x35),
				hex(0x00, 0x18, 0x03),
				hex(0x00, 0x0f, 0x02),
				hex(0x00, 0x00, 0x00),
			},
			Warning: tokens.Ramp{
				hex(0xff, 0xf5, 0xe4),
				hex(0xff, 0xe5, 0xb7),
				hex(0xff, 0xcd, 0x6e),
				hex(0xff, 0xa0, 0x31),
				hex(0xed, 0x78, 0x19),
				hex(0xc1, 0x5d, 0x00),
				hex(0x25, 0x0c, 0x00),
				hex(0x19, 0x06, 0x00),
				hex(0x00, 0x00, 0x00),
			},
			Info: tokens.Ramp{
				hex(0xf0, 0xf7, 0xff),
				hex(0xd8, 0xea, 0xff),
				hex(0xb5, 0xd8, 0xff),
				hex(0x7c, 0xba, 0xff),
				hex(0x3a, 0x9c, 0xfd),
				hex(0x09, 0x7b, 0xd9),
				hex(0x00, 0x14, 0x2c),
				hex(0x00, 0x0b, 0x1e),
				hex(0x00, 0x00, 0x00),
			},
		},
		Primary:     hex(0x72, 0x3a, 0xd4),
		OnPrimary:   hex(0xff, 0xff, 0xff),
		Secondary:   hex(0x62, 0x58, 0x86),
		OnSecondary: hex(0xff, 0xff, 0xff),
		Tertiary:    hex(0x91, 0x45, 0x67),
		OnTertiary:  hex(0xff, 0xff, 0xff),
		Error:       hex(0xb1, 0x22, 0x26),
		OnError:     hex(0xff, 0xff, 0xff),
		Success:     hex(0x00, 0x6b, 0x24),
		OnSuccess:   hex(0xff, 0xff, 0xff),
		Warning:     hex(0x94, 0x46, 0x00),
		OnWarning:   hex(0xff, 0xff, 0xff),
		Info:        hex(0x00, 0x5d, 0xa8),
		OnInfo:      hex(0xff, 0xff, 0xff),
		Background:  hex(0xf6, 0xf6, 0xf6),
		Text:        hex(0x00, 0x00, 0x00),
	}
	dark = tokens.ColorTokens{
		Ramps: tokens.RampSet{
			Neutral: tokens.Ramp{
				hex(0x18, 0x18, 0x18),
				hex(0x22, 0x22, 0x22),
				hex(0x2e, 0x2e, 0x2e),
				hex(0x47, 0x47, 0x47),
				hex(0x6d, 0x6d, 0x6d),
				hex(0x9b, 0x9b, 0x9b),
				hex(0xee, 0xee, 0xee),
				hex(0xf6, 0xf6, 0xf6),
				hex(0xff, 0xff, 0xff),
			},
			Primary: tokens.Ramp{
				hex(0x22, 0x00, 0x4e),
				hex(0x2f, 0x00, 0x66),
				hex(0x3f, 0x00, 0x85),
				hex(0x5c, 0x15, 0xb7),
				hex(0x80, 0x4b, 0xe5),
				hex(0xa7, 0x87, 0xff),
				hex(0xef, 0xec, 0xff),
				hex(0xf7, 0xf5, 0xff),
				hex(0xff, 0xff, 0xff),
			},
			Secondary: tokens.Ramp{
				hex(0x1c, 0x10, 0x37),
				hex(0x26, 0x1a, 0x43),
				hex(0x32, 0x27, 0x51),
				hex(0x4a, 0x40, 0x6c),
				hex(0x71, 0x66, 0x96),
				hex(0x9f, 0x95, 0xc7),
				hex(0xef, 0xec, 0xff),
				hex(0xf7, 0xf5, 0xff),
				hex(0xff, 0xff, 0xff),
			},
			Tertiary: tokens.Ramp{
				hex(0x37, 0x00, 0x1d),
				hex(0x49, 0x00, 0x29),
				hex(0x59, 0x10, 0x36),
				hex(0x76, 0x2c, 0x4f),
				hex(0xa2, 0x54, 0x76),
				hex(0xd5, 0x82, 0xa5),
				hex(0xff, 0xe8, 0xf0),
				hex(0xff, 0xf4, 0xf8),
				hex(0xff, 0xff, 0xff),
			},
			Error: tokens.Ramp{
				hex(0x3a, 0x00, 0x03),
				hex(0x4d, 0x00, 0x05),
				hex(0x64, 0x00, 0x08),
				hex(0x93, 0x00, 0x12),
				hex(0xc6, 0x39, 0x38),
				hex(0xfd, 0x6d, 0x65),
				hex(0xff, 0xe9, 0xe6),
				hex(0xff, 0xf4, 0xf3),
				hex(0xff, 0xff, 0xff),
			},
			Success: tokens.Ramp{
				hex(0x00, 0x1d, 0x05),
				hex(0x00, 0x28, 0x08),
				hex(0x00, 0x37, 0x0e),
				hex(0x00, 0x53, 0x1a),
				hex(0x00, 0x7f, 0x2c),
				hex(0x41, 0xb0, 0x57),
				hex(0xb7, 0xff, 0xbe),
				hex(0xde, 0xff, 0xe0),
				hex(0xff, 0xff, 0xff),
			},
			Warning: tokens.Ramp{
				hex(0x2c, 0x10, 0x00),
				hex(0x3b, 0x18, 0x00),
				hex(0x4e, 0x22, 0x00),
				hex(0x74, 0x35, 0x00),
				hex(0xae, 0x53, 0x00),
				hex(0xf0, 0x7b, 0x1e),
				hex(0xff, 0xec, 0xc9),
				hex(0xff, 0xf5, 0xe4),
				hex(0xff, 0xff, 0xff),
			},
			Info: tokens.Ramp{
				hex(0x00, 0x18, 0x33),
				hex(0x00, 0x22, 0x44),
				hex(0x00, 0x2f, 0x5a),
				hex(0x00, 0x48, 0x84),
				hex(0x00, 0x6e, 0xc6),
				hex(0x3e, 0x9e, 0xff),
				hex(0xe2, 0xef, 0xff),
				hex(0xf0, 0xf7, 0xff),
				hex(0xff, 0xff, 0xff),
			},
		},
		Primary:     hex(0xd0, 0xc4, 0xff),
		OnPrimary:   hex(0x00, 0x00, 0x00),
		Secondary:   hex(0xd0, 0xc5, 0xfa),
		OnSecondary: hex(0x00, 0x00, 0x00),
		Tertiary:    hex(0xff, 0xb8, 0xd4),
		OnTertiary:  hex(0x00, 0x00, 0x00),
		Error:       hex(0xff, 0xbb, 0xb4),
		OnError:     hex(0x00, 0x00, 0x00),
		Success:     hex(0x77, 0xe3, 0x87),
		OnSuccess:   hex(0x00, 0x00, 0x00),
		Warning:     hex(0xff, 0xc2, 0x42),
		OnWarning:   hex(0x00, 0x00, 0x00),
		Info:        hex(0xa6, 0xd0, 0xff),
		OnInfo:      hex(0x00, 0x00, 0x00),
		Background:  hex(0x18, 0x18, 0x18),
		Text:        hex(0xff, 0xff, 0xff),
	}
	// Surface, Divider and the inverse pair are recorded as the resolutions
	// they are — the first two off this scheme's neutral ramp, the inverse
	// pair off the counterpart scheme's, which is what makes a light
	// scheme's inverse chip dark and a dark scheme's light.
	fill := func(t, counterpart tokens.ColorTokens) tokens.ColorTokens {
		n, o := t.Ramps.Neutral, counterpart.Ramps.Neutral
		t.Surface = n.Step(200)
		t.Divider = n.Step(500)
		t.InverseSurface = o.Step(200)
		t.OnInverseSurface = o.Step(900)
		return t
	}
	return fill(light, dark), fill(dark, light)
}

// TestFromSeedHighContrastGoldenPalette pins the high-contrast variant
// derived from the default seed — both schemes, all eight ramps, every pin
// and semantic colour — to the recording in hcGolden.
func TestFromSeedHighContrastGoldenPalette(t *testing.T) {
	wantLight, wantDark := hcGolden()
	light, dark := tokens.FromSeedHighContrast(tokens.DefaultSeed)
	diffTokens(t, "FromSeedHighContrast light", light, wantLight)
	diffTokens(t, "FromSeedHighContrast dark", dark, wantDark)
}

// TestFromSeedHighContrastSharesGrounds verifies the variant widens tone
// separation without moving the grounds or the light pins: steps 100–600 of
// every ramp are byte-identical to FromSeed's, and the light pinned bases —
// the seed-exact Primary included — and their White on-colours carry over
// unchanged.
func TestFromSeedHighContrastSharesGrounds(t *testing.T) {
	seeds := []color.NRGBA{
		tokens.DefaultSeed,
		{0x3b, 0x82, 0xf6, 0xff},
	}
	for _, seed := range seeds {
		light, dark := tokens.FromSeed(seed)
		hcLight, hcDark := tokens.FromSeedHighContrast(seed)
		for i := range namedRamps(light) {
			for step := 100; step <= 600; step += 100 {
				if g, w := namedRamps(hcLight)[i].ramp.Step(step), namedRamps(light)[i].ramp.Step(step); g != w {
					t.Errorf("seed %v: HC light %s.Step(%d) = %v, want the default's %v",
						seed, namedRamps(light)[i].name, step, g, w)
				}
				if g, w := namedRamps(hcDark)[i].ramp.Step(step), namedRamps(dark)[i].ramp.Step(step); g != w {
					t.Errorf("seed %v: HC dark %s.Step(%d) = %v, want the default's %v",
						seed, namedRamps(dark)[i].name, step, g, w)
				}
			}
		}
		lightPins := [][2]color.NRGBA{
			{hcLight.Primary, light.Primary}, {hcLight.OnPrimary, light.OnPrimary},
			{hcLight.Secondary, light.Secondary}, {hcLight.OnSecondary, light.OnSecondary},
			{hcLight.Tertiary, light.Tertiary}, {hcLight.OnTertiary, light.OnTertiary},
			{hcLight.Error, light.Error}, {hcLight.OnError, light.OnError},
			{hcLight.Success, light.Success}, {hcLight.OnSuccess, light.OnSuccess},
			{hcLight.Warning, light.Warning}, {hcLight.OnWarning, light.OnWarning},
			{hcLight.Info, light.Info}, {hcLight.OnInfo, light.OnInfo},
		}
		for i, p := range lightPins {
			if p[0] != p[1] {
				t.Errorf("seed %v: HC light pin/on #%d = %v, want the default's %v", seed, i, p[0], p[1])
			}
		}
		if hcLight.Primary != light.Primary {
			t.Errorf("seed %v: HC light Primary = %v, want FromSeed's lifted base %v",
				seed, hcLight.Primary, light.Primary)
		}
	}
}

// diffTokens reports every field and ramp step where got differs from want.
func diffTokens(t *testing.T, scheme string, got, want tokens.ColorTokens) {
	t.Helper()
	if got == want {
		return
	}
	gotRamps, wantRamps := namedRamps(got), namedRamps(want)
	for i := range gotRamps {
		for step := 100; step <= 900; step += 100 {
			if g, w := gotRamps[i].ramp.Step(step), wantRamps[i].ramp.Step(step); g != w {
				t.Errorf("%s: %s.Step(%d) = %v, want %v", scheme, gotRamps[i].name, step, g, w)
			}
		}
	}
	fields := func(c tokens.ColorTokens) map[string]color.NRGBA {
		return map[string]color.NRGBA{
			"Primary": c.Primary, "OnPrimary": c.OnPrimary,
			"Secondary": c.Secondary, "OnSecondary": c.OnSecondary,
			"Tertiary": c.Tertiary, "OnTertiary": c.OnTertiary,
			"Error": c.Error, "OnError": c.OnError,
			"Success": c.Success, "OnSuccess": c.OnSuccess,
			"Warning": c.Warning, "OnWarning": c.OnWarning,
			"Info": c.Info, "OnInfo": c.OnInfo,
			"Background": c.Background, "Text": c.Text,
			"Surface": c.Surface, "Divider": c.Divider,
			"InverseSurface": c.InverseSurface, "OnInverseSurface": c.OnInverseSurface,
		}
	}
	g, w := fields(got), fields(want)
	for name := range w {
		if g[name] != w[name] {
			t.Errorf("%s: %s = %v, want %v", scheme, name, g[name], w[name])
		}
	}
}

// TestFromSeedGoldenPalette pins the entire palette derived from the
// default seed — both schemes, all eight ramps, every pin and semantic
// colour — to the recording in defaultGolden, and verifies DefaultLight
// and DefaultDark are exactly that derivation.
func TestFromSeedGoldenPalette(t *testing.T) {
	wantLight, wantDark := defaultGolden()
	light, dark := tokens.FromSeed(tokens.DefaultSeed)
	diffTokens(t, "FromSeed light", light, wantLight)
	diffTokens(t, "FromSeed dark", dark, wantDark)
	diffTokens(t, "DefaultLight", tokens.DefaultLight, wantLight)
	diffTokens(t, "DefaultDark", tokens.DefaultDark, wantDark)
}

func TestWCAGAAContrast(t *testing.T) {
	schemes := []struct {
		name   string
		tokens tokens.ColorTokens
	}{
		{"DefaultLight", tokens.DefaultLight},
		{"DefaultDark", tokens.DefaultDark},
	}
	for _, s := range schemes {
		for _, p := range tokenPairs(s.tokens) {
			cr := contrastRatio(p.bg, p.fg)
			if cr < wcagAA {
				t.Errorf("%s %s: contrast ratio %.2f:1 < %.1f:1 (WCAG AA)",
					s.name, p.name, cr, wcagAA)
			}
		}
	}
}
