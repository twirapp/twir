package short_links_custom_domains

import "errors"

var ErrNotFound = errors.New("custom domain not found")
var ErrDomainAlreadyExists = errors.New("domain already exists")
var ErrUserAlreadyHasDomain = errors.New("user already has a custom domain")
