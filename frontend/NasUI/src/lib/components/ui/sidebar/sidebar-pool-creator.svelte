<script lang="ts">
    import {getDriveManagerContext, DriveManager} from "$lib/state/driveManager.svelte.js";
    import {getPoolManagerContext, type PoolManager} from "$lib/state/poolManager.svelte.js";
    import {Button} from "$lib/components/ui/button/index.js";
    import {json} from "@sveltejs/kit";

    let name = $state('');
    let raidLevel = $state(10);
    let format = $state('ext4');
    let build = $state(false);

    let driveManager: DriveManager = getDriveManagerContext()
    let poolManager: PoolManager = getPoolManagerContext()

    const raidOptions = [0, 1, 5, 10];
    const formatOptions = ["ext4", "xfs", "btrfs"];

    let selectedDriveItems = $derived(
        Object.values(driveManager.adoptedDrives).filter((adoptedDrive) =>
            driveManager.selectedDrives.includes(adoptedDrive.drive.uuid)
        )
    );
    let selectedDriveCount = $derived(driveManager.selectedDrives.length);
    let selectedDriveTotalBytes = $derived(
        selectedDriveItems.reduce((total, adoptedDrive) => total + (adoptedDrive?.drive?.size_bytes ?? 0), 0)
    );

    function formatBytes(bytes: number): string {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        const index = Math.min(i, sizes.length - 1);
        return parseFloat((bytes / Math.pow(k, index)).toFixed(2)) + ' ' + sizes[index];
    }

    async function createPool() {
        if (driveManager.selectedDrives.length === 0) {
            console.warn('No drives selected for pool creation');
            return;
        }
        console.log('Selected drives for pool creation', driveManager.selectedDrives);
        const payload = {
            name,
            raidLevel,
            drives: driveManager.getSelectedDrives(),
            format,
            build,
        };
        console.log('Creating pool with payload', payload);
        try {
            await poolManager.postPool(JSON.stringify(payload));
            driveManager.removeSelectedDrivesFromAdopted();
            name = '';
        } catch (e) {
            console.error('Failed to create pool', e);
        }
    }
</script>

<div class="flex h-full flex-col gap-4 bg-canvas p-5 text-canvas-foreground">
    <div class="relative overflow-hidden rounded-2xl border border-brand/50 bg-surface/80 p-5 shadow-sm">
        <div class="absolute inset-0 opacity-60" style="background: radial-gradient(circle at top, rgba(59, 130, 246, 0.22), transparent 60%);"></div>
        <div class="relative flex items-start justify-between gap-3">
            <div>
                <p class="text-xs font-semibold uppercase tracking-[0.25em] text-muted-foreground">Pool Creator</p>
                <h3 class="text-2xl font-semibold tracking-tight">Build a new storage pool</h3>
                <p class="text-sm text-muted-foreground">Pick drives from the list and dial in your layout.</p>
            </div>
            <div class="shrink-0 rounded-full border border-brand/50 bg-surface-muted/40 px-3 py-1 text-xs font-semibold">
                {selectedDriveCount} selected
            </div>
        </div>
    </div>

    <div class="grid gap-3">
        <div class="grid gap-2">
            <label class="text-sm font-semibold uppercase tracking-wider text-muted-foreground">Pool name</label>
            <input
                class="w-full rounded-lg border border-surface-border/70 bg-surface-muted/60 px-3 py-2 text-sm text-surface-foreground outline-none transition focus:border-brand/60 focus:ring-2 focus:ring-brand/30"
                bind:value={name}
                placeholder="Fast storage pool"
            />
        </div>

        <div class="grid gap-2">
            <div class="flex items-center justify-between">
                <label class="text-sm font-semibold uppercase tracking-wider text-muted-foreground">RAID level</label>
                <span class="text-sm text-muted-foreground">Current: RAID {raidLevel}</span>
            </div>
            <div class="grid grid-cols-4 gap-2">
                {#each raidOptions as level}
                    <button
                        type="button"
                        class={`rounded-lg border px-2 py-2 text-xs font-semibold transition ${
                            raidLevel === level
                                ? "border-brand/70 bg-brand/20 text-brand-foreground"
                                : "border-surface-border/60 bg-surface-muted/40 text-surface-foreground hover:bg-surface-muted/70"
                        }`}
                        onclick={() => (raidLevel = level)}
                    >
                        RAID {level}
                    </button>
                {/each}
            </div>
        </div>

        <div class="grid gap-2">
            <label class="text-sm font-semibold uppercase tracking-wider text-muted-foreground">Format</label>
            <select
                class="w-full rounded-lg border border-surface-border/70 bg-surface-muted/60 px-3 py-2 text-sm text-surface-foreground outline-none transition focus:border-brand/60 focus:ring-2 focus:ring-brand/30"
                bind:value={format}
            >
                {#each formatOptions as option}
                    <option value={option}>{option}</option>
                {/each}
            </select>
        </div>

        <label class="flex items-center justify-between gap-3 rounded-xl border border-brand/40 bg-surface-muted/40 px-3 py-3">
            <div class="space-y-1">
                <div class="text-base font-semibold">Build immediately</div>
                <div class="text-sm text-muted-foreground">Start the build process right after creation.</div>
            </div>
            <span class="relative inline-flex h-6 w-11 shrink-0 items-center">
                <input type="checkbox" class="peer sr-only" bind:checked={build} />
                <span class="absolute inset-0 rounded-full bg-surface-border/60 transition peer-checked:bg-success/70"></span>
                <span class="absolute left-1 top-1 h-4 w-4 rounded-full bg-surface-foreground transition peer-checked:translate-x-5"></span>
            </span>
        </label>
    </div>

    <div class="grid gap-2 rounded-2xl border border-brand/40 bg-surface/60 p-4">
        <div class="flex items-center justify-between">
            <span class="text-sm font-semibold uppercase tracking-wider text-muted-foreground">Selected drives</span>
            <span class="text-sm text-muted-foreground">
                {selectedDriveCount} drives / {formatBytes(selectedDriveTotalBytes)}
            </span>
        </div>

        {#if selectedDriveItems.length === 0}
            <div class="rounded-lg border border-dashed border-brand/40 bg-surface-muted/30 px-3 py-4 text-sm text-muted-foreground">
                Select drives from the list to see them here.
            </div>
        {:else}
            <div class="grid max-h-56 gap-2 overflow-y-auto pr-1">
                {#each selectedDriveItems as adoptedDrive}
                    <div class="flex items-center justify-between gap-3 rounded-lg border border-brand/30 bg-panel/60 px-3 py-3">
                        <div class="min-w-0">
                            <div class="truncate text-base font-semibold">{adoptedDrive.drive.name}</div>
                            <div class="truncate text-sm text-muted-foreground">
                                {adoptedDrive.drive.model || "Unknown model"} - {adoptedDrive.drive.is_rotational ? "HDD" : "SSD"}
                            </div>
                        </div>
                        <div class="shrink-0 text-sm text-muted-foreground">{formatBytes(adoptedDrive.drive.size_bytes || 0)}</div>
                        <button
                            type="button"
                            class="shrink-0 rounded-md border border-brand/40 px-2 py-1 text-[11px] font-semibold uppercase tracking-wide text-muted-foreground transition hover:text-surface-foreground"
                            onclick={() => driveManager.toggleSelectedDrive(adoptedDrive.drive.uuid)}
                        >
                            Remove
                        </button>
                    </div>
                {/each}
            </div>
        {/if}
    </div>

    <div class="mt-auto flex gap-2">
        <Button
            class="flex-1 shadow-sm"
            variant="green"
            disabled={driveManager.selectedDrives.length === 0}
            onclick={createPool}
        >
            Create pool
        </Button>
        <Button class="flex-1" variant="secondary" onclick={driveManager.clearSelectedDrives}>
            Clear selection
        </Button>
    </div>
</div>
