package mcp_oauth

type ErrorCode string

const (
	ErrorInvalidRequest        ErrorCode = "invalid_request"
	ErrorInvalidClientMetadata ErrorCode = "invalid_client_metadata"
	ErrorInvalidClient         ErrorCode = "invalid_client"
	ErrorInvalidGrant          ErrorCode = "invalid_grant"
	ErrorInvalidScope          ErrorCode = "invalid_scope"
	ErrorInvalidToken          ErrorCode = "invalid_token"
	ErrorAccessDenied          ErrorCode = "access_denied"
)

type OAuthError struct {
	Code        ErrorCode
	Description string
}

func (err *OAuthError) Error() string { return string(err.Code) + ": " + err.Description }
func oauthError(code ErrorCode, description string) error {
	return &OAuthError{Code: code, Description: description}
}
