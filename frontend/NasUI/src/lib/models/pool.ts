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

export async function fetchPools(timeoutMs: number = 5000): Promise<Record<string, Pool>> {
    const url = "http://localhost:8080/api/v1/pools";

    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), timeoutMs);

    try {
        const res = await fetch(url, { signal: controller.signal });

        if (!res.ok) {
            throw new Error(`Failed to load pools: ${res.status}`);
        }

        const data = await res.json();
        return data.data as Record<string, Pool>;
    } catch (err: any) {
        if (err && err.name === 'AbortError') {
            throw new Error(`Request timed out after ${timeoutMs} ms`);
        }
        throw err;
    } finally {
        clearTimeout(timer);
    }
}
