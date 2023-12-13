package config

import (
	"errors"
	"fmt"
)

var ErrAliasNotFound = errors.New("not found")

type ErrInvalidAliasKey struct {
	value string
}

func (e ErrInvalidAliasKey) Error() string {
	return fmt.Sprintf("Key is invalid. Did you mean `" + e.value + "`?")
}
