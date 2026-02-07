package api

import (
	"errors"
	"goNAS/DB"
	"goNAS/helper"
	"testing"
)

func TestValidatePoolPatchName(t *testing.T) {
	n := &Nas{}

	t.Run("rejects invalid name", func(t *testing.T) {
		patch := &DB.PoolPatch{Name: "dog Food"}
		err := n.ValidatePoolPatch(patch)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, helper.ErrInvalidRaidName) {
			t.Fatalf("expected ErrInvalidRaidName, got %v", err)
		}
	})

	t.Run("sanitizes valid name", func(t *testing.T) {
		patch := &DB.PoolPatch{Name: "  dogFood  "}
		err := n.ValidatePoolPatch(patch)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if patch.Name != "dogFood" {
			t.Fatalf("expected sanitized patch name dogFood, got %q", patch.Name)
		}
	})
}
