<script lang="ts">
    import {UIPool} from "$lib/components/ui/pool/index.js";
    import type {Pool} from "$lib/models/pool.ts";
    import {fetchPools} from "$lib/models/pool.ts";
    import { onMount } from 'svelte';
    import { scale } from 'svelte/transition';
    import List from '$lib/components/ui/drive/list.svelte';
    import {getPoolManagerContext, PoolManager} from "$lib/state/poolManager.svelte.js";

    let manager:PoolManager = getPoolManagerContext()
    let error = $state<string | null>(null);


    async function loadPools() {
        error = null;
        try {
            await manager.fetchPools();
        } catch (e) {
            error = `Failed to load pools: ${e}`;
        }
    }
</script>

<List label="Storage Pools" loading={manager.loadingPools} error={error} onRefresh={loadPools} maxColumns={2}>
    {#snippet children({ maxItems })}
        {#if !manager.pools || Object.keys(manager.pools).length === 0}
            <div class="p-4 text-center text-sm text-muted-foreground">
                No storage pools available.
            </div>
        {/if}
        {@const poolEntries = Object.entries(manager.pools)}
        {#each poolEntries.slice(0, maxItems ?? poolEntries.length) as [id, pool], i (id)}
            <div in:scale={{ duration: 300, delay: i * 50, start: 0.8 }}>
                <UIPool {pool} {id} />
            </div>
        {/each}
    {/snippet}
</List>
