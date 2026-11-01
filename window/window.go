// Package window scopes a theme to a window. A [Window] is an [mvu.Window]
// plus the rx.Observable[theme.Theme] that drives what it renders, which is
// what lets two windows in one process show different themes at once —
// light and dark, brand A and brand B, one document's theme and another's.
//
// It is the application's entry point to the whole theme runtime, and
// three lines wide: build an mvu window, wrap it with a theme, render.
//
//	w := window.New(mvu.NewWindow(app.Title("My app")), system.LiveTheme(time.Second))
//	err := w.Render(buildLayers(modelObs)).Wait()
//
// Below the window the theme contract is per component — every components
// component takes an rx.Observable[theme.Theme] — and [Window.Render]
// exists to hand that one observable to the builder that constructs them,
// so no application code threads it manually.
//
// Isolation is a property of the observables, not of this type. Two Windows
// given two different observables share nothing — including their OS poll
// loops, since each system.LiveTheme call builds its own. Two Windows given
// the same one still render independently, and since FX.5 they also share
// its poll loops: the theme streams are multicast, so however many windows
// and layers subscribe to one observable, each OS source is polled once per
// interval.
//
// [Window.Render] shadows the embedded [mvu.Window.Render] and takes a
// different argument: a build function, not layers. The embedded one is
// still reachable as w.Window.Render for a caller that has plain layers and
// no use for the theme, at the cost of the scoping this package exists for.
// Either way it runs until the window is destroyed and cannot be restarted
// — a second Render call on the same window is not supported; build a new
// window instead.
package window

import (
	"image/color"

	"gioui.org/layout"

	"github.com/reactivego/rx"

	"github.com/vibrantgio/mvu"
	"github.com/vibrantgio/theme/theme"
	"github.com/vibrantgio/theme/tokens"
	"github.com/vibrantgio/theme/typeset"
)

// Window pairs an [mvu.Window] with the theme observable that scopes
// its rendering. The Theme field is the only path by which the wrapped
// window's content learns of theme changes; constructing two Window
// values with two different Theme observables yields two fully isolated
// theme paths.
type Window struct {
	*mvu.Window
	Theme rx.Observable[theme.Theme]
}

// New wraps an [mvu.Window] with a theme observable. The returned Window
// holds theme by reference; later emissions on theme reach the build
// callback passed to [Window.Render] and to no other Window.
func New(w *mvu.Window, theme rx.Observable[theme.Theme]) *Window {
	return &Window{Window: w, Theme: theme}
}

// Render starts the wrapped [mvu.Window] event loop with layers built
// from this window's theme. The build callback receives this window's
// own Theme observable; sibling windows constructed with their own
// themes do not share state with it.
//
// Render shadows the embedded [mvu.Window.Render]; callers that want
// the raw layer-only signature can still reach it via w.Window.Render.
//
// Every layer lays out on the window background of the theme as its fill,
// [typeset.WithFill], so text that stands on the window is drawn at the stem
// weight of the text against the window background.
func (w *Window) Render(build func(theme rx.Observable[theme.Theme]) []rx.Observable[layout.Widget]) rx.Subscription {
	return w.Window.Render(OnWindowBackground(w.Theme, build(w.Theme)...)...)
}

// OnWindowBackground returns layers with the window background of th set as
// the fill each one lays out on, [typeset.WithFill].
func OnWindowBackground(th rx.Observable[theme.Theme], layers ...rx.Observable[layout.Widget]) []rx.Observable[layout.Widget] {
	fills := rx.Map(rx.SwitchMap(th, func(t theme.Theme) rx.Observable[tokens.PlatformColors] {
		return t.Platform
	}), func(p tokens.PlatformColors) color.NRGBA {
		return p.WindowBackground
	})
	out := make([]rx.Observable[layout.Widget], len(layers))
	for i, l := range layers {
		out[i] = rx.Map(rx.CombineLatest2(l, fills), func(n rx.Tuple2[layout.Widget, color.NRGBA]) layout.Widget {
			draw, fill := n.First, n.Second
			if draw == nil {
				return nil
			}
			return func(gtx layout.Context) layout.Dimensions {
				return draw(typeset.WithFill(gtx, fill))
			}
		})
	}
	return out
}
