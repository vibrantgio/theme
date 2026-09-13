// Command vg-tokens writes the Claude Design project layout — theme.json,
// styles.css, readme.md and the foundation pages under foundations/ — for a
// theme theme into a target directory.
//
// Usage:
//
//	vg-tokens [-out dir] [-color #rrggbb]
//
// With no flags it serialises theme.Default() into ./design. -color rebuilds
// the platform set's accent rows for another theme colour; everything else
// stays the default.
package main

import (
	"flag"
	"fmt"
	stdcolor "image/color"
	"os"
	"strings"

	"github.com/reactivego/rx"
	"github.com/vibrantgio/theme/export"
	"github.com/vibrantgio/theme/theme"
	"github.com/vibrantgio/theme/tokens"
)

func main() {
	out := flag.String("out", "design", "target directory for the generated project")
	colorHex := flag.String("color", "", "theme colour as #rrggbb (default: the platform's own accent)")
	flag.Parse()

	th := theme.Default()
	if *colorHex != "" {
		themeColor, err := parseColor(*colorHex)
		if err != nil {
			fatal(err)
		}
		th.Platform = rx.Of(tokens.PlatformLight.WithAccent(themeColor))
	}

	snap, err := export.Capture(th)
	if err != nil {
		fatal(err)
	}
	if err := export.Write(*out, snap); err != nil {
		fatal(err)
	}
	fmt.Printf("wrote %s: theme.json, styles.css, readme.md, foundations/{color,type,layout}.html\n", *out)
}

// parseColor accepts #rrggbb or rrggbb.
func parseColor(s string) (stdcolor.NRGBA, error) {
	hex, _ := strings.CutPrefix(s, "#")
	var r, g, b uint8
	if n, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b); n != 3 || err != nil || len(hex) != 6 {
		return stdcolor.NRGBA{}, fmt.Errorf("vg-tokens: -color %q: want #rrggbb", s)
	}
	return stdcolor.NRGBA{R: r, G: g, B: b, A: 0xff}, nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
