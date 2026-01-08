package storage

import (
	"goNAS/helper"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

var DevFolder = "/dev/"

type DriveKey struct {
	Kind  string
	Value string
}

func (k DriveKey) String() string { return k.Kind + ":" + k.Value }

type DriveInfo struct {
	Name              string       `json:"name"`
	Uuid              string       `json:"uuid"`
	DriveKey          DriveKey     `json:"drive_key"`
	ByIds             []string     `json:"by_ids"`
	Wwid              string       `json:"wwid"`
	Path              string       `json:"path"`
	SizeSectors       uint64       `json:"size_sectors"`
	LogicalBlockSize  uint64       `json:"logical_block_size"`
	PhysicalBlockSize uint64       `json:"physical_block_size"`
	SizeBytes         uint64       `json:"size_bytes"`
	IsRotational      bool         `json:"is_rotational"`
	Model             string       `json:"model"`
	Vendor            string       `json:"vendor"`
	Serial            string       `json:"serial"`
	Type              string       `json:"type"`
	MountPoint        string       `json:"mountpoint"`
	Partitions        []*Partition `json:"partitions"`
	FsType            string       `json:"fstype"`
	FsAvail           uint64       `json:"fsavail"`
}

type Partition struct {
	Device     string
	MountPoint string
	FsType     string
	FsAvail    uint64
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

func (d *DriveInfo) GetUuid() string { return d.Uuid }

func (d *DriveInfo) SetUuid(uuid string) { d.Uuid = uuid }

func FilterFor(f DriveFilter, d ...*DriveInfo) []*DriveInfo {
	// Precompute small things to avoid recomputing inside loop
	hasNames := len(f.Names) > 0
	checkMountPrefix := f.MountPrefix != ""

	result := make([]*DriveInfo, 0, len(d))

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
		isMounted := len(drive.MountPoint) > 0
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
			for _, p := range drive.Partitions {
				if strings.HasPrefix(p.MountPoint, f.MountPrefix) {
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

func symlinksPointingToDev(dir string, devBase string) ([]string, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range ents {
		full := filepath.Join(dir, e.Name())
		target, err := os.Readlink(full)
		if err != nil {
			continue
		}
		// target often like "../../sdb" or "../../nvme0n1"
		if filepath.Base(target) == devBase {
			out = append(out, e.Name()) // store symlink NAME, not full path
		}
	}
	sort.Strings(out)
	return out, nil
}
func pickBestByID(byIDs []string) (DriveKey, bool) {
	preferPrefixes := []string{
		"wwn-",       // excellent for rotational/SAS/SATA
		"nvme-eui.",  // excellent for NVMe
		"nvme-uuid.", // good for NVMe (varies)
		"nvme-",      // usually includes model+serial
		"ata-",       // SATA model+serial
		"scsi-",      // may include WWN/serial
	}
	for _, p := range preferPrefixes {
		for _, id := range byIDs {
			if strings.HasPrefix(id, p) {
				return DriveKey{Kind: "by-id", Value: id}, true
			}
		}
	}
	return DriveKey{}, false
}

func (d *DriveInfo) generateDriveKey() {
	key, ok := pickBestByID(d.ByIds)
	if !ok && len(d.Wwid) > 0 {
		key = DriveKey{Kind: "wwid", Value: d.Wwid}
		ok = true
	}
	if !ok && len(d.Serial) > 0 {
		key = DriveKey{Kind: "serial", Value: d.Serial}
		ok = true
	}
	if !ok {
		keyString := d.Name + "_" + d.Model + "_" + d.Vendor + "_" + strconv.FormatUint(d.SizeBytes, 10)
		key = DriveKey{Kind: "hash", Value: keyString}
	}
	d.DriveKey = key
}

func GetDrives() ([]*DriveInfo, error) {
	basePath := "/sys/block"
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	partitions := parsePartitions("/proc/mounts")
	var drives []*DriveInfo

	for _, e := range entries {
		name := e.Name()
		path := filepath.Join(DevFolder, name)

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
		serial := readString(filepath.Join(basePath, name, "device/serial"))
		devType := readString(filepath.Join(basePath, name, "device/type"))
		byIDs, _ := symlinksPointingToDev("/dev/disk/by-id", name)
		wwid := readString(filepath.Join(basePath, name, "device/wwid"))
		if len(devType) == 0 || devType == "0" {
			devType = "disk"
		} else {
			devType = "ssd"
		}
		drive := &DriveInfo{
			ByIds:             byIDs,
			DriveKey:          DriveKey{},
			Wwid:              wwid,
			Name:              name,
			Path:              path,
			SizeSectors:       sizeSectors,
			LogicalBlockSize:  logicalBlockSize,
			PhysicalBlockSize: physicalBlockSize,
			SizeBytes:         sizeBytes,
			IsRotational:      isRotational,
			Model:             model,
			Vendor:            vendor,
			Serial:            serial,
			Type:              devType,
		}
		drive.generateDriveKey()

		// Check for p info
		for devPath, p := range partitions {
			if strings.HasPrefix(devPath, path) {
				drive.MountPoint = p[0].MountPoint
				drive.FsType = p[0].FsType
				drive.Partitions = p
				drive.FsAvail = getFsAvailable(p[0].MountPoint)
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

// todo: optimize
func parsePartitions(path string) map[string][]*Partition {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var partitions = make(map[string][]*Partition)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 || !strings.HasPrefix(fields[0], DevFolder) {
			continue
		}
		parentDrive := fields[0]
		kernelDrive := DevFolder + "sd"
		if strings.Contains(parentDrive, kernelDrive) {
			parentDrive = helper.StripTrailingDigits(parentDrive)
		}

		partitions[parentDrive] = append(partitions[parentDrive], &Partition{
			Device:     fields[0],
			MountPoint: fields[1],
			FsType:     fields[2],
			FsAvail:    getFsAvailable(fields[1]),
		})
	}
	return partitions
}

func getFsAvailable(mount string) uint64 {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(mount, &stat); err != nil {
		return 0
	}
	return stat.Bavail * uint64(stat.Bsize)
}
