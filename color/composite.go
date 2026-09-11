// Alpha compositing: the opaque colour a translucent fill actually lands as
// on the surface under it, which is the only colour a contrast measurement
// can be taken against.
package color

import (
	stdcolor "image/color"
	"math"
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
