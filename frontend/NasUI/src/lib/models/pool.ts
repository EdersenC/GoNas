import type { Drive } from './drive.ts';

export type PoolDrive = {
    drive: Drive;
    uuid: string;
    poolID: string;
    createdAt: string;
}

export type Pool = {
    name: string;
    uuid: string;
    status: string;
    mountPoint: string;
    mdDevice: string;
    type: string;
    totalCapacity: number;
    availableCapacity: number;
    format: string;
    createdAt: string;
    AdoptedDrives?: Record<string, PoolDrive>;
}

export async function fetchPools(): Promise<Record<string, Pool>> {
    const res = await fetch("http://localhost:8080/api/v1/pools");

    if (!res.ok) {
        throw new Error(`Failed to load pools: ${res.status}`);
    }

    const data = await res.json();
    return data.data as Record<string, Pool>;
}