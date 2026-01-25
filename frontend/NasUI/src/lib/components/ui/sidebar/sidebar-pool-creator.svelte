<script lang="ts">
    import {getDriveManagerContext, DriveManager} from "$lib/state/driveManager.svelte.js";
    import {getPoolManagerContext, type PoolManager} from "$lib/state/poolManager.svelte.js";

    let name = $state('');
    let raidLevel = $state(10);
    let format = $state('ext4');
    let build = $state(false);

    let driveManager: DriveManager = getDriveManagerContext()
    let poolManager: PoolManager = getPoolManagerContext()


    async function createPool() {
        if (driveManager.selectedDrives.length === 0) {
            console.warn('No drives selected for pool creation');
            return;
        }
        console.log('Selected drives for pool creation', driveManager.selectedDrives);
        const payload = {
            name,
            raidLevel,
            drives: driveManager.selectedDrives,
            format,
            build,
        };
        console.log('Creating pool with payload', payload);
        try {
            await poolManager.postPool(payload);
            driveManager.removeSelectedDrivesFromAdopted();
            name = '';
        } catch (e) {
            console.error('Failed to create pool', e);
        }
    }
</script>

<div class="p-4 bg-surface text-surface-foreground min-h-full">
    <h3 class="text-sm font-semibold mb-2">Create Pool</h3>
    <div class="flex flex-col gap-2">
        <label class="text-xs">Name</label>
        <input class="p-2 rounded bg-surface-muted text-surface-foreground border border-surface-border" bind:value={name} placeholder="Pool name" />

        <label class="text-xs">RAID Level</label>
        <input type="number" class="p-2 rounded bg-surface-muted text-surface-foreground border border-surface-border" bind:value={raidLevel} min={0} />

        <label class="text-xs">Format</label>
        <select class="p-2 rounded bg-surface-muted text-surface-foreground border border-surface-border" bind:value={format}>
            <option value="ext4">ext4</option>
            <option value="xfs">xfs</option>
            <option value="btrfs">btrfs</option>
        </select>

        <label class="flex items-center gap-2"><input type="checkbox" bind:checked={build} /> Build</label>

        <div class="text-xs text-muted-foreground">Selected drives: {driveManager.selectedDrives.length}</div>
        {#if true}
            <ul class="text-xs list-disc ml-4 text-muted-foreground">
                {#each driveManager.selectedDrives as d}
                    <li>{d}</li>
                {/each}
            </ul>
        {/if}

        <div class="flex gap-2">
            <button class="px-3 py-2 rounded bg-success text-success-foreground hover:bg-success/90"
                    onclick={createPool}>Create</button>
            <button class="px-3 py-2 rounded bg-surface-muted text-surface-foreground hover:bg-surface-muted/80"
                    onclick={driveManager.clearSelectedDrives}>Clear</button>
        </div>
    </div>
</div>
