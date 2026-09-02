// Seed-derived palettes. FromSeed turns one brand colour into the complete
// paired light and dark ramp sets, and DefaultLight/DefaultDark are FromSeed
// of the default seed.
//
// Derivation rules and the measurements behind them:
//
//   - The shared lightness scale is CIELAB L* per step, light and paired
//     dark, and it is one curve read twice rather than two tables. A scheme
//     names four depths of its own — its ground (step 100), the deepest of
//     its window surfaces (400), its accent pin (700) and its body ink
//     (900) — and toneCurve places the rungs between them as fractions of
//     each run. The fractions are the light scale's own, so light realizes
//     its measured depths exactly: 97, 92, 85, 74, 63, 51, 39, 28, 6. Read
//     between the dark anchors 8, 30, 82 and 94 the same curve gives 8, 13,
//     19, 30, 46, 64, 82, 86, 94, and the dark surface stack lands on its
//     own measured greys — #222222 and #2e2e2e — without being told to.
//
//     The 900 stop is L* 6 rather than the measured 18 because APCA's soft
//     black clamp caps even pure black near Lc 92 over the L* 92 step-200
//     ground, so Lc ≥ 90 needs L* 6 (min Lc 90.7 across the seven default
//     ramps; L* 18 measures Lc 85–87).
//
//     Only the anchors are read off the measured dark column; its 500 stop
//     is not. A platform's greys are its window surfaces and its ink with
//     nothing between them — read as a nine-step ramp they leave 35 L*
//     between 400 and 500, and a ramp with no mid tones has neither a
//     boundary tone nor a text tone to offer. Derived, the dark scale's
//     worst gap is 18 L* against light's
//     22, its closest neighbours measure 1.11:1 against each other, and its
//     500 and 600 rungs measure 3.4:1 and 6.4:1 over the dark page — one
//     rung in the 3:1 non-text band and one in the 4.5:1 text band, where
//     light's 600 and 700 sit at 4.0:1 and 6.2:1 over the light page. Every
//     one of those is gated in this package's contrast tests.
//
//     Both scales are swept at constant OKLCh hue and chroma via
//     color.Tone, which gamut-maps by chroma reduction.
//
//   - The neutral ramps carry no hue: chroma 0.000, measured 0.0000 at every
//     step in both modes. A surface is not where a brand belongs. Light
//     Background, Surface, Divider and Text measure #f6f6f6, #e8e8e8,
//     #d4d4d4 and #131313; dark #181818, #222222, #2e2e2e and #eeeeee.
//
//   - Accent chroma is a dial, not the seed's own measurement: 0.22 OKLCh
//     chroma. A brand colour between grey (chroma 0.020) and the dial is
//     rendered at the dial, one already past it keeps its own chroma, and a
//     grey one stays grey — a dial must not invent a hue where the brand has
//     none. Secondary and tertiary keep MD3's ratios to primary, 16 : 48 and
//     24 : 48, so they follow the dial: at it they land at 0.073 and 0.110.
//     Tertiary keeps MD3's +60° hue rotation. Conversion anchor for reading
//     MD3's numbers on the OKLCh chroma axis: #6750A4 has HCT chroma 48 and
//     measures OKLCh 0.1305, so one HCT chroma unit ≈ 0.00272 OKLCh chroma,
//     and the dial is HCT ≈ 81.
//
//     Measured on the canonical seed: the light primary base goes 0.1305 →
//     0.2196 (#6750a4 → #723ad4) at the same depth and the same white
//     on-colour, so the pair's contrast is unmoved (Lc 86.9 → 86.5, WCAG
//     6.44 → 6.46). Light steps 500 → 800 go 0.131 → 0.176/0.220/0.220/0.220
//     and dark 200 → 400 go 0.130 → 0.150/0.181/0.220. Light 100–400 and
//     dark 700–900 are gamut-limited at this hue and do not move.
//
//   - The four status roles are hue-anchored, not seed-derived: a semantic
//     colour must not rotate with the brand. Each anchor takes the OKLCh hue
//     and chroma of a canonical Material colour, measured with this module's
//     own converters:
//
//     error    hue  28.7°, chroma 0.178  — #B3261E, L* 39.7
//     success  hue 144.2°, chroma 0.162  — #4CAF50, L* 63.98
//     warning  hue  84.9°, chroma 0.172  — #FFC107, L* 81.52
//     info     hue 248.8°, chroma 0.169  — #2196F3, L* 60.43
//
//     Sources are MD3's canonical error base and Material Green, Amber and
//     Blue 500 — the 500 shade is each family's palette anchor. Only hue and
//     chroma are taken; depths come from the shared lightness scale. The four
//     land 56.2°, 59.3°, 104.6° and 139.9° apart around the OKLCh hue circle,
//     far enough that no status colour reads as its neighbour. Three hold
//     that hue at every depth; warning's is a track, two bullets down.
//
//   - The seed tints the status anchors, and only tints them: each rotates
//     toward the accent hue along the shorter arc by at most statusTint, 3°,
//     and never in chroma. The bound is what makes it a tint: 3° is a
//     twentieth of the smallest gap between two anchors (56.2°), so tinting
//     every anchor to its limit still leaves 50.2° between the closest pair,
//     and no seed can make one status role converge on, reorder past, or
//     leave the family of another. Error stays inside 25.7°–31.7°. A seed at
//     or under greyChroma tints nothing.
//
//     3° rather than 5° because a tint's cost is paid not in hue, where the
//     bound holds it, but in the chroma the gamut happens to hold at the
//     rotated hue, where nothing does. Measured: at 5° a red brand takes the
//     light success mark from #006B13 to #226A00 and a purple brand takes the
//     dark warning pin's chroma from 0.1675 to 0.1458; at 3° the same two
//     measure #136B00 and 0.1536.
//
//     Over the 414-seed sweep the accent is never closer to the fixed error
//     anchor than the error role is, in either scheme (see the
//     accent-versus-error gate in this package's contrast tests).
//
//   - Warning's hue is a function of the tone being realized, because a dark
//     yellow is not read as a dark yellow: at the light scale's step-700
//     depth the flat anchor realizes #785600, an olive-brown, and a brown
//     mark carries no warning. So the hue rotates toward orange as the tone
//     deepens: amber at and above L* 82, then warningBendSlope degrees per L*
//     of further depth, stopped at warningBend.
//
//     The pivot L* 82 is where the anchor itself sits (Amber 500 measures
//     L* 81.52) and is also the dark scale's step-700 depth, so a dark
//     scheme's bright pin is amber by construction. The slope 2.178°/L* is
//     the secant of the amber family's own hue-versus-lightness track:
//     #FFC107 measures h 84.93 at L* 81.52 and #FF6F00 h 46.46 at L* 63.86,
//     and the four shades between hold that slope to within 0.10°/L* (600:
//     2.086, 700: 2.258, 800: 2.275, 900: 2.117).
//
//     The rotation is one-signed: toward orange with depth, never the other
//     way, even though the family's own track keeps rising above the pivot
//     (Amber 300 h 91.2, Amber 50 h 92.9). What the bend fixes is at depth;
//     following the track upward would move every pale warning ground in the
//     system.
//
//     The error family sets the bound, not amber. 30° takes the anchor to
//     54.9°, leaving 26.2° between the two families' deep hues and 20.2° once
//     both tints are spent against each other (a seed between the anchors
//     tints error to 31.7° and warning to 81.9°, whence the bend takes it to
//     51.9°). At the depths a warning is painted at, the bent warning
//     realizes #944600 against the error's #b0250f at L* 39, #6d3100 against
//     #861100 at L* 28, #ed7819 against #f96c54 at L* 63. 35° measures
//     #9a4100 beside #b0250f and the two begin to read as one family; 20°
//     measures #894d00, still a brown.
//
//     The bend is chroma-positive wherever it acts, because sRGB starves
//     amber at depth and holds more of an orange. Asked for the anchor's own
//     chroma, the realized chroma goes 0.0977 → 0.1206 at L* 39, 0.0781 →
//     0.0962 at L* 28, 0.0620 → 0.0763 at L* 19 and 0.0391 → 0.0478 at L* 6.
//     warningChroma therefore stays 0.172: the bend raises the ceiling
//     wherever it acts.
//
//     The bend is a rule for marks and does not reach the washes. A status
//     container reads its hue at the pale tint depth, where the anchor is
//     unrotated, because at the container dial the rotation buys no
//     legibility and costs the status set its separation — it brought a dark
//     scheme's warning wash within 19.27° of its error wash (containers.go).
//
//     Composition with the seed tint is tint first, bend second: the seed
//     tints the anchor and the bend rotates the tinted anchor, so the whole
//     hue track is rigid under the tint, every rung moving by the same ≤ 3°.
//     Bending first and tinting each realized hue afterwards is wrong — the
//     tint takes the shorter arc, so an accent inside the bend's own swing
//     would pull light and deep rungs opposite ways and the ramp would wobble
//     in hue.
//
//     The rungs, the pin and the deep on-ink are all realized through the
//     same hue-at-tone, at the tone each sits at. A status container is not:
//     it takes its step-300 rung's tone and the container dial's chroma, but
//     reads its hue at the pale tint depth, so a wash carries the family's
//     anchor at every depth it is drawn at (see containers.go).
//
//   - Status containers are tonal, not blended: the role's own anchor hue at
//     containerChroma, 0.055, realized at the role ramp's step-300 depth —
//     StatusContainer, with OnStatusContainer for the mark read on it.
//     Alpha-compositing the pinned base over the neutral Surface instead
//     interpolates in non-linear sRGB, which preserves neither hue nor
//     chroma: at 12% the four fills come out at chroma 0.0155–0.0212, near
//     enough to grey to be indistinguishable, and the error fill's hue drags
//     28.7° → 21.6° toward magenta. Realized at a tone, the container keeps
//     its parent's hue exactly and all four carry the same measured chroma,
//     so they differ in hue and nothing else.
//
//     0.055 is the dial the sRGB gamut allows at both container depths for
//     every anchor across its whole tint window: the binding case is amber at
//     the dark step-300 depth, which holds 0.0637 at the worst hue in that
//     window, and the dial keeps 14% of headroom under it so quantization can
//     never clip one container and not another.
//
//     OnStatusContainer takes the most chromatic rung of the role's own ramp
//     that reaches graphicFloor over the container — WCAG 1.4.11's 3:1 for a
//     non-text graphic, which is what a status mark is (MarkOn, in
//     containers.go, is the general form). Asking for the most chromatic rung
//     rather than naming one is what keeps four hues equally saturated: sRGB
//     holds a red only at mid depths and an amber only at high ones, so a
//     fixed rung serves one hue at the cost of the others. Light schemes land
//     on step 700 and dark on 500, except amber, which takes 600 or 700 on
//     the dark scale; the worst mark-on-container pairing over the whole seed
//     sweep measures 4.47:1 and the default seed's eight measure 4.52 and up.
//     Body text on a container is not this pairing — the neutral Text token
//     measures 11.6:1 or better over all eight containers.
//
//   - Pins. The light primary base is the seed at its own hue and CIELAB
//     depth with the accent dial applied to its chroma; only its alpha is
//     forced opaque. Reading it off the ramp instead would lighten it. A
//     brand colour the dial leaves alone comes back byte-for-byte, so a
//     palette seeded from a desktop's accent colour still matches that accent
//     exactly (six of the nine macOS system colours are already past the dial
//     or below the grey threshold). Because the base is what a palette
//     publishes, the derivation reads the accent family's hue and chroma back
//     off the base rather than off the seed, and the dial is a projection:
//     deriving a palette from its own primary base reproduces that palette
//     byte-for-byte, which is what lets a serialized theme name one colour
//     and be rebuilt from it.
//
//     The other light bases are their role's hue and chroma at tone 40, the
//     depth MD3 pins accent bases at and the depth the default seed sits at
//     (L* 40.08). Dark bases are the same hue and chroma re-toned to L* 82 —
//     the dark scale's step-700 depth, beside MD3's dark accent-base tone 80,
//     making the dark pin byte-identical to its ramp's step 700. L* 82 is
//     pinned from both directions: an L* 65 mid-tone is a ground no text
//     reaches Lc 60 over (black tops out near Lc 52, white near 57), and
//     L* 82 is also the shallowest depth the increased-contrast variant's
//     Lc ≥ 75 floor allows (L* 80 reaches only 73.5 against pure black). The
//     solid state walk lands on exact rungs from there (hover 800, pressed
//     900). The accent dial buys the dark primary base nothing: at L* 82 the
//     seed's hue holds chroma 0.0822 in sRGB and no more. It is the one
//     accent surface the dial cannot reach.
//
//   - The inverse pair. Each scheme's InverseSurface and OnInverseSurface are
//     the *other* scheme's Surface and Text — its neutral ramp's steps 200
//     and 900 — so the pair's separation is the counterpart scheme's own
//     body-text separation rather than a second measurement (WCAG 13.75:1
//     light, 15.06:1 dark on the default seed; the high-contrast variant
//     widens both to 15.99:1 and 17.11:1). Both schemes are derived in one
//     pass here, so neither needs anything the other has not computed.
//
//   - On-colours are measured, not assumed. Each pinned base is read in the
//     ink that reaches 4.5:1 over it — WCAG AA for body text — with the
//     scheme's usual ink preferred and the other end of the tonal axis taken
//     when the usual one falls short (see onColour). In the light scheme the
//     pair on offer is White and Black; in the dark scheme the role's own
//     step-100 depth and White.
//
//     The rule is a no-op for almost every base: a light base at tone 40
//     carries White at Lc ≥ 85 (WCAG ≈ 6.4:1) and a dark base at L* 82
//     carries its deep ink at Lc ≥ 73 (WCAG ≈ 11:1). It exists for the one
//     base pinned to no depth — the primary base is the brand colour itself,
//     and a light brand colour under white text measures as little as 2.1:1.
//     Its ink flips and the colour does not move, so the accent stays true to
//     the seed. Across a 414-seed sweep 269 of the light schemes' primary
//     inks flip, nothing else in either scheme does, and no pinned pairing
//     any seed produces measures under the floor. Ramp-step pairings are
//     unaffected: a step is realized at a fixed depth, so 700-and-900 text
//     over 100 and 200 grounds measures 5.4:1 and up whatever the seed.
//
//     The state walk under a solid fill is not part of this and cannot be: a
//     fill walks toward its ramp's 900 end whichever depth its pin sits at
//     (see states.go), so on a mid-depth accent one ink reads at rest and the
//     other reads pressed, and one token cannot be both. Choosing each ink
//     for the whole walk instead of for the resting pair was measured over
//     the same sweep and buys four pairings back under the pointer at the
//     cost of eighty-three at rest.
//
// FromSeedHighContrast derives the increased-contrast variant from the same
// seed by the same machinery — a FromSeed option, not a third hand-written
// scheme. Three widenings, each computed against the APCA gate:
//
//   - The 700 text step deepens to the default scale's 900 depth in both
//     modes — light 700 L* 39 → 6, dark 700 L* 82 → 94 — so 700 text meets
//     the same Lc ≥ 90 bar the default asks only of 900 (light min Lc 90.7,
//     dark 93.0 across the role ramps; APCA's soft black clamp caps lighter
//     choices below 90). The 800 and 900 stops slide outward — light 3 and 0,
//     dark 97 and 100 — keeping the ladder strictly monotonic and the 900
//     gate clear with margin (light Lc 92.3, dark 104.4). Steps 100–600 are
//     the default scale unchanged: the grounds stay, the text pulls away.
//
//   - Divider resolves from Neutral step 500 instead of 300: the separator
//     moves from the subtle-border rung to the strong-border rung.
//
//   - Each pinned base's on-colour is pushed further from its base. The dark
//     pins' on-colours drop from their ramp's step 100 (L* 8, Lc ≈ 74, just
//     under the variant's Lc ≥ 75 floor) to tone 0 (Lc ≥ 76.3). The light
//     pins keep White wherever White is the better of the two ends: it is
//     already the far end of the axis and already clears the floor
//     (Lc ≥ 85.7).
//
//     The on-colour rule follows the variant to a stricter floor: the light
//     ink stands only while it reaches 7:1 rather than 4.5:1. What it can do
//     with the answer is bounded by the axis, which has no ink further out
//     than its two ends — where neither reaches the floor the better of the
//     two stands in both derivations, so the variant's flipped set is the
//     default's plus the sliver where the light ink clears AA and the dark
//     ink still reads higher. Every variant pairing therefore measures at
//     least what the default's does, which is the property its gate holds.
package tokens

import (
	stdcolor "image/color"
	"math"

	"github.com/vibrantgio/theme/color"
)

// toneAnchors are the four CIELAB depths a scheme names for itself, in
// order: its ground (the step-100 window floor), the deepest of its window
// surfaces (step 400), its accent pin (step 700) and its body ink (step
// 900). They are the whole of what separates the two schemes' scales.
type toneAnchors [4]float64

// toneCurve is the shared shape of a nine-step ramp: for each rung, which
// of the three anchored runs it sits in and how far along that run it lies.
// The three runs are the surface stack (100–400, the ground and the
// storeys a tinted surface hovers and presses onto), the mark run (400–700,
// from the deepest surface up to the accent pin, where boundary tones and
// text tones live) and the ink run (700–900, the pin's own hover and
// pressed rungs).
//
// The fractions are the light scale's, measured off the platform's own
// window greys: read between the light anchors they return 97, 92, 85, 74,
// 63, 51, 39, 28, 6 exactly. Stating the shape once is what makes the dark
// scale light's counterpart rather than a second table — a rung does the
// same job in both schemes because it sits at the same place in the same
// run.
var toneCurve = [9]struct {
	run int     // which anchored run the rung sits in
	at  float64 // how far along that run, ground end to ink end
}{
	{0, 0}, {0, 5.0 / 23}, {0, 12.0 / 23}, {0, 1},
	{1, 11.0 / 35}, {1, 23.0 / 35}, {1, 1},
	{2, 11.0 / 33}, {2, 1},
}

// scale realizes toneCurve between one scheme's anchors, to the nearest
// L*. Index i holds step (i+1)*100, matching Ramp.
func (a toneAnchors) scale() [9]int {
	var s [9]int
	for i, rung := range toneCurve {
		from, to := a[rung.run], a[rung.run+1]
		s[i] = int(math.Round(from + (to-from)*rung.at))
	}
	return s
}

// highContrast widens a scale for the increased-contrast variant: steps
// 100–600 stand, so the grounds stay where they are; 700 deepens to the
// scale's own ink depth, so 700 text meets the Lc ≥ 90 bar the default asks
// only of 900; and 800 and 900 slide to the end of the tonal axis the ink
// runs toward, in two equal steps, keeping the ladder strictly monotonic.
// See the file header.
func highContrast(s [9]int) [9]int {
	axisEnd := 100
	if s[8] < s[0] {
		axisEnd = 0
	}
	s[6] = s[8]
	s[8] = axisEnd
	s[7] = (s[6] + s[8]) / 2
	return s
}

// The two schemes' anchors, and the four scales they derive. The light
// ground is the platform's paper and the dark ground its window floor; the
// step-400 anchors are the deepest surface each scheme stacks; the step-700
// anchors are the pin depths, so a pin and its 700 rung are one colour; the
// step-900 anchors are the body inks.
var (
	lightAnchors = toneAnchors{97, 74, statusPinTone, 6}
	darkAnchors  = toneAnchors{8, 30, darkPinTone, 94}

	lightTones = lightAnchors.scale()
	darkTones  = darkAnchors.scale()

	hcLightTones = highContrast(lightTones)
	hcDarkTones  = highContrast(darkTones)
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
	infoHue          = 248.8     // OKLCh hue of Material Blue 500 #2196F3
	infoChroma       = 0.169     // OKLCh chroma of #2196F3
	statusTint       = 3.0       // degrees a status anchor may rotate toward the accent

	// Warning's hue-versus-depth track; see the file header for the whole
	// derivation and the renders the bound was chosen from.
	//
	// warningBendFrom is the L* at and above which warning is amber
	// outright: Amber 500's own depth (L* 81.52), which is also the dark
	// scale's step-700 depth, so the dark pin is amber by construction.
	// warningBendSlope is the secant of the amber family's own hue track
	// from #FFC107 (h 84.93, L* 81.52) to #FF6F00 (h 46.46, L* 63.86), a
	// slope its four intervening shades hold to within 0.17°/L*.
	// warningBend is the most the hue may rotate toward orange, which
	// leaves 20.2° between the deepest warning and the reddest error any
	// seed can ask for.
	warningBendFrom  = 82.0
	warningBendSlope = 2.178
	warningBend      = 30.0

	lightPinTone  = 40 // MD3's accent-base tone; the default seed's own depth
	statusPinTone = 39 // the light scheme's pin depth, and so its step-700 anchor:
	// a status pin IS its ramp's 700 stop rather than landing 3/255 beside it
	darkPinTone = 82 // the dark scheme's pin depth, and so its step-700 anchor;
	// no on-colour reaches Lc 60 over an L* 65 mid-tone, and the
	// increased-contrast variant's Lc ≥ 75 floor allows nothing shallower
	// (L* 80 reaches only 73.5 against pure black)
	darkOnTone   = 8   // dark pins' on-colour depth: the dark scale's step-100 L*
	hcDarkOnTone = 0   // high contrast pushes the dark on-colours to the axis floor
	onFloor      = 4.5 // WCAG AA body text: the ratio an on-colour has to reach
	hcOnFloor    = 7.0 // the increased-contrast variant asks AAA of the same pair

	// The status container dial and the floor its mark is chosen against;
	// see the file header for both derivations and their measurements.
	containerChroma = 0.055 // every status container's measured OKLCh chroma
	containerStep   = 300   // the ramp step a status container is realized at
	graphicFloor    = 3.0   // WCAG 1.4.11 non-text contrast: what a status mark owes its container
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
// darker of the two. Flipping the ink rather than deepening the colour is
// what keeps a palette true to the colour it was seeded with.
//
// The two ends are pure White and pure Black, and that is load-bearing:
// over any colour whatever, the better of white and black reaches 4.58:1,
// so no seed can produce a pinned pairing under the floor. An ink one rung
// short of the axis end — the ramp's own 900 stop — gives that guarantee up
// (it bottoms out at 4.31:1 across a seed sweep).
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

// tintToward rotates a fixed status anchor toward the accent hue along the
// shorter arc of the OKLCh hue circle, by at most statusTint degrees. It is
// the whole of the seed's influence on a status role: the anchor's chroma
// is never touched, so no brand can wash a status colour out or light one
// up, and the rotation is bounded far below the gap between two anchors, so
// no brand can make one status role read as another. A brand at or under
// greyChroma has no hue to lend and lends none.
//
// The bound is what makes the tint safe to state as a property rather than
// to check case by case: the smallest gap between two anchors is 56.2°, and
// tinting both to their limits toward each other closes it by at most 10°.
func tintToward(anchor, accentHue, accentChroma float64) float64 {
	if accentChroma <= greyChroma {
		return anchor
	}
	// The signed shorter arc from the anchor to the accent, in (-180, 180].
	delta := math.Mod(accentHue-anchor+540, 360) - 180
	if delta > statusTint {
		delta = statusTint
	} else if delta < -statusTint {
		delta = -statusTint
	}
	return math.Mod(anchor+delta+360, 360)
}

// hueRule answers a role's OKLCh hue at one realized CIELAB depth. It is the
// whole of the derivation's hue vocabulary: every surface a role has — each
// rung of both ramps, the light pin, the dark pin, the deep on-ink — is
// realized by asking the role's rule for the hue at the tone that surface
// sits at, so a role that varies its hue with depth varies every one of them
// by the one rule and none of its consumers has to know.
type hueRule func(tone int) float64

// flatHue is the rule of every role but warning: one hue, at every depth.
func flatHue(h float64) hueRule {
	return func(int) float64 { return h }
}

// bendingHue is warning's rule: the anchor down to warningBendFrom, then
// warningBendSlope degrees of rotation toward orange per L* of further
// depth, stopped at warningBend. The anchor handed in is the tinted one —
// the seed rotates the anchor and the bend rotates what the seed left, so
// the family's hue track is rigid under the tint (see the file header).
//
// The rotation is one-signed: hues only ever fall from the anchor, never
// rise past it, even though the amber family's own track keeps rising above
// the pivot — what the bend exists to fix is at depth, and a light amber is
// already read as amber.
func bendingHue(anchor float64) hueRule {
	return func(tone int) float64 {
		rotate := warningBendSlope * (warningBendFrom - float64(tone))
		if rotate <= 0 {
			return anchor
		}
		if rotate > warningBend {
			rotate = warningBend
		}
		return math.Mod(anchor-rotate+360, 360)
	}
}

// rampOf sweeps one role's hue rule and chroma across a lightness scale,
// asking the rule for the hue at each step's own depth.
func rampOf(tones [9]int, hue hueRule, chroma float64) Ramp {
	var r Ramp
	for i, tone := range tones {
		r[i] = color.Tone(hue(tone), chroma, tone)
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
// primary base is the same lifted seed FromSeed pins). The full rules are in
// the file header; the derived default-seed variant is recorded in this
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

	// The status anchors are fixed hues the seed may only tint (see
	// tintToward and the file header); the accent family rotates with the
	// brand outright.
	//
	// Warning is the one role whose hue depends on the depth it is realized
	// at; the rest answer one hue at every depth. Both are the same kind of
	// thing here — a hueRule — so nothing downstream of this table has a
	// case for either.
	roles := []struct {
		hue    hueRule
		chroma float64
		// pinTone is the light scheme's pin depth. The accent roles take
		// MD3's tone 40; the status roles take their ramp's own step-700
		// depth, so a status pin and its 700 rung are one colour.
		pinTone int
	}{
		{flatHue(hue), neutralChroma, lightPinTone},
		{flatHue(hue), accent, lightPinTone},
		{flatHue(hue), accent * secondaryShare, lightPinTone},
		{flatHue(hue + tertiaryHueShift), accent * tertiaryShare, lightPinTone},
		{flatHue(tintToward(errorHue, hue, accent)), errorChroma, statusPinTone},
		{flatHue(tintToward(successHue, hue, accent)), successChroma, statusPinTone},
		{bendingHue(tintToward(warningHue, hue, accent)), warningChroma, statusPinTone},
		{flatHue(tintToward(infoHue, hue, accent)), infoChroma, statusPinTone},
	}
	var lr, dr [8]Ramp
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
	var lightBase, darkBase, lightInk, darkInk [8]stdcolor.NRGBA
	for i := 1; i < len(roles); i++ {
		lightBase[i] = color.Tone(roles[i].hue(roles[i].pinTone), roles[i].chroma, roles[i].pinTone)
		darkBase[i] = color.Tone(roles[i].hue(darkPinTone), roles[i].chroma, darkPinTone)
	}
	lightBase[1] = primary // the lifted seed, never read off a ramp step
	for i := 1; i < len(roles); i++ {
		deep := color.Tone(roles[i].hue(d.darkOnTone), roles[i].chroma, d.darkOnTone)
		lightInk[i] = onColour(lightBase[i], White, Black, d.onFloor)
		darkInk[i] = onColour(darkBase[i], deep, White, d.onFloor)
	}

	// Each scheme's inverse pair resolves off the other scheme's neutral
	// ramp — both are derived here, so the pair costs nothing and needs no
	// second rule (see ColorTokens.InverseSurface).
	light = resolveAliases(ColorTokens{
		Ramps: RampSet{
			Neutral: lr[0], Primary: lr[1], Secondary: lr[2], Tertiary: lr[3],
			Error: lr[4], Success: lr[5], Warning: lr[6], Info: lr[7],
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
		Info:        lightBase[7],
		OnInfo:      lightInk[7],
		Background:  lr[0].Step(100),
		Text:        lr[0].Step(900),
	}, d.dividerStep, dr[0])
	dark = resolveAliases(ColorTokens{
		Ramps: RampSet{
			Neutral: dr[0], Primary: dr[1], Secondary: dr[2], Tertiary: dr[3],
			Error: dr[4], Success: dr[5], Warning: dr[6], Info: dr[7],
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
		Info:        darkBase[7],
		OnInfo:      darkInk[7],
		Background:  dr[0].Step(100),
		Text:        dr[0].Step(900),
	}, d.dividerStep, lr[0])
	return light, dark
}

// DefaultSeed is the brand seed DefaultLight and DefaultDark derive from:
// #6750A4, the seed every measurement in this package was made against.
var DefaultSeed = stdcolor.NRGBA{R: 0x67, G: 0x50, B: 0xA4, A: 0xff}

// DefaultLight and DefaultDark are the canonical colour token sets:
// FromSeed(DefaultSeed), light and paired dark. The exact derived palette
// is recorded byte-for-byte in this package's golden test.
var DefaultLight, DefaultDark = FromSeed(DefaultSeed)
