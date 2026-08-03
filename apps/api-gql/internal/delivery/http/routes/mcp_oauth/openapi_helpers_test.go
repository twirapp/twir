package mcp_oauth

import (
	"sort"
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

type responseSetCheck struct {
	path     string
	method   string
	statuses []string
}

func assertResponseStatuses(t *testing.T, openapi *huma.OpenAPI, check responseSetCheck) {
	t.Helper()
	operation := mustOperation(t, openapi, check.path, check.method)
	got := map[string]struct{}{}
	for status := range operation.Responses {
		got[status] = struct{}{}
	}
	if len(got) != len(check.statuses) {
		t.Fatalf("%s %s response statuses = %v; want %v", check.method, check.path, sortedKeys(got), check.statuses)
	}
	for _, status := range check.statuses {
		if _, ok := got[status]; !ok {
			t.Fatalf("%s %s missing response %s; got %v", check.method, check.path, status, sortedKeys(got))
		}
	}
}

func sortedKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
