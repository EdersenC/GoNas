<script lang="ts">
    import {Button} from "$lib/components/ui/button/index.ts";
    import {Status} from "$lib/components/ui/drive/index.ts";
    import type {Drive} from "$lib/models/drive.ts";
    import * as Card from "$lib/components/ui/card/index.js";
    export let drive: Drive;
    export let id : string;


    function formatBytes(bytes: number): string {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

   export async function Post(driveId: string) {
        // Placeholder function for adopting a drive
        const res = await fetch(`http://localhost:8080/api/v1/drives/adopt/${driveId}`, {
            method: 'POST',
        });
        console.log(await res.json())
        console.log(`Adopting drive with ID: ${driveId}`);
    }
</script>

<div class="drive-card" id="drive-{id}">
    <Card.Root class="h-full flex flex-col bg-zinc-900 text-zinc-100
">
        <Card.Header>
            <Card.Title>{drive.name}
                {#if drive.name.includes("loop")}
                    <Status degraded = {false} offline ={false}></Status>
                {:else}
                    <Status degraded = {false} offline ={drive.is_rotational}></Status>
                {/if}
            </Card.Title>
        </Card.Header>
        <Card.Content class="flex-1 flex flex-col gap-2 text-sm">
            <div class="flex justify-between">
                <span class="text-muted-foreground">Type</span>
                <span>{drive.type}</span>
            </div>
            <div class="flex justify-between">
                <span class="text-muted-foreground">Size</span>
                <span>{formatBytes(drive.size_bytes)}</span>
            </div>
            <div class="flex justify-between">
                <span class="text-muted-foreground">Available</span>
                <span>{formatBytes(drive.fsavail)}</span>
            </div>
            <div class="flex justify-between">
                <span class="text-muted-foreground">Model</span>
                <span class="truncate max-w-[120px]">{drive.model || 'N/A'}</span>
            </div>
            <div class="flex justify-between">
                <span class="text-muted-foreground">Mount Point</span>
                <span class="truncate max-w-[120px]">{drive.mountpoint || 'Not mounted'}</span>
            </div>
            <div class="flex justify-between">
                <span class="text-muted-foreground">Rotational</span>
                <span>{drive.is_rotational ? 'HDD' : 'SSD'}</span>
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
