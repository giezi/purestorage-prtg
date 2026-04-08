package main

import (
	"fmt"
	"io"
	"os"

	"github.com/purestorage-prtg/purestorage-sensor/internal/params"
	"github.com/purestorage-prtg/purestorage-sensor/internal/pureapi"
	"github.com/purestorage-prtg/purestorage-sensor/internal/sensor"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(version)
		os.Exit(0)
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		sensor.Fatal("failed to read stdin: %v", err)
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
