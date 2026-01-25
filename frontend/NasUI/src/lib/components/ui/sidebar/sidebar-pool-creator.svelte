<script lang="ts">
    import {getDriveManagerContext, DriveManager} from "$lib/state/pool.svelte.js";

    let name = $state('');
    let raidLevel = $state(10);
    let format = $state('ext4');
    let build = $state(false);

    let driveManager:DriveManager = getDriveManagerContext()


    async function createPool() {
        if (driveManager.selectedDrives.length === 0) {
            console.warn('No drives selected for pool creation');
            return;
        }
        console.log('Selected drives for pool creation', driveManager.getSelectedDrives());
        const payload = {
            name,
            raidLevel,
            drives: driveManager.getSelectedDrives(),
            format,
            build,
        };

        console.log('Creating pool with payload', payload);

        // Example POST - adjust URL to your API
        try {
            const res = await fetch('http://localhost:8080/api/v1/pool', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload),
            });
            const data = await res.json();
            console.log('Pool creation response', data);
            // Clear form and selection on success
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

        <div class="text-xs text-muted-foreground">Selected drives: {driveManager.getSelectedDrives().length}</div>
        {#if true}
            <ul class="text-xs list-disc ml-4 text-muted-foreground">
                {#each driveManager.getSelectedDrives() as d}
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
