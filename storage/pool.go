package storage

import (
	"goNAS/network"
)

type PoolType string

var Mirrored PoolType = "mirrored"
var Standard PoolType = "standard"
var Raid5 PoolType = "raid5"

type Status string

var Healthy Status = "healthy"
var Degraded Status = "degraded"
var Offline Status = "offline"

type Pool struct {
	Name              string
	Status            Status
	Drives            []*DriveInfo
	Type              PoolType
	TotalCapacity     uint64
	AvailableCapacity uint64
	Network           *network.Interface
}

func (p *Pool) AddDrive(drive ...*DriveInfo) {
	for i, _ := range drive {
		p.Drives = append(p.Drives, drive[i])
	}
}

func (p *Pool) RemoveDrive(drive *DriveInfo) {
	var updatedDrives []*DriveInfo
	for _, d := range p.Drives {
		if d.Name != drive.Name {
			updatedDrives = append(updatedDrives, d)
		}
	}
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

func NewPool(name string, poolType PoolType, network *network.Interface, drives ...*DriveInfo) *Pool {
	pool := &Pool{
		Name:    name,
		Drives:  drives,
		Type:    poolType,
		Network: network,
	}
	pool.CalculateTotalCapacity()
	pool.CalculateAvailableCapacity()
	return pool
}
