package model

import (
	"errors"
	"fmt"
)

var (
	ErrBadCredentials = errors.New("bad credentials")

	ErrNoSecrets       = errors.New("no secrets were added yet")
	ErrWrongSecretType = errors.New("invalid secret type")
)

type ErrUserExists struct {
	Login string
}

func (e *ErrUserExists) Error() string {
	return fmt.Sprintf("user with login '%s' has already registered", e.Login)
}

func (e *ErrUserExists) Is(tgt error) bool {
	target, ok := tgt.(*ErrUserExists)
	if !ok {
		return false
	}
	return e.Login == target.Login
}
