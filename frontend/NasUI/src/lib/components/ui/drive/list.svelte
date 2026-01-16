<script lang="ts">
    import {Header} from "$lib/components/ui/card/index.ts";
    import {ReloadButton} from "$lib/components/ui/button/index.js";
    import {AspectRatio} from "$lib/components/ui/aspect-ratio/index.js";
    import { Spinner } from "$lib/components/ui/spinner/index.js";
    import { fade } from 'svelte/transition';

    interface Props {
        label?: string;
        loading?: boolean;
        error?: string | null;
        onRefresh?: () => void | Promise<void>;
        children?: import('svelte').Snippet;
    }

    let {
        label = "Content",
        loading = false,
        error = null,
        onRefresh,
        children
    }: Props = $props();

    async function handleRefresh() {
        if (onRefresh) {
            await onRefresh();
        }
    }
</script>

<div class="grid grid-cols-1 gap-y-6 sm:gap-y-10">
    <div class="flex items-center justify-center px-2 sm:px-40 py-0">
        <div class="w-full corner-container">
            <div class="corner-content">
                <Header class="p-2 bg-transparent text-blue-500
                flex items-center justify-between gap-2">
                    <span class="text-2xl font-bold tracking-wide">{label}</span>
                    {#if onRefresh}
                        <ReloadButton
                            class="ml-auto text-blue-500 hover:text-blue-400"
                            onclick={handleRefresh}
                            isSpinning={loading}
                        />
                    {/if}
                </Header>
                <AspectRatio ratio={2} class="bg-[#121212] overflow-hidden">
                    <div class="h-full flex flex-col min-h-0 bg-[#121212]">
                        {#if loading}
                            <div class="flex-1 flex items-center justify-center" in:fade={{ duration: 200 }}>
                                <Spinner class="size-20 text-blue-500" />
                            </div>
                        {:else if error}
                            <div class="flex-1 flex items-center justify-center" in:fade={{ duration: 300 }}>
                                <p class="error text-red-500">Error: {error}</p>
                            </div>
                        {:else}
                            <div class="flex-1 overflow-y-auto min-h-0 bg-[#121212]" in:fade={{ duration: 300 }}>
                                <div class="content-grid">
                                    {#if children}
                                        {@render children()}
                                    {:else}
                                        <div class="empty-state" in:fade={{ duration: 300 }}>
                                            <p>No content available</p>
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        {/if}
                    </div>
                </AspectRatio>
            </div>
        </div>
    </div>
</div>






<style>
    .corner-container {
        position: relative;
    }

    .corner-content {
        position: relative;
    }

    .corner-container::before,
    .corner-container::after {
        content: '';
        position: absolute;
        width: 30px;
        height: 30px;
        border-color: #3b82f6;
        border-style: solid;
        z-index: 10;
        filter: drop-shadow(0 0 8px #3b82f6) drop-shadow(0 0 12px #3b82f680);
    }

    /* Top-left corner */
    .corner-container::before {
        top: 0;
        left: 0;
        border-width: 4px 0 0 4px;
    }

    /* Top-right corner */
    .corner-container::after {
        top: 0;
        right: 0;
        border-width: 4px 4px 0 0;
    }

    /* Bottom corners */
    .corner-content::before,
    .corner-content::after {
        content: '';
        position: absolute;
        width: 30px;
        height: 30px;
        border-color: #3b82f6;
        border-style: solid;
        pointer-events: none;
        z-index: 10;
        filter: drop-shadow(0 0 8px #3b82f6) drop-shadow(0 0 12px #3b82f680);
    }

    /* Bottom-left corner */
    .corner-content::before {
        bottom: 0;
        left: 0;
        border-width: 0 0 4px 4px;
    }

    /* Bottom-right corner */
    .corner-content::after {
        bottom: 0;
        right: 0;
        border-width: 0 4px 4px 0;
    }

    .content-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
        gap: 1.5rem;
        padding: 1.5rem;
        max-width: 100%;
        background-color: #121212;
    }

    @media (min-width: 960px) {
        .content-grid {
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
