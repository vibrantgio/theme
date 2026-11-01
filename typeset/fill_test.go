package typeset_test

import (
	"image/color"
	"testing"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"

	vgcolor "github.com/vibrantgio/theme/color"
	"github.com/vibrantgio/theme/tokens"
	"github.com/vibrantgio/theme/typeset"
)

// TestWithFillSetsTheFillOnACopy: WithFill sets the fill on the context it
// returns. The caller's context keeps the fill it had.
func TestWithFillSetsTheFillOnACopy(t *testing.T) {
	gtx := layout.Context{Ops: new(op.Ops)}
	on := typeset.WithFill(gtx, tokens.PlatformLight.WindowBackground)
	if got, ok := typeset.Fill(on); !ok || got != tokens.PlatformLight.WindowBackground {
		t.Errorf("the fill is %v %v, want the window background %v", got, ok, tokens.PlatformLight.WindowBackground)
	}
	if _, ok := typeset.Fill(gtx); ok {
		t.Error("the caller's context took the fill")
	}
	inner := typeset.WithFill(on, tokens.PlatformLight.DefaultButtonFill)
	if got, _ := typeset.Fill(inner); got != tokens.PlatformLight.DefaultButtonFill {
		t.Errorf("an opaque inner fill is %v, want %v", got, tokens.PlatformLight.DefaultButtonFill)
	}
	if got, _ := typeset.Fill(on); got != tokens.PlatformLight.WindowBackground {
		t.Errorf("the outer context's fill changed to %v", got)
	}
}

// TestATranslucentFillIsFlattened: a translucent fill is flattened onto the
// fill under it. With no fill under it, it is flattened onto the grey half
// way between black and white.
func TestATranslucentFillIsFlattened(t *testing.T) {
	veil := color.NRGBA{A: 0x40}
	grey := color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}
	under := tokens.PlatformLight.WindowBackground
	gtx := layout.Context{Ops: new(op.Ops)}
	if got, _ := typeset.Fill(typeset.WithFill(typeset.WithFill(gtx, under), veil)); got != vgcolor.Flatten(veil, under) {
		t.Errorf("on the window background: %v, want %v", got, vgcolor.Flatten(veil, under))
	}
	if got, _ := typeset.Fill(typeset.WithFill(gtx, veil)); got != vgcolor.Flatten(veil, grey) {
		t.Errorf("on no fill: %v, want %v", got, vgcolor.Flatten(veil, grey))
	}
}

// TestFontOnPicksTheStemWeightBySizeAndFill: through the fill on gtx, the
// font of text in the label colour, in the secondary label colour and of
// the white title on the default button is the font of its role's size and
// of its colour against that fill, in both schemes. At 13 pt the label
// colour on the light window background draws in a widened face.
func TestFontOnPicksTheStemWeightBySizeAndFill(t *testing.T) {
	white := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	typ := tokens.DefaultTypography
	for _, role := range []tokens.TextStyle{typ.Caption1, typ.SystemFont, typ.LargeTitle} {
		for _, tc := range []struct {
			name string
			fg   color.NRGBA
			fill color.NRGBA
		}{
			{"light: label colour on the window background", tokens.PlatformLight.Label, tokens.PlatformLight.WindowBackground},
			{"dark: label colour on the window background", tokens.PlatformDark.Label, tokens.PlatformDark.WindowBackground},
			{"light: secondary label colour on the window background", tokens.PlatformLight.SecondaryLabel, tokens.PlatformLight.WindowBackground},
			{"dark: secondary label colour on the window background", tokens.PlatformDark.SecondaryLabel, tokens.PlatformDark.WindowBackground},
			{"light: white title on the default button", white, tokens.PlatformLight.DefaultButtonFill},
			{"dark: white title on the default button", white, tokens.PlatformDark.DefaultButtonFill},
		} {
			gtx := typeset.WithFill(layout.Context{Ops: new(op.Ops)}, tc.fill)
			got := typeset.FontOn(gtx, role, font.Normal, tc.fg)
			if want := tokens.StemFont(typeset.Font(role, font.Normal), role.Size, tc.fg, tc.fill); got != want {
				t.Errorf("%v pt, %s: %v, want %v", role.Size, tc.name, got, want)
			}
		}
	}
	role := typ.SystemFont
	gtx := typeset.WithFill(layout.Context{Ops: new(op.Ops)}, tokens.PlatformLight.WindowBackground)
	if got := typeset.FontOn(gtx, role, font.Normal, tokens.PlatformLight.Label); got.Typeface == font.Typeface(role.Typeface) {
		t.Errorf("the label colour on the light window background at 13 pt draws in %q, the face as shipped", got.Typeface)
	}
}

// TestStemFontOnKeepsTheFontsOwnStyle: StemFontOn picks the stem weight of
// a ready font.Font by the fill on gtx and keeps its style and weight.
func TestStemFontOnKeepsTheFontsOwnStyle(t *testing.T) {
	p := tokens.PlatformLight
	italic := font.Font{Style: font.Italic, Weight: font.Bold}
	gtx := typeset.WithFill(layout.Context{Ops: new(op.Ops)}, p.WindowBackground)
	got := typeset.StemFontOn(gtx, italic, 13, p.Label)
	if want := tokens.StemFont(italic, 13, p.Label, p.WindowBackground); got != want {
		t.Errorf("on the window background: %v, want %v", got, want)
	}
	if got.Style != font.Italic || got.Weight != font.Bold {
		t.Errorf("the style and weight became %v %v", got.Style, got.Weight)
	}
	if got.Typeface == italic.Typeface {
		t.Errorf("the label colour on the light window background kept typeface %q, want a widened face's", got.Typeface)
	}
}
