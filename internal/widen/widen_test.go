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

// TestWhatMoves: every glyph of a face from Sideways has the shipped face's
// segments with the same y values, and the shipped advance. So the cap
// height, the x-height and the spacing of the text stay. A negative amount
// narrows the stems under the same rule. A face from Outward with a
// vertical amount also moves y values, and it keeps every advance.
func TestWhatMoves(t *testing.T) {
	for _, em := range []float64{0.03, -0.015} {
		t.Run(strconv.FormatFloat(em, 'f', -1, 64), func(t *testing.T) {
			wide, err := Sideways(normal.TTF, em)
			if err != nil {
				t.Fatal(err)
			}
			whatMoves(t, wide, false)
		})
	}
	t.Run("outward", func(t *testing.T) {
		wide, err := Outward(normal.TTF, 0.03, 0.027)
		if err != nil {
			t.Fatal(err)
		}
		whatMoves(t, wide, true)
	})
}

func whatMoves(t *testing.T, wide []byte, yMoves bool) {
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
	moved, movedY := 0, 0
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
					movedY++
				}
				if !yMoves && sa[j].Args[k].Y != sb[j].Args[k].Y {
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
	if yMoves && movedY == 0 {
		t.Error("no y value moved")
	}
}

// TestTheHeightsMoveByHalfTheAmount: the top of the T and of the x is a
// horizontal stroke with two corners, so a vertical amount ey moves it up
// by half of ey, to the rounding of one font unit. The cap height and the
// x-height of the drawn text therefore move up by half of ey: 0.18 px at
// 13 pt and 0.14 px at 10 pt for ey = 0.027 em.
func TestTheHeightsMoveByHalfTheAmount(t *testing.T) {
	const ey = 0.027
	tall, err := Outward(normal.TTF, 0.03, ey)
	if err != nil {
		t.Fatal(err)
	}
	top := func(ttf []byte, r rune) fixed.Int26_6 {
		f, err := sfnt.Parse(ttf)
		if err != nil {
			t.Fatal(err)
		}
		var buf sfnt.Buffer
		g, err := f.GlyphIndex(&buf, r)
		if err != nil {
			t.Fatal(err)
		}
		b, _, err := f.GlyphBounds(&buf, g, fixed.I(2048), font.HintingNone)
		if err != nil {
			t.Fatal(err)
		}
		return -b.Min.Y
	}
	want := int(math.Round(ey * 2048 / 2))
	for _, r := range []rune{'T', 'x'} {
		if got := (top(tall, r) - top(normal.TTF, r)).Round(); got < want-1 || got > want+1 {
			t.Errorf("the top of the %c moves up %d font units, want %d", r, got, want)
		}
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
