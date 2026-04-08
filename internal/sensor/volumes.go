package sensor

import (
	"fmt"
	"strings"

	"github.com/purestorage-prtg/purestorage-sensor/internal/pureapi"
)

const maxVolumesPerSensor = 10

// RunVolumes queries space data for explicitly configured volumes.
func RunVolumes(client *pureapi.Client, volumes []string) *Result {
	if len(volumes) == 0 {
		return NewErrorResult("no volumes specified (use --volumes vol1,vol2,...)")
	}
	if len(volumes) > maxVolumesPerSensor {
		return NewErrorResult(fmt.Sprintf("too many volumes (%d), maximum is %d per sensor",
			len(volumes), maxVolumesPerSensor))
	}

	resp, err := client.GetVolumesSpace(volumes)
	if err != nil {
		return NewErrorResult(fmt.Sprintf("volumes query failed: %v", err))
	}
	if len(resp.Items) == 0 {
		return NewErrorResult(fmt.Sprintf("no data returned for volumes: %s",
			strings.Join(volumes, ", ")))
	}

	r := NewResult(fmt.Sprintf("Monitoring %d volumes", len(resp.Items)))

	baseID := 40
	for i, vol := range resp.Items {
		offset := i * 3
		r.AddIntChannel(baseID+offset, vol.Name+" Used", "size_bytes_disk", vol.Space.TotalPhysical)
		r.AddIntChannel(baseID+offset+1, vol.Name+" Provisioned", "size_bytes_disk", vol.Space.TotalProvisioned)
		r.AddCustomFloatChannel(baseID+offset+2, vol.Name+" Data Reduction", "x:1", vol.Space.DataReduction)
	}

	return r
}
