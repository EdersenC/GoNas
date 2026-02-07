import type { Drive } from './drive.ts';
import { fetchWithTimeout } from "$lib/utils/fetch.js";
import { AppErrorCode, responseError } from "$lib/errors.js";

export type PoolDrive = {
    drive: Drive;
    uuid: string;
    poolID: string;
    createdAt: string;
}

export type Type = {
    name: string;
    Level : string;
}

export type Pool = {
    name: string;
    uuid: string;
    status: string;
    mountPoint: string;
    mdDevice: string;
    type: Type;
    totalCapacity: number;
    availableCapacity: number;
    format: string;
    createdAt: string;
    AdoptedDrives?: Record<string, PoolDrive>;
}

export async function fetchPools(timeoutMs: number = 5000): Promise<Record<string, Pool>> {
    const url = "http://localhost:8080/api/v1/pools";
    const res = await fetchWithTimeout(url, {}, timeoutMs);

    if (!res.ok) {
        throw await responseError(res, AppErrorCode.FETCH_POOLS_FAILED, "Failed to load pools");
    }

    const data = await res.json();
    return data.data as Record<string, Pool>;
}
