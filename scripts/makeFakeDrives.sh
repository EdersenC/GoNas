#!/bin/bash
# --- create_drives.sh ---
# Creates and attaches file-backed "fake drives" as loop devices
# ready for mdadm use (no formatting, just raw).
#
# Usage: sudo ./create_drives.sh <basename> <count> <size>
# Example: sudo ./create_drives.sh mydisk 3 20G
#
# After this script, you'll have /dev/loopX devices you can use with mdadm.

# Require root privileges
if [ "$EUID" -ne 0 ]; then
  echo "Please run this script with sudo."
  exit 1
fi

# Validate args
if [ "$#" -ne 3 ]; then
  echo "Usage: $0 <basename> <number_of_drives> <size_of_each_drive>"
  echo "Example: $0 myvol 3 20G"
  exit 1
fi

BASENAME=$1
NUM_DRIVES=$2
DRIVE_SIZE=$3
WORKDIR=$(pwd)

echo "==> Cleaning up old loop devices..."
losetup -D
rm -f "${WORKDIR}/${BASENAME}_"*.img

echo "==> Creating $NUM_DRIVES drives with base name '$BASENAME', size $DRIVE_SIZE each."

for i in $(seq 1 $NUM_DRIVES); do
  NAME="${BASENAME}_${i}"
  IMG_FILE="${WORKDIR}/${NAME}.img"

  echo "--- Creating drive $i ($IMG_FILE) ---"

  # Remove old file if exists
  [ -f "$IMG_FILE" ] && rm -f "$IMG_FILE"

  # Create raw image file
  truncate -s "$DRIVE_SIZE" "$IMG_FILE"

  # Attach to next available loop device
  LOOP=$(losetup -f --show "$IMG_FILE")

  echo "Attached $IMG_FILE -> $LOOP"
done

echo "==> All drives created and attached."
echo "You can now create your RAID array using mdadm, for example:"
echo "sudo mdadm --create --verbose /dev/md0 --level=0 --raid-devices=$NUM_DRIVES /dev/loop*"

echo "==> Verify devices with: losetup -a"
