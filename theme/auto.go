package theme

import (
	"time"

	"github.com/reactivego/rx"
	"github.com/vibrantgio/theme/tokens"
)

// AutoLightDark returns an Observable that emits a new Theme every minute.
// Hours 7–17 (inclusive) use PlatformLight; all other hours use PlatformDark.
// The typography draws at the stem weight of the scheme in force.
func AutoLightDark() rx.Observable[Theme] {
	return rx.Map(rx.Ticker(0, time.Minute), atHour)
}

// atHour is the theme AutoLightDark emits at the time t.
func atHour(t time.Time) Theme {
	dark := t.Hour() <= 6 || t.Hour() >= 18
	platform := tokens.PlatformLight
	if dark {
		platform = tokens.PlatformDark
	}
	return Theme{
		Platform:   rx.Of(platform),
		Typography: rx.Of(tokens.DefaultTypography.WithStemWeight(tokens.SchemeStemWeight(dark))),
		Density:    rx.Of(tokens.Comfortable),
		Motion:     rx.Of(tokens.Motion),
		Spacing:    rx.Of(tokens.Spacing),
		Radius:     rx.Of(tokens.Radius),
	}
}
