// Tonal elevation: elevation is a surface step, and the step goes toward the
// light.
//
// A raised surface separates from the surface beneath it primarily by colour
// and only secondarily by a cast shadow. IN BOTH SCHEMES, A SURFACE NEARER THE VIEWER
// IS LIGHTER. Elevation reads as elevation because a surface nearer the
// viewer catches more light, and reflectance does not invert when the room
// goes dark; what a dark scheme inverts is the ink, which the paired ramps
// already handle. There is no second rule for dark mode and no mirror.
//
// The consequence that orders the levels: chrome furniture is the window's
// BACKDROP, not a level above the paper. A sidebar is the desk the document
// lies on, so it is darker than the document in both schemes. Five levels,
// counted from the backdrop up, away from the desk and toward the reader:
//
//	LevelBackdrop chrome furniture — sidebars, rails, toolbars, inspectors
//	Level0        the paper: the content surface, the Background pin
//	Level1        raised: filled insets on the paper — cards, fences, fields
//	Level2        floating: dialogs and toasts
//	Level3        floating, the topmost level: menus and popovers
//
// Read that down and lightness increases, in the light scheme and in the dark
// one. The dp shadow is the secondary cue and is what effects/depth renders —
// opt-in vibrancy marking what floats and can leave, never a substitute for a
// level.
//
// # Where the levels land, and why they are not ramp steps
//
// A level cannot be a neutral-ramp step, because under one direction the ramp
// can only express half of them: the light ramp descends in lightness from
// its 100 stop, so a light scheme's levels ABOVE the paper are off it
// entirely, and the dark ramp ascends, so a dark scheme's backdrop is off it
// below. Reading a level off a ramp index only works if the numbering is
// allowed to mirror, which it is not.
//
// So the levels are anchored on the Background pin and measured in CIELAB
// L\*, not in ramp indices — but their shape is still read off the ramp
// rather than invented, which is what keeps it generative:
//
//   - THE BACKDROP sits below the pin, toward the scheme's dark extreme, by a
//     step that is MEASURED off the platform rather than derived from the
//     ramp — one measurement per scheme, because the platform itself takes a
//     different step in each. Both numbers, and the captures they came off,
//     are in "The backdrop's two measurements" below.
//
//   - THE LEVELS ABOVE THE PIN take the surface band's own shape — the
//     relative spacing of neutral 100/200/300/400 — scaled into whatever
//     headroom the scheme has above its pin. The ceiling is the lightest tone
//     the band itself reaches; where the band offers none because the pin IS
//     its lightest tone, the ceiling is the tonal axis, L\* 100. In the dark
//     scheme the band's own top is the ceiling, so the three levels land back
//     exactly on neutral 200, 300 and 400. In the light scheme the pin has
//     spent almost all of the axis already, so the same shape compresses into
//     the 3.1 L\* that remain: the levels above the paper are whispers.
//
// # The backdrop's two measurements
//
// The backdrop is the one level the ramp does not get to place. Everything
// above the pin takes the band's own shape, because the band is what the
// scheme has above its paper; below the pin the ramp has nothing to say in
// one scheme and says too much in the other, so the step is read off the
// platform instead. Two schemes, two captures, two numbers — and the code
// tells them apart the way headroom does, off the direction of the band,
// never off a mode flag.
//
// THE LIGHT MEASUREMENT AND THE BAND STEP ARE THE SAME NUMBER. macOS light
// panes sit about 4.9 L\* under their paper, and the ramp's own first surface
// interval |L*(200) − L*(100)| is 4.89, so the backdrop lands byte-for-byte
// on neutral 200 (#E8E8E8 under the #F6F6F6 paper) — the arrangement the stored
// macOS references measure. That half is written as the interval because the
// derivation and the measurement agree, not because the derivation outranks
// the measurement.
//
// THE DARK MEASUREMENT IS 1.48 L\*, AND THE BAND STEP OVERSHOOTS IT BY MORE
// THAN THREE TIMES. A full band step in the dark scheme is 4.98 L\* and
// realizes #0C0C0C under the #181818 paper — 4.93 L\* of separation, which
// reads as pure black where the platform reads dark grey. Three dark
// references, measured 2026-08-28, each the step from a window's content down
// to its furniture:
//
//	Voice Memos, sidebar panel under content   1.50 L*   #1B1B1B under #1E1E1E
//	the reference chat application             1.71 L*
//	macOS Settings, sidebar under content      3.81 L*   #1C2123 under #23292C
//	                                                     (wallpaper tint on)
//
// So the dark step is darkBackdropStep, 1.48 L\*, the Voice Memos reading
// almost
// exactly; on the default dark palette it realizes #151515, 1.47 L\* under
// the paper once eight bits have had their say. The asymmetry between 4.89 and
// 1.48 is the platform's own and not a hand-pick: a light window separates its
// furniture with a step the ramp happens to carry, a dark window with a
// whisper, and this file records both rather than picking one and mirroring
// it.
//
// # The light scheme's headroom, and what the whisper costs
//
// The light scheme takes WHISPER STEPS TOWARD WHITE, WITH THE DERIVED
// HAIRLINES CARRYING THE VISIBLE EDGE.
//
// The obligation that hands the libraries: a light-scheme card, fence, field
// or dialog has almost no fill signal — 0.7, 1.6 and 3.1 L\* above the paper —
// so what says where it is, is its MarkOn-derived border and its corner
// radius. Any construct that takes a level without taking a hairline (an
// inline code chip, say) has to find its separation elsewhere: a tint, a
// border of its own, or the mono face. What it buys: no light-mode surface is
// ever darker than what it lies on, and every light resting arrangement keeps
// its pixels. The reference application measures 0.4 L\* between its page and
// its fence and 0.4 and 1.1 L\* across sidebar, content and composer, with
// hairlines doing the separating.
package tokens

import (
	"fmt"
	"image/color"
	"math"

	vgcolor "github.com/vibrantgio/theme/color"
)

// ElevationScale carries each level's shadow depth in dp — the secondary
// cue, and the only thing a scale holds.
//
// There are five levels, LevelBackdrop through Level3. There are no MD3
// levels 4 and 5: a desktop surface has no six-deep stack to describe, and
// a call site that means "as raised as it gets" names Level3.
//
// The StepN fields describe the DARK scheme's coincidence — where the levels
// above the pin do land on neutral 200, 300 and 400 — and they cannot
// describe the light scheme, whose levels above the pin are off the ramp
// entirely. Ask a palette for a level: [ColorTokens.SurfaceAt] for the fill
// and [ColorTokens.StateAt] for a state walked from it. Prefer the Dp
// accessor over field access in new code.
type ElevationScale struct {
	// Shadow depths in device-independent pixels, following Material
	// Design 3 elevation levels 0–3, with the backdrop casting nothing: the
	// window's desk is behind everything and has nothing to cast onto.
	Backdrop float32 // 0 dp
	Level0   float32 // 0 dp
	Level1   float32 // 1 dp
	Level2   float32 // 3 dp
	Level3   float32 // 6 dp

	// Legacy surface steps on the neutral ramp; see SurfaceStep. Step0 is
	// not a ramp step: its zero value marks the Background pin. The backdrop
	// has no field at all — it is off the ramp in the dark scheme and on
	// it in the light one.
	Step0 int // 0 — sentinel: the Background pin, not a ramp step
	Step1 int // 200
	Step2 int // 300
	Step3 int // 400
}

// Elevation is the default scale instance.
var Elevation = ElevationScale{
	Backdrop: 0,
	Level0:   0,
	Level1:   1,
	Level2:   3,
	Level3:   6,

	Step0: 0, // Background pin
	Step1: 200,
	Step2: 300,
	Step3: 400,
}

// ElevationLevel names one of the five elevation levels. The fill and the dp
// for a given level are read from a palette and from the [Elevation] instance
// respectively.
//
// The levels run from the backdrop up, and the numbering is anchored on the
// paper, which is why the backdrop is −1 rather than a renumbering of
// everything above it: the paper is level 0.
const (
	// LevelBackdrop is chrome furniture: the window's backdrop, one
	// measured step under the paper, in both schemes.
	//
	// A chrome pane that FLOATS — a sidebar a button slides out of the
	// window, an inspector that detaches — is still chrome and still
	// fills here. Its depth is semantic, not geometric: what says it is a
	// floating object is its own hairline edge and its shadow, never a
	// lighter fill. The platform paints even the floating panel darker
	// than the content it sits beside — Voice Memos outlines its panel at
	// #3A3A3A on a #1B1B1B fill, a 1.51:1 whisper of a seam, while the
	// content beside it stays #1E1E1E — so a pane that leaves the wall
	// does not leave the backdrop.
	LevelBackdrop ElevationLevel = iota - 1
	// Level0 is the paper: the content surface, filled with the Background
	// pin.
	Level0
	// Level1 is a raised inset on the paper — a card, a code fence, a
	// text field, a band's own controls.
	Level1
	// Level2 is a floating surface — a dialog, a toast.
	Level2
	// Level3 is the topmost level, nearest the scheme's light
	// extreme — a menu, a popover.
	Level3
)

// ElevationLevel selects an elevation level by name.
type ElevationLevel int

// Raised is the level one step nearer the viewer than level: the fill of
// something that stands on a surface at level rather than on the window's own
// surface. "One step up" means one step toward the scheme's LIGHT extreme, in
// the light scheme and the dark one alike — the walk's direction does not
// depend on which scheme is on, only its destination does.
//
// Levels are walked from the surface a thing is lying on, never from an
// absolute level. So a thing raised over a level-0 plane fills at level 1 and
// the same thing over a level-1 plane fills at level 2, a control on
// furniture fills at the paper's level, and a caller that names the surface
// it stands on gets the separation the grammar asks for wherever it is
// placed. Reaching for a fixed level instead holds only while every surface
// in a window happens to be the one that level was chosen against; move the
// surface one level and the raised thing sits two levels off it, which the
// grammar calls a mistake in either direction.
//
// Level3 is the ceiling. The levels end there, so raising Level3 returns
// Level3 rather than naming a level the scale does not have: a caller
// already at the top asked for "as raised as it gets" and is given it, and
// the clamp is the one place this walk stops instead of stepping. A level
// outside the vocabulary — anything under LevelBackdrop or above Level3 —
// panics, matching [ElevationScale.Dp] and [Ramp.Step].
func (level ElevationLevel) Raised() ElevationLevel {
	switch level {
	case LevelBackdrop, Level0, Level1, Level2:
		return level + 1
	case Level3:
		return Level3
	}
	panic(fmt.Sprintf("tokens: unknown ElevationLevel %d", level))
}

// Dp returns level's shadow depth in device-independent pixels. An
// out-of-vocabulary level panics, matching [Ramp.Step].
func (e ElevationScale) Dp(level ElevationLevel) float32 {
	switch level {
	case LevelBackdrop:
		return e.Backdrop
	case Level0:
		return e.Level0
	case Level1:
		return e.Level1
	case Level2:
		return e.Level2
	case Level3:
		return e.Level3
	}
	panic(fmt.Sprintf("tokens: unknown ElevationLevel %d", level))
}

// SurfaceStep returns the neutral-ramp step a level's state walk starts
// from, or 0 for a level whose fill is not a ramp step at all.
//
// It is only half true. A level is a depth against the Background pin, and
// the three levels above the pin land back on neutral 200, 300 and 400 in
// the DARK scheme alone. In the light scheme they are off the ramp, and so is
// the backdrop in the dark scheme. One scheme-blind integer cannot say that,
// so this one is right where the dark scheme is concerned and stale
// elsewhere.
//
// The backdrop answers 0 — the sentinel this accessor has always used for
// "not a ramp step" — because it is off the ramp in the dark scheme and
// on it in the light one, and callers already handle 0 by falling back.
//
// A caller feeding this into a state walk wants [ColorTokens.StateAt],
// which walks from the level's own fill and needs no ramp index; a caller
// that wanted the fill wants [ColorTokens.SurfaceAt]. An
// out-of-vocabulary level panics, matching [Ramp.Step].
//
// Deprecated: ask a palette, not a scale — the levels left the ramp.
func (e ElevationScale) SurfaceStep(level ElevationLevel) int {
	switch level {
	case LevelBackdrop:
		return 0
	case Level0:
		return e.Step0
	case Level1:
		return e.Step1
	case Level2:
		return e.Step2
	case Level3:
		return e.Step3
	}
	panic(fmt.Sprintf("tokens: unknown ElevationLevel %d", level))
}

// validate panics on a level outside the vocabulary, which is what every
// accessor here does with one and what [Ramp.Step] does with a step number
// the ramp does not carry.
func (level ElevationLevel) validate() {
	if level < LevelBackdrop || level > Level3 {
		panic(fmt.Sprintf("tokens: unknown ElevationLevel %d", level))
	}
}

// SurfaceAt resolves the surface colour of an elevated component: the fill
// of the given level on t.
//
// Level0 is the Background pin exactly — the paper is the pin and the pin is
// off the ramp. Every other level is realized at the pin's own hue and
// chroma, at the CIELAB depth the level sits at: LevelBackdrop the
// platform's own measured step below the pin, Level1 through Level3 the
// surface band's own shape scaled into the headroom the scheme has above it.
// The file header carries both halves — the backdrop's two measurements and
// their captures, and what the light scheme's headroom costs above the paper.
//
// In the dark scheme the three levels above the paper land byte-for-byte
// on Neutral 200, 300 and 400, and in the light scheme the backdrop lands
// byte-for-byte on Neutral 200. A level whose depth coincides with a band
// step answers with that step's own colour rather than a re-realization of
// it, so those identities are exact and not approximate. The dark backdrop
// coincides with nothing: its step is measured, not read off the ramp, and
// on the default palette it realizes #151515.
//
// A state walked from a level is [ColorTokens.StateAt]: the fills above
// the pin are off-ramp in one scheme or the other, so there is no ramp
// index to walk from and the walk is taken on the neutral ramp from the
// level's own colour.
func (t ColorTokens) SurfaceAt(level ElevationLevel) color.NRGBA {
	if level == Level0 {
		return t.Background
	}
	level.validate()
	band, tone := t.surfaceBand()
	pin, _, _ := vgcolor.LabFromNRGBA(t.Background)
	var target float64
	if level == LevelBackdrop {
		// The measured step down, toward the scheme's dark extreme.
		target = pin - backdropStep(pin, tone)
		if target < 0 {
			target = 0
		}
	} else {
		target = pin + t.headroom(pin, tone)*bandShare(level, tone)
	}
	return t.realizeSurface(target, band, tone)
}

// StateAt resolves a level's fill under an interaction state, taken from the
// level's own colour: the wash a control with no fill of its own paints
// on the surface it stands on.
//
// The walk is [ColorTokens.PinnedStateColor]'s — a level above the paper is
// off the ramp in one scheme or the other, so it is exactly the case that
// walk exists for: a fill the scheme's ramps do not carry, walked on the
// neutral ramp because every ramp in a scheme sweeps the same lightness
// scale and only the neutral sweeps it at zero chroma. What it adds is
// [StateFloor]: both colours in this pairing come off that one scale, so a
// step of it can be too small to see, and a wash nobody can see has stopped
// being feedback. Every caller of this — a ghost button's wash, a sidebar
// row's, a tree row's — asks the same question and owes the same minimum,
// so the floor is carried here rather than by each of them.
//
// The walk's direction is deliberately independent of elevation's. A
// state walk heads toward the ramp's 900 end — darker in a light scheme,
// lighter in a dark one — because it says SOMETHING HAPPENED HERE, which
// is feedback rather than depth. A level says how near a surface is; a
// state says something happened.
func (t ColorTokens) StateAt(level ElevationLevel, state State) color.NRGBA {
	return t.washOn(t.SurfaceAt(level), state)
}

// surfaceBand returns the neutral ramp's four surface steps — 100 through
// 400, the tinted-fill half — and their measured CIELAB depths. Everything
// elevation knows about a scheme it reads here: which way the scheme's
// lightness runs, how big one step of it is, and how far the band reaches.
func (t ColorTokens) surfaceBand() (band [4]color.NRGBA, tone [4]float64) {
	for i := range band {
		band[i] = t.Ramps.Neutral.Step((i + 1) * 100)
		tone[i], _, _ = vgcolor.LabFromNRGBA(band[i])
	}
	return band, tone
}

// darkBackdropStep is how far under the paper the backdrop sits where the
// pin is the darkest surface the ramp carries, in CIELAB L\*. It is a
// MEASUREMENT, not a derivation: 1.48 L\*, the step three dark platform
// references measure between a window's furniture and its content, and on
// the default dark palette it realizes #151515 under the #181818 paper.
// The file header quotes the three captures and says why the ramp's own
// band step — which is the measurement in the other scheme — overshoots
// here.
const darkBackdropStep = 1.48

// backdropStep is how far below the pin the backdrop is drawn, in CIELAB
// L\*. Both answers are measurements of the platform, one per scheme, and
// the scheme is never named: which one applies is read off the band, the way
// headroom reads it.
//
// Where the band reaches above the pin, the pin is the darkest surface the
// ramp carries and the platform's measured step is darkBackdropStep. Where
// it does not, the pin is the lightest surface the ramp carries and the
// platform's measured step and the ramp's own first surface interval are
// the same number — which is why that half is written as the interval and
// lands the backdrop byte-for-byte on the band's own 200 step.
func backdropStep(pin float64, tone [4]float64) float64 {
	if bandTop(tone) > pin {
		return darkBackdropStep
	}
	return math.Abs(tone[1] - tone[0])
}

// bandTop is the lightest tone the surface band reaches. Whether it lies
// above the pin is the one fact that tells the two schemes apart in this
// file — a scheme whose band climbs away from its 100 stop has ramp
// headroom above its paper and a pin that is its darkest surface; one
// whose band descends has neither — and it is read off the band rather
// than carried as a mode flag.
func bandTop(tone [4]float64) float64 {
	top := tone[0]
	for _, l := range tone {
		if l > top {
			top = l
		}
	}
	return top
}

// headroom is how much CIELAB lightness the scheme has above its pin for
// the levels to spend: the band's own lightest tone where the band reaches
// past the pin, and the tonal axis where it does not.
//
// The two cases are the two schemes, and neither is named. A dark scheme's
// band climbs away from its 100 stop, so the band's top IS the ceiling and
// the levels land back on the band's own steps. A light scheme's band
// descends from its 100 stop, so the pin is already the lightest surface
// the ramp carries and the only room left is what remains to white.
func (t ColorTokens) headroom(pin float64, tone [4]float64) float64 {
	ceiling := bandTop(tone)
	if ceiling <= pin {
		ceiling = 100 // the band offers no room above the pin; the axis does
	}
	if ceiling <= pin {
		return 0
	}
	return ceiling - pin
}

// bandShare is a level's position as a share of the headroom: the surface
// band's own spacing, normalized against its full reach, so the levels above
// the pin keep the ramp's proportions whatever room the scheme has for them.
func bandShare(level ElevationLevel, tone [4]float64) float64 {
	span := math.Abs(tone[3] - tone[0])
	if span == 0 {
		return 0 // a band with no reach places every level on the pin
	}
	switch level {
	case Level1:
		return math.Abs(tone[1]-tone[0]) / span
	case Level2:
		return math.Abs(tone[2]-tone[0]) / span
	}
	return 1 // Level3 takes the whole headroom
}

// realizeSurface renders a level's depth as a colour at the Background
// pin's own hue and chroma — every level is the pin, lightened or darkened,
// so it carries whatever tint the pin carries and none of its own.
//
// A depth that coincides with one of the band's own steps answers with
// that step verbatim rather than re-realizing it. The realization would
// agree to the byte for a neutral palette, since both ask the same solver
// for the same depth at the same chroma; answering with the step makes
// that an identity instead of a coincidence, and saves the solver a run on
// the levels a window paints most.
func (t ColorTokens) realizeSurface(target float64, band [4]color.NRGBA, tone [4]float64) color.NRGBA {
	for i, l := range tone {
		if math.Abs(target-l) < 1e-6 {
			return band[i]
		}
	}
	_, chroma, hue := vgcolor.OKLChFromNRGBA(t.Background)
	return vgcolor.NRGBAFromToneChromaHue(target, chroma, hue)
}
