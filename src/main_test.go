package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSearchSongs(t *testing.T) {
	// Set up the Echo server
	e := echo.New()
	e.GET("/search", searchSongs)

	// Create a new HTTP request to the search endpoint
	req := httptest.NewRequest(http.MethodGet, "/search?query=test", nil)

	// Record the response
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Use Testify's assert package for assertions
	assert.Equal(t, http.StatusOK, rec.Code, "Status code should be 200")
}
