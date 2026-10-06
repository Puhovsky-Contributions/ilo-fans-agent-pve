package version

import (
	"strings"
	"testing"
)

func TestGetDefault(t *testing.T) {
	info := Get()
	if info.Version == "" {
		t.Fatalf("expected non-empty version, got empty")
	}
	if info.GoVersion == "" {
		t.Fatalf("expected non-empty GoVersion")
	}
	if info.Platform == "" {
		t.Fatalf("expected non-empty Platform")
	}
	if !strings.Contains(info.Platform, "/") {
		t.Fatalf("expected platform format os/arch, got %s", info.Platform)
	}
}

func TestVersionFormatting(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origDate := Date
	defer func() {
		Version = origVersion
		Commit = origCommit
		Date = origDate
	}()

	Version = "v1.2.3"
	Commit = "0123456789abcdef0123456789abcdef01234567"
	Date = "2026-10-06T20:00:00Z"

	info := Get()
	if info.Short() != "1.2.3" {
		t.Fatalf("expected Short() to be 1.2.3, got %s", info.Short())
	}
	if info.Version != "1.2.3" {
		t.Fatalf("expected Version to be 1.2.3, got %s", info.Version)
	}

	str := info.String()
	if !strings.Contains(str, "ilo-fans-agent-pve 1.2.3") {
		t.Fatalf("expected String() to contain app name and version, got %s", str)
	}
	if !strings.Contains(str, "built: 2026-10-06T20:00:00Z") {
		t.Fatalf("expected String() to contain build date, got %s", str)
	}
	if !strings.Contains(str, "commit: 0123456") {
		t.Fatalf("expected String() to contain commit sha, got %s", str)
	}
}

func TestShortWithoutVPrefix(t *testing.T) {
	origVersion := Version
	defer func() {
		Version = origVersion
	}()

	Version = "2.0.1"
	info := Get()
	if info.Short() != "2.0.1" {
		t.Fatalf("expected Short() to be 2.0.1, got %s", info.Short())
	}
}
