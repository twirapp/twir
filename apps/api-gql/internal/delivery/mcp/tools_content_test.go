package mcp

import (
	"strings"
	"testing"
)

func TestVariableScriptGuideCoversSandboxUsage(t *testing.T) {
	required := []string{
		"top-level await",
		"Always return",
		"fetch(url, options)",
		"storage and twir.storage",
		"twir.secrets.get(name)",
		"lodash as _",
		"5 seconds",
		"Node/Bun APIs",
		`"$(sender.displayName)"`,
		"list_builtin_variables",
		"evaluate_variable",
	}

	for _, text := range required {
		if !strings.Contains(variableScriptGuide, text) {
			t.Errorf("SCRIPT variable guide is missing %q", text)
		}
	}
}
