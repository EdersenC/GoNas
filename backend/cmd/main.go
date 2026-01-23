package main

import (
	"context"
	"fmt"
	"goNAS/DB"
	"goNAS/api"
	"goNAS/helper"
	"goNAS/storage"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

// main wires dependencies, initializes the database, and starts the API server.
func main() {
	db := DB.NewDB("Drives.db")
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing database: %v", err)
		}
	}()

	err := db.InitSchema(context.Background())
	if err != nil {
		log.Fatalf("Error initializing database schema: %v", err)
	}

	server := api.NewAPIServer(parsePort(), db)
	if err = run(server); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}

// parsePort returns the HTTP listen address from argv or the default :8080.
func parsePort() string {
	port := ":8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
		hasColon := false
		for i := 0; i < len(port); i++ {
			if port[i] == ':' {
				hasColon = true
				break
			}
		}
		if !hasColon && len(port) > 0 {
			port = ":" + port
		}
	}
	return port
}

// run prepares loop devices, starts the server, and blocks for shutdown signals.
func run(server *api.Server) error {
	if err := helper.CreateLoopDevice("100G", 4); err != nil {
		return err
	}
	//pool, err := createSystemPool(server.Nas.POOLS, 5)
	//if err != nil {
	//	return err
	//}
	//defer func() {
	//	fmt.Println("Deleting pool:", pool.Name)
	//	pool.Status = storage.Offline
	//	_ = pool.Delete()
	//	displayDrives()
	//}()

	go runServer(server)
	sigCh := make(chan os.Signal, 4)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for s := range sigCh {
			log.Printf("received signal: %v", s)
		}
	}()
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	defer stop()
	graceFull(server, ctx)
	return nil
}

// runServer launches the HTTP server in a background goroutine.
func runServer(server *api.Server) {
	go func() {
		err := server.Start()
		if err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
}

// graceFull waits for a stop signal and shuts down the server with a timeout.
func graceFull(server *api.Server, ctx context.Context) {
	<-ctx.Done()
	fmt.Println("Shutting down gracefully, press Ctrl+C again to force")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

// createSystemPool builds a pool from system drives with the requested RAID level.
func createSystemPool(pools *storage.Pools, level int) (*storage.Pool, error) {
	drives := storage.GetSystemDrives("l")
	raid := storage.Raid{Level: level} //Todo Raid 1 needs to be fixed
	myPool, err := pools.CreateAndAddPool("DezNuts", &raid, "", drives...)
	if err != nil {
		return nil, err
	}
	myPool.SetFormat("mkfs.ext4")
	err = myPool.Build()
	if err != nil {
		return nil, err
	}
	fmt.Println("Pool created:", myPool.Name)
	displayDrives()
	return myPool, nil

}

// displayDrives prints detected drive and partition details to stdout.
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
