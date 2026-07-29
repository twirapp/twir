package oauth

import (
	"reflect"
	"slices"
)

func isNil(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

func cloneCredential(credential Credential) Credential {
	credential.Scopes = slices.Clone(credential.Scopes)
	return credential
}

func cloneRefreshResult(result RefreshResult) RefreshResult {
	result.Scopes = slices.Clone(result.Scopes)
	return result
}
