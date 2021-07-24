package main

import (
	"net/http"
	"net/http/httptest"
	"shadeless-api/main/config"
	"testing"

	"github.com/go-playground/assert/v2"
)

var router = spawnApp()

func TestHealthCheckRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Health check ok", w.Body.String())
}

func TestResponseHeader(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dscnjdskcnjsdkcndfsiucdfnuidfncui", nil)
	router.ServeHTTP(w, req)

	acao := w.Header().Get("Access-Control-Allow-Origin")
	acam := w.Header().Get("Access-Control-Allow-Methods")
	acah := w.Header().Get("Access-Control-Allow-Headers")
	assert.Equal(t, acao, config.GetInstance().GetFrontendUrl())
	assert.Equal(t, acam, "POST, GET, OPTIONS, PUT, DELETE")
	assert.Equal(t, acah, "Content-Type, Content-Length, Accept-Encoding")
}

func TestRequestMethodOptions(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/dscnjdskcnjsdkcndfsiucdfnuidfncui", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, "", w.Body.String())
	assert.Equal(t, http.StatusNoContent, w.Code)
}
