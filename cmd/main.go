package main

import (
	"fmt"
	"goNAS/helper"
	"goNAS/storage"
)

func main() {
	drives, err := storage.GetDrives()
	drives = storage.Filter(storage.DriveFilter{
		Names:      []string{},
		MinFsAvail: 1,
	}, drives...)
	if err != nil {
		panic(err)
	}

	for _, d := range drives {
		fmt.Printf("Device: %-6s | Type: %-5s | Size: %-8s | Model: %-20s | FSAvail: %s\n",
			d.Name, d.Type, helper.HumanSize(d.SizeBytes), d.Model, helper.HumanSize(d.FsAvail))
		for _, mp := range d.MountPoints {
			fmt.Printf("└─>%s\n", mp.MountPoint)
		}
	}

	fmt.Println("done")
}
