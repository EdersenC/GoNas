import {baseUrl, type Drive, fetchDrives} from "$lib/models/drive.js";
import {getContext, setContext} from "svelte";

export class DriveManager {

    adoptedDrives: Record<string, Drive> = $state({});
    loadingAdoptedDrives:boolean = $state(true);
    addAdoptedDrive= (driveId:string, drive:Drive) =>{
        this.adoptedDrives = {...this.adoptedDrives, [driveId]: drive};
    }

    systemDrives: Record<string, Drive> = $state({});
    loadingSystemDrives:boolean = $state(true);


    selectedDrives: string[] = $state([]);
    getSelectedDrives = (): string[] => {
        return this.selectedDrives;
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
        this.systemDrives = await fetchDrives(`http://${baseUrl}/drives`);
        this.loadingSystemDrives = false;
    }
    fetchAdoptedDrives = async () =>{
        this.loadingAdoptedDrives = true;
        this.adoptedDrives = await fetchDrives(`http://${baseUrl}/drives/adopted`);
        this.loadingAdoptedDrives = false;
    }

}
export const DriveManagerKey:Symbol = Symbol("DriveManagerType");

export function setDriveManagerContext(){
    return setContext(DriveManagerKey, new DriveManager());
}

export function getDriveManagerContext(){
    return getContext<ReturnType<typeof  setDriveManagerContext>>(DriveManagerKey);
}
