package sensor

import (
	"fmt"

	"github.com/purestorage-prtg/purestorage-sensor/internal/pureapi"
)

// RunCapacity queries array space data and builds the PRTG result.
func RunCapacity(client *pureapi.Client, warning, critical float64) *Result {
	resp, err := client.GetArraySpace()
	if err != nil {
		return NewErrorResult(fmt.Sprintf("capacity query failed: %v", err))
	}
	if len(resp.Items) == 0 {
		return NewErrorResult("capacity query returned no data")
	}

	item := resp.Items[0]
	usedPct := 0.0
	if item.Capacity > 0 {
		usedPct = float64(item.Space.TotalPhysical) / float64(item.Capacity) * 100
	}

	r := NewResult(fmt.Sprintf("FlashArray %s: Capacity %.1f%% used", item.Name, usedPct))

	r.AddIntChannel(10, "Total Capacity", "size_bytes_disk", item.Capacity)
	r.AddIntChannel(11, "Used Space", "size_bytes_disk", item.Space.TotalPhysical)
	r.AddFloatChannelWithLimits(12, "Used Percentage", "percent", usedPct, warning, critical)
	r.AddCustomFloatChannel(13, "Data Reduction Ratio", "x:1", item.Space.DataReduction)
	r.AddCustomFloatChannel(14, "Total Reduction Ratio", "x:1", item.Space.TotalReduction)
	r.AddIntChannel(15, "Snapshot Space", "size_bytes_disk", item.Space.Snapshots)
	r.AddIntChannel(16, "Shared Space", "size_bytes_disk", item.Space.Shared)
	r.AddIntChannel(17, "System Space", "size_bytes_disk", item.Space.System)
	r.AddIntChannel(18, "Volume Space", "size_bytes_disk", item.Space.Unique)
	r.AddIntChannel(19, "Total Provisioned", "size_bytes_disk", item.Space.TotalProvisioned)

	return r
}
