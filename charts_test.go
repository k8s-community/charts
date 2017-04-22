package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k8s-community/charts/version"
	"github.com/takama/router"
)

func TestHealthHandler(t *testing.T) {
	r := router.New()
	r.GET("/healthz", healthz)
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	if trw.Body.String() != "Ok" {
		t.Error("Expected", "Ok", "got", trw.Body.String())
	}
	if trw.Code != 200 {
		t.Error("Expected status:", 200, "got", trw.Body.String())
	}
}

func TestInfoHandler(t *testing.T) {
	r := router.New()
	r.GET("/info", info)
	req, err := http.NewRequest("GET", "/info", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	inf := "{\n  \"commit\": \"" + version.COMMIT + "\",\n  \"repo\": \"" + version.REPO +
		"\",\n  \"version\": \"" + version.RELEASE + "\"\n}"
	if trw.Body.String() != inf {
		t.Error("Expected", inf, "got", trw.Body.String())
	}
	if trw.Code != 200 {
		t.Error("Expected status:", 200, "got", trw.Body.String())
	}
}
