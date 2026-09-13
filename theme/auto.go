package theme

import (
	"time"

	"github.com/reactivego/rx"
	"github.com/vibrantgio/theme/tokens"
)

// AutoLightDark returns an Observable that emits a new Theme every minute.
// Hours 7–17 (inclusive) use PlatformLight; all other hours use PlatformDark.
func AutoLightDark() rx.Observable[Theme] {
	return rx.Map(rx.Ticker(0, time.Minute), func(t time.Time) Theme {
		platform := tokens.PlatformLight
		if t.Hour() <= 6 || t.Hour() >= 18 {
			platform = tokens.PlatformDark
		}
		return Theme{
			Platform:   rx.Of(platform),
			Typography: rx.Of(tokens.DefaultTypography),
			Density:    rx.Of(tokens.Comfortable),
			Motion:     rx.Of(tokens.Motion),
			Spacing:    rx.Of(tokens.Spacing),
			Radius:     rx.Of(tokens.Radius),
			Elevation:  rx.Of(tokens.Elevation),
		}
	})
}
