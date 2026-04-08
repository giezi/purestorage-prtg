package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/purestorage-prtg/purestorage-sensor/internal/params"
	"github.com/purestorage-prtg/purestorage-sensor/internal/pureapi"
	"github.com/purestorage-prtg/purestorage-sensor/internal/sensor"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version":
			fmt.Println(version)
			os.Exit(0)
		case "--help", "-h":
			printUsage()
			os.Exit(0)
		}
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		sensor.Fatal("failed to read stdin: %v", err)
	}

	if len(strings.TrimSpace(string(input))) == 0 {
		printUsage()
		os.Exit(0)
	}

	p, err := params.Parse(string(input))
	if err != nil {
		sensor.Fatal("parameter error: %v", err)
	}

	client := pureapi.NewClient(p.Endpoint, p.APIToken, p.APIVersion, p.Insecure)

	if err := client.Login(); err != nil {
		sensor.Fatal("FlashArray login failed: %v", err)
	}
	defer client.Logout()

	var result *sensor.Result

	switch p.Scope {
	case "capacity":
		result = sensor.RunCapacity(client, p.Warning, p.Critical)
	case "performance":
		result = sensor.RunPerformance(client)
	case "hardware":
		result = sensor.RunHardware(client)
	case "volumes":
		result = sensor.RunVolumes(client, p.Volumes)
	default:
		sensor.Fatal("unknown scope: %s", p.Scope)
	}

	result.Print()
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `purestorage-sensor %s
Pure Storage FlashArray PRTG Sensor

Monitors Pure Storage FlashArrays via REST API 2.x.
Designed for the PRTG Script v2 Sensor (parameters via stdin).

USAGE:
  echo "<parameters>" | purestorage-sensor

REQUIRED PARAMETERS:
  --endpoint <ip-or-fqdn>    FlashArray management IP or FQDN
  --apitoken <token>          API token for authentication
  --scope <scope>             Monitoring scope (see below)

SCOPES:
  capacity      Array space: total, used, used %%, data reduction, snapshots
  performance   IOPS, bandwidth, latency, queue depth
  hardware      Health summary: controllers, drives, PSUs, fans
  volumes       Per-volume space (requires --volumes)

OPTIONAL PARAMETERS:
  --warning <pct>             Capacity warning threshold (default: 80)
  --critical <pct>            Capacity error threshold (default: 90)
  --volumes <v1,v2,...>       Comma-separated volume names (max 10)
  --secure                    Enable TLS certificate verification (default: insecure)
  --insecure                  Skip TLS certificate verification (default)
  --auth-mode <mode>          Authentication mode (default: session)
  --apiversion <version>      API version, e.g. 2.26 (default: auto-negotiate)

OTHER:
  --help, -h                  Show this help
  --version                   Show version

EXAMPLES:
  echo "--endpoint 10.0.0.1 --apitoken abc123 --scope capacity" | purestorage-sensor
  echo "--endpoint fa.local --apitoken abc123 --scope performance" | purestorage-sensor
  echo "--endpoint 10.0.0.1 --apitoken abc123 --scope hardware" | purestorage-sensor
  echo "--endpoint 10.0.0.1 --apitoken abc123 --scope volumes --volumes vol1,vol2" | purestorage-sensor

PRTG CONFIGURATION:
  In PRTG, add a Script v2 sensor and set the Parameters field to:
    --endpoint %%host --apitoken %%scriptplaceholder1 --scope capacity
  PRTG passes this string to the binary via stdin.

`, version)
}
