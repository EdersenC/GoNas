package test

import (
	"goNAS/helper"
	"goNAS/storage"
	"testing"
)

var drives = []*storage.DriveInfo{
	{Name: "sda", SizeBytes: 500 * helper.Gigabyte, FsAvail: 400 * helper.Megabyte, Type: "HDD", Model: "Seagate ST500DM002"},
	{Name: "sdb", SizeBytes: 500 * helper.Gigabyte, FsAvail: 400 * helper.Megabyte, Type: "HDD", Model: "Seagate ST500DM002"},
	{Name: "sdc", SizeBytes: 1000 * helper.Gigabyte, FsAvail: 400 * helper.Megabyte, Type: "HDD", Model: "WD WD10EZEX"},
}

// TestPoolSize verifies capacity calculations on a new pool.
func TestPoolSize(t *testing.T) {
	var pools = storage.Pools{}
	testPool, _ := pools.NewPool("TestPool", storage.Standard, "", drives...)
	expectedCapacity := 2000 * helper.Gigabyte
	expectedAvailable := 1200 * helper.Megabyte
	if testPool.TotalCapacity != expectedCapacity {
		t.Errorf("Expected total capacity %d, got %d", expectedCapacity, testPool.TotalCapacity)
	}
	if testPool.AvailableCapacity != expectedAvailable {
		t.Errorf("Expected available capacity %d, got %d", expectedAvailable, testPool.AvailableCapacity)
	}
}

// TestGetAndRemoveDrive checks pool drive add/remove behavior.
func TestGetAndRemoveDrive(t *testing.T) {
	var pools = storage.Pools{}
	var err error
	newDrive := &storage.DriveInfo{Name: "sdd", SizeBytes: 2000 * helper.Gigabyte, FsAvail: 1500 * helper.Megabyte, Type: "HDD", Model: "Seagate ST2000DM008"}
	testPool, _ := pools.NewPool("TestPool", storage.Standard, "", newDrive)

	testPool.AddDrives(newDrive)
	myDrives := testPool.GetDrives(newDrive.Uuid)
	if len(myDrives) != 1 || myDrives[0].Uuid != newDrive.Uuid {
		t.Errorf("Expected to find 1 drive after adding, found %d", len(myDrives))
	}
	err = testPool.RemoveDrives(newDrive.Uuid)
	if err != nil {
		t.Errorf("Expected no error when removing drive, got %v", err)
	}
	myDrives = testPool.GetDrives(newDrive.Uuid)
	if len(myDrives) != 0 {
		t.Errorf("Expected to find 0 drives after removal, found %d", len(myDrives))
	}
}

// TestGetAndDeletePool verifies pool lookup and deletion.
func TestGetAndDeletePool(t *testing.T) {
	var pools = storage.Pools{}
	var err error
	testPool, _ := pools.NewPool("TestPool", storage.Standard, "", drives...)
	_, err = pools.GetPool(testPool.Uuid)
	if err != nil {
		t.Errorf("Expected no error when getting existing pool, got %v", err)
	}

	err = pools.DeletePool(testPool.Uuid)
	if err != nil {
		t.Errorf("Expeted no error when deleting pool, got %v ", err)
	}

	_, err = pools.GetPool(testPool.Uuid)
	if err == nil {
		t.Errorf("Expected error when getting deleted pool, got none %v", err)
	}

}
