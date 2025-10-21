package main

import (
	"fmt"
	"goNAS/helper"
	"goNAS/storage"
)

func main() {
	drives, err := storage.GetDrives()
	drives = storage.FilterDrives(drives, "loop")
	if err != nil {
		panic(err)
	}

	for _, d := range drives {
		fmt.Printf("Device: %-6s | Type: %-5s | Size: %-8s | Model: %-20s | Mount: %s\n| FSAvail: %s\n",
			d.Name, d.Type, helper.HumanSize(d.Size), d.Model, d.Mountpoint, helper.HumanSize(d.FSAvail))

		for _, c := range d.Children {
			fmt.Printf("  └─ Partition: %-6s | Size: %-8s | Mount: %s\n",
				c.Name, c.Size, c.Mountpoint)
		}
	}
	fmt.Println("done")
}
