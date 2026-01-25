import {getContext, setContext} from "svelte";
import {fetchPools, type Pool} from "$lib/models/pool.js";

export class PoolManager{
    pools: Record<string, Pool> = $state({});
    loadingPools:boolean = $state(true);

    addPool= (poolId:string, pool:Pool) =>{
        this.pools = {...this.pools, [poolId]: pool};
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

    postPool = async (poolData: FormData, timeoutMs: number = 5000) => {
        const url = `http://localhost:8080/api/v1/pools`;

        const controller = new AbortController();
        const timer = setTimeout(() => controller.abort(), timeoutMs);

        try {
            const res = await fetch(url, {
                method: 'POST',
                body: poolData,
                signal: controller.signal
            });

            if (!res.ok) {
                throw new Error(`Failed to create pool: ${res.status}`);
            }

            const data = await res.json();
            const newPool: Pool = data.data;
            this.addPool(newPool.uuid, newPool); // todo make backend return the pool also
        } catch (err: any) {
            if (err && err.name === 'AbortError') {
                throw new Error(`Request timed out after ${timeoutMs} ms`);
            }
            console.error("Error creating pool:", err);
            throw err;
        } finally {
            clearTimeout(timer);
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
