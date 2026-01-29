<script lang="ts">
    import {Button} from "$lib/components/ui/button/index.ts";
    import {Spinner} from "$lib/components/ui/spinner/index.js";
    import type {Pool} from "$lib/models/pool.ts";
    import {getPoolManagerContext} from "$lib/state/poolManager.svelte.js";

    interface Props {
		pool: Pool;
		onSuccess?: () => void;
	}

	let { pool, onSuccess }: Props = $props();

	const poolManager = getPoolManagerContext();
	let buildPending = $state(false);
	let deletePending = $state(false);
	let actionError = $state<string | null>(null);

	let canBuild = $derived(!pool.mountPoint || pool.mountPoint.trim().length === 0);

	export function resetState() {
		actionError = null;
		buildPending = false;
		deletePending = false;
	}

	const handleBuild = async () => {
		if (buildPending || deletePending || !canBuild) return;
		actionError = null;
		buildPending = true;
		try {
			await poolManager.buildPool(pool.uuid);
			onSuccess?.();
		} catch (err) {
			actionError = err instanceof Error ? err.message : "Failed to build pool.";
		} finally {
			buildPending = false;
		}
	};

	const handleDelete = async () => {
		if (deletePending || buildPending) return;
		if (
			typeof window !== "undefined" &&
			!window.confirm(`Delete ${pool.name || "this pool"}? This cannot be undone.`)
		) {
			return;
		}
		actionError = null;
		deletePending = true;
		try {
			await poolManager.deletePool(pool.uuid);
			onSuccess?.();
		} catch (err) {
			actionError = err instanceof Error ? err.message : "Failed to delete pool.";
		} finally {
			deletePending = false;
		}
	};
</script>

<div class="flex flex-col gap-4">
	{#if canBuild}
		<div class="rounded-xl border border-brand/40 bg-surface-muted/40 px-4 py-3 text-sm">
			<div class="font-semibold">Build needed</div>
			<p class="text-xs text-panel-foreground/70">
				This pool does not have a mount point. Build it to create one.
			</p>
		</div>
		<Button
			variant="green"
			class="w-full shadow-sm"
			onclick={handleBuild}
			disabled={buildPending || deletePending}
		>
			{#if buildPending}
				<Spinner class="size-4" />
				Building...
			{:else}
				Build Pool
			{/if}
		</Button>
	{:else}
		<div class="rounded-xl border border-success/30 bg-success/10 px-4 py-3 text-sm text-success">
			Pool is mounted and ready.
		</div>
	{/if}

	<div class="grid gap-2 rounded-2xl border border-brand/40 bg-surface/60 p-4 text-sm">
		<span class="text-xs font-semibold uppercase tracking-wider text-panel-foreground/70">
			Danger zone
		</span>
		<p class="text-xs text-panel-foreground/70">
			Deleting this pool removes the storage pool and its configuration.
		</p>
		<Button
			variant="destructive"
			class="w-full shadow-sm"
			onclick={handleDelete}
			disabled={deletePending || buildPending}
		>
			{#if deletePending}
				<Spinner class="size-4" />
				Deleting...
			{:else}
				Delete Pool
			{/if}
		</Button>
	</div>

	{#if actionError}
		<div class="rounded-xl border border-destructive/50 bg-destructive/10 p-3 text-xs text-destructive">
			{actionError}
		</div>
	{/if}
</div>
