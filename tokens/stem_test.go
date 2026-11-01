package tokens

import (
	"image/color"
	"testing"

	"gioui.org/font"
	"gioui.org/text"
	"golang.org/x/image/math/fixed"

	"github.com/vibrantgio/font/roboto"
)

// TestTheDarkerTextFamilyWidensOnlyRobotoRegular: the shaper's collection
// is the faces as shipped, then a copy of every Roboto entry under the
// darker text family. Of the copies, only Roboto Regular upright has a face
// of its own.
func TestTheDarkerTextFamilyWidensOnlyRobotoRegular(t *testing.T) {
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
	copies := got[len(faces):]
	if len(copies) != len(shipped) {
		t.Fatalf("%d darker text faces, want one per Roboto face, %d", len(copies), len(shipped))
	}
	for _, c := range copies {
		if c.Font.Typeface != darkerTextTypeface {
			t.Fatalf("a copy is in %q, want %q", c.Font.Typeface, darkerTextTypeface)
		}
		f := c.Font
		f.Typeface = roboto.RegularNormal.Typeface
		widened := c.Face != shipped[f]
		if widened != (f == roboto.RegularNormal) {
			t.Errorf("%v: widened %v, want widened only for Roboto Regular upright", f, widened)
		}
	}
}

// TestEveryCopySharesTheShaper: copies of a typography hand back one shaper
// of each kind, and each holds both Roboto families.
func TestEveryCopySharesTheShaper(t *testing.T) {
	a, b := DefaultTypography, DefaultTypography
	if a.DeterministicShaper() != b.DeterministicShaper() || a.Shaper() != b.Shaper() {
		t.Error("two copies built two shapers")
	}
}

// TestTheStemWeightFollowsTheTextAgainstItsFill: the stem weight follows the
// text against the fill it stands on, in either scheme. The label colour on
// the window background is darker text in light and lighter text in dark.
// The white title on the default button is lighter text in both. The
// secondary label colour, a mid-grey, is darker text on the light window
// background and lighter text on the dark one, and so is one fixed grey.
func TestTheStemWeightFollowsTheTextAgainstItsFill(t *testing.T) {
	white := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	grey := color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}
	for _, tc := range []struct {
		name     string
		fg, fill color.NRGBA
		want     StemWeight
	}{
		{"light: label colour on the window background", PlatformLight.Label, PlatformLight.WindowBackground, StemWeightDarkerText},
		{"dark: label colour on the window background", PlatformDark.Label, PlatformDark.WindowBackground, StemWeightLighterText},
		{"light: white title on the default button", white, PlatformLight.DefaultButtonFill, StemWeightLighterText},
		{"dark: white title on the default button", white, PlatformDark.DefaultButtonFill, StemWeightLighterText},
		{"light: secondary label colour on the window background", PlatformLight.SecondaryLabel, PlatformLight.WindowBackground, StemWeightDarkerText},
		{"dark: secondary label colour on the window background", PlatformDark.SecondaryLabel, PlatformDark.WindowBackground, StemWeightLighterText},
		{"light: grey on the window background", grey, PlatformLight.WindowBackground, StemWeightDarkerText},
		{"dark: grey on the window background", grey, PlatformDark.WindowBackground, StemWeightLighterText},
		{"the fill's own colour", grey, grey, StemWeightLighterText},
	} {
		if got := TextStemWeight(tc.fg, tc.fill); got != tc.want {
			t.Errorf("%s: stem weight %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestStemFontNamesTheDarkerTextFamilyForRobotoOnly: darker text in Roboto,
// named or as the empty typeface, shapes in the darker text family. Lighter
// text, and text in any other typeface, keeps its font.
func TestStemFontNamesTheDarkerTextFamilyForRobotoOnly(t *testing.T) {
	light := PlatformLight
	dark := PlatformDark
	mono := font.Font{Typeface: CodeFaceRoboto}
	for _, tc := range []struct {
		name     string
		f        font.Font
		fg, fill color.NRGBA
		want     font.Typeface
	}{
		{"Roboto, darker text", roboto.RegularMedium, light.Label, light.WindowBackground, darkerTextTypeface},
		{"empty typeface, darker text", font.Font{}, light.Label, light.WindowBackground, darkerTextTypeface},
		{"Roboto, lighter text", roboto.RegularNormal, dark.Label, dark.WindowBackground, roboto.RegularNormal.Typeface},
		{"empty typeface, lighter text", font.Font{}, dark.Label, dark.WindowBackground, ""},
		{"Roboto Mono, darker text", mono, light.Label, light.WindowBackground, CodeFaceRoboto},
	} {
		got := StemFont(tc.f, tc.fg, tc.fill)
		if got.Typeface != tc.want {
			t.Errorf("%s: typeface %q, want %q", tc.name, got.Typeface, tc.want)
		}
		if got.Weight != tc.f.Weight || got.Style != tc.f.Style {
			t.Errorf("%s: %v changed its weight or style to %v", tc.name, tc.f, got)
		}
	}
}

// TestTheDarkerTextFamilyShapesInItsOwnFace: through both shapers, Roboto
// Regular upright in the darker text family shapes in a face of its own,
// with the advances of the face as shipped. Roboto Medium in the darker
// text family shapes in the face as shipped.
func TestTheDarkerTextFamilyShapesInItsOwnFace(t *testing.T) {
	typ := DefaultTypography
	shape := func(sh *text.Shaper, f font.Font) []text.Glyph {
		sh.LayoutString(text.Parameters{Font: f, PxPerEm: fixed.I(13), MaxLines: 1, MaxWidth: 10000}, "Vibrant Gio")
		var gs []text.Glyph
		for {
			g, ok := sh.NextGlyph()
			if !ok {
				return gs
			}
			gs = append(gs, g)
		}
	}
	darker := func(f font.Font) font.Font {
		f.Typeface = darkerTextTypeface
		return f
	}
	for name, sh := range map[string]*text.Shaper{"Shaper": typ.Shaper(), "DeterministicShaper": typ.DeterministicShaper()} {
		shipped, widened := shape(sh, roboto.RegularNormal), shape(sh, darker(roboto.RegularNormal))
		if len(shipped) != len(widened) || len(shipped) == 0 {
			t.Fatalf("%s: %d glyphs as shipped, %d in the darker text family", name, len(shipped), len(widened))
		}
		for i := range shipped {
			if shipped[i].ID == widened[i].ID {
				t.Errorf("%s: glyph %d shapes in the face as shipped in the darker text family", name, i)
			}
			if shipped[i].Advance != widened[i].Advance {
				t.Errorf("%s: glyph %d advances %v, want %v as shipped", name, i, widened[i].Advance, shipped[i].Advance)
			}
		}
		medium, mediumDarker := shape(sh, roboto.RegularMedium), shape(sh, darker(roboto.RegularMedium))
		for i := range medium {
			if medium[i].ID != mediumDarker[i].ID {
				t.Errorf("%s: Roboto Medium glyph %d shapes in another face in the darker text family", name, i)
			}
		}
	}
}
