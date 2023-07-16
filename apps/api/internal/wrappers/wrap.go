package wrappers

import "net/http"

type Wrapper func(http.Handler) http.Handler

func Wrap(base http.Handler, wrappers ...Wrapper) http.Handler {
	for _, wrapper := range wrappers {
		base = wrapper(base)
	}

	return base
}
