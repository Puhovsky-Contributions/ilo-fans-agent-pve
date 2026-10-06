package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/version"
)

func TestHandleVersion(t *testing.T) {
	origVersion := version.Version
	origCommit := version.Commit
	origDate := version.Date
	defer func() {
		version.Version = origVersion
		version.Commit = origCommit
		version.Date = origDate
	}()

	version.Version = "1.5.0"
	version.Commit = "abcdef123456"
	version.Date = "2026-10-06T12:00:00Z"

	tests := []struct {
		name       string
		args       []string
		wantHandled bool
		wantSubstr  string
		wantExact   string
	}{
		{
			name:        "subcommand version",
			args:        []string{"version"},
			wantHandled: true,
			wantSubstr:  "ilo-fans-agent-pve 1.5.0 (commit: abcdef1, built: 2026-10-06T12:00:00Z",
		},
		{
			name:        "subcommand version --short",
			args:        []string{"version", "--short"},
			wantHandled: true,
			wantExact:   "1.5.0\n",
		},
		{
			name:        "subcommand version -s",
			args:        []string{"version", "-s"},
			wantHandled: true,
			wantExact:   "1.5.0\n",
		},
		{
			name:        "flag -v",
			args:        []string{"-v"},
			wantHandled: true,
			wantSubstr:  "ilo-fans-agent-pve 1.5.0",
		},
		{
			name:        "flag --version",
			args:        []string{"--version"},
			wantHandled: true,
			wantSubstr:  "ilo-fans-agent-pve 1.5.0",
		},
		{
			name:        "flag -version",
			args:        []string{"-version"},
			wantHandled: true,
			wantSubstr:  "ilo-fans-agent-pve 1.5.0",
		},
		{
			name:        "flag -v --short",
			args:        []string{"-v", "--short"},
			wantHandled: true,
			wantExact:   "1.5.0\n",
		},
		{
			name:        "flag --version -s",
			args:        []string{"--version", "-s"},
			wantHandled: true,
			wantExact:   "1.5.0\n",
		},
		{
			name:        "not version",
			args:        []string{"token", "show"},
			wantHandled: false,
		},
		{
			name:        "empty args",
			args:        []string{},
			wantHandled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handled := handleVersion(tt.args, &buf)
			if handled != tt.wantHandled {
				t.Fatalf("expected handled=%v, got %v", tt.wantHandled, handled)
			}
			if !handled {
				return
			}
			out := buf.String()
			if tt.wantExact != "" && out != tt.wantExact {
				t.Errorf("got exact %q, want %q", out, tt.wantExact)
			}
			if tt.wantSubstr != "" && !strings.Contains(out, tt.wantSubstr) {
				t.Errorf("expected substring %q in output %q", tt.wantSubstr, out)
			}
		})
	}
}
