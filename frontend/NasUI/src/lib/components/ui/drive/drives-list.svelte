<script lang="ts">
    import {UIDrive} from "$lib/components/ui/drive/index.js";
    import type {Drive} from "$lib/models/drive.ts";
    import {fetchDrives} from "$lib/models/drive.ts";
    import { onMount } from 'svelte';
    import { scale } from 'svelte/transition';
    import List from './list.svelte';

    let drives = $state<Record<string, Drive>>({});
    let loading = $state(true);
    let error = $state<string | null>(null);

    onMount(async () => {
        await loadDrives();
    });

    async function loadDrives() {
        loading = true;
        error = null;
        try {
            drives = await fetchDrives();
        } catch (e) {
            error = e instanceof Error ? e.message : 'Failed to fetch drives';
        } finally {
            loading = false;
        }
    }
</script>

<List label="Drives" {loading} {error} onRefresh={loadDrives}>
    {#each Object.entries(drives) as [id, drive], i (id)}
        <div in:scale={{ duration: 300, delay: i * 50, start: 0.8 }}>
            <UIDrive drive={drive} id={i} />
        </div>
    {/each}
</List>

