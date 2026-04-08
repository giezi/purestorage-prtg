# Pure Storage FlashArray PRTG Sensor

A Go-based custom sensor for PRTG Network Monitor that monitors Pure Storage FlashArrays via REST API 2.x. Built as a single static binary for the **PRTG Script v2 Sensor** -- no Python, no dependencies.

## Features

| Scope | Description | Channels |
|-------|-------------|----------|
| `capacity` | Array capacity, used space, data reduction | 10 channels |
| `performance` | IOPS, bandwidth, latency, queue depth | 9 channels |
| `hardware` | Controller, shelf, drive, PSU, fan health summary | 8 channels |
| `volumes` | Per-volume used space, provisioned, data reduction | 3 per volume (max 10 volumes) |

## Requirements

- PRTG 25.x or newer with Script v2 Sensor support
- Pure Storage FlashArray with REST API 2.x
- API Token for FlashArray authentication

## Build

Requires Go 1.21+.

```bash
# Windows binary (for PRTG on Windows)
make build-windows

# Linux binary (for Multi-Platform Probe)
make build-linux

# Both
make build-all
```

The output is a single static binary (`purestorage-sensor.exe` / `purestorage-sensor`).

## Installation

### Windows PRTG Server / Remote Probe

1. Copy `purestorage-sensor.exe` to:
   ```
   C:\Program Files (x86)\PRTG Network Monitor\Custom Sensors\scripts\
   ```
2. Restart the PRTG Probe Service (only needed after first deployment).

### Linux Multi-Platform Probe

1. Copy `purestorage-sensor` to:
   ```
   /opt/paessler/share/scripts/
   ```
2. Make executable: `chmod +x /opt/paessler/share/scripts/purestorage-sensor`

## PRTG Sensor Configuration

### Prerequisites

1. Create or obtain an API Token on the FlashArray (System > Users > API Tokens).
2. In PRTG, add a Device for the FlashArray (IP/FQDN).
3. On the Device, configure **Credentials for Script Sensors**:
   - Set **Script Placeholder 1** to the API Token.

### Adding Sensors

Add a **Script v2** sensor for each scope. Select `purestorage-sensor.exe` as the script.

#### Capacity Sensor

**Parameters:**
```
--endpoint %host --apitoken %scriptplaceholder1 --scope capacity --warning 80 --critical 90
```

**Recommended interval:** 5 minutes

Monitors total capacity, used space, used percentage (with warning/error thresholds), data reduction ratios, snapshot/shared/system/volume space, and total provisioned.

#### Performance Sensor

**Parameters:**
```
--endpoint %host --apitoken %scriptplaceholder1 --scope performance
```

**Recommended interval:** 1 minute

Monitors read/write IOPS, read/write bandwidth, read/write latency (microseconds), average I/O size, and queue latency.

#### Hardware Summary Sensor

**Parameters:**
```
--endpoint %host --apitoken %scriptplaceholder1 --scope hardware
```

**Recommended interval:** 5 minutes

Reports critical component counts by type (controllers, shelves, drives, PSUs, fans, temperature). Sensor status is automatically set to error/warning based on hardware health.

#### Volume Space Sensor (optional)

**Parameters:**
```
--endpoint %host --apitoken %scriptplaceholder1 --scope volumes --volumes vol1,vol2,vol3
```

**Recommended interval:** 5 minutes

Monitors used space, provisioned space, and data reduction per volume. Maximum 10 volumes per sensor instance. Create multiple sensors for more volumes.

## Parameters Reference

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `--endpoint` | Yes | | FlashArray IP or FQDN |
| `--apitoken` | Yes | | API Token for authentication |
| `--scope` | Yes | | `capacity`, `performance`, `hardware`, or `volumes` |
| `--insecure` | No | `true` | Skip TLS certificate verification (default) |
| `--secure` | No | | Enable TLS certificate verification |
| `--warning` | No | `80` | Capacity warning threshold (%) |
| `--critical` | No | `90` | Capacity error threshold (%) |
| `--volumes` | No | | Comma-separated volume names (required for `volumes` scope) |
| `--auth-mode` | No | `session` | Authentication mode (`session`) |
| `--apiversion` | No | auto | Specific API version (e.g., `2.26`) |

## How It Works

1. PRTG starts the binary and writes the parameter string to stdin.
2. The binary authenticates to the FlashArray (API Token -> Session Token).
3. It queries the relevant API endpoint based on `--scope`.
4. It builds PRTG Script v2 JSON (Schema v3) with channels and writes it to stdout.
5. The session is logged out and the binary exits with code 0.

Errors are always reported via JSON output (`"status": "error"`) with exit code 0, per PRTG convention.

## Example Output

```json
{
  "version": 3,
  "status": "ok",
  "message": "FlashArray myarray: Capacity 37.2% used",
  "channels": [
    {
      "id": 10,
      "name": "Total Capacity",
      "type": "integer",
      "kind": "size_bytes_disk",
      "value": 29511391379456
    },
    {
      "id": 12,
      "name": "Used Percentage",
      "type": "float",
      "kind": "percent",
      "value": 37.2,
      "limits": {
        "warning": { "upper": 80 },
        "error": { "upper": 90 }
      }
    }
  ]
}
```

## TLS / Certificates

TLS certificate verification is **disabled by default** because FlashArrays typically use self-signed certificates. Use `--secure` to enable strict TLS verification against the system CA store.

## License

Internal use.
