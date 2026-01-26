<script lang="ts">
    import {UIDrive} from "$lib/components/ui/drive/index.js";
    import { scale } from 'svelte/transition';
    import List from './list.svelte';
    import {DriveManagerKey, getDriveManagerContext, DriveManager} from "$lib/state/driveManager.svelte.js";
    import type {AdoptedDrive} from "$lib/models/drive.js";

    // Cache drives data to prevent refetching on every mount

    let {
        ratio,
        poolCreatorMode,
    } = $props();

    let manager: DriveManager = getDriveManagerContext(DriveManagerKey);
    let adopted= true
    let error = $state<string | null>(null);
    let loading = $derived<boolean>(adopted ? manager.loadingAdoptedDrives : manager.loadingSystemDrives);


    async function loadDrives() {
        try {
            error = null;
            await manager.fetchAdoptedDrives()
        } catch (e) {
            console.error('Failed to load drives', e);
            error = 'Failed to load drives';
        }
    }

</script>

<List label="Drives" {loading} {error} onRefresh={loadDrives} ratio={ratio} compact={poolCreatorMode} >
    {@const drivesEntries = Object.entries(manager.adoptedDrives)}
    {#if drivesEntries.length === 0}
        <div class="p-4 text-center text-sm text-muted-foreground">

            No drives available for adoption.
        </div>
    {/if}
    {#each Object.entries(manager.adoptedDrives) as [id, adoptedDrive], i (id)}
        <div in:scale={{ duration: 300, delay: i * 50, start: 0.8 }}>
            <UIDrive
                    adopted={adopted}
                    drive={adoptedDrive.drive}
                    id={i}
                    poolCreatorMode={poolCreatorMode}
            />
        </div>
    {/each}
</List>
