// Package system publishes the operating system's appearance — dark mode
// and the accent colour — as a reactive stream, and bridges it to the theme
// the components above read. A per-OS shim reads the live state behind a
// [Source]; [FromSource] turns a Source plus a poll interval into an
// rx.Observable that emits only when the value changes; [Live] wires the
// current platform's shim, and [LiveTheme] maps that stream to
// [theme.Theme] values whose colour set matches the OS setting. LiveTheme
// also composes the OS accessibility preferences (theme/a11y): reduce
// motion zeroes the emitted motion scale's durations so animated components
// snap. High contrast needs no branch of its own — the platform answers it
// in the values it reports.
//
// Reach for it as the theme argument of a window: LiveTheme(time.Second) is
// what every workbench application hands to theme/window, and from there
// an appearance change reaches every component with no application code.
// Pass your own Source to [FromSource] or [FromSourceTheme] to stub the OS
// out in a test. The package never imports Gio — it speaks to the OS
// directly, so it is usable with or without a window.
//
// Platform support is uneven, and the matrix below is the contract; where
// a cell says "no", the shim reports the zero value for that dimension and
// an application that looks like it is ignoring the setting is not
// misconfigured — the source has nothing to read.
//
//	platform  dark mode                    accent colour
//	macOS     yes — AppleInterfaceStyle    yes — AppleAccentColor index,
//	          via `defaults read -g`       normalized to [Accent] (throttled)
//	Windows   no (always light)            yes — HKCU\Software\Microsoft\
//	                                       Windows\DWM AccentColor, an
//	                                       arbitrary colour → AccentSeed
//	Linux     no (always light)            GNOME 47+: the named accent via
//	                                       `gsettings`, mapped to libadwaita's
//	                                       published colour → AccentSeed
//	                                       KDE Plasma: kdeglobals [General]
//	                                       AccentColor r,g,b → AccentSeed
//	                                       other desktops, older GNOME, or a
//	                                       KDE scheme with no explicit accent:
//	                                       none
//	other     no (always light)            no
//
// With nothing chosen at all — no [WithThemeColor], and a platform
// reporting no colour — the accent a stream emits is the platform's own. On
// macOS that is what AppKit reports for controlAccentColor: the Multicolour
// setting (the absent AppleAccentColor key) and a failed read both mean "no
// accent override", and the live reader answers with the platform's blue.
// On Windows and Linux, whose desktops publish an accent colour but no
// colour set, the recorded set's accent rows are rebuilt for it.
//
// The platform's own colour set — the tokens.PlatformColors every emission
// carries — is read off AppKit on macOS: every name in
// the set under the aqua and darkAqua appearances, resolved to sRGB with
// its alpha through a small Objective-C shim, on the same cadence as the
// accent key. So the accent rows are the platform's own reading rather than
// a derivation, and a settings change arrives within a poll or two. Windows
// and Linux publish no such set: they carry the recorded one with the
// accent rows rebuilt from the colour their desktop reports.
//
// Dark-mode sources for Windows and Linux are a later milestone. The two
// accent shapes are deliberate: macOS's accent is one of eight named
// choices, carried as the [Accent] enum; Windows and Linux accents are
// arbitrary colours, carried raw in Appearance.AccentSeed. Both feed the
// same rebuild of the set's accent rows.
//
// The streams are shared. One [FromSource]/[Live]/[LiveTheme] value
// runs one poll loop no matter how many subscribers attach: the loop starts
// with the first subscriber, later subscribers immediately replay the
// latest value and then track changes, and the loop stops when the
// subscriber count drops to zero (restarting, latest-first, on the next
// subscription). A LiveTheme handed to n layers therefore polls each of its
// two sources — appearance and a11y — once per interval, not n times.
// Distinct calls still get distinct loops: sharing is per observable value,
// so build the stream once and hand the same value around. Keep the
// interval at the intended one second; the OS caches these values and will
// not report a change much sooner.
//
// Errors are invisible by design: a failing Read is folded into the zero
// Appearance rather than an error emission, so a broken source is
// indistinguishable from light mode with no accent. The accent is not just
// carried: with no theme colour chosen, LiveTheme follows it — each
// [Accent] maps to Apple's published colour and the emitted set's accent
// rows are rebuilt for it. An explicit [WithThemeColor] beats the OS
// accent: the application chose its colour, so the accent is ignored
// entirely; with neither, and no colour reported, the platform's own colour
// above stands. [PlatformColor] answers that last question on its own, for
// an application that offers the colour it would rebuild from as a choice.
package system

import (
	"image/color"
	"time"

	"github.com/reactivego/rx"
	"github.com/vibrantgio/theme/a11y"
	"github.com/vibrantgio/theme/internal/poll"
	"github.com/vibrantgio/theme/theme"
	"github.com/vibrantgio/theme/tokens"
)

// Appearance is the OS-level appearance state we observe. All fields are
// comparable so the value can be used with rx.DistinctUntilChanged.
type Appearance struct {
	// Dark is true iff the OS reports a dark interface style.
	Dark bool

	// Accent is the OS accent colour, normalized to this package's
	// [Accent] enum — the shape for platforms whose accent is one of a
	// small named set. On macOS the darwin shim maps the raw
	// AppleAccentColor key (-1 graphite, 0..6 red through pink, absent =
	// multicolour) onto it; platforms without an enum-shaped accent report
	// the zero value. The zero value, AccentDefault, means "no accent
	// override", so the zero Appearance keeps the platform's own accent.
	Accent Accent

	// AccentSeed is the OS accent as a raw colour, for platforms whose
	// accent is an arbitrary colour rather than a named choice: the
	// Windows shim decodes the DWM AccentColor registry value into it, and
	// the Linux shim the GNOME named accent or the KDE kdeglobals RGB.
	// It is meaningful only when AccentSeedSet is true; when set it takes
	// precedence over Accent when the accent rows are rebuilt (an explicit
	// WithThemeColor still beats both).
	AccentSeed color.NRGBA

	// AccentSeedSet reports whether AccentSeed carries a value. A separate
	// flag rather than a sentinel colour keeps every colour — including
	// black — representable, and keeps Appearance comparable for
	// rx.DistinctUntilChanged.
	AccentSeedSet bool
}

// Source reads the current OS appearance state.
// Implement this interface to provide a custom or test-double backend.
type Source interface {
	Read() (Appearance, error)
}

// FromSource returns a shared Observable that polls src every interval,
// emitting Appearance only when the value changes. The first read is
// scheduled immediately (no initial delay).
//
// The returned observable is multicast (FX.5): all subscribers to this one
// value share a single poll loop, a subscriber arriving after the first
// read immediately observes the latest Appearance before tracking changes,
// and the loop stops when the last subscriber unsubscribes (restarting on
// the next subscription). Each FromSource call builds its own loop —
// sharing is per returned value, not per Source.
//
// Read errors are folded into the zero-value Appearance — the stream is
// never an error stream. This keeps the contract simple for consumers
// that only care about the last good value, and matches a11y.FromSource.
func FromSource(src Source, interval time.Duration) rx.Observable[Appearance] {
	return poll.Shared(func() Appearance {
		a, _ := src.Read()
		return a
	}, interval)
}

// Live returns an Observable backed by the current OS's appearance APIs,
// polling every interval and emitting whenever a value changes. Like
// [FromSource] it is shared: n subscribers to one Live value cost one poll
// loop, not n.
//
// Recommended interval: 100–250 ms. One second between an external
// `defaults write` and the corresponding emission is the outer budget, but
// most desktop UIs prefer to feel snappier than that.
func Live(interval time.Duration) rx.Observable[Appearance] {
	return FromSource(defaultSource(), interval)
}

// Option customizes a theme stream. [WithThemeColor] pins the colour the
// emitted set's accent rows are rebuilt for; with no such option the stream
// follows the OS accent instead. Pinning it means the application chose its
// colour, so the OS accent is ignored. It changes only WHICH values are
// emitted; it never affects when emissions happen, so OS dark-mode tracking
// keeps working with a chosen theme colour. [WithTypography] chooses
// the type roles the stream emits; the default is tokens.EmojiTypography().
// [WithA11ySource] chooses where the accessibility preferences composed
// into the emissions are read from.
type Option func(*config)

// config is everything the options configure: the theme colour, the
// typography the stream emits, and the accessibility-preference source
// the stream composes on top of it.
type config struct {
	// themeColor is the colour the accent rows of every emitted set are
	// rebuilt for, and pinned says one was chosen. With nothing chosen the
	// platform's own accent stands.
	themeColor color.NRGBA
	pinned     bool

	// typ is the type roles every emission carries. The default is
	// tokens.EmojiTypography(); [WithTypography] replaces it.
	typ tokens.Typography

	// a11ySrc overrides where accessibility preferences come from. nil
	// means the per-constructor default: the live OS source for
	// [LiveTheme], a constant all-off source for [FromSourceTheme].
	a11ySrc a11y.Source
}

// newConfig applies opts over the defaults. When several theme-colour
// options or several [WithTypography] options are given, the last one of
// each wins.
func newConfig(opts []Option) *config {
	c := &config{typ: tokens.EmojiTypography()}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// a11yStream resolves the accessibility observable for one theme stream:
// the configured source if [WithA11ySource] was given, else fallback.
func (c *config) a11yStream(interval time.Duration, fallback rx.Observable[a11y.A11yPrefs]) rx.Observable[a11y.A11yPrefs] {
	if c.a11ySrc != nil {
		return a11y.FromSource(c.a11ySrc, interval)
	}
	return fallback
}

// WithThemeColor pins the theme colour: the colour the emitted set's accent
// rows are rebuilt for, through [tokens.PlatformColors.WithAccent]. It is
// what the themer keeps and what a brand carries, and it beats the accent
// the OS reports — the application chose its colour.
//
// Nothing else in the set moves. The accent is the one thing the platform
// derives from a colour of the user's choosing; every other name is the
// platform's own answer for the appearance.
func WithThemeColor(c color.NRGBA) Option {
	return func(cfg *config) {
		cfg.themeColor, cfg.pinned = c, true
	}
}

// WithA11ySource overrides where the stream reads accessibility
// preferences. [LiveTheme] defaults to the OS ([a11y.Live]);
// [FromSourceTheme] defaults to a constant all-off source so a test that
// stubs the appearance is hermetic by default — pass a fake [a11y.Source]
// here to exercise the reduce-motion and high-contrast composition.
func WithA11ySource(src a11y.Source) Option {
	return func(c *config) {
		c.a11ySrc = src
	}
}

// WithTypography supplies the type roles the stream emits. The default is
// tokens.EmojiTypography() — DefaultTypography with Noto Color Emoji
// appended as fallback. A brand that names a code face uses this so
// every emission wears it (and still applies WithEmoji). Goldens and
// DeterministicShaper stay on DefaultTypography.
func WithTypography(t tokens.Typography) Option {
	return func(c *config) {
		c.typ = t
	}
}

// LiveTheme bridges system-appearance changes to a theme.Theme stream.
// Each emission is a fresh theme.Theme whose Platform field matches the OS
// dark-mode setting; Typography is [WithTypography]'s value or
// tokens.EmojiTypography(); the remaining token categories use their
// package defaults, modulated by the OS accessibility preferences below.
//
// Which accent the set carries is decided by precedence: an explicit
// [WithThemeColor] wins outright — the application chose its colour, and
// the OS accent is ignored. With none, macOS carries what AppKit reports,
// and Windows and Linux rebuild the recorded set's accent rows for the
// colour their desktop publishes — a raw Appearance.AccentSeed, else the
// colour the [Accent] enum carries. An accent change re-emits the theme.
//
// The stream also composes the OS accessibility preferences
// ([a11y.Live] at the same interval, or [WithA11ySource]'s source): while
// ReduceMotion is on, Motion emits tokens.Motion.Reduced() — every duration
// zero, so duration-driven components snap to their targets. A preference
// toggle re-emits the theme just as an appearance change does.
//
// The two streams it composes are shared: however many layers subscribe to
// one LiveTheme value, the appearance source and the a11y source are each
// polled by exactly one loop.
func LiveTheme(interval time.Duration, opts ...Option) rx.Observable[theme.Theme] {
	c := newConfig(opts)
	prefs := c.a11yStream(interval, a11y.Live(interval))
	return rx.Map(rx.CombineLatest2(Live(interval), prefs), c.theme)
}

// FromSourceTheme is the test-friendly variant of LiveTheme: it lets a
// caller plug in a fake Source while exercising the same Appearance →
// theme.Theme bridge, including any options. Unlike LiveTheme it does NOT
// read the OS accessibility preferences by default — the a11y stream is a
// constant all-off value, so a test's emissions cannot depend on the
// machine it runs on; pass [WithA11ySource] to drive that half too.
func FromSourceTheme(src Source, interval time.Duration, opts ...Option) rx.Observable[theme.Theme] {
	c := newConfig(opts)
	prefs := c.a11yStream(interval, rx.Of(a11y.A11yPrefs{}))
	return rx.Map(rx.CombineLatest2(FromSource(src, interval), prefs), c.theme)
}

// theme maps one (Appearance, A11yPrefs) combination to a theme.Theme
// value: the platform's set for the appearance, with the accent rows
// rebuilt for a pinned theme colour where one was chosen, and ReduceMotion
// swapping the motion scale for its zero-duration variant.
//
// There is no high-contrast branch. The platform answers that preference
// itself: on macOS the live reader asks AppKit for every name while
// "Increase Contrast" is on and gets the values the platform paints under
// it, so the set an emission carries is already the high-contrast one.
func (c *config) theme(v rx.Tuple2[Appearance, a11y.A11yPrefs]) theme.Theme {
	a, prefs := v.First, v.Second
	platform := platformColors(a)
	if c.pinned {
		platform = platform.WithAccent(c.themeColor)
	}
	motion := tokens.Motion
	if prefs.ReduceMotion {
		motion = motion.Reduced()
	}
	return theme.Theme{
		Platform:   rx.Of(platform),
		Typography: rx.Of(c.typ),
		Density:    rx.Of(tokens.Comfortable),
		Motion:     rx.Of(motion),
		Spacing:    rx.Of(tokens.Spacing),
		Radius:     rx.Of(tokens.Radius),
		Elevation:  rx.Of(tokens.Elevation),
	}
}

// PlatformColor is the colour a stream with nothing chosen derives its pair
// from for this appearance, and false where there is nothing to derive from:
// a raw Appearance.AccentSeed (the arbitrary colour a Windows or Linux
// desktop reports), else the seed the [Accent] enum carries, else the colour
// this platform paints an application that has chosen none — systemBlue on
// macOS, nothing on Windows and Linux, where such a stream keeps the
// package's own pair.
//
// It is the fallthrough a stream with no theme colour applies, exported so
// an application can offer that colour as a choice and draw it. Reading
// Appearance.Accent alone is not the same question and answers it wrongly on
// the setting most Macs are on: Multicolour is AccentDefault, which carries
// no seed of its own.
func PlatformColor(a Appearance) (seed color.NRGBA, ok bool) {
	if a.AccentSeedSet {
		return a.AccentSeed, true
	}
	if seed, ok := a.Accent.Seed(); ok {
		return seed, true
	}
	return platformSeed()
}
