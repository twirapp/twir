package helpers

import "fmt"

type BusinessError struct {
	Err error

	StatusCode int
	Messages   []string
}

func (r BusinessError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func CreateBusinessErrorWithMessage(code int, message string) BusinessError {
	return BusinessError{
		StatusCode: code,
		Messages:   []string{message},
	}
}
