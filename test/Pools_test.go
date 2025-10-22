package test

import (
	"goNAS/storage"
	"testing"
)

var gigabyte = uint64(1024 * 1024 * 1024)
var megabyte = uint64(1024 * 1024)

var drives = []*storage.DriveInfo{
	{Name: "sda", SizeBytes: 500 * gigabyte, FsAvail: 400 * megabyte, Type: "HDD", Model: "Seagate ST500DM002"},
	{Name: "sdb", SizeBytes: 500 * gigabyte, FsAvail: 400 * megabyte, Type: "HDD", Model: "Seagate ST500DM002"},
	{Name: "sdc", SizeBytes: 1000 * gigabyte, FsAvail: 400 * megabyte, Type: "HDD", Model: "WD WD10EZEX"},
}

func TestPoolSize(t *testing.T) {
	pool := storage.NewPool("TestPool", storage.Standard, nil, drives...)
	expectedCapacity := 2000 * gigabyte
	expectedAvailable := 1200 * megabyte
	if pool.TotalCapacity != expectedCapacity {
		t.Errorf("Expected total capacity %d, got %d", expectedCapacity, pool.TotalCapacity)
	}
	if pool.AvailableCapacity != expectedAvailable {
		t.Errorf("Expected available capacity %d, got %d", expectedAvailable, pool.AvailableCapacity)
	}
}
