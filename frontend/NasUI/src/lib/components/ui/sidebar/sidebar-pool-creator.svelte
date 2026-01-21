<script lang="ts">
    import { selectedDrives, selectedDrivesActions } from "$lib/stores/selectedDrives.ts";
    import { onDestroy } from 'svelte';

    let name = $state('');
    let raidLevel = $state(10);
    let format = $state('ext4');
    let build = $state(false);

    let selected: string[] = [];
    const unsubscribe = selectedDrives.subscribe(v => selected = v);
    onDestroy(unsubscribe);

    async function createPool() {
        const payload = {
            name,
            raidLevel,
            drives: selected,
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
            selectedDrivesActions.clear();
        } catch (e) {
            console.error('Failed to create pool', e);
        }
    }
</script>

<div class="p-4 bg-zinc-800 text-zinc-100 min-h-full">
    <h3 class="text-sm font-semibold mb-2">Create Pool</h3>
    <div class="flex flex-col gap-2">
        <label class="text-xs">Name</label>
        <input class="p-2 rounded bg-zinc-900 text-white border border-zinc-700" bind:value={name} placeholder="Pool name" />

        <label class="text-xs">RAID Level</label>
        <input type="number" class="p-2 rounded bg-zinc-900 text-white border border-zinc-700" bind:value={raidLevel} min={0} />

        <label class="text-xs">Format</label>
        <select class="p-2 rounded bg-zinc-900 text-white border border-zinc-700" bind:value={format}>
            <option value="ext4">ext4</option>
            <option value="xfs">xfs</option>
            <option value="btrfs">btrfs</option>
        </select>

        <label class="flex items-center gap-2"><input type="checkbox" bind:checked={build} /> Build</label>

        <div class="text-xs text-zinc-400">Selected drives: {selected.length}</div>
        {#if selected.length}
            <ul class="text-xs list-disc ml-4 text-zinc-300">
                {#each selected as d}
                    <li>{d}</li>
                {/each}
            </ul>
        {/if}

        <div class="flex gap-2">
            <button class="px-3 py-2 rounded bg-green-600" on:click={createPool}>Create</button>
            <button class="px-3 py-2 rounded bg-zinc-700" on:click={() => selectedDrivesActions.clear()}>Clear</button>
        </div>
    </div>
</div>
