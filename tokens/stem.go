package tokens

import (
	"fmt"
	"image/color"
	"math"
	"slices"
	"sync"

	"gioui.org/font"
	"gioui.org/font/opentype"
	"github.com/vibrantgio/font/roboto"
	"github.com/vibrantgio/font/roboto/regular/bold"
	"github.com/vibrantgio/font/roboto/regular/normal"
	"github.com/vibrantgio/font/robotomono"
	mononormal "github.com/vibrantgio/font/robotomono/regular/normal"

	vgcolor "github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/internal/widen"
)

// StemWeight is the stem weight typography role text is drawn at. Its value
// is the width the library adds to each stem of Roboto Regular upright, in
// em, in total over the two sides of the stem. Zero is the face as shipped,
// and a negative value makes the stems narrower. No advance changes, so a
// text's width does not follow its stem weight. Roboto Bold upright and
// Roboto Mono Regular upright have stem weights of their own, under
// [StemFont]. The other Roboto and Roboto Mono faces draw as shipped at
// every stem weight.
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

// The stem weights the shaper holds for Roboto Regular upright: the
// multiples of stemStep from stemStepMin to stemStepMax. Text whose stem weight falls outside them
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

// The platform's stems for semibold and bold text, as stored, in px at 1x,
// for text darker than its fill and for text lighter than it. Roboto ships
// no semibold face. Weight 600 and weight 700 both draw Roboto Bold: go-text
// takes the next heavier face for a weight over 500, and that face is Bold.
// The platform draws the two weights apart, so the font's weight, and not
// the face it draws, selects the rows. The reference's section "What the
// stem measures for the semibold, bold and monospaced faces" holds them:
//
//   - Weight 600 at 13 pt is System Settings' group title, MEASURED.
//   - Weight 600 at 11 pt is Mail's sidebar section heading, MEASURED. It is
//     the same in both schemes.
//   - Weight 700 at 13 pt is TextEdit's light window title, MEASURED. No
//     stored capture holds bold text lighter than its fill whose weight is
//     read, so lighter text at weight 700 has no row.
var boldStems = map[font.Weight][]boldStemRow{
	font.SemiBold: {{size: 11, darker: 1.52, lighter: 1.52, hasLighter: true}, {size: 13, darker: 1.73, lighter: 1.99, hasLighter: true}},
	font.Bold:     {{size: 13, darker: 2.02}},
}

type boldStemRow struct {
	size            float64
	darker, lighter float64
	hasLighter      bool
}

// boldPlatformStem returns the platform's stem as stored, in px at 1x, for
// text in weight w at size pt. Between two rows it is the straight line
// through them. Outside the rows the stem keeps the nearest row's share of
// the size. It reports false for a weight or a side of the fill with no
// row.
func boldPlatformStem(w font.Weight, size float64, lighter bool) (float64, bool) {
	rows := boldStems[w]
	if len(rows) == 0 || lighter && !rows[0].hasLighter {
		return 0, false
	}
	at := func(i int) float64 {
		if lighter {
			return rows[i].lighter
		}
		return rows[i].darker
	}
	last := len(rows) - 1
	switch {
	case size <= rows[0].size:
		return at(0) * size / rows[0].size, true
	case size >= rows[last].size:
		return at(last) * size / rows[last].size, true
	}
	i := 1
	for size > rows[i].size {
		i++
	}
	s0, s1 := rows[i-1].size, rows[i].size
	return at(i-1) + (size-s0)*(at(i)-at(i-1))/(s1-s0), true
}

// shippedBoldStem is the stem of Roboto Bold upright as shipped, in em,
// MEASURED from its l, n and i rendered at 1000 px by FreeType.
const shippedBoldStem = 0.141

// boldTargets are the roles whose bold text has a platform row, in the
// text colour and on the fill each is measured on. The section heading and
// the headline use the label colour on the window. The small section
// heading uses the secondary label colour on the sidebar.
func boldTargets() []boldTarget {
	t := DefaultTypography
	var out []boldTarget
	for _, p := range []PlatformColors{PlatformLight, PlatformDark} {
		out = append(out,
			boldTarget{t.SectionHeading, p.Label, p.WindowBackground},
			boldTarget{t.Headline, p.Label, p.WindowBackground},
			boldTarget{t.SmallSectionHeading, p.SecondaryLabel, p.SidebarMaterial},
		)
	}
	return out
}

type boldTarget struct {
	role     TextStyle
	fg, fill color.NRGBA
}

// boldSteps are the stem weights of Roboto Bold upright that the shaper
// holds, as multiples of stemStep: the steps of the targets, and zero, the
// face as shipped. Each held face takes about 3.1 MB in both shapers, so
// the shaper holds no other step.
var boldSteps = func() []int {
	held := []int{0}
	for _, tg := range boldTargets() {
		k := boldStepFor(FontWeight(tg.role.Weight), tg.role.Size, tg.fg, tg.fill)
		if !slices.Contains(held, k) {
			held = append(held, k)
		}
	}
	slices.Sort(held)
	return held
}()

// boldStepFor is the stem weight of Roboto Bold upright text in weight w
// at size pt in fg on fill, as a multiple of stemStep, before it is held:
// the platform's stem less the blend correction, over the size, less the
// stem as shipped. Text with no platform row has the step zero.
func boldStepFor(w font.Weight, size float32, fg, fill color.NRGBA) int {
	if size <= 0 {
		return 0
	}
	fill.A = 0xff
	fg = vgcolor.Flatten(fg, fill)
	lighter := vgcolor.RelativeLuminance(fg) >= vgcolor.RelativeLuminance(fill)
	s := float64(size)
	stem, ok := boldPlatformStem(w, s, lighter)
	if !ok {
		return 0
	}
	return int(math.Round(((stem-blendCorrection(fg, fill))/s - shippedBoldStem) / stemStep))
}

// boldStepOf is the held stem weight of Roboto Bold upright text: the
// held step nearest boldStepFor. Bold text in a role with no target, as
// bold text in reading text and the document headings, draws at it.
func boldStepOf(w font.Weight, size float32, fg, fill color.NRGBA) int {
	return nearestHeld(boldSteps, boldStepFor(w, size, fg, fill))
}

func abs(k int) int {
	if k < 0 {
		return -k
	}
	return k
}

// The platform's stem for monospaced text, in em, for text darker than its
// fill and for text lighter than it. The reference's section "What the
// stem measures for the semibold, bold and monospaced faces" holds the
// rows, MEASURED from TextEdit's 11 pt plain-text body: 1.11 px in black on
// white and 1.34 px in white on the dark window. No stored capture holds
// monospaced text at another size. The platform's regular face keeps one
// stem em from 10 to 13 pt, so the library keeps the monospaced em at every
// size. At the Code role's 13 pt the stem is DERIVED: 1.31 px and 1.59 px.
const (
	monoStemDarker  = 0.101
	monoStemLighter = 0.122
)

// shippedMonoStem is the stem of Roboto Mono Regular upright as shipped,
// in em, MEASURED from its l, n and i rendered at 1000 px by FreeType.
const shippedMonoStem = 0.090

// monoTargets are the text colours and fills that draw code in the Code
// role, in both schemes: the label colour on the window, and the text
// colour on the fill of a fenced code block and of an inline code span.
// That fill is the alternating content background over the text
// background.
func monoTargets() []monoTarget {
	var out []monoTarget
	for _, p := range []PlatformColors{PlatformLight, PlatformDark} {
		fence := vgcolor.Flatten(p.AlternatingContentBackground, p.TextBackground)
		out = append(out,
			monoTarget{p.Label, p.WindowBackground},
			monoTarget{p.Text, fence},
		)
	}
	return out
}

type monoTarget struct {
	fg, fill color.NRGBA
}

// monoSteps are the stem weights of Roboto Mono Regular upright that the
// shaper holds, as multiples of stemStep: the steps of the targets at the
// Code role's size, and zero, the face as shipped. Each held face takes
// heap in both shapers, so the shaper holds no other step.
var monoSteps = func() []int {
	held := []int{0}
	for _, tg := range monoTargets() {
		k := monoStepFor(DefaultTypography.Code.Size, tg.fg, tg.fill)
		if !slices.Contains(held, k) {
			held = append(held, k)
		}
	}
	slices.Sort(held)
	return held
}()

// monoStepFor is the stem weight of Roboto Mono Regular upright text at
// size pt in fg on fill, as a multiple of stemStep, before it is held: the
// platform's stem less the blend correction, over the size, less the stem
// as shipped.
func monoStepFor(size float32, fg, fill color.NRGBA) int {
	if size <= 0 {
		return 0
	}
	fill.A = 0xff
	fg = vgcolor.Flatten(fg, fill)
	em := monoStemDarker
	if vgcolor.RelativeLuminance(fg) >= vgcolor.RelativeLuminance(fill) {
		em = monoStemLighter
	}
	s := float64(size)
	return int(math.Round(((em*s-blendCorrection(fg, fill))/s - shippedMonoStem) / stemStep))
}

// monoStepOf is the held stem weight of Roboto Mono Regular upright text:
// the held step nearest monoStepFor. Code in a colour or on a fill with no
// target, as highlighted code, draws at it.
func monoStepOf(size float32, fg, fill color.NRGBA) int {
	return nearestHeld(monoSteps, monoStepFor(size, fg, fill))
}

// nearestHeld returns the step of held nearest k. Of two steps at the same
// distance it returns the one nearer zero.
func nearestHeld(held []int, k int) int {
	best := 0
	for _, h := range held {
		if d, bd := abs(h-k), abs(best-k); d < bd || d == bd && abs(h) < abs(best) {
			best = h
		}
	}
	return best
}

// drawsMonoRegular reports whether f draws Roboto Mono Regular upright.
// Roboto Mono ships Regular and Bold, and go-text takes Bold for a weight
// over 500 only.
func drawsMonoRegular(f font.Font) bool {
	return f.Typeface == robotomono.RegularNormal.Typeface && f.Style == font.Regular && f.Weight < font.SemiBold
}

// monoTypeface returns the key the shaper holds Roboto Mono under at the
// stem weight k × stemStep. It is the shaper's key and not a typeface, as
// for [stemTypeface].
func monoTypeface(k int) font.Typeface {
	return font.Typeface(fmt.Sprintf("Roboto Mono Stem %+d", k))
}

// drawsBold reports whether f, in Roboto, draws Roboto Bold upright.
func drawsBold(f font.Font) bool {
	return f.Style == font.Regular && (f.Weight == font.SemiBold || f.Weight == font.Bold)
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
// at its stem weight. Text in Roboto Bold upright, weight 600 or 700, has
// the stem weight of its weight's platform rows. Other text in Roboto has
// [TextStemWeight]. Text in Roboto Mono Regular upright has the stem weight
// of the platform's monospaced row. Roboto Mono's bold and italic faces
// have no platform row and return f unchanged. JetBrains Mono, the other
// code face, and any other typeface return f unchanged: they draw as
// shipped. The shapers of [Typography] hold every family.
//
// The empty typeface is the collection's default family, its first faces.
// Every library collection starts with Roboto: [DefaultTypography] does, and
// [Typography.WithFaces] appends. A collection that starts with another
// family must name Roboto in its roles, or its text draws in Roboto.
func StemFont(f font.Font, size float32, fg, fill color.NRGBA) font.Font {
	if drawsMonoRegular(f) {
		if k := monoStepOf(size, fg, fill); k != 0 {
			f.Typeface = monoTypeface(k)
		}
		return f
	}
	if f.Typeface != roboto.RegularNormal.Typeface && f.Typeface != "" {
		return f
	}
	k := 0
	if drawsBold(f) {
		k = boldStepOf(f.Weight, size, fg, fill)
	} else {
		k = stemStepOf(size, fg, fill)
	}
	if k != 0 {
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
// weight's key. In a family, Roboto Regular upright is widened when the
// step is in the regular range, and Roboto Bold upright is widened when the
// step is in boldSteps. Then follows one family per step of monoSteps but
// zero: a copy of every Roboto Mono entry under that step's key, with Roboto
// Mono Regular upright widened. The other copies share the faces as
// shipped. Gio's
// text.NewShaper reads every face of its collection when it builds the
// shaper, so the first shaper of a process builds every widened face.
func stemFaces(faces []font.FontFace) []font.FontFace {
	steps := heldSteps()
	out := make([]font.FontFace, 0, len(faces)*(len(steps)+1))
	out = append(out, faces...)
	for _, k := range steps {
		for _, f := range faces {
			if f.Font.Typeface != roboto.RegularNormal.Typeface {
				continue
			}
			switch {
			case f.Font == roboto.RegularNormal && k >= stemStepMin && k <= stemStepMax:
				f.Face = stemFace(normal.TTF, k)
			case f.Font == roboto.RegularBold && slices.Contains(boldSteps, k):
				f.Face = stemFace(bold.TTF, k)
			}
			f.Font.Typeface = stemTypeface(k)
			out = append(out, f)
		}
	}
	for _, k := range monoSteps {
		if k == 0 {
			continue
		}
		for _, f := range faces {
			if f.Font.Typeface != robotomono.RegularNormal.Typeface {
				continue
			}
			if f.Font == robotomono.RegularNormal {
				f.Face = stemFace(mononormal.TTF, k)
			}
			f.Font.Typeface = monoTypeface(k)
			out = append(out, f)
		}
	}
	return out
}

// heldSteps returns every non-zero step a family exists for, in order: the
// regular range and boldSteps.
func heldSteps() []int {
	var steps []int
	for k := stemStepMin; k <= stemStepMax; k++ {
		if k != 0 {
			steps = append(steps, k)
		}
	}
	for _, k := range boldSteps {
		if k != 0 && !slices.Contains(steps, k) {
			steps = append(steps, k)
		}
	}
	slices.Sort(steps)
	return steps
}

type stemFaceKey struct {
	ttf *byte
	k   int
}

var (
	stemFaceMu sync.Mutex
	stemFaceAt = map[stemFaceKey]font.Face{}
)

// stemFace returns the face ttf at the stem weight k × stemStep, parsed
// once per process for each face and each k.
func stemFace(ttf []byte, k int) font.Face {
	stemFaceMu.Lock()
	defer stemFaceMu.Unlock()
	key := stemFaceKey{&ttf[0], k}
	if f, ok := stemFaceAt[key]; ok {
		return f
	}
	w := float64(k) * stemStep
	wide, err := widen.Sideways(ttf, w)
	if err != nil {
		panic(fmt.Sprintf("tokens: face at stem weight %v: %v", w, err))
	}
	f, err := opentype.Parse(wide)
	if err != nil {
		panic(fmt.Sprintf("tokens: face at stem weight %v: %v", w, err))
	}
	stemFaceAt[key] = f
	return f
}
