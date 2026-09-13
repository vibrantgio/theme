package theme

import (
	"github.com/reactivego/rx"
	"github.com/vibrantgio/theme/tokens"
)

// Theme carries an rx.Observable for each token category. Consumers subscribe
// to the fields they need and react to updates without polling.
type Theme struct {
	// Platform carries the platform's colour set — AppKit's semantic
	// colours under the emitted appearance, with the accent rows following
	// the theme colour. It is the theme's only colour vocabulary;
	// theme/system builds it.
	Platform rx.Observable[tokens.PlatformColors]

	Typography rx.Observable[tokens.Typography]
	Density    rx.Observable[tokens.Density]
	Motion     rx.Observable[tokens.MotionScale]
	Spacing    rx.Observable[tokens.SpacingScale]
	Radius     rx.Observable[tokens.RadiusScale]
	Elevation  rx.Observable[tokens.ElevationScale]
}

// Default returns a Theme whose every field emits the package-level default
// value once and then completes. It is the canonical starting point for static
// or test scenarios that do not need live token switching.
func Default() Theme {
	return Theme{
		Platform:   rx.Of(tokens.PlatformLight),
		Typography: rx.Of(tokens.DefaultTypography),
		Density:    rx.Of(tokens.Comfortable),
		Motion:     rx.Of(tokens.Motion),
		Spacing:    rx.Of(tokens.Spacing),
		Radius:     rx.Of(tokens.Radius),
		Elevation:  rx.Of(tokens.Elevation),
	}
}
