package DB

import (
	"context"
	"goNAS/storage"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestDatabaseOperations(t *testing.T) {
	// Create a temporary database file using t.TempDir() for automatic cleanup
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/test_gonas.db"

	// Initialize database
	db := NewDB(tmpFile)
	defer db.Close()

	ctx := context.Background()
	if err := db.InitSchema(ctx); err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	t.Run("Pool Operations", func(t *testing.T) {
		// Create a test pool
		poolID := uuid.New().String()
		pool := &storage.Pool{
			Uuid:          poolID,
			Name:          "TestPool",
			MdDevice:      "/dev/md0",
			Status:        storage.Offline,
			Type:          &storage.Raid{Level: 5},
			Format:        "ext4",
			AdoptedDrives: make(map[string]*storage.AdoptedDrive),
		}

		// Test Insert
		createdAt := time.Now().UTC().Format(time.RFC3339Nano)
		if err := db.InsertPool(ctx, pool, createdAt); err != nil {
			t.Fatalf("Failed to insert pool: %v", err)
		}

		// Test Query All
		pools, err := db.QueryAllPools(ctx)
		if err != nil {
			t.Fatalf("Failed to query pools: %v", err)
		}
		if len(pools) != 1 {
			t.Errorf("Expected 1 pool, got %d", len(pools))
		}
		if retrievedPool, ok := pools[poolID]; !ok {
			t.Error("Pool not found in query results")
		} else {
			if retrievedPool.Name != "TestPool" {
				t.Errorf("Expected pool name 'TestPool', got '%s'", retrievedPool.Name)
			}
			if retrievedPool.Status != storage.Offline {
				t.Errorf("Expected status 'offline', got '%s'", retrievedPool.Status)
			}
		}

		// Test Update (Patch)
		patch := &PoolPatch{
			Name:   "UpdatedPool",
			Status: storage.Healthy,
		}
		updatedPool, err := db.PatchPool(ctx, pool, patch)
		if err != nil {
			t.Fatalf("Failed to patch pool: %v", err)
		}
		if updatedPool.Name != "UpdatedPool" {
			t.Errorf("Expected updated name 'UpdatedPool', got '%s'", updatedPool.Name)
		}
		if updatedPool.Status != storage.Healthy {
			t.Errorf("Expected updated status 'healthy', got '%s'", updatedPool.Status)
		}

		// Test PatchPoolMount
		if err := db.PatchPoolMount(poolID, "/mnt/test"); err != nil {
			t.Fatalf("Failed to patch pool mount: %v", err)
		}

		// Verify mount point was updated
		pools, _ = db.QueryAllPools(ctx)
		if pools[poolID].MountPoint != "/mnt/test" {
			t.Errorf("Expected mount point '/mnt/test', got '%s'", pools[poolID].MountPoint)
		}

		// Test Delete
		if err := db.DeletePool(ctx, poolID); err != nil {
			t.Fatalf("Failed to delete pool: %v", err)
		}

		// Verify deletion
		pools, _ = db.QueryAllPools(ctx)
		if len(pools) != 0 {
			t.Errorf("Expected 0 pools after deletion, got %d", len(pools))
		}
	})

	t.Run("Drive Operations", func(t *testing.T) {
		// Create a test pool first (needed for foreign key)
		poolID := uuid.New().String()
		pool := &storage.Pool{
			Uuid:          poolID,
			Name:          "TestPoolForDrive",
			MdDevice:      "/dev/md0",
			Status:        storage.Offline,
			Type:          &storage.Raid{Level: 1},
			Format:        "ext4",
			AdoptedDrives: make(map[string]*storage.AdoptedDrive),
		}
		poolCreatedAt := time.Now().UTC().Format(time.RFC3339Nano)
		if err := db.InsertPool(ctx, pool, poolCreatedAt); err != nil {
			t.Fatalf("Failed to insert pool: %v", err)
		}

		// Create a test drive
		driveID := uuid.New().String()
		drive := &storage.DriveInfo{
			DriveKey: storage.DriveKey{
				Kind:  "serial",
				Value: "ABC123",
			},
			Uuid: driveID,
		}

		// Test Insert
		createdAt := time.Now().UTC().Format(time.RFC3339Nano)
		if err := db.InsertDrive(ctx, drive, createdAt); err != nil {
			t.Fatalf("Failed to insert drive: %v", err)
		}

		// Test QueryDriveByKey
		adoptedDrive, found, err := db.QueryDriveByKey(ctx, drive.DriveKey)
		if err != nil {
			t.Fatalf("Failed to query drive by key: %v", err)
		}
		if !found {
			t.Error("Drive not found by key")
		}
		if adoptedDrive.GetUuid() != driveID {
			t.Errorf("Expected drive UUID '%s', got '%s'", driveID, adoptedDrive.GetUuid())
		}

		// Test QueryAllAdoptedDrives
		drives, err := db.QueryAllAdoptedDrives(ctx)
		if err != nil {
			t.Fatalf("Failed to query all drives: %v", err)
		}
		if len(drives) < 1 {
			t.Errorf("Expected at least 1 drive, got %d", len(drives))
		}

		// Test PatchDrive (associate with pool)
		patch := DrivePatch{PoolID: &poolID}
		if err := db.PatchDrive(ctx, driveID, patch); err != nil {
			t.Fatalf("Failed to patch drive: %v", err)
		}

		// Verify patch
		adoptedDrive, _, _ = db.QueryDriveByKey(ctx, drive.DriveKey)
		if adoptedDrive.GetPoolID() != poolID {
			t.Errorf("Expected poolID '%s', got '%s'", poolID, adoptedDrive.GetPoolID())
		}
	})

	t.Run("Foreign Key Constraint", func(t *testing.T) {
		// Create a pool
		poolID := uuid.New().String()
		pool := &storage.Pool{
			Uuid:          poolID,
			Name:          "FKTestPool",
			MdDevice:      "/dev/md1",
			Status:        storage.Offline,
			Type:          &storage.Raid{Level: 1},
			Format:        "ext4",
			AdoptedDrives: make(map[string]*storage.AdoptedDrive),
		}
		createdAt := time.Now().UTC().Format(time.RFC3339Nano)
		db.InsertPool(ctx, pool, createdAt)

		// Create a drive associated with the pool
		driveID := uuid.New().String()
		drive := &storage.DriveInfo{
			DriveKey: storage.DriveKey{
				Kind:  "serial",
				Value: "XYZ789",
			},
			Uuid: driveID,
		}
		db.InsertDrive(ctx, drive, createdAt)
		patch := DrivePatch{PoolID: &poolID}
		db.PatchDrive(ctx, driveID, patch)

		// Delete the pool - drive's poolID should be set to NULL due to foreign key
		db.DeletePool(ctx, poolID)

		// Verify drive still exists but poolID is empty
		drives, _ := db.QueryAllAdoptedDrives(ctx)
		foundDrive := false
		for _, d := range drives {
			if d.GetUuid() == driveID {
				foundDrive = true
				if d.GetPoolID() != "" {
					t.Errorf("Expected poolID to be empty after pool deletion, got '%s'", d.GetPoolID())
				}
			}
		}
		if !foundDrive {
			t.Error("Drive should still exist after pool deletion")
		}
	})
}
