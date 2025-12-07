package main

import (
	"fmt"
	"goNAS/helper"
	"goNAS/storage"
	"log"
	"os/exec"
)

func main() {
	displayDrives()
	err := helper.CreateLoopDevice("100G", 4)
	if err != nil {
		fmt.Println("Error creating loop devices:", err)
		return
	}
	displayDrives()
	pool, err := createSystemPool()
	if err != nil {
		fmt.Println("Error creating pool:", err)
		return
	}
	displayDrives()
	pool.Status = storage.Offline
	err = pool.Delete()
	if err != nil {
		fmt.Println("Error deleting pool:", err)
		return
	}
	displayDrives()
}

func createSystemPool() (*storage.Pool, error) {
	var pools = storage.Pools{}
	drives := getSystemDrives("l")
	raid := storage.Raid{Level: 10}
	myPool := pools.NewPool("DezNuts", &raid, nil, drives...)
	err := myPool.Build("mkfs.ext4")
	if err != nil {
		return nil, err
	}
	fmt.Println("Pool created:", myPool.Name)
	return myPool, nil
}

func getSystemDrives(names ...string) []*storage.DriveInfo {
	drives, err := storage.GetDrives()
	drives = storage.Filter(storage.DriveFilter{
		Names:   names,
		MinSize: 1 * helper.Gigabyte,
	}, drives...)
	if err != nil {
		panic(err)
	}
	return drives
}

func displayDrives() {
	drives := getSystemDrives()
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
