package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func parseSize(size string) (int64, error) {
	s := strings.TrimSpace(size)
	re := regexp.MustCompile(`(?i)^(\d+)([kmgt])?b?$`)
	m := re.FindStringSubmatch(s)
	if len(m) < 2 {
		return 0, fmt.Errorf("invalid size: %s", size)
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

func sanitizeBasename(b string) (string, error) {
	if b == "" {
		return "", fmt.Errorf("empty basename")
	}
	ok, _ := regexp.MatchString(`^[A-Za-z0-9_-]+$`, b)
	if !ok {
		return "", fmt.Errorf("basename contains invalid characters")
	}
	return b, nil
}

// CreateLoopImages creates `count` file-backed images in workdir named
// <basename>_N.img, truncates them to `size`, and attaches them as loop devices.
func CreateLoopImages(basename, workdir string, count int, size string) ([]string, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be > 0")
	}
	if strings.TrimSpace(size) == "" {
		return nil, fmt.Errorf("size must be provided")
	}
	if os.Geteuid() != 0 {
		return nil, fmt.Errorf("operation requires root privileges")
	}
	if !commandExists("losetup") {
		return nil, fmt.Errorf("losetup not found")
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
		return nil, fmt.Errorf("size must be > 0")
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
			return nil, fmt.Errorf("create %s: %w", img, err)
		}
		if err := f.Truncate(sizeBytes); err != nil {
			_ = f.Close()
			return nil, fmt.Errorf("truncate %s: %w", img, err)
		}
		if err := f.Close(); err != nil {
			return nil, fmt.Errorf("close %s: %w", img, err)
		}

		cmd := exec.Command("losetup", "-f", "--show", img)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("losetup %s failed: %w; output: %s", img, err, strings.TrimSpace(string(out)))
		}
		dev := strings.TrimSpace(string(out))
		loops = append(loops, dev)
	}
	return loops, nil
}
