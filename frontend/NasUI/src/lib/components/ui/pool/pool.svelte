<script lang="ts">
    import {Button} from "$lib/components/ui/button/index.ts";
    import {Status} from "$lib/components/ui/drive/index.ts";
    import type {Pool} from "$lib/models/pool.ts";
    import * as Card from "$lib/components/ui/card/index.js";

    interface Props {
        pool: Pool;
        id: string;
    }

    let { pool, id }: Props = $props();

    function formatBytes(bytes: number): string {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    function getUsagePercent(total: number, available: number): number {
        if (total === 0) return 0;
        return Math.round(((total - available) / total) * 100);
    }

    function formatDate(dateStr: string): string {
        if (!dateStr) return 'N/A';
        const date = new Date(dateStr);
        return date.toLocaleDateString();
    }

    let usagePercent = $derived(getUsagePercent(pool.totalCapacity, pool.availableCapacity));
    let driveCount = $derived(pool.AdoptedDrives ? Object.keys(pool.AdoptedDrives).length : 0);
</script>

<div class="pool-card" id="pool-{id}">
    <Card.Root class="h-full flex flex-col bg-zinc-700 text-zinc-100">
        <Card.Header class="pb-2">
            <Card.Title class="flex items-center justify-between gap-2">
                <span class="truncate">{pool.name || 'Unnamed Pool'}</span>
                <Status degraded={pool.status === 'degraded'} offline={pool.status === 'offline' || pool.status === 'failed'} />
            </Card.Title>
            <Card.Description class="text-zinc-400 text-xs">
                {pool.type || 'RAID'} • {pool.format || 'Unknown'}
            </Card.Description>
        </Card.Header>
        <Card.Content class="flex-1 flex flex-col gap-3 text-sm overflow-hidden">
            <!-- Capacity Bar -->
            <div class="space-y-1">
                <div class="flex justify-between text-xs">
                    <span class="text-muted-foreground">Capacity</span>
                    <span>{usagePercent}% used</span>
                </div>
                <div class="h-2 bg-zinc-600 rounded-full overflow-hidden">
                    <div
                        class="h-full transition-all duration-300 {usagePercent > 90 ? 'bg-red-500' : usagePercent > 70 ? 'bg-yellow-500' : 'bg-blue-500'}"
                        style="width: {usagePercent}%"
                    ></div>
                </div>
                <div class="flex justify-between text-xs text-zinc-400">
                    <span>{formatBytes(pool.totalCapacity - pool.availableCapacity)} used</span>
                    <span>{formatBytes(pool.totalCapacity)} total</span>
                </div>
            </div>

            <!-- Pool Info -->
            <div class="space-y-1.5">
                <div class="flex justify-between">
                    <span class="text-muted-foreground">Available</span>
                    <span>{formatBytes(pool.availableCapacity)}</span>
                </div>
                <div class="flex justify-between">
                    <span class="text-muted-foreground">Mount Point</span>
                    <span class="truncate max-w-[200px]">{pool.mountPoint || 'Not mounted'}</span>
                </div>
                <div class="flex justify-between">
                    <span class="text-muted-foreground">Created</span>
                    <span>{formatDate(pool.createdAt)}</span>
                </div>
            </div>

            <!-- Drives List -->
            {#if pool.AdoptedDrives && driveCount > 0}
                <div class="mt-auto min-h-0 flex flex-col">
                    <span class="text-xs text-muted-foreground mb-1 block shrink-0">Drives ({driveCount})</span>
                    <div class="flex flex-wrap gap-2 overflow-y-auto max-h-[80px]">
                        {#each Object.values(pool.AdoptedDrives).slice(0, 6) as poolDrive}
                            <span class="text-xs px-3 py-1.5 rounded bg-zinc-600 text-zinc-300 inline-flex items-center gap-2 max-w-[180px]">
                                <Status
                                    degraded={false}
                                    offline={poolDrive.drive.is_rotational}
                                />
                                <span class="truncate">{poolDrive.drive.name}</span>
                                <span class="text-green-500 shrink-0">({formatBytes(poolDrive.drive.size_bytes)})</span>
                            </span>
                        {/each}
                        {#if driveCount > 6}
                            <span class="text-xs px-3 py-1.5 rounded bg-zinc-600 text-zinc-400">
                                +{driveCount - 6} more
                            </span>
                        {/if}
                    </div>
                </div>
            {/if}
        </Card.Content>

        <Card.Footer class="pt-2 mt-auto shrink-0">
            <Button variant="outline" class="w-full text-zinc-300 border-zinc-500 hover:bg-zinc-600">
                Manage Pool
            </Button>
        </Card.Footer>
    </Card.Root>
</div>

<style>
    .pool-card {
        width: 90%;
        min-width: 600px;
        height: 380px;
        display: block;
    }
</style>

