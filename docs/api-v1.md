# GET /v1/thermal

Bearer token (same value as `/var/lib/ilo-fans-agent-pve/token`).

```http
GET /v1/thermal HTTP/1.1
Host: 172.16.x.x:9847
Authorization: Bearer <token>
```

Response:

```json
{
  "cpu": [{"name": "Package id 0", "temp": 42}, {"name": "Core 0", "temp": 38}],
  "disks": [{"devpath": "/dev/sda", "label": "sda (MODEL)", "temp": 35, "model": "MODEL"}],
  "meta": {
    "cpu": {
      "attempted": true,
      "ok": true,
      "method": "local",
      "readingsCount": 2,
      "error": null
    },
    "collectedAt": "2026-08-06T10:00:00Z"
  }
}
```

# GET /health

No auth. `{"status":"ok"}`
