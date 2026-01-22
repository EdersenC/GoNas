<script lang="ts">
    import type { Drive } from "$lib/models/drive";
    import {UIDrive} from "$lib/components/ui/drive/index.js";
    export let data: {
        drives: Record<string, Drive>;
    };

    console.log(data.drives)
</script>

<div class="drives-page">
    <header class="page-header">
        <h1>Drives</h1>
        <p class="subtitle">Manage your storage devices</p>
    </header>

    <div class="drives-grid">
        {#each Object.entries(data.drives) as [id, drive],i}
            <UIDrive drive ={drive} id={i} />
        {/each}
    </div>

    {#if Object.keys(data.drives).length === 0}
        <div class="empty-state">
            <p>No drives detected</p>
        </div>
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
        color: var(--color-muted-foreground);
        margin: 0;
    }

    .drives-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
        gap: 1.5rem;
        max-width: 100%;
    }

    @media (min-width: 960px) {
        .drives-grid {
            grid-template-columns: repeat(4, 1fr);
        }
    }

    .empty-state {
        text-align: center;
        padding: 3rem;
        color: var(--color-muted-foreground);
        background: var(--color-card);
        border-radius: 0.5rem;
        border: 1px dashed var(--color-border);
    }
</style>
