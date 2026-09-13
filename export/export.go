package export

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vibrantgio/theme/theme"
	"github.com/vibrantgio/theme/tokens"
)

// Snapshot is one resolved theme: the first emission of each theme.Theme
// observable, with the paired dark appearance's colour set. It is the input
// Write serialises.
type Snapshot struct {
	// PlatformLight and PlatformDark are the platform's colour set in
	// the two appearances: the set the theme emitted, and its counterpart
	// under the same accent. The accent is the only row a machine's own
	// settings move, so the pair is the two recorded sets carrying the
	// emitted set's accent — which on the default accent is the two recorded
	// sets exactly.
	PlatformLight, PlatformDark tokens.PlatformColors

	Typography tokens.Typography
	Density    tokens.Density
	Motion     tokens.MotionScale
	Spacing    tokens.SpacingScale
	Radius     tokens.RadiusScale
	Elevation  tokens.ElevationScale
}

// Capture collects the first emission of each observable a serialisation
// needs — Platform, Typography, Density, Motion, Spacing, Radius and
// Elevation — into a Snapshot. (Type is not consumed: it duplicates
// Typography's sizes.)
//
// The colour emission is the platform's set for the light appearance; its
// dark counterpart is the recorded dark set carrying the emitted accent,
// since the accent is the only row a machine's own settings move.
// The density emission must be one of the two published settings —
// tokens.Comfortable or tokens.Compact — because theme.json records density
// as a named setting plus both settings' metrics, not as free-form numbers.
func Capture(th theme.Theme) (Snapshot, error) {
	var s Snapshot
	if th.Platform == nil || th.Typography == nil || th.Density == nil || th.Motion == nil || th.Spacing == nil || th.Radius == nil || th.Elevation == nil {
		return s, fmt.Errorf("export: Capture: theme has nil observables; every consumed field of theme.Theme must be set")
	}
	var err error
	if s.PlatformLight, err = th.Platform.First(); err != nil {
		return s, fmt.Errorf("export: Capture: Platform: %w", err)
	}
	s.PlatformDark = tokens.PlatformDark.WithAccent(s.PlatformLight.ControlAccent)
	if s.Typography, err = th.Typography.First(); err != nil {
		return s, fmt.Errorf("export: Capture: Typography: %w", err)
	}
	if s.Density, err = th.Density.First(); err != nil {
		return s, fmt.Errorf("export: Capture: Density: %w", err)
	}
	if s.Motion, err = th.Motion.First(); err != nil {
		return s, fmt.Errorf("export: Capture: Motion: %w", err)
	}
	if s.Spacing, err = th.Spacing.First(); err != nil {
		return s, fmt.Errorf("export: Capture: Spacing: %w", err)
	}
	if s.Radius, err = th.Radius.First(); err != nil {
		return s, fmt.Errorf("export: Capture: Radius: %w", err)
	}
	if s.Elevation, err = th.Elevation.First(); err != nil {
		return s, fmt.Errorf("export: Capture: Elevation: %w", err)
	}

	if _, ok := densitySetting(s.Density); !ok {
		return s, fmt.Errorf("export: Capture: the density emission is neither tokens.Comfortable nor tokens.Compact; theme.json records density as a named setting, so only the published settings are reproducible")
	}
	return s, nil
}

// densitySetting names a density emission: the setting string theme.json
// records, and whether the emission is one of the two published settings.
func densitySetting(d tokens.Density) (string, bool) {
	switch d {
	case tokens.Comfortable:
		return "comfortable", true
	case tokens.Compact:
		return "compact", true
	}
	return "", false
}

// Write renders s into dir as the full Claude Design project layout —
// theme.json, styles.css, readme.md and the foundation pages under
// foundations/ — creating directories as needed. Existing files are
// overwritten: the tree is generated output, regenerated whole.
func Write(dir string, s Snapshot) error {
	if err := os.MkdirAll(filepath.Join(dir, "foundations"), 0o755); err != nil {
		return fmt.Errorf("export: Write: %w", err)
	}
	js, err := themeJSON(s)
	if err != nil {
		return fmt.Errorf("export: Write: %w", err)
	}
	files := []struct {
		name    string
		content []byte
	}{
		{"theme.json", js},
		{"styles.css", []byte(stylesCSS(s))},
		{"readme.md", []byte(readmeMD(s))},
		{filepath.Join("foundations", "color.html"), []byte(colorHTML(s))},
		{filepath.Join("foundations", "type.html"), []byte(typeHTML(s))},
		{filepath.Join("foundations", "layout.html"), []byte(layoutHTML(s))},
	}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f.name), f.content, 0o644); err != nil {
			return fmt.Errorf("export: Write: %w", err)
		}
	}
	return nil
}
