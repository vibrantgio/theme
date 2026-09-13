// Package color is the colour mathematics the theme measures and composites
// with. It converts between sRGB and the perceptual spaces that work is done
// in — no colour values live here, only the mathematics.
//
// Lightness, hue and chroma come from OKLCh: OKLab, OKLCh and their inverses
// live in oklab.go, after Björn Ottosson's reference formulation.
//
// Out-of-gamut requests are gamut mapped, never channel-clamped: a
// conversion that produces 8-bit sRGB reduces OKLCh chroma at constant
// lightness and constant hue until the colour fits (gamut.go). The
// float-level converters stay raw: out-of-gamut input yields out-of-range
// channels, documented per function.
//
// The package name collides with image/color in an importer's import list;
// alias one of them — theme code aliases the standard library one.
package color
