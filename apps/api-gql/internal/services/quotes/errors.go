package quotes

import "github.com/twirapp/twir/libs/errors"

var ErrQuoteNotFound = errors.NewNotFoundError("Quote with this ID was not found")
