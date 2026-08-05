package cpu

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Reading struct {
	Name string `json:"name"`
	Temp int    `json:"temp"`
}

type Meta struct {
	Attempted      bool   `json:"attempted"`
	OK             bool   `json:"ok"`
	Method         string `json:"method"`
	ReadingsCount  int    `json:"readingsCount"`
	Error          *string `json:"error"`
}

type Result struct {
	Readings []Reading `json:"readings"`
	Meta     Meta      `json:"meta"`
}

var (
	reCoreJSON = regexp.MustCompile(`(?i)^(Package id \d+|Core \d+)$`)
	reCoreText = regexp.MustCompile(`(?i)^(Package id \d+|Core \d+):\s+\+?([\d.]+)\s*(?:°| )?\s*C`)
)

func Collect(sensorsPath string) Result {
	meta := Meta{
		Attempted: true,
		Method:    "local",
	}
	raw, err := exec.Command(sensorsPath, "-j").Output()
	if err != nil {
		raw, err = exec.Command(sensorsPath).Output()
	}
	if err != nil {
		e := "sensors_failed"
		meta.Error = &e
		return Result{Meta: meta}
	}
	trimmed := strings.TrimLeft(string(raw), " \t\r\n")
	var readings []Reading
	if len(trimmed) > 0 && trimmed[0] == '{' {
		readings = parseJSON(raw)
	}
	if len(readings) == 0 {
		readings = parseText(string(raw))
	}
	meta.ReadingsCount = len(readings)
	if len(readings) == 0 {
		e := "sensors_output_unparsed"
		meta.Error = &e
	} else {
		meta.OK = true
	}
	return Result{Readings: readings, Meta: meta}
}

func parseJSON(raw []byte) []Reading {
	var data map[string]json.RawMessage
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil
	}
	var out []Reading
	for _, chipRaw := range data {
		var chip map[string]json.RawMessage
		if err := json.Unmarshal(chipRaw, &chip); err != nil {
			continue
		}
		for label, propsRaw := range chip {
			if !reCoreJSON.MatchString(label) {
				continue
			}
			var props map[string]any
			if err := json.Unmarshal(propsRaw, &props); err != nil {
				continue
			}
			for key, val := range props {
				if !strings.HasPrefix(key, "temp") || !strings.HasSuffix(key, "_input") {
					continue
				}
				switch v := val.(type) {
				case float64:
					out = append(out, Reading{Name: label, Temp: int(v + 0.5)})
				case json.Number:
					if f, err := v.Float64(); err == nil {
						out = append(out, Reading{Name: label, Temp: int(f + 0.5)})
					}
				}
				break
			}
		}
	}
	return out
}

func parseText(text string) []Reading {
	var out []Reading
	for _, line := range strings.Split(text, "\n") {
		m := reCoreText.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		f, _ := strconv.ParseFloat(m[2], 64)
		out = append(out, Reading{Name: m[1], Temp: int(f + 0.5)})
	}
	return out
}
