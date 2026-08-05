#!/usr/bin/env bash
set -euo pipefail

PREFIX="${PREFIX:-/opt/ilo-fans-agent-pve}"
ETC="/etc/ilo-fans-agent-pve"
VAR="/var/lib/ilo-fans-agent-pve"
BIN_SRC="${1:-./ilo-fans-agent-pve}"

if [[ ! -f "$BIN_SRC" ]]; then
  echo "usage: $0 [path/to/ilo-fans-agent-pve-binary]" >&2
  exit 1
fi

install -d -m 0755 "$PREFIX/bin" "$ETC" "$VAR"
install -m 0755 "$BIN_SRC" "$PREFIX/bin/ilo-fans-agent-pve"

if [[ ! -f "$ETC/config.yaml" ]]; then
  install -m 0644 config.yaml.example "$ETC/config.yaml"
fi

RUN_USER="root"
RUN_GROUP="root"
if [[ -f "$ETC/config.yaml" ]]; then
  RUN_USER=$(grep -E '^\s*user:' "$ETC/config.yaml" | head -1 | awk '{print $2}' || true)
  RUN_GROUP=$(grep -E '^\s*group:' "$ETC/config.yaml" | head -1 | awk '{print $2}' || true)
  RUN_USER="${RUN_USER:-root}"
  RUN_GROUP="${RUN_GROUP:-root}"
fi

UNIT="/etc/systemd/system/ilo-fans-agent-pve.service"
sed -e "s|%RUN_USER%|$RUN_USER|g" -e "s|%RUN_GROUP%|$RUN_GROUP|g" deploy/ilo-fans-agent-pve.service > "$UNIT.tmp"
if [[ "$RUN_USER" == "root" ]]; then
  sed -i '/^# User=/d; /^# Group=/d' "$UNIT.tmp" 2>/dev/null || sed -i '' '/^# User=/d; /^# Group=/d' "$UNIT.tmp"
else
  sed -i 's/^# User=/User=/; s/^# Group=/Group=/' "$UNIT.tmp" 2>/dev/null || sed -i '' 's/^# User=/User=/; s/^# Group=/Group=/' "$UNIT.tmp"
fi
mv "$UNIT.tmp" "$UNIT"

systemctl daemon-reload
echo "Installed. Next: $PREFIX/bin/ilo-fans-agent-pve -config $ETC/config.yaml token init"
echo "Then: systemctl enable --now ilo-fans-agent-pve"
