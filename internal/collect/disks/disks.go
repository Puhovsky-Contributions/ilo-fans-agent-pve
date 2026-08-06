package disks

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/collect/executil"
)

type DiskReading struct {
	Devpath string `json:"devpath"`
	Label   string `json:"label"`
	Temp    int    `json:"temp"`
	Model   string `json:"model"`
	Serial  string `json:"serial,omitempty"`
	WWN     string `json:"wwn,omitempty"`
}

type CollectOpts struct {
	SmartctlPath    string
	SmartctlUseSudo bool
}

var (
	reTemp1   = regexp.MustCompile(`(?i)Temperature:\s+(\d+)\s+Celsius`)
	reTemp2   = regexp.MustCompile(`(?i)Current Drive Temperature:\s*(\d+)\s*C`)
	reTemp3   = regexp.MustCompile(`(?i)Temperature Sensor \d+:\s+(\d+)\s+Celsius`)
	rePart    = regexp.MustCompile(`^(sd[a-z]+[0-9]+|nvme[0-9]+n[0-9]+p[0-9]+)$`)
	reSerial  = regexp.MustCompile(`(?i)^Serial [Nn]umber:\s*(.+)$`)
	reWWN     = regexp.MustCompile(`(?i)^(?:LU )?WWN(?: Device Id)?:\s*(.+)$`)
)

func Collect(opts CollectOpts) []DiskReading {
	var result []DiskReading
	entries, err := os.ReadDir("/sys/block")
	if err != nil {
		return result
	}
	for _, e := range entries {
		name := e.Name()
		if rePart.MatchString(name) {
			continue
		}
		if !strings.HasPrefix(name, "sd") && !strings.HasPrefix(name, "nvme") {
			continue
		}
		devpath := "/dev/" + name
		out, err := runSmartctl(opts, devpath)
		if err != nil || out == "" {
			continue
		}
		temp := parseSmartTemp(out)
		if temp == nil {
			continue
		}
		model := readModel(name)
		if model == "" {
			model = "unknown"
		}
		serial, wwn := parseSmartIdentity(out)
		short := filepath.Base(devpath)
		id := diskIdentity(serial, wwn, model)
		result = append(result, DiskReading{
			Devpath: devpath,
			Label:   short + " (" + id + ")",
			Temp:    *temp,
			Model:   model,
			Serial:  serial,
			WWN:     wwn,
		})
	}
	return result
}

func diskIdentity(serial, wwn, model string) string {
	s := strings.TrimSpace(serial)
	if s != "" && !strings.EqualFold(s, "unknown") {
		return s
	}
	w := strings.TrimSpace(wwn)
	if w != "" && !strings.EqualFold(w, "unknown") {
		return normalizeWWN(w)
	}
	if len(model) > 28 {
		return model[:12] + "…" + model[len(model)-8:]
	}
	if model != "" {
		return model
	}
	return "unknown"
}

func normalizeWWN(w string) string {
	w = strings.TrimSpace(w)
	w = strings.ReplaceAll(w, " ", "")
	if strings.HasPrefix(strings.ToLower(w), "0x") {
		return strings.ToLower(w)
	}
	return w
}

func parseSmartIdentity(text string) (serial, wwn string) {
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if m := reSerial.FindStringSubmatch(line); m != nil {
			serial = strings.TrimSpace(m[1])
			continue
		}
		if m := reWWN.FindStringSubmatch(line); m != nil {
			wwn = normalizeWWN(strings.TrimSpace(m[1]))
		}
	}
	return serial, wwn
}

func runSmartctl(opts CollectOpts, devpath string) (string, error) {
	args := []string{"-a", devpath}
	if opts.SmartctlUseSudo {
		return executil.Run("sudo", append([]string{opts.SmartctlPath}, args...)...)
	}
	return executil.Run(opts.SmartctlPath, args...)
}

func readModel(blockName string) string {
	data, err := os.ReadFile(filepath.Join("/sys/block", blockName, "device", "model"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func parseSmartTemp(text string) *int {
	if m := reTemp1.FindStringSubmatch(text); m != nil {
		v, _ := strconv.Atoi(m[1])
		return &v
	}
	if m := reTemp2.FindStringSubmatch(text); m != nil {
		v, _ := strconv.Atoi(m[1])
		return &v
	}
	if m := reTemp3.FindAllStringSubmatch(text, -1); len(m) > 0 {
		max := 0
		for _, match := range m {
			v, _ := strconv.Atoi(match[1])
			if v > max {
				max = v
			}
		}
		return &max
	}
	return nil
}
