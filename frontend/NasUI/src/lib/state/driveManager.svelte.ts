import {type AdoptedDrive, type Drive, fetchAdoptedDrives, fetchSystemDrives} from "$lib/models/drive.js";
import {getContext, setContext} from "svelte";

export class DriveManager {

    adoptedDrives: Record<string, AdoptedDrive> = $state({});
    loadingAdoptedDrives:boolean = $state(true);
    addAdoptedDrive= (driveId:string, drive:AdoptedDrive) =>{
        this.adoptedDrives = {...this.adoptedDrives, [driveId]: drive};
    }
    addAdoptedDrives= (drives:Record<string, AdoptedDrive>) =>{
        this.adoptedDrives = {...this.adoptedDrives, ...drives};
    }

    removeAdoptedDrive= (driveId:string) =>{
        const {[driveId]: _, ...rest} = this.adoptedDrives;
        this.adoptedDrives = rest;
    }

    removeSelectedDrivesFromAdopted = () => {
        this.selectedDrives.forEach(driveId => {
            this.removeAdoptedDrive(driveId);
        });
        this.clearSelectedDrives();
    }

    systemDrives: Record<string, Drive> = $state({});
    loadingSystemDrives:boolean = $state(true);


    selectedDrives: string[] = $state([]);
    getSelectedDrives = (): string[] => {
        return $state.snapshot(this.selectedDrives);
    }

    isSelected = (driveId: string): boolean => {
        return this.selectedDrives.includes(driveId);
    }

    clearSelectedDrives = () => {
        this.selectedDrives = [];
    }

    toggleSelectedDrive = (driveId: string) => {
        console.log("Toggling drive:", driveId);
        if (this.selectedDrives.includes(driveId)) {
            this.selectedDrives = this.selectedDrives.filter(id => id !== driveId);
        } else {
            this.selectedDrives = [...this.selectedDrives, driveId];
        }
    }

    fetchSystemDrives = async () =>{
        this.loadingSystemDrives = true;
        try {
            this.systemDrives = await fetchSystemDrives();
            this.loadingSystemDrives = false;
        } catch (e) {
            console.error("Error fetching system drives:", e);
            this.loadingSystemDrives = false;
            throw e
        }
    }
    fetchAdoptedDrives = async () =>{
        this.loadingAdoptedDrives = true;
        try {
            this.adoptedDrives = await fetchAdoptedDrives();
            this.loadingAdoptedDrives = false;
        } catch (e) {
            console.error("Error fetching adopted drives:", e);
            this.loadingAdoptedDrives = false;
            throw e
        }
    }

    adopt = async (driveId: string) => {
        if (!driveId) return;
        let url = `http://localhost:8080/api/v1/drives/adopt/${driveId}`;
        if (window.location.pathname.startsWith('/pools')) {
            url = `http://localhost:8080/api/v1/pools/adopt/${driveId}`;
        }
        const res = await fetch(url, {
            method: 'POST',
        });
        console.log(await res.json())
        console.log(`Adopting drive with ID: ${driveId}`);
    }

}
export const DriveManagerKey:Symbol = Symbol("DriveManager");

export function setDriveManagerContext(){
    return setContext(DriveManagerKey, new DriveManager());
}

export function getDriveManagerContext(){
    return getContext<ReturnType<typeof  setDriveManagerContext>>(DriveManagerKey);
}
