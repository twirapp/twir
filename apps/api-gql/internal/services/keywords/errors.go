package keywords

import (
	"github.com/twirapp/twir/libs/errors"
)

var ErrKeywordNotFound = errors.NewNotFoundError("Keyword with this ID was not found")
