package tokens

import (
	"strings"

	"gioui.org/font"
	fontapi "github.com/go-text/typesetting/font"
)

// FaceMetrics are the two measurements of a type role that something drawn
// beside its text has to match: the height of its capitals and the width of
// its upright stem, both in the dp [TextStyle.Size] is in.
//
// They exist because a mark set beside a label is read as part of that label's
// line. The band the words occupy runs from the baseline to the cap height, so
// a mark rises no higher and hangs no lower; and the strokes the words are
// made of are one stem wide, so a stroked mark beside them is stroked at that
// width. Measured on the platform's own marks — the plus, the check and the
// cross set against a system-font label, rendered offscreen — the mark's drawn
// box runs 1.11 to 1.21 times the label's cap height and its stroke band
// equals the label's stem; the excess over the cap band is the optical licence
// a stroke straddling that band takes on its own. See the measured macOS
// reference.
//
// Both come from the face rather than from a ratio pinned in this table,
// because they are the face's property and not the role's: a collection
// carrying a different typeface answers differently for the same role.
type FaceMetrics struct {
	// CapHeight is the distance from the baseline to the top of a flat
	// capital, in dp.
	CapHeight float32

	// Stem is the width of an upright stem, in dp, read as the outline
	// extents of the face's own 'I'.
	Stem float32
}

// Fallback ratios for a role whose typeface this collection cannot answer for:
// the shipped face's own cap height and Medium-weight stem as fractions of the
// type size, so an unresolvable role is still drawn beside its text at the
// right relation rather than at nothing.
const (
	fallbackCapHeightEm = 0.711
	fallbackStemEm      = 0.123
)

// FaceMetrics resolves the face a role names in this collection and reports
// its cap height and stem width at that role's size.
//
// Resolution is by the first family the role names and by weight, nearest
// wins, upright faces only — the same family the shaper picks for the role. A
// role whose family this collection does not carry falls back to the shipped
// face's own ratios, which is a guess and is documented as one: it keeps the
// relation right for every role in this system and merely plausible for a role
// set in a face nobody handed over.
func (t Typography) FaceMetrics(style TextStyle) FaceMetrics {
	fallback := FaceMetrics{
		CapHeight: style.Size * fallbackCapHeightEm,
		Stem:      style.Size * fallbackStemEm,
	}

	family := style.Typeface
	if i := strings.IndexByte(family, ','); i >= 0 {
		family = family[:i]
	}
	family = strings.Trim(strings.TrimSpace(family), `"'`)
	if family == "" {
		return fallback
	}

	want := font.Normal
	if style.Weight != 0 {
		want = FontWeight(style.Weight)
	}
	var chosen font.Face
	best := 1 << 30
	for _, ff := range t.Faces {
		if ff.Face == nil || ff.Font.Style != font.Regular {
			continue
		}
		if !strings.EqualFold(string(ff.Font.Typeface), family) {
			continue
		}
		delta := int(ff.Font.Weight) - int(want)
		if delta < 0 {
			delta = -delta
		}
		if chosen == nil || delta < best {
			chosen, best = ff.Face, delta
		}
	}
	if chosen == nil {
		return fallback
	}

	face := chosen.Face()
	upem := float32(face.Upem())
	if upem == 0 {
		return fallback
	}
	m := fallback
	if cap := face.LineMetric(fontapi.CapHeight); cap > 0 {
		m.CapHeight = style.Size * cap / upem
	}
	if gid, ok := face.NominalGlyph('I'); ok {
		if ext, ok := face.GlyphExtents(gid); ok && ext.Width > 0 {
			m.Stem = style.Size * ext.Width / upem
		}
	}
	return m
}

// FaceMetrics reports this role's cap height and stem width against the
// collection this package ships. A theme carrying faces of its own should ask
// its own [Typography] instead, which is what a component holding one does.
func (s TextStyle) FaceMetrics() FaceMetrics { return DefaultTypography.FaceMetrics(s) }
