// @ts-ignore
import type { PageLoad } from "./$types";
// @ts-ignore
import type { Pool } from "$lib/models/pool.js";
import { PoolErrorCode, responseError } from "$lib/errors.js";

// @ts-ignore
export const load: PageLoad = async ({ params }) => {
    const res = await fetch(`http://localhost:8080/api/v1/pool/${params.uuid}`);

    if (!res.ok) {
        throw await responseError(res, PoolErrorCode.FETCH_POOL_FAILED, "Failed to load pool");
    }

    const data = await res.json();
    const pool:Pool = data.data;

    return {
        pool
    };
};
