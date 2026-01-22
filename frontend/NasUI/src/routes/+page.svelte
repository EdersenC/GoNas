<script lang="ts">
    import {DrivesList} from "$lib/components/ui/drive/index.ts"
    import {PoolsList} from "$lib/components/ui/pool/index.ts"
    import {Separator} from "$lib/components/ui/separator/index.js";
    import { Provider as SidebarProvider, Root as SidebarRoot, Content as SidebarContent } from "$lib/components/ui/sidebar/index.js";
    import {Button} from "$lib/components/ui/button/index.js";
    import PoolCreator from "$lib/components/ui/sidebar/sidebar-pool-creator.svelte";
    import { onMount } from 'svelte';
    import {Footer} from "$lib/components/ui/card/index.js";
    import {poolSelection} from "$lib/state/pool.svelte.js";

    let isSideBarOpened = $state(false);
    let poolCreatorMode = $derived(isSideBarOpened);
    let ratio = $state(2);


    function openSidebar() {
        isSideBarOpened = !isSideBarOpened;
        ratio = isSideBarOpened ? 1.5 : 2;
        if (!isSideBarOpened)
            poolSelection.clearSelectedDrives()
        if (isSideBarOpened) {
            document.documentElement.style.overflow = "hidden";
            document.body.style.overflow = "hidden";
        } else {
            document.documentElement.style.overflow = "";
            document.body.style.overflow = "";
        }
    }





</script>
<div
        class="min-h-screen overflow-x-hidden bg-canvas text-canvas-foreground"
        style="--sb: clamp(14rem, 30vw, 22rem);"
>
    <div
            class="grid min-h-screen transition-[grid-template-columns] duration-300 ease-in-out"
            style={isSideBarOpened
      ? "grid-template-columns: var(--sb) 1fr;"
      : "grid-template-columns: 0 1fr;"}
    >
        <!-- Sidebar -->
        <aside class="h-screen overflow-hidden bg-canvas">
            <SidebarProvider bind:open={isSideBarOpened} style="--sidebar-width: var(--sb);">
                <SidebarRoot>
                    <SidebarContent>
                        <PoolCreator poolSelection={poolSelection} />
                    </SidebarContent>
                </SidebarRoot>
            </SidebarProvider>
        </aside>

        <!-- Main content -->
        <main class="min-w-0 bg-canvas">
            <div class="h-screen overflow-y-auto overscroll-contain overflow-x-hidden">
                {poolSelection.getSelectedDrives().length}
                <div class="pt-10 pb-48 min-w-0">
                    <DrivesList ratio={ratio} poolCreatorMode={poolCreatorMode} poolSelection={poolSelection} />

                    {#if !isSideBarOpened}
                        <PoolsList />
                    {/if}

                    <Footer class="pt-10 pb-[env(safe-area-inset-bottom)]">
                        <div class="text-center text-sm text-muted-foreground">
                            &copy; 2026 NasUI. All rights reserved.
                        </div>
                    </Footer>
                </div>
            </div>
        </main>
    </div>

    <Button onclick={openSidebar} class="fixed bottom-4 right-4 z-50">
        {isSideBarOpened ? 'Close Sidebar' : 'Open Sidebar'}
    </Button>
</div>
