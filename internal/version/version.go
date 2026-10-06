package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

var (
	// Version is the application version, optionally injected at build time via -ldflags.
	Version = "dev"
	// Commit is the git commit SHA, optionally injected at build time via -ldflags.
	Commit = "none"
	// Date is the build timestamp, optionally injected at build time via -ldflags.
	Date = "unknown"
)

// Info holds version and build metadata.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"goVersion"`
	Platform  string `json:"platform"`
}

// Get returns the populated Info structure.
// If Commit or Date were not injected via ldflags, it attempts to read
// them from Go's runtime/debug build info.
func Get() Info {
	v := Version
	c := Commit
	d := Date

	if info, ok := debug.ReadBuildInfo(); ok {
		if v == "dev" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			v = info.Main.Version
		}
		var modified bool
		var revision, buildTime string
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.time":
				buildTime = setting.Value
			case "vcs.modified":
				modified = setting.Value == "true"
			}
		}
		if c == "none" && revision != "" {
			c = revision
			if modified {
				c += "-dirty"
			}
		}
		if d == "unknown" && buildTime != "" {
			d = buildTime
		}
	}

	displayCommit := c
	if len(displayCommit) > 7 && !strings.Contains(displayCommit, " ") {
		if strings.HasSuffix(displayCommit, "-dirty") && len(displayCommit) > 13 {
			displayCommit = displayCommit[:7] + "-dirty"
		} else if !strings.HasSuffix(displayCommit, "-dirty") {
			displayCommit = displayCommit[:7]
		}
	}

	return Info{
		Version:   strings.TrimPrefix(v, "v"),
		Commit:    displayCommit,
		Date:      d,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// Short returns just the clean version string without 'v' prefix.
func (i Info) Short() string {
	return i.Version
}

// String returns a single-line summary of version and build details.
func (i Info) String() string {
	return fmt.Sprintf("ilo-fans-agent-pve %s (commit: %s, built: %s, go: %s, platform: %s)",
		i.Version, i.Commit, i.Date, i.GoVersion, i.Platform)
}
