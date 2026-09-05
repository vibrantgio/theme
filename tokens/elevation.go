// Tonal elevation: elevation is a surface step, and the step goes toward the
// light.
//
// A surface separates from the one beneath it primarily by colour and only
// secondarily by a cast shadow. IN BOTH SCHEMES A SURFACE NEARER THE VIEWER
// IS NEVER DARKER THAN THE ONE BENEATH IT. A surface nearer the viewer
// catches more light, and reflectance does not invert when the room goes
// dark; what a dark scheme inverts is the foreground, which the paired ramps
// already handle. There is no second rule for the dark scheme and no mirror.
//
// Six levels, counted from the backdrop up:
//
//	LevelBackdrop the bare window plane. Nothing is drawn at it; it shows
//	              wherever nothing stands — around an inset pane, between
//	              regions — and it is the darkest region in either scheme.
//	LevelChrome   the window's furniture: navbar, toolbar, sidebar,
//	              inspector, status bar, pane.
//	Level0        the content itself, filled with the Background pin.
//	Level1        raised on the content — cards, filled insets, fields.
//	Level2        floating — dialogs and toasts.
//	Level3        floating, the top of the elevation — menus, popovers and
//	              tooltips.
//
// Chrome is window-scale only: the trim inside a component or a pattern — a
// card's header, a dialog's footer, a table's header row — is that thing's
// structure and takes no level of its own. The dp shadow is the secondary
// cue and is what effects/depth renders — opt-in vibrancy marking what
// floats and can leave, never a substitute for a level.
//
// # A raise is walked, not read off the table
//
// Only four of the six are placed by this file: the backdrop, the chrome
// level, the content, and the two floating levels, which are absolute
// because a dialog and a menu are detached from whatever they cover. Level1
// is not a table entry at all — it is [ColorTokens.RaisedOn] taken from the
// content, and any other raise is that same walk taken from whatever the
// thing actually stands on (raise.go). A card on a modal is one step above
// the modal, a field in that card one step above the card, and neither
// names a level.
//
// One step is the surface band's own first interval, |L*(200) − L*(100)|,
// in both schemes: 4.88 L* in the light band and 4.98 L* in the dark one.
// The floating levels are held no darker than the first raise off the
// content, which is what "above everything raised beneath them" means when
// the scheme has run out of room.
//
// # Where the levels land, and why they are not ramp steps
//
// A level cannot be a neutral-ramp step, because under one direction the
// ramp can only express half of them: the light ramp descends in lightness
// from its 100 stop, so a light scheme's levels ABOVE the content are off
// it entirely, and the dark ramp ascends, so a dark scheme's levels BELOW
// the content are off it. Reading a level off a ramp index only works if
// the numbering is allowed to mirror, which it is not.
//
// So the levels are anchored on the Background pin and measured in CIELAB
// L\*, not in ramp indices — but their shape is still read off the ramp
// rather than invented, which is what keeps it generative. On both sides of
// the pin that shape is the surface band's own: the relative spacing of
// neutral 100/200/300/400. The one fact that tells the two schemes apart is
// whether that band CLIMBS away from its own 100 stop; nothing here carries
// a mode flag.
//
//   - ABOVE THE PIN the two floating levels take the band's spacing scaled
//     into whatever headroom the scheme has. The ceiling is the band's own
//     top where the band climbs, and the tonal axis, L\* 100, where it does
//     not. In the dark scheme the band's top is the ceiling, so the two
//     floating levels land back exactly on neutral 300 and 400.
//
//   - BELOW THE PIN chrome takes one step down, and the backdrop takes the
//     band's second interval expressed as a multiple of that step.
//
// # The content pin keeps headroom above it
//
// Where the band climbs away from its 100 stop the pin IS that stop: the
// dark content is the band's darkest surface and has the whole band above
// it. Where the band descends there is nothing above the 100 stop at all,
// so pinning the content there would leave a raise nowhere to go — which is
// what the light scheme used to do, spending 3.1 L\* of axis on three
// levels and answering a card a 0.7 L\* whisper nobody could see.
//
// So the content pin STANDS ONE BAND STEP BELOW THE CEILING where the band
// offers nothing above it: L\* 100 − |L*(200) − L*(100)|, realized at the
// band's own hue and chroma, which on the default light palette is
// #F1F1F1 — and white is then the first raise on it, at 1.13:1. That is the
// arrangement the platform ships: macOS light stands grouped content on an
// off-white plane and fills the cells raised on it white.
//
// Chrome and the backdrop keep their measured relation to the pin, so they
// move with it: on the default light palette #E3E3E3 and #CFCFCF, one and
// two-and-a-half band steps under the content rather than landing on the
// band's own 200 and 300 stops as they did while the pin was the 100 stop.
//
// # The chrome step's two measurements
//
// Chrome is the one level the ramp does not get to place in both schemes.
// Two schemes, two captures, two numbers — and the code tells them apart
// off the direction of the band, never off a mode flag.
//
// THE LIGHT MEASUREMENT AND THE BAND STEP ARE THE SAME NUMBER. macOS light
// panes sit about 4.9 L\* under their content, and the ramp's own first
// surface interval |L*(200) − L*(100)| is 4.88, so the light chrome step is
// written as that interval — because the derivation and the measurement
// agree, not because the derivation outranks the measurement.
//
// THE DARK MEASUREMENT IS 1.48 L\*, AND THE BAND STEP OVERSHOOTS IT BY MORE
// THAN THREE TIMES. A full band step in the dark scheme is 4.98 L\* and
// realizes #0C0C0C under the #181818 content — 4.93 L\* of separation,
// which reads as pure black where the platform reads dark grey. Three dark
// references, measured 2026-08-28, each the step from a window's content
// down to its furniture:
//
//	Voice Memos, sidebar panel under content   1.50 L*   #1B1B1B under #1E1E1E
//	the reference chat application             1.71 L*
//	macOS Settings, sidebar under content      3.81 L*   #1C2123 under #23292C
//	                                                     (wallpaper tint on)
//
// So the dark chrome step is darkChromeStep, 1.48 L\*, the Voice Memos
// reading almost exactly; on the default dark palette it realizes #151515
// under the #181818 content. The asymmetry between 4.88 and 1.48 is the
// platform's own and not a hand-pick: a light window separates its
// furniture with a step the ramp happens to carry, a dark window with a
// whisper, and this file records both rather than picking one and
// mirroring it.
//
// # The backdrop has no platform capture
//
// A macOS window paints its furniture edge to edge, so none of the stored
// references shows a window plane beneath it: the backdrop's step has
// nothing to be measured against and is DERIVED, which the comment says
// rather than dressing it as a measurement. It is the chrome step scaled by
// the surface band's own proportion, |L*(300) − L*(100)| over
// |L*(200) − L*(100)| — the same shape the levels above the pin take.
//
// # What the scheme has room for, and what the seam carries
//
// The dark scheme has four band steps above its content and the light
// scheme has one. So a card is told by its fill in both — that is what
// moving the pin bought — and the raise ABOVE that card is where the light
// scheme runs out: a field on a card, a card on a modal. There the raise is
// told by a seam at its own edge instead, which raise.go derives and
// [Raise.Seamed] reports. A raise never vanishes.
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
// There are six levels, LevelBackdrop through Level3. There are no MD3
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
	// Design 3 elevation levels 0–3. Neither level under the content casts
	// anything: the backdrop is the plane everything else stands on, and
	// chrome lies flat on it.
	Backdrop float32 // 0 dp
	Chrome   float32 // 0 dp
	Level0   float32 // 0 dp
	Level1   float32 // 1 dp
	Level2   float32 // 3 dp
	Level3   float32 // 6 dp

	// Legacy surface steps on the neutral ramp; see SurfaceStep. Step0 is
	// not a ramp step: its zero value marks the Background pin. The two
	// levels under the content have no field at all — they are off the ramp
	// in the dark scheme and on it in the light one.
	Step0 int // 0 — sentinel: the Background pin, not a ramp step
	Step1 int // 200
	Step2 int // 300
	Step3 int // 400
}

// Elevation is the default scale instance.
var Elevation = ElevationScale{
	Backdrop: 0,
	Chrome:   0,
	Level0:   0,
	Level1:   1,
	Level2:   3,
	Level3:   6,

	Step0: 0, // Background pin
	Step1: 200,
	Step2: 300,
	Step3: 400,
}

// ElevationLevel names one of the six elevation levels. The fill and the dp
// for a given level are read from a palette and from the [Elevation] instance
// respectively.
//
// The levels run from the backdrop up, and the numbering is anchored on the
// content, which is why the two levels under it are −2 and −1 rather than a
// renumbering of everything above them: the content is level 0.
const (
	// LevelBackdrop is the bare window plane, one derived step under the
	// chrome level in both schemes. Nothing is drawn at it; it is what
	// shows wherever nothing stands.
	LevelBackdrop ElevationLevel = iota - 2
	// LevelChrome is the window's furniture — navbar, toolbar, sidebar,
	// inspector, status bar, pane — one measured step under the content in
	// both schemes. It is window-scale only: the trim inside a component
	// or a pattern takes no level of its own.
	//
	// A chrome region that FLOATS — a sidebar a button slides out of the
	// window, an inspector that detaches — is still chrome and still
	// fills here. Its depth is semantic, not geometric: what says it is a
	// floating object is its own hairline edge and its shadow, never a
	// lighter fill. The platform paints even the floating panel darker
	// than the content it sits beside — Voice Memos outlines its panel at
	// #3A3A3A on a #1B1B1B fill, a 1.51:1 whisper of a seam, while the
	// content beside it stays #1E1E1E — so a pane that leaves the wall
	// does not leave the chrome level.
	LevelChrome
	// Level0 is the content surface, filled with the Background pin.
	Level0
	// Level1 is a raised inset on the content — a card, a code fence, a
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

// Dp returns level's shadow depth in device-independent pixels. An
// out-of-vocabulary level panics, matching [Ramp.Step].
func (e ElevationScale) Dp(level ElevationLevel) float32 {
	switch level {
	case LevelBackdrop:
		return e.Backdrop
	case LevelChrome:
		return e.Chrome
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
// the DARK scheme alone. In the light scheme they are off the ramp, and so
// are the backdrop and the chrome level in the dark scheme. One scheme-blind
// integer cannot say that, so this one is right where the dark scheme is
// concerned and stale elsewhere.
//
// The backdrop and the chrome level answer 0 — the sentinel this accessor
// has always used for "not a ramp step" — because they are off the ramp in
// the dark scheme and on it in the light one, and callers already handle 0
// by falling back.
//
// A caller feeding this into a state walk wants [ColorTokens.StateAt],
// which walks from the level's own fill and needs no ramp index; a caller
// that wanted the fill wants [ColorTokens.SurfaceAt]. An
// out-of-vocabulary level panics, matching [Ramp.Step].
//
// Deprecated: ask a palette, not a scale — the levels left the ramp.
func (e ElevationScale) SurfaceStep(level ElevationLevel) int {
	switch level {
	case LevelBackdrop, LevelChrome:
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

// SurfaceAt resolves the surface colour of a named level on t.
//
// Level0 is the Background pin exactly — the content is the pin. Level1 is
// not a table entry: it is [ColorTokens.RaisedOn] taken from the content,
// which is what "raised on the content" means, and a thing raised on
// anything else asks that walk directly rather than naming a level. Every
// other level is realized at the pin's own hue and chroma, at the CIELAB
// depth the level sits at: LevelChrome the platform's own measured step
// below the pin, LevelBackdrop that step scaled by the surface band's own
// proportion, and the two floating levels the band's shape scaled into the
// headroom the scheme has above the content.
//
// The floating levels are held NO DARKER THAN THE FIRST RAISE off the
// content. That is the whole of "floating: above everything raised beneath
// them" — a dialog that filled under the card it covers would be reporting
// depth backwards — and it binds only where the scheme has run out of room,
// which is the light scheme, where the first raise is already white.
//
// In the dark scheme the two floating levels land byte-for-byte on Neutral
// 300 and 400 and the first raise on Neutral 200. In the light scheme
// nothing above the content lands on the band at all: the pin stands one
// band step under the axis and every level above it is white.
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
	if level == Level1 {
		return t.RaisedOn(t.Background).Fill
	}
	band, tone := t.surfaceBand()
	pin, _, _ := vgcolor.LabFromNRGBA(t.Background)
	var target float64
	switch level {
	case LevelBackdrop:
		target = pin - backdropStep(tone)
	case LevelChrome:
		target = pin - chromeStep(tone)
	default:
		target = pin + t.headroom(tone)*bandShare(level, tone)
		if raised, _, _ := vgcolor.LabFromNRGBA(t.RaisedOn(t.Background).Fill); target < raised {
			target = raised
		}
	}
	// The two levels under the content walk toward the scheme's dark
	// extreme, which the tonal axis ends at.
	if target < 0 {
		target = 0
	}
	return t.realizeSurface(target, band, tone)
}

// StateAt resolves a level's fill under an interaction state, taken from the
// level's own colour: the fill a control with no fill of its own paints
// on the surface it stands on.
//
// The walk is [ColorTokens.PinnedStateColor]'s — a level above the content is
// off the ramp in one scheme or the other, so it is exactly the case that
// walk exists for: a fill the scheme's ramps do not carry, walked on the
// neutral ramp because every ramp in a scheme sweeps the same lightness
// scale and only the neutral sweeps it at zero chroma. What it adds is
// [StateFloor]: both colours in this pairing come off that one scale, so a
// step of it can be too small to see, and a fill nobody can see has stopped
// being feedback. Every caller of this — a ghost button's fill, a sidebar
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

// darkChromeStep is how far under the content the chrome level sits where
// the pin is the darkest surface the ramp carries, in CIELAB L\*. It is a
// MEASUREMENT, not a derivation: 1.48 L\*, the step three dark platform
// references measure between a window's content and its furniture, and on
// the default dark palette it realizes #151515 under the #181818 content.
// The file header quotes the three captures and says why the ramp's own
// band step — which is the measurement in the other scheme — overshoots
// here.
const darkChromeStep = 1.48

// chromeStep is how far below the pin the chrome level is drawn, in CIELAB
// L\*. Both answers are measurements of the platform, one per scheme, and
// the scheme is never named: which one applies is read off the direction of
// the band, the way every other rule in this file reads it.
//
// Where the band climbs away from its 100 stop, the content is the darkest
// surface the ramp carries and the platform's measured step is
// darkChromeStep. Where it descends, the platform's measured step and the
// ramp's own first surface interval are the same number — which is why that
// half is written as the interval.
func chromeStep(tone [4]float64) float64 {
	if bandClimbs(tone) {
		return darkChromeStep
	}
	return raiseStep(tone)
}

// backdropStep is how far below the pin the backdrop is drawn, in CIELAB
// L\*. No stored platform capture shows a window plane beneath its
// furniture, so this one is DERIVED rather than measured: the chrome step
// scaled by the surface band's own proportion, its second interval over its
// first, which is the same shape the levels above the pin take.
//
// A band whose first interval is zero has no proportion to give, and the
// backdrop then takes a second chrome step — the least that keeps every
// level lighter than the one beneath it.
func backdropStep(tone [4]float64) float64 {
	chrome := chromeStep(tone)
	first := raiseStep(tone)
	if first == 0 {
		return 2 * chrome
	}
	return chrome * math.Abs(tone[2]-tone[0]) / first
}

// bandTop is the lightest tone the surface band reaches.
func bandTop(tone [4]float64) float64 {
	top := tone[0]
	for _, l := range tone {
		if l > top {
			top = l
		}
	}
	return top
}

// bandClimbs reports whether the surface band reaches above its own 100
// stop. It is the ONE fact that tells the two schemes apart in this
// package — a band that climbs has room above its 100 stop and pins the
// content there; a band that descends has none, and the content pin steps
// down off the axis to make some — and it is read off the band rather than
// carried as a mode flag. It asks the band about the band, never about the
// pin, so it keeps answering after the pin has moved off the 100 stop.
func bandClimbs(tone [4]float64) bool {
	return bandTop(tone) > tone[0]
}

// raiseStep is one step of the elevation, in CIELAB L\*: the surface band's
// own first interval, |L*(200) − L*(100)|, in both schemes. Every raise is
// this far above the surface it stands on, and the chrome step below the
// pin is the same interval in the scheme whose platform measurement agrees
// with it (see chromeStep).
func raiseStep(tone [4]float64) float64 {
	return math.Abs(tone[1] - tone[0])
}

// ceiling is the lightest tone the scheme's surfaces may reach: the band's
// own top where the band climbs away from its 100 stop, and the tonal axis
// where it does not, the axis being the only room a descending band leaves.
func ceiling(tone [4]float64) float64 {
	if bandClimbs(tone) {
		return bandTop(tone)
	}
	return axisTop
}

// axisTop is the light extreme of the tonal axis, L\* 100 — white.
const axisTop = 100.0

// contentPin is the fill of the content plane, the Background pin, read off
// the neutral surface band alone.
//
// Where the band climbs away from its 100 stop the pin IS that stop: the
// content is the darkest surface the ramp carries and the whole band stands
// above it. Where the band descends the 100 stop is already the lightest
// surface the ramp carries, so pinning the content there would leave a
// raise nowhere to go; the pin stands one band step under the ceiling
// instead, realized at the band's own hue and chroma, and white is the
// first raise on it. The file header carries what that buys and what the
// platform does.
func contentPin(neutral Ramp) color.NRGBA {
	var tone [4]float64
	for i := range tone {
		tone[i], _, _ = vgcolor.LabFromNRGBA(neutral.Step((i + 1) * 100))
	}
	if bandClimbs(tone) {
		return neutral.Step(100)
	}
	_, chroma, hue := vgcolor.OKLChFromNRGBA(neutral.Step(100))
	return vgcolor.NRGBAFromToneChromaHue(ceiling(tone)-raiseStep(tone), chroma, hue)
}

// headroom is how much CIELAB lightness the scheme has above its pin for
// the floating levels to spend: the ceiling less the pin.
func (t ColorTokens) headroom(tone [4]float64) float64 {
	pin, _, _ := vgcolor.LabFromNRGBA(t.Background)
	if room := ceiling(tone) - pin; room > 0 {
		return room
	}
	return 0
}

// bandShare is a floating level's position as a share of the headroom: the
// surface band's own spacing, normalized against its full reach, so the two
// levels above the raised ones keep the ramp's proportions whatever room
// the scheme has for them. Level1 is not among them — it is walked, not
// placed (see [ColorTokens.SurfaceAt]).
func bandShare(level ElevationLevel, tone [4]float64) float64 {
	span := math.Abs(tone[3] - tone[0])
	if span == 0 {
		return 0 // a band with no reach places every level on the pin
	}
	if level == Level2 {
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
