package storage

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Define structs matching the JSON schema from `lsblk -J`
type LSBLKOutput struct {
	Blockdevices []BlockDevice `json:"blockdevices"`
}

// todo maybe add id field
type BlockDevice struct {
	Name       string        `json:"name"`
	Kname      string        `json:"kname"`
	Type       string        `json:"type"`
	Size       uint64        `json:"size"`
	FSAvail    uint64        `json:"fsavail,omitempty"`
	Mountpoint string        `json:"mountpoint"`
	Model      string        `json:"model"`
	Children   []BlockDevice `json:"children,omitempty"`
}

// GetDrives runs lsblk and parses the JSON output
func GetDrives() ([]BlockDevice, error) {
	cmd := exec.Command("lsblk", "--json", "--bytes", "-o", "NAME,KNAME,TYPE,SIZE,MOUNTPOINT,MODEL,FSAVAIL")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run lsblk: %w", err)
	}

	var data LSBLKOutput
	if err := json.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("failed to parse lsblk output: %w", err)
	}

	return data.Blockdevices, nil
}

func FilterDrives(devices []BlockDevice, names ...string) []BlockDevice {
	var result []BlockDevice
	if len(names) == 0 {
		return devices
	}
	for _, d := range devices {
		for _, name := range names {
			isDiskAndMatches := d.Type == "disk" && !strings.HasPrefix(d.Name, name)
			if isDiskAndMatches && d.FSAvail != 0 {
				result = append(result, d)
			}
		}
	}
	return result
}
