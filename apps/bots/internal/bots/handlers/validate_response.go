package handlers

import (
	"errors"
	"strings"
)

func ValidateResponseSlashes(response string) error {
	if !strings.HasPrefix(response, "/") || strings.HasPrefix(response, "/me") {
		return nil
	} else if strings.HasPrefix(response, "/") {
		return errors.New("slash commands except /me and is disallowed. This response won't be ever sent")
	} else if strings.HasPrefix(response, ".") {
		return errors.New(`message cannot start from "." symbol`)
	} else {
		return nil
	}
}
