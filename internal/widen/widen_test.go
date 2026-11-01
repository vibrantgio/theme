package widen

import (
	"math"
	"strconv"
	"testing"

	"github.com/vibrantgio/font/roboto/regular/normal"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// TestOnlyTheXValuesMove: every glyph of the widened face has the shipped
// face's segments with the same y values, and the shipped advance. So the
// cap height, the x-height and the spacing of the text stay. A negative
// amount narrows the stems under the same rule.
func TestOnlyTheXValuesMove(t *testing.T) {
	for _, em := range []float64{0.03, -0.015} {
		t.Run(strconv.FormatFloat(em, 'f', -1, 64), func(t *testing.T) { onlyTheXValuesMove(t, em) })
	}
}

func onlyTheXValuesMove(t *testing.T, em float64) {
	wide, err := Sideways(normal.TTF, em)
	if err != nil {
		t.Fatal(err)
	}
	a, err := sfnt.Parse(normal.TTF)
	if err != nil {
		t.Fatal(err)
	}
	b, err := sfnt.Parse(wide)
	if err != nil {
		t.Fatal(err)
	}
	if a.NumGlyphs() != b.NumGlyphs() {
		t.Fatalf("the widened face has %d glyphs, want %d", b.NumGlyphs(), a.NumGlyphs())
	}
	ppem := fixed.I(2048)
	var ba, bb sfnt.Buffer
	moved := 0
	for i := 0; i < a.NumGlyphs(); i++ {
		g := sfnt.GlyphIndex(i)
		adva, err := a.GlyphAdvance(&ba, g, ppem, font.HintingNone)
		if err != nil {
			t.Fatal(err)
		}
		advb, err := b.GlyphAdvance(&bb, g, ppem, font.HintingNone)
		if err != nil {
			t.Fatal(err)
		}
		if adva != advb {
			t.Errorf("glyph %d: advance %v, want %v", i, advb, adva)
		}
		sa, err := a.LoadGlyph(&ba, g, ppem, nil)
		if err != nil {
			t.Fatal(err)
		}
		sa = append(sfnt.Segments(nil), sa...)
		sb, err := b.LoadGlyph(&bb, g, ppem, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(sa) != len(sb) {
			t.Errorf("glyph %d: %d segments, want %d", i, len(sb), len(sa))
			continue
		}
		for j := range sa {
			if sa[j].Op != sb[j].Op {
				t.Errorf("glyph %d segment %d: op %v, want %v", i, j, sb[j].Op, sa[j].Op)
			}
			for k := range sa[j].Args {
				if sa[j].Args[k].Y != sb[j].Args[k].Y {
					t.Errorf("glyph %d segment %d: y %v, want %v", i, j, sb[j].Args[k].Y, sa[j].Args[k].Y)
				}
				if sa[j].Args[k].X != sb[j].Args[k].X {
					moved++
				}
			}
		}
	}
	if moved == 0 {
		t.Error("no x value moved")
	}
}

// TestAStemIsWiderByTheAmount: the l is one vertical stem, so its box is
// wider by the amount, to the rounding of one font unit on each side.
func TestAStemIsWiderByTheAmount(t *testing.T) {
	const em = 0.03
	wide, err := Sideways(normal.TTF, em)
	if err != nil {
		t.Fatal(err)
	}
	width := func(ttf []byte) fixed.Int26_6 {
		f, err := sfnt.Parse(ttf)
		if err != nil {
			t.Fatal(err)
		}
		var buf sfnt.Buffer
		g, err := f.GlyphIndex(&buf, 'l')
		if err != nil {
			t.Fatal(err)
		}
		r, _, err := f.GlyphBounds(&buf, g, fixed.I(2048), font.HintingNone)
		if err != nil {
			t.Fatal(err)
		}
		return r.Max.X - r.Min.X
	}
	got := (width(wide) - width(normal.TTF)).Round()
	if want := int(math.Round(em * 2048)); got < want-2 || got > want+2 {
		t.Errorf("the l is %d font units wider, want %d", got, want)
	}
}
