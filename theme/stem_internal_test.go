package theme

import (
	"testing"
	"time"

	"github.com/vibrantgio/theme/tokens"
)

// TestEachHourDrawsAtItsSchemesStemWeight: AutoLightDark's typography draws
// at the light scheme's stem weight in the light hours and at the dark
// scheme's in the dark hours, with the colours of that scheme.
func TestEachHourDrawsAtItsSchemesStemWeight(t *testing.T) {
	for _, tc := range []struct {
		hour int
		set  tokens.PlatformColors
		stem tokens.StemWeight
	}{
		{12, tokens.PlatformLight, tokens.StemWeightLight},
		{23, tokens.PlatformDark, tokens.StemWeightDark},
	} {
		th := atHour(time.Date(2026, 1, 1, tc.hour, 0, 0, 0, time.Local))
		set, err := th.Platform.First()
		if err != nil {
			t.Fatalf("hour %d: platform: %v", tc.hour, err)
		}
		typ, err := th.Typography.First()
		if err != nil {
			t.Fatalf("hour %d: typography: %v", tc.hour, err)
		}
		if set != tc.set {
			t.Errorf("hour %d: the colours are not the scheme's set", tc.hour)
		}
		if typ.StemWeight() != tc.stem {
			t.Errorf("hour %d: the typography draws at stem weight %v, want %v", tc.hour, typ.StemWeight(), tc.stem)
		}
	}
}

// TestTheDefaultThemeDrawsAtTheLightStemWeight: Default's scheme is light,
// and its typography draws at the light scheme's stem weight.
func TestTheDefaultThemeDrawsAtTheLightStemWeight(t *testing.T) {
	typ, err := Default().Typography.First()
	if err != nil {
		t.Fatalf("typography: %v", err)
	}
	if typ.StemWeight() != tokens.StemWeightLight {
		t.Errorf("the default theme draws at stem weight %v, want %v", typ.StemWeight(), tokens.StemWeightLight)
	}
	light := tokens.DefaultTypography.WithStemWeight(tokens.StemWeightLight)
	if typ.Shaper() != light.Shaper() {
		t.Error("the default theme built its own shaper at the light stem weight")
	}
}
