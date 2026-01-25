package shortened_urls

import (
	"fmt"
)

var ErrNotFound = fmt.Errorf("shortened url not found")
var ErrShortIDAlreadyExists = fmt.Errorf("shortened url short id already exists")
