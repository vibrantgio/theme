// Seed-derived palettes. FromSeed turns one brand colour into ADR-007's
// complete paired light and dark ramp sets, and DefaultLight/DefaultDark
// are FromSeed of the default seed.
//
// Derivation rules, and where each comes from:
//
//   - The shared lightness scale is ADR-007's: CIELAB L* per step, measured
//     by the D0.1 spike from the Claude Design reference project's own
//     ramps. Light 100–900 = 97, 92, 85, 74, 63, 51, 39, 28, 6 — the
//     measured 900 was L* 18, but D2.4's APCA gate deepened it: APCA's soft
//     black clamp caps even pure black near Lc 92 over the L* 92 step-200
//     ground, so Lc ≥ 90 needs the 900 stop at L* 6 (min Lc 90.7 across the
//     seven default ramps; L* 18 measured Lc 85–87). The dark scale is the
//     paired scale measured from the same source's dark column (ADR-007's
//     evidence table): 8, 13, 19, 30, 65, —, 82, —, 94; the 600 and 800
//     stops the table has no surface for are interpolated to 74 and 88 (the
//     dark 900 already clears the gate at Lc 93–96, so it is untouched).
//     Both scales are swept at constant OKLCh hue and chroma via
//     color.Tone, which gamut-maps by chroma reduction (ADR-002).
//
//   - The neutral ramps carry no hue: chroma 0.000, measured 0.0000 at
//     every step in both modes. They used to carry the seed at chroma
//     0.010 — MD3's neutral 4 — which tinted every surface in the system
//     with the brand, and the platform's own opaque window fills are flat
//     greys (the stored macOS reference reads (30,30,30) for both the
//     Notes reading surface and the Voice Memos panes). A surface is not
//     where a brand belongs, so the tint is gone: light Background,
//     Surface, Divider and Text now measure #f6f6f6, #e8e8e8, #d4d4d4 and
//     #131313, and their dark counterparts #181818, #222222, #2e2e2e and
//     #eeeeee. The neutral scale's lightnesses are untouched, so every
//     contrast gate reads what it read before.
//
//   - Accent chroma is a dial, not the seed's own measurement. The dial is
//     0.22 OKLCh chroma: every brand colour between "grey" (chroma 0.020)
//     and the dial is rendered at the dial, one already past it keeps its
//     own chroma, and a grey one stays grey — a dial must not invent a hue
//     where the brand has none. Before this, primary took the seed's
//     measured chroma unchanged, so a washed-out brand colour made a
//     washed-out palette; the canonical seed #6750A4 measures chroma
//     0.1305, and every accent surface in the system inherited exactly
//     that. Secondary and tertiary keep MD3's ratios to primary — 16 : 48
//     and 24 : 48 of material-color-utilities' chromas — so they follow
//     the dial rather than being set by hand: at the dial they land at
//     0.073 and 0.110. Tertiary keeps MD3's +60° hue rotation. The
//     conversion anchor for reading MD3's numbers on the OKLCh chroma
//     axis: #6750A4 has HCT chroma 48 and measures OKLCh 0.1305, so one
//     HCT chroma unit ≈ 0.00272 OKLCh chroma, and the dial is HCT ≈ 81.
//
//     What the dial buys, measured on the canonical seed: the light
//     primary base goes 0.1305 → 0.2196 (#6750a4 → #723ad4) at the same
//     depth and the same white on-colour, so the pair's contrast is
//     unmoved (Lc 86.9 → 86.5, WCAG 6.44 → 6.46). The ramp steps the
//     containers and selection fills resolve from follow the same dial:
//     light 500 → 800 go 0.131 → 0.176/0.220/0.220/0.220, dark 200 → 400
//     go 0.130 → 0.150/0.181/0.220. Where a step is already against the
//     sRGB boundary the dial buys nothing and the step does not move —
//     light 100–400 and dark 700–900 are gamut-limited at this hue, the
//     dark primary base among them (see the pins below).
//
//   - The three status roles are hue-fixed, not seed-derived. A semantic
//     colour must not rotate with the brand: a purple "success" says
//     nothing. Each takes the OKLCh hue and chroma of a canonical Material
//     colour, measured with this module's own converters and recorded here
//     (measurements from F4.6; the error pair predates it, from D2.2):
//
//       error    hue  28.7°, chroma 0.178  — MD3's canonical error base
//                                            #B3261E (its "hue 25,
//                                            chroma 84"), L* 39.7
//       success  hue 144.2°, chroma 0.162  — Material Green 500 #4CAF50,
//                                            L* 63.98
//       warning  hue  84.9°, chroma 0.172  — Material Amber 500 #FFC107,
//                                            L* 81.52
//
//     The palette anchor of each family is its 500 shade, so that is what
//     is measured; only the hue and chroma are taken, since the depths come
//     from the shared lightness scale like every other role. The three land
//     56–59° apart on the OKLCh hue circle — far enough that warning is not
//     read as error at a glance, which is the whole point of a status
//     colour.
//
//   - Pins. The light primary base is the seed at its own hue and CIELAB
//     depth with the accent dial applied to its chroma (ADR-007: "the seed
//     sits deep, so bases are pins" — reading it off the ramp would lighten
//     it); only its alpha is forced opaque. A brand colour the dial leaves
//     alone comes back byte-for-byte, so a palette seeded from a desktop's
//     accent colour still matches that accent exactly — six of the nine
//     macOS system colours are already past the dial or below the grey
//     threshold and reproduce exactly. Because the base is what a palette
//     publishes, the derivation reads the accent family's hue and chroma
//     back off the base rather than off the seed, and the dial is a
//     projection: deriving a palette from its own primary base reproduces
//     that palette byte-for-byte, which is what lets a serialized theme
//     name one colour and be rebuilt from it. The other light
//     bases are their role's hue and chroma at tone 40, the depth MD3 pins
//     accent bases at and the depth the default seed itself sits at
//     (L* 40.08). Dark bases are the same hue and chroma re-toned to L* 82
//     — the dark scale's step-700 depth, right beside MD3's dark
//     accent-base tone 80, making the dark pin byte-identical to its ramp's
//     step 700. The D0.1 spike sat them at L* 65, the step-500 depth
//     reproducing ADR-007's recorded dark fill #a690ea, but D2.4's APCA
//     gate showed an L* 65 mid-tone is a ground no text can reach Lc 60
//     over (black tops out near Lc 52, white near 57), so the pins moved up
//     two rungs — the default seed's dark primary is now #d0c4ff — and the
//     solid state walk still lands on exact rungs (hover 800, pressed 900).
//     L* 82 is also the shallowest depth the increased-contrast variant's
//     Lc ≥ 75 floor allows (L* 80 reaches only 73.5 against pure black),
//     so it is pinned from both directions. The accent dial buys the dark
//     primary base nothing: at L* 82 the seed's hue holds chroma 0.0822 in
//     sRGB and no more, and it was already there. The dark base is the one
//     accent surface in the system the dial cannot reach — everything the
//     dark scheme fills with a primary ramp step does move.
//
//   - The inverse pair. Each scheme's InverseSurface and OnInverseSurface
//     are the *other* scheme's Surface and Text — its neutral ramp's steps
//     200 and 900 — so a light scheme's inverse chip is dark and a dark
//     scheme's is light, and the pair's separation is the counterpart
//     scheme's own body-text separation rather than a second measurement
//     (WCAG 13.75:1 light, 15.06:1 dark on the default seed; the
//     high-contrast variant widens both to 15.99:1 and 17.11:1 by
//     deepening the 900 stop the on-colour reads off). Both
//     schemes are derived in one pass here, so neither needs anything the
//     other has not already computed.
//
//   - On-colours are measured, not assumed. Each pinned base is read in the
//     ink that reaches 4.5:1 over it — WCAG AA for body text — with the
//     scheme's usual ink preferred and the other end of the tonal axis
//     taken when the usual one falls short (see onColour). In the light
//     scheme the pair on offer is White and Black; in the dark scheme it is
//     the role's own step-100 depth and White.
//
//     Almost every base keeps the ink it always had, and the rule is a
//     no-op there: a light base at tone 40 carries White at Lc ≥ 85, WCAG
//     ≈ 6.4:1, and a dark base at L* 82 carries its deep ink at Lc ≥ 73,
//     WCAG ≈ 11:1. What the rule is for is the one base pinned to no
//     depth — the primary base is the brand colour itself, so a light brand
//     colour used to come back under white text at as little as 2.1:1. Now
//     its ink flips and the colour does not move: an accent stays true to
//     the seed, the way the design language pairs a high tone with a dark
//     ink everywhere else. Across a 411-seed sweep 266 of the light
//     schemes' primary inks flip — nothing else in either scheme does, the
//     other bases being pinned to depths their usual ink clears — and no
//     pinned pairing any seed produces measures under the floor. The
//     container and fill pairings the ramps carry are untouched by all of
//     this and stay where they were: a ramp step is realized at a fixed
//     depth, so its 700-and-900 text over its 100 and 200 grounds measures
//     the same 5.4:1 and up whatever the seed. The APCA gate above still
//     enforces its own Lc ≥ 60 on top of all of it.
//
//     The state walk under a solid fill is not part of this, and cannot
//     be: a fill walks toward its ramp's 900 end whichever depth its pin
//     sits at (see states.go), so on a mid-depth accent one ink reads at
//     rest and the other reads pressed, and one token cannot be both.
//     Choosing each ink for the whole walk instead of for the resting pair
//     was measured over the same sweep and buys four pairings back under
//     the pointer at the cost of eighty-three at rest — the wrong trade,
//     since resting is where a surface is read. The walk keeps its own
//     rule, and what it needs is a walk that knows which way its pin
//     faces.
//
// FromSeedHighContrast (task E3.3) derives the increased-contrast variant
// from the same seed by the same machinery — it is a FromSeed option, not a
// third hand-written scheme. Three widenings, each computed against the
// APCA gate rather than guessed:
//
//   - The 700 text step deepens to the default scale's 900 depth in both
//     modes — light 700 L* 39 → 6, dark 700 L* 82 → 94 — so 700 text meets
//     the same Lc ≥ 90 bar the default asks only of 900 (light min Lc 90.7,
//     dark 93.0 across the seven ramps; APCA's soft black clamp caps lighter
//     choices below 90, the same wall D2.4 hit). The 800 and 900 stops
//     slide outward — light 3 and 0, dark 97 and 100 — keeping the ladder
//     strictly monotonic and the 900 gate clear with margin (light Lc 92.3,
//     dark 104.4). Steps 100–600 are the default scale unchanged: the
//     grounds stay, the text pulls away.
//
//   - Divider resolves from Neutral step 500 instead of 300: the separator
//     jumps from the subtle-border rung to ADR-007's strong-border rung.
//
//   - Each pinned base's on-colour is pushed further from its base. The
//     dark pins' on-colours drop from their ramp's step 100 (L* 8, Lc ≈ 74
//     — just under the variant's Lc ≥ 75 floor) to tone 0, the scale's
//     floor (Lc ≥ 76.3). The light pins keep White wherever White is the
//     better of the two ends: it is already the far end of the axis and
//     already clears the floor (Lc ≥ 85.7), so the pins do not move — the
//     light primary base is the same lifted seed FromSeed pins, the same
//     contract.
//
//     The on-colour rule follows the variant to a stricter floor: the
//     light ink stands only while it reaches 7:1 rather than 4.5:1, so the
//     variant questions an ink the default is satisfied with. What it can
//     do with the answer is bounded by the axis, which has no ink further
//     out than its two ends — where neither reaches the floor the better
//     of the two stands in both derivations, so the variant's flipped set
//     is the default's plus the sliver where the light ink clears AA and
//     the dark ink still reads higher. Every variant pairing therefore
//     measures at least what the default's does, which is the property its
//     gate holds.
package tokens

import (
	stdcolor "image/color"

	"github.com/vibrantgio/theme/color"
)

// lightTones and darkTones are the shared perceptual lightness scale:
// CIELAB L* for steps 100–900, light and paired dark, per ADR-007. Index i
// holds step (i+1)*100, matching Ramp.
var (
	lightTones = [9]int{97, 92, 85, 74, 63, 51, 39, 28, 6}
	darkTones  = [9]int{8, 13, 19, 30, 65, 74, 82, 88, 94}
)

// hcLightTones and hcDarkTones are the high-contrast variant's scales:
// steps 100–600 are the default scale unchanged, 700 deepens to the default
// 900 depth so 700 text meets the Lc ≥ 90 bar, and 800/900 slide to the
// axis ends to keep the ladder strictly monotonic. See the file header.
var (
	hcLightTones = [9]int{97, 92, 85, 74, 63, 51, 6, 3, 0}
	hcDarkTones  = [9]int{8, 13, 19, 30, 65, 74, 94, 97, 100}
)

// Accent-derivation constants; see the file header for provenance.
const (
	neutralChroma    = 0.000     // the neutral ramps carry no hue at all
	accentChroma     = 0.22      // the accent dial; see the file header
	accentChromaSnap = 0.005     // a chroma this close to the dial counts as at it
	greyChroma       = 0.020     // at or under this a brand colour is grey and stays grey
	secondaryShare   = 1.0 / 3.0 // MD3's secondary : primary chroma ratio, 16 : 48
	tertiaryShare    = 0.5       // MD3's tertiary : primary chroma ratio, 24 : 48
	tertiaryHueShift = 60        // MD3: tertiary is the seed hue rotated +60°
	errorHue         = 28.7      // OKLCh hue of MD3's error base #B3261E
	errorChroma      = 0.178     // OKLCh chroma of #B3261E
	successHue       = 144.2     // OKLCh hue of Material Green 500 #4CAF50
	successChroma    = 0.162     // OKLCh chroma of #4CAF50
	warningHue       = 84.9      // OKLCh hue of Material Amber 500 #FFC107
	warningChroma    = 0.172     // OKLCh chroma of #FFC107
	lightPinTone     = 40        // MD3's accent-base tone; the default seed's own depth
	darkPinTone      = 82        // the dark scale's step-700 L*; D2.4 raised it from the
	// spike's 65 — no on-colour reaches Lc 60 over an L* 65 mid-tone
	darkOnTone   = 8   // dark pins' on-colour depth: the dark scale's step-100 L*
	hcDarkOnTone = 0   // high contrast pushes the dark on-colours to the axis floor
	onFloor      = 4.5 // WCAG AA body text: the ratio an on-colour has to reach
	hcOnFloor    = 7.0 // the increased-contrast variant asks AAA of the same pair
)

// derivation is the knob set that separates FromSeed from its high-contrast
// variant: the two lightness scales, the ramp step Divider resolves from,
// and the CIELAB L* the dark pins' on-colours are realized at. Everything
// else — hues, chromas, pin depths, the lifted light primary base — is
// shared, which is what makes the variant a FromSeed option rather than a
// third hand-written scheme.
type derivation struct {
	lightTones, darkTones [9]int
	dividerStep           int     // ramp step Divider resolves from
	darkOnTone            int     // L* of the dark pins' on-colours
	onFloor               float64 // the ratio an on-colour has to reach over its base
}

var (
	defaultDerivation = derivation{lightTones, darkTones, 300, darkOnTone, onFloor}
	hcDerivation      = derivation{hcLightTones, hcDarkTones, 500, hcDarkOnTone, hcOnFloor}
)

// onColour picks the ink one pinned base is read in. The preferred ink
// stands while it reaches the floor over that base; below it the ink flips
// to the other end of the tonal axis, unless that end reads worse still —
// a base no ink can carry keeps the better of the two rather than the
// darker of the two.
//
// Which is the whole of the rule the light scheme needed. Its bases used to
// take White unconditionally, which is right for a base at the depth MD3
// pins accents at and wrong for the one base that is not pinned to a depth
// at all: the primary base is the brand colour itself, and a light brand
// colour under white text measured as little as 2.1:1 where the floor is
// 4.5. Flipping the ink rather than deepening the colour is what keeps a
// palette true to the colour it was seeded with.
//
// The two ends are pure White and pure Black, and that is load-bearing
// rather than tidy: over any colour whatever, the better of white and black
// reaches 4.58:1, so no seed can produce a pinned pairing under the floor.
// An ink one rung short of the axis end — the ramp's own 900 stop, say —
// gives that guarantee up (it bottoms out at 4.31:1 across a seed sweep),
// and an on-colour is text, where a tint buys nothing anyway.
func onColour(base, preferred, other stdcolor.NRGBA, floor float64) stdcolor.NRGBA {
	got := color.ContrastRatio(preferred, base)
	if got >= floor {
		return preferred
	}
	if color.ContrastRatio(other, base) > got {
		return other
	}
	return preferred
}

// liftChroma turns a brand colour's own measured chroma into the chroma
// the accent family is rendered at. A brand colour under greyChroma is
// grey and comes back untouched — the dial must not invent a hue where the
// brand has none. One already at or past the dial (within accentChromaSnap,
// which absorbs the eight-bit quantization of a realized colour) keeps its
// own chroma, so a vivid brand is never pulled down. Everything between is
// rendered at the dial, which is the whole point: a washed-out brand colour
// does not make a washed-out palette.
//
// The function is a projection — every value it returns is a value it
// leaves alone — and that is load-bearing, not tidiness: re-deriving a
// palette from its own primary base has to land on the same palette, or
// nothing downstream can recover a palette from the base it published.
func liftChroma(c float64) float64 {
	if c <= greyChroma || c >= accentChroma+accentChromaSnap {
		return c
	}
	return accentChroma
}

// liftSeed realizes the light primary base: the brand colour's own hue and
// CIELAB depth at liftChroma of its chroma. A brand colour the dial leaves
// alone comes back byte-for-byte, which is what keeps a palette seeded from
// a desktop's accent colour matching that accent exactly; a washed-out one
// comes back at the same hue and depth with the dial's chroma, so nothing
// the contrast gates pin about the base moves.
func liftSeed(seed stdcolor.NRGBA) stdcolor.NRGBA {
	realize := func(c stdcolor.NRGBA) stdcolor.NRGBA {
		tone, _, _ := color.LabFromNRGBA(c)
		_, chroma, hue := color.OKLChFromNRGBA(c)
		return color.NRGBAFromToneChromaHue(tone, liftChroma(chroma), hue)
	}
	// Settle. Realizing a request in eight bits moves the measured hue and
	// depth a hair off what was asked for, so one pass can land somewhere
	// that would realize differently if it were asked again — and then a
	// palette could not be re-derived from the base it published.
	// Re-realizing walks that out. The walk is deterministic over a finite
	// set, so it always closes on a cycle; almost always the cycle is one
	// colour, and where two colours trade places the smallest is taken so
	// that every way into the cycle settles on the same answer. A handful
	// of solver runs, once per palette.
	trail := []stdcolor.NRGBA{realize(seed)}
	for {
		next := realize(trail[len(trail)-1])
		at := -1
		for i, c := range trail {
			if c == next {
				at = i
				break
			}
		}
		if at < 0 {
			trail = append(trail, next)
			continue
		}
		base := trail[at]
		for _, c := range trail[at:] {
			if packRGB(c) < packRGB(base) {
				base = c
			}
		}
		return base
	}
}

// packRGB orders colours for liftSeed's cycle canonicalization; any total
// order does, and this one is the obvious reading of the three bytes.
func packRGB(c stdcolor.NRGBA) uint32 {
	return uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
}

// rampOf sweeps one role's hue and chroma across a lightness scale.
func rampOf(tones [9]int, hue, chroma float64) Ramp {
	var r Ramp
	for i, tone := range tones {
		r[i] = color.Tone(hue, chroma, tone)
	}
	return r
}

// FromSeed derives the complete paired light and dark colour token sets
// from one brand seed: for every role a nine-step ramp on the shared
// lightness scale in both modes — the same step keeps the same job — plus
// the pinned bases, on-colours and semantic layer, per the rules in the
// file header. The light primary base is the seed at its own hue and
// depth with the accent dial on its chroma (alpha forced opaque), and a
// seed the dial leaves alone comes back byte-for-byte; every other value is
// generated. Deriving from a pair's own light primary base reproduces the
// pair exactly, so a palette can be rebuilt from the one colour it names.
//
// DefaultLight and DefaultDark are FromSeed(DefaultSeed). Applications
// re-brand by calling FromSeed with their own colour and handing the pair
// to a theme.
func FromSeed(seed stdcolor.NRGBA) (light, dark ColorTokens) {
	return fromSeed(seed, defaultDerivation)
}

// FromSeedHighContrast derives the increased-contrast variant of FromSeed's
// pair from the same seed: same roles, hues, chromas and pin depths, with
// the tone separation widened where it counts — the 700 text step deepened
// to the default 900 depth (Lc ≥ 90 where the default asks 60), Divider
// resolved from Neutral step 500 instead of 300, and the dark pins'
// on-colours pushed to the tonal axis floor (Lc ≥ 75 over their bases; the
// light pins keep White, which already clears that floor, so the light
// primary base is the same lifted seed FromSeed pins). The full rules are in the
// file header; the derived default-seed variant is recorded in this
// package's high-contrast golden test.
//
// It is the palette theme/system swaps in while the OS reports increased
// contrast (see system.HighContrastVariant).
func FromSeedHighContrast(seed stdcolor.NRGBA) (light, dark ColorTokens) {
	return fromSeed(seed, hcDerivation)
}

// fromSeed is the shared derivation engine behind FromSeed and
// FromSeedHighContrast; d selects which variant's knobs apply.
func fromSeed(seed stdcolor.NRGBA, d derivation) (light, dark ColorTokens) {
	// The primary base is realized first and the whole accent family is
	// measured off it, not off the seed: the base is what a palette
	// publishes, so reading the family's hue and chroma back off the base
	// is what lets a derivation reproduce itself from its own output.
	seed.A = 0xff
	primary := liftSeed(seed)
	_, accent, hue := color.OKLChFromNRGBA(primary)
	accent = liftChroma(accent)

	roles := []struct {
		hue, chroma float64
	}{
		{hue, neutralChroma},
		{hue, accent},
		{hue, accent * secondaryShare},
		{hue + tertiaryHueShift, accent * tertiaryShare},
		{errorHue, errorChroma},
		{successHue, successChroma},
		{warningHue, warningChroma},
	}
	var lr, dr [7]Ramp
	for i, role := range roles {
		lr[i] = rampOf(d.lightTones, role.hue, role.chroma)
		dr[i] = rampOf(d.darkTones, role.hue, role.chroma)
	}
	// The pinned bases and the ink each is read in. Index 0 is the neutral
	// role, which carries surfaces rather than a solid fill and has no pin.
	// Every ink is measured over the base it sits on rather than assumed
	// from the scheme (see onColour); the light scheme's alternative is the
	// far end of the tonal axis and the dark scheme's is White, so in each
	// scheme the pair on offer is the ramp's own dark end and its light one.
	var lightBase, darkBase, lightInk, darkInk [7]stdcolor.NRGBA
	for i := 1; i < len(roles); i++ {
		lightBase[i] = color.Tone(roles[i].hue, roles[i].chroma, lightPinTone)
		darkBase[i] = color.Tone(roles[i].hue, roles[i].chroma, darkPinTone)
	}
	lightBase[1] = primary // the lifted seed, never read off a ramp step
	for i := 1; i < len(roles); i++ {
		deep := color.Tone(roles[i].hue, roles[i].chroma, d.darkOnTone)
		lightInk[i] = onColour(lightBase[i], White, Black, d.onFloor)
		darkInk[i] = onColour(darkBase[i], deep, White, d.onFloor)
	}

	// Each scheme's inverse pair resolves off the other scheme's neutral
	// ramp — both are derived here, so the pair costs nothing and needs no
	// second rule (see ColorTokens.InverseSurface).
	light = resolveAliases(ColorTokens{
		Ramps: RampSet{
			Neutral: lr[0], Primary: lr[1], Secondary: lr[2], Tertiary: lr[3],
			Error: lr[4], Success: lr[5], Warning: lr[6],
		},
		Primary:     lightBase[1], // the lifted seed, never read off a ramp step
		OnPrimary:   lightInk[1],
		Secondary:   lightBase[2],
		OnSecondary: lightInk[2],
		Tertiary:    lightBase[3],
		OnTertiary:  lightInk[3],
		Error:       lightBase[4],
		OnError:     lightInk[4],
		Success:     lightBase[5],
		OnSuccess:   lightInk[5],
		Warning:     lightBase[6],
		OnWarning:   lightInk[6],
		Background:  lr[0].Step(100),
		Text:        lr[0].Step(900),
	}, d.dividerStep, dr[0])
	dark = resolveAliases(ColorTokens{
		Ramps: RampSet{
			Neutral: dr[0], Primary: dr[1], Secondary: dr[2], Tertiary: dr[3],
			Error: dr[4], Success: dr[5], Warning: dr[6],
		},
		Primary:     darkBase[1],
		OnPrimary:   darkInk[1],
		Secondary:   darkBase[2],
		OnSecondary: darkInk[2],
		Tertiary:    darkBase[3],
		OnTertiary:  darkInk[3],
		Error:       darkBase[4],
		OnError:     darkInk[4],
		Success:     darkBase[5],
		OnSuccess:   darkInk[5],
		Warning:     darkBase[6],
		OnWarning:   darkInk[6],
		Background:  dr[0].Step(100),
		Text:        dr[0].Step(900),
	}, d.dividerStep, lr[0])
	return light, dark
}

// DefaultSeed is the brand seed DefaultLight and DefaultDark derive from:
// #6750A4, the seed every ADR-002/ADR-007 measurement was made against.
var DefaultSeed = stdcolor.NRGBA{R: 0x67, G: 0x50, B: 0xA4, A: 0xff}

// DefaultLight and DefaultDark are the canonical colour token sets:
// FromSeed(DefaultSeed), light and paired dark. The exact derived palette
// is recorded byte-for-byte in this package's golden test.
var DefaultLight, DefaultDark = FromSeed(DefaultSeed)
