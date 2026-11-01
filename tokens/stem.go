package tokens

import (
	"fmt"
	"image/color"
	"sync"

	"gioui.org/font"
	"gioui.org/font/opentype"
	"github.com/vibrantgio/font/roboto"
	"github.com/vibrantgio/font/roboto/regular/normal"

	vgcolor "github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/internal/widen"
)

// StemWeight is the stem weight typography role text is drawn at. Its value
// is the width the library adds to each stem of Roboto Regular upright, in
// em, in total over the two sides of the stem. Zero is the face as shipped.
// The other Roboto faces draw as shipped at every stem weight.
type StemWeight float32

// The stem weight of text darker than the fill it stands on, and of text
// lighter than it. Each value makes the mean stem of "Vibrant Gio" over four
// horizontal offsets, 0, ¼, ½ and ¾ px, the platform's at 13 pt and 10 pt.
// The reference's section "What the same text measures in a platform window
// and a library window" holds the platform's rows and the library's.
//
// Each amount was measured at 13 pt and at 10 pt, and it is the same
// amount in em at both sizes. Darker text 0.031 em: 13 pt 1.29 px against the platform's 1.26, and 10 pt
// 0.92 px against 0.97. Lighter text 0 em: the face as shipped gives 13 pt
// 1.49 px against 1.51, and 10 pt 1.20 px against 1.17. The white title of
// the default button reads 1.72 px in both schemes. That is the lighter
// text's weight: the a in "Format" reads 1.70 px in dark and 1.45 px in light.
const (
	StemWeightDarkerText  StemWeight = 0.031
	StemWeightLighterText StemWeight = 0
)

// darkerTextTypeface is the key the shaper holds Roboto under for text darker
// than its fill: every Roboto face, with Regular upright at
// [StemWeightDarkerText]. It is the shaper's key and not a typeface. No
// caller chooses it, [Typography.Faces] does not hold it, and no text on
// screen names it. Only [StemFont] sets it. The family "Roboto" holds the
// faces as shipped, which is the stem weight of lighter text.
const darkerTextTypeface = "Roboto Darker Text"

// TextStemWeight returns the stem weight of text in fg on fill. Text whose
// luminance is below the fill's is darker text. fg is flattened onto fill
// first, so a translucent label colour compares as it is drawn. Text at the
// fill's own luminance is lighter text. fill is taken as opaque.
func TextStemWeight(fg, fill color.NRGBA) StemWeight {
	fill.A = 0xff
	if vgcolor.RelativeLuminance(vgcolor.Flatten(fg, fill)) < vgcolor.RelativeLuminance(fill) {
		return StemWeightDarkerText
	}
	return StemWeightLighterText
}

// StemFont returns f as the font that draws text in fg on fill at its stem
// weight, [TextStemWeight]. Only Roboto, named or as the empty typeface,
// has a darker text family; any other typeface returns f unchanged. The
// shapers of [Typography] hold both families.
//
// The empty typeface is the collection's default family, its first faces.
// Every library collection starts with Roboto: [DefaultTypography] does, and
// [Typography.WithFaces] appends. A collection that starts with another
// family must name Roboto in its roles, or its darker text draws in Roboto.
func StemFont(f font.Font, fg, fill color.NRGBA) font.Font {
	if f.Typeface != roboto.RegularNormal.Typeface && f.Typeface != "" {
		return f
	}
	if TextStemWeight(fg, fill) == StemWeightDarkerText {
		f.Typeface = darkerTextTypeface
	}
	return f
}

// Dark reports whether p is a dark scheme's set: its label is light.
func (p PlatformColors) Dark() bool {
	return int(p.Label.R)+int(p.Label.G)+int(p.Label.B) > 3*127
}

// stemFaces returns faces followed by the darker text family: a copy of
// every Roboto entry under [darkerTextTypeface], with the face of Roboto
// Regular upright at [StemWeightDarkerText]. The other copies share the
// faces as shipped.
func stemFaces(faces []font.FontFace) []font.FontFace {
	out := make([]font.FontFace, 0, 2*len(faces))
	out = append(out, faces...)
	for _, f := range faces {
		if f.Font.Typeface != roboto.RegularNormal.Typeface {
			continue
		}
		if f.Font == roboto.RegularNormal {
			f.Face = stemFace(StemWeightDarkerText)
		}
		f.Font.Typeface = darkerTextTypeface
		out = append(out, f)
	}
	return out
}

var (
	stemFaceMu sync.Mutex
	stemFaceAt = map[StemWeight]font.Face{}
)

// stemFace returns Roboto Regular upright at the stem weight w, parsed once
// per process for each w.
func stemFace(w StemWeight) font.Face {
	stemFaceMu.Lock()
	defer stemFaceMu.Unlock()
	if f, ok := stemFaceAt[w]; ok {
		return f
	}
	ttf, err := widen.Sideways(normal.TTF, float64(w))
	if err != nil {
		panic(fmt.Sprintf("tokens: Roboto Regular at stem weight %v: %v", w, err))
	}
	f, err := opentype.Parse(ttf)
	if err != nil {
		panic(fmt.Sprintf("tokens: Roboto Regular at stem weight %v: %v", w, err))
	}
	stemFaceAt[w] = f
	return f
}
