package tokens

import (
	"math"
	"sort"
	"testing"

	"github.com/vibrantgio/theme/color"
)

// markerYellow is the OKLCh hue of Material Yellow 500 #FFEB3B, the
// canonical anchor of the yellow family. It is the witness the derivation
// below picks its arc with, and nothing derives from it.
const markerYellow = 102.50

// TestHighlightHueIsTheMidpointOfTheArcTheYellowsOccupy recomputes the
// reserved hue from the four status anchors themselves: sort them around
// the OKLCh hue circle, take the run between two neighbours that the
// yellows fall in, and the reserved hue is its midpoint — the point in the
// yellow furthest from either status beside it. Pinning the constant
// against its own derivation is what keeps the reservation true if an
// anchor ever moves: the constant would stop matching rather than quietly
// drift toward a status.
func TestHighlightHueIsTheMidpointOfTheArcTheYellowsOccupy(t *testing.T) {
	anchors := []float64{errorHue, successHue, warningHue, infoHue}
	sort.Float64s(anchors)
	run, midpoint := 0.0, 0.0
	found := false
	for i, a := range anchors {
		b := anchors[(i+1)%len(anchors)]
		length := math.Mod(b-a+360, 360)
		if math.Mod(markerYellow-a+360, 360) >= length {
			continue // the yellow is not in this run
		}
		run, midpoint, found = length, math.Mod(a+length/2, 360), true
	}
	if !found {
		t.Fatalf("no run between two status anchors holds the marker's yellow at %.2f° — a status anchor has moved into it", markerYellow)
	}
	if math.Abs(run-80.15) > 0.01 {
		t.Errorf("the run the yellows occupy measures %.4f°, want 80.15° (warning %.2f° to success %.2f°)",
			run, warningHue, successHue)
	}
	if math.Abs(midpoint-highlightHue) > 0.005 {
		t.Errorf("highlightHue = %.4f°, want %.4f° — the midpoint of the run the yellows occupy",
			highlightHue, midpoint)
	}
	if gap := math.Abs(highlightHue - markerYellow); gap > 2.0 {
		t.Errorf("the reserved hue sits %.2f° from Material Yellow 500's %.2f° — that is no longer the yellow a palette would have named",
			gap, markerYellow)
	}
}

// TestTheHighlightDialFitsTheGamutAtBothFillDepths verifies sRGB holds the
// container dial at the reserved hue at both depths the fill is realized
// at, with the headroom the file header records: a clipped chroma would
// take the reservation's distance from the status hues away silently.
func TestTheHighlightDialFitsTheGamutAtBothFillDepths(t *testing.T) {
	for _, d := range []struct {
		name string
		tone int
		want float64 // the most chroma sRGB holds at highlightHue at this depth
	}{
		{"light step 300", 85, 0.1850},
		{"dark step 300", 19, 0.0650},
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
			t.Errorf("%s (L* %d): sRGB holds only chroma %.4f at hue %.2f°, under the %.3f dial the fill is realized at",
				d.name, d.tone, got, highlightHue, containerChroma)
		}
	}
}
