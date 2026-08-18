package imageseed

import (
	"image"
	stdcolor "image/color"
	"math"
	"testing"

	"github.com/vibrantgio/theme/color"
)

// The fixtures are painted here rather than stored: an image whose palette
// is written in the test is an image whose expected answer is written in the
// test, and a photograph in a repository is a palette nobody can check.

// band is one horizontal stripe of a fixture: a flat colour and how many
// rows of the image it covers.
type band struct {
	c    stdcolor.NRGBA
	rows int
}

// bands paints a w-wide image out of the given stripes, flat colour each, so
// every distinct colour in the result is one the test named and every share
// is one the test can count. The image is as tall as the stripes together.
func bands(w int, parts ...band) *image.NRGBA {
	h := 0
	for _, p := range parts {
		h += p.rows
	}
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	y := 0
	for _, p := range parts {
		for ; y < len(img.Pix) && p.rows > 0; p.rows-- {
			for x := 0; x < w; x++ {
				img.SetNRGBA(x, y, p.c)
			}
			y++
		}
	}
	return img
}

// flat paints one colour over the whole image.
func flat(w, h int, c stdcolor.NRGBA) *image.NRGBA {
	return bands(w, band{c, h})
}

func rgb(r, g, b uint8) stdcolor.NRGBA { return stdcolor.NRGBA{R: r, G: g, B: b, A: 0xff} }

// Fixture colours, with the OKLCh chroma each measures at through the
// module's own converters — the axis the ranking weighs.
var (
	grey   = rgb(0x9a, 0x9a, 0x9a) // chroma 0.000
	beige  = rgb(0xd8, 0xcf, 0xc0) // chroma 0.023: drab, but not neutral
	orange = rgb(0xff, 0x6a, 0x00) // chroma 0.201
	blue   = rgb(0x00, 0x50, 0xd0) // chroma 0.209
	green  = rgb(0x2f, 0xb5, 0x4a) // chroma 0.186
	red    = rgb(0xe8, 0x11, 0x2d) // chroma 0.234
)

func chromaOf(c stdcolor.NRGBA) float64 {
	_, chroma, _ := color.OKLChFromNRGBA(c)
	return chroma
}

func colors(cs []Candidate) []stdcolor.NRGBA {
	out := make([]stdcolor.NRGBA, len(cs))
	for i, c := range cs {
		out[i] = c.Color
	}
	return out
}

func indexOf(cs []Candidate, want stdcolor.NRGBA) int {
	for i, c := range cs {
		if c.Color == want {
			return i
		}
	}
	return -1
}

// TestFixtureChromas pins what the fixture colours measure, so a failure in
// a ranking test below reads as a ranking defect rather than as a fixture
// that quietly stopped being vivid.
func TestFixtureChromas(t *testing.T) {
	for _, tc := range []struct {
		name string
		c    stdcolor.NRGBA
		want float64
	}{
		{"grey", grey, 0.000},
		{"beige", beige, 0.023},
		{"orange", orange, 0.201},
		{"blue", blue, 0.209},
		{"green", green, 0.186},
		{"red", red, 0.234},
	} {
		if got := chromaOf(tc.c); math.Abs(got-tc.want) > 0.002 {
			t.Errorf("%s chroma = %.3f, want %.3f", tc.name, got, tc.want)
		}
	}
}

// TestExtractRanksKnownPalette is the core fixture: an image of four flat
// colours in known proportions, extracted and asserted on rank AND on
// share. The three vivid stripes lead the grey that covers more of the
// image than all three together, and among themselves they fall in the
// order their prominence implies.
func TestExtractRanksKnownPalette(t *testing.T) {
	img := bands(64,
		band{grey, 60},
		band{orange, 20},
		band{blue, 12},
		band{green, 8},
	)
	got := Extract(img)
	want := []stdcolor.NRGBA{orange, blue, green, grey}
	if len(got) != len(want) {
		t.Fatalf("got %d candidates %v, want %d", len(got), colors(got), len(want))
	}
	for i, w := range want {
		if got[i].Color != w {
			t.Errorf("candidate %d = %v, want %v (full order %v)", i, got[i].Color, w, colors(got))
		}
	}
	shares := []float64{0.20, 0.12, 0.08, 0.60}
	for i, w := range shares {
		if math.Abs(got[i].Share-w) > 0.02 {
			t.Errorf("candidate %d share = %.3f, want %.2f", i, got[i].Share, w)
		}
	}
	if got[0].Weight <= got[3].Weight {
		t.Errorf("leading weight %.5f does not beat the grey majority's %.5f", got[0].Weight, got[3].Weight)
	}
}

// TestVividMinorityBeatsDullMajority is the ranking rule stated as an
// image: a tenth of the frame in a vivid red, the rest in a drab beige that
// is not grey enough for the floor to rescue it.
func TestVividMinorityBeatsDullMajority(t *testing.T) {
	img := bands(64, band{beige, 90}, band{red, 10})
	got := Extract(img)
	if len(got) < 2 {
		t.Fatalf("got %d candidates %v, want both colours", len(got), colors(got))
	}
	if got[0].Color != red {
		t.Errorf("leading candidate = %v, want the vivid minority %v", got[0].Color, red)
	}
	if got[1].Color != beige {
		t.Errorf("second candidate = %v, want the drab majority %v", got[1].Color, beige)
	}
	if got[1].Share <= got[0].Share {
		t.Errorf("shares %.2f/%.2f do not show the leader as the minority", got[0].Share, got[1].Share)
	}
}

// TestGreyscaleYieldsGreysByShare: nothing in the image has a hue, so the
// floor carries the ordering and the answer is the image's own greys, most
// common first — not noise, and not an invented colour.
func TestGreyscaleYieldsGreysByShare(t *testing.T) {
	dark, mid, light := rgb(0x22, 0x22, 0x22), rgb(0x88, 0x88, 0x88), rgb(0xdd, 0xdd, 0xdd)
	img := bands(64, band{mid, 50}, band{dark, 30}, band{light, 20})
	got := Extract(img)
	if len(got) != 3 {
		t.Fatalf("got %d candidates %v, want 3", len(got), colors(got))
	}
	for i, w := range []stdcolor.NRGBA{mid, dark, light} {
		if got[i].Color != w {
			t.Errorf("candidate %d = %v, want %v (full order %v)", i, got[i].Color, w, colors(got))
		}
	}
	for i, c := range got {
		if c.Chroma > greyChroma {
			t.Errorf("candidate %d %v chroma %.4f: a grey image yielded a coloured candidate", i, c.Color, c.Chroma)
		}
	}
}

// TestSingleColourImage: one colour in, one candidate out, at the whole
// share — no runners-up split off a flat field.
func TestSingleColourImage(t *testing.T) {
	got := Extract(flat(120, 90, blue))
	if len(got) != 1 {
		t.Fatalf("got %d candidates %v, want exactly 1", len(got), colors(got))
	}
	if got[0].Color != blue {
		t.Errorf("candidate = %v, want %v", got[0].Color, blue)
	}
	if got[0].Share != 1 {
		t.Errorf("share = %v, want 1", got[0].Share)
	}
}

// TestTinyImages: an image smaller than any sampling stride is read whole,
// down to the single pixel.
func TestTinyImages(t *testing.T) {
	one := Extract(flat(1, 1, red))
	if len(one) != 1 || one[0].Color != red {
		t.Errorf("1x1 image gave %v, want [%v]", colors(one), red)
	}
	two := Extract(bands(3, band{red, 1}, band{blue, 1}))
	if len(two) != 2 {
		t.Fatalf("3x2 image gave %d candidates %v, want 2", len(two), colors(two))
	}
	if indexOf(two, red) < 0 || indexOf(two, blue) < 0 {
		t.Errorf("3x2 image gave %v, want both %v and %v", colors(two), red, blue)
	}
}

// TestNoColourAtAll: nothing to cluster is no candidates, not a panic and
// not an invented black.
func TestNoColourAtAll(t *testing.T) {
	if got := Extract(nil); got != nil {
		t.Errorf("nil image gave %v, want nil", colors(got))
	}
	if got := Extract(image.NewNRGBA(image.Rectangle{})); got != nil {
		t.Errorf("empty image gave %v, want nil", colors(got))
	}
	clear := image.NewNRGBA(image.Rect(0, 0, 32, 32)) // every pixel alpha 0
	if got := Extract(clear); got != nil {
		t.Errorf("fully transparent image gave %v, want nil", colors(got))
	}
}

// TestSamplingIsSizeIndependent is the cost cap stated as behaviour: the
// same picture painted at eight megapixels and at four thousand yields the
// same candidates, because the large one is read on a stride rather than
// whole. A change that made extraction scan every pixel would still pass
// this; a change that made the ANSWER depend on the scan would not.
func TestSamplingIsSizeIndependent(t *testing.T) {
	small := bands(64, band{grey, 60}, band{orange, 20}, band{blue, 12}, band{green, 8})
	large := bands(2900, band{grey, 1800}, band{orange, 600}, band{blue, 360}, band{green, 240})
	if got, want := colors(Extract(large)), colors(Extract(small)); !sameColors(got, want) {
		t.Errorf("8 Mpx answer %v differs from the 4 Kpx answer %v", got, want)
	}
}

// TestExtractionIsDeterministic: no randomness in the clustering, so a
// textured image — not flat blocks — answers identically twice.
func TestExtractionIsDeterministic(t *testing.T) {
	img := textured(400, 300)
	first, second := Extract(img), Extract(img)
	if !sameColors(colors(first), colors(second)) {
		t.Errorf("two runs disagreed: %v then %v", colors(first), colors(second))
	}
	for i := range first {
		if first[i] != second[i] {
			t.Errorf("candidate %d differed between runs: %+v then %+v", i, first[i], second[i])
		}
	}
}

// TestMaxAndSeparation: the maximum is honoured, and near-duplicates do not
// fill the row. The fixture is six barely distinguishable blues and one red;
// the blues must collapse to about one entry, leaving the red visible.
func TestMaxAndSeparation(t *testing.T) {
	parts := []band{{red, 10}}
	for i := 0; i < 6; i++ {
		parts = append(parts, band{rgb(0x20+uint8(i), 0x60+uint8(i), 0xc8+uint8(i)), 15})
	}
	got := Extract(bands(64, parts...))
	if len(got) > 2 {
		t.Errorf("got %d candidates %v; the near-identical blues should have collapsed", len(got), colors(got))
	}
	if indexOf(got, red) < 0 {
		t.Errorf("got %v, want the red among them", colors(got))
	}
	capped := ExtractWith(textured(400, 300), Options{Max: 3})
	if len(capped) != 3 {
		t.Errorf("Max 3 gave %d candidates %v", len(capped), colors(capped))
	}
	for i := range capped {
		for j := i + 1; j < len(capped); j++ {
			if d := math.Sqrt(oklabDistance(capped[i].Color, capped[j].Color)); d < DefaultSeparation {
				t.Errorf("candidates %d and %d are %.3f apart, under the %.2f separation", i, j, d, DefaultSeparation)
			}
		}
	}
}

// TestCandidatesAreImagePixels: every candidate is a colour the image
// contains, which is what makes a swatch honest and keeps the answer in
// gamut without a mapping step.
func TestCandidatesAreImagePixels(t *testing.T) {
	img := textured(200, 150)
	present := map[stdcolor.NRGBA]bool{}
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			present[img.NRGBAAt(x, y)] = true
		}
	}
	for i, c := range Extract(img) {
		if !present[c.Color] {
			t.Errorf("candidate %d %v is not a pixel of the image", i, c.Color)
		}
	}
}

// TestWeightOrderingIsMonotonic: whatever the picture, the answer comes
// back in descending prominence.
func TestWeightOrderingIsMonotonic(t *testing.T) {
	for name, img := range map[string]image.Image{
		"textured": textured(300, 200),
		"bands":    bands(64, band{grey, 60}, band{orange, 20}, band{blue, 12}, band{green, 8}),
		"grey":     bands(64, band{rgb(0x30, 0x30, 0x30), 40}, band{rgb(0xc0, 0xc0, 0xc0), 60}),
	} {
		got := Extract(img)
		for i := 1; i < len(got); i++ {
			if got[i].Weight > got[i-1].Weight {
				t.Errorf("%s: candidate %d weight %.6f exceeds its predecessor's %.6f", name, i, got[i].Weight, got[i-1].Weight)
			}
		}
	}
}

// textured paints a fixture with no flat regions at all: a vertical
// lightness gradient in a blue-green over most of the frame, a warm band
// across the lower third, and a small vivid patch — the shape of a
// photograph, built from arithmetic so no picture has to be stored.
func textured(w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		t := float64(y) / float64(h)
		for x := 0; x < w; x++ {
			u := float64(x) / float64(w)
			var c stdcolor.NRGBA
			switch {
			case t < 0.62:
				c = rgb(uint8(60+90*t+20*u), uint8(130+70*t), uint8(200-40*t))
			default:
				c = rgb(uint8(190-30*t), uint8(140+40*u), uint8(70+30*t))
			}
			if x > w*7/10 && x < w*8/10 && y > h/4 && y < h/3 {
				c = rgb(0xd8, 0x1e, 0x2c)
			}
			img.SetNRGBA(x, y, c)
		}
	}
	return img
}

func sameColors(a, b []stdcolor.NRGBA) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestStrideHonoursTheBudgetAtEveryShape: the sampling cap is a cap, not an
// estimate. A long thin image is the case an area-based step gets wrong —
// the step that would cover a square of the same area leaves a
// one-pixel-wide strip sampled on every row of its length.
func TestStrideHonoursTheBudgetAtEveryShape(t *testing.T) {
	for _, wh := range [][2]int{
		{1, 1}, {1, 20_000_000}, {20_000_000, 1}, {5000, 4000}, {141, 141}, {1, 20001},
	} {
		w, h := wh[0], wh[1]
		s := stride(w, h, DefaultSamples)
		count := ((w + s - 1) / s) * ((h + s - 1) / s)
		if count > DefaultSamples {
			t.Errorf("%dx%d at stride %d visits %d pixels, over the %d budget", w, h, s, count, DefaultSamples)
		}
		if s > 1 && count*4 < DefaultSamples {
			t.Errorf("%dx%d at stride %d visits only %d pixels, far under the %d budget", w, h, s, count, DefaultSamples)
		}
	}
}

// One hue at four depths, from a photograph of a lake sampled before this
// gathering existed: the row it produced then was four blues and offered one
// choice dressed as four.
var (
	skyDeep = rgb(0x1f, 0x65, 0xa1) // hue 249, chroma 0.118 — the most chromatic
	skyMid  = rgb(0x39, 0x84, 0xbe) // hue 245, chroma 0.116
	skyPale = rgb(0x61, 0xa2, 0xcb) // hue 238, chroma 0.090
	skyDark = rgb(0x19, 0x4a, 0x69) // hue 240, chroma 0.074
)

// TestOneHueAtSeveralDepthsIsOneCandidate: a sky is one seed however many
// depths it was sampled at, it holds the share of all of them, and the
// swatch offered for it is the member that carries the hue most clearly.
func TestOneHueAtSeveralDepthsIsOneCandidate(t *testing.T) {
	img := bands(64, band{skyMid, 40}, band{skyPale, 25}, band{skyDeep, 20}, band{skyDark, 15})
	got := Extract(img)
	if len(got) != 1 {
		t.Fatalf("four depths of one blue gave %d candidates %v, want 1", len(got), colors(got))
	}
	if got[0].Color != skyDeep {
		t.Errorf("the gathered blue is %v, want the most chromatic member %v", got[0].Color, skyDeep)
	}
	if math.Abs(got[0].Share-1) > 0.02 {
		t.Errorf("share = %.3f, want the four depths' shares together", got[0].Share)
	}
}

// TestHueTellsColoursApart: the gathering is a hue question, so colours a
// hue apart survive it and colours a shade apart do not — asserted at both
// sides of the threshold rather than on one example.
func TestHueTellsColoursApart(t *testing.T) {
	near := rgb(0x2d, 0x7a, 0xb0) // ~10 degrees off skyMid
	far := rgb(0x2f, 0xb5, 0x4a)  // green, most of a quarter-turn away
	if got := Extract(bands(64, band{skyMid, 50}, band{near, 50})); len(got) != 1 {
		t.Errorf("two blues %v apart gave %d candidates %v, want 1",
			hueGap(skyMid, near), len(got), colors(got))
	}
	if got := Extract(bands(64, band{skyMid, 50}, band{far, 50})); len(got) != 2 {
		t.Errorf("a blue and a green %v apart gave %d candidates %v, want 2",
			hueGap(skyMid, far), len(got), colors(got))
	}
}

// TestGatheredSharesAccountForThePicture: with room for every colour in it,
// the shares add up to the whole picture. A row whose percentages sum to a
// third of the image is a row whose percentages mean nothing.
func TestGatheredSharesAccountForThePicture(t *testing.T) {
	img := bands(64,
		band{skyMid, 30}, band{skyPale, 20}, band{skyDark, 10}, // one blue, three depths
		band{grey, 20}, band{orange, 12}, band{green, 8},
	)
	total := 0.0
	for _, c := range Extract(img) {
		total += c.Share
	}
	if math.Abs(total-1) > 0.02 {
		t.Errorf("the candidates account for %.0f%% of the picture, want all of it", total*100)
	}
}
