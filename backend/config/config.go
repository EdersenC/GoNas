package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidSize         = errors.New("invalid size")
	ErrEmptyBasename       = errors.New("empty basename")
	ErrInvalidBasename     = errors.New("invalid basename")
	ErrInvalidImageCount   = errors.New("invalid image count")
	ErrMissingSize         = errors.New("missing size")
	ErrRootRequired        = errors.New("root privileges required")
	ErrLosetupMissing      = errors.New("losetup not found")
	ErrImageCreateFailed   = errors.New("failed to create loop image")
	ErrImageTruncateFailed = errors.New("failed to truncate loop image")
	ErrImageCloseFailed    = errors.New("failed to close loop image")
	ErrLosetupAttachFailed = errors.New("failed to attach loop device")
)

// commandExists reports whether a command is available in PATH.
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// parseSize parses sizes like 10G or 512M into bytes.
func parseSize(size string) (int64, error) {
	s := strings.TrimSpace(size)
	re := regexp.MustCompile(`(?i)^(\d+)([kmgt])?b?$`)
	m := re.FindStringSubmatch(s)
	if len(m) < 2 {
		return 0, fmt.Errorf("%w: %s", ErrInvalidSize, size)
	}
	v, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return 0, err
	}
	suffix := strings.ToLower(m[2])
	switch suffix {
	case "k":
		return v * 1024, nil
	case "m":
		return v * 1024 * 1024, nil
	case "g":
		return v * 1024 * 1024 * 1024, nil
	case "t":
		return v * 1024 * 1024 * 1024 * 1024, nil
	default:
		return v, nil
	}
}

// sanitizeBasename ensures the basename is safe for filenames.
func sanitizeBasename(b string) (string, error) {
	if b == "" {
		return "", ErrEmptyBasename
	}
	ok, _ := regexp.MatchString(`^[A-Za-z0-9_-]+$`, b)
	if !ok {
		return "", ErrInvalidBasename
	}
	return b, nil
}

// CreateLoopImages creates `count` file-backed images in workdir named
// <basename>_N.img, truncates them to `size`, and attaches them as loop devices.
func CreateLoopImages(basename, workdir string, count int, size string) ([]string, error) {
	if count <= 0 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidImageCount, count)
	}
	if strings.TrimSpace(size) == "" {
		return nil, ErrMissingSize
	}
	if os.Geteuid() != 0 {
		return nil, ErrRootRequired
	}
	if !commandExists("losetup") {
		return nil, ErrLosetupMissing
	}
	bn, err := sanitizeBasename(basename)
	if err != nil {
		return nil, err
	}
	wd, err := filepath.Abs(workdir)
	if err != nil {
		return nil, err
	}
	sizeBytes, err := parseSize(size)
	if err != nil {
		return nil, err
	}
	if sizeBytes <= 0 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidSize, size)
	}

	// detach loop devices that reference previous images with this basename
	if out, err := exec.Command("losetup", "-a").CombinedOutput(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			if strings.Contains(line, bn+"_") && strings.Contains(line, ".img") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) > 0 {
					dev := strings.TrimSpace(parts[0])
					if dev != "" {
						_ = exec.Command("losetup", "-d", dev).Run()
					}
				}
			}
		}
	}

	// remove old image files
	pat := filepath.Join(wd, bn+"_"+"*.img")
	old, _ := filepath.Glob(pat)
	for _, f := range old {
		_ = os.Remove(f)
	}

	loops := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		img := filepath.Join(wd, fmt.Sprintf("%s_%d.img", bn, i))
		f, err := os.OpenFile(img, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %v", ErrImageCreateFailed, img, err)
		}
		if err := f.Truncate(sizeBytes); err != nil {
			_ = f.Close()
			return nil, fmt.Errorf("%w: %s: %v", ErrImageTruncateFailed, img, err)
		}
		if err := f.Close(); err != nil {
			return nil, fmt.Errorf("%w: %s: %v", ErrImageCloseFailed, img, err)
		}

		cmd := exec.Command("losetup", "-f", "--show", img)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %v; output=%s", ErrLosetupAttachFailed, img, err, strings.TrimSpace(string(out)))
		}
		dev := strings.TrimSpace(string(out))
		loops = append(loops, dev)
	}
	return loops, nil
}
