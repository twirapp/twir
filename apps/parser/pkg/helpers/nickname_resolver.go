package helpers

import (
	"strings"
)

func ResolveDisplayName(login string, displayName string) string {
	if strings.ToLower(login) == strings.ToLower(displayName) {
		return displayName
	}

	return login
}
