<script lang="ts">
    import {UIDrive} from "$lib/components/ui/drive/index.js";
    import {type AdoptedDrive, type Drive, fetchAdoptedDrives} from "$lib/models/drive.ts";
    import {fetchSystemDrives} from "$lib/models/drive.ts";
    import { onMount } from 'svelte';
    import { scale } from 'svelte/transition';
    import List from './list.svelte';
    import type {PoolSelection} from "$lib/models/pool.js";

    // Cache drives data to prevent refetching on every mount

    let {
        ratio,
        poolCreatorMode,
        poolSelection,
    } = $props();


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
                adopted = false;
            }else {
                let adoptedDrives:AdoptedDrive = await fetchAdoptedDrives();
                for (const [key, drive] of Object.entries(adoptedDrives)) {
                    drives[key] = drive.drive;
                }
            }
            // Add fake drives for testing
            for (let i = 0; i < 20; i++) {
                drives[`fake-${i}`] = {
                    name: `Fake Drive ${i}`,
                    uuid: `fake-uuid-${i}`,
                    drive_key: { kind: 'hash', value: `fake-value-${i}` },
                    by_ids: null,
                    wwid: '',
                    path: `/dev/fake${i}`,
                    size_sectors: 209715200,
                    logical_block_size: 512,
                    physical_block_size: 512,
                    size_bytes: 107374182400 + i * 1000000000,
                    is_rotational: i % 2 === 0,
                    model: `Fake Model ${i}`,
                    vendor: `Fake Vendor ${i}`,
                    serial: `Fake Serial ${i}`,
                    type: 'disk',
                    mountpoint: i < 5 ? `/mnt/fake${i}` : '',
                    partitions: null,
                    fstype: i < 5 ? 'ext4' : '',
                    fsavail: i < 5 ? 50000000 : 0
                };
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

<List label="Drives" {loading} {error} onRefresh={loadDrives} ratio={ratio} >
    {#each Object.entries(drives) as [id, drive], i (id)}
        <div in:scale={{ duration: 300, delay: i * 50, start: 0.8 }}>
            <UIDrive
                    adopted={adopted}
                    drive={drive} ]
                    id={i}
                    poolCreatorMode={poolCreatorMode}
                    poolSelection={poolSelection}/>
        </div>
    {/each}
</List>