package tokens_test

import (
	"math"
	"testing"

	"github.com/vibrantgio/theme/tokens"
)

// The shipped face's own numbers, read from its tables: Roboto's OS/2
// sCapHeight is 1456 units of a 2048 unit em, and the Medium weight's 'I'
// outline is 252 units wide. Every role in this system is set in that face, so
// these two ratios are what FaceMetrics answers with everywhere.
const (
	robotoCapEm       = 1456.0 / 2048.0
	robotoMediumStem  = 252.0 / 2048.0
	robotoRegularStem = 193.0 / 2048.0
)

func closeTo(got, want float64) bool { return math.Abs(got-want) < 0.001 }

func TestFaceMetricsReadsTheShippedFace(t *testing.T) {
	for _, tc := range []struct {
		name    string
		style   tokens.TextStyle
		wantCap float64
		wantSte float64
	}{
		{
			"LabelLarge is Medium",
			tokens.DefaultTypography.LabelLarge,
			14 * robotoCapEm, 14 * robotoMediumStem,
		},
		{
			"BodyLarge is Regular",
			tokens.DefaultTypography.BodyLarge,
			float64(tokens.DefaultTypography.BodyLarge.Size) * robotoCapEm,
			float64(tokens.DefaultTypography.BodyLarge.Size) * robotoRegularStem,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := tc.style.FaceMetrics()
			if !closeTo(float64(m.CapHeight), tc.wantCap) {
				t.Errorf("CapHeight = %v, want %v", m.CapHeight, tc.wantCap)
			}
			if !closeTo(float64(m.Stem), tc.wantSte) {
				t.Errorf("Stem = %v, want %v", m.Stem, tc.wantSte)
			}
		})
	}
}

// A role naming a family nobody handed over still answers, at the shipped
// face's ratios — the documented guess, pinned so it cannot become a zero.
func TestFaceMetricsFallsBackForAnUnknownFamily(t *testing.T) {
	m := tokens.TextStyle{Typeface: "Nothing Here", Weight: tokens.WeightMedium, Size: 20}.FaceMetrics()
	if m.CapHeight <= 0 || m.Stem <= 0 {
		t.Fatalf("unknown family gave %+v, want both positive", m)
	}
	if m.CapHeight >= 20 || m.Stem >= m.CapHeight {
		t.Errorf("unknown family gave %+v, want a cap under the size and a stem under the cap", m)
	}
}

// The cap band scales with the role, which is the whole point of asking the
// face rather than pinning a number: a larger role gets a taller band.
func TestFaceMetricsScalesWithTheRole(t *testing.T) {
	small := tokens.DefaultTypography.LabelSmall.FaceMetrics()
	large := tokens.DefaultTypography.HeadlineLarge.FaceMetrics()
	if !(small.CapHeight < large.CapHeight) {
		t.Errorf("LabelSmall cap %v is not under HeadlineLarge cap %v", small.CapHeight, large.CapHeight)
	}
}
