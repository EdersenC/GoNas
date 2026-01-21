import { writable } from 'svelte/store';

// Store holds an array of selected drive UUIDs for pool creation
const { subscribe, set, update } = writable<string[]>([]);

function toggleDrive(uuid: string | undefined) {
    if (!uuid) return;
    update(list => {
        const idx = list.indexOf(uuid);
        if (idx === -1) {
            return [...list, uuid];
        } else {
            const copy = [...list];
            copy.splice(idx, 1);
            return copy;
        }
    });
}

function addDrive(uuid: string) {
    update(list => (list.includes(uuid) ? list : [...list, uuid]));
}

function removeDrive(uuid: string) {
    update(list => list.filter(id => id !== uuid));
}

function clear() {
    set([]);
}

export const selectedDrives = { subscribe };
export const selectedDrivesActions = { toggleDrive, addDrive, removeDrive, clear };

export default selectedDrives;
