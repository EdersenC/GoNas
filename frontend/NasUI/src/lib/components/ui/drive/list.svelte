<script lang="ts">
    import {UIDrive} from "$lib/components/ui/drive/index.js";
    import {Footer} from "$lib/components/ui/card/index.ts";
    import type {Drive} from "$lib/models/drive.ts";
    import  {fetchDrives} from "$lib/models/drive.ts";
    import { onMount } from 'svelte';
    import {Button} from "$lib/components/ui/button/index.js";

    let drives = $state<Record<string, Drive>>({});
    let loading = $state(true);
    let error = $state<string | null>(null);
    onMount(async () => {
        try {
            drives = await fetchDrives();
        } catch (e) {
            error = e instanceof Error ? e.message : 'Failed to fetch drives';
        } finally {
            loading = false;
        }
    });


    async function refresh() {
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


<div class="">
    {#if loading}
    <p>Loading...</p>
    {:else if error}
    <p class="error">Error: {error}</p>
    {:else}
        <header class="page-header">
            <h1>Drives</h1>
            <p class="subtitle">Manage your storage devices</p>
        </header>

        <div class="drives-grid">
            {#each Object.entries(drives) as [id, drive],i}
                <UIDrive drive ={drive} id={i} />
            {/each}
        </div>

        {#if Object.keys(drives).length === 0}
            <div class="empty-state">
                <p>No drives detected</p>
            </div>
        {/if}
        <Footer class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 p-4">
            <Button variant="green" onclick={refresh}>Nice Dog</Button>

            <Button
                    class="col-start-1 sm:col-start-2 md:col-start-3 lg:col-start-4 justify-self-end"
                    variant="green"
                    onclick={refresh}
            >
               Refresh Drives
            </Button>
        </Footer>
    {/if}
</div>


<style>
    .drives-page {
        padding: 2rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    .page-header {
        margin-bottom: 2rem;
    }

    .page-header h1 {
        font-size: 2rem;
        font-weight: 700;
        margin: 0 0 0.5rem 0;
    }

    .subtitle {
        color: #6b7280;
        margin: 0;
    }

    .drives-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
        gap: 1.5rem;
        max-width: 100%;
        background-color: #121212;
    }

    @media (min-width: 960px) {
        .drives-grid {
            grid-template-columns: repeat(4, 1fr);
        }
    }

    .empty-state {
        text-align: center;
        padding: 3rem;
        color: #6b7280;
        background: #f9fafb;
        border-radius: 0.5rem;
        border: 1px dashed #d1d5db;
    }
</style>
