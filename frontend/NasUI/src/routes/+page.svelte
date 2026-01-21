<script lang="ts">
    import {DrivesList} from "$lib/components/ui/drive/index.ts"
    import {PoolsList} from "$lib/components/ui/pool/index.ts"
    import {Separator} from "$lib/components/ui/separator/index.js";
    import { Provider as SidebarProvider, Root as SidebarRoot, Content as SidebarContent } from "$lib/components/ui/sidebar/index.js";
    import {Button} from "$lib/components/ui/button/index.js";
    import PoolCreator from "$lib/components/ui/sidebar/sidebar-pool-creator.svelte";
    import { onMount } from 'svelte';
    import {Footer} from "$lib/components/ui/card/index.js";

    let myOpen = $state(false);
    let ratio = $state(2);


    function openSidebar() {
        myOpen = !myOpen;
        ratio = myOpen ? 1.5 : 2;
        if (myOpen) {
            document.documentElement.style.overflow = "hidden";
            document.body.style.overflow = "hidden";
        } else {
            document.documentElement.style.overflow = "";
            document.body.style.overflow = "";
        }
    }




</script>
<div
        class="min-h-screen overflow-x-hidden bg-[#121212] text-zinc-100"
        style="--sb: clamp(14rem, 30vw, 22rem);"
>
    <div
            class="grid min-h-screen transition-[grid-template-columns] duration-300 ease-in-out"
            style={myOpen
      ? "grid-template-columns: var(--sb) 1fr;"
      : "grid-template-columns: 0 1fr;"}
    >
        <!-- Sidebar -->
        <aside class="h-screen overflow-hidden bg-[#121212]">
            <SidebarProvider bind:open={myOpen} style="--sidebar-width: var(--sb);">
                <SidebarRoot>
                    <SidebarContent>
                        <PoolCreator />
                    </SidebarContent>
                </SidebarRoot>
            </SidebarProvider>
        </aside>

        <!-- Main content -->
        <main class="min-w-0 bg-[#121212]">
            <div class="h-screen overflow-y-auto overscroll-contain overflow-x-hidden">
                <div class="pt-10 pb-48 min-w-0">
                    <DrivesList ratio={ratio} />

                    {#if !myOpen}
                        <PoolsList />
                    {/if}

                    <Footer class="pt-10 pb-[env(safe-area-inset-bottom)]">
                        <div class="text-center text-sm text-zinc-500">
                            &copy; 2026 NasUI. All rights reserved.
                        </div>
                    </Footer>
                </div>
            </div>
        </main>
    </div>

    <Button onclick={openSidebar} class="fixed bottom-4 right-4 z-50">
        {myOpen ? 'Close Sidebar' : 'Open Sidebar'}
    </Button>
</div>
