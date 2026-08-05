# Install (tarball)

1. Extract `ilo-fans-agent-pve_<ver>_linux_amd64_build.tar.gz` on the PVE host.
2. Edit `config.yaml` — set `listen` to the **SDN IP** IFC can reach (plain HTTP, no TLS).
3. `sudo ./scripts/install.sh ./ilo-fans-agent-pve`
4. `sudo /opt/ilo-fans-agent-pve/bin/ilo-fans-agent-pve token init` — save token for IFC.
5. `sudo systemctl enable --now ilo-fans-agent-pve`
6. From IFC host: `curl -H "Authorization: Bearer $TOKEN" http://SDN_IP:9847/v1/thermal`

Dependencies on PVE: `lm-sensors`, `smartmontools`.

## Building .deb (optional)

Same paths as `install.sh`. Example with [nfpm](https://nfpm.goreleaser.com/): add an `nfpm.yaml` mapping `/opt`, `/etc`, systemd unit, `depends: smartmontools`.
