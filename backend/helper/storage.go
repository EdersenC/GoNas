package helper

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"goNAS/config"
)

var Gigabyte = uint64(1024 * 1024 * 1024)
var Megabyte = uint64(1024 * 1024)

// RAID-related errors
var (
	ErrRaid0RequiresDrives   = errors.New("raid0 requires at least 2 drives")
	ErrRaid1RequiresDrives   = errors.New("raid1 requires at least 2 drives")
	ErrRaid5RequiresDrives   = errors.New("raid5 requires at least 3 drives")
	ErrRaid6RequiresDrives   = errors.New("raid6 requires at least 4 drives")
	ErrRaid10RequiresDrives  = errors.New("raid10 requires at least 4 drives and an even number of drives")
	ErrUnsupportedRaidLevel  = errors.New("unsupported raid level")
	ErrInvalidSizeInput      = errors.New("invalid size input")
	ErrInvalidAmountInput    = errors.New("invalid amount input")
	ErrRootPrivilegesNeeded  = errors.New("root privileges required")
	ErrWorkdirResolve        = errors.New("failed to determine workdir")
	ErrPackageManagerMissing = errors.New("supported package manager not found")
	ErrMdadmInstall          = errors.New("failed to install mdadm")
	ErrMdadmInstallVerify    = errors.New("mdadm installation verification failed")
	ErrMdadmArgsEmpty        = errors.New("mdadm arguments are required")
	ErrMdadmBuild            = errors.New("mdadm command failed")
	ErrMountPointCreate      = errors.New("failed to create mount point")
	ErrMountRaidDevice       = errors.New("failed to mount raid device")
	ErrFormatRaidDevice      = errors.New("failed to format raid device")
	ErrInvalidRaidName       = errors.New("invalid raid name")
)

// Contains reports whether val is contained within any element of list.
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

// CreateLoopDevice creates file-backed loop devices for testing.
func CreateLoopDevice(size string, amount int) error {
	name := "testDrive"
	if strings.TrimSpace(size) == "" {
		return fmt.Errorf("%w: empty string", ErrInvalidSizeInput)
	}
	if amount <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidAmountInput, amount)
	}

	if os.Geteuid() != 0 {
		return ErrRootPrivilegesNeeded
	}

	workdir, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWorkdirResolve, err)
	}

	loops, err := config.CreateLoopImages(name, workdir, amount, size)
	if err != nil {
		return err
	}

	if len(loops) > 0 {
		fmt.Printf("Loop devices created:\n")
		for _, l := range loops {
			fmt.Printf("  %s\n", l)
		}
	}
	return nil
}

// installMdadm ensures mdadm is installed on the system.
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
		return ErrPackageManagerMissing
	}

	fmt.Printf("Installing mdadm using %s...\n", pm)

	cmd := exec.Command("bash", "-c", "sudo "+installCmd)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: manager=%s err=%v", ErrMdadmInstall, pm, err)
	}

	// Verify installation succeeded
	if _, err := exec.LookPath("mdadm"); err != nil {
		return ErrMdadmInstallVerify
	}

	fmt.Println("mdadm successfully installed ✅")
	return nil
}

// commandExists checks if a command is available in PATH.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

var raidNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// SanitizeRaidName validates and normalizes a RAID array name for mdadm.
func SanitizeRaidName(name string) (string, error) {
	sanitized := strings.TrimSpace(name)
	if sanitized == "" {
		return "", fmt.Errorf("%w: name is empty", ErrInvalidRaidName)
	}
	if !raidNamePattern.MatchString(sanitized) {
		return "", fmt.Errorf("%w: %q must match %s", ErrInvalidRaidName, sanitized, raidNamePattern.String())
	}
	return sanitized, nil
}

// CheckRaidLevel validates RAID level constraints against drive count.
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

// BuildMdadm runs mdadm with the provided args to create a RAID array.
func BuildMdadm(args []string) error {
	if len(args) == 0 {
		return ErrMdadmArgsEmpty
	}

	if err := installMdadm(); err != nil {
		return err
	}

	cmd := exec.Command("mdadm", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	outStr := strings.TrimSpace(stdout.String())
	errStr := strings.TrimSpace(stderr.String())

	if err != nil {
		// Try to extract exit code when available
		exitCode := -1
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			// ExitError exposes ExitCode() in recent Go versions
			exitCode = ee.ExitCode()
		}
		cmdStr := "mdadm " + strings.Join(args, " ")
		return fmt.Errorf("%w: cmd=%q exit=%d err=%v stdout=%s stderr=%s", ErrMdadmBuild, cmdStr, exitCode, err, outStr, errStr)
	}

	// Print any non-empty output for debugging
	if outStr != "" {
		fmt.Printf("mdadm stdout: %s\n", outStr)
	}
	if errStr != "" {
		fmt.Printf("mdadm stderr: %s\n", errStr)
	}
	return nil
}

// CreateMountPoint creates a mount point directory and mounts the given mdDevice there.
func CreateMountPoint(uuid string, mdDevice string) error {
	// Create and mount directory
	mountPoint := fmt.Sprintf("%s/%s", DefaultMountPoint, uuid)
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return fmt.Errorf("%w: %v", ErrMountPointCreate, err)
	}

	if err := exec.Command("mount", mdDevice, mountPoint).Run(); err != nil {
		return fmt.Errorf("%w: %v", ErrMountRaidDevice, err)
	}
	return nil
}

// FormatPool formats the given mdDevice with the specified format command.
func FormatPool(format string, mdDevice string) error {
	if err := exec.Command("mkfs."+format, "-F", mdDevice).Run(); err != nil {
		return fmt.Errorf("%w: %v", ErrFormatRaidDevice, err)
	}
	return nil
}
