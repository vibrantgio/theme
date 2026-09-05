package tokens

import (
	"fmt"
	"image/color"
)

// ColorScale holds the eleven Tailwind shade stops for one hue family (50–950).
type ColorScale struct {
	C50, C100, C200, C300, C400, C500, C600, C700, C800, C900, C950 color.NRGBA
}

// Optional named palettes, taken verbatim from the Tailwind CSS v3 default
// config. They are exactly that — palettes an application may reach for by
// name — and no part of the semantic layer resolves from them: every role
// ramp, pin and semantic colour derives from a seed (see FromSeed). The
// Tailwind values may survive only in this arrangement, never behind a role
// name.
var (
	Slate = ColorScale{
		C50:  color.NRGBA{0xf8, 0xfa, 0xfc, 0xff},
		C100: color.NRGBA{0xf1, 0xf5, 0xf9, 0xff},
		C200: color.NRGBA{0xe2, 0xe8, 0xf0, 0xff},
		C300: color.NRGBA{0xcb, 0xd5, 0xe1, 0xff},
		C400: color.NRGBA{0x94, 0xa3, 0xb8, 0xff},
		C500: color.NRGBA{0x64, 0x74, 0x8b, 0xff},
		C600: color.NRGBA{0x47, 0x55, 0x69, 0xff},
		C700: color.NRGBA{0x33, 0x41, 0x55, 0xff},
		C800: color.NRGBA{0x1e, 0x29, 0x3b, 0xff},
		C900: color.NRGBA{0x0f, 0x17, 0x2a, 0xff},
		C950: color.NRGBA{0x02, 0x06, 0x17, 0xff},
	}
	Blue = ColorScale{
		C50:  color.NRGBA{0xef, 0xf6, 0xff, 0xff},
		C100: color.NRGBA{0xdb, 0xea, 0xfe, 0xff},
		C200: color.NRGBA{0xbf, 0xdb, 0xfe, 0xff},
		C300: color.NRGBA{0x93, 0xc5, 0xfd, 0xff},
		C400: color.NRGBA{0x60, 0xa5, 0xfa, 0xff},
		C500: color.NRGBA{0x3b, 0x82, 0xf6, 0xff},
		C600: color.NRGBA{0x25, 0x63, 0xeb, 0xff},
		C700: color.NRGBA{0x1d, 0x4e, 0xd8, 0xff},
		C800: color.NRGBA{0x1e, 0x40, 0xaf, 0xff},
		C900: color.NRGBA{0x1e, 0x3a, 0x8a, 0xff},
		C950: color.NRGBA{0x17, 0x25, 0x54, 0xff},
	}
	Red = ColorScale{
		C50:  color.NRGBA{0xfe, 0xf2, 0xf2, 0xff},
		C100: color.NRGBA{0xfe, 0xe2, 0xe2, 0xff},
		C200: color.NRGBA{0xfe, 0xca, 0xca, 0xff},
		C300: color.NRGBA{0xfc, 0xa5, 0xa5, 0xff},
		C400: color.NRGBA{0xf8, 0x71, 0x71, 0xff},
		C500: color.NRGBA{0xef, 0x44, 0x44, 0xff},
		C600: color.NRGBA{0xdc, 0x26, 0x26, 0xff},
		C700: color.NRGBA{0xb9, 0x1c, 0x1c, 0xff},
		C800: color.NRGBA{0x99, 0x1b, 0x1b, 0xff},
		C900: color.NRGBA{0x7f, 0x1d, 0x1d, 0xff},
		C950: color.NRGBA{0x45, 0x0a, 0x0a, 0xff},
	}
)

// White and Black are the two ends of the tonal axis. Unlike the named
// palettes above they are part of the semantic layer: they are the pair
// FromSeed chooses a light-scheme on-colour from, by measuring both over
// the base rather than assuming one. White carries almost every accent —
// a pinned base sits deep enough to read white text — and Black is what an
// accent light enough to lose its white foreground takes instead.
var (
	White = color.NRGBA{0xff, 0xff, 0xff, 0xff}
	Black = color.NRGBA{0x00, 0x00, 0x00, 0xff}
)

// Ramp is one colour role's nine-step functional ramp. Steps run
// 100–900 in hundreds and the step number carries the meaning: 100–300 are
// tinted fills, hovers and subtle borders, 500 is the mid-value reference,
// 700–900 are text over tinted fills and pressed states. Index i holds step
// (i+1)*100; use Step to address a ramp by its step number.
//
// Light and dark ramps are paired scales, not two role tables: the same step
// keeps the same job in both modes, so a component asking for neutral 200
// gets a light card on a light surface and a dark card on a dark one.
type Ramp [9]color.NRGBA

// Step returns the colour at step n, where n is one of 100, 200, … 900.
// Any other n is a programming error and panics.
func (r Ramp) Step(n int) color.NRGBA {
	if n < 100 || n > 900 || n%100 != 0 {
		panic(fmt.Sprintf("tokens: Ramp.Step(%d): step must be 100–900 in hundreds", n))
	}
	return r[n/100-1]
}

// RampSet holds the colour-role ramps. Neutral carries
// every surface, border and text shade; the accent ramps carry each role's
// tints and text shades, while the role's base colour is pinned separately
// on ColorTokens (see ColorTokens.Primary).
//
// Error, Success, Warning and Info are the semantic status roles. Unlike
// Primary, Secondary and Tertiary they do not rotate with the seed — a
// purple "success" would be useless, and an "info" wearing the accent says
// whatever the brand happens to say — so each is anchored at a fixed hue
// and chroma the seed may only tint, by at most a few degrees of hue. See
// the seed.go file header for the four measurements and the bound.
type RampSet struct {
	Neutral   Ramp
	Primary   Ramp
	Secondary Ramp
	Tertiary  Ramp
	Error     Ramp
	Success   Ramp
	Warning   Ramp
	Info      Ramp
}

// ColorTokens holds a scheme's whole colour vocabulary: the nine-step
// functional ramps, the pinned role bases, and a thin semantic layer resolved
// from ramp steps. Each "On" field is the recommended text/icon colour
// rendered on top of its companion pinned base.
//
// The pinned bases exist because a brand seed rarely sits on the shared
// lightness scale: the pin reproduces the role's base colour exactly instead
// of reading a lightened approximation off a ramp step. Dark mode pins a
// dark-appropriate base rather than reusing the light pin.
//
// There are no MD3-only alias fields. Where the M3 role is a fixed resolution
// off the neutral ramp it is reachable by asking the ramp directly:
// OnBackground is Text, OnSurface Ramps.Neutral.Step(900), SurfaceVariant
// Step(300), Outline Step(500) — or FocusRing(). Where naming a step would
// state one colour and two measurements, the token is a derivation instead:
// [ColorTokens.OutlineVariant] and [ColorTokens.OnSurfaceVariant] choose their
// step against a floor, and SurfaceContainerLow is a level, so
// [ColorTokens.SurfaceAt] answers it.
type ColorTokens struct {
	// Ramps holds the functional ramps, fully populated: nine steps per
	// role, generated on the shared lightness scale by FromSeed.
	Ramps RampSet

	// Pinned accent bases and their on-colours: the solid fills.
	Primary     color.NRGBA // pinned primary base — in a light scheme, the lifted seed
	OnPrimary   color.NRGBA // text/icon over Primary
	Secondary   color.NRGBA // pinned secondary base
	OnSecondary color.NRGBA // text/icon over Secondary
	Tertiary    color.NRGBA // pinned tertiary base
	OnTertiary  color.NRGBA // text/icon over Tertiary
	Error       color.NRGBA // pinned error base
	OnError     color.NRGBA // text/icon over Error
	Success     color.NRGBA // pinned success base
	OnSuccess   color.NRGBA // text/icon over Success
	Warning     color.NRGBA // pinned warning base
	OnWarning   color.NRGBA // text/icon over Warning
	Info        color.NRGBA // pinned info base
	OnInfo      color.NRGBA // text/icon over Info

	// The thin semantic layer. Background and Text are pins; Surface and
	// Divider resolve from Neutral ramp steps at construction.
	Background color.NRGBA // pinned app background
	Text       color.NRGBA // pinned body text over Background
	// Surface is the neutral ramp's step 200 — one step off the app's own
	// background. It is a RAMP ALIAS, not a level: the elevation is anchored
	// on the Background pin and placed in CIELAB L*, so which level this step
	// happens to carry depends on the scheme (light furniture wears it; dark
	// raised surfaces do). Ask [ColorTokens.SurfaceAt] for a level.
	Surface color.NRGBA
	// Divider is the subtle border / separator — Ramps.Neutral.Step(300),
	// except in the high-contrast variant, which resolves it from the
	// strong-border step 500 (see FromSeedHighContrast).
	Divider color.NRGBA

	// The inverse pair: a surface deliberately on the wrong side of the
	// scheme, dark in a light scheme and light in a dark one, with the
	// on-colour to read on it. It is what a transient message stands on —
	// a message that can appear over any surface separates by being the
	// one thing built from the opposite scheme, not by out-elevating what
	// it covers.
	//
	// Both resolve from the *counterpart* scheme's neutral ramp, which is
	// the whole of the derivation: a light scheme's inverse pair is the
	// dark scheme's Surface and Text, and a dark scheme's is the light
	// scheme's. So the pair carries the same measured separation the
	// counterpart scheme's own reading pair does — nothing about it is an
	// approximation — and it re-derives with the seed and the
	// high-contrast variant like every other role.
	InverseSurface   color.NRGBA // counterpart Neutral.Step(200)
	OnInverseSurface color.NRGBA // text/icon over it — counterpart Neutral.Step(900)

	// Highlight is the reserved highlighter: the fill marking content the
	// user was brought to. It is not a colour role and reports no status,
	// so its hue is reserved outside the role table and no status hue may
	// serve it; highlight.go carries the reservation and the distances
	// that hold it. This is the fill resolved against the Background pin,
	// the surface content stands on;
	// [ColorTokens.HighlightOn] answers for any other surface.
	Highlight color.NRGBA
}

// resolveAliases fills every field defined as a resolution of a ramp step:
// Surface and Divider off this scheme's own neutral ramp, the inverse
// pair off the counterpart scheme's, and Highlight off the reserved hue
// against the Background pin. dividerStep is the Neutral step
// Divider resolves from: 300 in the default derivation, 500 in the
// high-contrast variant. counterpart is the other scheme's neutral ramp —
// the dark one while building the light scheme and the light one while
// building the dark. Constructing tokens through it is what keeps each
// field byte-identical to its documented resolution; FromSeed and
// FromSeedHighContrast build both schemes through it.
func resolveAliases(t ColorTokens, dividerStep int, counterpart Ramp) ColorTokens {
	t.Surface = t.Ramps.Neutral.Step(200)
	t.Divider = t.Ramps.Neutral.Step(dividerStep)
	t.InverseSurface = counterpart.Step(200)
	t.OnInverseSurface = counterpart.Step(900)
	t.Highlight = t.HighlightOn(t.Background)
	return t
}
