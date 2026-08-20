package imageseed

import (
	stdcolor "image/color"
	"math"
	"testing"

	"github.com/vibrantgio/theme/color"
)

// The palette fixtures reuse the image fixtures' colours, which is the point
// of the exercise: the two entry points are meant to rank the same colours the
// same way, so a test that ranked a palette by different colours would prove
// nothing about that.

// repeat is a palette entry listed n times — the shape a curated palette
// arrives in, where one ink is worn by eight kinds of token and another by
// one.
func repeat(c stdcolor.NRGBA, n int) []stdcolor.NRGBA {
	out := make([]stdcolor.NRGBA, n)
	for i := range out {
		out[i] = c
	}
	return out
}

func palette(parts ...[]stdcolor.NRGBA) []stdcolor.NRGBA {
	var out []stdcolor.NRGBA
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

// TestExtractPaletteRanksLikeTheImagePath is the core fixture, and it is the
// image path's core fixture written as a palette: the same four colours in the
// same proportions, ranked the same way. The three vivid entries lead the grey
// that fills more of the palette than all three together, and among themselves
// they fall in the order their prominence implies.
func TestExtractPaletteRanksLikeTheImagePath(t *testing.T) {
	got := ExtractPalette(palette(
		repeat(grey, 60),
		repeat(orange, 20),
		repeat(blue, 12),
		repeat(green, 8),
	))
	want := []stdcolor.NRGBA{orange, blue, green, grey}
	if len(got) != len(want) {
		t.Fatalf("got %d candidates %v, want %d %v", len(got), colors(got), len(want), want)
	}
	for i := range want {
		if got[i].Color != want[i] {
			t.Errorf("rank %d = %v, want %v (whole order %v)", i, got[i].Color, want[i], colors(got))
		}
	}
	for i, share := range []float64{0.20, 0.12, 0.08, 0.60} {
		if math.Abs(got[i].Share-share) > 0.001 {
			t.Errorf("rank %d share = %.3f, want %.3f", i, got[i].Share, share)
		}
	}
}

// TestExtractPaletteWeighsRepeats pins how a caller weights an entry: by
// listing it more than once. One blue against one slightly less vivid orange
// is the blue's answer; the same two with the orange worn twice is the
// orange's, because share is the other half of what prominence measures.
func TestExtractPaletteWeighsRepeats(t *testing.T) {
	once := ExtractPalette([]stdcolor.NRGBA{orange, blue})
	if len(once) != 2 || once[0].Color != blue {
		t.Errorf("one each: leader %v, want %v (whole order %v)", once[0].Color, blue, colors(once))
	}
	twice := ExtractPalette(palette(repeat(orange, 2), repeat(blue, 1)))
	if len(twice) != 2 || twice[0].Color != orange {
		t.Errorf("orange twice: leader %v, want %v (whole order %v)", twice[0].Color, orange, colors(twice))
	}
	if math.Abs(twice[0].Share-2.0/3.0) > 0.001 {
		t.Errorf("orange share = %.3f, want %.3f", twice[0].Share, 2.0/3.0)
	}
}

// TestExtractPaletteGathersOneHue holds the gathering step, which a palette
// needs as much as a photograph does: a style that draws its strings pale and
// its keywords deep in one blue has one blue, not three, and an answer that
// spent three of its places on it would offer a choice it does not have. The
// three entries are built at one OKLCh hue so the fixture cannot drift into
// being three hues by accident.
func TestExtractPaletteGathersOneHue(t *testing.T) {
	l, chroma, hue := color.OKLChFromNRGBA(blue)
	pale := color.NRGBAFromOKLCh(l+0.20, chroma*0.6, hue)
	deep := color.NRGBAFromOKLCh(l-0.15, chroma*0.8, hue)
	got := ExtractPalette([]stdcolor.NRGBA{pale, blue, deep})
	if len(got) != 1 {
		t.Fatalf("got %d candidates %v, want 1", len(got), colors(got))
	}
	if got[0].Color != blue {
		t.Errorf("swatch = %v, want the most chromatic member %v", got[0].Color, blue)
	}
	if math.Abs(got[0].Share-1) > 0.001 {
		t.Errorf("share = %.3f, want 1", got[0].Share)
	}
}

// TestExtractPaletteNearGrey is the degenerate palette the style gallery will
// actually meet: a base whose whole scheme is greys, told apart by weight
// rather than by colour. It has to come back as its own greys ranked by how
// much of the palette they are — muted candidates — rather than as whichever
// grey rounding gave the largest hue to.
func TestExtractPaletteNearGrey(t *testing.T) {
	var (
		black = rgb(0x00, 0x00, 0x00)
		mid   = rgb(0x66, 0x66, 0x66)
		pale  = rgb(0xcc, 0xcc, 0xcc)
	)
	got := ExtractPalette(palette(repeat(black, 1), repeat(mid, 4), repeat(pale, 2)))
	want := []stdcolor.NRGBA{mid, pale, black}
	if len(got) != len(want) {
		t.Fatalf("got %d candidates %v, want %d %v", len(got), colors(got), len(want), want)
	}
	for i := range want {
		if got[i].Color != want[i] {
			t.Errorf("rank %d = %v, want %v (whole order %v)", i, got[i].Color, want[i], colors(got))
		}
		if got[i].Chroma > greyChroma {
			t.Errorf("rank %d chroma = %.4f, want no more than the grey floor %.4f",
				i, got[i].Chroma, greyChroma)
		}
	}
}

// TestExtractPaletteDegenerate walks the cases that have no colours to rank,
// or only one. None of them is a special case in the code and none of them may
// need one.
func TestExtractPaletteDegenerate(t *testing.T) {
	for _, tc := range []struct {
		name    string
		in      []stdcolor.NRGBA
		want    int
		leader  stdcolor.NRGBA
		wantAll float64
	}{
		{name: "nil", in: nil},
		{name: "empty", in: []stdcolor.NRGBA{}},
		{name: "transparent", in: []stdcolor.NRGBA{{R: 0xff, G: 0x00, B: 0x00, A: 0x00}}},
		{name: "one colour", in: []stdcolor.NRGBA{red}, want: 1, leader: red, wantAll: 1},
		{name: "one colour repeated", in: repeat(red, 9), want: 1, leader: red, wantAll: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ExtractPalette(tc.in)
			if len(got) != tc.want {
				t.Fatalf("got %d candidates %v, want %d", len(got), colors(got), tc.want)
			}
			if tc.want == 0 {
				return
			}
			if got[0].Color != tc.leader {
				t.Errorf("leader = %v, want %v", got[0].Color, tc.leader)
			}
			if math.Abs(got[0].Share-tc.wantAll) > 0.001 {
				t.Errorf("share = %.3f, want %.3f", got[0].Share, tc.wantAll)
			}
		})
	}
}

// TestExtractPaletteSkipsTransparent holds the alpha rule the image path has:
// an entry nobody can see is not a colour the palette has, and it does not
// dilute the shares of the ones that are.
func TestExtractPaletteSkipsTransparent(t *testing.T) {
	clear := stdcolor.NRGBA{R: 0x00, G: 0xff, B: 0x00, A: 0x02}
	got := ExtractPalette([]stdcolor.NRGBA{red, clear, blue})
	if len(got) != 2 {
		t.Fatalf("got %d candidates %v, want 2", len(got), colors(got))
	}
	for _, c := range got {
		if math.Abs(c.Share-0.5) > 0.001 {
			t.Errorf("%v share = %.3f, want 0.5 — the invisible entry was counted", c.Color, c.Share)
		}
	}
}

// TestExtractPaletteOrderIndependent is the determinism the whole feature
// rests on: a palette read off a style comes out of a map somewhere, and the
// answer must not depend on which order it happened to arrive in.
func TestExtractPaletteOrderIndependent(t *testing.T) {
	in := palette(repeat(grey, 5), repeat(orange, 3), repeat(blue, 2), repeat(beige, 4))
	want := ExtractPalette(in)
	// A fixed stride rather than a random shuffle — a test that fails only
	// sometimes is a test nobody can act on. The stride is coprime with the
	// length, so it visits every entry exactly once and the two lists hold
	// the same colours in a thoroughly different order.
	shuffled := make([]stdcolor.NRGBA, len(in))
	for i := range in {
		shuffled[i] = in[(i*5)%len(in)]
	}
	got := ExtractPalette(shuffled)
	if len(got) != len(want) {
		t.Fatalf("shuffled gave %d candidates %v, want %d %v",
			len(got), colors(got), len(want), colors(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("rank %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestExtractPaletteWithOptions holds that the options a palette can honour
// are honoured, and that the one it cannot is harmless.
func TestExtractPaletteWithOptions(t *testing.T) {
	in := []stdcolor.NRGBA{red, orange, green, blue, grey, beige}
	if got := ExtractPaletteWith(in, Options{Max: 2}); len(got) != 2 {
		t.Errorf("Max 2 gave %d candidates %v", len(got), colors(got))
	}
	// Samples is an image's pixel budget and a palette has no pixels; naming
	// one must not change the answer.
	want := ExtractPalette(in)
	got := ExtractPaletteWith(in, Options{Samples: 4})
	if len(got) != len(want) {
		t.Fatalf("Samples changed the answer: %v vs %v", colors(got), colors(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Samples changed rank %d: %+v vs %+v", i, got[i], want[i])
		}
	}
}
