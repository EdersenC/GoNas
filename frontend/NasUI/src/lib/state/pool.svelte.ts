export class PoolSelection {
    selectedDrives: string[] = $state([]);

    getSelectedDrives = (): string[] => {
        return this.selectedDrives;
    }

    isSelected = (driveId: string): boolean => {
        return this.selectedDrives.includes(driveId);
    }

    clearSelectedDrives = () => {
        this.selectedDrives = [];
        console.log("Cleared selected drives", $state.snapshot(this.selectedDrives));
    }

    addSelectedDrive = (driveId: string) => {
        if (!this.selectedDrives.includes(driveId)) {
            this.selectedDrives = [...this.selectedDrives, driveId];
        }
    }

    toggleSelectedDrive = (driveId: string) => {
        console.log("Toggling drive:", driveId);
        if (this.selectedDrives.includes(driveId)) {
            this.selectedDrives = this.selectedDrives.filter(id => id !== driveId);
        } else {
            this.selectedDrives = [...this.selectedDrives, driveId];
        }
        console.log("Selected drives after toggle:", $state.snapshot(this.selectedDrives));
    }
}

export const poolSelection = new PoolSelection();
