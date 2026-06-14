package system

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Partition matches the frontend Partition interface
type Partition struct {
	Device     string   `json:"device"`
	Size       uint64   `json:"size"`
	Start      uint64   `json:"start"`
	End        uint64   `json:"end"`
	Type       string   `json:"type"`
	Filesystem string   `json:"filesystem"`
	MountPoint string   `json:"mountPoint"`
	Flags      []string `json:"flags"`
}

// PartitionTable matches the frontend PartitionTable interface
type PartitionTable struct {
	Type       string      `json:"type"` // gpt, mbr, none
	Partitions []Partition `json:"partitions"`
	FreeSpace  uint64      `json:"freeSpace"`
}

// GetPartitionTable gets the partition table for a device
func GetPartitionTable(device string) (*PartitionTable, error) {
	devPath := device
	if !strings.HasPrefix(devPath, "/dev/") {
		devPath = "/dev/" + devPath
	}

	// Get disk type using parted
	cmdType := exec.Command("parted", "-s", devPath, "print")
	outputType, _ := cmdType.Output()
	tableType := "none"
	if strings.Contains(string(outputType), "Partition Table: gpt") {
		tableType = "gpt"
	} else if strings.Contains(string(outputType), "Partition Table: msdos") {
		tableType = "mbr"
	}

	// Get partitions using lsblk
	cmd := exec.Command("lsblk", "-J", "-b", "-o", "NAME,SIZE,START,FSTYPE,MOUNTPOINT,PARTTYPE,PARTFLAGS", devPath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("lsblk failed: %w", err)
	}

	var data struct {
		BlockDevices []struct {
			Name       string `json:"name"`
			Size       uint64 `json:"size"`
			Children   []struct {
				Name       string `json:"name"`
				Size       uint64 `json:"size"`
				Start      uint64 `json:"start"`
				FsType     string `json:"fstype"`
				MountPoint string `json:"mountpoint"`
				PartType   string `json:"parttype"`
				PartFlags  string `json:"partflags"`
			} `json:"children"`
		} `json:"blockdevices"`
	}

	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("failed to parse lsblk output: %w", err)
	}

	table := &PartitionTable{
		Type:       tableType,
		Partitions: []Partition{},
	}

	if len(data.BlockDevices) > 0 {
		deviceSize := data.BlockDevices[0].Size
		usedSpace := uint64(0)
		for _, child := range data.BlockDevices[0].Children {
			p := Partition{
				Device:     "/dev/" + child.Name,
				Size:       child.Size,
				Start:      child.Start,
				End:        child.Start + child.Size,
				Type:       child.PartType,
				Filesystem: child.FsType,
				MountPoint: child.MountPoint,
			}
			if child.PartFlags != "" {
				p.Flags = strings.Split(child.PartFlags, ",")
			}
			table.Partitions = append(table.Partitions, p)
			usedSpace += child.Size
		}
		if deviceSize > usedSpace {
			table.FreeSpace = deviceSize - usedSpace
		}
	}

	return table, nil
}

// CreatePartition creates a new partition
func CreatePartition(device string, start, end uint64, partType, fs string) error {
	devPath := device
	if !strings.HasPrefix(devPath, "/dev/") {
		devPath = "/dev/" + devPath
	}

	// Default to primary if not specified
	if partType == "" {
		partType = "primary"
	}

	// parted -s /dev/sdx mkpart primary ext4 startMB endMB
	// Converting bytes to MB for parted
	startMB := start / (1024 * 1024)
	endMB := end / (1024 * 1024)
	
	cmd := exec.Command("parted", "-s", devPath, "mkpart", partType, fs, fmt.Sprintf("%dMB", startMB), fmt.Sprintf("%dMB", endMB))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create partition: %v, output: %s", err, string(output))
	}
	return nil
}

// DeletePartition deletes a partition by number
func DeletePartition(device string, number int) error {
	devPath := device
	if !strings.HasPrefix(devPath, "/dev/") {
		devPath = "/dev/" + devPath
	}

	cmd := exec.Command("parted", "-s", devPath, "rm", fmt.Sprintf("%d", number))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete partition: %v, output: %s", err, string(output))
	}
	return nil
}
