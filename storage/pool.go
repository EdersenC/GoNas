package storage

import (
	"fmt"
	"goNAS/helper"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

// ============================================================================
// Constants and Type Definitions
// ============================================================================

const SHORTUUIDLEN = 16

var Mirrored PoolType
var Standard PoolType

type Pool struct {
	Name              string
	Uuid              string
	Status            Status
	MountPoint        string
	MdDevice          string
	Type              PoolType
	TotalCapacity     uint64
	AvailableCapacity uint64
	Format            string
	CreatedAt         string
	AdoptedDrives     map[string]*AdoptedDrive
}

type PoolType interface {
	Build(*Pool) error
	Value() string
}

type Raid struct {
	Level int
}

type Status string

var Healthy Status = "healthy"
var Degraded Status = "degraded"
var Offline Status = "offline"

type Pools map[string]*Pool

// ============================================================================
// PoolType Interface Implementations
// ============================================================================

func (r *Raid) Value() string {
	return fmt.Sprintf("raid%d", r.Level)
}

func (r *Raid) Build(p *Pool) error {
	if err := helper.CheckRaidLevel(r.Level, len(p.AdoptedDrives)); err != nil {
		return err
	}
	if p.Format == "" {
		return ErrPoolFormatRequired
	}

	drives := make([]string, 0, len(p.AdoptedDrives))
	for _, d := range p.AdoptedDrives {
		drives = append(drives, DevFolder+d.Drive.Name)
	}
	args := append(
		[]string{
			"--create",
			"--verbose",
			p.MdDevice,
			fmt.Sprintf("--level=%d", r.Level),
			fmt.Sprintf("--raid-devices=%d", len(p.AdoptedDrives)),
			fmt.Sprintf("--name=%s", p.Name), //Todo check if name is valid or sanitize
		},
		drives...,
	)

	err := helper.BuildMadam(args)
	if err != nil {
		return err
	}

	// Format the RAID device
	if err = helper.FormatPool(p.Format, p.MdDevice); err != nil {
		return err
	}

	// Create and mount the mount point
	if err = helper.CreateMountPoint(p.Uuid, p.MdDevice); err != nil {
		return err
	}

	p.MountPoint = fmt.Sprintf("%s/%s", helper.DefaultMountPoint, p.Uuid)
	p.Status = Healthy
	return nil
}

// ============================================================================
// Status Type Methods
// ============================================================================

func (s Status) ToLower() Status {
	return Status(strings.ToLower(string(s)))
}

// ============================================================================
// Pools Collection Methods
// ============================================================================

func (p *Pools) GetPool(uuid string) (*Pool, error) {
	pool, exists := (*p)[uuid]
	if !exists {
		return nil, ErrPoolNotFound
	}
	return pool, nil
}

// AddPool adds a new pool to the Pools collection.
func (p *Pools) AddPool(pool *Pool) error {
	if _, exists := (*p)[pool.Uuid]; exists {
		return ErrPoolAlreadyExists
	}
	(*p)[pool.Uuid] = pool
	return nil
}

func (p *Pools) DeletePool(uuid string) error {
	pool, err := p.GetPool(uuid)
	if err != nil {
		return err
	}

	if pool.MountPoint != "" {
		if err = pool.Delete(); err != nil {
			return err
		}
	}

	delete(*p, uuid)
	return nil
}

// NewPool creates a new pool instance.
func (p *Pools) NewPool(name string, poolType PoolType, format string, drives ...*DriveInfo) (*Pool, error) {
	return NewPool(name, poolType, format, drives...)
}

// CreateAndAddPool creates a new pool and adds it to the Pools collection.
func (p *Pools) CreateAndAddPool(name string, poolType PoolType, format string, drives ...*DriveInfo) (*Pool, error) {
	pool, err := p.NewPool(name, poolType, format, drives...)
	if err != nil {
		return nil, err
	}
	err = p.AddPool(pool)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// ============================================================================
// Pool Lifecycle Methods
// ============================================================================

// NewPool creates a new pool instance.
func NewPool(name string, poolType PoolType, format string, drives ...*DriveInfo) (*Pool, error) {
	poolMap := make(map[string]*AdoptedDrive)
	poolId := uuid.New().String()
	for i := range drives {
		adoptedDrive := NewAdoptedDrive(drives[i])
		adoptedDrive.SetPoolID(poolId)
		poolMap[adoptedDrive.GetUuid()] = adoptedDrive
	}
	shortID, err := ShortUuid(SHORTUUIDLEN, poolId)
	if err != nil {
		return nil, err
	}
	pool := Pool{
		Name:          name,
		Uuid:          poolId,
		Status:        Offline,
		AdoptedDrives: poolMap,
		MdDevice:      "/dev/md/" + shortID,
		Type:          poolType,
		Format:        format,
		CreatedAt:     CreationTime(),
	}
	pool.CalculateTotalCapacity()
	pool.CalculateAvailableCapacity()
	return &pool, nil
}

// Build constructs the storage pool based on its type and format.
func (p *Pool) Build() error {
	return p.Type.Build(p)
}

func (p *Pool) Delete() error {
	if p.Status != Offline {
		return ErrPoolNotOffline
	}
	if err := p.UnmountDrive(); err != nil {
		return err
	}

	args := []string{"mdadm", "--remove", p.MdDevice}
	mdamRemove := exec.Command("sudo", args...)
	if err := mdamRemove.Run(); err != nil {
		return err
	}
	args = []string{"mdadm", "--stop", p.MdDevice}
	mdamStop := exec.Command("sudo", args...)
	if err := mdamStop.Run(); err != nil {
		return err
	}

	var drivePaths []string
	for _, d := range p.AdoptedDrives {
		drivePaths = append(drivePaths, d.Drive.Path)
	}
	args = append([]string{"mdadm", "--zero-superblock"}, drivePaths...)
	zeroOut := exec.Command("sudo", args...)
	zeroOut.Stderr = os.Stderr
	if err := zeroOut.Run(); err != nil {
		return err
	}
	return nil
}

// ============================================================================
// Pool Management Methods
// ============================================================================

// Clone creates a shadow copy of the pool.
func (p *Pool) Clone() *Pool {
	return &Pool{
		Name:              p.Name,
		Uuid:              p.Uuid,
		Status:            p.Status,
		AdoptedDrives:     p.AdoptedDrives,
		MountPoint:        p.MountPoint,
		MdDevice:          p.MdDevice,
		Type:              p.Type,
		TotalCapacity:     p.TotalCapacity,
		AvailableCapacity: p.AvailableCapacity,
		Format:            p.Format,
		CreatedAt:         p.CreatedAt,
	}
}

// ============================================================================
// Drive Management Methods
// ============================================================================

// AddDrives adds drives to the pool.
func (p *Pool) AddDrives(drive ...*DriveInfo) {
	for i := range drive {
		adoptedDrive := NewAdoptedDrive(drive[i])
		adoptedDrive.SetPoolID(p.Uuid)
		p.AdoptedDrives[adoptedDrive.GetUuid()] = adoptedDrive
	}
}

// GetDrives retrieves drives from the pool by their UUIDs.
func (p *Pool) GetDrives(uuids ...string) []*DriveInfo {
	var drives = make([]*DriveInfo, 0)
	for i := range p.AdoptedDrives {
		for _, id := range uuids {
			if p.AdoptedDrives[i].Drive.Uuid != id {
				continue
			}
			drives = append(drives, p.AdoptedDrives[i].Drive)
		}
	}
	return drives
}

// RemoveDrives removes drives from the pool by their UUIDs.
func (p *Pool) RemoveDrives(uuids ...string) error {
	toRemove := make(map[string]bool, len(uuids))
	for _, id := range uuids {
		toRemove[id] = true
	}

	if len(toRemove) == 0 {
		return ErrNoDrivesToRemove
	}

	for name, d := range p.AdoptedDrives {
		if toRemove[d.Drive.Uuid] {
			delete(p.AdoptedDrives, name)
		}
	}
	return nil
}

// ============================================================================
// Pool Property Setters
// ============================================================================

func (p *Pool) SetName(name string) {
	p.Name = name
}

func (p *Pool) SetStatus(status Status) {
	p.Status = status
}

func (p *Pool) SetFormat(format string) {
	p.Format = format
}

// ============================================================================
// Capacity Calculation Methods
// ============================================================================

func (p *Pool) CalculateTotalCapacity() {
	var total uint64
	for _, d := range p.AdoptedDrives {
		total += d.Drive.SizeBytes
	}
	p.TotalCapacity = total
}

func (p *Pool) CalculateAvailableCapacity() {
	var available uint64
	for _, d := range p.AdoptedDrives {
		available += d.Drive.FsAvail
	}
	p.AvailableCapacity = available
}

// ============================================================================
// Validation Functions
// ============================================================================

func ValidateStatus(status Status) error {
	switch status.ToLower() {
	case Healthy.ToLower(), Degraded.ToLower(), Offline.ToLower():
		return nil
	default:
		return ErrInvalidStatus
	}
}

func ValidatePoolFormat(format string) error {
	supportedFormats := []string{"ext4", "xfs", "btrfs"}
	for _, f := range supportedFormats {
		if format == f {
			return nil
		}
	}
	return ErrUnsupportedFormat
}

// ============================================================================
// Utility Functions
// ============================================================================

func ShortUuid(length int, uuid string) (string, error) {
	if len(uuid) < length {
		return "", ErrUuidTooShort
	}
	return uuid[:length], nil
}

func ParsePoolType(value string) (PoolType, error) {
	switch value {
	case "standard":
		return Standard, nil
	case "mirrored":
		return Mirrored, nil
	case "raid0":
		return &Raid{Level: 0}, nil
	case "raid1":
		return &Raid{Level: 1}, nil
	case "raid5":
		return &Raid{Level: 5}, nil
	case "raid6":
		return &Raid{Level: 6}, nil
	case "raid10":
		return &Raid{Level: 10}, nil
	default:
		return nil, ErrInvalidPoolType
	}
}

func (p *Pool) UnmountDrive() error {
	if err := exec.Command("sudo", "umount", p.MountPoint).Run(); err != nil {
		return fmt.Errorf("failed to unmount pool: %w", err)
	}
	if err := exec.Command("sudo", "rmdir", p.MountPoint).Run(); err != nil {
		return fmt.Errorf("failed to remove dir: %w", err)
	}
	return nil
}
