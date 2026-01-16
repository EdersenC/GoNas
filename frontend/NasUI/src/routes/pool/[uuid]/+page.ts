import { error } from '@sveltejs/kit';
// @ts-ignore
import type { PageLoad } from "./$types";
// @ts-ignore
import type { Pool } from "$lib/models/pool.js";

export const load: PageLoad = async ({ params }) => {
    const res = await fetch(`http://localhost:8080/api/v1/pool/${params.uuid}`);

    if (!res.ok) {
        throw new Error(`Failed to load pool: ${res.status}`);
    }

    const data = await res.json();
    const pool:Pool = data.data;

    return {
        pool
    };
};