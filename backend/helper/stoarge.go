package helper

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
	"unicode/utf8"
)

var Gigabyte = uint64(1024 * 1024 * 1024)
var Megabyte = uint64(1024 * 1024)

// RAID-related errors
var (
	ErrRaid0RequiresDrives  = errors.New("raid0 requires at least 2 drives")
	ErrRaid1RequiresDrives  = errors.New("raid1 requires at least 2 drives")
	ErrRaid5RequiresDrives  = errors.New("raid5 requires at least 3 drives")
	ErrRaid6RequiresDrives  = errors.New("raid6 requires at least 4 drives")
	ErrRaid10RequiresDrives = errors.New("raid10 requires at least 4 drives and an even number of drives")
	ErrUnsupportedRaidLevel = errors.New("unsupported raid level")
)

func Contains(list []string, val string) bool {
	for _, v := range list {
		if strings.Contains(val, v) {
			return true
		}
	}
	return false
}

// StripTrailingDigits returns s with any trailing digit characters removed.
func StripTrailingDigits(s string) string {
	i := len(s)
	for i > 0 {
		r, size := utf8.DecodeLastRuneInString(s[:i])
		if unicode.IsDigit(r) {
			i -= size
			continue
		}
		break
	}
	return s[:i]
}

func CreateLoopDevice(size string, amount int) error {
	name := "testDrive"
	amountStr := fmt.Sprintf("%d", amount)
	filePath := "./backend/makeFakeDrives.sh"
	cmd := exec.Command("sudo", filePath, name, amountStr, size)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create Loop devices with size %s: %w", size, err)
	}
	return nil
}

// InstallMdadm ensures mdadm is installed on the system.
func installMdadm() error {
	// Check if mdadm already exists
	if _, err := exec.LookPath("mdadm"); err == nil {
		fmt.Println("mdadm already installed ✅")
		return nil
	}

	fmt.Println("mdadm not found — attempting to install...")

	// Detect available package manager
	var pm, installCmd string
	switch {
	case commandExists("apt"):
		pm = "apt"
		installCmd = "apt update -y && apt install -y mdadm"
	case commandExists("dnf"):
		pm = "dnf"
		installCmd = "dnf install -y mdadm"
	case commandExists("yum"):
		pm = "yum"
		installCmd = "yum install -y mdadm"
	case commandExists("pacman"):
		pm = "pacman"
		installCmd = "pacman -Sy --noconfirm mdadm"
	default:
		return fmt.Errorf("no supported package manager found (apt, dnf, yum, pacman)")
	}

	fmt.Printf("Installing mdadm using %s...\n", pm)

	cmd := exec.Command("bash", "-c", "sudo "+installCmd)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install mdadm with %s: %w", pm, err)
	}

	// Verify installation succeeded
	if _, err := exec.LookPath("mdadm"); err != nil {
		return fmt.Errorf("mdadm installation appears to have failed")
	}

	fmt.Println("mdadm successfully installed ✅")
	return nil
}

// commandExists checks if a command is available in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func CheckRaidLevel(level int, drives int) error {
	switch level {
	case 0:
		if drives < 2 {
			return ErrRaid0RequiresDrives
		}
	case 1:
		if drives < 2 {
			return ErrRaid1RequiresDrives
		}
	case 5:
		if drives < 3 {
			return ErrRaid5RequiresDrives
		}
	case 6:
		if drives < 4 {
			return ErrRaid6RequiresDrives
		}
	case 10:
		if drives < 4 || drives%2 != 0 {
			return ErrRaid10RequiresDrives
		}
	default:
		return ErrUnsupportedRaidLevel
	}
	return nil
}

var DefaultMountPoint = "/mnt/pools"

func BuildMadam(args []string) error {
	if err := installMdadm(); err != nil {
		return err
	}

	cmd := exec.Command("mdadm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create RAID array: %w", err)
	}
	return nil
}

// CreateMountPoint creates a mount point directory and mounts the given mdDevice there.
func CreateMountPoint(uuid string, mdDevice string) error {
	// Create and mount directory
	mountPoint := fmt.Sprintf("%s/%s", DefaultMountPoint, uuid)
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return fmt.Errorf("failed to create mount directory: %w", err)
	}

	if err := exec.Command("mount", mdDevice, mountPoint).Run(); err != nil {
		return fmt.Errorf("failed to mount RAID device: %w", err)
	}
	return nil
}

// FormatPool formats the given mdDevice with the specified format command.
func FormatPool(format string, mdDevice string) error {
	if err := exec.Command("mkfs."+format, "-F", mdDevice).Run(); err != nil {
		return fmt.Errorf("failed to format RAID device: %w", err)
	}
	return nil
}
