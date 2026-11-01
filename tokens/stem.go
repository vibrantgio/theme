package tokens

import (
	"fmt"
	"image/color"
	"math"
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
// em, in total over the two sides of the stem. Zero is the face as shipped,
// and a negative value makes the stems narrower. No advance changes, so a
// text's width does not follow its stem weight. The other Roboto faces
// draw as shipped at every stem weight.
type StemWeight float32

// platformStems are the platform's lowercase stem, as stored, in px at 1x,
// for text darker than its fill and for text lighter than it. The rows at
// 10 to 13 pt are MEASURED in the reference's section "What the platform's
// small text measures at 1x", from Finder's file names. The 28 pt row is
// MEASURED in the reference's section "What the stem measures at the title
// sizes and at mid contrast", from Mail's empty-state line in the tertiary
// label colour. The contrast between the text and its fill does not move
// the platform's stem, so one row holds for every colour at its size.
var platformStems = []struct {
	size            float64
	darker, lighter float64
}{
	{10, 0.97, 1.16},
	{11, 1.06, 1.29},
	{12, 1.17, 1.40},
	{13, 1.26, 1.52},
	{28, 2.15, 2.44},
}

// platformStem returns the platform's stem as stored, in px at 1x, for text
// at size pt that is lighter than its fill or darker than it. From 10 to
// 13 pt it is the measured rows, joined by straight lines. Above 13 pt it
// is the line through the 13 pt and 28 pt rows. That line is DERIVED: no
// stored capture checks a size between the two, and the line gives about
// 1.08 px at 10 pt where 0.97 px is measured, so it holds above 13 pt
// only. Under 10 pt the stem keeps the 10 pt row's share of the size.
func platformStem(size float64, lighter bool) float64 {
	at := func(i int) float64 {
		if lighter {
			return platformStems[i].lighter
		}
		return platformStems[i].darker
	}
	if size <= platformStems[0].size {
		return at(0) * size / platformStems[0].size
	}
	i := 1
	for i < len(platformStems)-1 && size > platformStems[i].size {
		i++
	}
	s0, s1 := platformStems[i-1].size, platformStems[i].size
	return at(i-1) + (size-s0)*(at(i)-at(i-1))/(s1-s0)
}

// shippedStem is the stem of Roboto Regular upright as shipped, in em. The
// reference reads the library's dark lines linearised at 0.0907 em at every
// size, in its section "What the stem measures at the title sizes and at
// mid contrast".
const shippedStem = 0.0907

// The stem weights the shaper holds: the multiples of stemStep from
// stemStepMin to stemStepMax. Text whose stem weight falls outside them
// draws at the nearest end. The ends hold the platform's stem from 8 pt to
// about 40 pt for every colour on its fill.
const (
	stemStep    = 0.005
	stemStepMin = -3
	stemStepMax = 8
)

// minDrawnStem is the narrowest stem the library draws, in px at 1x. A
// stem drawn narrower than 1 px covers at most its own width of any pixel,
// so its darkest pixel stays under the text colour at every offset. The
// platform's lighter text at 10 pt reaches the text colour. Measured as
// stored, Roboto in the dark scheme's label colour on its control
// background at 10 pt reaches 95 % of the way to the text colour drawn at
// 0.907 px, and 94 % drawn at 0.882 px.
const minDrawnStem = 0.9

// blendCorrection returns the amount, in px, by which a straight stem in fg
// on fill reads wider as stored than its drawn width. Gio blends the text
// onto its fill in linear light, and the platform blends in encoded sRGB.
// A stored reading of a pixel at the share a of the way from the fill to
// the text is therefore not a: it is the encoded blend's share, per
// channel. Each of a stem's two edge pixels holds a share that is spread
// evenly from 0 to 1 over the stem's horizontal offsets. So the correction
// is twice the mean of the stored share less a over a from 0 to 1, with the
// mean over the channels in which fg and fill differ. It is negative for
// text darker than its fill and positive for text lighter than it. It
// grows with the contrast between them.
//
// The model gives the stored stems of the library's goldens to 0.03 px,
// read the way the reference reads them, for the label colour and the
// secondary label colour in both schemes and for white on the accent.
func blendCorrection(fg, fill color.NRGBA) float64 {
	key := [2]color.NRGBA{fg, fill}
	if v, ok := blendCorrections.Load(key); ok {
		return v.(float64)
	}
	f := [3]uint8{fill.R, fill.G, fill.B}
	t := [3]uint8{fg.R, fg.G, fg.B}
	const samples = 64
	sum, channels := 0.0, 0
	for c := range 3 {
		if f[c] == t[c] {
			continue
		}
		channels++
		ef, et := float64(f[c])/255, float64(t[c])/255
		lf, lt := vgcolor.LinearFromSRGB(ef), vgcolor.LinearFromSRGB(et)
		for k := range samples {
			a := (float64(k) + 0.5) / samples
			stored := (ef - vgcolor.SRGBFromLinear(lf+(lt-lf)*a)) / (ef - et)
			sum += stored - a
		}
	}
	v := 0.0
	if channels > 0 {
		v = 2 * sum / float64(channels*samples)
	}
	blendCorrections.Store(key, v)
	return v
}

// blendCorrections holds each blendCorrection the process has computed,
// keyed by the text colour and the fill.
var blendCorrections sync.Map

// TextStemWeight returns the stem weight of text at size pt in fg on fill.
// Text whose luminance is below the fill's is darker text, and any other
// text is lighter text. fg is flattened onto fill first, so a translucent
// label colour compares as it is drawn, and fill is taken as opaque.
//
// The platform's stem follows the size and whether the text is darker or
// lighter than its fill, and not the contrast between them. The library's
// stem as stored also moves with the contrast, because Gio blends in linear
// light. The stem weight is the amount that makes the library's stem as
// stored read the platform's: the platform's stem less the blend
// correction, over the size, less the stem as shipped. It is rounded to the
// nearest stem weight the shaper holds, and it is never less than the
// stem weight that draws the stem minDrawnStem wide.
func TextStemWeight(size float32, fg, fill color.NRGBA) StemWeight {
	return StemWeight(stemStepOf(size, fg, fill)) * stemStep
}

// stemStepOf is TextStemWeight as a multiple of stemStep.
func stemStepOf(size float32, fg, fill color.NRGBA) int {
	if size <= 0 {
		return 0
	}
	fill.A = 0xff
	fg = vgcolor.Flatten(fg, fill)
	lighter := vgcolor.RelativeLuminance(fg) >= vgcolor.RelativeLuminance(fill)
	s := float64(size)
	w := (platformStem(s, lighter)-blendCorrection(fg, fill))/s - shippedStem
	k := max(int(math.Round(w/stemStep)), int(math.Ceil((minDrawnStem/s-shippedStem)/stemStep)))
	return min(max(k, stemStepMin), stemStepMax)
}

// stemTypeface returns the key the shaper holds Roboto under at the stem
// weight k × stemStep. It is the shaper's key and not a typeface. No caller
// chooses it, [Typography.Faces] does not hold it, and no text on screen
// names it. Only [StemFont] sets it. The family "Roboto" holds the faces as
// shipped, which is the stem weight zero.
func stemTypeface(k int) font.Typeface {
	return font.Typeface(fmt.Sprintf("Roboto Stem %+d", k))
}

// StemFont returns f as the font that draws text at size pt in fg on fill
// at its stem weight, [TextStemWeight]. Only Roboto, named or as the empty
// typeface, has a family per stem weight; any other typeface returns f
// unchanged. The shapers of [Typography] hold every family.
//
// The empty typeface is the collection's default family, its first faces.
// Every library collection starts with Roboto: [DefaultTypography] does, and
// [Typography.WithFaces] appends. A collection that starts with another
// family must name Roboto in its roles, or its text draws in Roboto.
func StemFont(f font.Font, size float32, fg, fill color.NRGBA) font.Font {
	if f.Typeface != roboto.RegularNormal.Typeface && f.Typeface != "" {
		return f
	}
	if k := stemStepOf(size, fg, fill); k != 0 {
		f.Typeface = stemTypeface(k)
	}
	return f
}

// Dark reports whether p is a dark scheme's set: its label is light.
func (p PlatformColors) Dark() bool {
	return int(p.Label.R)+int(p.Label.G)+int(p.Label.B) > 3*127
}

// stemFaces returns faces followed by one family per stem weight the
// shaper holds, but zero: a copy of every Roboto entry under that stem
// weight's key, with the face of Roboto Regular upright at that stem
// weight. The other copies share the faces as shipped. Gio's text.NewShaper
// reads every face of its collection when it builds the shaper, so the
// first shaper of a process builds every widened face.
func stemFaces(faces []font.FontFace) []font.FontFace {
	out := make([]font.FontFace, 0, len(faces)*(stemStepMax-stemStepMin+1))
	out = append(out, faces...)
	for k := stemStepMin; k <= stemStepMax; k++ {
		if k == 0 {
			continue
		}
		for _, f := range faces {
			if f.Font.Typeface != roboto.RegularNormal.Typeface {
				continue
			}
			if f.Font == roboto.RegularNormal {
				f.Face = stemFace(k)
			}
			f.Font.Typeface = stemTypeface(k)
			out = append(out, f)
		}
	}
	return out
}

var (
	stemFaceMu sync.Mutex
	stemFaceAt = map[int]font.Face{}
)

// stemFace returns Roboto Regular upright at the stem weight k × stemStep,
// parsed once per process for each k.
func stemFace(k int) font.Face {
	stemFaceMu.Lock()
	defer stemFaceMu.Unlock()
	if f, ok := stemFaceAt[k]; ok {
		return f
	}
	w := float64(k) * stemStep
	ttf, err := widen.Sideways(normal.TTF, w)
	if err != nil {
		panic(fmt.Sprintf("tokens: Roboto Regular at stem weight %v: %v", w, err))
	}
	f, err := opentype.Parse(ttf)
	if err != nil {
		panic(fmt.Sprintf("tokens: Roboto Regular at stem weight %v: %v", w, err))
	}
	stemFaceAt[k] = f
	return f
}
