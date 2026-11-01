package tokens

import (
	"image/color"
	"math"
	"slices"
	"testing"

	"gioui.org/font"
	"gioui.org/text"
	"golang.org/x/image/math/fixed"

	"github.com/vibrantgio/font/roboto"

	vgcolor "github.com/vibrantgio/theme/color"
)

// TestEachStemFamilyWidensOnlyItsHeldFaces: the shaper's collection is the
// faces as shipped, then a copy of every Roboto entry under each stem
// weight the shaper holds but zero. Of the copies, Roboto Regular upright
// has a face of its own at each step of the regular range, and Roboto Bold
// upright at each step of boldSteps. Every other copy is the face as
// shipped.
func TestEachStemFamilyWidensOnlyItsHeldFaces(t *testing.T) {
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
	steps := heldSteps()
	copies := got[len(faces):]
	if len(copies) != len(steps)*len(shipped) {
		t.Fatalf("%d stem faces, want one per Roboto face for each of %d stem weights, %d", len(copies), len(steps), len(steps)*len(shipped))
	}
	stepOf := map[font.Typeface]int{}
	for _, k := range steps {
		stepOf[stemTypeface(k)] = k
	}
	seen := map[font.Face]bool{}
	widened := 0
	for _, c := range copies {
		k, ok := stepOf[c.Font.Typeface]
		if !ok {
			t.Fatalf("%v: no held stem weight has this key", c.Font)
		}
		f := c.Font
		f.Typeface = roboto.RegularNormal.Typeface
		want := f == roboto.RegularNormal && k >= stemStepMin && k <= stemStepMax ||
			f == roboto.RegularBold && slices.Contains(boldSteps, k)
		if got := c.Face != shipped[f]; got != want {
			t.Errorf("%v: widened %v, want %v", c.Font, got, want)
		}
		if want {
			if seen[c.Face] {
				t.Errorf("%v: shares its widened face with another step", c.Font)
			}
			seen[c.Face] = true
			widened++
		}
	}
	if want := stemStepMax - stemStepMin + len(boldSteps) - 1; widened != want {
		t.Errorf("%d widened faces, want %d", widened, want)
	}
}

// TestTheBoldStemFollowsTheWeight: Roboto Bold upright draws weight 600
// and weight 700, and each weight has its own platform rows. For each
// role with a row, the stem as stored that the step gives, the stem as
// shipped plus the step at the size, plus the blend correction, is the
// platform's stem to half a step. The shaper holds those steps and zero,
// and no other. The Headline lighter than its fill has no platform row
// whose weight is read, so it draws as shipped.
func TestTheBoldStemFollowsTheWeight(t *testing.T) {
	held := []int{0}
	for _, tg := range boldTargets() {
		w := FontWeight(tg.role.Weight)
		k := boldStepOf(w, tg.role.Size, tg.fg, tg.fill)
		if k != boldStepFor(w, tg.role.Size, tg.fg, tg.fill) {
			t.Errorf("weight %v at %v pt: the held step %d is not the target's step", tg.role.Weight, tg.role.Size, k)
		}
		if !slices.Contains(held, k) {
			held = append(held, k)
		}
		flat := vgcolor.Flatten(tg.fg, tg.fill)
		lighter := vgcolor.RelativeLuminance(flat) >= vgcolor.RelativeLuminance(tg.fill)
		size := float64(tg.role.Size)
		want, ok := boldPlatformStem(w, size, lighter)
		if !ok {
			if k != 0 {
				t.Errorf("weight %v at %v pt, lighter: step %d, want zero with no platform row", tg.role.Weight, size, k)
			}
			continue
		}
		stored := (shippedBoldStem+float64(k)*stemStep)*size + blendCorrection(flat, tg.fill)
		if math.Abs(stored-want) > stemStep*size/2+1e-6 {
			t.Errorf("weight %v at %v pt: step %d gives %.3f px, want the platform's %.3f px", tg.role.Weight, size, k, stored, want)
		}
	}
	slices.Sort(held)
	if !slices.Equal(held, boldSteps) {
		t.Errorf("the shaper holds the bold steps %v, want the targets' %v", boldSteps, held)
	}
	l := PlatformLight
	if a, b := boldStepOf(font.SemiBold, 13, l.Label, l.WindowBackground), boldStepOf(font.Bold, 13, l.Label, l.WindowBackground); a >= b {
		t.Errorf("13 pt in the label colour: weight 600 at step %d, weight 700 at step %d", a, b)
	}
	for _, r := range DefaultTypography.DocumentHeadings {
		for _, p := range []PlatformColors{PlatformLight, PlatformDark} {
			if k := boldStepOf(FontWeight(r.Weight), r.Size, p.Label, p.WindowBackground); !slices.Contains(boldSteps, k) {
				t.Errorf("a document heading at %v pt: step %d, not held", r.Size, k)
			}
		}
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
// correction, is the platform's stem to half a step. The one exception is a
// stem that the floor raises: it is drawn minDrawnStem wide to one step,
// and no stem is drawn narrower.
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
			drawn := (shippedStem + float64(w)) * float64(size)
			stored := drawn + blendCorrection(flat, p.fill)
			want := platformStem(float64(size), p.lighter)
			if drawn < minDrawnStem-1e-6 {
				t.Errorf("%s at %v pt: weight %v draws the stem %.3f px wide, under %v px", p.name, size, w, drawn, minDrawnStem)
			}
			if raised := drawn-stemStep*float64(size) < minDrawnStem-1e-6; raised && stored > want {
				continue
			}
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
// stem weight zero in Roboto as shipped. Roboto Bold upright takes the
// stem weight of its own weight's rows. Text in any other typeface keeps
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
		{"Roboto weight 600, darker text", font.Font{Typeface: "Roboto", Weight: font.SemiBold}, 13, light.Label, light.WindowBackground},
		{"Roboto weight 700, darker text", roboto.RegularBold, 13, light.Label, light.WindowBackground},
		{"Roboto weight 700, lighter text", roboto.RegularBold, 13, dark.Label, dark.WindowBackground},
		{"Roboto Bold Italic", roboto.ItalicBold, 13, light.Label, light.WindowBackground},
	} {
		got := StemFont(tc.f, tc.size, tc.fg, tc.fill)
		want := tc.f.Typeface
		k := stemStepOf(tc.size, tc.fg, tc.fill)
		if drawsBold(tc.f) {
			k = boldStepOf(tc.f.Weight, tc.size, tc.fg, tc.fill)
		}
		if k != 0 {
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
// of the face as shipped. Roboto Bold upright shapes in a face of its own
// in a family of boldSteps, with the advances as shipped, and as shipped in
// any other family. Roboto Medium in a stem family shapes in the face as
// shipped.
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
	for _, k := range []int{stemStepMin, 1, 3, stemStepMax} {
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
			bold, boldIn := shape(sh, roboto.RegularBold), shape(sh, in(roboto.RegularBold))
			for i := range bold {
				if (bold[i].ID != boldIn[i].ID) != slices.Contains(boldSteps, k) {
					t.Errorf("%s, step %d: Roboto Bold glyph %d widened %v, want %v", name, k, i, bold[i].ID != boldIn[i].ID, slices.Contains(boldSteps, k))
				}
				if bold[i].Advance != boldIn[i].Advance {
					t.Errorf("%s, step %d: Roboto Bold glyph %d advances %v, want %v as shipped", name, k, i, boldIn[i].Advance, bold[i].Advance)
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
