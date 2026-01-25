<script lang="ts">
    import {setDriveManagerContext} from "$lib/state/driveManager.svelte.js";

    let { children } = $props();
    import '../app.css'
    import {setPoolManagerContext} from "$lib/state/poolManager.svelte.js";
    setDriveManagerContext()
    setPoolManagerContext()
</script>

<style>
    :global(body) {
        background-color: var(--color-canvas);
        color: var(--color-canvas-foreground);
    }

    nav {
        background-color: var(--color-brand);
        padding: 1rem;
        display: flex;
        gap: 1rem;
    }

    nav a {
        color: var(--color-brand-foreground);
        text-decoration: none;
        font-weight: bold;
    }

    nav a:hover {
        text-decoration: underline;
    }


</style>


{@render children()}
