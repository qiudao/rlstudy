package env

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer_Info(t *testing.T) {
	s := NewServer(10, 42)
	req := httptest.NewRequest("GET", "/info", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var info InfoResponse
	json.NewDecoder(w.Body).Decode(&info)
	if info.Arms != 10 {
		t.Errorf("expected 10 arms, got %d", info.Arms)
	}
}

func TestServer_ResetAndStep(t *testing.T) {
	s := NewServer(10, 42)

	// Reset
	req := httptest.NewRequest("POST", "/reset", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("reset: expected 200, got %d", w.Code)
	}

	// Step
	body := strings.NewReader(`{"action":0}`)
	req = httptest.NewRequest("POST", "/step", body)
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("step: expected 200, got %d", w.Code)
	}

	var resp StepResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Reward != resp.Reward { // NaN check
		t.Error("reward is NaN")
	}
}

func TestServer_StepInvalidAction(t *testing.T) {
	s := NewServer(10, 42)

	body := strings.NewReader(`{"action":10}`)
	req := httptest.NewRequest("POST", "/step", body)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestServer_StepBadJSON(t *testing.T) {
	s := NewServer(10, 42)

	body := strings.NewReader(`not json`)
	req := httptest.NewRequest("POST", "/step", body)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
