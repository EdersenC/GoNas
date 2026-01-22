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
        ratio = 2,
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
        <div class="w-full">
            <Header
                    class="p-2 bg-transparent text-brand flex items-center justify-between gap-2"
            >
                <span class="text-2xl font-bold tracking-wide">{label}</span>

                {#if onRefresh}
                    <ReloadButton
                            class="ml-auto text-brand hover:text-brand/80"
                            onclick={handleRefresh}
                            isSpinning={loading}
                    />
                {/if}
            </Header>

            <AspectRatio ratio={ratio} class="bg-canvas overflow-hidden rounded-lg">
                <!-- Ratio box -->
                <div class="h-full flex flex-col min-h-0 bg-canvas">
                    <!-- Scroll surface (ALWAYS present) -->
                    <div
                            class="flex-1 min-h-0 overflow-y-auto overscroll-contain bg-canvas"
                            style="padding-bottom: calc(6rem + env(safe-area-inset-bottom));"
                    >
                        {#if loading}
                            <div
                                    class="min-h-full flex items-center justify-center"
                                    in:fade={{ duration: 200 }}
                            >
                                <Spinner class="size-20 text-brand" />
                            </div>

                        {:else if error}
                            <div
                                    class="min-h-full flex items-center justify-center"
                                    in:fade={{ duration: 300 }}
                            >
                                <p class="error text-red-500">Error: {error}</p>
                            </div>

                        {:else}
                            <div class="min-h-0" in:fade={{ duration: 300 }}>
                                <!-- Your grid/list wrapper -->
                                <div class="content-grid min-h-0 pb-2">
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
                </div>
            </AspectRatio>
        </div>
    </div>
</div>

<style>
.content-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
        gap: 1.5rem;
        padding: 1.5rem;
        max-width: 100%;
        background-color: var(--color-canvas);
    }

    @media (min-width: 960px) {
        .content-grid {
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
