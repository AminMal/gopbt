package gopbt

import (
	"testing"

	"github.com/AminMal/gopbt/gen"
)

func getGeneratorsLen(s *Session) int {
	return len(s.mapping.generatorMapping)
}

func TestSessionGenerators(t *testing.T) {
	s := NewSession()

	numGens := 0

	if getGeneratorsLen(s) > 0 {
		t.Error("empty new session should not contain any generator by default")
	}

	stringGen := gen.OneOf("first", "second", "third", "4th")
	SetGen(s, stringGen)
	numGens++

	intGen := gen.Between(-1000, 1000)
	SetGen(s, intGen)
	numGens++

	if getGeneratorsLen(s) != numGens {
		t.Errorf("added %d generators to session, session has %d generators", numGens, getGeneratorsLen(s))
	}
}

func TestSessionGeneratorAutoAppend(t *testing.T) {
	s := NewSession()
	s.SupportAdhocGenerators = true

	if getGeneratorsLen(s) != 0 {
		t.Error("empty new session should not contain any generator by default")
	}

	randomPropertyOfSomethingCheck := func (i int, r rune, s string) bool {
		// we should have int, rune and string generators now
		return true
	}

	if err := s.Check(randomPropertyOfSomethingCheck, nil); err != nil {
		t.Fatal("session check returned error for a propery which always returns true")
	}

	if getGeneratorsLen(s) != 3 {
		t.Errorf("added 3 generators to session, session has %d generators", getGeneratorsLen(s))
	}
}
