package tokens

import (
	"image/color"
	"math"
	"testing"

	"gioui.org/font"
	"gioui.org/text"
	"golang.org/x/image/math/fixed"

	"github.com/vibrantgio/font/roboto"

	vgcolor "github.com/vibrantgio/theme/color"
)

// TestEachStemFamilyWidensOnlyRobotoRegular: the shaper's collection is the
// faces as shipped, then a copy of every Roboto entry under each stem
// weight the shaper holds but zero. Of the copies, only Roboto Regular
// upright has a face of its own, and each stem weight has its own.
func TestEachStemFamilyWidensOnlyRobotoRegular(t *testing.T) {
	faces := DefaultTypography.Faces
	got := stemFaces(faces)
	for i := range faces {
		if got[i] != faces[i] {
			t.Fatalf("face %d is %v, want the face as shipped %v", i, got[i].Font, faces[i].Font)
		}
	}
	shipped := map[font.Font]font.Face{}
	for _, f := range faces {
		if f.Font.Typeface == roboto.RegularNormal.Typeface {
			shipped[f.Font] = f.Face
		}
	}
	families := stemStepMax - stemStepMin
	copies := got[len(faces):]
	if len(copies) != families*len(shipped) {
		t.Fatalf("%d stem faces, want one per Roboto face for each of %d stem weights, %d", len(copies), families, families*len(shipped))
	}
	widenedFaces := map[font.Typeface]font.Face{}
	for _, c := range copies {
		f := c.Font
		f.Typeface = roboto.RegularNormal.Typeface
		widened := c.Face != shipped[f]
		if widened != (f == roboto.RegularNormal) {
			t.Errorf("%v: widened %v, want widened only for Roboto Regular upright", c.Font, widened)
		}
		if widened {
			widenedFaces[c.Font.Typeface] = c.Face
		}
	}
	seen := map[font.Face]bool{}
	for k := stemStepMin; k <= stemStepMax; k++ {
		if k == 0 {
			continue
		}
		f, ok := widenedFaces[stemTypeface(k)]
		if !ok || seen[f] {
			t.Errorf("stem weight %d × %v has no face of its own", k, stemStep)
		}
		seen[f] = true
	}
}

// TestEveryCopySharesTheShaper: copies of a typography hand back one shaper
// of each kind.
func TestEveryCopySharesTheShaper(t *testing.T) {
	a, b := DefaultTypography, DefaultTypography
	if a.DeterministicShaper() != b.DeterministicShaper() || a.Shaper() != b.Shaper() {
		t.Error("two copies built two shapers")
	}
}

// TestTheStemWeightFollowsSizeAndTheTextAgainstItsFill: at one size, the
// stem weight of darker text falls as its contrast with the fill falls, and
// the stem weight of lighter text rises. Above 13 pt the stem weight falls
// as the size grows. In each case the stem as stored that the weight gives,
// the stem as shipped plus the weight at the size, plus the blend
// correction, is the platform's stem to half a step.
func TestTheStemWeightFollowsSizeAndTheTextAgainstItsFill(t *testing.T) {
	white := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	l, d := PlatformLight, PlatformDark
	type pair struct {
		name     string
		fg, fill color.NRGBA
		lighter  bool
	}
	light := []pair{
		{"light: label colour", l.Label, l.WindowBackground, false},
		{"light: secondary label colour", l.SecondaryLabel, l.WindowBackground, false},
		{"light: tertiary label colour", l.TertiaryLabel, l.WindowBackground, false},
	}
	dark := []pair{
		{"dark: label colour", d.Label, d.WindowBackground, true},
		{"dark: secondary label colour", d.SecondaryLabel, d.WindowBackground, true},
		{"dark: white on the default button", white, d.DefaultButtonFill, true},
	}
	for _, ps := range [][]pair{light, dark} {
		for i := 1; i < len(ps); i++ {
			a, b := TextStemWeight(13, ps[i-1].fg, ps[i-1].fill), TextStemWeight(13, ps[i].fg, ps[i].fill)
			if ps[i].lighter && a > b || !ps[i].lighter && a < b {
				t.Errorf("13 pt: %s %v, then %s %v", ps[i-1].name, a, ps[i].name, b)
			}
		}
	}
	for _, p := range append(light, dark...) {
		prev := StemWeight(math.MaxFloat32)
		for _, size := range []float32{10, 11, 12, 13, 15, 17, 22, 26, 28} {
			w := TextStemWeight(size, p.fg, p.fill)
			if size > 13 && w > prev {
				t.Errorf("%s: %v at %v pt, over %v at the size before", p.name, w, size, prev)
			}
			prev = w
			flat := vgcolor.Flatten(p.fg, p.fill)
			stored := (shippedStem+float64(w))*float64(size) + blendCorrection(flat, p.fill)
			want := platformStem(float64(size), p.lighter)
			if math.Abs(stored-want) > stemStep*float64(size)/2+1e-6 {
				t.Errorf("%s at %v pt: weight %v gives %.3f px, want the platform's %.3f px", p.name, size, w, stored, want)
			}
		}
	}
}

// TestThePlatformStemFollowsTheRowsAndTheLine: the platform's stem is the
// measured rows at 10 to 13 pt, and above 13 pt the DERIVED line through
// the 13 pt and 28 pt rows.
func TestThePlatformStemFollowsTheRowsAndTheLine(t *testing.T) {
	for _, tc := range []struct {
		size    float64
		lighter bool
		want    float64
	}{
		{10, false, 0.97}, {11, false, 1.06}, {13, false, 1.26},
		{10, true, 1.16}, {11, true, 1.29}, {13, true, 1.52},
		{20.5, false, 1.705}, {28, false, 2.15}, {20.5, true, 1.98}, {28, true, 2.44},
	} {
		if got := platformStem(tc.size, tc.lighter); math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("%v pt, lighter %v: %v px, want %v", tc.size, tc.lighter, got, tc.want)
		}
	}
}

// TestStemFontNamesAStemFamilyForRobotoOnly: text in Roboto, named or as
// the empty typeface, shapes in the family of its stem weight, and at the
// stem weight zero in Roboto as shipped. Text in any other typeface keeps
// its font. No font changes its weight or style.
func TestStemFontNamesAStemFamilyForRobotoOnly(t *testing.T) {
	light := PlatformLight
	dark := PlatformDark
	mono := font.Font{Typeface: CodeFaceRoboto}
	for _, tc := range []struct {
		name     string
		f        font.Font
		size     float32
		fg, fill color.NRGBA
	}{
		{"Roboto, darker text", roboto.RegularMedium, 13, light.Label, light.WindowBackground},
		{"empty typeface, darker text", font.Font{}, 13, light.Label, light.WindowBackground},
		{"Roboto, lighter text", roboto.RegularNormal, 13, dark.Label, dark.WindowBackground},
		{"empty typeface, lighter text", font.Font{}, 28, dark.Label, dark.WindowBackground},
		{"Roboto, no size", roboto.RegularNormal, 0, light.Label, light.WindowBackground},
	} {
		got := StemFont(tc.f, tc.size, tc.fg, tc.fill)
		want := tc.f.Typeface
		if k := stemStepOf(tc.size, tc.fg, tc.fill); k != 0 {
			want = stemTypeface(k)
		}
		if got.Typeface != want {
			t.Errorf("%s: typeface %q, want %q", tc.name, got.Typeface, want)
		}
		if got.Weight != tc.f.Weight || got.Style != tc.f.Style {
			t.Errorf("%s: %v changed its weight or style to %v", tc.name, tc.f, got)
		}
	}
	if got := StemFont(mono, 13, light.Label, light.WindowBackground); got != mono {
		t.Errorf("Roboto Mono, darker text: %v, want %v", got, mono)
	}
	if k := stemStepOf(13, light.Label, light.WindowBackground); k == 0 {
		t.Error("the label colour on the light window background at 13 pt draws as shipped")
	}
}

// TestAStemFamilyShapesInItsOwnFace: through both shapers, Roboto Regular
// upright in a stem family shapes in a face of its own, with the advances
// of the face as shipped. Roboto Medium in a stem family shapes in the face
// as shipped.
func TestAStemFamilyShapesInItsOwnFace(t *testing.T) {
	typ := DefaultTypography
	const px = 13
	shape := func(sh *text.Shaper, f font.Font) []text.Glyph {
		sh.LayoutString(text.Parameters{Font: f, PxPerEm: fixed.I(px), MaxLines: 1, MaxWidth: 10000}, "Vibrant Gio")
		var gs []text.Glyph
		for {
			g, ok := sh.NextGlyph()
			if !ok {
				return gs
			}
			gs = append(gs, g)
		}
	}
	for _, k := range []int{stemStepMin, 1, stemStepMax} {
		in := func(f font.Font) font.Font {
			f.Typeface = stemTypeface(k)
			return f
		}
		for name, sh := range map[string]*text.Shaper{"Shaper": typ.Shaper(), "DeterministicShaper": typ.DeterministicShaper()} {
			shipped, widened := shape(sh, roboto.RegularNormal), shape(sh, in(roboto.RegularNormal))
			if len(shipped) != len(widened) || len(shipped) == 0 {
				t.Fatalf("%s: %d glyphs as shipped, %d in the family of step %d", name, len(shipped), len(widened), k)
			}
			for i := range shipped {
				if shipped[i].ID == widened[i].ID {
					t.Errorf("%s, step %d: glyph %d shapes in the face as shipped", name, k, i)
				}
				if shipped[i].Advance != widened[i].Advance {
					t.Errorf("%s, step %d: glyph %d advances %v, want %v as shipped", name, k, i, widened[i].Advance, shipped[i].Advance)
				}
			}
			medium, mediumIn := shape(sh, roboto.RegularMedium), shape(sh, in(roboto.RegularMedium))
			for i := range medium {
				if medium[i].ID != mediumIn[i].ID {
					t.Errorf("%s, step %d: Roboto Medium glyph %d shapes in another face", name, k, i)
				}
			}
		}
	}
}
