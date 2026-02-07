import {getContext, setContext} from "svelte";
import {fetchPools, type Pool} from "$lib/models/pool.js";
import {fetchWithTimeout} from "$lib/utils/fetch.js";
import { AppErrorCode, responseError } from "$lib/errors.js";

export class PoolManager{
    pools: Record<string, Pool> = $state({});
    loadingPools:boolean = $state(true);
    creatingPool: boolean = $state(false);

    addPool= (poolId:string, pool:Pool) =>{
        this.pools = {...this.pools, [poolId]: pool};
    }

    removePool = (poolId: string) => {
        const {[poolId]: _, ...rest} = this.pools;
        this.pools = rest;
    }

    fetchPools = async () => {
        this.loadingPools = true;
        try {
            this.pools = await fetchPools()
            this.loadingPools = false;
        }catch (e) {
            console.error("Error fetching pools:", e);
            this.loadingPools = false;
            throw e
        }
    }

    createPool = async (poolData: FormData, timeoutMs: number = 5000) => {
        const url = `http://localhost:8080/api/v1/pool`;
        try {
            this.creatingPool = true;
            const res = await fetchWithTimeout(url, {
                method: 'POST',
                body: poolData,
            }, timeoutMs);

            if (!res.ok) {
                throw await responseError(res, AppErrorCode.CREATE_POOL_FAILED, "Failed to create pool");
            }

            const data = await res.json();
            const newPool: Pool = data.data;
            this.addPool(newPool.uuid, newPool); // todo make backend return the pool also
            this.creatingPool = false;
        } catch (err: any) {
            console.error("Error creating pool:", err);
            this.creatingPool = false;
            throw err;
        }
    }

    buildPool = async (poolId: string, timeoutMs: number = 10000) => {
        const url = `http://localhost:8080/api/v1/pool/${poolId}/build`;
        try {
            const res = await fetchWithTimeout(url, {
                method: 'POST',
            }, timeoutMs);

            if (!res.ok) {
                throw await responseError(res, AppErrorCode.BUILD_POOL_FAILED, "Failed to build pool");
            }

            await this.fetchPools();
        } catch (err: any) {
            console.error("Error building pool:", err);
            throw err;
        }
    }

    deletePool = async (poolId: string, timeoutMs: number = 10000) => {
        const url = `http://localhost:8080/api/v1/pool/${poolId}`;
        try {
            const res = await fetchWithTimeout(url, {
                method: 'DELETE',
            }, timeoutMs);

            if (!res.ok) {
                throw await responseError(res, AppErrorCode.DELETE_POOL_FAILED, "Failed to delete pool");
            }

            this.removePool(poolId);
        } catch (err: any) {
            console.error("Error deleting pool:", err);
            throw err;
        }
    }

}

export const PoolManagerKey:Symbol = Symbol("PoolManager");

export function setPoolManagerContext(){
    return setContext(PoolManagerKey, new PoolManager());
}


export function getPoolManagerContext(){
    return getContext<ReturnType<typeof setPoolManagerContext>>(PoolManagerKey);

}
