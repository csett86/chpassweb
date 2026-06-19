package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandleStart_NoAuthHeader(t *testing.T) {
	originalStore := store
	defer func() { store = originalStore }()
	store = NewSessionStore(1 * time.Minute)

	req := httptest.NewRequest("GET", "/start", nil)
	w := httptest.NewRecorder()
	handleStart(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Authentication required") {
		t.Error("Expected 'Authentication required' in response")
	}
}

func TestHandleStart_EmptyUsername(t *testing.T) {
	originalStore := store
	defer func() { store = originalStore }()
	store = NewSessionStore(1 * time.Minute)

	req := httptest.NewRequest("GET", "/start", nil)
	req.Header.Set("REMOTE_USER", "")
	w := httptest.NewRecorder()
	handleStart(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleRespond_NoSession(t *testing.T) {
	originalStore := store
	defer func() { store = originalStore }()
	store = NewSessionStore(1 * time.Minute)

	req := httptest.NewRequest("POST", "/respond", strings.NewReader("response=test"))
	w := httptest.NewRecorder()
	handleRespond(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Session expired") {
		t.Error("Expected 'Session expired' in response")
	}
}

func TestRenderPrompt(t *testing.T) {
	w := httptest.NewRecorder()
	renderPrompt(w, "Enter current password:")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Enter current password:") {
		t.Error("Expected prompt in response")
	}
	if !strings.Contains(body, "<form") {
		t.Error("Expected form in response")
	}
}

func TestRenderError(t *testing.T) {
	w := httptest.NewRecorder()
	renderError(w, "Password too short")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Password too short") {
		t.Error("Expected error message in response")
	}
	if !strings.Contains(body, "class=\"error\"") {
		t.Error("Expected error class in response")
	}
}

func TestRenderSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	renderSuccess(w)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Success") {
		t.Error("Expected success message in response")
	}
	if !strings.Contains(body, "class=\"success\"") {
		t.Error("Expected success class in response")
	}
}
