<script lang="ts">
    import {UIPool} from "$lib/components/ui/pool/index.js";
    import type {Pool} from "$lib/models/pool.ts";
    import {fetchPools} from "$lib/models/pool.ts";
    import { onMount } from 'svelte';
    import { scale } from 'svelte/transition';
    import List from '$lib/components/ui/drive/list.svelte';

    let pools = $state<Record<string, Pool>>({});
    let loading = $state(true);
    let error = $state<string | null>(null);

    onMount(async () => {
        await loadPools();
    });

    async function loadPools() {
        loading = true;
        error = null;
        try {
            pools = await fetchPools();
        } catch (e) {
            error = e instanceof Error ? e.message : 'Failed to fetch pools';
        } finally {
            loading = false;
        }
    }
</script>

<List label="Storage Pools" {loading} {error} onRefresh={loadPools}>
    {#each Object.entries(pools) as [id, pool], i (id)}
        <div in:scale={{ duration: 300, delay: i * 50, start: 0.8 }}>
            <UIPool {pool} {id} />
        </div>
    {/each}
</List>
