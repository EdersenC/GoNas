<script lang="ts">
    import type {Pool} from "$lib/models/pool.ts";

    interface Props {
		pool: Pool;
	}

	let { pool }: Props = $props();

	let usagePercent = $derived(getUsagePercent(pool.totalCapacity, pool.availableCapacity));

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

<div class="flex flex-col gap-4">
	<div class="grid gap-3 md:grid-cols-2">
		<div class="rounded-2xl border border-brand/40 bg-surface/60 p-4 text-sm">
			<span class="text-xs font-semibold uppercase tracking-wider text-panel-foreground/70">
				Health
			</span>
			<div class="mt-2 flex items-center justify-between">
				<span class="capitalize">{pool.status || "unknown"}</span>
				<span class="text-xs text-panel-foreground/70">
					{pool.type.Level || "RAID"} • {pool.format || "Unknown"}
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
</div>
