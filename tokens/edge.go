// What a measured edge is: the line around a floating surface, its corner
// pixel and the rim inside it.
package tokens

import "image/color"

// Edge is the one-pixel line the platform draws around a floating surface,
// outside its bounds, with the rim the platform draws inside that line.
//
// Each value is a coverage and not a fill: the edge is black at an alpha and
// the rim white at an alpha, and each lands as that alpha over whatever lies
// beneath it. The edge's coverage is read with the shadow under it taken
// out, so a caller that draws the shadow first and the edge over it lands
// the platform's pixel.
//
// The zero value draws nothing: a surface with no edge answers it.
type Edge struct {
	// Side is the edge's coverage along each of the four sides.
	Side color.NRGBA
	// Corner is the edge's coverage at its corner pixel, the one diagonal
	// pixel where two sides meet. It is lighter than Side.
	Corner color.NRGBA
	// Rim holds the rim's coverage per pixel inward from the edge, the
	// outermost first. The rim is square at the corners. A zero entry draws
	// nothing, so a one-pixel rim leaves the second entry zero.
	Rim [2]color.NRGBA
}
