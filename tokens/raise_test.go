package tokens_test

import (
	"image/color"
	"math"
	"testing"

	vgcolor "github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/tokens"
)

// schemes is every palette the raise is asserted over: both schemes of both
// derivations, so a rule that holds "in either scheme" is read against the
// high-contrast variant too.
func schemes(seed color.NRGBA) []struct {
	name string
	tok  tokens.ColorTokens
} {
	light, dark := tokens.FromSeed(seed)
	hcLight, hcDark := tokens.FromSeedHighContrast(seed)
	return []struct {
		name string
		tok  tokens.ColorTokens
	}{{"light", light}, {"dark", dark}, {"light-hc", hcLight}, {"dark-hc", hcDark}}
}

// standable is every named level a thing can stand on. The backdrop is not
// among them: nothing stands on the backdrop, so no raise is ever walked
// from it and nothing is ever derived against it.
var standable = []tokens.ElevationLevel{
	tokens.LevelChrome, tokens.Level0, tokens.Level1, tokens.Level2, tokens.Level3,
}

// TestARaiseIsLighterOrCarriesASeam takes the rule over the whole seed
// sweep, in both schemes of both derivations, walking from every level a
// thing can stand on and then from the walk's own answer four times over: a
// raise is lighter than what it stands on, or it carries a seam. Never
// darker, and never both unsaid.
func TestARaiseIsLighterOrCarriesASeam(t *testing.T) {
	worstFill, quietestSeam := 99.0, 99.0
	for _, seed := range sweepSeeds() {
		for _, scheme := range schemes(seed) {
			for _, level := range standable {
				surface := scheme.tok.SurfaceAt(level)
				for depth := 1; depth <= 4; depth++ {
					raise := scheme.tok.RaisedOn(surface)
					if lstar(raise.Fill) < lstar(surface) {
						t.Fatalf("seed %v %s: raise %d off level %d fills L*%.2f under the L*%.2f it stands on — the walk inverted",
							seed, scheme.name, depth, level, lstar(raise.Fill), lstar(surface))
					}
					got := vgcolor.ContrastRatio(raise.Fill, surface)
					if !raise.Seamed {
						if got < tokens.RaiseFloor {
							t.Fatalf("seed %v %s: raise %d off level %d measures %.4f:1, under RaiseFloor, and reports no seam",
								seed, scheme.name, depth, level, got)
						}
						if got < worstFill {
							worstFill = got
						}
					} else {
						// A seam is only discharged if it is findable
						// against BOTH fills, which is the whole of what
						// the Seam entry asks of it.
						below := vgcolor.ContrastRatio(raise.Seam, surface)
						above := vgcolor.ContrastRatio(raise.Seam, raise.Fill)
						if min(below, above) < tokens.SeamRatio {
							t.Fatalf("seed %v %s: raise %d off level %d seams at %.3f:1 below and %.3f:1 above, under SeamRatio %.2f",
								seed, scheme.name, depth, level, below, above, tokens.SeamRatio)
						}
						if q := min(below, above); q < quietestSeam {
							quietestSeam = q
						}
					}
					surface = raise.Fill
				}
			}
		}
	}
	t.Logf("over %d seeds × 4 schemes: faintest fill-told raise %.4f:1, faintest seam %.3f:1",
		len(sweepSeeds()), worstFill, quietestSeam)
}

// TestARaiseOnARaiseOnAModalResolves walks the stack the ruling names — a
// card on a modal, a field in that card — and asserts each step is told, by
// its fill or by its seam, in both schemes. It is the case an absolute level
// cannot answer: there is no level above the two floating ones.
func TestARaiseOnARaiseOnAModalResolves(t *testing.T) {
	for _, seed := range sweepSeeds() {
		for _, scheme := range schemes(seed) {
			modal := scheme.tok.SurfaceAt(tokens.Level2)
			card := scheme.tok.RaisedOn(modal)
			field := scheme.tok.RaisedOn(card.Fill)
			for _, step := range []struct {
				name    string
				beneath color.NRGBA
				raise   tokens.Raise
			}{{"card on the modal", modal, card}, {"field in the card", card.Fill, field}} {
				told := !step.raise.Seamed ||
					vgcolor.ContrastRatio(step.raise.Seam, step.beneath) >= tokens.SeamRatio
				if !told {
					t.Fatalf("seed %v %s: the %s is told by nothing", seed, scheme.name, step.name)
				}
			}
		}
	}
	light, dark := tokens.DefaultLight, tokens.DefaultDark
	t.Logf("default light: modal %v → card %v (seamed %v) → field %v (seamed %v)",
		light.SurfaceAt(tokens.Level2), light.RaisedOn(light.SurfaceAt(tokens.Level2)).Fill,
		light.RaisedOn(light.SurfaceAt(tokens.Level2)).Seamed,
		light.RaisedOn(light.RaisedOn(light.SurfaceAt(tokens.Level2)).Fill).Fill,
		light.RaisedOn(light.RaisedOn(light.SurfaceAt(tokens.Level2)).Fill).Seamed)
	t.Logf("default dark: modal %v → card %v (seamed %v)",
		dark.SurfaceAt(tokens.Level2), dark.RaisedOn(dark.SurfaceAt(tokens.Level2)).Fill,
		dark.RaisedOn(dark.SurfaceAt(tokens.Level2)).Seamed)
}

// TestACardOnTheContentIsToldByItsFill pins what moving the content pin
// bought: in BOTH schemes the first raise off the content clears
// [tokens.RaiseFloor] on its fill alone and owes no seam. The light scheme
// used to answer a 0.7 L* whisper here.
func TestACardOnTheContentIsToldByItsFill(t *testing.T) {
	for _, seed := range sweepSeeds() {
		for _, scheme := range schemes(seed) {
			raise := scheme.tok.RaisedOn(scheme.tok.Background)
			if raise.Seamed {
				t.Fatalf("seed %v %s: a card on the content is seamed at %.4f:1 — the content owes its first raise a whole step",
					seed, scheme.name, vgcolor.ContrastRatio(raise.Fill, scheme.tok.Background))
			}
		}
	}
}

// TestTheContentPinKeepsHeadroomAboveIt pins the pin itself. Where the
// neutral surface band climbs away from its 100 stop the content IS that
// stop; where it descends the content stands one band step under white, and
// white is the first raise on it.
func TestTheContentPinKeepsHeadroomAboveIt(t *testing.T) {
	white := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	for _, seed := range sweepSeeds() {
		light, dark := tokens.FromSeed(seed)
		if got, want := dark.Background, dark.Ramps.Neutral.Step(100); got != want {
			t.Fatalf("seed %v dark: the content is %v, want the band's own 100 stop %v", seed, got, want)
		}
		step := math.Abs(lstar(light.Ramps.Neutral.Step(200)) - lstar(light.Ramps.Neutral.Step(100)))
		if got, want := lstar(light.Background), 100-step; math.Abs(got-want) > 0.05 {
			t.Fatalf("seed %v light: the content stands at L*%.2f, want L*%.2f — one band step under the axis",
				seed, got, want)
		}
		if got := light.RaisedOn(light.Background).Fill; got != white {
			t.Fatalf("seed %v light: the first raise off the content is %v, want white", seed, got)
		}
	}
	t.Logf("default light: content %v, chrome %v, backdrop %v",
		tokens.DefaultLight.Background,
		tokens.DefaultLight.SurfaceAt(tokens.LevelChrome),
		tokens.DefaultLight.SurfaceAt(tokens.LevelBackdrop))
}

// TestTheFloatingLevelsStayAboveWhatIsRaisedBeneathThem pins the one thing
// the fixed table still owes the walk: a dialog and a menu are never darker
// than a card on the content they cover.
func TestTheFloatingLevelsStayAboveWhatIsRaisedBeneathThem(t *testing.T) {
	for _, seed := range sweepSeeds() {
		for _, scheme := range schemes(seed) {
			raised := lstar(scheme.tok.RaisedOn(scheme.tok.Background).Fill)
			for _, level := range []tokens.ElevationLevel{tokens.Level2, tokens.Level3} {
				if got := lstar(scheme.tok.SurfaceAt(level)); got < raised {
					t.Fatalf("seed %v %s: floating level %d fills L*%.2f under the L*%.2f raised on the content",
						seed, scheme.name, level, got, raised)
				}
			}
			if lstar(scheme.tok.SurfaceAt(tokens.Level3)) < lstar(scheme.tok.SurfaceAt(tokens.Level2)) {
				t.Fatalf("seed %v %s: the top floating level is under the one below it", seed, scheme.name)
			}
		}
	}
}

// TestNothingIsWalkedFromTheBackdrop pins the backdrop's exclusion: nothing
// stands on it, so the level above it is placed by the platform measurement
// the file records and not by a raise off it. A chrome level that happened
// to be the backdrop's raise would be deriving against the one surface the
// language says nothing derives against.
func TestNothingIsWalkedFromTheBackdrop(t *testing.T) {
	for _, seed := range sweepSeeds() {
		for _, scheme := range schemes(seed) {
			backdrop := scheme.tok.SurfaceAt(tokens.LevelBackdrop)
			if got, walked := scheme.tok.SurfaceAt(tokens.LevelChrome), scheme.tok.RaisedOn(backdrop).Fill; got == walked {
				t.Fatalf("seed %v %s: chrome %v is the backdrop's raise — chrome is measured, not walked",
					seed, scheme.name, got)
			}
		}
	}
}

// TestSeamIsFilledInWhetherOrNotItIsOwed pins the struct's contract: a
// caller reading Seam never reads a zero colour, so a component that wants
// the hairline for a reason of its own — an outline it draws in every
// scheme, say — can take it without asking whether the raise owed one.
func TestSeamIsFilledInWhetherOrNotItIsOwed(t *testing.T) {
	for _, scheme := range schemes(tokens.DefaultSeed) {
		raise := scheme.tok.RaisedOn(scheme.tok.Background)
		if raise.Seam.A == 0 {
			t.Fatalf("%s: an unowed seam came back unset", scheme.name)
		}
	}
}

// TestSeamOnIsFindableAgainstTheSurfaceItPartsFromItself takes the group's
// hairline over the whole seed sweep, in both schemes of both derivations,
// at every level a thing can stand on: a group takes the fill of the
// surface it is in, so both sides of its line are that one fill, and the
// line owes SeamRatio against it. Nothing else says where a group ends, so
// a hairline that misses the ratio is a group nobody can find.
func TestSeamOnIsFindableAgainstTheSurfaceItPartsFromItself(t *testing.T) {
	quietest := 99.0
	for _, seed := range sweepSeeds() {
		for _, scheme := range schemes(seed) {
			for _, level := range standable {
				surface := scheme.tok.SurfaceAt(level)
				got := vgcolor.ContrastRatio(scheme.tok.SeamOn(surface), surface)
				if got < tokens.SeamRatio {
					t.Fatalf("seed %v %s: the hairline on level %d measures %.3f:1 against the fill it parts, under SeamRatio %.2f:1",
						seed, scheme.name, level, got, tokens.SeamRatio)
				}
				if got < quietest {
					quietest = got
				}
			}
		}
	}
	t.Logf("over %d seeds × 4 schemes × %d levels: faintest group hairline %.3f:1",
		len(sweepSeeds()), len(standable), quietest)
}

// TestSeamOnIsTheDirectionTheSchemeReads holds the direction the derivation
// claims: a hairline goes toward the scheme's own foreground, so a light
// scheme's is darker than the fill it parts and a dark scheme's is lighter.
// A line that went the other way would read as a second fill rather than as
// a boundary.
func TestSeamOnIsTheDirectionTheSchemeReads(t *testing.T) {
	for _, seed := range sweepSeeds() {
		for _, scheme := range schemes(seed) {
			toward := lstar(scheme.tok.Text) > lstar(scheme.tok.Background)
			for _, level := range standable {
				surface := scheme.tok.SurfaceAt(level)
				lighter := lstar(scheme.tok.SeamOn(surface)) > lstar(surface)
				if lighter != toward {
					t.Fatalf("seed %v %s: the hairline on level %d is L*%.2f against a fill of L*%.2f; the foreground is L*%.2f against a background of L*%.2f, so it went the wrong way",
						seed, scheme.name, level, lstar(scheme.tok.SeamOn(surface)), lstar(surface),
						lstar(scheme.tok.Text), lstar(scheme.tok.Background))
				}
			}
		}
	}
}
