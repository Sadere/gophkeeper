package screens

import (
	"testing"

	"github.com/knz/catwalk"
)

func TestWelcomeModel_Login(t *testing.T) {
	t.Run("select login", func(t *testing.T) {
		m := NewWelcomeModel(&State{})

		catwalk.RunModel(t, "testdata/welcome/login_test", m)
	})

	t.Run("select register", func(t *testing.T) {
		m := NewWelcomeModel(&State{})

		catwalk.RunModel(t, "testdata/welcome/register_test", m)
	})
}
