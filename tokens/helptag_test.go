package tokens_test

import (
	"image/color"
	"math"
	"testing"

	"github.com/vibrantgio/theme/tokens"
)

// helpTagRow is what the reference's MEASURED help tag rows record for one
// scheme, every value an alpha out of 255.
type helpTagRow struct {
	name      string
	set       tokens.PlatformColors
	edge      uint8
	corner    uint8
	rim       [2]uint8
	shadow    []uint8
	diagonal  [4]uint8 // the shadow on the diagonal out from the corner, at dx = dy = 2..5
	offsetDp  float32
	edgeShown [2]uint8 // the edge as the capture shows it, the shadow under it: top, side
}

var helpTagRows = []helpTagRow{
	{
		name: "light", set: tokens.PlatformLight,
		edge: 45, corner: 30, rim: [2]uint8{23, 0},
		shadow:    []uint8{28, 23, 18, 14, 10, 7, 5, 3, 2, 1},
		diagonal:  [4]uint8{6, 4, 2, 1},
		offsetDp:  2,
		edgeShown: [2]uint8{56, 64},
	},
	{
		name: "dark", set: tokens.PlatformDark,
		edge: 115, corner: 90, rim: [2]uint8{39, 8},
		shadow:    []uint8{35, 29, 23, 17, 13, 9, 6, 4, 2, 1, 1},
		diagonal:  [4]uint8{8, 5, 2, 1},
		offsetDp:  2,
		edgeShown: [2]uint8{124, 131},
	},
}

// TestHelpTagEdgeAndShadowAreTheMeasuredRows pins the help tag's edge, its
// corner pixel, its rim and its shadow curve to the reference's MEASURED
// rows in both schemes. The edge and the shadow are black coverages and the
// rim a white one.
func TestHelpTagEdgeAndShadowAreTheMeasuredRows(t *testing.T) {
	black := func(a uint8) color.NRGBA { return color.NRGBA{A: a} }
	white := func(a uint8) color.NRGBA {
		if a == 0 {
			return color.NRGBA{}
		}
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: a}
	}
	for _, r := range helpTagRows {
		e := r.set.HelpTagEdge
		if e.Side != black(r.edge) {
			t.Errorf("%s: HelpTagEdge.Side = %v, the reference reads black at %d", r.name, e.Side, r.edge)
		}
		if e.Corner != black(r.corner) {
			t.Errorf("%s: HelpTagEdge.Corner = %v, the reference reads black at %d", r.name, e.Corner, r.corner)
		}
		for i, a := range r.rim {
			if e.Rim[i] != white(a) {
				t.Errorf("%s: HelpTagEdge.Rim[%d] = %v, the reference reads white at %d", r.name, i, e.Rim[i], a)
			}
		}

		s := r.set.HelpTagShadow
		if s.Offset != r.offsetDp {
			t.Errorf("%s: HelpTagShadow.Offset = %v, the reference reads %v", r.name, s.Offset, r.offsetDp)
		}
		for d := range s.Alpha {
			var want color.NRGBA
			if d < len(r.shadow) {
				want = black(r.shadow[d])
			}
			if s.Alpha[d] != want {
				t.Errorf("%s: HelpTagShadow.Alpha[%d] = %v, the reference reads %v", r.name, d, s.Alpha[d], want)
			}
		}
		if s.Reach() != len(r.shadow) {
			t.Errorf("%s: HelpTagShadow.Reach() = %d, the reference reads %d pixels", r.name, s.Reach(), len(r.shadow))
		}

		// The corner follows the product rule: the measured diagonal is
		// within 1 of 255 of Alpha[d]·Alpha[d] / (2·Alpha[0]).
		for i, a := range r.diagonal {
			d := i + 2
			if got := s.CornerAlpha(d, d); math.Abs(got-float64(a)) > 1 {
				t.Errorf("%s: the shadow's corner rule at dx = dy = %d gives %.2f, the capture reads %d", r.name, d, got, a)
			}
		}
	}
}

// TestHelpTagEdgeIsTheCaptureWithTheShadowOut holds the arithmetic of the
// reference rows: the edge laid over the shadow under it lands the pixel the
// capture shows, within 1 of 255, on the top and the sides. The shadow under
// the top edge is index 3 of the curve and under a side edge index 1, as the
// edge is one pixel outside the tag and the shadow's rectangle is the tag
// moved 2 px down. The bottom edge is inside that rectangle, where the curve
// does not reach, so its shadow comes from the reference's fit and not from
// this field.
func TestHelpTagEdgeIsTheCaptureWithTheShadowOut(t *testing.T) {
	for _, r := range helpTagRows {
		e := float64(r.set.HelpTagEdge.Side.A)
		for _, c := range []struct {
			where string
			under int
			shown uint8
		}{
			{"top", 3, r.edgeShown[0]},
			{"side", 1, r.edgeShown[1]},
		} {
			s := float64(r.set.HelpTagShadow.Alpha[c.under].A)
			got := s + e*(255-s)/255
			if math.Abs(got-float64(c.shown)) > 1 {
				t.Errorf("%s %s: the edge over the shadow lands %.2f, the capture reads %d", r.name, c.where, got, c.shown)
			}
		}
	}
}
