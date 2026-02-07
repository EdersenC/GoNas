package storage

import (
	"errors"
	"goNAS/helper"
	"testing"
)

func TestNewPoolSanitizesName(t *testing.T) {
	pool, err := NewPool("  dogFood  ", &Raid{Level: 5}, "ext4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pool.Name != "dogFood" {
		t.Fatalf("expected sanitized name %q, got %q", "dogFood", pool.Name)
	}
}

func TestNewPoolRejectsInvalidName(t *testing.T) {
	_, err := NewPool("dog Food", &Raid{Level: 5}, "ext4")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, helper.ErrInvalidRaidName) {
		t.Fatalf("expected ErrInvalidRaidName, got %v", err)
	}
}
