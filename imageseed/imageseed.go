// Package imageseed extracts brand-seed candidates from an image: the
// colours the picture is actually made of, clustered in the same perceptual
// space a palette derivation reads a seed in, and ranked so a vivid minority
// leads a dull majority.
//
// Extract returns candidates, never a palette. Each Candidate carries a
// colour that occurs in the image — a pixel, not a computed average that
// nothing in the picture ever had — together with the share of sampled
// pixels it stands for and the chroma-weighted prominence it was ranked by.
// Deriving tokens from one of them is the caller's step; every candidate is
// a legal seed.
//
// # Placement
//
// This package sits beside color rather than carrying conversions of its
// own. The OKLab and OKLCh transforms it clusters and ranks in are that
// package's exported surface, called from here: a second copy of the
// mathematics would be a second answer to a question the module already
// answers once, free to drift from the answer palettes are actually derived
// with. Clustering in the space the derivation reads a seed in is the whole
// point of the exercise, so the space has to be the same space.
//
// # The steps
//
// Sample. Pixels are read on a square stride solved so that at most
// Options.Samples of them are visited, whatever the image measures — a
// twenty-megapixel photograph costs what a thumbnail costs, because colour
// statistics converge long before a full scan does. Pixels that are all but
// transparent carry no colour anyone sees and are skipped; an image with
// none left yields no candidates at all.
//
// Cluster. Lloyd's algorithm in OKLab, where plain Euclidean distance is
// the perceptual difference the space was built to make it. The starting
// centroids are not random: sampled colours are counted into a coarse OKLab
// grid and the most populated cells, ordered by population and then by
// cell so ties cannot swing, seed the iteration. Extraction is therefore
// deterministic — the same image yields the same candidates in the same
// order on every machine and every run — which is what lets a fixture test
// assert on rank rather than on membership. More clusters are formed than
// candidates are returned, so the gathering has something to work with and
// the ranking has something to choose between.
//
// Gather. Clusters of one hue are one colour, however many depths they
// were sampled at. A sky is pale at the horizon and deep overhead and
// clusters honestly into several, but it is one seed, and an answer that
// spent four of its six places on it would offer a choice it does not have.
// So two colours that both carry a hue are the same colour when their hues
// are within a quarter-turn-to-the-complement of each other, whatever their
// depths; where a colour is too grey to have a hue that question has no
// answer and sameness falls back to distance in OKLab, which is how a
// greyscale picture still comes back as its own several greys. A gathered
// colour holds its members' shares together and takes its swatch from the
// most chromatic of them, since that is the member that names the hue most
// clearly.
//
// Rank. Prominence is the share of sampled pixels a colour holds times an
// emphasis on its chroma: the chroma lifted off zero by a floor, squared.
// The floor is what keeps the ordering meaningful when nothing in the
// picture is vivid — a wholly grey image still ranks its greys by share.
// The square is what lets a vivid tenth of an image outrank a drab half of
// it, which is the ordering a person picking a brand colour out of a
// photograph wants: a colour ten times more common leads a rival only
// while that rival is less than about three times more colourful.
//
// Degenerate images need no special case and get none. A greyscale image
// clusters into greys and, their chromas being equal, ranks them by share.
// A single-colour image yields that one colour at share 1. An image smaller
// than the sample budget is read whole. An image of one pixel yields one
// candidate.
//
// # A palette instead of an image
//
// [ExtractPalette] takes a list of colours somebody already chose — a syntax
// style's colours, a brand sheet, a scheme copied off a screenshot — and
// returns the same kind of answer. Only the first two steps differ: there is
// nothing to sample and nothing to cluster, because a palette's entries
// already are the decisions clustering exists to recover from a photograph.
// Gathering, ranking and the degenerate cases are the image path's own, so a
// candidate row derived from a palette and one derived from a picture are
// comparable things.
package imageseed

import (
	"image"
	stdcolor "image/color"
	"math"
	"sort"

	"github.com/vibrantgio/theme/color"
)

// Defaults for Options, applied to any field left at zero.
const (
	// DefaultMax is the number of candidates Extract returns at most: a
	// row a person can compare at a glance rather than a catalogue.
	DefaultMax = 6
	// DefaultSamples is the pixel budget. It is what bounds the cost of a
	// large image, and it is well past the point where a dominant colour
	// stops moving: the sample is a regular grid over the whole frame, not
	// a corner of it.
	DefaultSamples = 20000
	// DefaultSeparation is the smallest OKLab distance between two
	// returned candidates that hue cannot tell apart. Below it two greys
	// read as one grey, and a row of near-duplicates offers no choice at
	// all.
	DefaultSeparation = 0.08
)

const (
	// greyChroma is the chroma at or under which a colour has no usable
	// hue left. It is the floor the ranking emphasis is lifted off: a
	// cluster this drab is ranked almost entirely on how much of the image
	// it covers, and a wholly grey image is ranked on share alone rather
	// than on nothing at all.
	greyChroma = 0.020
	// initCell is the edge of one cell of the coarse OKLab grid the
	// starting centroids are counted into. Roughly a dozen cells across
	// each axis of the space colours actually occupy: fine enough to
	// separate the picture's real colours, coarse enough that the count in
	// a cell means something.
	initCell = 0.08
	// maxIterations bounds Lloyd's algorithm. Colour clustering settles in
	// a handful of passes; the bound is a guard, not a schedule, and the
	// loop leaves early whenever a pass moves nothing.
	maxIterations = 12
	// hueSeparation is how far apart two hues must be, in degrees, to be
	// two colours rather than one colour at two depths. A quarter of the
	// way to a complement: wide enough that a sky sampled from horizon to
	// zenith comes back as one blue, narrow enough that a blue and a teal
	// stay two answers.
	hueSeparation = 25
	// clustersPerCandidate and minClusters size the clustering relative to
	// the answer: forming several times more clusters than are returned
	// gives the ranking room to prefer a vivid cluster over a large one,
	// and gives the separation filter something to fall back on when it
	// rejects a near-duplicate.
	clustersPerCandidate = 3
	minClusters          = 8
	// minAlpha is the alpha below which a pixel is treated as carrying no
	// colour. Fully transparent pixels are common in exported artwork and
	// are nobody's idea of the image's colour.
	minAlpha = 8
)

// Candidate is one extracted seed: a colour the image contains, with the
// evidence for it.
type Candidate struct {
	// Color is a sampled pixel — a colour the image really has, in gamut
	// by construction because it was read out of the picture rather than
	// computed into it: the pixel nearest the centre of the most chromatic
	// cluster gathered into this colour.
	Color stdcolor.NRGBA
	// Share is the fraction of sampled pixels this colour holds, 0 to 1,
	// counting every cluster gathered into it. Shares over one answer sum
	// to at most 1 — to less than 1 when the picture holds more colours
	// than were asked for, since the ones left out keep their share.
	Share float64
	// Chroma is Color's OKLCh chroma — how much colour it has, on the same
	// axis a palette derivation reads a seed's chroma on. It is the
	// highest chroma among the clusters gathered into this colour.
	Chroma float64
	// Weight is the chroma-weighted prominence the candidates were ordered
	// by: Share times the square of Chroma lifted off zero. It is
	// comparable within one answer and meaningless between two.
	Weight float64
}

// Options tunes an extraction. Every zero field takes its default, so the
// zero Options is what Extract uses.
type Options struct {
	// Max is the number of candidates returned at most; DefaultMax when
	// zero or negative.
	Max int
	// Samples is the pixel budget; DefaultSamples when zero or negative.
	// Fewer samples cost less and read a small or flat image just as well;
	// a photograph with a small vivid subject wants the default or more.
	Samples int
	// Separation is the smallest OKLab distance between two returned
	// candidates that hue cannot tell apart — between two greys, or
	// between a grey and a colour; DefaultSeparation when zero or
	// negative. Candidates that both carry a hue are told apart by hue
	// instead, and this does not reach them.
	Separation float64
}

func (o Options) withDefaults() Options {
	if o.Max <= 0 {
		o.Max = DefaultMax
	}
	if o.Samples <= 0 {
		o.Samples = DefaultSamples
	}
	if o.Separation <= 0 {
		o.Separation = DefaultSeparation
	}
	return o
}

// Extract returns the image's seed candidates under the default options,
// most prominent first. See ExtractWith for the options and the package
// documentation for what "prominent" measures.
func Extract(img image.Image) []Candidate {
	return ExtractWith(img, Options{})
}

// ExtractWith returns the image's seed candidates, most prominent first,
// under the given options. It returns nil for a nil image, an empty one, or
// one whose pixels are all transparent; it never returns two candidates
// closer together than the separation, and never more than the maximum.
func ExtractWith(img image.Image, o Options) []Candidate {
	o = o.withDefaults()
	pixels := sample(img, o.Samples)
	if len(pixels) == 0 {
		return nil
	}
	groups := cluster(pixels, clusterCount(o.Max))
	return choose(rank(groups, len(pixels)), o.Max, o.Separation)
}

// clusterCount is how many clusters an answer of n candidates is formed
// from; see clustersPerCandidate.
func clusterCount(n int) int {
	return max(n*clustersPerCandidate, minClusters)
}

// pixel is one sampled colour, kept both as it will be shown and as it is
// clustered. Holding the sRGB value alongside the OKLab one is what lets a
// candidate be an actual pixel of the image instead of a centroid that may
// sit outside the gamut it was averaged from.
type pixel struct {
	c       stdcolor.NRGBA
	l, a, b float64
}

// sample reads the image on a square stride sized so that at most budget
// pixels are visited, and converts each to OKLab. Near-transparent
// pixels are skipped, so an image can yield fewer samples than the stride
// implies — or none.
func sample(img image.Image, budget int) []pixel {
	if img == nil {
		return nil
	}
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= 0 || h <= 0 {
		return nil
	}
	step := stride(w, h, budget)
	out := make([]pixel, 0, min(budget, w*h))
	for y := bounds.Min.Y; y < bounds.Max.Y; y += step {
		for x := bounds.Min.X; x < bounds.Max.X; x += step {
			c := stdcolor.NRGBAModel.Convert(img.At(x, y)).(stdcolor.NRGBA)
			if c.A < minAlpha {
				continue
			}
			c.A = 0xff
			l, a, b := color.OKLabFromNRGBA(c)
			out = append(out, pixel{c: c, l: l, a: a, b: b})
		}
	}
	return out
}

// stride is the square step that visits at most budget pixels of a w by h
// image. It is solved rather than estimated, because the estimate is wrong
// on a long thin image: the step that would cover a square of the same area
// leaves a one-pixel-wide strip sampled on every row of its length. Each
// pass scales the step by the square root of how far over budget it still
// is, so it settles in a couple of turns whatever the shape.
func stride(w, h, budget int) int {
	step := 1
	for {
		count := ((w + step - 1) / step) * ((h + step - 1) / step)
		if count <= budget {
			return step
		}
		next := int(math.Ceil(float64(step) * math.Sqrt(float64(count)/float64(budget))))
		if next <= step {
			next = step + 1
		}
		step = next
	}
}

// group is one settled cluster: its centre, the pixels that landed in it,
// and the real pixel nearest that centre.
type group struct {
	count  int
	medoid pixel
}

// centroid is a cluster centre during the iteration.
type centroid struct{ l, a, b float64 }

func (c centroid) distance(p pixel) float64 {
	dl, da, db := c.l-p.l, c.a-p.a, c.b-p.b
	return dl*dl + da*da + db*db // squared: only comparisons are made
}

// cluster runs Lloyd's algorithm over the samples and returns the clusters
// that kept at least one pixel. k is an upper bound: an image with fewer
// distinct coarse colours than k yields fewer clusters, which is how a
// single-colour image yields exactly one candidate.
func cluster(pixels []pixel, k int) []group {
	centres := seedCentroids(pixels, k)
	if len(centres) == 0 {
		return nil
	}
	assign := make([]int, len(pixels))
	for i := range assign {
		assign[i] = -1
	}
	// Assignment is the last thing every path does, so the clusters
	// collected below are always the ones the final centres imply — even
	// when the iteration stops on the bound rather than on convergence.
	for iter := 0; ; iter++ {
		moved := assignAll(pixels, centres, assign)
		if !moved || iter == maxIterations-1 {
			break
		}
		recentre(pixels, assign, centres)
	}
	return collect(pixels, assign, centres)
}

// assignAll puts every pixel in its nearest centre's cluster and reports
// whether any pixel changed hands.
func assignAll(pixels []pixel, centres []centroid, assign []int) bool {
	moved := false
	for i, p := range pixels {
		best, bestDist := 0, math.Inf(1)
		for j, c := range centres {
			if d := c.distance(p); d < bestDist {
				best, bestDist = j, d
			}
		}
		if assign[i] != best {
			assign[i], moved = best, true
		}
	}
	return moved
}

// recentre moves every centre to the mean of the pixels assigned to it. An
// emptied centre stays where it is: it costs nothing and may reclaim pixels
// once its neighbours move.
func recentre(pixels []pixel, assign []int, centres []centroid) {
	sums := make([]centroid, len(centres))
	counts := make([]int, len(centres))
	for i, p := range pixels {
		j := assign[i]
		sums[j].l += p.l
		sums[j].a += p.a
		sums[j].b += p.b
		counts[j]++
	}
	for j := range centres {
		if counts[j] == 0 {
			continue
		}
		n := float64(counts[j])
		centres[j] = centroid{sums[j].l / n, sums[j].a / n, sums[j].b / n}
	}
}

// collect turns a settled assignment into groups: the pixel count per
// cluster and, as the cluster's colour, the sampled pixel nearest its
// centre. Empty clusters are dropped.
func collect(pixels []pixel, assign []int, centres []centroid) []group {
	best := make([]float64, len(centres))
	groups := make([]group, len(centres))
	for i := range best {
		best[i] = math.Inf(1)
	}
	for i, p := range pixels {
		j := assign[i]
		groups[j].count++
		if d := centres[j].distance(p); d < best[j] {
			best[j], groups[j].medoid = d, p
		}
	}
	out := groups[:0]
	for _, g := range groups {
		if g.count > 0 {
			out = append(out, g)
		}
	}
	return out
}

// cell is one coarse OKLab grid cell, used only to start the iteration.
type cell struct {
	key   [3]int
	count int
	sum   centroid
}

// seedCentroids picks the starting centres deterministically: sampled
// colours are counted into a coarse OKLab grid, the cells are ordered by
// population with the cell coordinates breaking every tie, and the means of
// the k most populated become the centres. No randomness enters, so two
// runs over one image cannot disagree.
func seedCentroids(pixels []pixel, k int) []centroid {
	byCell := make(map[[3]int]*cell, k*4)
	order := make([][3]int, 0, k*4)
	for _, p := range pixels {
		key := [3]int{
			int(math.Floor(p.l / initCell)),
			int(math.Floor(p.a / initCell)),
			int(math.Floor(p.b / initCell)),
		}
		c, ok := byCell[key]
		if !ok {
			c = &cell{key: key}
			byCell[key] = c
			order = append(order, key)
		}
		c.count++
		c.sum.l += p.l
		c.sum.a += p.a
		c.sum.b += p.b
	}
	cells := make([]*cell, len(order))
	for i, key := range order {
		cells[i] = byCell[key]
	}
	sort.Slice(cells, func(i, j int) bool {
		if cells[i].count != cells[j].count {
			return cells[i].count > cells[j].count
		}
		return lessKey(cells[i].key, cells[j].key)
	})
	if k > len(cells) {
		k = len(cells)
	}
	centres := make([]centroid, k)
	for i, c := range cells[:k] {
		n := float64(c.count)
		centres[i] = centroid{c.sum.l / n, c.sum.a / n, c.sum.b / n}
	}
	return centres
}

// lessKey orders two grid cells; any total order serves, and this one reads
// the coordinates in turn.
func lessKey(a, b [3]int) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

// rank turns settled clusters into candidates ordered by chroma-weighted
// prominence. The order is total: weight first, then share, then the colour
// itself, so no two runs can order one answer differently.
func rank(groups []group, samples int) []Candidate {
	out := make([]Candidate, 0, len(groups))
	for _, g := range groups {
		share := float64(g.count) / float64(samples)
		chroma := math.Hypot(g.medoid.a, g.medoid.b)
		emphasis := greyChroma + chroma
		out = append(out, Candidate{
			Color:  g.medoid.c,
			Share:  share,
			Chroma: chroma,
			Weight: share * emphasis * emphasis,
		})
	}
	sortByWeight(out)
	return out
}

// sortByWeight orders candidates by prominence, most first. The order is
// total — weight, then share, then the colour itself — so no two runs can
// order one answer differently.
func sortByWeight(cs []Candidate) {
	sort.Slice(cs, func(i, j int) bool {
		if cs[i].Weight != cs[j].Weight {
			return cs[i].Weight > cs[j].Weight
		}
		if cs[i].Share != cs[j].Share {
			return cs[i].Share > cs[j].Share
		}
		return packRGB(cs[i].Color) < packRGB(cs[j].Color)
	})
}

// choose gathers the ranked clusters into colour families and returns the
// strongest families, up to max.
//
// A family is one colour of the picture, however many depths it was sampled
// at. A sky is light at the horizon and deep overhead, and clustering
// faithfully reports that as several clusters — but they are one seed, and a
// row that spent four of its six places on one sky would be offering a
// choice it does not have. So sameness is judged on hue, not on distance:
// two colours that both carry a hue are the same colour when their hues are
// within hueSeparation of each other, whatever their depths. Where a colour
// is too grey to have a hue the question has no answer on that axis, and
// sameness falls back to Options.Separation in OKLab — which is how a
// greyscale picture still comes back as its own several greys.
//
// A family's share is its members' shares together, so the percentages mean
// what they look like they mean. Its colour is its MOST chromatic member,
// not its largest: the members differ mostly in depth, and the one that
// carries the family's hue most clearly is the one worth offering as a
// seed. Its rank follows from the two.
func choose(ranked []Candidate, max int, separation float64) []Candidate {
	limit := separation * separation
	families := make([]Candidate, 0, max)
	for _, c := range ranked {
		at := -1
		for i, f := range families {
			if sameColour(f, c, limit) {
				at = i
				break
			}
		}
		if at < 0 {
			families = append(families, c)
			continue
		}
		families[at].Share += c.Share
		if c.Chroma > families[at].Chroma {
			families[at].Color, families[at].Chroma = c.Color, c.Chroma
		}
	}
	for i := range families {
		emphasis := greyChroma + families[i].Chroma
		families[i].Weight = families[i].Share * emphasis * emphasis
	}
	sortByWeight(families)
	if len(families) > max {
		families = families[:max]
	}
	return families
}

// sameColour reports whether two candidates name one colour of the picture.
// Both having a hue, it is a question about their hues alone; either being
// grey, it is a question about how far apart they sit in OKLab.
func sameColour(a, b Candidate, limit float64) bool {
	if a.Chroma > greyChroma && b.Chroma > greyChroma {
		return hueGap(a.Color, b.Color) < hueSeparation
	}
	return oklabDistance(a.Color, b.Color) < limit
}

// hueGap is the angle between two colours' OKLCh hues, in degrees, the
// short way round the circle.
func hueGap(a, b stdcolor.NRGBA) float64 {
	_, _, ah := color.OKLChFromNRGBA(a)
	_, _, bh := color.OKLChFromNRGBA(b)
	gap := math.Abs(ah - bh)
	if gap > 180 {
		gap = 360 - gap
	}
	return gap
}

// oklabDistance is the squared OKLab distance between two sRGB colours —
// squared because only comparisons against a squared threshold are made.
func oklabDistance(a, b stdcolor.NRGBA) float64 {
	al, aa, ab := color.OKLabFromNRGBA(a)
	bl, ba, bb := color.OKLabFromNRGBA(b)
	dl, da, db := al-bl, aa-ba, ab-bb
	return dl*dl + da*da + db*db
}

// packRGB orders colours for rank's final tie-break; any total order does,
// and this one is the obvious reading of the three bytes.
func packRGB(c stdcolor.NRGBA) uint32 {
	return uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
}
