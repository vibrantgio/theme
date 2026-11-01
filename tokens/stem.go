package tokens

import (
	"fmt"
	"sync"

	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/text"
	"github.com/vibrantgio/font/roboto"
	"github.com/vibrantgio/font/roboto/regular/normal"

	"github.com/vibrantgio/theme/internal/widen"
)

// StemWeight is the stem weight a scheme draws typography role text at. Its
// value is the width the library adds to each stem of Roboto Regular
// upright, in em, in total over the two sides of the stem. Zero is the face
// as shipped. The other Roboto faces draw as shipped at every stem weight.
type StemWeight float32

// The stem weight of each scheme. Each value makes the mean stem of
// "Vibrant Gio" over four horizontal offsets, 0, ¼, ½ and ¾ px, the
// platform's at 13 pt and 10 pt. The reference's section "What the same
// text measures in a platform window and a library window" holds the
// platform's rows and the library's.
//
// Light 0.031 em: 13 pt 1.29 px against the platform's 1.26, and 10 pt
// 0.92 px against 0.97. Dark 0 em: the face as shipped gives 13 pt 1.49 px
// against 1.51, and 10 pt 1.20 px against 1.17.
const (
	StemWeightLight StemWeight = 0.031
	StemWeightDark  StemWeight = 0
)

// SchemeStemWeight returns the stem weight of the dark scheme when dark is
// true, and of the light scheme when it is false.
func SchemeStemWeight(dark bool) StemWeight {
	if dark {
		return StemWeightDark
	}
	return StemWeightLight
}

// Dark reports whether p is a dark scheme's set: its label is light.
func (p PlatformColors) Dark() bool {
	return int(p.Label.R)+int(p.Label.G)+int(p.Label.B) > 3*127
}

// WithStemWeight returns t with its text drawn at the stem weight w. The
// roles, their sizes and their line heights stay the same. The copy shares
// t's shaper cache, so every copy at one stem weight hands back one shaper.
// A theme calls it with the stem weight of its scheme, [SchemeStemWeight].
func (t Typography) WithStemWeight(w StemWeight) Typography {
	t.stem = w
	return t
}

// StemWeight returns the stem weight t draws its text at.
func (t Typography) StemWeight() StemWeight {
	return t.stem
}

// shaperAt is the shaper Shaper hands back for the stem weight w.
func (t *Typography) shaperAt(w StemWeight) *text.Shaper {
	shaperMu.Lock()
	defer shaperMu.Unlock()
	c := t.cache()
	if w == 0 {
		if c.shaper == nil {
			c.shaper = text.NewShaper(text.WithCollection(t.Faces))
		}
		return c.shaper
	}
	if c.stems[w] == nil {
		if c.stems == nil {
			c.stems = map[StemWeight]*text.Shaper{}
		}
		c.stems[w] = text.NewShaper(text.WithCollection(stemFaces(t.Faces, w)))
	}
	return c.stems[w]
}

// pinnedShaperAt is the shaper DeterministicShaper hands back for the stem
// weight w.
func (t *Typography) pinnedShaperAt(w StemWeight) *text.Shaper {
	shaperMu.Lock()
	defer shaperMu.Unlock()
	c := t.cache()
	if w == 0 {
		if c.pinnedShaper == nil {
			c.pinnedShaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(t.Faces))
		}
		return c.pinnedShaper
	}
	if c.pinnedStems[w] == nil {
		if c.pinnedStems == nil {
			c.pinnedStems = map[StemWeight]*text.Shaper{}
		}
		c.pinnedStems[w] = text.NewShaper(text.NoSystemFonts(), text.WithCollection(stemFaces(t.Faces, w)))
	}
	return c.pinnedStems[w]
}

// stemFaces returns a copy of faces with the face of Roboto Regular upright
// at the stem weight w. The other faces are the same entries.
func stemFaces(faces []font.FontFace, w StemWeight) []font.FontFace {
	out := make([]font.FontFace, len(faces))
	for i, f := range faces {
		if f.Font == roboto.RegularNormal {
			f.Face = stemFace(w)
		}
		out[i] = f
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
