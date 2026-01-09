package storage

import (
	"errors"
	"fmt"
	"goNAS/helper"
	"goNAS/network"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

var Mirrored PoolType
var Standard PoolType

type Pool struct {
	Name              string
	Uuid              string
	Status            Status
	AdoptedDrives     map[string]*AdoptedDrive
	MountPoint        string
	MdDevice          string
	Type              PoolType
	TotalCapacity     uint64
	AvailableCapacity uint64
	Network           *network.Interface
	CreatedAt         string
}

type PoolType interface {
	Build(*Pool, string) error
	Value() string
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
		return nil, errors.New("invalid pool type")
	}
}

type Raid struct {
	Level int
}

func (r *Raid) Value() string {
	return fmt.Sprintf("raid%d", r.Level)
}

func (r *Raid) Build(p *Pool, format string) error {
	if err := helper.CheckRaidLevel(r.Level, len(p.AdoptedDrives)); err != nil {
		return err
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
			fmt.Sprintf("--name=%s", p.Name),
		},
		drives...,
	)

	err := helper.BuildMadam(args)
	if err != nil {
		return err
	}

	// Format the RAID device
	if err = helper.FormatPool(format, p.MdDevice); err != nil {
		return err
	}

	// Create and mount the mount point
	if err = helper.CreateMountPoint(p.Name, p.MdDevice); err != nil {
		return err
	}

	p.MountPoint = fmt.Sprintf("%s/%s", helper.DefaultMountPoint, p.Name)
	p.Status = Healthy
	return nil
}

type Status string

var Healthy Status = "healthy"
var Degraded Status = "degraded"
var Offline Status = "offline"

type Pools map[string]*Pool

func (p *Pools) GetPool(uuid string) (*Pool, error) {
	pool, exists := (*p)[uuid]
	if !exists {
		return nil, errors.New("pool not found")
	}
	return pool, nil
}

func (p *Pool) UnmountDrive() error {
	if err := exec.Command("sudo", "umount", p.MountPoint).Run(); err != nil {
		return fmt.Errorf("failed to unmount drive: %w", err)
	}
	return nil
}

func (p *Pool) Delete() error {
	if p.Status != Offline {
		return errors.New("cannot delete a pool that is not offline")
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

func (p *Pools) DeletePool(uuid string) error {
	pool, err := p.GetPool(uuid)
	if err != nil {
		return err
	}
	if err = pool.Delete(); err != nil {
		return err
	}

	delete(*p, uuid)
	return nil
}

// NewPool creates a new pool instance.
func NewPool(name string, poolType PoolType, network *network.Interface, drives ...*DriveInfo) (*Pool, error) {
	poolMap := make(map[string]*AdoptedDrive)
	poolId := uuid.New().String()
	for i, _ := range drives {
		adoptedDrive := NewAdoptedDrive(drives[i])
		adoptedDrive.SetPoolID(poolId)
		poolMap[adoptedDrive.GetUuid()] = adoptedDrive
	}
	pool := Pool{
		Name:          name,
		Uuid:          poolId,
		Status:        Offline,
		AdoptedDrives: poolMap,
		MdDevice:      "/dev/md/" + name,
		Type:          poolType,
		Network:       network,
		CreatedAt:     CreationTime(),
	}
	pool.CalculateTotalCapacity()
	pool.CalculateAvailableCapacity()
	return &pool, nil
}
func (p *Pools) NewPool(name string, poolType PoolType, network *network.Interface, drives ...*DriveInfo) (*Pool, error) {
	return NewPool(name, poolType, network, drives...)
}

// CreateAndAddPool creates a new pool and adds it to the Pools collection.
func (p *Pools) CreateAndAddPool(name string, poolType PoolType, network *network.Interface, drives ...*DriveInfo) (*Pool, error) {
	pool, err := p.NewPool(name, poolType, network, drives...)
	if err != nil {
		return nil, err
	}
	err = p.AddPool(pool)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// AddPool adds a new pool to the Pools collection.
func (p *Pools) AddPool(pool *Pool) error {
	if _, exists := (*p)[pool.Uuid]; exists {
		return errors.New("pool with the same UUID already exists")
	}
	(*p)[pool.Uuid] = pool
	return nil
}

// Build constructs the storage pool based on its type and format.
func (p *Pool) Build(format string) error {
	return p.Type.Build(p, format)
}

func (p *Pool) AddDrives(drive ...*DriveInfo) {
	for i, _ := range drive {
		adoptedDrive := NewAdoptedDrive(drive[i])
		adoptedDrive.SetPoolID(p.Uuid)
		p.AdoptedDrives[adoptedDrive.GetUuid()] = adoptedDrive
	}
}

func (p *Pool) GetDrives(uuids ...string) []*DriveInfo {
	var drives = make([]*DriveInfo, 0)
	for i, _ := range p.AdoptedDrives {
		for _, id := range uuids {
			if p.AdoptedDrives[i].Drive.Uuid != id {
				continue
			}
			drives = append(drives, p.AdoptedDrives[i].Drive)
		}
	}
	return drives
}

func (p *Pool) RemoveDrives(uuids ...string) error {
	toRemove := make(map[string]bool, len(uuids))
	for _, id := range uuids {
		toRemove[id] = true
	}

	if len(toRemove) == 0 {
		return errors.New("no drives to remove")
	}

	for name, d := range p.AdoptedDrives {
		if toRemove[d.Drive.Uuid] {
			delete(p.AdoptedDrives, name)
		}
	}
	return nil
}

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
