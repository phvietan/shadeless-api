package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"shadeless-api/main/config"
	"shadeless-api/main/libs"
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

func createRandomFile() (string, string) {
	fileName := libs.RandomString(16) + ".txt"
	content := libs.RandomString(200)
	if err := ioutil.WriteFile("./files/"+fileName, []byte(content), 0755); err != nil {
		panic("Unable to write file")
	}
	return fileName, content
}

func removeFile(name string) {
	os.Remove(name)
}

// File exist test
func TestFileServeExist(t *testing.T) {
	for i := 0; i < 1; i++ {
		fileName, content := createRandomFile()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/files/"+fileName, nil)
		router.ServeHTTP(w, req)

		fmt.Println(w.Code)
		fmt.Println(w.Body.String())
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, content, w.Body.String())

		acao := w.Header().Get("Access-Control-Allow-Origin")
		ct := w.Header().Get("Content-Type")
		assert.Equal(t, acao, config.GetInstance().GetFrontendUrl())
		assert.Equal(t, ct, "application/octet-stream")

		removeFile("./files/" + fileName)
	}
}

// None file exist test
func TestFileServeNonExist(t *testing.T) {
	for i := 0; i < 10; i++ {
		fileName := libs.RandomString(32)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/files/"+fileName, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Equal(t, "", w.Body.String())
	}
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
