<script lang="ts">
    import {Button} from "$lib/components/ui/button/index.ts";
    import {Status, Size, ActionDropdown} from "$lib/components/ui/drive/index.ts";
    import type {Drive} from "$lib/models/drive.ts";
    import * as Card from "$lib/components/ui/card/index.js";
    import { page } from '$app/stores';

    let{
        drive,
        id
    } = $props();

    function formatBytes(bytes: number): string {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

   export async function Post(driveId: string) {
        let url = `http://localhost:8080/api/v1/drives/adopt/${driveId}`;
        if ($page.url.pathname.startsWith('/pools')) {
            url = `http://localhost:8080/api/v1/pools/adopt/${driveId}`;
        }
        const res = await fetch(url, {
            method: 'POST',
        });
        console.log(await res.json())
        console.log(`Adopting drive with ID: ${driveId}`);
    }

    let used = $derived(drive.size_bytes - drive.fsavail);
    let percent = $derived(drive.size_bytes > 0 ? (used / drive.size_bytes) * 100 : 0);
</script>

<div class="drive-card" id="drive-{id}">
    <Card.Root class="h-full flex flex-col bg-zinc-700 text-zinc-100">
        <Card.Header>
            <Card.Title class="flex items-center gap-2">
                {#if drive.name.includes("loop")}
                    <Status degraded={false} offline={false}></Status>
                {:else}
                    <Status degraded={false} offline={drive.is_rotational}></Status>
                {/if}
                <span class="flex-1 text-center font-bold truncate">{drive.name}</span>
                <ActionDropdown />
            </Card.Title>
        </Card.Header>
        <Card.Content class="flex-1 flex flex-col gap-1 text-sm">
            <div class="flex justify-center mb-1">
                <div class="flex flex-col gap-0 text-center">
                    <span class="text-zinc-100 font-bold underline text-xs">Size</span>
                    <span class="font-bold text-sm truncate max-w-[100px]">{formatBytes(drive.size_bytes)}</span>
                </div>
            </div>
            <div class="flex items-center">
                <div class="flex-1 text-left">
                    <div class="flex flex-col gap-1">
                        <span class="text-zinc-100 font-bold underline text-xs">Used</span>
                        <span class="font-bold truncate max-w-[80px]">{formatBytes(used)}</span>
                    </div>
                </div>
                <div class="flex justify-center">
                    <Size percent={percent} />
                </div>
                <div class="flex-1 text-center">
                    <div class="flex flex-col gap-1">
                        <span class="text-zinc-100 font-bold underline text-xs">Available</span>
                        <span class="font-bold truncate max-w-[80px]">{formatBytes(drive.fsavail)}</span>
                    </div>
                </div>
            </div>
        </Card.Content>
        <Card.Footer>
            <Button variant="green" class="w-full" onclick={()=>Post(drive.drive_key.kind+":"+drive.drive_key.value)}>Adopt{drive.name} </Button>
        </Card.Footer>
    </Card.Root>
</div>

<style>
    .drive-card {
        width: 280px;
        height: 360px;
        display: block;
    }
</style>