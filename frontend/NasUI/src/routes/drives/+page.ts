// src/routes/drives/+page.ts
// @ts-ignore
import type { PageLoad } from "./$types";
import { fetchSystemDrives} from "$lib/models/drive.js";

export const load: PageLoad = async () => {
    const drives = await fetchSystemDrives();

    return {
        drives
    };
};
