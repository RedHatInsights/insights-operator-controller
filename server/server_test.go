package server

import (
	"net/http"
	"testing"
)

func TestMainEndpoint(t *testing.T) {
	response, err := http.Get(API_URL)
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}
}
