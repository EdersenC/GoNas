package helper

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseHumanSize converts strings like "500G", "1.5T", "200M" → bytes
func ParseHumanSize(s string) uint64 {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return 0
	}

	units := map[string]uint64{
		"B": 1,
		"K": 1024,
		"M": 1024 * 1024,
		"G": 1024 * 1024 * 1024,
		"T": 1024 * 1024 * 1024 * 1024,
		"P": 1024 * 1024 * 1024 * 1024 * 1024,
	}

	// find unit suffix
	for unit, mult := range units {
		if strings.HasSuffix(s, unit) {
			num := strings.TrimSuffix(s, unit)
			f, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return 0
			}
			return uint64(f * float64(mult))
		}
	}

	// default: assume bytes
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return uint64(f)
}

func HumanSize(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
