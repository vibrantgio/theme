package tokens

import (
	"math"
	"sort"
	"testing"

	"github.com/vibrantgio/theme/color"
)

// TestHighlightHueIsTheWidestArcTheStatusAnchorsLeave recomputes the
// reserved hue from the four status anchors themselves: sort them around
// the OKLCh hue circle, take the widest run between two neighbours, and
// the reserved hue is its midpoint. Pinning the constant against its own
// derivation is what keeps the reservation true if an anchor ever moves —
// the constant would stop matching rather than quietly drift toward a
// status.
func TestHighlightHueIsTheWidestArcTheStatusAnchorsLeave(t *testing.T) {
	anchors := []float64{errorHue, successHue, warningHue, infoHue}
	sort.Float64s(anchors)
	widest, midpoint := 0.0, 0.0
	for i, a := range anchors {
		b := anchors[(i+1)%len(anchors)]
		run := math.Mod(b-a+360, 360)
		if run > widest {
			widest, midpoint = run, math.Mod(a+run/2, 360)
		}
	}
	if math.Abs(widest-139.9) > 1e-9 {
		t.Errorf("widest run between two status anchors = %.4f°, want 139.9° (info %.1f° to error %.1f°)",
			widest, infoHue, errorHue)
	}
	if math.Abs(midpoint-highlightHue) > 1e-9 {
		t.Errorf("highlightHue = %.4f°, want %.4f° — the midpoint of the widest arc the anchors leave",
			highlightHue, midpoint)
	}
}

// TestTheHighlightDialFitsTheGamutAtBothWashDepths verifies sRGB holds the
// container dial at the reserved hue at both depths the wash is realized
// at, with the headroom the file header records: a clipped chroma would
// take the reservation's distance from the status hues away silently.
func TestTheHighlightDialFitsTheGamutAtBothWashDepths(t *testing.T) {
	for _, d := range []struct {
		name string
		tone int
		want float64 // the most chroma sRGB holds at highlightHue at this depth
	}{
		{"light step 300", 85, 0.0935},
		{"dark step 300", 19, 0.1575},
	} {
		got := 0.0
		for c := 0.0; c < 0.4; c += 0.0005 {
			_, realized, _ := color.OKLChFromNRGBA(color.Tone(highlightHue, c, d.tone))
			if realized < c-0.002 { // past the gamut the solver reduces chroma
				break
			}
			got = c
		}
		if math.Abs(got-d.want) > 0.001 {
			t.Errorf("%s (L* %d): sRGB holds chroma %.4f at hue %.2f°, want %.4f", d.name, d.tone, got, highlightHue, d.want)
		}
		if got <= containerChroma {
			t.Errorf("%s (L* %d): sRGB holds only chroma %.4f at hue %.2f°, under the %.3f dial the wash is realized at",
				d.name, d.tone, got, highlightHue, containerChroma)
		}
	}
}
