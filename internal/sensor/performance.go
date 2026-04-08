package sensor

import (
	"fmt"

	"github.com/purestorage-prtg/purestorage-sensor/internal/pureapi"
)

// RunPerformance queries array performance data and builds the PRTG result.
func RunPerformance(client *pureapi.Client) *Result {
	resp, err := client.GetArrayPerformance()
	if err != nil {
		return NewErrorResult(fmt.Sprintf("performance query failed: %v", err))
	}
	if len(resp.Items) == 0 {
		return NewErrorResult("performance query returned no data")
	}

	item := resp.Items[0]
	totalIOPS := item.ReadsPerSec + item.WritesPerSec + item.OthersPerSec

	r := NewResult(fmt.Sprintf("FlashArray %s: %d IOPS, Read Latency %.0f us",
		item.Name, totalIOPS, item.UsecPerReadOp))

	r.AddIntChannel(20, "Read IOPS", "count", item.ReadsPerSec)
	r.AddIntChannel(21, "Write IOPS", "count", item.WritesPerSec)
	r.AddIntChannel(22, "Read Bandwidth", "size_bytes-per-second_network", item.ReadBytesPerSec)
	r.AddIntChannel(23, "Write Bandwidth", "size_bytes-per-second_network", item.WriteBytesPerSec)
	r.AddCustomFloatChannel(24, "Read Latency", "us", item.UsecPerReadOp)
	r.AddCustomFloatChannel(25, "Write Latency", "us", item.UsecPerWriteOp)
	r.AddIntChannel(26, "Avg I/O Size", "size_bytes_disk", item.BytesPerOp)
	r.AddCustomFloatChannel(27, "Queue Latency Read", "us", item.QueueUsecPerReadOp)
	r.AddCustomFloatChannel(28, "Queue Latency Write", "us", item.QueueUsecPerWriteOp)

	return r
}
