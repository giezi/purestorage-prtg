package sensor

import (
	"fmt"
	"strings"

	"github.com/purestorage-prtg/purestorage-sensor/internal/pureapi"
)

// RunHardware queries hardware status and builds a summary PRTG result.
func RunHardware(client *pureapi.Client) *Result {
	resp, err := client.GetHardware()
	if err != nil {
		return NewErrorResult(fmt.Sprintf("hardware query failed: %v", err))
	}
	if len(resp.Items) == 0 {
		return NewErrorResult("hardware query returned no data")
	}

	var (
		totalCritical    int64
		totalNonOK       int64
		ctCritical       int64
		shCritical       int64
		drvCritical      int64
		pwrCritical      int64
		fanCritical      int64
		tempAlerts       int64
		criticalNames    []string
	)

	for _, hw := range resp.Items {
		if isNotInstalled(hw.Status) {
			continue
		}

		if hw.Status == "critical" {
			totalCritical++
			criticalNames = append(criticalNames, hw.Name)

			switch normalizeType(hw.Type) {
			case "controller":
				ctCritical++
			case "shelf":
				shCritical++
			case "drive":
				drvCritical++
			case "psu":
				pwrCritical++
			case "fan":
				fanCritical++
			}
		}

		if hw.Status != "ok" {
			totalNonOK++
		}

		if normalizeType(hw.Type) == "temp" && hw.Status != "ok" {
			tempAlerts++
		}
		if hw.Details != "" && strings.Contains(strings.ToLower(hw.Details), "temperature") {
			tempAlerts++
		}
	}

	status := "ok"
	message := "All hardware components healthy"
	if totalCritical > 0 {
		status = "error"
		if len(criticalNames) > 5 {
			message = fmt.Sprintf("%d critical: %s, ...", totalCritical,
				strings.Join(criticalNames[:5], ", "))
		} else {
			message = fmt.Sprintf("%d critical: %s", totalCritical,
				strings.Join(criticalNames, ", "))
		}
	} else if totalNonOK > 0 {
		status = "warning"
		message = fmt.Sprintf("%d components not OK", totalNonOK)
	}

	r := &Result{
		Version: 3,
		Status:  status,
		Message: message,
	}

	r.AddIntChannel(30, "Critical Components", "count", totalCritical)
	r.AddIntChannel(31, "Non-OK Components", "count", totalNonOK)
	r.AddIntChannel(32, "Controllers Critical", "count", ctCritical)
	r.AddIntChannel(33, "Shelves Critical", "count", shCritical)
	r.AddIntChannel(34, "Drives Critical", "count", drvCritical)
	r.AddIntChannel(35, "Power Supplies Critical", "count", pwrCritical)
	r.AddIntChannel(36, "Fans Critical", "count", fanCritical)
	r.AddIntChannel(37, "Temperature Alerts", "count", tempAlerts)

	return r
}

// isNotInstalled handles both API formats: "not installed" (old) and "not_installed" (2.x).
func isNotInstalled(status string) bool {
	return status == "not installed" || status == "not_installed"
}

// normalizeType maps both short (old docs) and long (real API 2.x) type names to a canonical form.
func normalizeType(t string) string {
	switch t {
	case "ct", "controller":
		return "controller"
	case "ch", "chassis", "sh", "storage_shelf":
		return "shelf"
	case "bay", "drive_bay":
		return "drive"
	case "pwr", "power_supply":
		return "psu"
	case "fan", "cooling":
		return "fan"
	case "tmp", "temp_sensor":
		return "temp"
	case "eth", "eth_port":
		return "eth"
	case "fc", "fc_port":
		return "fc"
	case "sas", "sas_port":
		return "sas"
	case "nvb", "nvram_bay":
		return "nvram"
	case "iom", "io_module":
		return "iom"
	default:
		return t
	}
}
