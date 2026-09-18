package export

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/vibrantgio/theme/tokens"
)

// lowerCamel spells a Go field name the way theme.json spells it: the same
// word run with its first letter lowered.
func lowerCamel(name string) string {
	return strings.ToLower(name[:1]) + name[1:]
}

// TestEveryDensityFieldReachesTheSheetAndTheFile walks tokens.Density by
// reflection and fails if a field reaches neither output. A metric the token
// set carries and the bundle drops is a length a consumer cannot ask for and
// a number theme.json cannot reproduce — and the bundle has lost three of
// them that way, each added to the Go struct through the field walkers and
// each left out of the two emitters here because both spell their fields out
// by hand. Spelling them out is deliberate: the order and the prose per field
// are the sheet's, not reflection's. What this walk adds is that the hand
// cannot forget one.
func TestEveryDensityFieldReachesTheSheetAndTheFile(t *testing.T) {
	cssNames := map[string]func(tokens.Density) float32{}
	for _, m := range densityMetrics {
		cssNames[m.name] = m.pick
	}

	// theme.json's metrics, read back off the emitted file rather than off the
	// struct, so a field that carries no json tag is caught too.
	raw, err := json.Marshal(densityMetricsOf(tokens.Comfortable))
	if err != nil {
		t.Fatalf("marshal DensityMetrics: %v", err)
	}
	var emitted map[string]float64
	if err := json.Unmarshal(raw, &emitted); err != nil {
		t.Fatalf("unmarshal DensityMetrics: %v", err)
	}

	rt := reflect.TypeOf(tokens.Density{})
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Type.Kind() != reflect.Float32 {
			t.Fatalf("tokens.Density.%s is a %s; this walk reads the metrics as dp lengths",
				field.Name, field.Type)
		}
		name := kebab(field.Name)
		pick, ok := cssNames[name]
		if !ok {
			t.Errorf("tokens.Density.%s has no --density-%s in the sheet: a metric the token set carries and styles.css drops is a length no rule can reference",
				field.Name, name)
			continue
		}
		// Both settings, because the .compact block is generated from the
		// same entry and a pick reading the wrong field would show up in
		// exactly one of them.
		for _, d := range []struct {
			label string
			d     tokens.Density
		}{{"comfortable", tokens.Comfortable}, {"compact", tokens.Compact}} {
			want := reflect.ValueOf(d.d).Field(i).Interface().(float32)
			if got := pick(d.d); got != want {
				t.Errorf("--density-%s reads %v at %s, want tokens.Density.%s's %v",
					name, got, d.label, field.Name, want)
			}
		}

		key := lowerCamel(field.Name)
		got, ok := emitted[key]
		if !ok {
			t.Errorf("tokens.Density.%s has no %q in theme.json: a metric the file drops is a number the file cannot reproduce",
				field.Name, key)
			continue
		}
		if want := float64(reflect.ValueOf(tokens.Comfortable).Field(i).Interface().(float32)); got != want {
			t.Errorf("theme.json %s = %v, want tokens.Comfortable's %v", key, got, want)
		}
	}

	// The chip height is the one entry with no field behind it: it is a
	// method, so that the relation to the control height is stated once and
	// no Density value can carry a chip height that has come loose from it.
	if _, ok := cssNames["chip-height"]; !ok {
		t.Error("the sheet drops --density-chip-height, which tokens.Density states as a method rather than a field")
	}
}
