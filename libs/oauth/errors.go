package oauth

import "errors"

var (
	ErrInvalidOption     = errors.New("oauth: invalid option")
	ErrInvalidCredential = errors.New("oauth: invalid credential")
	ErrLeaseLost         = errors.New("oauth: lease lost")
	ErrClosed            = errors.New("oauth: closed")
	ErrLoad              = errors.New("oauth: load failed")
	ErrRefresh           = errors.New("oauth: refresh failed")
	ErrCommit            = errors.New("oauth: commit failed")
	ErrCoordinator       = errors.New("oauth: coordinator failed")
)
