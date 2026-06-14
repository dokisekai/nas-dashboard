package system

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// SMARTInfo matches the frontend SMARTInfo interface
type SMARTInfo struct {
	OverallHealth string         `json:"overallHealth"`
	Attributes    []SMARTAttr    `json:"attributes"`
	LastTest      *SMARTTest     `json:"lastTest"`
	Temperature   int            `json:"temperature"`
	PowerOnHours  int            `json:"powerOnHours"`
	ErrorLog      []SMARTError   `json:"errorLog"`
}

type SMARTAttr struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Value     int    `json:"value"`
	Worst     int    `json:"worst"`
	Threshold int    `json:"threshold"`
	Status    string `json:"status"`
}

type SMARTTest struct {
	Type      string `json:"type"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	Remaining int    `json:"remaining"`
	Result    string `json:"result"`
}

type SMARTError struct {
	Critical bool     `json:"critical"`
	Count    int      `json:"count"`
	Details  []string `json:"details"`
}

// GetSMARTInfo gets SMART info for a device using smartctl
func GetSMARTInfo(device string) (*SMARTInfo, error) {
	// smartctl -j -a /dev/sdx
	// We use the full path for device if it doesn't start with /dev/
	devPath := device
	if !strings.HasPrefix(devPath, "/dev/") {
		devPath = "/dev/" + devPath
	}

	cmd := exec.Command("smartctl", "-j", "-a", devPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// smartctl returns non-zero if some bits are set in the exit status
		// but we might still have valid JSON. Exit code 0-2 are usually fine.
		if len(output) == 0 {
			return nil, fmt.Errorf("smartctl failed: %w", err)
		}
	}

	var data map[string]interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("failed to parse smartctl output: %w", err)
	}

	info := &SMARTInfo{
		OverallHealth: "good",
		Attributes:    []SMARTAttr{},
		ErrorLog:      []SMARTError{},
	}

	// Extract health
	if smartStatus, ok := data["smart_status"].(map[string]interface{}); ok {
		if passed, ok := smartStatus["passed"].(bool); ok && !passed {
			info.OverallHealth = "failed"
		}
	}

	// Extract Temperature
	if temperature, ok := data["temperature"].(map[string]interface{}); ok {
		if current, ok := temperature["current"].(float64); ok {
			info.Temperature = int(current)
		}
	}

	// Extract Power On Hours
	if powerOnTime, ok := data["power_on_time"].(map[string]interface{}); ok {
		if hours, ok := powerOnTime["hours"].(float64); ok {
			info.PowerOnHours = int(hours)
		}
	}

	// Extract Attributes
	if ataSmartAttr, ok := data["ata_smart_attributes"].(map[string]interface{}); ok {
		if table, ok := ataSmartAttr["table"].([]interface{}); ok {
			for _, item := range table {
				attrMap := item.(map[string]interface{})
				attr := SMARTAttr{
					ID:    int(attrMap["id"].(float64)),
					Name:  attrMap["name"].(string),
					Value: int(attrMap["value"].(float64)),
					Worst: int(attrMap["worst"].(float64)),
					Status: "ok",
				}
				if thresh, ok := attrMap["thresh"].(float64); ok {
					attr.Threshold = int(thresh)
				}
				
				// Basic status check based on threshold
				if attr.Threshold > 0 && attr.Value <= attr.Threshold {
					attr.Status = "warning"
					if info.OverallHealth == "good" {
						info.OverallHealth = "warning"
					}
				}
				info.Attributes = append(info.Attributes, attr)
			}
		}
	}

	// Extract Last Test
	if selfTest, ok := data["ata_smart_self_test"].(map[string]interface{}); ok {
		if statusMap, ok := selfTest["status"].(map[string]interface{}); ok {
			info.LastTest = &SMARTTest{
				Status: strings.ToLower(statusMap["string"].(string)),
				Result: statusMap["string"].(string),
			}
			// Estimate progress if possible
			if remaining, ok := statusMap["remaining_minutes"].(float64); ok {
				info.LastTest.Remaining = int(remaining)
			}
		}
	}

	return info, nil
}

// RunSMARTTest initiates a SMART self-test
func RunSMARTTest(device string, testType string) error {
	// testType: short, long, conveyance
	validTypes := map[string]bool{"short": true, "long": true, "conveyance": true}
	if !validTypes[testType] {
		return fmt.Errorf("invalid test type: %s", testType)
	}

	devPath := device
	if !strings.HasPrefix(devPath, "/dev/") {
		devPath = "/dev/" + devPath
	}

	cmd := exec.Command("smartctl", "-t", testType, devPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start SMART test: %v, output: %s", err, string(output))
	}

	return nil
}
