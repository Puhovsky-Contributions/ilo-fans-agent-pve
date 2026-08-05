# ilo-fans-agent-pve + IFC agent transport

## Agent repo (`../ilo-fans-agent-pve`)
- [x] Go layout: cmd, internal/*, deploy, scripts, docs
- [x] GET /health, GET /v1/thermal + Bearer
- [x] CPU: sensors -j/text; disks: smartctl
- [x] token file + CLI init/rotate/show + SIGHUP
- [x] install.sh + systemd template
- [x] GHA release tarball linux-amd64

## IFC
- [x] proxmox-disks.php → HTTP agent (+ request cache)
- [x] .env.example, servers.json.example, docker-compose, Dockerfile
- [x] README Proxmox section
- [x] fan-daemon log meta source=agent

## Verify
- [x] go build ./...
- [x] php -l lib/proxmox-disks.php

## Review
Agent sibling at `/Users/andreypu/playtika/ilo-fans-agent-pve`. IFC keeps `pve:cpu` / `pve:disks` prefixes; transport = agent only.

### How to test
1. PVE: build/install agent, `token init`, curl `/v1/thermal`
2. IFC: set `PROXMOX_AGENT_*` or servers.json fields, UI Proxmox blocks + daemon JSON `proxmox.agentUrl`
