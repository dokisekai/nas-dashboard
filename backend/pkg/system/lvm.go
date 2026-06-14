package system

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type PhysicalVolume struct {
	Device string `json:"device"`
	VGName string `json:"vgName"`
	Size   uint64 `json:"size"`
	Free   uint64 `json:"free"`
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

type VolumeGroup struct {
	Name    string `json:"name"`
	Size    uint64 `json:"size"`
	Free    uint64 `json:"free"`
	PVCount int    `json:"pvCount"`
	LVCount int    `json:"lvCount"`
	UUID    string `json:"uuid"`
}

type LogicalVolume struct {
	Name       string `json:"name"`
	VGName     string `json:"vgName"`
	Size       uint64 `json:"size"`
	Path       string `json:"path"`
	UUID       string `json:"uuid"`
	MountPoint string `json:"mountPoint"`
	Status     string `json:"status"`
}

// GetPhysicalVolumes returns list of PVs
func GetPhysicalVolumes() ([]PhysicalVolume, error) {
	cmd := exec.Command("pvs", "--reportformat", "json", "--units", "b", "--nosuffix")
	output, err := cmd.Output()
	if err != nil {
		// If LVM is not installed or no PVs, it might return error, but let's be graceful
		return []PhysicalVolume{}, nil
	}

	var data struct {
		Report []struct {
			PV []struct {
				PVName string `json:"pv_name"`
				VGName string `json:"vg_name"`
				PVSize string `json:"pv_size"`
				PVFree string `json:"pv_free"`
				PVUUID string `json:"pv_uuid"`
				PVAttr string `json:"pv_attr"`
			} `json:"pv"`
		} `json:"report"`
	}

	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("failed to parse pvs output: %w", err)
	}

	pvs := []PhysicalVolume{}
	if len(data.Report) > 0 {
		for _, pv := range data.Report[0].PV {
			var size, free uint64
			fmt.Sscanf(pv.PVSize, "%d", &size)
			fmt.Sscanf(pv.PVFree, "%d", &free)
			pvs = append(pvs, PhysicalVolume{
				Device: pv.PVName,
				VGName: pv.VGName,
				Size:   size,
				Free:   free,
				UUID:   pv.PVUUID,
				Status: pv.PVAttr,
			})
		}
	}
	return pvs, nil
}

// GetVolumeGroups returns list of VGs
func GetVolumeGroups() ([]VolumeGroup, error) {
	cmd := exec.Command("vgs", "--reportformat", "json", "--units", "b", "--nosuffix")
	output, err := cmd.Output()
	if err != nil {
		return []VolumeGroup{}, nil
	}

	var data struct {
		Report []struct {
			VG []struct {
				VGName  string `json:"vg_name"`
				VGSize  string `json:"vg_size"`
				VGFree  string `json:"vg_free"`
				PVCount string `json:"pv_count"`
				LVCount string `json:"lv_count"`
				VGUUID  string `json:"vg_uuid"`
			} `json:"vg"`
		} `json:"report"`
	}

	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("failed to parse vgs output: %w", err)
	}

	vgs := []VolumeGroup{}
	if len(data.Report) > 0 {
		for _, vg := range data.Report[0].VG {
			var size, free uint64
			var pvc, lvc int
			fmt.Sscanf(vg.VGSize, "%d", &size)
			fmt.Sscanf(vg.VGFree, "%d", &free)
			fmt.Sscanf(vg.PVCount, "%d", &pvc)
			fmt.Sscanf(vg.LVCount, "%d", &lvc)
			vgs = append(vgs, VolumeGroup{
				Name:    vg.VGName,
				Size:    size,
				Free:    free,
				PVCount: pvc,
				LVCount: lvc,
				UUID:    vg.VGUUID,
			})
		}
	}
	return vgs, nil
}

// GetLogicalVolumes returns list of LVs
func GetLogicalVolumes(vgName string) ([]LogicalVolume, error) {
	args := []string{"--reportformat", "json", "--units", "b", "--nosuffix"}
	if vgName != "" {
		args = append(args, vgName)
	}
	cmd := exec.Command("lvs", args...)
	output, err := cmd.Output()
	if err != nil {
		return []LogicalVolume{}, nil
	}

	var data struct {
		Report []struct {
			LV []struct {
				LVName string `json:"lv_name"`
				VGName string `json:"vg_name"`
				LVSize string `json:"lv_size"`
				LVPath string `json:"lv_path"`
				LVUUID string `json:"lv_uuid"`
				LVAttr string `json:"lv_attr"`
			} `json:"lv"`
		} `json:"report"`
	}

	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("failed to parse lvs output: %w", err)
	}

	lvs := []LogicalVolume{}
	if len(data.Report) > 0 {
		for _, lv := range data.Report[0].LV {
			var size uint64
			fmt.Sscanf(lv.LVSize, "%d", &size)
			lvs = append(lvs, LogicalVolume{
				Name:   lv.LVName,
				VGName: lv.VGName,
				Size:   size,
				Path:   lv.LVPath,
				UUID:   lv.LVUUID,
				Status: lv.LVAttr,
			})
		}
	}
	return lvs, nil
}

// CreatePhysicalVolume creates a PV
func CreatePhysicalVolume(device string) error {
	cmd := exec.Command("pvcreate", "-f", device)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("pvcreate failed: %v, output: %s", err, string(output))
	}
	return nil
}

// CreateVolumeGroup creates a VG
func CreateVolumeGroup(name string, devices []string) error {
	args := append([]string{name}, devices...)
	cmd := exec.Command("vgcreate", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("vgcreate failed: %v, output: %s", err, string(output))
	}
	return nil
}

// DeleteVolumeGroup deletes a VG
func DeleteVolumeGroup(name string) error {
	cmd := exec.Command("vgremove", "-f", name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("vgremove failed: %v, output: %s", err, string(output))
	}
	return nil
}

// CreateLogicalVolume creates an LV
func CreateLogicalVolume(vgName, name string, size uint64) error {
	// size is in bytes, convert to MB for lvcreate -L
	sizeMB := size / (1024 * 1024)
	if sizeMB == 0 && size > 0 {
		sizeMB = 1 // Minimum 1MB
	}
	cmd := exec.Command("lvcreate", "-L", fmt.Sprintf("%dM", sizeMB), "-n", name, vgName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("lvcreate failed: %v, output: %s", err, string(output))
	}
	return nil
}

// DeleteLogicalVolume deletes an LV
func DeleteLogicalVolume(vgName, name string) error {
	cmd := exec.Command("lvremove", "-f", vgName+"/"+name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("lvremove failed: %v, output: %s", err, string(output))
	}
	return nil
}
