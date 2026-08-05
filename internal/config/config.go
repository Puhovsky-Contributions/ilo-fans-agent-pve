package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type RunAs struct {
	User  string `yaml:"user"`
	Group string `yaml:"group"`
}

type Collectors struct {
	SmartctlPath     string `yaml:"smartctl_path"`
	SmartctlUseSudo  bool   `yaml:"smartctl_use_sudo"`
	SensorsPath      string `yaml:"sensors_path"`
}

type Config struct {
	Listen     string     `yaml:"listen"`
	TokenFile  string     `yaml:"token_file"`
	RunAs      RunAs      `yaml:"run_as"`
	Collectors Collectors `yaml:"collectors"`
}

func Default() Config {
	return Config{
		Listen:    "127.0.0.1:9847",
		TokenFile: "/var/lib/ilo-fans-agent-pve/token",
		RunAs: RunAs{
			User:  "root",
			Group: "root",
		},
		Collectors: Collectors{
			SmartctlPath:    "/usr/sbin/smartctl",
			SmartctlUseSudo: false,
			SensorsPath:     "/usr/bin/sensors",
		},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Listen == "" {
		cfg.Listen = Default().Listen
	}
	if cfg.TokenFile == "" {
		cfg.TokenFile = Default().TokenFile
	}
	if cfg.Collectors.SmartctlPath == "" {
		cfg.Collectors.SmartctlPath = Default().Collectors.SmartctlPath
	}
	if cfg.Collectors.SensorsPath == "" {
		cfg.Collectors.SensorsPath = Default().Collectors.SensorsPath
	}
	if cfg.RunAs.User == "" {
		cfg.RunAs.User = "root"
	}
	if cfg.RunAs.Group == "" {
		cfg.RunAs.Group = "root"
	}
	return cfg, nil
}
