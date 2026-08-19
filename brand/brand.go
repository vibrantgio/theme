// Package brand keeps the brand colour a person chose, so a palette worth
// keeping outlives the application that generated it, and hands it back as
// the options a live theme stream is built with.
//
// What is kept is one colour. [tokens.FromSeed] is a pure function of its
// seed and reproduces its own output from the primary it pins, so the seed
// alone regenerates every ramp, pin and on-colour of both schemes exactly —
// there is nothing else to store, and storing the generated colours instead
// would freeze a palette that the generator is still entitled to improve.
// What is kept alongside it is provenance, not input: where the colour came
// from and when it was kept, so a file found six months later can say what
// it is.
//
// # The file
//
// One JSON object:
//
//	{
//	  "seed": "#e8112d",
//	  "base": {
//	    "light": "catppuccin-latte",
//	    "dark": "catppuccin-mocha"
//	  },
//	  "source": "harbour.jpg",
//	  "saved": "2026-08-19T11:04:31Z"
//	}
//
// It sits in an OS-appropriate config directory:
//
//   - darwin:  ~/Library/Application Support/vibrantgio/theme.json
//   - linux:   $XDG_CONFIG_HOME/vibrantgio/theme.json (or ~/.config/...)
//   - windows: %AppData%\vibrantgio\theme.json
//
// That is [os.UserConfigDir], the same root theme/preferences resolves
// against, because a chosen brand is config rather than data. The path is
// per user and NOT per application: the point of keeping a brand is that
// everything the person opens wears it, so the directory is the design
// system's, and one file serves every application that asks.
//
// The seed is spelled the way theme/export's theme.json spells it —
// lowercase #rrggbb under the key "seed" — so the two files agree on how a
// seed is written, and an exported theme.json dropped in as this file loads
// without translation. Keys this package does not know are ignored.
//
// "base" is the second thing a person chooses and the only other thing the
// file holds: the name of the syntax palette code is coloured from, one per
// appearance. It is a name and not a palette for the same reason the seed is
// not a set of ramps — the styling is derived from it, and the derivation is
// entitled to improve. This package neither resolves the name nor judges it:
// it does not know what styles exist, so an empty or unrecognised name is the
// reader's to fall back on, and the fallback is whatever that reader's default
// base is.
//
// It is a pair because a syntax palette is fitted to a ground: a set of inks
// somebody balanced against a near-white page is not the set they would
// balance against a near-black one, and the two appearances of one theme
// therefore call for two names rather than one name and a rule. So the light
// appearance and the dark appearance each name their own, and a person moving
// between them moves between both.
//
// A file whose "base" is a plain string is the spelling that predates the
// pair — one name, with no appearance attached to it. It loads with that name
// in both members, because it is the only name the file has; which appearance
// it was actually fitted to is measured off the style itself, and measuring a
// style is exactly what this package cannot do.
//
// # The styles folder
//
// [StylesDir] names a folder beside the file, where a person can drop style
// files of their own for a highlighter to load — the same shared directory,
// so a style added once is offered by every application that looks. This
// package only says where it is. Nothing here reads it, creates it, or has an
// opinion about what a style file contains.
//
// # Adopting it
//
// [Kept] reads the file and folds every way it can go wrong into the zero
// [Brand], which is "nothing was kept". A zero Brand's [Brand.Options] is
// nil and its [Brand.Colors] are the package defaults, so the adopting call
// is one line that behaves exactly as it did before this package existed
// when there is nothing to adopt:
//
//	th := system.LiveTheme(time.Second, brand.Kept().Options()...)
//
// The kept brand pins the palette pair; which side of it shows is still the
// OS's to decide, and still changes live. Adoption replaces the seed, never
// the light/dark switching.
//
// A missing file and an unreadable one are deliberately the same answer to
// [Kept]: an application asking for a brand it does not have must draw
// something, and what it draws is the default palette either way. Use
// [Load] when the difference matters — it separates "no file" from "a file
// that would not parse" — and note that neither is worth interrupting a
// person over.
package brand

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/vibrantgio/theme/system"
	"github.com/vibrantgio/theme/tokens"
)

// dirName is the directory the file sits in under the user's config
// directory — the design system's own name, because the file is shared by
// every application that adopts the brand rather than owned by one of them.
const dirName = "vibrantgio"

// fileName is the file itself. It is theme/export's name for the same fact
// on purpose: both files answer "which seed is this theme", and a reader
// that finds either knows what it is holding.
const fileName = "theme.json"

// stylesDirName is the folder beside the file, holding style files a person
// added themselves. It is one folder for the same reason the file is one
// file: a style dropped in once is available to everything that looks.
const stylesDirName = "styles"

// Brand is a kept brand colour with the provenance that explains it.
//
// The zero Brand is "nothing kept": Seed's alpha is zero, which no kept
// colour has, so [Brand.Chosen] can tell the two apart without a second
// field and every method degrades to the package defaults.
type Brand struct {
	// Seed is the colour the palette derives from, opaque. It is the
	// whole input: tokens.FromSeed(Seed) is both schemes.
	Seed color.NRGBA

	// Base names the syntax palettes code is coloured from, one per
	// appearance. It is the one other choice the file carries, and it is
	// carried as a name: what resolves it is the highlighting package the
	// reader uses, and what an unknown name means is that reader's default.
	// Empty is the ordinary state — nothing chosen, the reader's default
	// applies — so a file written before this field existed reads as a brand
	// with no base and behaves exactly as it did.
	Base BasePair

	// Source names where the colour was found — a picture's file name, a
	// hand-typed hex, whatever the chooser can honestly say. It is
	// provenance and nothing reads it back as input; empty is allowed and
	// means the chooser had nothing to say.
	Source string

	// Saved is when the colour was kept. [Save] fills it with the current
	// time when it is zero, so a caller that does not care about the clock
	// still writes an honest file.
	Saved time.Time
}

// BasePair is the syntax palette names a kept theme carries: the one code is
// coloured from under a light appearance and the one it is coloured from under
// a dark one.
//
// The two are held apart and never derived from each other. A name is an
// artifact somebody fitted to a ground, and the pair is a person's answer to
// "what should code look like on each of my two grounds" — the reader applies
// the member the appearance on screen calls for, and swaps to the other when
// the appearance changes.
//
// Either member may be empty, which means nothing was chosen for that
// appearance and the reader's own default stands in. So may both, which is a
// brand that says nothing about code at all.
type BasePair struct {
	Light string
	Dark  string
}

// Names returns the pair as the two names a reader resolves, in light, dark
// order. It exists so that resolving a kept pair is one call whatever the file
// looked like: the spelling that predates the pair names one base with no
// appearance attached, and it arrives here in both members, for the reader to
// sort out by measuring the style — which this package cannot do.
func (p BasePair) Names() (light, dark string) { return p.Light, p.Dark }

// Chosen reports whether either member names a palette. It is what a caller
// asks before writing the pair out or comparing it with another.
func (p BasePair) Chosen() bool { return p.Light != "" || p.Dark != "" }

// Chosen reports whether this Brand carries a colour. It is false for the
// zero Brand — the value [Kept] returns when there is no file, or none that
// parses.
func (b Brand) Chosen() bool { return b.Seed.A != 0 }

// Colors returns the pair of schemes the brand generates, or the package
// defaults when nothing was kept. It is what a caller needs before a theme
// stream has emitted anything — the palette to draw the first frame in, so
// that frame is already wearing the kept brand rather than flashing the
// default one at the person who chose against it.
func (b Brand) Colors() (light, dark tokens.ColorTokens) {
	if !b.Chosen() {
		return tokens.DefaultLight, tokens.DefaultDark
	}
	return tokens.FromSeed(b.Seed)
}

// Options returns the theme-stream options that put the kept brand on the
// stream, and nil when nothing was kept — so splatting the result into a
// stream constructor adopts a brand when there is one and changes nothing
// when there is not.
//
// The option pins the palette pair, which means the OS accent colour no
// longer overrides it: a deliberately chosen brand outranks the desktop's.
// Light and dark still follow the OS.
func (b Brand) Options() []system.Option {
	if !b.Chosen() {
		return nil
	}
	return []system.Option{system.WithSeed(b.Seed)}
}

// Path returns the file's path for this user. It creates nothing.
func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("brand: path: %w", err)
	}
	return filepath.Join(dir, dirName, fileName), nil
}

// StylesDir returns the folder beside the file where a person's own style
// files live, for whatever knows how to read one. It creates nothing and does
// not report whether the folder is there: a folder nobody has made yet holds
// no styles, which is the same answer as an empty one.
func StylesDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("brand: styles dir: %w", err)
	}
	return filepath.Join(dir, dirName, stylesDirName), nil
}

// Kept is the forgiving read: the kept brand, or the zero [Brand] when
// there is none to be had for any reason at all — no file, an unreadable
// directory, a file that is not JSON, a seed that is not a colour. Nothing
// an application does with a brand is worth failing to start over, and the
// answer it needs in every one of those cases is the same one.
func Kept() Brand {
	b, _, err := Load()
	if err != nil {
		return Brand{}
	}
	return b
}

// KeptFrom is [Kept] against an explicit path, which is what a test points
// at a temporary directory.
func KeptFrom(path string) Brand {
	b, _, err := LoadFrom(path)
	if err != nil {
		return Brand{}
	}
	return b
}

// Load reads the kept brand from this user's file. The bool reports whether
// a brand was found; the error reports why one could not be read. A missing
// file is (zero, false, nil) — not having chosen a brand is not a failure.
func Load() (Brand, bool, error) {
	path, err := Path()
	if err != nil {
		return Brand{}, false, err
	}
	return LoadFrom(path)
}

// LoadFrom is [Load] against an explicit path.
func LoadFrom(path string) (Brand, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Brand{}, false, nil
	}
	if err != nil {
		return Brand{}, false, fmt.Errorf("brand: load: %w", err)
	}
	var f file
	if err := json.Unmarshal(data, &f); err != nil {
		return Brand{}, false, fmt.Errorf("brand: load %s: %w", path, err)
	}
	seed, err := parseHex(f.Seed)
	if err != nil {
		return Brand{}, false, fmt.Errorf("brand: load %s: %w", path, err)
	}
	b := Brand{Seed: seed, Base: f.Base.pair(), Source: f.Source}
	if f.Saved != "" {
		// An unreadable timestamp costs the provenance, not the brand: the
		// colour is what the file is for, and it parsed.
		if ts, err := time.Parse(time.RFC3339, f.Saved); err == nil {
			b.Saved = ts
		}
	}
	return b, true, nil
}

// Save writes the brand to this user's file, creating the directory. A
// brand with no colour is refused: writing one would leave a file that
// [Kept] reads back as nothing kept, which is a slower way of deleting it.
func Save(b Brand) error {
	path, err := Path()
	if err != nil {
		return err
	}
	return SaveTo(path, b)
}

// SaveTo is [Save] against an explicit path.
func SaveTo(path string, b Brand) error {
	if !b.Chosen() {
		return errors.New("brand: save: the brand carries no colour")
	}
	if b.Saved.IsZero() {
		b.Saved = time.Now()
	}
	data, err := json.MarshalIndent(file{
		Seed:   hexRGB(b.Seed),
		Base:   baseFrom(b.Base),
		Source: b.Source,
		Saved:  b.Saved.UTC().Format(time.RFC3339),
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("brand: save: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("brand: save: %w", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("brand: save: %w", err)
	}
	return nil
}

// file is the JSON shape on disk. It is separate from [Brand] so the file's
// spelling — a hex string, an RFC 3339 timestamp — is a decision this
// package makes once rather than a shape every caller has to hold a colour
// in.
type file struct {
	Seed   string     `json:"seed"`
	Base   *baseField `json:"base,omitempty"`
	Source string     `json:"source,omitempty"`
	Saved  string     `json:"saved,omitempty"`
}

// baseField is how the pair is spelled on disk, and the only place the two
// spellings it has ever had are known: an object with a member per appearance,
// and — from a file written before a theme carried a pair — a plain string
// naming one base with no appearance attached.
//
// The two are told apart by what JSON says they are rather than by a version
// number: a string is a string, an object is an object, and a file cannot be
// both. What is written is always the object, because a theme kept now knows
// both members; the string is read and never emitted.
type baseField struct {
	Light string `json:"light,omitempty"`
	Dark  string `json:"dark,omitempty"`
}

// baseFrom is the pair as it goes to disk, or nothing at all when neither
// member was chosen — an unchosen base leaves no key behind, so the file says
// what was chosen and nothing else.
func baseFrom(p BasePair) *baseField {
	f := baseField{Light: strings.TrimSpace(p.Light), Dark: strings.TrimSpace(p.Dark)}
	if f.Light == "" && f.Dark == "" {
		return nil
	}
	return &f
}

// pair reads the field back as the pair a caller holds. A field that is not
// there at all — a file from before code had a base in it — is the empty pair,
// which is "nothing chosen" and behaves as it always did.
func (f *baseField) pair() BasePair {
	if f == nil {
		return BasePair{}
	}
	return BasePair{Light: f.Light, Dark: f.Dark}
}

// UnmarshalJSON accepts both spellings. A bare string fills both members with
// the one name the file has: it is not a claim about which appearance the
// style was fitted to — nothing here can measure that — it is the whole of
// what was kept, offered to whichever appearance the reader ends up asking
// for.
func (f *baseField) UnmarshalJSON(data []byte) error {
	var one string
	if err := json.Unmarshal(data, &one); err == nil {
		one = strings.TrimSpace(one)
		*f = baseField{Light: one, Dark: one}
		return nil
	}
	// A named type without the method, so unmarshalling it does not call this
	// one again.
	type object baseField
	var o object
	if err := json.Unmarshal(data, &o); err != nil {
		return fmt.Errorf("base is neither a name nor a light/dark pair: %w", err)
	}
	*f = baseField{Light: strings.TrimSpace(o.Light), Dark: strings.TrimSpace(o.Dark)}
	return nil
}

// hexRGB writes a colour as lowercase #rrggbb. A kept seed is opaque, so
// alpha is never written.
func hexRGB(c color.NRGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// parseHex reads #rrggbb, in either case, and returns it opaque. Anything
// else — an empty string, a name, a short form, trailing rubbish — is an
// error, because a file whose seed cannot be read is a file with no brand
// in it and saying so is more useful than guessing a colour.
func parseHex(s string) (color.NRGBA, error) {
	h := strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(h) != 6 {
		return color.NRGBA{}, fmt.Errorf("seed %q is not a #rrggbb colour", s)
	}
	v, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return color.NRGBA{}, fmt.Errorf("seed %q is not a #rrggbb colour", s)
	}
	return color.NRGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xff}, nil
}
