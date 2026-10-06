package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/auth"
	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/config"
	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/version"
)

func TestHandleHealth(t *testing.T) {
	origVersion := version.Version
	defer func() {
		version.Version = origVersion
	}()
	version.Version = "1.2.3"

	srv := New(config.Config{}, auth.NewStore("/nonexistent"))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %s", body["status"])
	}
	if body["version"] != "1.2.3" {
		t.Errorf("expected version 1.2.3, got %s", body["version"])
	}
}
