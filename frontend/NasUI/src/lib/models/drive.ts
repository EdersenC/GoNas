import { fetchWithTimeout } from "$lib/utils/fetch.js";

export type DriveKey = {
    kind: string;
    value: string;
}

export type Partition = {
    device: string;
    mountPoint: string;
    fsType: string;
    fsAvail: number;
}

export type Drive = {
    name: string;
    uuid: string;
    drive_key: DriveKey;
    by_ids: string[];
    wwid: string;
    path: string;
    size_sectors: number;
    logical_block_size: number;
    physical_block_size: number;
    size_bytes: number;
    is_rotational: boolean;
    model: string;
    vendor: string;
    serial: string;
    type: string;
    mountpoint: string;
    partitions: Partition[];
    fstype: string;
    fsavail: number;
}


export type AdoptedDrive = {
    drive : Drive;
    uuid: string;
    created_at: string;
    poolId : string;
}


export const baseUrl = "localhost:8080/api/v1";


export async function fetchSystemDrives(timeoutMs: number = 5000): Promise<Record<string, Drive>> {
    const url = `http://${baseUrl}/drives`
    const res = await fetchWithTimeout(url, {}, timeoutMs);

    if (!res.ok) {
        throw new Error(`Failed to load drives: ${res.status}`);
    }

    const data = await res.json();

    return data.data as Record<string, Drive>;
}


export async function fetchAdoptedDrives(timeoutMs: number = 5000): Promise<Record<string, AdoptedDrive>> {
    const url = `http://${baseUrl}/drives/adopted`
    const res = await fetchWithTimeout(url, {}, timeoutMs);

    if (!res.ok) {
        throw new Error(`Failed to load adopted drives: ${res.status}`);
    }
    const data = await res.json();

    return data.data as Record<string, AdoptedDrive>;
}
import { fetchWithTimeout } from "$lib/utils/fetch.js";
