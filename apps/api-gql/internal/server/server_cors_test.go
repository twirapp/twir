package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGlobalCORS_allows_explicit_OPTIONS_routes_to_run(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(newGlobalCORS(router))

	called := false
	router.OPTIONS("/oauth/register", func(c *gin.Context) {
		called = true
		c.Header("Access-Control-Allow-Origin", "https://client.example")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodOptions, "/oauth/register", nil)
	req.Header.Set("Origin", "https://client.example")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	rec := httptest.NewRecorder()

	// When
	router.ServeHTTP(rec, req)

	// Then
	if !called {
		t.Fatalf("explicit OPTIONS handler was not reached: status=%d headers=%#v", rec.Code, rec.Header())
	}
	if rec.Code != http.StatusNoContent || rec.Header().Get("Access-Control-Allow-Origin") != "https://client.example" || rec.Header().Get("Access-Control-Allow-Methods") != "POST, OPTIONS" || rec.Header().Get("Access-Control-Allow-Headers") != "Content-Type" || rec.Header().Get("Access-Control-Allow-Credentials") != "" {
		t.Fatalf("explicit OPTIONS response = %d, headers = %#v", rec.Code, rec.Header())
	}
}

func TestGlobalCORS_keeps_fallback_preflight_for_unmatched_OPTIONS_routes(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(newGlobalCORS(router))

	req := httptest.NewRequest(http.MethodOptions, "/fallback", nil)
	req.Header.Set("Origin", "https://client.example")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	rec := httptest.NewRecorder()

	// When
	router.ServeHTTP(rec, req)

	// Then
	if rec.Code != http.StatusNoContent || rec.Header().Get("Access-Control-Allow-Origin") != "*" || rec.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatalf("fallback OPTIONS response = %d, headers = %#v", rec.Code, rec.Header())
	}
}
