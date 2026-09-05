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

// TestTheFillTakesAllTheChromaItsDepthHolds verifies the two halves of the
// marker's chroma: the depth the fill is realized at is the shallowest step
// that holds markerChroma, and the fill realized there carries all the
// chroma sRGB holds at that depth — a highlighter owes its own chroma, not
// a dial borrowed from the containers. The margin is eight-bit
// quantization: the realization is rounded to a byte per channel, which
// costs the measured chroma a little of what the solver found.
func TestTheFillTakesAllTheChromaItsDepthHolds(t *testing.T) {
	light, dark := FromSeed(DefaultSeed)
	for _, d := range []struct {
		name string
		tok  ColorTokens
		step int     // the depth the fill is realized at
		tone float64 // that depth's L*
		held float64 // the most chroma sRGB holds at highlightHue there
	}{
		{"light", light, 300, 84.91, 0.1845},
		{"dark", dark, 400, 30.16, 0.0850},
	} {
		if got := d.tok.highlightStep(); got != d.step {
			t.Errorf("%s: the fill is realized at step %d, want step %d — the shallowest depth that holds the marker's chroma", d.name, got, d.step)
		}
		tone, _, _ := color.LabFromNRGBA(d.tok.Ramps.Neutral.Step(d.step))
		if math.Abs(tone-d.tone) > 0.01 {
			t.Errorf("%s: step %d stands at L* %.2f, want %.2f", d.name, d.step, tone, d.tone)
		}
		held := 0.0
		for c := 0.0; c < highlightChroma; c += 0.0005 {
			_, realized, _ := color.OKLChFromNRGBA(color.NRGBAFromToneChromaHue(tone, c, highlightHue))
			if realized < c-0.002 { // past the gamut the solver reduces chroma
				break
			}
			held = c
		}
		if math.Abs(held-d.held) > 0.001 {
			t.Errorf("%s: sRGB holds chroma %.4f at hue %.2f° at L* %.2f, want %.4f", d.name, held, highlightHue, tone, d.held)
		}
		if held < markerChroma {
			t.Errorf("%s: the depth the fill is realized at holds only chroma %.4f, under the %.3f a marker's yellow asks for",
				d.name, held, markerChroma)
		}
		_, chroma, hue := color.OKLChFromNRGBA(d.tok.Highlight)
		if chroma < held-0.002 {
			t.Errorf("%s: the fill %v carries chroma %.4f where sRGB holds %.4f — it is not taking the yellow its depth can hold",
				d.name, d.tok.Highlight, chroma, held)
		}
		if math.Abs(hue-highlightHue) > 0.5 {
			t.Errorf("%s: the fill %v wears hue %.2f°, want the reserved %.2f°", d.name, d.tok.Highlight, hue, highlightHue)
		}
	}
}
