package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	t.Run("health check should return 200", func(t *testing.T) {
		a := api{}

		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		rec := httptest.NewRecorder()

		a.Routes().ServeHTTP(rec, req)

		statusCode := rec.Result().StatusCode

		if statusCode != 200 {
			t.Errorf("got %d, want %d", statusCode, 200)
		}
	})
}
