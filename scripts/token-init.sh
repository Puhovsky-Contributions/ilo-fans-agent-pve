#!/usr/bin/env bash
set -euo pipefail
CONFIG="${CONFIG:-/etc/ilo-fans-agent-pve/config.yaml}"
exec /opt/ilo-fans-agent-pve/bin/ilo-fans-agent-pve -config "$CONFIG" token init
