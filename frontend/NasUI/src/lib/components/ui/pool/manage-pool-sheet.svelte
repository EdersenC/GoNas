<script lang="ts">
	import { Button } from "$lib/components/ui/button/index.ts";
	import * as Sheet from "$lib/components/ui/sheet/index.js";
	import * as Tabs from "$lib/components/ui/tabs/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import type { Pool } from "$lib/models/pool.ts";
	import { getPoolManagerContext } from "$lib/state/poolManager.svelte.js";

	interface Props {
		pool: Pool;
		triggerClass?: string;
	}

	let { pool, triggerClass = "" }: Props = $props();

	const poolManager = getPoolManagerContext();
	let manageOpen = $state(false);
	let buildPending = $state(false);
	let deletePending = $state(false);
	let actionError = $state<string | null>(null);

	let canBuild = $derived(!pool.mountPoint || pool.mountPoint.trim().length === 0);
	let usagePercent = $derived(getUsagePercent(pool.totalCapacity, pool.availableCapacity));

	const resetActionState = () => {
		actionError = null;
		buildPending = false;
		deletePending = false;
	};

	const handleBuild = async () => {
		if (buildPending || deletePending || !canBuild) return;
		actionError = null;
		buildPending = true;
		try {
			await poolManager.buildPool(pool.uuid);
			manageOpen = false;
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
			manageOpen = false;
		} catch (err) {
			actionError = err instanceof Error ? err.message : "Failed to delete pool.";
		} finally {
			deletePending = false;
		}
	};

	function formatBytes(bytes: number): string {
		if (bytes === 0) return "0 B";
		const k = 1024;
		const sizes = ["B", "KB", "MB", "GB", "TB"];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
	}

	function getUsagePercent(total: number, available: number): number {
		if (total === 0) return 0;
		return Math.round(((total - available) / total) * 100);
	}

	function formatDate(dateStr: string): string {
		if (!dateStr) return "N/A";
		const date = new Date(dateStr);
		return date.toLocaleDateString();
	}
</script>

<Sheet.Root bind:open={manageOpen}>
	<Sheet.Trigger asChild>
		<Button
			variant="outline"
			class={`w-full text-surface-foreground border-surface-border hover:bg-surface-muted ${triggerClass}`}
			onclick={resetActionState}
		>
			Manage Pool
		</Button>
	</Sheet.Trigger>
	<Sheet.Content side="right" class="bg-canvas text-canvas-foreground w-full p-0 sm:max-w-lg">
		<Sheet.Header class="sr-only">
			<Sheet.Title>Manage {pool.name || "Pool"}</Sheet.Title>
			<Sheet.Description>Configure storage pool actions and settings.</Sheet.Description>
		</Sheet.Header>
		<div class="flex h-full flex-col gap-4 bg-canvas p-5 text-canvas-foreground">
			<div class="relative overflow-hidden rounded-2xl border border-brand/50 bg-surface/80 p-5 shadow-sm">
				<div
					class="absolute inset-0 opacity-60"
					style="background: radial-gradient(circle at top, rgba(59, 130, 246, 0.22), transparent 60%);"
				></div>
				<div class="relative">
					<p class="text-xs font-semibold uppercase tracking-[0.25em] text-panel-foreground/70">
						Pool Manager
					</p>
					<h3 class="text-2xl font-semibold tracking-tight text-panel-foreground">
						Manage {pool.name || "Pool"}
					</h3>
					<p class="text-sm text-panel-foreground/70">Configure storage pool actions and settings.</p>
				</div>
			</div>

			<Tabs.Root value="stats" class="flex flex-1 flex-col gap-4">
				<Tabs.List>
					<Tabs.Trigger value="stats">Stats</Tabs.Trigger>
					<Tabs.Trigger value="shares">Shares</Tabs.Trigger>
					<Tabs.Trigger value="actions">Actions</Tabs.Trigger>
				</Tabs.List>

				<Tabs.Content value="stats" class="flex flex-col gap-4">
					<div class="grid gap-3 md:grid-cols-2">
						<div class="rounded-2xl border border-brand/40 bg-surface/60 p-4 text-sm">
							<span class="text-xs font-semibold uppercase tracking-wider text-panel-foreground/70">
								Health
							</span>
							<div class="mt-2 flex items-center justify-between">
								<span class="capitalize">{pool.status || "unknown"}</span>
								<span class="text-xs text-panel-foreground/70">
									{pool.type || "RAID"} • {pool.format || "Unknown"}
								</span>
							</div>
						</div>
						<div class="rounded-2xl border border-brand/40 bg-surface/60 p-4 text-sm">
							<span class="text-xs font-semibold uppercase tracking-wider text-panel-foreground/70">
								Mount point
							</span>
							<div class="mt-2 truncate">{pool.mountPoint || "Not mounted"}</div>
							<div class="text-xs text-panel-foreground/70">{formatDate(pool.createdAt)} created</div>
						</div>
					</div>

					<div class="rounded-2xl border border-brand/40 bg-surface/60 p-4 text-sm">
						<div class="flex items-center justify-between text-xs text-panel-foreground/70">
							<span class="font-semibold uppercase tracking-wider">Capacity</span>
							<span>{usagePercent}% used</span>
						</div>
						<div class="mt-2 h-2 rounded-full bg-surface-muted/80">
							<div
								class="h-full rounded-full transition-all duration-300 {usagePercent > 90
									? 'bg-danger'
									: usagePercent > 70
										? 'bg-warning'
										: 'bg-brand'}"
								style="width: {usagePercent}%"
							></div>
						</div>
						<div class="mt-2 flex items-center justify-between text-xs text-panel-foreground/70">
							<span>{formatBytes(pool.totalCapacity - pool.availableCapacity)} used</span>
							<span>{formatBytes(pool.availableCapacity)} free</span>
							<span>{formatBytes(pool.totalCapacity)} total</span>
						</div>
					</div>
				</Tabs.Content>

				<Tabs.Content value="shares" class="flex flex-col gap-4">
					<div class="rounded-2xl border border-brand/40 bg-surface/60 p-4 text-sm">
						<span class="text-xs font-semibold uppercase tracking-wider text-panel-foreground/70">Shares</span>
						<p class="mt-2 text-sm text-panel-foreground/70">
							Create shares within this pool for VM access, permissions, and quotas.
						</p>
						<div class="mt-3 rounded-lg border border-dashed border-brand/40 bg-surface-muted/30 px-3 py-4 text-xs text-panel-foreground/70">
							Share management UI coming next.
						</div>
					</div>
				</Tabs.Content>

				<Tabs.Content value="actions" class="flex flex-col gap-4">
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
				</Tabs.Content>
			</Tabs.Root>
		</div>
	</Sheet.Content>
</Sheet.Root>
