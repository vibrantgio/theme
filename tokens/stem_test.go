package tokens

import (
	"testing"

	"github.com/vibrantgio/font/roboto"
)

// TestOnlyRobotoRegularTakesTheStemWeight: at a stem weight, one face of the
// collection changes, Roboto Regular upright, and the rest are the entries
// as shipped.
func TestOnlyRobotoRegularTakesTheStemWeight(t *testing.T) {
	faces := DefaultTypography.Faces
	got := stemFaces(faces, StemWeightLight)
	changed := 0
	for i := range faces {
		if got[i].Font != faces[i].Font {
			t.Fatalf("face %d is %v, want %v", i, got[i].Font, faces[i].Font)
		}
		if got[i].Face != faces[i].Face {
			changed++
			if got[i].Font != roboto.RegularNormal {
				t.Errorf("face %v changed, want only %v", got[i].Font, roboto.RegularNormal)
			}
		}
	}
	if changed != 1 {
		t.Errorf("%d faces changed, want 1", changed)
	}
}

// TestEveryCopySharesTheStemShaper: a copy of a typography hands back the
// shaper of its source for each stem weight, and stem weight zero is the
// shaper as shipped.
func TestEveryCopySharesTheStemShaper(t *testing.T) {
	a := DefaultTypography
	b := DefaultTypography
	if a.DeterministicShaperFor(StemWeightLight) != b.DeterministicShaperFor(StemWeightLight) {
		t.Error("two copies built two shapers at one stem weight")
	}
	if a.DeterministicShaperFor(0) != a.DeterministicShaper() {
		t.Error("stem weight zero is not the shaper as shipped")
	}
	if a.DeterministicShaperFor(StemWeightLight) == a.DeterministicShaper() {
		t.Error("the light stem weight shares the shaper as shipped")
	}
	if SchemeStemWeight(false) != StemWeightLight || SchemeStemWeight(true) != StemWeightDark {
		t.Error("SchemeStemWeight does not pick by scheme")
	}
}
