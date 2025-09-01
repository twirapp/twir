package channels_commands_prefix

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("prefix not found")
	ErrAlreadyExists = errors.New("prefix already exists")
)
