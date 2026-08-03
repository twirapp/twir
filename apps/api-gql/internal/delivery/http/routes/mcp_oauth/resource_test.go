package mcp_oauth

import (
	"net/http"
	"testing"
)

func TestHandler_token_openapi_marks_resource_optional(t *testing.T) {
	// Given
	handler := newTestHandler(t)
	_ = handler.router()
	openapi := handler.api.OpenAPI()

	// When
	operation := mustOperation(t, openapi, "/oauth/token", http.MethodPost)
	schema := resolveSchema(t, openapi, operation.RequestBody.Content["application/x-www-form-urlencoded"].Schema)

	// Then
	for _, required := range schema.Required {
		if required == "resource" {
			t.Fatal("token resource form field is required in OpenAPI")
		}
	}
}
