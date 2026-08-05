# ilo-fans-agent-pve

Small HTTP agent on each Proxmox node. Exposes host CPU (`sensors`) and disk SMART temps for [IFC](https://github.com/playtika/ifc) fan control — no Proxmox API token or SSH from IFC.

- `GET /health` — liveness
- `GET /v1/thermal` — Bearer auth, JSON `cpu` + `disks`

See [docs/install.md](docs/install.md), [docs/api-v1.md](docs/api-v1.md), [docs/run-as.md](docs/run-as.md).

## Build

```bash
go build -trimpath -o ilo-fans-agent-pve ./cmd/ilo-fans-agent-pve
```

Static Linux amd64 (release CI on tag `v*`):

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o ilo-fans-agent-pve ./cmd/ilo-fans-agent-pve
```

GitHub Release assets:

- `ilo-fans-agent-pve_<ver>_linux_amd64_build.tar.gz` — binary + `config.yaml.example`, `deploy/`, `scripts/` (install/upgrade; see [docs/install.md](docs/install.md) in repo)
- `SHA256SUMS` — checksum for the build archive

## Token

```bash
ilo-fans-agent-pve -config /etc/ilo-fans-agent-pve/config.yaml token init
ilo-fans-agent-pve token rotate
ilo-fans-agent-pve token show
```

Send **SIGHUP** to reload token from disk without restart.
