// Alpha compositing: the opaque colour a translucent fill actually lands as
// on the ground under it, which is the only colour a contrast measurement
// can be taken against.
package color

import (
	stdcolor "image/color"
)

// Over returns src composited over dst, blended in linear light.
//
// Linear light is not a refinement here, it is where the blend happens: Gio
// converts every colour to premultiplied *linear* RGBA before it reaches the
// rasterizer, so the pixel written is the linear mix of the two and the
// eight-bit average of the two hex codes is a different colour altogether.
// The design system's own overlay scrollbar is the demonstration — the
// low-contrast-text step at 39% coverage over the light page lands on
// #CCCCCC by this route and on #BABABA by the naive one, 1.49:1 against the
// page rather than 1.80:1. Measuring a translucent ink against the wrong
// composite is how an ink no reader can find comes to be believed legible.
//
// src's alpha is its coverage. dst is a ground — what is already on the
// screen when the fill is drawn on it — so its alpha is ignored and the
// result is opaque, which is what makes the result something [ContrastRatio]
// can be handed. Coverage 0 returns dst and coverage 255 returns src, both
// exactly.
func Over(src, dst stdcolor.NRGBA) stdcolor.NRGBA {
	a := float64(src.A) / 255
	mix := func(s, d uint8) float64 {
		return a*LinearFromSRGB(float64(s)/255) + (1-a)*LinearFromSRGB(float64(d)/255)
	}
	r, g, b := quantizeLinear(mix(src.R, dst.R), mix(src.G, dst.G), mix(src.B, dst.B))
	return stdcolor.NRGBA{R: r, G: g, B: b, A: 0xff}
}
