// palette.go — seed candidates from a palette somebody already chose, rather
// than from a photograph.
//
// A picture has to be clustered before it can be ranked: a photograph holds
// tens of thousands of distinct values and none of them is a decision. A
// palette is the opposite — a few dozen colours, each one deliberate — so
// there is nothing to cluster and clustering it would be a way of losing
// detail rather than of finding it. The entry point below therefore skips
// straight to the two steps that were never about pixels: gathering colours
// of one hue into one answer, and ranking those answers by chroma-weighted
// prominence. Everything from that point on is the image path's own code,
// which is what keeps a candidate row derived from a palette and a candidate
// row derived from a photograph the same kind of thing.
//
// Weight comes from the list. A colour listed once counts once and a colour
// listed twice counts twice, so a caller holding a palette where some colours
// carry more of it than others — a syntax style that draws eight kinds of
// name in one blue and one number in one orange — weights the answer by
// handing that fact over as repeats, and a caller with a plain list of
// distinct colours gets the uniform reading.

package imageseed

import (
	stdcolor "image/color"

	"github.com/vibrantgio/theme/color"
)

// ExtractPalette returns seed candidates for a curated palette under the
// default options, most prominent first. See ExtractPaletteWith.
func ExtractPalette(palette []stdcolor.NRGBA) []Candidate {
	return ExtractPaletteWith(palette, Options{})
}

// ExtractPaletteWith returns seed candidates for a curated palette, most
// prominent first, under the given options.
//
// The colours enter as they are given: there is no clustering step, because a
// palette's entries are already the decisions clustering exists to recover.
// What follows is the image path exactly — colours carrying one hue are
// gathered into one candidate, a gathered candidate takes its swatch from its
// most chromatic member, and the answer is ranked by share times the square
// of chroma lifted off a floor. So a vivid accent used for one token leads a
// drab colour used for ten, on the same terms a vivid tenth of a photograph
// leads a drab half of it.
//
// Share is counted over entries, so an entry appearing twice in the list
// weighs twice as much as one appearing once, and the shares of one answer
// sum to at most 1. Near-transparent entries carry no colour anyone sees and
// are skipped, exactly as near-transparent pixels are.
//
// A degenerate palette needs no special case and gets none. A palette with no
// colour in it at all — a monochrome style's blacks and greys — ranks its
// greys by share and comes back as those greys, muted candidates rather than
// hues invented out of rounding, because the floor under the chroma emphasis
// is what carries the ordering when nothing is vivid. A palette of one colour
// yields that one colour at share 1. An empty palette, or one that is
// entirely transparent, yields no candidates at all.
//
// [Options].Samples has nothing to read and is ignored; Max and Separation
// mean what they mean for an image.
func ExtractPaletteWith(palette []stdcolor.NRGBA, o Options) []Candidate {
	o = o.withDefaults()
	groups, entries := tally(palette)
	if entries == 0 {
		return nil
	}
	return choose(rank(groups, entries), o.Max, o.Separation)
}

// tally counts a palette into one group per distinct colour — the shape the
// clustering step hands on — and returns the number of entries counted.
//
// The groups come out in first-seen order and not in a map's order, and the
// ranking that reads them is a total order regardless, so two runs over one
// palette cannot disagree about the answer or about the order of it.
func tally(palette []stdcolor.NRGBA) (groups []group, entries int) {
	at := make(map[uint32]int, len(palette))
	groups = make([]group, 0, len(palette))
	for _, c := range palette {
		if c.A < minAlpha {
			continue
		}
		c.A = 0xff
		entries++
		if i, seen := at[packRGB(c)]; seen {
			groups[i].count++
			continue
		}
		l, a, b := color.OKLabFromNRGBA(c)
		at[packRGB(c)] = len(groups)
		groups = append(groups, group{count: 1, medoid: pixel{c: c, l: l, a: a, b: b}})
	}
	return groups, entries
}
