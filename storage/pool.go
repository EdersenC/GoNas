package storage

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"goNAS/helper"
	"goNAS/network"
)

var Mirrored PoolType
var Standard PoolType

type PoolType interface {
	Build(*Pool) error
}

type Raid struct {
	Level int
}

func (r *Raid) Build(p *Pool) error {
	if err := helper.CheckRaidLevel(r.Level, len(p.Drives)); err != nil {
		return err
	}
	mdDevice := "/dev/md/" + p.Name
	drives := make([]string, 0, len(p.Drives))
	for _, d := range p.Drives {
		drives = append(drives, DevFolder+d.Name)
	}
	args := append(
		[]string{
			"--create",
			"--verbose",
			mdDevice,
			fmt.Sprintf("--level=%d", r.Level),
			fmt.Sprintf("--raid-devices=%d", len(p.Drives)),
			fmt.Sprintf("--name=%s", p.Name),
		},
		drives...,
	)

	err := helper.BuildMadam(args)
	if err != nil {
		return err
	}

	// Format the RAID device
	if err = helper.FormatPool("mkfs.ext4", mdDevice); err != nil {
		return err
	}

	// Create and mount the mount point
	if err = helper.CreateMountPoint(p.Name, mdDevice); err != nil {
		return err
	}

	p.MdDevice = mdDevice
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

func (p *Pools) DeletePool(uuid string) error {
	//todo: implement pool deletion logic in linux
	delete(*p, uuid)
	return nil
}

func (p *Pools) NewPool(name string, poolType PoolType, network *network.Interface, drives ...*DriveInfo) *Pool {
	poolMap := make(map[string]*DriveInfo)
	poolId := uuid.New().String()
	for i, _ := range drives {
		poolMap[drives[i].Name] = drives[i]
	}
	pool := Pool{
		Name:    name,
		Uuid:    poolId,
		Status:  Offline,
		Drives:  poolMap,
		Type:    poolType,
		Network: network,
	}
	pool.CalculateTotalCapacity()
	pool.CalculateAvailableCapacity()
	(*p)[poolId] = &pool
	return &pool
}

type Pool struct {
	Name              string
	Uuid              string
	Status            Status
	Drives            map[string]*DriveInfo
	MountPoint        string
	MdDevice          string
	Type              PoolType
	TotalCapacity     uint64
	AvailableCapacity uint64
	Network           *network.Interface
}

func (p *Pool) BuildPool() error {
	return p.Type.Build(p)
}

func (p *Pool) AddDrives(drive ...*DriveInfo) {
	for i, _ := range drive {
		p.Drives[drive[i].Name] = drive[i]
	}
}

func (p *Pool) GetDrives(uuids ...string) []*DriveInfo {
	var drives = make([]*DriveInfo, 0)
	for i, _ := range p.Drives {
		for _, id := range uuids {
			if p.Drives[i].Uuid != id {
				continue
			}
			drives = append(drives, p.Drives[i])
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

	for name, d := range p.Drives {
		if toRemove[d.Uuid] {
			delete(p.Drives, name)
		}
	}
	return nil
}

func (p *Pool) CalculateTotalCapacity() {
	var total uint64
	for _, d := range p.Drives {
		total += d.SizeBytes
	}
	p.TotalCapacity = total
}

func (p *Pool) CalculateAvailableCapacity() {
	var available uint64
	for _, d := range p.Drives {
		available += d.FsAvail
	}
	p.AvailableCapacity = available
}
