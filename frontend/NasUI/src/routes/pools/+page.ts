// @ts-ignore
import type { PageLoad } from "./$types";
// @ts-ignore
import type { Pool } from "$lib/models/pool.js";
import { AppErrorCode, responseError } from "$lib/errors.js";

// @ts-ignore
export const load: PageLoad = async ({fetch}) => {
    const res = await fetch("http://localhost:8080/api/v1/pools");

    if (!res.ok) {
        throw await responseError(res, AppErrorCode.FETCH_POOLS_FAILED, "Failed to load pools");
    }


    const data = await res.json();
    const pools: Map<string, Pool> = data.data;

    return {
       pools: pools
    };
};
