package vkvideoprobe

import (
	"bytes"
	"regexp"
	"sort"
)

var sensitiveJSONFieldPattern = regexp.MustCompile(`(?i)("(?:[a-z0-9_-]*(?:token|secret|password|authorization|credential)[a-z0-9_-]*|api[_-]?key|cookie)"\s*:\s*)"(?:\\.|[^"\\])*"`)
var sensitiveURLValuePattern = regexp.MustCompile(`(?i)([?&](?:[a-z0-9_-]*(?:token|secret|password|authorization|credential)[a-z0-9_-]*|api[_-]?key|cookie)=)[^&#"\\\s]*`)
var authorizationValuePattern = regexp.MustCompile(`(?i)(\b(?:authorization|bearer|basic)\b\s*[:=]?\s*)(?:bearer\s+|basic\s+)?[^\s,;"'}]+`)
var jwtPattern = regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{2,}\.[A-Za-z0-9_-]{2,}\.[A-Za-z0-9_-]{2,}\b`)

func RedactFrame(frame []byte, secrets []string) []byte {
	redacted := redactKnownSecrets(frame, secrets)
	redacted = sensitiveJSONFieldPattern.ReplaceAll(redacted, []byte(`${1}"[REDACTED]"`))
	redacted = sensitiveURLValuePattern.ReplaceAll(redacted, []byte(`${1}[REDACTED]`))
	redacted = authorizationValuePattern.ReplaceAll(redacted, []byte(`${1}[REDACTED]`))

	return jwtPattern.ReplaceAll(redacted, []byte("[REDACTED]"))
}

func redactKnownSecrets(frame []byte, secrets []string) []byte {
	sortedSecrets := append([]string(nil), secrets...)
	sort.Slice(sortedSecrets, func(left, right int) bool {
		return len(sortedSecrets[left]) > len(sortedSecrets[right])
	})

	redacted := bytes.Clone(frame)
	for _, secret := range sortedSecrets {
		if secret != "" {
			redacted = bytes.ReplaceAll(redacted, []byte(secret), []byte("[REDACTED]"))
		}
	}

	return redacted
}
