<script lang="ts">
    import {DrivesList} from "$lib/components/ui/drive/index.ts"
    import {PoolsList} from "$lib/components/ui/pool/index.ts"
    import {Separator} from "$lib/components/ui/separator/index.js";
    import { Provider as SidebarProvider, Root as SidebarRoot, Content as SidebarContent } from "$lib/components/ui/sidebar/index.js";
    import {Button} from "$lib/components/ui/button/index.js";
    import PoolCreator from "$lib/components/ui/sidebar/sidebar-pool-creator.svelte";
    import { onMount } from 'svelte';
    import {Footer} from "$lib/components/ui/card/index.js";
    import {getDriveManagerContext} from "$lib/state/driveManager.svelte.js";
    import {getPoolManagerContext} from "$lib/state/poolManager.svelte.js";
    import type {AdoptedDrive} from "$lib/models/drive.js";

    let isSideBarOpened = $state(false);
    let poolCreatorMode = $derived(isSideBarOpened);
    let poolManager = getPoolManagerContext()
    let driveManager = getDriveManagerContext()
    let ratio = $state(2);

    function openSidebar() {
        isSideBarOpened = !isSideBarOpened;
        ratio = isSideBarOpened ? 1.2 : 2;
        if (isSideBarOpened) {
            document.documentElement.style.overflow = "hidden";
            document.body.style.overflow = "hidden";
        } else {
            document.documentElement.style.overflow = "";
            document.body.style.overflow = "";
        }
    }



onMount(async () => {
    await Promise.all([
        poolManager.fetchPools(),
        driveManager.fetchAdoptedDrives(),
    ]);

    loadFakeDrives(100)

});

    function loadFakeDrives(count: number){
        let fakeDrive: Record<string, AdoptedDrive> = {};
        for (let i = 1; i <= count; i++) {

            const size = 500 * 1024 * 1024 * 1024 * (i * Math.random());
            const id = `fake-${i}`;
            fakeDrive[id] = {
                uuid: `fake-uuid-${i}`,
                drive :{
                    name: `FakeDrive${i}`,
                    uuid: `fake-uuid-${i}`,
                    drive_key: { kind: 'by-path', value: `/dev/sd${i}` },
                    by_ids: [`/dev/disk/by-id/fake-${i}`],
                    wwid: `wwid-${i}`,
                    path: `/dev/sd${i}`,
                    size_sectors: 976773168,
                    logical_block_size: 512,
                    physical_block_size: 4096,
                    size_bytes: size,
                    is_rotational: true,
                    model: `FAKE_MODEL_${i}`,
                    vendor: `FAKE_VENDOR`,
                    serial: `SNFAKE${1000 + i}`,
                    type: (i % 2 === 0 ? 'SSD' : 'HDD'),
                    mountpoint: '',
                    partitions: [],
                    fstype: '',
                    fsavail: size * Math.random(),
                }
            };
        }
        driveManager.addAdoptedDrives(fakeDrive)
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
                        <PoolCreator />
                    </SidebarContent>
                </SidebarRoot>
            </SidebarProvider>
        </aside>

        <!-- Main content -->
        <main class="min-w-0 bg-canvas">
            <div class="h-screen overflow-y-auto overscroll-contain overflow-x-hidden">
                <div class="pt-10 pb-48 min-w-0">
                    <DrivesList ratio={ratio} poolCreatorMode={poolCreatorMode}  />


                    {#if !isSideBarOpened}
                        <PoolsList />
                    {/if}

                    <Footer class="pt-10 pb-[env(safe-area-inset-bottom)]">
                        <div class="text-center text-sm text-muted-foreground">
                            &copy; 2026 GONAS. All rights reserved.
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
