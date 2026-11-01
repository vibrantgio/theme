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

// TestEveryCopySharesTheStemShaper: copies of a typography at one stem
// weight hand back one shaper, and stem weight zero is the shaper as
// shipped.
func TestEveryCopySharesTheStemShaper(t *testing.T) {
	a := DefaultTypography.WithStemWeight(StemWeightLight)
	b := DefaultTypography.WithStemWeight(StemWeightLight)
	if a.DeterministicShaper() != b.DeterministicShaper() || a.Shaper() != b.Shaper() {
		t.Error("two copies built two shapers at one stem weight")
	}
	zero := DefaultTypography.WithStemWeight(0)
	if zero.DeterministicShaper() != DefaultTypography.DeterministicShaper() {
		t.Error("stem weight zero is not the shaper as shipped")
	}
	if a.DeterministicShaper() == DefaultTypography.DeterministicShaper() {
		t.Error("the light stem weight shares the shaper as shipped")
	}
	if a.StemWeight() != StemWeightLight || DefaultTypography.StemWeight() != 0 {
		t.Error("StemWeight does not report the stem weight set")
	}
	if SchemeStemWeight(false) != StemWeightLight || SchemeStemWeight(true) != StemWeightDark {
		t.Error("SchemeStemWeight does not pick by scheme")
	}
}

// TestTheStemWeightStaysThroughTheFaceOptions: adding emoji or a code face
// to a typography keeps the stem weight it is drawn at.
func TestTheStemWeightStaysThroughTheFaceOptions(t *testing.T) {
	light := DefaultTypography.WithStemWeight(StemWeightLight)
	for name, got := range map[string]Typography{
		"WithEmoji":                   light.WithEmoji(),
		"WithCodeFace JetBrains":      light.WithCodeFace(CodeFaceJetBrains),
		"WithCodeFace then WithEmoji": light.WithCodeFace(CodeFaceJetBrains).WithEmoji(),
		"WithFaces":                   light.WithFaces(),
	} {
		if got.StemWeight() != StemWeightLight {
			t.Errorf("%s draws at stem weight %v, want %v", name, got.StemWeight(), StemWeightLight)
		}
	}
}
