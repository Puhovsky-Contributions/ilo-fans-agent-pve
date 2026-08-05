# Run-as user (smartctl / sensors)

Process user is set in **systemd** (`User=` / `Group=`) from `run_as` in `config.yaml` via `install.sh`.

Pick one:

## 1. Root (default v0.1)

No extra steps.

## 2. Dedicated user + disk group

```bash
useradd -r -s /usr/sbin/nologin ifc-pve-agent
usermod -aG disk ifc-pve-agent
```

Set in `config.yaml`:

```yaml
run_as:
  user: ifc-pve-agent
  group: ifc-pve-agent
```

Re-run `install.sh` or edit the unit, then `systemctl daemon-reload && systemctl restart ilo-fans-agent-pve`.

## 3. setcap on smartctl (site policy)

```bash
setcap cap_sys_rawio+ep /usr/sbin/smartctl
```

Use only if your security policy allows it.

## 4. sudoers (optional)

`/etc/sudoers.d/ilo-fans-agent-pve`:

```
ifc-pve-agent ALL=(root) NOPASSWD: /usr/sbin/smartctl
```

In `config.yaml`:

```yaml
collectors:
  smartctl_use_sudo: true
```

Sensors (`lm-sensors`) may still require root or appropriate groups on some hosts.
