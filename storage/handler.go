package storage

import (
	"goNAS/helper"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type DriveInfo struct {
	Name              string      `json:"name"`
	Path              string      `json:"path"`
	SizeSectors       uint64      `json:"size_sectors"`
	LogicalBlockSize  uint64      `json:"logical_block_size"`
	PhysicalBlockSize uint64      `json:"physical_block_size"`
	SizeBytes         uint64      `json:"size_bytes"`
	IsRotational      bool        `json:"is_rotational"`
	Model             string      `json:"model"`
	Vendor            string      `json:"vendor"`
	Type              string      `json:"type"`
	MountPoints       []MountInfo `json:"mountpoint"`
	FsType            string      `json:"fstype"`
	FsAvail           uint64      `json:"fsavail"`
}
type MountInfo struct {
	Device     string
	MountPoint string
	FsType     string
}
type DriveFilter struct {
	Names        []string
	IsRotational *bool
	MinSize      uint64
	MaxSize      uint64
	Mounted      *bool
	MountPrefix  string
	MinFsAvail   uint64
	MaxFsAvail   uint64
}

func Filter(f DriveFilter, d ...DriveInfo) []DriveInfo {
	// Precompute small things to avoid recomputing inside loop
	hasNames := len(f.Names) > 0
	checkMountPrefix := f.MountPrefix != ""

	// Preallocate result capacity for small optimization
	result := make([]DriveInfo, 0, len(d))

	for _, drive := range d {
		// --- Names ---
		if hasNames && !helper.Contains(f.Names, drive.Name) {
			continue
		}

		// --- Type ---
		if f.IsRotational != nil && *f.IsRotational && drive.IsRotational {
			if *f.IsRotational && !drive.IsRotational {
				continue
			}
		}

		// --- Size ---
		if (f.MinSize > 0 && drive.SizeBytes < f.MinSize) ||
			(f.MaxSize > 0 && drive.SizeBytes > f.MaxSize) ||
			(f.MinFsAvail > 0 && drive.FsAvail < f.MinFsAvail) ||
			(f.MaxFsAvail > 0 && drive.FsAvail > f.MaxFsAvail) {
			continue
		}

		// --- Mounted ---

		// --- Mounted filters ---
		isMounted := len(drive.MountPoints) > 0
		if f.Mounted != nil {
			if *f.Mounted && !isMounted {
				continue
			}
			if !*f.Mounted && isMounted {
				continue
			}
		}

		// --- Mount prefix filter ---
		if checkMountPrefix {
			match := false
			for _, mp := range drive.MountPoints {
				if strings.HasPrefix(mp.MountPoint, f.MountPrefix) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		result = append(result, drive)
	}

	return result
}

func GetDrives() ([]DriveInfo, error) {
	basePath := "/sys/block"
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	mountInfos := parseMounts("/proc/mounts")
	var drives []DriveInfo

	for _, e := range entries {
		name := e.Name()
		path := filepath.Join("/dev", name)

		blockDir := filepath.Join(basePath, name, "queue")

		sizeSectors := readUint(filepath.Join(basePath, name, "size"))
		logicalBlockSize := readUint(filepath.Join(blockDir, "logical_block_size"))
		physicalBlockSize := readUint(filepath.Join(blockDir, "physical_block_size"))

		if logicalBlockSize == 0 {
			logicalBlockSize = 512
		}
		if physicalBlockSize == 0 {
			physicalBlockSize = logicalBlockSize
		}

		sizeBytes := sizeSectors * logicalBlockSize

		isRotational := readUint(filepath.Join(blockDir, "rotational")) == 1
		model := readString(filepath.Join(basePath, name, "device/model"))
		vendor := readString(filepath.Join(basePath, name, "device/vendor"))
		devType := readString(filepath.Join(basePath, name, "device/type"))
		if len(devType) == 0 || devType == "0" {
			devType = "disk"
		} else {
			devType = "ssd"
		}
		drive := DriveInfo{
			Name:              name,
			Path:              path,
			SizeSectors:       sizeSectors,
			LogicalBlockSize:  logicalBlockSize,
			PhysicalBlockSize: physicalBlockSize,
			SizeBytes:         sizeBytes,
			IsRotational:      isRotational,
			Model:             model,
			Vendor:            vendor,
			Type:              devType,
		}

		// Check for mount info
		for devPath, mount := range mountInfos {
			if strings.HasPrefix(devPath, path) {
				drive.MountPoints = mount
				drive.FsType = mount[0].FsType
				drive.FsAvail = getFsAvailable(mount[0].MountPoint)
				break
			}
		}

		drives = append(drives, drive)
	}

	return drives, nil
}
func readUint(path string) uint64 {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	val, _ := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	return val
}

func readString(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func parseMounts(path string) map[string][]MountInfo {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	mounts := make(map[string][]MountInfo)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			mounts[fields[0]] = append(mounts[fields[0]], MountInfo{
				MountPoint: fields[1],
				FsType:     fields[2],
			})
		}
	}
	return mounts
}

func getFsAvailable(mount string) uint64 {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(mount, &stat); err != nil {
		return 0
	}
	return stat.Bavail * uint64(stat.Bsize)
}
