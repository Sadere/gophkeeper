package screens

import (
	"testing"

	"github.com/knz/catwalk"
)

func TestRootModel(t *testing.T) {
	m := NewRootModel(&State{})

	catwalk.RunModel(t, "testdata/root/tests", m)
}
