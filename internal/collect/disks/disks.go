package disks

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/playtika/ilo-fans-agent-pve/internal/collect/executil"
)

type DiskReading struct {
	Devpath string `json:"devpath"`
	Label   string `json:"label"`
	Temp    int    `json:"temp"`
	Model   string `json:"model"`
}

type CollectOpts struct {
	SmartctlPath    string
	SmartctlUseSudo bool
}

var (
	reTemp1 = regexp.MustCompile(`(?i)Temperature:\s+(\d+)\s+Celsius`)
	reTemp2 = regexp.MustCompile(`(?i)Current Drive Temperature:\s*(\d+)\s*C`)
	reTemp3 = regexp.MustCompile(`(?i)Temperature Sensor \d+:\s+(\d+)\s+Celsius`)
	rePart  = regexp.MustCompile(`^(sd[a-z]+[0-9]+|nvme[0-9]+n[0-9]+p[0-9]+)$`)
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
		short := filepath.Base(devpath)
		result = append(result, DiskReading{
			Devpath: devpath,
			Label:   short + " (" + model + ")",
			Temp:    *temp,
			Model:   model,
		})
	}
	return result
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
