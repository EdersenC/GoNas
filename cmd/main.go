package main

import (
	"fmt"
	"goNAS/api"
	"goNAS/helper"
	"goNAS/storage"
	"log"
	"os/exec"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("Fatal:", err)
	}
}

func run() error {
	nas := &api.Nas{POOLS: &storage.Pools{}}
	displayDrives()
	if err := helper.CreateLoopDevice("100G", 4); err != nil {
		return err
	}
	pool, err := createSystemPool(nas.POOLS)
	if err != nil {
		return err
	}
	err = api.Run(nas, ":8080")
	if err != nil {
		return err
	}
	defer func() {
		fmt.Println("Deleting pool:", pool.Name)
		pool.Status = storage.Offline
		_ = pool.Delete()
	}()
	displayDrives()
	return nil
}

func createSystemPool(pools *storage.Pools) (*storage.Pool, error) {
	drives := getSystemDrives("l")
	raid := storage.Raid{Level: 0} //Todo Raid 1 needs to be fixed
	myPool := pools.NewPool("DezNuts", &raid, nil, drives...)
	err := myPool.Build("mkfs.ext4")
	if err != nil {
		return nil, err
	}
	fmt.Println("Pool created:", myPool.Name)
	return myPool, nil

}

func getSystemDrives(names ...string) []*storage.DriveInfo {
	drives := getSystemDrives()
	drives = storage.FilterFor(storage.DriveFilter{
		Names:   names,
		MinSize: 1 * helper.Gigabyte,
	}, drives...)
	return drives
}

func displayDrives() {
	drives, _ := storage.GetDrives()
	for _, d := range drives {
		fmt.Printf("Device: %-6s | Type: %-5s | Size: %-8s | Model: %-20s | FSAvail: %s\n",
			d.Name, d.Type, helper.HumanSize(d.SizeBytes), d.Model, helper.HumanSize(d.FsAvail))
		for _, p := range d.Partitions {
			fmt.Printf("└─>%-6s |FSAvail: %s\n", p.MountPoint, helper.HumanSize(p.FsAvail))
		}
	}

	fmt.Println("done")

}

// zeroSuperblocks clears the RAID metadata (superblock) from a set of loop devices.
func zeroSuperblocks(deviceNumbers ...int) error {
	// List of loop devices to target (0 through 3 in your case)
	mdadmPath, err := exec.LookPath("mdadm")
	if err != nil {
		return fmt.Errorf("mdadm command not found: %w", err)
	}

	for _, i := range deviceNumbers {
		deviceName := fmt.Sprintf("/dev/loop%d", i)

		// The full command to execute: sudo mdadm --zero-superblock /dev/loopX
		cmd := exec.Command("sudo", mdadmPath, "--zero-superblock", deviceName)

		fmt.Printf("Executing: %s...\n", cmd.Args)

		// Run the command and capture the output and error
		output, er := cmd.CombinedOutput()
		if er != nil {
			// Print the output (which often contains the sudo error or mdadm error details)
			log.Printf("Error clearing superblock on %s: %s", deviceName, string(output))
			return fmt.Errorf("failed to execute mdadm on %s: %w", deviceName, er)
		}

		fmt.Printf("Successfully zeroed superblock on %s.\n", deviceName)
	}

	return nil
}
