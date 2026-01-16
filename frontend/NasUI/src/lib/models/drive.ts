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

export async function fetchDrives(): Promise<Record<string, Drive>> {
    const res = await fetch("http://localhost:8080/api/v1/drives");

    if (!res.ok) {
        throw new Error(`Failed to load drives: ${res.status}`);
    }

    const data = await res.json();
    return data.data as Record<string, Drive>;
}