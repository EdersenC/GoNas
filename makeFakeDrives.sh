#!/bin/bash
# --- recreate_drives.sh ---
# Fully cleans up mdadm + loop devices, removes old images,
# then recreates fresh file-backed loop devices for mdadm testing.
#
# Usage: sudo ./recreate_drives.sh <basename> <count> <size>
# Example: sudo ./recreate_drives.sh mydisk 4 20G

set -euo pipefail

# ----- Root check -----
if [ "$EUID" -ne 0 ]; then
  echo "ERROR: Please run this script with sudo."
  exit 1
fi

# ----- Args -----
if [ "$#" -ne 3 ]; then
  echo "Usage: $0 <basename> <number_of_drives> <size>"
  echo "Example: $0 mydisk 4 20G"
  exit 1
fi

BASENAME="$1"
NUM_DRIVES="$2"
DRIVE_SIZE="$3"
WORKDIR="$(pwd)"

echo "==> STEP 1: Stopping any active mdadm arrays"
mdadm --stop --scan 2>/dev/null || true

echo "==> STEP 2: Detaching all loop devices"
losetup -D

echo "==> STEP 3: Zeroing RAID superblocks (if any)"
for dev in /dev/loop*; do
  [ -b "$dev" ] || continue
  mdadm --zero-superblock "$dev" 2>/dev/null || true
done

echo "==> STEP 4: Removing old image files"
rm -f "${WORKDIR}/${BASENAME}_"*.img

echo "==> STEP 5: Creating $NUM_DRIVES fresh images ($DRIVE_SIZE each)"

LOOPS=()

for i in $(seq 1 "$NUM_DRIVES"); do
  IMG="${WORKDIR}/${BASENAME}_${i}.img"

  echo "  - Creating $IMG"
  truncate -s "$DRIVE_SIZE" "$IMG"

  LOOP=$(losetup -f --show "$IMG")
  LOOPS+=("$LOOP")

  echo "    Attached -> $LOOP"
done

echo
echo "==> SUCCESS"
echo "Loop devices created:"
printf '  %s\n' "${LOOPS[@]}"

echo
echo "You can now safely create your RAID array, e.g.:"
echo "  sudo mdadm --create /dev/md0 --level=0 --raid-devices=$NUM_DRIVES ${LOOPS[*]}"

echo
echo "Verify with:"
echo "  lsblk"
echo "  losetup -a"
