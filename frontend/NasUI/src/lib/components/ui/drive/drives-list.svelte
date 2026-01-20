<script lang="ts">
    import {UIDrive} from "$lib/components/ui/drive/index.js";
    import {type AdoptedDrive, type Drive, fetchAdoptedDrives} from "$lib/models/drive.ts";
    import {fetchSystemDrives} from "$lib/models/drive.ts";
    import { onMount } from 'svelte';
    import { scale } from 'svelte/transition';
    import List from './list.svelte';

    // Cache drives data to prevent refetching on every mount
    let cachedDrives: Record<string, Drive> | null = null;
    let cacheLoading = false;

    let drives = $state<Record<string, Drive>>({});
    let adopted= true
    let loading = $state(true);
    let error = $state<string | null>(null);

    onMount(async () => {
        // Only fetch if we don't have cached data and aren't already loading
        if (!cachedDrives && !cacheLoading) {
            await loadDrives();
        } else if (cachedDrives) {
            // Use cached data
            drives = cachedDrives;
            loading = false;
        }
    });

    async function loadDrives() {
        cacheLoading = true;
        loading = true;
        error = null;
        try {
            if (window.location.pathname.endsWith("/drives")) {
                drives = await fetchSystemDrives();
            }else {
                let adoptedDrives:AdoptedDrive = await fetchAdoptedDrives();
                for (const [key, drive] of Object.entries(adoptedDrives)) {
                    drives[key] = drive.drive;
                }
            }
            cachedDrives = drives;
        } catch (e) {
            error = e instanceof Error ? e.message : 'Failed to fetch drives';
        } finally {
            loading = false;
            cacheLoading = false;
        }
    }
</script>

<List label={adopted ? "Adopted Drives" : "System Drives"} {loading} {error} onRefresh={loadDrives}>
    {#each Object.entries(drives) as [id, drive], i (id)}
        <div in:scale={{ duration: 300, delay: i * 50, start: 0.8 }}>
            <UIDrive drive={drive} id={i} />
        </div>
    {/each}
</List>