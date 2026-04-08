package sensor

import (
	"fmt"
	"time"

	"github.com/purestorage-prtg/purestorage-sensor/internal/pureapi"
)

// RunSnapshots checks for stale volume snapshots using server-side filtering.
// Two lightweight API calls are made instead of fetching all snapshots:
//  1. Total snapshot count (limit=1 + total_item_count)
//  2. Old snapshots filtered by created < threshold (limit=1 + total_item_count + total)
func RunSnapshots(client *pureapi.Client, snapAgeDays int) *Result {
	totalCount, err := client.GetVolumeSnapshotCount()
	if err != nil {
		return NewErrorResult(fmt.Sprintf("snapshot count query failed: %v", err))
	}

	thresholdMs := time.Now().Add(-time.Duration(snapAgeDays) * 24 * time.Hour).UnixMilli()

	oldResp, err := client.GetOldVolumeSnapshots(thresholdMs)
	if err != nil {
		return NewErrorResult(fmt.Sprintf("old snapshots query failed: %v", err))
	}

	oldCount := oldResp.TotalItemCount

	var oldestAgeDays float64
	var oldestName string
	if len(oldResp.Items) > 0 {
		oldest := oldResp.Items[0]
		oldestName = oldest.Name
		created := time.UnixMilli(oldest.Created)
		oldestAgeDays = time.Since(created).Hours() / 24
	}

	var oldSpace int64
	if len(oldResp.Total) > 0 {
		oldSpace = oldResp.Total[0].Space.TotalPhysical
	}

	status := "ok"
	message := fmt.Sprintf("All %d snapshots within %d-day threshold", totalCount, snapAgeDays)

	if oldCount > 0 {
		status = "warning"
		message = fmt.Sprintf("%d of %d snapshots older than %d days (oldest: %s, %.0f days)",
			oldCount, totalCount, snapAgeDays, oldestName, oldestAgeDays)
	}

	if totalCount == 0 {
		message = "No snapshots found"
	}

	r := &Result{
		Version: 3,
		Status:  status,
		Message: message,
	}

	r.AddIntChannel(50, "Total Snapshots", "count", totalCount)

	warnUpper := 0.0
	r.Channels = append(r.Channels, Channel{
		ID:    51,
		Name:  "Old Snapshots",
		Type:  "integer",
		Kind:  "count",
		Value: oldCount,
		Limits: &Limits{
			Warning: &Threshold{Upper: &warnUpper},
		},
	})

	r.AddCustomFloatChannel(52, "Oldest Snapshot Age", "days", oldestAgeDays)
	r.AddIntChannel(53, "Old Snapshot Space", "size_bytes_disk", oldSpace)
	r.AddIntChannel(54, "Snapshot Threshold", "count", int64(snapAgeDays))

	return r
}
