package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/auth"
	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/collect/cpu"
	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/collect/disks"
	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/config"
	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/version"
)

type Server struct {
	cfg   config.Config
	auth  *auth.Store
	mux   *http.ServeMux
}

type thermalResponse struct {
	CPU    []cpu.Reading       `json:"cpu"`
	Disks  []disks.DiskReading `json:"disks"`
	Meta   thermalMeta         `json:"meta"`
}

type thermalMeta struct {
	CPU   cpu.Meta `json:"cpu"`
	At    string   `json:"collectedAt"`
}

func New(cfg config.Config, store *auth.Store) *Server {
	s := &Server{cfg: cfg, auth: store, mux: http.NewServeMux()}
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /v1/thermal", s.handleThermal)
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"version": version.Get().Short(),
	})
}

func (s *Server) handleThermal(w http.ResponseWriter, r *http.Request) {
	if !s.auth.Valid(r.Header.Get("Authorization")) {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	cpuResult := cpu.Collect(s.cfg.Collectors.SensorsPath)
	diskList := disks.Collect(disks.CollectOpts{
		SmartctlPath:    s.cfg.Collectors.SmartctlPath,
		SmartctlUseSudo: s.cfg.Collectors.SmartctlUseSudo,
	})
	resp := thermalResponse{
		CPU:   cpuResult.Readings,
		Disks: diskList,
		Meta: thermalMeta{
			CPU: cpuResult.Meta,
			At:  time.Now().UTC().Format(time.RFC3339),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
