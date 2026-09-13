// sRGB relative luminance: how much light a colour puts out, which is what
// the palette's lightness questions are asked in — which side of the scheme
// a surface sits on, and how far a seam stands from the fills it parts.
// Contrast is not asked here: the contrast metric is APCA (apca.go).
package color

import (
	stdcolor "image/color"
	"math"
)

// RelativeLuminance returns the sRGB relative luminance of an image/color
// NRGBA value, in [0,1] — 0 for black, 1 for white. The alpha channel is
// ignored: NRGBA is non-premultiplied, so the colour channels are the
// colour regardless of coverage.
//
// This is a luminance and not a contrast measure. It answers questions
// about light output — is this scheme's background above or below the
// middle of the axis, what luminance must a hairline land at to stand a
// measured distance from two fills — and nothing about whether a reader
// can resolve one colour laid on another, which is [APCA]'s question
// alone.
//
// The linearization uses the piecewise threshold V ≤ 0.03928 rather than
// the 0.04045 of the standard sRGB EOTF ([LinearFromSRGB]); for 8-bit
// channels the two never straddle a representable value, so the results
// are identical. The coefficients are the four-significant-digit row, not
// the seven-digit Lindbloom row in linearY.
func RelativeLuminance(c stdcolor.NRGBA) float64 {
	lin := func(V8 uint8) float64 {
		V := float64(V8) / 255.0
		if V <= 0.03928 {
			return V / 12.92
		}
		return math.Pow((V+0.055)/1.055, 2.4)
	}
	return 0.2126*lin(c.R) + 0.7152*lin(c.G) + 0.0722*lin(c.B)
}
