package system

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type RAIDArray struct {
	Name       string       `json:"name"`
	Level      string       `json:"level"`
	Devices    []RAIDDevice `json:"devices"`
	Status     string       `json:"status"`
	Size       uint64       `json:"size"`
	Used       uint64       `json:"used"`
	UUID       string       `json:"uuid"`
	MountPoint string       `json:"mountPoint"`
}

type RAIDDevice struct {
	Device string `json:"device"`
	Size   uint64 `json:"size"`
	Status string `json:"status"`
	Role   string `json:"role"`
}

// GetRAIDArrays parses /proc/mdstat and mdadm --detail
func GetRAIDArrays() ([]RAIDArray, error) {
	// 1. Read /proc/mdstat to find active arrays
	file, err := os.Open("/proc/mdstat")
	if err != nil {
		// If mdstat doesn't exist, maybe RAID is not supported/configured
		return []RAIDArray{}, nil
	}
	defer file.Close()

	arrays := []RAIDArray{}
	scanner := bufio.NewScanner(file)
	// Regex to find array names: mdX : active raid1 ...
	re := regexp.MustCompile(`^(md\d+)\s+:\s+active\s+raid(\d+)`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 3 {
			name := matches[1]
			
			// 2. Get more details using mdadm --detail
			detail, err := GetRAIDDetail("/dev/" + name)
			if err == nil {
				arrays = append(arrays, *detail)
			}
		}
	}

	return arrays, nil
}

func GetRAIDDetail(device string) (*RAIDArray, error) {
	cmd := exec.Command("mdadm", "--detail", "--export", device)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("mdadm --detail failed: %w", err)
	}

	array := &RAIDArray{
		Name:    device,
		Devices: []RAIDDevice{},
		Status:  "active",
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" { continue }
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 { continue }
		key := parts[0]
		val := parts[1]

		switch key {
		case "MD_LEVEL":
			array.Level = strings.TrimPrefix(val, "raid")
		case "MD_ARRAY_SIZE":
			// mdadm reports in KB
			size, _ := strconv.ParseUint(val, 10, 64)
			array.Size = size * 1024 
		case "MD_UUID":
			array.UUID = val
		case "MD_STATE":
			if strings.Contains(val, "degraded") {
				array.Status = "degraded"
			} else if strings.Contains(val, "failed") {
				array.Status = "failed"
			} else if strings.Contains(val, "rebuilding") || strings.Contains(val, "resync") {
				array.Status = "rebuilding"
			}
		}
	}

	// Get member devices
	cmdDetailed := exec.Command("mdadm", "--detail", device)
	outputDetailed, _ := cmdDetailed.Output()
	// Parse the table at the end of mdadm --detail output
	// Number   Major   Minor   RaidDevice State
	//   0       8       33        0      active sync   /dev/sdc1
	deviceLineRe := regexp.MustCompile(`\s+\d+\s+\d+\s+\d+\s+\d+\s+([\w\s]+)\s+(/dev/\S+)`)
	linesDetailed := strings.Split(string(outputDetailed), "\n")
	for _, line := range linesDetailed {
		m := deviceLineRe.FindStringSubmatch(line)
		if len(m) >= 3 {
			state := strings.TrimSpace(m[1])
			devName := m[2]
			array.Devices = append(array.Devices, RAIDDevice{
				Device: devName,
				Status: state,
				Role:   "data", // Simplified
			})
		}
	}

	return array, nil
}

// CreateRAID creates a new RAID array
func CreateRAID(name string, level string, devices []string) error {
	// mdadm --create /dev/md0 --level=1 --raid-devices=2 /dev/sda1 /dev/sdb1
	// level might be "1" or "raid1"
	if !strings.HasPrefix(level, "raid") {
		level = "raid" + level
	}

	devPath := name
	if !strings.HasPrefix(name, "/dev/") {
		devPath = "/dev/" + name
	}

	args := []string{"--create", devPath, "--level=" + level, fmt.Sprintf("--raid-devices=%d", len(devices))}
	args = append(args, devices...)
	
	// We use --run to avoid prompt if devices were previously used
	args = append(args, "--run")

	cmd := exec.Command("mdadm", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create RAID: %v, output: %s", err, string(output))
	}
	return nil
}

// DeleteRAID stops a RAID array
func DeleteRAID(name string) error {
	devPath := name
	if !strings.HasPrefix(name, "/dev/") {
		devPath = "/dev/" + name
	}
	
	cmdStop := exec.Command("mdadm", "--stop", devPath)
	if output, err := cmdStop.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop RAID: %v, output: %s", err, string(output))
	}

	return nil
}

// AddDiskToRAID adds a disk to an existing RAID array
func AddDiskToRAID(name, device string) error {
	devPath := name
	if !strings.HasPrefix(name, "/dev/") {
		devPath = "/dev/" + name
	}

	cmd := exec.Command("mdadm", "--manage", devPath, "--add", device)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add disk to RAID: %v, output: %s", err, string(output))
	}
	return nil
}

// RemoveDiskFromRAID removes a disk from an existing RAID array
func RemoveDiskFromRAID(name, device string) error {
	devPath := name
	if !strings.HasPrefix(name, "/dev/") {
		devPath = "/dev/" + name
	}

	// Usually need to fail it first before removing
	exec.Command("mdadm", "--manage", devPath, "--fail", device).Run()

	cmd := exec.Command("mdadm", "--manage", devPath, "--remove", device)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove disk from RAID: %v, output: %s", err, string(output))
	}
	return nil
}
