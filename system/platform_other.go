//go:build !darwin

package system

import "github.com/vibrantgio/theme/tokens"

// platformColors builds the platform's own colour set for an appearance:
// the recorded set for the appearance's side, with the accent-following
// rows taken from the colour [PlatformColor] resolves. Where the platform
// reports no colour to derive from — Windows and Linux today — the recorded
// rows stand, which is the platform's own blue. macOS reads the whole set
// off AppKit instead (platform_darwin.go).
//
// It is deliberately independent of the theme bridge's palette precedence:
// this set is the platform's answer, not the application's brand, so
// [WithSeed] and [WithPalette] do not reach it.
func platformColors(a Appearance) tokens.PlatformColors {
	p := tokens.PlatformLight
	if a.Dark {
		p = tokens.PlatformDark
	}
	if accent, ok := PlatformColor(a); ok {
		p = p.WithAccent(accent)
	}
	return p
}
