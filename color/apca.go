// APCA (Accessible Perceptual Contrast Algorithm), the one contrast metric
// the palette is gated on and reported in. This implements the published
// APCA-W3 version 0.0.98G-4g formula — the version documented in the apca-w3
// reference implementation (github.com/Myndex/apca-w3) and used by WCAG 3
// drafts — with its standard constants, verbatim.
package color

import (
	stdcolor "image/color"
	"math"
)

// The APCA-W3 0.0.98G-4g constants, as published in apca-w3. Names follow
// the reference implementation so the two can be compared line for line.
const (
	// Estimated screen luminance: a simple 2.4-exponent power curve per
	// channel (no piecewise toe — APCA models flare separately) with
	// APCA's own sRGB coefficients.
	apcaTRC = 2.4
	apcaRco = 0.2126729
	apcaGco = 0.7151522
	apcaBco = 0.0721750

	// Soft black clamp: luminances below blkThrs are raised toward it,
	// modeling screen flare and the eye's response near black.
	apcaBlkThrs = 0.022
	apcaBlkClmp = 1.414

	// Contrast-curve exponents: normal polarity (dark text on light
	// background) uses normBG/normTXT, reverse polarity (light on dark)
	// uses revBG/revTXT.
	apcaNormBG  = 0.56
	apcaNormTXT = 0.57
	apcaRevTXT  = 0.62
	apcaRevBG   = 0.65

	// Output scaling and low-contrast handling.
	apcaScale    = 1.14
	apcaOffset   = 0.027
	apcaLoClip   = 0.1
	apcaDeltaMin = 0.0005
)

// apcaLuminance is APCA's estimated screen luminance Ys of an sRGB colour:
// per-channel 2.4-exponent linearization weighted by APCA's coefficients.
// This is deliberately not [RelativeLuminance] — APCA specifies its own
// transfer curve and coefficients.
func apcaLuminance(c stdcolor.NRGBA) float64 {
	ch := func(v uint8) float64 {
		return math.Pow(float64(v)/255.0, apcaTRC)
	}
	return apcaRco*ch(c.R) + apcaGco*ch(c.G) + apcaBco*ch(c.B)
}

// apcaClamp applies the soft black clamp to an estimated screen luminance.
func apcaClamp(y float64) float64 {
	if y < apcaBlkThrs {
		return y + math.Pow(apcaBlkThrs-y, apcaBlkClmp)
	}
	return y
}

// APCA returns the APCA-W3 (0.0.98G-4g) lightness contrast Lc between text
// and background colours, in the published signed convention: positive for
// dark text on a light background, negative for light text on a dark
// background. Alpha is ignored (NRGBA channels are non-premultiplied);
// pairs too close to distinguish return 0.
//
// The sign is polarity and not strength, so a floor is a floor on the
// magnitude: [Magnitude] answers that, and the theme's floors are stated
// in it. APCA's own published levels are Lc 75 for body text and Lc 45 for
// a mark that carries meaning without being text, each with a higher
// increased-contrast step.
func APCA(text, background stdcolor.NRGBA) float64 {
	ytxt := apcaClamp(apcaLuminance(text))
	ybg := apcaClamp(apcaLuminance(background))
	if math.Abs(ybg-ytxt) < apcaDeltaMin {
		return 0
	}
	var sapc float64
	if ybg > ytxt { // normal polarity: dark text on light background
		sapc = (math.Pow(ybg, apcaNormBG) - math.Pow(ytxt, apcaNormTXT)) * apcaScale
		if sapc < apcaLoClip {
			return 0
		}
		return (sapc - apcaOffset) * 100
	}
	// reverse polarity: light text on dark background
	sapc = (math.Pow(ybg, apcaRevBG) - math.Pow(ytxt, apcaRevTXT)) * apcaScale
	if sapc > -apcaLoClip {
		return 0
	}
	return (sapc + apcaOffset) * 100
}

// Magnitude returns |Lc|: [APCA]'s reading with the polarity dropped, which
// is what a floor is compared against. A pairing is legible or not by how
// far apart it reads, not by which of the two is the lighter.
func Magnitude(text, background stdcolor.NRGBA) float64 {
	return math.Abs(APCA(text, background))
}

// BestOn returns the candidate foreground that reads furthest over fill:
// the one with the greatest [Magnitude], the first of equals when two tie.
// Candidates are offered in preference order, so a fill that two colours
// serve equally keeps the one its derivation named first.
//
// This is how every on-colour in the palette is chosen. Measuring both ends
// of the available axis and keeping the better one — rather than taking one
// end while it clears a floor — is what puts white on a saturated mid-tone
// blue, where a luminance-only measure ranks black above it and every eye
// disagrees.
//
// No candidates yields the zero NRGBA; a caller with none to offer has
// nothing to choose.
func BestOn(fill stdcolor.NRGBA, candidates ...stdcolor.NRGBA) stdcolor.NRGBA {
	var best stdcolor.NRGBA
	bestLc := -1.0
	for _, c := range candidates {
		if lc := Magnitude(c, fill); lc > bestLc {
			best, bestLc = c, lc
		}
	}
	return best
}
