// Alpha compositing: the opaque colour a translucent fill actually lands as
// on the surface under it, which is the only colour a contrast measurement
// can be taken against.
package color

import (
	stdcolor "image/color"
	"math"
	"sync"
)

// Flatten returns src laid over the opaque surface dst, composited per
// channel in the encoded sRGB space — a·src + (1−a)·dst on the eight-bit
// channels, rounded — and opaque.
//
// Encoded sRGB is where the platform composites an alpha colour, measured
// off the Save dialog captures in the organization's macOS reference:
// labelColor, black at 0.85, lands on #242424 over a push button's #ececec
// fill and its white lands on #e0e1e2 over the dark button's #333a3f, and
// secondaryLabelColor's white lands on #9c9fa1 over the dark sheet's
// #232a2f — every channel of every one reproduced by this blend and by no
// other. Gio's rasterizer mixes a translucent fill it is handed in linear
// light instead, which puts that same label on #6c6c6c over white where the
// platform puts #272727: not a shade off, a different colour. So a caller
// flattens the platform's alpha here and hands the rasterizer an opaque
// fill.
//
// src's alpha is its coverage. dst is a surface — what is already on the
// screen when the fill is drawn on it — so its alpha is ignored and the
// result is opaque, which is also what makes the result something [APCA] can
// be handed. Coverage 0 returns dst and coverage 255 returns src, both
// exactly.
func Flatten(src, dst stdcolor.NRGBA) stdcolor.NRGBA {
	a := float64(src.A) / 255
	mix := func(s, d uint8) uint8 {
		return uint8(math.Round(a*float64(s) + (1-a)*float64(d)))
	}
	return stdcolor.NRGBA{
		R: mix(src.R, dst.R),
		G: mix(src.G, dst.G),
		B: mix(src.B, dst.B),
		A: 0xff,
	}
}

// LinearCoverage returns the coverage to hand a renderer that blends in
// linear light so that what it lands is what [Flatten] lands at coverage a
// on the same surface.
//
// Gio blends a translucent fill it is handed in linear light: it renders
// into an sRGB render target, so the hardware decodes the surface byte,
// mixes, and encodes the result again. The platform mixes the byte itself.
// One coverage cannot carry that difference exactly — the platform's blend
// is affine on the encoded byte and Gio's is affine on the linear value, and
// no single coverage maps one onto the other for every surface — so this is
// the coverage whose worst miss across all 256 surface bytes is least,
// smallest sum of squares among those that tie.
//
// The miss it leaves is measured against [Flatten] on every surface byte:
// one 255th at and below 0x13, the floating shadow's peak, two to 0x26,
// three at the scrim's 0x33, and seven at 0xcc. An overlay that can read the
// pixels beneath it flattens them per pixel instead and misses nothing; this
// is for the overlay that cannot.
//
// Coverage 0 and coverage 255 return themselves.
func LinearCoverage(a uint8) uint8 {
	linearCoverage.once.Do(func() {
		for d := range linearCoverage.surface {
			linearCoverage.surface[d] = linearFromSRGB(float64(d) / 255)
		}
	})
	linearCoverage.mu.Lock()
	defer linearCoverage.mu.Unlock()
	if c, ok := linearCoverage.fitted[a]; ok {
		return c
	}
	c := fitCoverage(a)
	linearCoverage.fitted[a] = c
	return c
}

var linearCoverage = struct {
	once    sync.Once
	surface [256]float64
	mu      sync.Mutex
	fitted  map[uint8]uint8
}{fitted: map[uint8]uint8{}}

// fitCoverage scans every coverage at or above a — a linear blend always
// needs at least as much coverage as an encoded one to reach the same
// place — and keeps the one whose predictions across all 256 surface bytes
// miss [Flatten] least.
func fitCoverage(a uint8) uint8 {
	if a == 0 || a == 0xff {
		return a
	}
	want := [256]uint8{}
	for d := range want {
		want[d] = uint8(math.Round((1 - float64(a)/255) * float64(d)))
	}
	best, bestWorst, bestSquares := a, math.MaxInt, math.MaxInt
	for c := int(a); c <= 0xff; c++ {
		keep := 1 - float64(c)/255
		worst, squares := 0, 0
		for d := range want {
			got := int(math.Round(sRGBFromLinear(keep*linearCoverage.surface[d]) * 255))
			e := got - int(want[d])
			if e < 0 {
				e = -e
			}
			if e > worst {
				worst = e
			}
			squares += e * e
		}
		if worst < bestWorst || (worst == bestWorst && squares < bestSquares) {
			best, bestWorst, bestSquares = uint8(c), worst, squares
		}
	}
	return best
}

// The sRGB transfer function, both ways, on the unit interval. These are the
// curve the hardware applies for us either side of a blend in linear light,
// written out here because the fit has to predict what that blend lands.
func linearFromSRGB(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func sRGBFromLinear(c float64) float64 {
	if c <= 0.0031308 {
		return c * 12.92
	}
	return 1.055*math.Pow(c, 1/2.4) - 0.055
}

// Fade returns c at coverage of the coverage it already carries: the colour
// a control's own paint takes when the control is drawn at less than full
// strength. An opaque fill comes back translucent, and a colour that already
// carries a coverage — a seam, a label — comes back at that coverage scaled
// down, so one call covers both. [Flatten] then lands it on the surface the
// control stands on.
//
// coverage is a fraction of 255, the way every coverage in this library is
// carried: 255 returns c unchanged and 0 returns it at no coverage at all.
// The rounding is the same half-up rounding [Flatten] uses, so the two
// compose without a second rounding rule.
func Fade(c stdcolor.NRGBA, coverage uint8) stdcolor.NRGBA {
	c.A = uint8(math.Round(float64(c.A) * float64(coverage) / 255))
	return c
}
