// @ts-ignore
import type { PageLoad } from "./$types";
// @ts-ignore
import type { Pool } from "$lib/models/pool.js";

// @ts-ignore
export const load: PageLoad = async ({fetch}) => {
    const res = await fetch("http://localhost:8080/api/v1/pools");

    if (!res.ok) {
        throw new Error(`Failed to load drives: ${res.status}`);
    }


    const data = await res.json();
    const pools: Map<string, Pool> = data.data;

    return {
       pools: pools
    };
};
