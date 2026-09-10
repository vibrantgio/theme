package tokens

import (
	stdcolor "image/color"
	"math"
	"sort"
	"testing"

	"github.com/vibrantgio/theme/color"
)

// TestTheMarkerYellowFallsInTheRunTheStatusesLeaveOpen recomputes, from the
// four status anchors themselves, the run of the hue circle the yellows
// occupy — sort the anchors, take the gap between the warning's orange and
// the success's green — and places the marker's own hue inside it. Deriving
// the run rather than asserting the colour is what keeps the reservation true
// if an anchor ever moves: the run would stop holding the yellow rather than
// quietly closing on it.
func TestTheMarkerYellowFallsInTheRunTheStatusesLeaveOpen(t *testing.T) {
	_, _, yellow := color.OKLChFromNRGBA(markerYellow)
	anchors := []float64{errorHue, successHue, warningHue, infoHue}
	sort.Float64s(anchors)
	run, from := 0.0, 0.0
	found := false
	for i, a := range anchors {
		b := anchors[(i+1)%len(anchors)]
		length := math.Mod(b-a+360, 360)
		if math.Mod(yellow-a+360, 360) >= length {
			continue // the marker is not in this run
		}
		run, from, found = length, a, true
	}
	if !found {
		t.Fatalf("no run between two status anchors holds the marker's yellow at %.2f° — a status anchor has moved into it", yellow)
	}
	if math.Abs(run-80.15) > 0.01 {
		t.Errorf("the run the yellows occupy measures %.4f°, want 80.15° (warning %.2f° to success %.2f°)",
			run, warningHue, successHue)
	}
	// The marker stands clear of both ends of the run, not merely inside it:
	// a yellow a few degrees off the warning's orange would report a status.
	const clearance = 20.0
	if gap := math.Mod(yellow-from+360, 360); gap < clearance || run-gap < clearance {
		t.Errorf("the marker's yellow at %.2f° stands %.2f° past the warning and %.2f° short of the success — under the %.1f° it is reserved by",
			yellow, gap, run-gap, clearance)
	}
	t.Logf("the marker's yellow measures %.2f° in the %.2f° run the statuses leave open", yellow, run)
}

// TestTheCoverageIsTheWholeOfTheBlend pins the two ends of yellowOver: no
// coverage is the surface untouched and full coverage is the marker yellow
// itself, so the coverage is the only thing between them and a highlight can
// never land on a colour that is neither.
func TestTheCoverageIsTheWholeOfTheBlend(t *testing.T) {
	for _, surface := range []stdcolor.NRGBA{
		{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		{R: 0x1e, G: 0x1e, B: 0x1e, A: 0xff},
		{R: 0x67, G: 0x50, B: 0xa4, A: 0xff},
	} {
		if got := yellowOver(0, surface); got != surface {
			t.Errorf("no coverage over %v landed on %v, want the surface itself", surface, got)
		}
		if got := yellowOver(1, surface); got != markerYellow {
			t.Errorf("full coverage over %v landed on %v, want the marker yellow %v", surface, got, markerYellow)
		}
	}
}
